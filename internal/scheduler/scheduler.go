package scheduler

import (
	"errors"
	"fmt"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

var (
	ErrAlreadyQueued   = errors.New("process already queued")
	ErrInvalidEnqueue  = errors.New("process must be in NEW state to enqueue")
	ErrInvalidSchedule = errors.New("scheduled process must be in READY state")
)

// Scheduler applies scheduling policy over the process table and queue.
// See docs/contracts/scheduler.md and tasks/010-scheduler.md.
type Scheduler interface {
	Enqueue(p process.Process) error
	Schedule() (process.Process, error)
	Preempt(pid string) error
	Boost() error
}

// FIFOScheduler is a FIFO scheduling policy skeleton ready for MLFQ migration.
type FIFOScheduler struct {
	table *process.ProcessTable
	queue Queue
}

// NewFIFOScheduler returns a scheduler backed by table and queue.
func NewFIFOScheduler(table *process.ProcessTable, queue Queue) *FIFOScheduler {
	return &FIFOScheduler{
		table: table,
		queue: queue,
	}
}

// Enqueue admits p to the queue and transitions NEW to READY.
func (s *FIFOScheduler) Enqueue(p process.Process) error {
	if p.PID == "" {
		return fmt.Errorf("PID is required")
	}

	stored, err := s.table.Get(p.PID)
	if err != nil {
		return err
	}

	switch stored.State {
	case process.New:
		if err := s.table.UpdateState(p.PID, process.Ready); err != nil {
			return err
		}
	case process.Ready:
		// already ready; may be re-queued after preemption in future tasks
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidEnqueue, p.PID, stored.State)
	}

	ready, err := s.table.Get(p.PID)
	if err != nil {
		return err
	}

	return s.queue.Enqueue(ready)
}

// Schedule dequeues the next process and transitions READY to RUNNING.
func (s *FIFOScheduler) Schedule() (process.Process, error) {
	queued, err := s.queue.Dequeue()
	if err != nil {
		return process.Process{}, err
	}

	stored, err := s.table.Get(queued.PID)
	if err != nil {
		return process.Process{}, err
	}
	if stored.State != process.Ready {
		return process.Process{}, fmt.Errorf("%w: %s in state %s", ErrInvalidSchedule, queued.PID, stored.State)
	}

	if err := s.table.UpdateState(queued.PID, process.Running); err != nil {
		return process.Process{}, err
	}

	return s.table.Get(queued.PID)
}

// Preempt marks pid as PREEMPTED in the process table.
func (s *FIFOScheduler) Preempt(pid string) error {
	if _, err := s.table.Get(pid); err != nil {
		return err
	}
	return s.table.UpdateState(pid, process.Preempted)
}

// Boost is a no-op until MLFQ boost/aging is implemented.
func (s *FIFOScheduler) Boost() error {
	return nil
}
