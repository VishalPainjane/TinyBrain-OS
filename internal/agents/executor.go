package agents

import (
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
)

// Executor runs agent plugins via the runtime API and publishes lifecycle events.
type Executor struct {
	rt  RuntimeAPI
	bus events.EventBus
	now func() time.Time
}

// NewExecutor returns an executor that delegates inference to rt.
func NewExecutor(rt RuntimeAPI, bus events.EventBus) *Executor {
	return &Executor{
		rt:  rt,
		bus: bus,
		now: time.Now,
	}
}

// Run executes agent for req, emitting AgentStarted and AgentStopped on bus.
func (e *Executor) Run(agent Agent, req TaskRequest) (TaskResult, error) {
	e.bus.Publish(events.NewEvent(events.TypeAgentStarted, events.AgentStartedPayload{
		AgentID: agent.ID(),
		PID:     req.PID,
	}, e.now()))

	result, err := agent.Execute(ExecuteContext{Runtime: e.rt}, req)
	if result.AgentID == "" {
		result.AgentID = agent.ID()
	}
	if result.TaskID == "" {
		result.TaskID = req.TaskID
	}

	e.bus.Publish(events.NewEvent(events.TypeAgentStopped, events.AgentStoppedPayload{
		AgentID: agent.ID(),
		PID:     req.PID,
	}, e.now()))

	return result, err
}
