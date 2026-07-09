package v2

import (
	"container/list"
	"errors"
	"sort"
	"time"
)

// BlockAllocator manages the physical VRAM blocks using Reference Counting and LRU Eviction.
type BlockAllocator struct {
	FreeBlocks  []int
	RefCounts   []int
	LRUList     *list.List
	LRUElements map[int]*list.Element
	EvictHook   func(blockID int)
}

func NewBlockAllocator(maxBlocks int) *BlockAllocator {
	blocks := make([]int, maxBlocks)
	for i := 0; i < maxBlocks; i++ {
		blocks[i] = i
	}
	return &BlockAllocator{
		FreeBlocks:  blocks,
		RefCounts:   make([]int, maxBlocks),
		LRUList:     list.New(),
		LRUElements: make(map[int]*list.Element),
	}
}

// Allocate pops a free physical block ID, or evicts from LRU if empty.
func (b *BlockAllocator) Allocate() (int, error) {
	if len(b.FreeBlocks) > 0 {
		block := b.FreeBlocks[len(b.FreeBlocks)-1]
		b.FreeBlocks = b.FreeBlocks[:len(b.FreeBlocks)-1]
		b.RefCounts[block] = 1
		return block, nil
	}

	// Evict from LRU
	if b.LRUList.Len() > 0 {
		elem := b.LRUList.Front() // head is oldest (eviction candidate)
		block := elem.Value.(int)
		b.LRUList.Remove(elem)
		delete(b.LRUElements, block)

		if b.EvictHook != nil {
			b.EvictHook(block)
		}

		b.RefCounts[block] = 1
		return block, nil
	}

	return -1, errors.New("out of VRAM blocks")
}

// AllocateN pops N free physical block IDs.
func (b *BlockAllocator) AllocateN(n int) ([]int, error) {
	if b.FreeCount() < n {
		return nil, errors.New("not enough VRAM blocks")
	}

	blocks := make([]int, n)
	for i := 0; i < n; i++ {
		blocks[i], _ = b.Allocate()
	}
	return blocks, nil
}

// Free decrements the reference count and pushes to LRU if 0.
func (b *BlockAllocator) Free(blocks []int) {
	for _, block := range blocks {
		b.RefCounts[block]--
		if b.RefCounts[block] == 0 {
			// Add to tail of LRU (most recently used)
			elem := b.LRUList.PushBack(block)
			b.LRUElements[block] = elem
		}
	}
}

// FreeCount returns the total number of blocks available for allocation (unallocated + evictable LRU).
func (b *BlockAllocator) FreeCount() int {
	return len(b.FreeBlocks) + b.LRUList.Len()
}

// IncRef increments reference count (used for cache hits).
func (b *BlockAllocator) IncRef(blocks []int) {
	for _, block := range blocks {
		if b.RefCounts[block] == 0 {
			// Remove from LRU as it is active again
			if elem, ok := b.LRUElements[block]; ok {
				b.LRUList.Remove(elem)
				delete(b.LRUElements, block)
			}
		}
		b.RefCounts[block]++
	}
}

// Scheduler handles the MLFQ-style token-boundary scheduling for sequences.
type Scheduler struct {
	Queue            *SchedulerQueue
	Worker           InferenceWorker
	BlockAllocator   *BlockAllocator
	PrefixCache      *PrefixCache
	IncomingReqs     <-chan *SequenceGroup
	CancelChan       chan string
	PrefillChunkSize int
	MaxBlocks        int

	// Engine Analytics
	TotalTTFT time.Duration
	CountTTFT int
	TotalITL  time.Duration
	CountITL  int

	// Anti-Thrashing metrics
	RecentSwappedBlocks      int
	StepsSinceThrashingReset int
	IsThrashingActive        bool
}

func NewScheduler(worker InferenceWorker, maxBlocks int, reqChan <-chan *SequenceGroup, prefillChunkSize int) *Scheduler {
	alloc := NewBlockAllocator(maxBlocks)
	cache := NewPrefixCache()
	alloc.EvictHook = cache.Evict

	return &Scheduler{
		Queue:            &SchedulerQueue{},
		Worker:           worker,
		BlockAllocator:   alloc,
		PrefixCache:      cache,
		IncomingReqs:     reqChan,
		CancelChan:       make(chan string, 1000),
		PrefillChunkSize: prefillChunkSize,
		MaxBlocks:        maxBlocks,
	}
}

