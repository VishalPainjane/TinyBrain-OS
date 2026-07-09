package v2

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Engine is the master orchestrator running on a dedicated goroutine.
// It manages the token scheduling loop continuously and maintains
// concurrent-safe client streaming channels.
type Engine struct {
	scheduler      *Scheduler
	submitChan     chan *SequenceGroup
	streamRegistry map[string]chan int32
	registryMu     sync.RWMutex
	seqCounter     uint64
}

// NewEngine creates a new Engine instance.
func NewEngine(scheduler *Scheduler, submitChan chan *SequenceGroup) *Engine {
	return &Engine{
		scheduler:      scheduler,
		submitChan:     submitChan,
		streamRegistry: make(map[string]chan int32),
	}
}

// Submit generates a unique SequenceID, registers a token streaming channel,
// and pushes the sequence onto the lock-free ingest channel for scheduling.
func (e *Engine) Submit(prompt []int32, maxTokens int32, eosTokenID int32) (<-chan int32, string) {
	seqID := fmt.Sprintf("seq-%d", atomic.AddUint64(&e.seqCounter, 1))

	seq := &Sequence{
		ID:           seqID,
		PromptTokens: prompt,
		State:        SequenceStatePending,
		MaxTokens:    int(maxTokens),
		EOSToken:     eosTokenID,
		Metrics: &SequenceMetrics{
			ArrivalTimestamp: time.Now(),
		},
	}

	group := &SequenceGroup{
		RequestID:        seqID,
		Sequences:        map[string]*Sequence{seqID: seq},
		ArrivalTimestamp: seq.Metrics.ArrivalTimestamp,
	}

	tokenStream := make(chan int32, 100)

	e.registryMu.Lock()
	e.streamRegistry[seqID] = tokenStream
	e.registryMu.Unlock()

	// Non-blocking handoff to the scheduler
	e.submitChan <- group

	return tokenStream, seqID
}

// Cancel removes a sequence from the active stream registry and notifies
// the scheduler to safely preempt the sequence and defer its block deallocation.
func (e *Engine) Cancel(seqID string) {
	e.registryMu.Lock()
	if ch, ok := e.streamRegistry[seqID]; ok {
		// Drain the channel completely before closing to prevent goroutine leaks
		for len(ch) > 0 {
			<-ch
		}
		close(ch)
		delete(e.streamRegistry, seqID)
	}
	e.registryMu.Unlock()
	
	e.scheduler.Cancel(seqID)
}

// Start runs the asynchronous orchestrator loop to drive the scheduler continuously.
// It uses an Adaptive Backoff Strategy to prevent CPU starvation.
func (e *Engine) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Clean up registry and shut down
			e.registryMu.Lock()
			for _, ch := range e.streamRegistry {
				close(ch)
			}
			e.streamRegistry = make(map[string]chan int32)
			e.registryMu.Unlock()
			return
		default:
			// Drive the scheduling step
			nextTokens, finishedSeqs, err := e.scheduler.Step()
			if err != nil {
				fmt.Printf("[Engine] Scheduler step failed: %v\n", err)
				// If hardware fails, backoff briefly
				time.Sleep(time.Millisecond * 10)
				continue
			}

			// Evaluate system queues for idling
			if nextTokens == nil && len(finishedSeqs) == 0 && e.scheduler.IsIdle() {
				// No work to do, sleep briefly to avoid spinning CPU at 100%
				time.Sleep(time.Millisecond * 5)
				continue
			}

			// Stream output tokens to their respective clients
			if nextTokens != nil {
				e.registryMu.RLock()
				for seqID, token := range nextTokens {
					if ch, ok := e.streamRegistry[seqID]; ok {
						// Non-blocking send
						select {
						case ch <- token:
						default:
							// Client isn't reading fast enough, token dropped in buffer
							// Production engines might buffer or disconnect here
						}
					}
				}
				e.registryMu.RUnlock()
			}

			// Identify any sequences that transitioned to Finished
			if len(finishedSeqs) > 0 {
				e.registryMu.Lock()
				for _, seqID := range finishedSeqs {
					if ch, ok := e.streamRegistry[seqID]; ok {
						close(ch)
						delete(e.streamRegistry, seqID)
					}
				}
				e.registryMu.Unlock()
			}
		}
	}
}
