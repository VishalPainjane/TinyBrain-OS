package scheduler

import (
	"errors"
	"fmt"
	"sync"
	"time"

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
	case process.Ready, process.Preempted:
		if stored.State == process.Preempted {
			if err := s.table.UpdateState(p.PID, process.Ready); err != nil {
				return err
			}
		}
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

// MLFQScheduler applies multi-level feedback queue policy with token quanta and preemption.
type MLFQScheduler struct {
	table *process.ProcessTable
	queue *MLFQQueue

	mu               sync.Mutex
	queueLevel       map[string]int
	quantumTokens    map[string]int
	tokensSinceBoost int
	lastBoost        time.Time
}

// NewMLFQScheduler returns a scheduler backed by table and MLFQ queue.
func NewMLFQScheduler(table *process.ProcessTable, queue *MLFQQueue) *MLFQScheduler {
	return &MLFQScheduler{
		table:         table,
		queue:         queue,
		queueLevel:    make(map[string]int),
		quantumTokens: make(map[string]int),
		lastBoost:     time.Now(),
	}
}

// Enqueue admits p at its MLFQ level and transitions NEW or PREEMPTED to READY.
func (s *MLFQScheduler) Enqueue(p process.Process) error {
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
	case process.Ready, process.Preempted:
		if stored.State == process.Preempted {
			if err := s.table.UpdateState(p.PID, process.Ready); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidEnqueue, p.PID, stored.State)
	}

	ready, err := s.table.Get(p.PID)
	if err != nil {
		return err
	}

	s.mu.Lock()
	level, ok := s.queueLevel[p.PID]
	if !ok {
		level = QueueLevelFromPriority(ready.Priority)
		s.queueLevel[p.PID] = level
		s.quantumTokens[p.PID] = 0
	}
	s.mu.Unlock()

	return s.queue.EnqueueAt(level, ready)
}

// Schedule dequeues the highest-priority ready process, preempting a lower-priority runner if needed.
func (s *MLFQScheduler) Schedule() (process.Process, error) {
	next, err := s.queue.Peek()
	if err != nil {
		return process.Process{}, err
	}

	s.mu.Lock()
	nextLevel := s.queueLevel[next.PID]
	s.mu.Unlock()

	if running, ok := s.findRunning(); ok {
		s.mu.Lock()
		runningLevel := s.queueLevel[running.PID]
		s.mu.Unlock()
		if nextLevel < runningLevel {
			if err := s.preemptAndRequeue(running); err != nil {
				return process.Process{}, err
			}
		}
	}

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

// Preempt marks pid PREEMPTED in the process table.
func (s *MLFQScheduler) Preempt(pid string) error {
	if _, err := s.table.Get(pid); err != nil {
		return err
	}
	return s.table.UpdateState(pid, process.Preempted)
}

// Boost moves all queued processes to Q0 and resets per-process levels.
func (s *MLFQScheduler) Boost() error {
	s.queue.BoostAll()

	s.mu.Lock()
	defer s.mu.Unlock()

	for pid := range s.queueLevel {
		s.queueLevel[pid] = 0
		s.quantumTokens[pid] = 0
	}
	s.tokensSinceBoost = 0
	s.lastBoost = time.Now()
	return nil
}

// RecordToken counts one generated token for pid and demotes when the level quantum is exceeded.
// It may trigger Boost per BoostEveryTokens or BoostEveryDuration.
func (s *MLFQScheduler) RecordToken(pid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	level, ok := s.queueLevel[pid]
	if !ok {
		return fmt.Errorf("unknown process %s", pid)
	}

	s.quantumTokens[pid]++
	if s.quantumTokens[pid] >= TokenQuantum(level) {
		s.demoteLocked(pid)
	}

	s.tokensSinceBoost++
	if s.tokensSinceBoost >= BoostEveryTokens || time.Since(s.lastBoost) >= BoostEveryDuration*time.Second {
		s.boostLocked()
	}
	return nil
}

// QueueDepths returns per-level queue depths for telemetry.
func (s *MLFQScheduler) QueueDepths() [NumLevels]int {
	return s.queue.Depths()
}

// QueueLevel returns the current MLFQ level for pid.
func (s *MLFQScheduler) QueueLevel(pid string) (int, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	level, ok := s.queueLevel[pid]
	return level, ok
}

func (s *MLFQScheduler) findRunning() (process.Process, bool) {
	for _, p := range s.table.List() {
		if p.State == process.Running {
			return p, true
		}
	}
	return process.Process{}, false
}

func (s *MLFQScheduler) preemptAndRequeue(p process.Process) error {
	if err := s.table.UpdateState(p.PID, process.Preempted); err != nil {
		return err
	}
	preempted, err := s.table.Get(p.PID)
	if err != nil {
		return err
	}
	return s.Enqueue(preempted)
}

func (s *MLFQScheduler) demoteLocked(pid string) {
	level := s.queueLevel[pid]
	if level < NumLevels-1 {
		level++
	}
	s.queueLevel[pid] = level
	s.quantumTokens[pid] = 0
}

func (s *MLFQScheduler) boostLocked() {
	s.queue.BoostAll()
	for pid := range s.queueLevel {
		s.queueLevel[pid] = 0
		s.quantumTokens[pid] = 0
	}
	s.tokensSinceBoost = 0
	s.lastBoost = time.Now()
}
