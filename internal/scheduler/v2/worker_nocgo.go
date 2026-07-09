//go:build !cgo

package v2

import (
	"sync"
)

// CGOWorker implements InferenceWorker fallback for non-CGO / non-CUDA environments.
type CGOWorker struct {
	mu                sync.Mutex
	maxBatch          int
	maxTokensPerChunk int
	maxBlocksPerSeq   int
}

func NewCGOWorker() *CGOWorker {
	return &CGOWorker{
		maxBatch:          128,
		maxTokensPerChunk: 1024,
		maxBlocksPerSeq:   512,
	}
}

func (w *CGOWorker) Init(modelPath string) error {
	return nil
}

func (w *CGOWorker) ExecuteStep(payloads []StepPayload) (map[string]int32, error) {
	results := make(map[string]int32, len(payloads))
	for _, p := range payloads {
		results[p.SeqID] = 100
	}
	return results, nil
}

func (w *CGOWorker) SwapOut(seqID string) error {
	return nil
}

func (w *CGOWorker) SwapIn(seqID string, newBlockIDs []int) error {
	return nil
}
