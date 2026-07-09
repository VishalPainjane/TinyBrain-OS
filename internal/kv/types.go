package kv

import "time"

// Tier is the storage location for a KV block.
// See docs/architecture/memory.md.
type Tier int

const (
	// TierVRAM holds active KV in GPU memory.
	TierVRAM Tier = iota
	// TierRAM holds compressed or warm KV in host memory.
	TierRAM
	// TierNVMe holds cold KV on local storage.
	TierNVMe
)

var tierNames = [...]string{
	TierVRAM: "VRAM",
	TierRAM:  "RAM",
	TierNVMe: "NVMe",
}

// String returns the canonical tier name.
func (t Tier) String() string {
	if t < TierVRAM || t > TierNVMe {
		return "UNKNOWN"
	}
	return tierNames[t]
}

// Block is a tracked KV cache block in the pool.
type Block struct {
	KVCacheID           string
	PID                 string
	SizeBytes           uint64
	CompressedSizeBytes uint64
	Tier                Tier
	LastAccess          time.Time
}
