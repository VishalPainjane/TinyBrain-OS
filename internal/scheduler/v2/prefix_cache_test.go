package v2

import (
	"testing"
	"time"
)

// MockWorker to satisfy InferenceWorker
type MockWorker struct {
	StepCount int
}

func (m *MockWorker) ExecuteStep(payloads []StepPayload) (map[string]int32, error) {
	m.StepCount++
	res := make(map[string]int32)
	for _, p := range payloads {
		if !p.IsPrefill {
			// Generate a token
			res[p.SeqID] = 100
		} else {
			// Actually wait, prefill also can generate a token if it's the final chunk
			res[p.SeqID] = 100
		}
	}
	return res, nil
}
func (m *MockWorker) SwapOut(seqID string) error              { return nil }
func (m *MockWorker) SwapIn(seqID string, blocks []int) error { return nil }
func (m *MockWorker) UpdateWeights(path string) error         { return nil }
func (m *MockWorker) Shutdown() error                         { return nil }

func TestPrefixCaching(t *testing.T) {
	reqChan := make(chan *SequenceGroup, 10)
	worker := &MockWorker{}
	
	// Create scheduler with MaxBlocks=100 and PrefillChunkSize=100
	scheduler := NewScheduler(worker, 100, reqChan, 100)

	// User 1 sends a prompt of 32 tokens (exactly 2 blocks)
	prompt1 := make([]int32, 32)
	for i := 0; i < 32; i++ {
		prompt1[i] = int32(i + 1000)
	}

	seq1 := &Sequence{
		ID:           "seq-1",
		PromptTokens: prompt1,
		State:        SequenceStatePending,
		Metrics:      &SequenceMetrics{ArrivalTimestamp: time.Now()},
	}

	group1 := &SequenceGroup{
		RequestID:        "group-1",
		ArrivalTimestamp: time.Now(),
		Sequences:        map[string]*Sequence{seq1.ID: seq1},
	}

	// Submit User 1
	reqChan <- group1
	
	// Execute Step 1 (Ingest and Prefill)
	scheduler.Step()
	
	if seq1.PrefillProcessedCount != 32 {
		t.Fatalf("Expected PrefillProcessedCount=32, got %d", seq1.PrefillProcessedCount)
	}
	
	// Step 2 to allow cacheFullyProcessedBlocks to capture the generated blocks
	scheduler.Step()

	// Now User 2 arrives with the EXACT same 32 token prompt + 5 more tokens
	prompt2 := make([]int32, 37)
	copy(prompt2, prompt1)
	for i := 32; i < 37; i++ {
		prompt2[i] = int32(i + 2000)
	}

	seq2 := &Sequence{
		ID:           "seq-2",
		PromptTokens: prompt2,
		State:        SequenceStatePending,
		Metrics:      &SequenceMetrics{ArrivalTimestamp: time.Now()},
	}

	group2 := &SequenceGroup{
		RequestID:        "group-2",
		ArrivalTimestamp: time.Now(),
		Sequences:        map[string]*Sequence{seq2.ID: seq2},
	}

	reqChan <- group2
	
	// Execute Step 3 (Ingest User 2)
	// User 2 should immediately match the 32 tokens in the cache during DrainLoop
	scheduler.Step()
	
	if seq2.PrefillProcessedCount != 32 {
		t.Fatalf("Prefix cache failed. Expected seq2 to skip 32 tokens, got %d", seq2.PrefillProcessedCount)
	}
	
	if len(seq2.LogicalKVBlocks) < 2 {
		t.Fatalf("Prefix cache failed. Expected seq2 to have 2 cached blocks mapped, got %d", len(seq2.LogicalKVBlocks))
	}
	
	if seq2.LogicalKVBlocks[0] != seq1.LogicalKVBlocks[0] || seq2.LogicalKVBlocks[1] != seq1.LogicalKVBlocks[1] {
		t.Fatalf("Physical blocks were not shared! seq1: %v, seq2: %v", seq1.LogicalKVBlocks, seq2.LogicalKVBlocks)
	}

	// Verify allocator refcounts
	alloc := scheduler.BlockAllocator
	if alloc.RefCounts[seq2.LogicalKVBlocks[0]] != 2 {
		t.Fatalf("Expected RefCount=2 for shared block, got %d", alloc.RefCounts[seq2.LogicalKVBlocks[0]])
	}
}
