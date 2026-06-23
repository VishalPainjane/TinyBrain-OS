package router

import (
	"fmt"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

// Router listens for TaskCreated events, resolves the agent, and spawns a process.
type Router struct {
	bus   events.EventBus
	reg   *registry.AgentRegistry
	ptab  *process.ProcessTable
	now   func() time.Time
	unsub func()
}

// NewRouter creates and starts a new Router.
func NewRouter(bus events.EventBus, reg *registry.AgentRegistry, ptab *process.ProcessTable) *Router {
	r := &Router{
		bus:  bus,
		reg:  reg,
		ptab: ptab,
		now:  time.Now,
	}
	r.start()
	return r
}

func (r *Router) start() {
	r.unsub = r.bus.Subscribe(events.TypeTaskCreated, r.handleTaskCreated)
}

// Stop unsubscribes the router from the event bus.
func (r *Router) Stop() {
	if r.unsub != nil {
		r.unsub()
	}
}

func (r *Router) handleTaskCreated(ev events.Event) {
	payload, ok := ev.Payload.(events.TaskCreatedPayload)
	if !ok {
		return
	}

	// Resolve agent. For now, we assume the user requested a specific AgentID.
	// In the future, this will use task capabilities to find an agent.
	agentID := payload.AgentID
	if agentID == "" {
		fmt.Printf("[router] error: AgentID required for Task %s\n", payload.TaskID)
		return
	}

	def, err := r.reg.GetAgent(agentID)
	if err != nil {
		fmt.Printf("[router] error: agent %s not found for Task %s\n", agentID, payload.TaskID)
		return
	}

	// Create process
	pid := fmt.Sprintf("p-%s", payload.TaskID)
	p := process.Process{
		PID:      pid,
		AgentRef: def.ID,
		Priority: def.Priority,
		TaskID:   payload.TaskID,
	}

	if err := r.ptab.Create(p); err != nil {
		fmt.Printf("[router] error: failed to create process for Task %s: %v\n", payload.TaskID, err)
		return
	}

	// Publish ProcessSpawned
	r.bus.Publish(events.NewEvent(events.TypeProcessSpawned, events.ProcessSpawnedPayload{
		PID:      pid,
		AgentRef: def.ID,
		TaskID:   payload.TaskID,
	}, r.now()))
}