// calculateDeltaBlocks calculates how many new physical blocks are required for the sequence's next iteration.
func (s *Scheduler) calculateDeltaBlocks(seq *Sequence) int {
	var expectedTokens int
	if seq.PrefillProcessedCount < len(seq.PromptTokens) {
		expectedTokens = seq.PrefillProcessedCount + s.PrefillChunkSize
		if expectedTokens > len(seq.PromptTokens) {
			expectedTokens = len(seq.PromptTokens)
		}
	} else {
		expectedTokens = seq.PrefillProcessedCount + len(seq.GeneratedTokens) + 1
	}

	targetBlocks := (expectedTokens + BlockSize - 1) / BlockSize
	delta := targetBlocks - len(seq.LogicalKVBlocks)
	if delta < 0 {
		return 0
	}
	return delta
}

// Step executes one iteration of the scheduler loop.
func (s *Scheduler) Step() (map[string]int32, []string, error) {
	s.Queue.mu.Lock()
	defer s.Queue.mu.Unlock()

	// 1. Drain Ingest (Lock-Free)
DrainLoop:
	for {
		select {
		case req := <-s.IncomingReqs:
			// Resolve Prefix Cache hits immediately upon arrival
			for _, seq := range req.Sequences {
				matchedBlocks, matchedTokens := s.PrefixCache.Match(seq.PromptTokens, BlockSize)
				if matchedTokens > 0 {
					seq.LogicalKVBlocks = append(seq.LogicalKVBlocks, matchedBlocks...)
					s.BlockAllocator.IncRef(matchedBlocks)
					seq.PrefillProcessedCount = matchedTokens
				}
			}
			s.Queue.Waiting = append(s.Queue.Waiting, req)
		case seqID := <-s.CancelChan:
			s.executeCancellation(seqID)
		default:
			break DrainLoop
		}
	}

	// 2. Schedule Phase
	
	// Ingress Hysteresis Safeguard
	s.StepsSinceThrashingReset++
	if s.StepsSinceThrashingReset >= 10 {
		if s.MaxBlocks > 0 && s.RecentSwappedBlocks > s.MaxBlocks/2 {
			s.IsThrashingActive = true
		}
		s.RecentSwappedBlocks = 0
		s.StepsSinceThrashingReset = 0
	}
	
	if s.IsThrashingActive {
		if s.MaxBlocks > 0 && s.BlockAllocator.FreeCount() >= int(float64(s.MaxBlocks)*0.7) {
			s.IsThrashingActive = false
		}
	}
	
	// Check Minimum Residence Quota for entire running queue
	allBelowQuota := false
	if len(s.Queue.Running) > 0 {
		allBelowQuota = true
		for _, group := range s.Queue.Running {
			for _, seq := range group.Sequences {
				if seq.State != SequenceStateFinished && seq.StepsExecutedCount >= 4 {
					allBelowQuota = false
				}
			}
		}
	}
	
	// A. Swap-In Protocol (FIFO from Swapped)
	if !s.IsThrashingActive && len(s.Queue.Swapped) > 0 {
		sort.Slice(s.Queue.Swapped, func(i, j int) bool {
			return s.Queue.Swapped[i].ArrivalTimestamp.Before(s.Queue.Swapped[j].ArrivalTimestamp)
		})

		var stillSwapped []*SequenceGroup
		for _, group := range s.Queue.Swapped {
			reqBlocks := 0
			for _, seq := range group.Sequences {
				if seq.State != SequenceStateFinished {
					reqBlocks += s.calculateDeltaBlocks(seq)
				}
			}
			
			if s.BlockAllocator.FreeCount() >= reqBlocks {
				for _, seq := range group.Sequences {
					if seq.State == SequenceStateFinished {
						continue
					}
					delta := s.calculateDeltaBlocks(seq)
					newBlocks, _ := s.BlockAllocator.AllocateN(delta)
					seq.LogicalKVBlocks = append(seq.LogicalKVBlocks, newBlocks...)
					s.Worker.SwapIn(seq.ID, newBlocks)
					seq.State = SequenceStateRunning
					seq.LastAllocatedTimestamp = time.Now()
					seq.StepsExecutedCount = 0
					
					// Track swap-in volume
					s.RecentSwappedBlocks += delta
				}
				s.Queue.Running = append(s.Queue.Running, group)
			} else {
				stillSwapped = append(stillSwapped, group)
			}
		}
		s.Queue.Swapped = stillSwapped
	}

	// B. Admit from Waiting
	if !s.IsThrashingActive && !allBelowQuota && len(s.Queue.Waiting) > 0 {
		var stillWaiting []*SequenceGroup
		for _, group := range s.Queue.Waiting {
			reqBlocks := 0
			for _, seq := range group.Sequences {
				if seq.State != SequenceStateFinished {
					reqBlocks += s.calculateDeltaBlocks(seq)
				}
			}
			
			if s.BlockAllocator.FreeCount() >= reqBlocks {
				for _, seq := range group.Sequences {
					delta := s.calculateDeltaBlocks(seq)
					newBlocks, _ := s.BlockAllocator.AllocateN(delta)
					seq.LogicalKVBlocks = append(seq.LogicalKVBlocks, newBlocks...)
					seq.State = SequenceStateRunning
					seq.LastAllocatedTimestamp = time.Now()
					seq.StepsExecutedCount = 0
				}
				s.Queue.Running = append(s.Queue.Running, group)
			} else {
				stillWaiting = append(stillWaiting, group)
			}
		}
		s.Queue.Waiting = stillWaiting
	}

	// C. Trigger OOM Check and Swap-Out Protocol (LIFO from Running)
	totalRequiredDelta := 0
	for _, group := range s.Queue.Running {
		for _, seq := range group.Sequences {
			if seq.State != SequenceStateFinished {
				totalRequiredDelta += s.calculateDeltaBlocks(seq)
			}
		}
	}

	for s.BlockAllocator.FreeCount() < totalRequiredDelta && len(s.Queue.Running) > 0 {
		victimIdx := -1
		for i := 0; i < len(s.Queue.Running); i++ {
			// Enforce Minimum Residence Quota: skip groups containing a sequence with < 4 steps
			canPreempt := true
			for _, seq := range s.Queue.Running[i].Sequences {
				if seq.State != SequenceStateFinished && seq.StepsExecutedCount < 4 {
					canPreempt = false
					break
				}
			}
			if !canPreempt {
				continue
			}

			if victimIdx == -1 || s.Queue.Running[i].ArrivalTimestamp.After(s.Queue.Running[victimIdx].ArrivalTimestamp) {
				victimIdx = i
			}
		}
		
		if victimIdx == -1 {
			// No valid victim could be found due to residence quota
			break
		}
		
		victim := s.Queue.Running[victimIdx]
		s.Queue.Running = append(s.Queue.Running[:victimIdx], s.Queue.Running[victimIdx+1:]...)
		
		groupDeltaSaved := 0
		for _, seq := range victim.Sequences {
			if seq.State == SequenceStateFinished {
				continue
			}
			s.Worker.SwapOut(seq.ID)
			s.BlockAllocator.Free(seq.LogicalKVBlocks)
			
			// Track swap-out volume
			s.RecentSwappedBlocks += len(seq.LogicalKVBlocks)
			
			// Recalculate what we saved by not running this sequence
			groupDeltaSaved += s.calculateDeltaBlocks(seq)
			seq.LogicalKVBlocks = nil
			seq.State = SequenceStatePreempted
			seq.StepsExecutedCount = 0
		}
		
		s.Queue.Swapped = append(s.Queue.Swapped, victim)
		totalRequiredDelta -= groupDeltaSaved
	}

	// D. Allocate safe blocks and build payloads
	stepStart := time.Now()
	var payloads []StepPayload
	for _, group := range s.Queue.Running {
		for _, seq := range group.Sequences {
			if seq.State == SequenceStateFinished {
				continue
			}

			if seq.Metrics == nil {
				seq.Metrics = &SequenceMetrics{
					ArrivalTimestamp: group.ArrivalTimestamp,
				}
			}

			if seq.PrefillProcessedCount == 0 && seq.Metrics.PrefillStart.IsZero() {
				seq.Metrics.PrefillStart = stepStart
			}

			// Safe to allocate now, as LIFO preemption guaranteed FreeCount >= totalRequiredDelta

			delta := s.calculateDeltaBlocks(seq)
			if delta > 0 {
				newBlocks, err := s.BlockAllocator.AllocateN(delta)
				if err != nil {
					// VRAM is exhausted and we couldn't preempt (Quota). Skip executing this sequence this iteration.
					continue
				}
				seq.LogicalKVBlocks = append(seq.LogicalKVBlocks, newBlocks...)
			}

			var history []int32
			history = append(history, seq.PromptTokens...)
			history = append(history, seq.GeneratedTokens...)

			payload := StepPayload{
				SeqID:           seq.ID,
				LogicalKVBlocks: seq.LogicalKVBlocks,
				History:         history,
			}

			if seq.PrefillProcessedCount < len(seq.PromptTokens) {
				payload.IsPrefill = true
				end := seq.PrefillProcessedCount + s.PrefillChunkSize
				if end > len(seq.PromptTokens) {
					end = len(seq.PromptTokens)
				}
				payload.Tokens = seq.PromptTokens[seq.PrefillProcessedCount:end]
				
				payload.AbsolutePositions = make([]int32, len(payload.Tokens))
				for idx := range payload.Tokens {
					payload.AbsolutePositions[idx] = int32(seq.PrefillProcessedCount + idx)
				}
			} else {
				payload.IsPrefill = false
				payload.Tokens = []int32{seq.GeneratedTokens[len(seq.GeneratedTokens)-1]}
				
				absolutePos := len(seq.PromptTokens) + len(seq.GeneratedTokens) - 1
				payload.AbsolutePositions = []int32{int32(absolutePos)}
			}
			
			payloads = append(payloads, payload)
		}
	}

	if len(payloads) == 0 {
		s.deallocateSweep()
		return nil, nil, nil
	}

	// 3. Execute - Unlock queue during hardware execution to allow concurrent cancellations
	s.Queue.mu.Unlock()
	nextTokens, err := s.Worker.ExecuteStep(payloads)
	s.Queue.mu.Lock()
	
	if err != nil {
		s.deallocateSweep()
		return nil, nil, err
	}

	// 4. Process Outputs
	var nextRunning []*SequenceGroup
	var finishedSeqs []string
	for _, group := range s.Queue.Running {
		groupFinished := true
		
		for _, seq := range group.Sequences {
			if seq.State == SequenceStateFinished {
				continue
			}
			
			// Cleanly increment the executed step count
			seq.StepsExecutedCount++
			
			if seq.PrefillProcessedCount < len(seq.PromptTokens) {
				// Processed a prefill chunk
				end := seq.PrefillProcessedCount + s.PrefillChunkSize
				if end > len(seq.PromptTokens) {
					end = len(seq.PromptTokens)
				}
				seq.PrefillProcessedCount = end
				
				if seq.PrefillProcessedCount == len(seq.PromptTokens) {
					seq.Metrics.PrefillEnd = time.Now()
					seq.Metrics.TTFT = seq.Metrics.PrefillEnd.Sub(seq.Metrics.ArrivalTimestamp)
					seq.Metrics.LastTokenTimestamp = seq.Metrics.PrefillEnd
					
					s.TotalTTFT += seq.Metrics.TTFT
					s.CountTTFT++

					if token, ok := nextTokens[seq.ID]; ok {
						seq.GeneratedTokens = append(seq.GeneratedTokens, token)
					}
				}
				
				// Keep it running for next prefill chunk or decode phase
				groupFinished = false
			} else {
				// Processed a decode step
				now := time.Now()
				
				if !seq.Metrics.LastTokenTimestamp.IsZero() {
					itl := now.Sub(seq.Metrics.LastTokenTimestamp)
					s.TotalITL += itl
					s.CountITL++
				}
				
				seq.Metrics.DecodeTokenTimes = append(seq.Metrics.DecodeTokenTimes, now)
				seq.Metrics.LastTokenTimestamp = now

				if token, ok := nextTokens[seq.ID]; ok {
					seq.GeneratedTokens = append(seq.GeneratedTokens, token)
				}
	
				if seq.IsFinished() {
					seq.State = SequenceStateFinished
					s.BlockAllocator.Free(seq.LogicalKVBlocks)
					seq.LogicalKVBlocks = nil
					finishedSeqs = append(finishedSeqs, seq.ID)
				} else {
					groupFinished = false
				}
			}
		}
		
		if !groupFinished {
			nextRunning = append(nextRunning, group)
		}
	}
	
	s.Queue.Running = nextRunning
	
	// Cache any newly formed full blocks
	s.cacheFullyProcessedBlocks()

	// Execute the Safe Deallocation Sweep
	s.deallocateSweep()
	
	return nextTokens, finishedSeqs, nil
}

