package kv

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
)

var (
	// ErrBlockNotFound is returned when a KV cache ID is unknown.
	ErrBlockNotFound = errors.New("kv block not found")
	// ErrBlockExists is returned when Allocate is called with a duplicate ID.
	ErrBlockExists = errors.New("kv block already exists")
	// ErrInvalidTier is returned when a tier transition is not allowed.
	ErrInvalidTier = errors.New("invalid kv tier transition")
)

// Manager coordinates KV block allocation and tier movement.
// See docs/architecture/memory.md and tasks/011-kv-manager.md.
type Manager interface {
	Allocate(kvCacheID, pid string, sizeBytes uint64) error
	Save(kvCacheID string) error
	Load(kvCacheID string) error
	Get(kvCacheID string) (Block, error)
	Delete(kvCacheID string) error
}

// StubManager is an in-memory KV block pool that emits KVStored/KVLoaded events.
// Real llama.cpp export is deferred; Save/Load update tier metadata only.
type StubManager struct {
	mu     sync.Mutex
	blocks map[string]Block
	bus    events.EventBus
	now    func() time.Time
}

// NewStubManager returns a stub KV manager backed by bus for lifecycle events.
func NewStubManager(bus events.EventBus) *StubManager {
	return &StubManager{
		blocks: make(map[string]Block),
		bus:    bus,
		now:    time.Now,
	}
}

// Allocate registers a new KV block in VRAM for pid.
func (m *StubManager) Allocate(kvCacheID, pid string, sizeBytes uint64) error {
	if kvCacheID == "" {
		return fmt.Errorf("kv cache ID is required")
	}
	if pid == "" {
		return fmt.Errorf("PID is required")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.blocks[kvCacheID]; exists {
		return ErrBlockExists
	}

	now := m.now()
	m.blocks[kvCacheID] = Block{
		KVCacheID:  kvCacheID,
		PID:        pid,
		SizeBytes:  sizeBytes,
		Tier:       TierVRAM,
		LastAccess: now,
	}
	return nil
}

// Save moves a block from VRAM to RAM and publishes KVStored.
func (m *StubManager) Save(kvCacheID string) error {
	m.mu.Lock()
	block, err := m.saveLocked(kvCacheID)
	m.mu.Unlock()
	if err != nil {
		return err
	}

	m.publish(events.TypeKVStored, events.KVStoredPayload{
		KVCacheID: block.KVCacheID,
		PID:       block.PID,
	})
	return nil
}

// Load moves a block from RAM to VRAM and publishes KVLoaded.
func (m *StubManager) Load(kvCacheID string) error {
	m.mu.Lock()
	block, err := m.loadLocked(kvCacheID)
	m.mu.Unlock()
	if err != nil {
		return err
	}

	m.publish(events.TypeKVLoaded, events.KVLoadedPayload{
		KVCacheID: block.KVCacheID,
		PID:       block.PID,
	})
	return nil
}

// Get returns block metadata for kvCacheID.
func (m *StubManager) Get(kvCacheID string) (Block, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	block, ok := m.blocks[kvCacheID]
	if !ok {
		return Block{}, ErrBlockNotFound
	}
	return block, nil
}

// Delete removes a block from the pool.
func (m *StubManager) Delete(kvCacheID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.blocks[kvCacheID]; !ok {
		return ErrBlockNotFound
	}
	delete(m.blocks, kvCacheID)
	return nil
}

func (m *StubManager) saveLocked(kvCacheID string) (Block, error) {
	block, ok := m.blocks[kvCacheID]
	if !ok {
		return Block{}, ErrBlockNotFound
	}
	if block.Tier != TierVRAM {
		return Block{}, fmt.Errorf("%w: save requires VRAM, got %s", ErrInvalidTier, block.Tier)
	}

	block.Tier = TierRAM
	block.LastAccess = m.now()
	m.blocks[kvCacheID] = block
	return block, nil
}

func (m *StubManager) loadLocked(kvCacheID string) (Block, error) {
	block, ok := m.blocks[kvCacheID]
	if !ok {
		return Block{}, ErrBlockNotFound
	}
	if block.Tier != TierRAM {
		return Block{}, fmt.Errorf("%w: load requires RAM, got %s", ErrInvalidTier, block.Tier)
	}

	block.Tier = TierVRAM
	block.LastAccess = m.now()
	m.blocks[kvCacheID] = block
	return block, nil
}

func (m *StubManager) publish(eventType events.Type, payload any) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(events.NewEvent(eventType, payload, m.now()))
}
