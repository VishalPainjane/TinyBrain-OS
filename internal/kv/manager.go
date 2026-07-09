package kv

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
	"github.com/klauspost/compress/zstd"
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
	VRAMUsage() uint64
}

// StubManager is an in-memory KV block pool that emits KVStored/KVLoaded events.
// Real llama.cpp export is deferred; Save/Load update tier metadata only.
type StubManager struct {
	mu       sync.Mutex
	blocks   map[string]Block
	ramCache map[string][]byte // simulates TierRAM physical storage
	bus      events.EventBus
	now      func() time.Time
	
	encoder *zstd.Encoder
	decoder *zstd.Decoder
}

// NewStubManager returns a stub KV manager backed by bus for lifecycle events.
func NewStubManager(bus events.EventBus) *StubManager {
	enc, _ := zstd.NewWriter(nil)
	dec, _ := zstd.NewReader(nil)
	
	return &StubManager{
		blocks:   make(map[string]Block),
		ramCache: make(map[string][]byte),
		bus:      bus,
		now:      time.Now,
		encoder:  enc,
		decoder:  dec,
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

// VRAMUsage returns the sum of SizeBytes for all blocks currently in VRAM.
func (m *StubManager) VRAMUsage() uint64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	var used uint64
	for _, block := range m.blocks {
		if block.Tier == TierVRAM {
			used += block.SizeBytes
		}
	}
	return used
}

func (m *StubManager) saveLocked(kvCacheID string) (Block, error) {
	block, ok := m.blocks[kvCacheID]
	if !ok {
		return Block{}, ErrBlockNotFound
	}
	if block.Tier != TierVRAM {
		return Block{}, fmt.Errorf("%w: save requires VRAM, got %s", ErrInvalidTier, block.Tier)
	}

	// 1. Simulate VRAM physical memory extraction
	// We create a dummy tensor block. Repeating pattern compresses well.
	dummyTensor := bytes.Repeat([]byte("tensor_data_"), int(block.SizeBytes/12)+1)
	if len(dummyTensor) > int(block.SizeBytes) {
		dummyTensor = dummyTensor[:block.SizeBytes]
	}

	start := m.now()

	// 2. Compress the block using Zstandard
	compressed := m.encoder.EncodeAll(dummyTensor, make([]byte, 0, len(dummyTensor)))

	latency := m.now().Sub(start).Milliseconds()
	ratio := float64(len(dummyTensor)) / float64(len(compressed))

	// 3. Store in RAM cache simulation
	m.ramCache[kvCacheID] = compressed

	block.Tier = TierRAM
	block.CompressedSizeBytes = uint64(len(compressed))
	block.LastAccess = m.now()
	m.blocks[kvCacheID] = block
	
	// Emit telemetry
	m.publish(events.TypeKVCompressed, events.KVCompressedPayload{
		KVCacheID:        block.KVCacheID,
		CompressionRatio: ratio,
		LatencyMs:        latency,
	})

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

	compressed, ok := m.ramCache[kvCacheID]
	if !ok {
		return Block{}, fmt.Errorf("physical RAM cache missing for block %s", kvCacheID)
	}

	start := m.now()

	// 1. Decompress from RAM
	decompressed, err := m.decoder.DecodeAll(compressed, nil)
	if err != nil {
		return Block{}, fmt.Errorf("failed to decompress block: %w", err)
	}
	
	// Ensure sizes match logic
	if uint64(len(decompressed)) != block.SizeBytes {
		// Log warning or handle mismatch
	}

	latency := m.now().Sub(start).Milliseconds()

	// 2. Clean up RAM cache
	delete(m.ramCache, kvCacheID)

	block.Tier = TierVRAM
	block.CompressedSizeBytes = 0
	block.LastAccess = m.now()
	m.blocks[kvCacheID] = block
	
	// Emit telemetry
	m.publish(events.TypeKVDecompressed, events.KVDecompressedPayload{
		KVCacheID: block.KVCacheID,
		LatencyMs: latency,
	})

	return block, nil
}

func (m *StubManager) publish(eventType events.Type, payload any) {
	if m.bus == nil {
		return
	}
	m.bus.Publish(events.NewEvent(eventType, payload, m.now()))
}