func (s *Scheduler) cacheFullyProcessedBlocks() {
	for _, group := range s.Queue.Running {
		for _, seq := range group.Sequences {
			totalTokens := len(seq.PromptTokens) + len(seq.GeneratedTokens)
			if seq.State != SequenceStateFinished && seq.PrefillProcessedCount < len(seq.PromptTokens) {
				totalTokens = seq.PrefillProcessedCount
			}
			
			numFullBlocks := totalTokens / BlockSize
			
			var allTokens []int32
			allTokens = append(allTokens, seq.PromptTokens...)
			allTokens = append(allTokens, seq.GeneratedTokens...)
			
			for i := 0; i < numFullBlocks; i++ {
				if i >= len(seq.LogicalKVBlocks) {
					break
				}
				
				var parentBlocks []int
				if i > 0 {
					parentBlocks = append([]int(nil), seq.LogicalKVBlocks[:i]...)
				}
				
				blockTokens := allTokens[i*BlockSize : (i+1)*BlockSize]
				blockID := seq.LogicalKVBlocks[i]
				
				s.PrefixCache.Insert(parentBlocks, blockTokens, blockID)
			}
		}
	}
}

// deallocateSweep returns blocks to the free list only after the GPU is done with them.
// Must be called while s.Queue.mu is locked.
func (s *Scheduler) deallocateSweep() {
	if len(s.Queue.DeferredFreeBlocks) > 0 {
		s.BlockAllocator.Free(s.Queue.DeferredFreeBlocks)
		s.Queue.DeferredFreeBlocks = nil
	}
}

