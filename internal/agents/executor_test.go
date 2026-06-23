package agents_test

import (
	"encoding/json"
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/agents"
	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func TestExecutor_SamplePluginRunsTask(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(8)
	provider := runtime.NewStubProvider()
	rt := runtime.NewModelRuntime(provider, bus)
	if err := rt.LoadModel("tinyllama-q4"); err != nil {
		t.Fatalf("LoadModel: %v", err)
	}

	exec := agents.NewExecutor(rt, bus)
	plugin := agents.NewSamplePlugin("sample-alpha", "tinyllama-q4")

	req := agents.TaskRequest{
		TaskID: "task-1",
		PID:    "pid-1",
		Input:  "summarize the design doc",
	}

	result, err := exec.Run(plugin, req)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.AgentID != "sample-alpha" {
		t.Errorf("AgentID = %q, want sample-alpha", result.AgentID)
	}
	if result.Output == "" {
		t.Fatal("expected non-empty structured JSON output")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(result.Output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if parsed["task_id"] != "task-1" {
		t.Errorf("task_id = %v, want task-1", parsed["task_id"])
	}
	if parsed["agent_id"] != "sample-alpha" {
		t.Errorf("agent_id = %v, want sample-alpha", parsed["agent_id"])
	}
}

func TestExecutor_PublishesAgentLifecycleEvents(t *testing.T) {
	t.Parallel()

	bus := events.NewChannelBus(8)
	var startedWG, stoppedWG sync.WaitGroup
	startedWG.Add(1)
	stoppedWG.Add(1)

	unsubStart := bus.Subscribe(events.TypeAgentStarted, func(ev events.Event) {
		defer startedWG.Done()
		payload, ok := ev.Payload.(events.AgentStartedPayload)
		if !ok {
			t.Errorf("AgentStarted payload type = %T", ev.Payload)
			return
		}
		if payload.AgentID != "sample-beta" || payload.PID != "p" {
			t.Errorf("AgentStarted = %+v, want sample-beta / p", payload)
		}
	})
	unsubStop := bus.Subscribe(events.TypeAgentStopped, func(ev events.Event) {
		defer stoppedWG.Done()
		payload, ok := ev.Payload.(events.AgentStoppedPayload)
		if !ok {
			t.Errorf("AgentStopped payload type = %T", ev.Payload)
			return
		}
		if payload.AgentID != "sample-beta" || payload.PID != "p" {
			t.Errorf("AgentStopped = %+v, want sample-beta / p", payload)
		}
	})
	defer unsubStart()
	defer unsubStop()

	provider := runtime.NewStubProvider()
	rt := runtime.NewModelRuntime(provider, bus)
	if err := rt.LoadModel("m1"); err != nil {
		t.Fatalf("LoadModel: %v", err)
	}

	exec := agents.NewExecutor(rt, bus)
	plugin := agents.NewSamplePlugin("sample-beta", "m1")

	_, err := exec.Run(plugin, agents.TaskRequest{
		TaskID: "t",
		PID:    "p",
		Input:  "hello",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	waitEvents(t, &startedWG, "AgentStarted")
	waitEvents(t, &stoppedWG, "AgentStopped")
}

func waitEvents(t *testing.T, wg *sync.WaitGroup, name string) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Errorf("timed out waiting for %s", name)
	}
}
