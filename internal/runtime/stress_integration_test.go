//go:build cgo && integration

package runtime_test

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/inference/llama"
	"github.com/VishalPainjane/TinyBrain-OS/internal/kv"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
	"github.com/VishalPainjane/TinyBrain-OS/internal/swap"
)

func TestRuntime_SwapBoundaryStress(t *testing.T) {
	t.Parallel()

	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	modelID := "stress-model"
	reg := registry.NewModelRegistry()
	if err := reg.RegisterModel(registry.ModelDefinition{ID: modelID, Path: path}); err != nil {
		t.Fatalf("RegisterModel() error = %v", err)
	}

	resolver := runtime.NewRegistryResolver(reg)
	cfg := llama.DefaultConfig()
	cfg.GreedySampler = true
	provider := llama.NewLlamaProvider(resolver, cfg)

	if err := provider.LoadModel(modelID); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	defer func() { _ = provider.UnloadModel(modelID) }()

	bus := events.NewChannelBus(100)
	kvm := kv.NewStubManager(bus)
	table := process.NewProcessTable()

	// Artificial VRAM cap of 8KB to force extreme memory pressure
	totalVRAM := uint64(8192)
	mgr := swap.NewStubManager(table, kvm, bus, totalVRAM)

	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(table, queue)

	numAgents := 10
	for i := 0; i < numAgents; i++ {
		pid := "agent-" + string(rune('0'+i))
		p := process.Process{
			PID:       pid,
			State:     process.Ready,
			KVCacheID: "kv-" + pid,
			Priority:  8, // Start in Q0
		}
		if err := table.Create(p); err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if err := kvm.Allocate(p.KVCacheID, p.PID, 4096); err != nil { // 4KB per process means only 2 fit in VRAM
			t.Fatalf("Allocate() error = %v", err)
		}
		if err := sched.Enqueue(p); err != nil {
			t.Fatalf("Enqueue() error = %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(numAgents)

	// Background worker simulating the event loop / runtime scheduler
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				p, err := sched.Schedule()
				if err != nil {
					time.Sleep(10 * time.Millisecond)
					continue
				}

				// Simulate inference at the CGO boundary
				_, _ = provider.Generate(runtime.GenerateRequest{
					ModelID: modelID,
					Prompt:  "Hello",
				})

				_ = sched.RecordToken(p.PID)

				// Re-fetch process state
				p, _ = table.Get(p.PID)

				// Artificially push LastExecution back so ShouldSwap thinks it's idle,
				// forcing the adaptive heuristic to trigger under our low TotalVRAM
				p.LastExecution = time.Now().Add(-100 * time.Second)
				_ = table.UpdateState(p.PID, process.Waiting)

				used := kvm.VRAMUsage()
				if scheduler.ShouldSwap(p, time.Now(), used, totalVRAM) {
					_ = mgr.SwapOut(p.PID)
				}
			}
		}
	}()

	// Simulate agents polling and swapping back in
	for i := 0; i < numAgents; i++ {
		go func(agentID int) {
			defer wg.Done()
			pid := "agent-" + string(rune('0'+agentID))

			for step := 0; step < 2; step++ {
				select {
				case <-ctx.Done():
					return
				case <-time.After(50 * time.Millisecond):
					p, _ := table.Get(pid)
					if p.State == process.Hibernated {
						_ = mgr.SwapIn(pid)
						_ = sched.Enqueue(p)
					}
				}
			}
		}(i)
	}

	wg.Wait()
}
