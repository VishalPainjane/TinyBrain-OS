package events_test

import (
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
)

func TestNewEvent(t *testing.T) {
	at := time.Unix(1_700_000_000, 0).UTC()
	payload := events.TaskCreatedPayload{TaskID: "task-1"}

	ev := events.NewEvent(events.TypeTaskCreated, payload, at)

	if ev.Type != events.TypeTaskCreated {
		t.Errorf("Type = %q, want TaskCreated", ev.Type)
	}
	if !ev.Timestamp.Equal(at) {
		t.Errorf("Timestamp = %v, want %v", ev.Timestamp, at)
	}

	got, ok := ev.Payload.(events.TaskCreatedPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want TaskCreatedPayload", ev.Payload)
	}
	if got.TaskID != "task-1" {
		t.Errorf("TaskID = %q, want task-1", got.TaskID)
	}
}

func TestAllTypes_ContainsCoreEvents(t *testing.T) {
	types := events.AllTypes()
	if len(types) != 13 {
		t.Fatalf("AllTypes() len = %d, want 13", len(types))
	}

	want := map[events.Type]bool{
		events.TypeTaskCreated:         true,
		events.TypeTaskAssigned:        true,
		events.TypeProcessSpawned:      true,
		events.TypeProcessStateChanged: true,
		events.TypeAgentStarted:        true,
		events.TypeAgentStopped:        true,
		events.TypeModelLoaded:         true,
		events.TypeModelUnloaded:       true,
		events.TypeSwapStarted:         true,
		events.TypeSwapCompleted:       true,
		events.TypeKVStored:            true,
		events.TypeKVLoaded:            true,
		events.TypeTaskCompleted:       true,
	}

	seen := make(map[events.Type]bool, len(types))
	for _, eventType := range types {
		if !want[eventType] {
			t.Errorf("unexpected event type: %q", eventType)
		}
		if seen[eventType] {
			t.Errorf("duplicate event type: %q", eventType)
		}
		seen[eventType] = true
	}
}

func TestEventTypeIdentification(t *testing.T) {
	tests := []struct {
		name    string
		typ     events.Type
		payload any
	}{
		{"TaskCreated", events.TypeTaskCreated, events.TaskCreatedPayload{TaskID: "t1"}},
		{"TaskAssigned", events.TypeTaskAssigned, events.TaskAssignedPayload{TaskID: "t1", AgentID: "a1"}},
		{"ProcessSpawned", events.TypeProcessSpawned, events.ProcessSpawnedPayload{PID: "p1"}},
		{"ModelLoaded", events.TypeModelLoaded, events.ModelLoadedPayload{ModelID: "m1"}},
		{"TaskCompleted", events.TypeTaskCompleted, events.TaskCompletedPayload{TaskID: "t1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := events.NewEvent(tt.typ, tt.payload, time.Now())
			if ev.Type != tt.typ {
				t.Errorf("Type = %q, want %q", ev.Type, tt.typ)
			}
			if ev.Payload == nil {
				t.Error("Payload is nil")
			}
		})
	}
}