// Cancel enqueues a cancellation request to the scheduler tick loop.
// This prevents blocking the HTTP handler or mutating shared state concurrently.
func (s *Scheduler) Cancel(seqID string) {
	select {
	case s.CancelChan <- seqID:
	default:
	}
}

// executeCancellation safely preempts a sequence from the system and defers its block deallocation.
// Must be called from the tick loop while s.Queue.mu is locked.
func (s *Scheduler) executeCancellation(seqID string) {
	var blocksToFree []int

	// Helper to remove sequence and collect its blocks
	removeSeq := func(groups []*SequenceGroup) []*SequenceGroup {
		var remaining []*SequenceGroup
		for _, group := range groups {
			if seq, ok := group.Sequences[seqID]; ok {
				blocksToFree = append(blocksToFree, seq.LogicalKVBlocks...)
				seq.LogicalKVBlocks = nil
				seq.State = SequenceStateFinished
				delete(group.Sequences, seqID)
			}
			if len(group.Sequences) > 0 {
				remaining = append(remaining, group)
			}
		}
		return remaining
	}

	s.Queue.Waiting = removeSeq(s.Queue.Waiting)
	s.Queue.Running = removeSeq(s.Queue.Running)
	s.Queue.Swapped = removeSeq(s.Queue.Swapped)

	if len(blocksToFree) > 0 {
		s.Queue.DeferredFreeBlocks = append(s.Queue.DeferredFreeBlocks, blocksToFree...)
	}
}

