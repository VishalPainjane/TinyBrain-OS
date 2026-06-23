package router

import (
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func TestRouter_HandleTaskCreated(t *testing.T) {
	bus := events.NewChannelBus(10)

	reg := registry.NewAgentRegistry()
	reg.RegisterAgent(registry.AgentDefinition{
		ID:       "test-agent",
		Priority: 1,
	})

	ptab := process.NewProcessTable()

	r := NewRouter(bus, reg, ptab)
	defer r.Stop()

	spawned := make(chan events.Event, 1)
	bus.Subscribe(events.TypeProcessSpawned, func(ev events.Event) {
		spawned <- ev
	})

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  "task-123",
		AgentID: "test-agent",
		Input:   "hello",
	}, time.Now()))

	select {
	case ev := <-spawned:
		payload, ok := ev.Payload.(events.ProcessSpawnedPayload)
		if !ok {
			t.Fatalf("expected ProcessSpawnedPayload")
		}
		if payload.TaskID != "task-123" {
			t.Errorf("expected TaskID task-123, got %s", payload.TaskID)
		}
		if payload.AgentRef != "test-agent" {
			t.Errorf("expected AgentRef test-agent, got %s", payload.AgentRef)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for ProcessSpawned event")
	}

	p, err := ptab.Get("p-task-123")
	if err != nil {
		t.Fatalf("process not created: %v", err)
	}
	if p.State != process.New {
		t.Errorf("expected new process state to be NEW, got %v", p.State)
	}
}

func TestRouter_MissingAgent(t *testing.T) {
	bus := events.NewChannelBus(10)

	reg := registry.NewAgentRegistry()
	ptab := process.NewProcessTable()

	r := NewRouter(bus, reg, ptab)
	defer r.Stop()

	spawned := make(chan events.Event, 1)
	bus.Subscribe(events.TypeProcessSpawned, func(ev events.Event) {
		spawned <- ev
	})

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{
		TaskID:  "task-456",
		AgentID: "missing-agent",
		Input:   "hello",
	}, time.Now()))

	select {
	case <-spawned:
		t.Fatalf("expected no process to be spawned for missing agent")
	case <-time.After(100 * time.Millisecond):
		// Expected
	}
}
