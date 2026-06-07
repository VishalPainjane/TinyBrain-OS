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
