package loader

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	ErrModelNotFound        = errors.New("model not found")
	ErrModelAlreadyLoaded   = errors.New("model already loaded")
	ErrInvalidState         = errors.New("invalid model state transition")
	ErrStubPathInaccessible = errors.New("stub model path not accessible")
)

type modelEntry struct {
	path  string
	state ModelState
}

// StubLoader is an in-memory loader for v0.4 testing without real mmap.
type StubLoader struct {
	mu       sync.Mutex
	models   map[string]modelEntry
	capacity int
	policy   EvictionPolicy
}

// StubLoaderOption configures a StubLoader.
type StubLoaderOption func(*StubLoader)

// WithCapacity sets the maximum number of in-memory models before eviction runs.
func WithCapacity(n int) StubLoaderOption {
	return func(l *StubLoader) {
		if n > 0 {
			l.capacity = n
		}
	}
}

// WithEvictionPolicy sets the policy used when capacity is exceeded.
func WithEvictionPolicy(p EvictionPolicy) StubLoaderOption {
	return func(l *StubLoader) {
		l.policy = p
	}
}

// NewStubLoader returns a stub model loader.
func NewStubLoader(opts ...StubLoaderOption) *StubLoader {
	l := &StubLoader{
		models: make(map[string]modelEntry),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// Load transitions NOT_LOADED or UNLOADED to ACTIVE.
func (l *StubLoader) Load(modelID, path string) error {
	if err := assertStubPath(path); err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.evictIfNeededLocked(modelID); err != nil {
		return err
	}

	entry, ok := l.models[modelID]
	if !ok || entry.state == StateNotLoaded || entry.state == StateUnloaded {
		l.models[modelID] = modelEntry{path: path, state: StateLoading}
		l.models[modelID] = modelEntry{path: path, state: StateActive}
		l.touchLocked(modelID)
		return nil
	}

	switch entry.state {
	case StateActive, StateWarm, StateLoading:
		return fmt.Errorf("%w: %s", ErrModelAlreadyLoaded, modelID)
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidState, modelID, entry.state)
	}
}

// Unload transitions ACTIVE or WARM to UNLOADED.
func (l *StubLoader) Unload(modelID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.models[modelID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrModelNotFound, modelID)
	}

	switch entry.state {
	case StateActive, StateWarm:
		entry.state = StateUnloading
		entry.state = StateUnloaded
		l.models[modelID] = entry
		return nil
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidState, modelID, entry.state)
	}
}

// Warm transitions ACTIVE to WARM.
func (l *StubLoader) Warm(modelID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.models[modelID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrModelNotFound, modelID)
	}
	if entry.state != StateActive {
		return fmt.Errorf("%w: %s in state %s", ErrInvalidState, modelID, entry.state)
	}

	entry.state = StateWarm
	l.models[modelID] = entry
	l.touchLocked(modelID)
	return nil
}

// Prefetch transitions NOT_LOADED or UNLOADED to WARM without activating inference.
func (l *StubLoader) Prefetch(modelID, path string) error {
	if err := assertStubPath(path); err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.evictIfNeededLocked(modelID); err != nil {
		return err
	}

	entry, ok := l.models[modelID]
	if !ok || entry.state == StateNotLoaded || entry.state == StateUnloaded {
		l.models[modelID] = modelEntry{path: path, state: StateLoading}
		l.models[modelID] = modelEntry{path: path, state: StateWarm}
		l.touchLocked(modelID)
		return nil
	}

	switch entry.state {
	case StateActive, StateWarm, StateLoading:
		return fmt.Errorf("%w: %s", ErrModelAlreadyLoaded, modelID)
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidState, modelID, entry.state)
	}
}

// Evict transitions ACTIVE or WARM to UNLOADED.
func (l *StubLoader) Evict(modelID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.models[modelID]
	if !ok {
		return fmt.Errorf("%w: %s", ErrModelNotFound, modelID)
	}

	switch entry.state {
	case StateActive, StateWarm:
		entry.state = StateUnloading
		entry.state = StateUnloaded
		l.models[modelID] = entry
		return nil
	default:
		return fmt.Errorf("%w: %s in state %s", ErrInvalidState, modelID, entry.state)
	}
}

// State returns the current lifecycle state for modelID.
func (l *StubLoader) State(modelID string) (ModelState, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry, ok := l.models[modelID]
	if !ok {
		return StateNotLoaded, nil
	}
	return entry.state, nil
}

// LRUEvictionPolicy evicts the least recently touched candidate.
type LRUEvictionPolicy struct {
	mu    sync.Mutex
	touch map[string]int64
	clock int64
}

// NewLRUEvictionPolicy returns an LRU eviction policy shell for capacity tests.
func NewLRUEvictionPolicy() *LRUEvictionPolicy {
	return &LRUEvictionPolicy{
		touch: make(map[string]int64),
	}
}

// Touch records modelID as recently used.
func (p *LRUEvictionPolicy) Touch(modelID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.clock++
	p.touch[modelID] = p.clock
}

// SelectVictim returns the least recently touched candidate.
func (p *LRUEvictionPolicy) SelectVictim(candidates []string) string {
	if len(candidates) == 0 {
		return ""
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	victim := candidates[0]
	oldest := p.touch[victim]
	for _, id := range candidates[1:] {
		if ts, ok := p.touch[id]; !ok || ts < oldest {
			victim = id
			oldest = ts
		}
	}
	return victim
}

func (l *StubLoader) evictIfNeededLocked(incomingID string) error {
	if l.capacity <= 0 || l.policy == nil {
		return nil
	}

	loaded := l.loadedIDsLocked()
	if len(loaded) < l.capacity {
		return nil
	}

	for _, id := range loaded {
		if id == incomingID {
			return nil
		}
	}

	victim := l.policy.SelectVictim(loaded)
	if victim == "" {
		return fmt.Errorf("%w: no eviction candidate", ErrInvalidState)
	}

	entry := l.models[victim]
	entry.state = StateUnloading
	entry.state = StateUnloaded
	l.models[victim] = entry
	return nil
}

func (l *StubLoader) loadedIDsLocked() []string {
	var ids []string
	for id, entry := range l.models {
		switch entry.state {
		case StateActive, StateWarm, StateLoading:
			ids = append(ids, id)
		}
	}
	return ids
}

func assertStubPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("%w: %s", ErrStubPathInaccessible, path)
	}
	return nil
}

func (l *StubLoader) touchLocked(modelID string) {
	if tracker, ok := l.policy.(interface{ Touch(string) }); ok {
		tracker.Touch(modelID)
	}
}
