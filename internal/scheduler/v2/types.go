package v2

import (
	"sync"
	"time"
)

// SequenceState represents the lifecycle state of a sequence during iteration-level scheduling.
type SequenceState int

const (
	SequenceStatePending SequenceState = iota
	SequenceStateRunning
	SequenceStatePreempted
	SequenceStateFinished
)

// BlockSize defines the number of tokens stored in a single physical VRAM block.
const BlockSize = 16

type SequenceMetrics struct {
	ArrivalTimestamp   time.Time
	PrefillStart       time.Time
	PrefillEnd         time.Time     // Set when PrefillProcessedCount == len(PromptTokens)
	TTFT               time.Duration // PrefillEnd.Sub(ArrivalTimestamp)
	DecodeTokenTimes   []time.Time   // Timestamps of each generated token
	LastTokenTimestamp time.Time
}

// Sequence represents the mutable execution context of a single request.
type Sequence struct {
	ID                    string
	PromptTokens          []int32
	GeneratedTokens       []int32
	PrefillProcessedCount int // Tracks how many prompt tokens have been processed
	State                 SequenceState
	LogicalKVBlocks       []int // Maps to physical blocks in the Data Plane
	
	// Anti-Thrashing Tracking
	StepsExecutedCount     int
	LastAllocatedTimestamp time.Time

	// Metrics tracking
	Metrics *SequenceMetrics

	// Storing limits on the struct avoids needing to pass them into IsFinished on every check
	MaxTokens int 
	EOSToken  int32
}


// IsFinished evaluates if the sequence has reached its max tokens or generated the EOS token.
func (s *Sequence) IsFinished() bool {
	if s.State == SequenceStateFinished {
		return true
	}
	
	if s.MaxTokens > 0 && len(s.GeneratedTokens) >= s.MaxTokens {
		return true
	}
	
	if len(s.GeneratedTokens) > 0 {
		lastToken := s.GeneratedTokens[len(s.GeneratedTokens)-1]
		if lastToken == s.EOSToken {
			return true
		}
	}
	
	return false
}

// SequenceGroup wraps one or more Sequences to support parallel sampling or beam search.
type SequenceGroup struct {
	RequestID        string
	Sequences        map[string]*Sequence
	ArrivalTimestamp time.Time
}

// SchedulerQueue represents the MLFQ-style queue for managing sequence states at the token boundary.
type SchedulerQueue struct {
	mu                 sync.Mutex
	Waiting            []*SequenceGroup // New requests, waiting for prefill
	Running            []*SequenceGroup // Currently allocated and decoding
	Swapped            []*SequenceGroup // Preempted to Host RAM due to VRAM pressure
	DeferredFreeBlocks []int            // Blocks scheduled for safe deallocation
}

type StepPayload struct {
	SeqID             string
	IsPrefill         bool
	Tokens            []int32 // PromptTokens for prefill, or single delta token for decode
	LogicalKVBlocks   []int   // Physical VRAM block mapping managed by Go
	AbsolutePositions []int32 // Absolute sequence positions for RoPE encoding
	History           []int32 // Complete sequence history for sampling penalty
}

// InferenceWorker defines the boundary contract between the Go scheduler and the CGO Data Plane.
type InferenceWorker interface {
	// ExecuteStep executes one forward pass (iteration) for a batch of payloads.
	ExecuteStep(payloads []StepPayload) (map[string]int32, error)
	
	// SwapOut moves a sequence's KV cache from VRAM to Host RAM (Preemption).
	SwapOut(seqID string) error
	
	// SwapIn moves a sequence's KV cache back to VRAM (Resumption).
	SwapIn(seqID string, newBlockIDs []int) error
}
