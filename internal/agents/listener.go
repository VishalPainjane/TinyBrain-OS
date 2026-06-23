package agents

import (
	"fmt"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

// EventListener bridges ProcessStateChanged to Agent execution.
type EventListener struct {
	bus   events.EventBus
	exec  *Executor
	ptab  *process.ProcessTable
	reg   *registry.AgentRegistry
	now   func() time.Time
	unsub []func()

	taskInputs map[string]string
}

// NewEventListener creates and starts a new EventListener.
func NewEventListener(bus events.EventBus, exec *Executor, ptab *process.ProcessTable, reg *registry.AgentRegistry) *EventListener {
	l := &EventListener{
		bus:        bus,
		exec:       exec,
		ptab:       ptab,
		reg:        reg,
		now:        time.Now,
		taskInputs: make(map[string]string),
	}
	l.start()
	return l
}

func (l *EventListener) start() {
	u1 := l.bus.Subscribe(events.TypeTaskCreated, l.handleTaskCreated)
	u2 := l.bus.Subscribe(events.TypeProcessStateChanged, l.handleStateChanged)
	l.unsub = append(l.unsub, u1, u2)
}

// Stop unsubscribes the listener.
func (l *EventListener) Stop() {
	for _, u := range l.unsub {
		u()
	}
}

func (l *EventListener) handleTaskCreated(ev events.Event) {
	payload, ok := ev.Payload.(events.TaskCreatedPayload)
	if !ok {
		return
	}
	// cache input for later
	l.taskInputs[payload.TaskID] = payload.Input
}

func (l *EventListener) handleStateChanged(ev events.Event) {
	payload, ok := ev.Payload.(events.ProcessStateChangedPayload)
	if !ok {
		return
	}

	if payload.NewState != process.Running.String() {
		return
	}

	p, err := l.ptab.Get(payload.PID)
	if err != nil {
		fmt.Printf("[agent-listener] error: process %s not found\n", payload.PID)
		return
	}

	def, err := l.reg.GetAgent(p.AgentRef)
	if err != nil {
		fmt.Printf("[agent-listener] error: agent definition %s not found\n", p.AgentRef)
		return
	}

	input := l.taskInputs[p.TaskID]

	plugin := NewSamplePlugin(def.ID, def.ModelProfile)
	req := TaskRequest{
		TaskID: p.TaskID,
		PID:    p.PID,
		Input:  input,
	}

	// Execute asynchronously to not block the event dispatcher
	go func() {
		res, err := l.exec.Run(plugin, req)
		if err != nil {
			fmt.Printf("[agent-listener] task %s failed: %v\n", p.TaskID, err)
		}

		_ = l.ptab.UpdateState(p.PID, process.Terminated)

		l.bus.Publish(events.NewEvent(events.TypeTaskCompleted, events.TaskCompletedPayload{
			TaskID: p.TaskID,
			Result: res.Output,
		}, l.now()))
	}()
}
