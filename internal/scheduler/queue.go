package scheduler

import (
	"errors"
	"sync"

	"github.com/VishalPainjane/TinyBrain-OS/internal/process"
)

var ErrQueueEmpty = errors.New("queue is empty")

// Queue admits and selects processes in scheduling order.
// See docs/contracts/scheduler.md.
type Queue interface {
	Enqueue(p process.Process) error
	Dequeue() (process.Process, error)
	Peek() (process.Process, error)
	Depth() int
}

// FIFOQueue is a first-in-first-out process queue for v1 scheduling.
type FIFOQueue struct {
	mu    sync.Mutex
	items []process.Process
}

// NewFIFOQueue returns an empty FIFO queue.
func NewFIFOQueue() *FIFOQueue {
	return &FIFOQueue{
		items: make([]process.Process, 0),
	}
}

// Enqueue appends p to the tail of the queue.
func (q *FIFOQueue) Enqueue(p process.Process) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, item := range q.items {
		if item.PID == p.PID {
			return ErrAlreadyQueued
		}
	}

	q.items = append(q.items, p)
	return nil
}

// Dequeue removes and returns the front process.
func (q *FIFOQueue) Dequeue() (process.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return process.Process{}, ErrQueueEmpty
	}

	p := q.items[0]
	q.items = q.items[1:]
	return p, nil
}

// Peek returns the front process without removing it.
func (q *FIFOQueue) Peek() (process.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return process.Process{}, ErrQueueEmpty
	}
	return q.items[0], nil
}

// Depth returns the number of queued processes.
func (q *FIFOQueue) Depth() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// MLFQQueue implements multi-level feedback queues Q0 (highest) through Q3.
type MLFQQueue struct {
	mu     sync.Mutex
	levels [NumLevels][]process.Process
}

// NewMLFQQueue returns an empty MLFQ queue set.
func NewMLFQQueue() *MLFQQueue {
	return &MLFQQueue{}
}

// Enqueue appends p using QueueLevelFromPriority(p.Priority).
func (q *MLFQQueue) Enqueue(p process.Process) error {
	return q.EnqueueAt(QueueLevelFromPriority(p.Priority), p)
}

// EnqueueAt appends p to the FIFO at level.
func (q *MLFQQueue) EnqueueAt(level int, p process.Process) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if level < 0 || level >= NumLevels {
		return errors.New("invalid MLFQ level")
	}

	for _, item := range q.levels[level] {
		if item.PID == p.PID {
			return ErrAlreadyQueued
		}
	}
	for i := 0; i < NumLevels; i++ {
		if i == level {
			continue
		}
		for _, item := range q.levels[i] {
			if item.PID == p.PID {
				return ErrAlreadyQueued
			}
		}
	}

	q.levels[level] = append(q.levels[level], p)
	return nil
}

// Dequeue removes and returns the highest-priority non-empty queue head.
func (q *MLFQQueue) Dequeue() (process.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for level := 0; level < NumLevels; level++ {
		if len(q.levels[level]) == 0 {
			continue
		}
		p := q.levels[level][0]
		q.levels[level] = q.levels[level][1:]
		return p, nil
	}
	return process.Process{}, ErrQueueEmpty
}

// Peek returns the next process that Dequeue would return.
func (q *MLFQQueue) Peek() (process.Process, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for level := 0; level < NumLevels; level++ {
		if len(q.levels[level]) == 0 {
			continue
		}
		return q.levels[level][0], nil
	}
	return process.Process{}, ErrQueueEmpty
}

// Depth returns the total number of queued processes across all levels.
func (q *MLFQQueue) Depth() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.totalDepthLocked()
}

// Depths returns per-level queue depths for telemetry (Q0..Q3).
func (q *MLFQQueue) Depths() [NumLevels]int {
	q.mu.Lock()
	defer q.mu.Unlock()

	var depths [NumLevels]int
	for i := 0; i < NumLevels; i++ {
		depths[i] = len(q.levels[i])
	}
	return depths
}

// BoostAll moves every queued process to Q0, preserving relative order within each level.
func (q *MLFQQueue) BoostAll() {
	q.mu.Lock()
	defer q.mu.Unlock()

	var boosted []process.Process
	for level := 0; level < NumLevels; level++ {
		boosted = append(boosted, q.levels[level]...)
		q.levels[level] = nil
	}
	q.levels[0] = boosted
}

func (q *MLFQQueue) totalDepthLocked() int {
	n := 0
	for i := 0; i < NumLevels; i++ {
		n += len(q.levels[i])
	}
	return n
}