// GetEngineStatus returns real-time engine telemetry for monitoring and visualization.
func (s *Scheduler) GetEngineStatus() map[string]interface{} {
	s.Queue.mu.Lock()
	waiting := len(s.Queue.Waiting)
	running := len(s.Queue.Running)
	swapped := len(s.Queue.Swapped)
	s.Queue.mu.Unlock()

	vramUtil := float64(0)
	if s.MaxBlocks > 0 {
		vramUtil = float64(s.MaxBlocks-s.BlockAllocator.FreeCount()) / float64(s.MaxBlocks) * 100.0
	}

	avgTTFT := time.Duration(0)
	if s.CountTTFT > 0 {
		avgTTFT = s.TotalTTFT / time.Duration(s.CountTTFT)
	}

	avgITL := time.Duration(0)
	if s.CountITL > 0 {
		avgITL = s.TotalITL / time.Duration(s.CountITL)
	}

	return map[string]interface{}{
		"vram_utilization_percent": vramUtil,
		"waiting_sequences":        waiting,
		"running_sequences":        running,
		"swapped_sequences":        swapped,
		"avg_ttft_ms":              avgTTFT.Milliseconds(),
		"avg_itl_ms":               avgITL.Milliseconds(),
	}
}

// IsIdle returns true if there are no sequences currently tracked in the queue.
func (s *Scheduler) IsIdle() bool {
	s.Queue.mu.Lock()
	defer s.Queue.mu.Unlock()
	return len(s.Queue.Waiting) == 0 && len(s.Queue.Running) == 0 && len(s.Queue.Swapped) == 0
}
