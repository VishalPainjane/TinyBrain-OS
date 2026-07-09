//go:build cgo

package v2

// #cgo CFLAGS: -I.
// #cgo LDFLAGS: -L. -ltinybrain -lcudart -lcublas -L/usr/local/cuda/lib64 -lstdc++ -lm
/*
#include <stdlib.h>

#include "tinybrain_backend.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"unsafe"
)

// CGOWorker implements InferenceWorker using CGO bridging to the CUDA Data Plane.
type CGOWorker struct {
	mu sync.Mutex
	arenaPtr unsafe.Pointer
	maxBatch int
	maxTokensPerChunk int
	maxBlocksPerSeq int
}

func NewCGOWorker() *CGOWorker {
	maxBatch := 128
	maxTokensPerChunk := 1024
	maxBlocksPerSeq := 512

	// Allocate a fixed-size continuous block to act as our pinned serialization space
	vocabSize := 32000
	arenaSize := (2 * maxBatch * maxTokensPerChunk * 4) + (maxBatch * maxBlocksPerSeq * 4) + (4 * maxBatch * 4) + (maxBatch * vocabSize * 4)
	arenaPtr := C.malloc(C.size_t(arenaSize))

	return &CGOWorker{
		arenaPtr:          arenaPtr,
		maxBatch:          maxBatch,
		maxTokensPerChunk: maxTokensPerChunk,
		maxBlocksPerSeq:   maxBlocksPerSeq,
	}
}

// Init initializes the CUDA backend with the given model path.
func (w *CGOWorker) Init(modelPath string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	cModelPath := C.CString(modelPath)
	defer C.free(unsafe.Pointer(cModelPath))

	var config C.tb_config_t
	// We leave config zeroed out so backend.cu uses defaults (or we could set them)
	config.vocab_size = 32000
	config.num_layers = 22
	config.hidden_dim = 2048
	config.num_heads = 32
	config.head_dim = 64
	config.block_size = 16
	config.max_blocks = 512

	errCode := C.tb_init(config, cModelPath)
	if errCode != C.TB_SUCCESS {
		return errors.New("failed to initialize CUDA backend")
	}

	return nil
}

// ExecuteStep implements InferenceWorker by marshaling StepPayloads into flat C arrays.
func (w *CGOWorker) ExecuteStep(payloads []StepPayload) (map[string]int32, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	numSeqs := len(payloads)
	if numSeqs == 0 {
		return nil, nil
	}
	if numSeqs > w.maxBatch {
		return nil, errors.New("batch size exceeds maxBatch")
	}

	// 1. Flattening Metrics Loop
	var totalTokens, totalBlocks int
	for _, p := range payloads {
		totalTokens += len(p.Tokens)
		totalBlocks += len(p.LogicalKVBlocks)
	}

	if totalTokens > w.maxBatch*w.maxTokensPerChunk {
		return nil, errors.New("token chunk size exceeds arena limit")
	}
	if totalBlocks > w.maxBatch*w.maxBlocksPerSeq {
		return nil, errors.New("block chunk size exceeds arena limit")
	}

	// 2. Buffer Allocation & Packing using Persistent Arena
	basePtr := uintptr(w.arenaPtr)

	tokenDataPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * w.maxTokensPerChunk * 4)

	tokenOffsetsPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * 4)

	tokenLengthsPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * 4)

	blockDataPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * w.maxBlocksPerSeq * 4)

	blockOffsetsPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * 4)

	blockLengthsPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * 4)

	tokenPositionsPtr := unsafe.Pointer(basePtr)
	basePtr += uintptr(w.maxBatch * w.maxTokensPerChunk * 4)

	outLogitsPtr := unsafe.Pointer(basePtr)

	// Create Go slice headers over the C memory for zero-copy filling
	var tokenData []int32
	var tokenPositions []int32
	if totalTokens > 0 {
		tokenData = unsafe.Slice((*int32)(tokenDataPtr), totalTokens)
		tokenPositions = unsafe.Slice((*int32)(tokenPositionsPtr), totalTokens)
	}
	
	var blockData []int32
	if totalBlocks > 0 {
		blockData = unsafe.Slice((*int32)(blockDataPtr), totalBlocks)
	}

	tokenOffsets := unsafe.Slice((*int32)(tokenOffsetsPtr), numSeqs)
	tokenLengths := unsafe.Slice((*int32)(tokenLengthsPtr), numSeqs)
	blockOffsets := unsafe.Slice((*int32)(blockOffsetsPtr), numSeqs)
	blockLengths := unsafe.Slice((*int32)(blockLengthsPtr), numSeqs)

	tIdx := 0
	bIdx := 0

	for i, p := range payloads {
		// Pack Tokens
		tokenOffsets[i] = int32(tIdx)
		tokenLengths[i] = int32(len(p.Tokens))
		fmt.Printf("[worker_cgo] Seq %s Phase: %v Tokens: %v Positions: %v\n", p.SeqID, p.IsPrefill, p.Tokens, p.AbsolutePositions)
		for j, t := range p.Tokens {
			tokenData[tIdx] = t
			tokenPositions[tIdx] = p.AbsolutePositions[j]
			tIdx++
		}

		// Pack Blocks
		blockOffsets[i] = int32(bIdx)
		blockLengths[i] = int32(len(p.LogicalKVBlocks))
		for _, b := range p.LogicalKVBlocks {
			blockData[bIdx] = int32(b)
			bIdx++
		}
	}

	// 3. The CGO Execution Call
	fmt.Printf("\n=== [DEBUG 1: GO EXIT] ===\n")
	n := 6
	if len(tokenData) < n { n = len(tokenData) }
	fmt.Printf("Tokens (First %d): %v\n", n, tokenData[:n])
	fmt.Printf("Positions (First %d): %v\n", n, tokenPositions[:n])
	fmt.Printf("==========================\n")

	errCode := C.tb_execute_step(
		C.int(numSeqs),
		(*C.int)(tokenDataPtr),
		(*C.int)(tokenOffsetsPtr),
		(*C.int)(tokenLengthsPtr),
		(*C.int)(blockDataPtr),
		(*C.int)(blockOffsetsPtr),
		(*C.int)(blockLengthsPtr),
		(*C.int)(tokenPositionsPtr),
		(*C.float)(outLogitsPtr),
	)
	
	if errCode != C.TB_SUCCESS {
		return nil, errors.New("CUDA backend execution failed")
	}

	// Extract Results mapping strictly back to Go
	vocabSize := 32000
	allLogits := unsafe.Slice((*float32)(outLogitsPtr), numSeqs*vocabSize)
	results := make(map[string]int32, numSeqs)
	
	for i, p := range payloads {
		logits := allLogits[i*vocabSize : (i+1)*vocabSize]
		
		// Repetition Penalty
		for _, tokenID := range p.History {
			if tokenID >= 0 && int(tokenID) < vocabSize {
				val := logits[tokenID]
				if val > 0 {
					logits[tokenID] = val / 1.2
				} else {
					logits[tokenID] = val * 1.2
				}
			}
		}

		// --- NEW: Temperature & Top-K Sampling ---
		temperature := float32(0.7)
		topK := 40

		// 1. Apply Temperature and track original indices
		type TokenLogit struct {
			ID    int32
			Value float32
		}
		tokenLogits := make([]TokenLogit, vocabSize)
		for j := 0; j < vocabSize; j++ {
			tokenLogits[j] = TokenLogit{
				ID:    int32(j),
				Value: logits[j] / temperature,
			}
		}

		// 2. Sort by highest logit value
		sort.Slice(tokenLogits, func(a, b int) bool {
			return tokenLogits[a].Value > tokenLogits[b].Value
		})

		// 3. Keep only Top-K and calculate Softmax denominator
		maxLogit := tokenLogits[0].Value
		var sumExp float32 = 0.0
		for j := 0; j < topK; j++ {
			// Subtract maxLogit for numerical stability before exp()
			sumExp += float32(math.Exp(float64(tokenLogits[j].Value - maxLogit)))
		}

		// 4. Sample from the probability distribution (Roulette Wheel)
		randomDraw := rand.Float32() // random float between 0.0 and 1.0
		var cumulativeProb float32 = 0.0
		selectedToken := tokenLogits[0].ID // default to best

		for j := 0; j < topK; j++ {
			prob := float32(math.Exp(float64(tokenLogits[j].Value-maxLogit))) / sumExp
			cumulativeProb += prob
			if randomDraw <= cumulativeProb {
				selectedToken = tokenLogits[j].ID
				break
			}
		}

		results[p.SeqID] = selectedToken
	}

	return results, nil
}

// SwapOut moves a sequence's KV cache from VRAM to Host RAM.
func (w *CGOWorker) SwapOut(seqID string) error {
	// Stub for Data Plane swap out call
	return nil
}

// SwapIn moves a sequence's KV cache back to VRAM.
func (w *CGOWorker) SwapIn(seqID string, newBlockIDs []int) error {
	// Stub for Data Plane swap in call
	return nil
}
