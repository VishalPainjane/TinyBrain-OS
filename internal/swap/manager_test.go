package swap_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/kv"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
	"github.com/VishalPainjane/TinyBrain-OS/internal/swap"
)

func testProcess(pid string) process.Process {
	return process.Process{
		PID:            pid,
		AgentRef:       "agent-1",
		Priority:       5,
		KVCacheID:      "kv-" + pid,
		LastExecution:  time.Unix(1_700_000_000, 0).UTC(),
		TokensProduced: 1,
		TaskID:         "task-1",
	}
}

func newSwapFixture(t *testing.T) (*process.ProcessTable, kv.Manager, events.EventBus, *swap.StubManager) {
	t.Helper()
	table := process.NewProcessTable()
	bus := events.NewChannelBus(8)
	kvm := kv.NewStubManager(bus)
	mgr := swap.NewStubManager(table, kvm, bus, 8000)
	return table, kvm, bus, mgr
}

func TestStubManager_SwapOutMovesKVAndHibernates(t *testing.T) {
	t.Parallel()

	table, kvm, bus, mgr := newSwapFixture(t)
	p := testProcess("pid-a")
	p.LastExecution = time.Now().Add(-100 * time.Second)

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := table.UpdateState("pid-a", process.Waiting); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}
	if err := kvm.Allocate(p.KVCacheID, p.PID, 2048); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}

	var started, completed sync.WaitGroup
	started.Add(1)
	completed.Add(1)
	unsubStart := bus.Subscribe(events.TypeSwapStarted, func(ev events.Event) {
		defer started.Done()
		payload, ok := ev.Payload.(events.SwapStartedPayload)
		if !ok {
			t.Errorf("SwapStarted payload type = %T", ev.Payload)
			return
		}
		if payload.FromModelID != "VRAM" || payload.ToModelID != "RAM" {
			t.Errorf("SwapStarted = %+v, want VRAM→RAM", payload)
		}
	})
	unsubDone := bus.Subscribe(events.TypeSwapCompleted, func(ev events.Event) {
		defer completed.Done()
		payload, ok := ev.Payload.(events.SwapCompletedPayload)
		if !ok {
			t.Errorf("SwapCompleted payload type = %T", ev.Payload)
			return
		}
		if payload.FromModelID != "VRAM" || payload.ToModelID != "RAM" {
			t.Errorf("SwapCompleted = %+v, want VRAM→RAM", payload)
		}
	})
	defer unsubStart()
	defer unsubDone()

	if err := mgr.SwapOut("pid-a"); err != nil {
		t.Fatalf("SwapOut() error = %v", err)
	}

	wait(t, &started)
	wait(t, &completed)

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Hibernated {
		t.Errorf("State = %s, want HIBERNATED", stored.State)
	}

	block, err := kvm.Get(p.KVCacheID)
	if err != nil {
		t.Fatalf("kv Get() error = %v", err)
	}
	if block.Tier != kv.TierRAM {
		t.Errorf("KV tier = %s, want RAM", block.Tier)
	}
}

func TestStubManager_SwapInRestoresKVAndReady(t *testing.T) {
	t.Parallel()

	table, kvm, _, mgr := newSwapFixture(t)
	p := testProcess("pid-a")
	p.LastExecution = time.Now().Add(-100 * time.Second)

	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := table.UpdateState("pid-a", process.Waiting); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}
	if err := kvm.Allocate(p.KVCacheID, p.PID, 2048); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}
	if err := mgr.SwapOut("pid-a"); err != nil {
		t.Fatalf("SwapOut() error = %v", err)
	}

	if err := mgr.SwapIn("pid-a"); err != nil {
		t.Fatalf("SwapIn() error = %v", err)
	}

	stored, err := table.Get("pid-a")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if stored.State != process.Ready {
		t.Errorf("State = %s, want READY", stored.State)
	}

	block, err := kvm.Get(p.KVCacheID)
	if err != nil {
		t.Fatalf("kv Get() error = %v", err)
	}
	if block.Tier != kv.TierVRAM {
		t.Errorf("KV tier = %s, want VRAM", block.Tier)
	}
}

func TestStubManager_SwapOutErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		prepare func(table *process.ProcessTable, kvm kv.Manager) string
		wantErr error
	}{
		{
			name: "running process",
			prepare: func(table *process.ProcessTable, kvm kv.Manager) string {
				p := testProcess("pid-run")
				p.LastExecution = time.Now().Add(-1 * time.Hour)
				_ = table.Create(p)
				_ = table.UpdateState("pid-run", process.Running)
				_ = kvm.Allocate(p.KVCacheID, p.PID, 512)
				return "pid-run"
			},
			wantErr: swap.ErrCannotSwapRunning,
		},
		{
			name: "not idle",
			prepare: func(table *process.ProcessTable, kvm kv.Manager) string {
				p := testProcess("pid-fresh")
				p.LastExecution = time.Now()
				_ = table.Create(p)
				_ = table.UpdateState("pid-fresh", process.Waiting)
				_ = kvm.Allocate(p.KVCacheID, p.PID, 512)
				return "pid-fresh"
			},
			wantErr: swap.ErrNotIdle,
		},
		{
			name: "missing kv id",
			prepare: func(table *process.ProcessTable, _ kv.Manager) string {
				p := testProcess("pid-nokv")
				p.KVCacheID = ""
				p.LastExecution = time.Now().Add(-1 * time.Hour)
				_ = table.Create(p)
				_ = table.UpdateState("pid-nokv", process.Waiting)
				return "pid-nokv"
			},
			wantErr: swap.ErrNoKVCache,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			table, kvm, bus, mgr := newSwapFixture(t)
			pid := tt.prepare(table, kvm)
			if err := mgr.SwapOut(pid); !errors.Is(err, tt.wantErr) {
				t.Fatalf("SwapOut() error = %v, want %v", err, tt.wantErr)
			}
			_ = bus
		})
	}
}

func TestStubManager_SwapOutRespectsSchedulerIdleHeuristic(t *testing.T) {
	t.Parallel()

	p := testProcess("pid-a")
	p.LastExecution = time.Now().Add(-100 * time.Second)

	table, kvm, _, mgr := newSwapFixture(t)
	if err := table.Create(p); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := table.UpdateState("pid-a", process.Waiting); err != nil {
		t.Fatalf("UpdateState() error = %v", err)
	}
	stored, _ := table.Get("pid-a")
	if !scheduler.ShouldSwap(stored, time.Now(), 0, 8000) {
		t.Fatal("fixture should be eligible for swap")
	}
	if err := kvm.Allocate(p.KVCacheID, p.PID, 1024); err != nil {
		t.Fatalf("Allocate() error = %v", err)
	}
	if err := mgr.SwapOut("pid-a"); err != nil {
		t.Fatalf("SwapOut() error = %v", err)
	}
}

func wait(t *testing.T, wg *sync.WaitGroup) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("event not received in time")
	}
}
