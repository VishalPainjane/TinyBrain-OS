package process

import (
	"errors"
	"fmt"
	"sync"
)

// ErrNotFound is returned when a process PID does not exist in the table.
var ErrNotFound = errors.New("process not found")

// ErrDuplicatePID is returned when Create is called with an existing PID.
var ErrDuplicatePID = errors.New("duplicate PID")

// ErrNotTerminated is returned when Delete is called on a non-terminated process.
var ErrNotTerminated = errors.New("process is not terminated")

// ProcessTable stores process records with O(1) lookup by PID.
// See docs/contracts/process.md.
type ProcessTable struct {
	mu    sync.RWMutex
	byPID map[string]Process
}

// NewProcessTable returns an empty process table.
func NewProcessTable() *ProcessTable {
	return &ProcessTable{
		byPID: make(map[string]Process),
	}
}

// Create inserts a new process in NEW state. PID must be unique.
func (t *ProcessTable) Create(p Process) error {
	if p.PID == "" {
		return fmt.Errorf("PID is required")
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.byPID[p.PID]; exists {
		return ErrDuplicatePID
	}

	p.State = New
	t.byPID[p.PID] = p
	return nil
}

// Get returns the process for pid. Lookup is O(1) via map.
func (t *ProcessTable) Get(pid string) (Process, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	p, ok := t.byPID[pid]
	if !ok {
		return Process{}, ErrNotFound
	}
	return p, nil
}

// List returns a snapshot of all processes in the table.
func (t *ProcessTable) List() []Process {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make([]Process, 0, len(t.byPID))
	for _, p := range t.byPID {
		out = append(out, p)
	}
	return out
}

// UpdateState sets the lifecycle state for pid. Transition validation is owned by the scheduler.
func (t *ProcessTable) UpdateState(pid string, state ProcessState) error {
	if !state.Valid() {
		return fmt.Errorf("invalid process state: %s", state)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	p, ok := t.byPID[pid]
	if !ok {
		return ErrNotFound
	}

	p.State = state
	t.byPID[pid] = p
	return nil
}

// Delete removes a terminated process from the table.
func (t *ProcessTable) Delete(pid string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	p, ok := t.byPID[pid]
	if !ok {
		return ErrNotFound
	}
	if p.State != Terminated {
		return ErrNotTerminated
	}

	delete(t.byPID, pid)
	return nil
}
