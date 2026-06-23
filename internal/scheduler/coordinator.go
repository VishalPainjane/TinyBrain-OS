package scheduler

import (
	"fmt"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

// EventCoordinator integrates the Scheduler with the EventBus.
type EventCoordinator struct {
	sched Scheduler
	bus   events.EventBus
	ptab  *process.ProcessTable
	now   func() time.Time
	unsub []func()
	quit  chan struct{}
}

// NewEventCoordinator creates and starts a new EventCoordinator.
func NewEventCoordinator(sched Scheduler, bus events.EventBus, ptab *process.ProcessTable) *EventCoordinator {
	c := &EventCoordinator{
		sched: sched,
		bus:   bus,
		ptab:  ptab,
		now:   time.Now,
		quit:  make(chan struct{}),
	}
	c.start()
	return c
}

func (c *EventCoordinator) start() {
	u1 := c.bus.Subscribe(events.TypeProcessSpawned, c.handleSpawned)
	c.unsub = append(c.unsub, u1)

	go c.loop()
}

// Stop halts the event loop and unsubscribes.
func (c *EventCoordinator) Stop() {
	close(c.quit)
	for _, u := range c.unsub {
		u()
	}
}

func (c *EventCoordinator) handleSpawned(ev events.Event) {
	payload, ok := ev.Payload.(events.ProcessSpawnedPayload)
	if !ok {
		return
	}

	p, err := c.ptab.Get(payload.PID)
	if err != nil {
		fmt.Printf("[scheduler] error finding spawned process %s: %v\n", payload.PID, err)
		return
	}

	if err := c.sched.Enqueue(p); err != nil {
		fmt.Printf("[scheduler] error enqueuing process %s: %v\n", payload.PID, err)
		return
	}
}

func (c *EventCoordinator) loop() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-c.quit:
			return
		case <-ticker.C:
			p, err := c.sched.Schedule()
			if err != nil {
				// queue is empty or error
				continue
			}
			// Scheduler updated state to Running.
			c.bus.Publish(events.NewEvent(events.TypeProcessStateChanged, events.ProcessStateChangedPayload{
				PID:      p.PID,
				OldState: process.Ready.String(),
				NewState: process.Running.String(),
			}, c.now()))
		}
	}
}
