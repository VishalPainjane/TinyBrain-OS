package main

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/agents"
	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
	"github.com/VishalPainjane/TinyBrain-OS/internal/inference/llama"
	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/router"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

func runBenchmark(mode, prompt string, stdout, stderr io.Writer) int {
	targetAgent := "bench-monolith"
	if mode == "swarm" {
		targetAgent = "bench-swarm-worker"
	}

	profile, err := hardware.ProbeAndClassify()
	if err != nil {
		fmt.Fprintf(stderr, "probe: %v\n", err)
		return 1
	}

	reg, err := openRegistry()
	if err != nil {
		fmt.Fprintf(stderr, "open registry: %v\n", err)
		return 1
	}
	defer reg.Close()

	// Use test models from root
	if err := registry.LoadModelsYAML(filepath.Join("testdata", "models.yaml"), reg); err != nil {
		fmt.Fprintf(stderr, "load models: %v\n", err)
	}

	areg := registry.NewAgentRegistry()
	if err := registry.LoadAgentsYAML(filepath.Join("testdata", "fleet.yaml"), areg); err != nil {
		fmt.Fprintf(stderr, "load fleet: %v\n", err)
		return 1
	}

	resolver := runtime.NewRegistryResolver(reg)
	cfg := llama.ConfigFromProbe(profile.Probe)
	if ng, ok := ngLayersFromEnv(); ok {
		cfg.NGLayers = ng
	}
	cfg.GreedySampler = true

	provider := llama.NewLlamaProvider(resolver, cfg)
	ld := loader.NewStubLoader()
	bus := events.NewChannelBus(16)
	rt := runtime.NewIntegratedModelRuntime(provider, ld, resolver, bus)

	ptab := process.NewProcessTable()
	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(ptab, queue)
	coord := scheduler.NewEventCoordinator(sched, bus, ptab)
	defer coord.Stop()

	exec := agents.NewExecutor(rt, bus)
	listener := agents.NewEventListener(bus, exec, ptab, areg)
	defer listener.Stop()

	rtr := router.NewRouter(bus, areg, ptab)
	defer rtr.Stop()

	fmt.Fprintf(stderr, "TinyBrain Benchmark | mode=%s profile=%s backend=%s\n", mode, profile.Name, profile.Probe.Backend)

	done := make(chan string)
	
	taskID := fmt.Sprintf("bench-%d", time.Now().UnixNano())

	var t0, tFirstToken, tDone time.Time
	var peakRAM, peakVRAM uint64

	bus.Subscribe(events.TypeTaskCreated, func(ev events.Event) { t0 = ev.Timestamp })
	bus.Subscribe(events.TypeProcessStateChanged, func(ev events.Event) {
		payload := ev.Payload.(events.ProcessStateChangedPayload)
		if payload.NewState == process.Running.String() {
			tFirstToken = ev.Timestamp // Approximate TTFT since adapters are sync
		}
	})

	bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
		payload := ev.Payload.(events.TaskCompletedPayload)
		if payload.TaskID == taskID {
			tDone = ev.Timestamp
			done <- payload.Result
		}
	})

	// Background monitor for peak RAM/VRAM
	stopMonitor := make(chan struct{})
	go func() {
		for {
			select {
			case <-stopMonitor:
				return
			case <-time.After(50 * time.Millisecond):
				for _, p := range ptab.List() {
					if p.MemoryUsage > peakRAM { peakRAM = p.MemoryUsage }
					if p.VRAMUsage > peakVRAM { peakVRAM = p.VRAMUsage }
				}
			}
		}
	}()

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  taskID,
		Input:   prompt,
		AgentID: targetAgent,
	}, time.Now()))

	var result string
	select {
	case result = <-done:
	case <-time.After(5 * time.Minute):
		fmt.Fprintln(stderr, "error: task timed out")
		close(stopMonitor)
		return 1
	}
	close(stopMonitor)

	var parsed struct {
		Text           string `json:"text"`
		TokensProduced int    `json:"tokens_produced"`
	}
	_ = json.Unmarshal([]byte(result), &parsed)

	ttft := tFirstToken.Sub(t0).Seconds()
	totalTime := tDone.Sub(tFirstToken).Seconds()
	tps := 0.0
	if totalTime > 0 {
		tps = float64(parsed.TokensProduced) / totalTime
	}

	fmt.Fprintf(stdout, "\n--- Benchmark Results ---\nMode: %s\nTTFT: %.3fs\nTPS: %.2f\nPeak RAM: %d MB\nPeak VRAM: %d MB\nTokens: %d\n",
		mode, ttft, tps, peakRAM/1024/1024, peakVRAM/1024/1024, parsed.TokensProduced)

	return 0
}
