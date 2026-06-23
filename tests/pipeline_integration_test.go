package tests

import (
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/agents"
	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/router"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
	"github.com/VishalPainjane/TinyBrain-OS/internal/scheduler"
)

func TestEventPipelineEndToEnd(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(32)
	areg := registry.NewAgentRegistry()
	if err := registry.LoadAgentsYAML(filepath.Join("..", "testdata", "fleet.yaml"), areg); err != nil {
		t.Fatalf("load fleet: %v", err)
	}

	ptab := process.NewProcessTable()
	queue := scheduler.NewMLFQQueue()
	sched := scheduler.NewMLFQScheduler(ptab, queue)
	coord := scheduler.NewEventCoordinator(sched, bus, ptab)
	defer coord.Stop()

	// Stub runtime
	provider := runtime.NewStubProvider()
	rt := runtime.NewModelRuntime(provider, bus)

	if err := rt.LoadModel("tinyllama-q4"); err != nil {
		t.Fatalf("LoadModel: %v", err)
	}
	if err := rt.LoadModel("qwen-coder"); err != nil {
		t.Fatalf("LoadModel: %v", err)
	}

	exec := agents.NewExecutor(rt, bus)
	listener := agents.NewEventListener(bus, exec, ptab, areg)
	defer listener.Stop()

	rtr := router.NewRouter(bus, areg, ptab)
	defer rtr.Stop()

	var trace []events.Type
	var mu sync.Mutex

	bus.Subscribe(events.TypeTaskCreated, func(ev events.Event) {
		mu.Lock()
		defer mu.Unlock()
		trace = append(trace, ev.Type)
	})
	bus.Subscribe(events.TypeProcessSpawned, func(ev events.Event) {
		mu.Lock()
		defer mu.Unlock()
		trace = append(trace, ev.Type)
	})
	bus.Subscribe(events.TypeProcessStateChanged, func(ev events.Event) {
		mu.Lock()
		defer mu.Unlock()
		payload, ok := ev.Payload.(events.ProcessStateChangedPayload)
		if ok && payload.NewState == process.Running.String() {
			trace = append(trace, ev.Type)
		}
	})
	bus.Subscribe(events.TypeAgentStarted, func(ev events.Event) {
		mu.Lock()
		defer mu.Unlock()
		trace = append(trace, ev.Type)
	})

	done := make(chan struct{})
	bus.Subscribe(events.TypeTaskCompleted, func(ev events.Event) {
		mu.Lock()
		defer mu.Unlock()
		trace = append(trace, ev.Type)
		close(done)
	})

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  "task-integration-1",
		Input:   "hello event pipeline",
		AgentID: "sample-alpha",
	}, time.Now()))

	select {
	case <-done:
		// Let concurrent subscribers flush their channels to the trace
		time.Sleep(50 * time.Millisecond)
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline timed out waiting for TaskCompleted")
	}

	mu.Lock()
	defer mu.Unlock()

	expected := []events.Type{
		events.TypeTaskCreated,
		events.TypeProcessSpawned,
		events.TypeProcessStateChanged,
		events.TypeAgentStarted,
		events.TypeTaskCompleted,
	}

	if len(trace) < len(expected) {
		t.Fatalf("trace length %d < expected %d. Trace: %v", len(trace), len(expected), trace)
	}

	found := make(map[events.Type]bool)
	for _, tr := range trace {
		found[tr] = true
	}

	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("trace missing expected event: %s", exp)
		}
	}
}
