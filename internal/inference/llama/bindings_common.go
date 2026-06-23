//go:build cgo

package llama

/*
#cgo CFLAGS: -I${SRCDIR}/../../../third_party/llama.cpp/include -I${SRCDIR}/../../../third_party/llama.cpp/ggml/include
#include "llama.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

var (
	backendOnce    sync.Once
	nativeMu       sync.Mutex
	nativeModels   = make(map[string]unsafe.Pointer)
	nativeContexts = make(map[string]unsafe.Pointer)
)

func initBackend() {
	backendOnce.Do(func() {
		C.llama_backend_init()
	})
}

func effectiveBatchSize(cfg LlamaConfig) uint32 {
	batch := cfg.BatchSize
	if batch == 0 {
		batch = 512
	}
	if cfg.ContextSize > 0 && batch > cfg.ContextSize {
		batch = cfg.ContextSize
	}
	if batch == 0 {
		batch = 512
	}
	return batch
}

func freeNativeHandles(modelID string) {
	if ctx, ok := nativeContexts[modelID]; ok {
		C.llama_free((*C.struct_llama_context)(ctx))
		delete(nativeContexts, modelID)
	}
	if model, ok := nativeModels[modelID]; ok {
		C.llama_model_free((*C.struct_llama_model)(model))
		delete(nativeModels, modelID)
	}
}

func unloadNativeModel(modelID string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()

	if _, ok := nativeModels[modelID]; !ok {
		return nil
	}
	freeNativeHandles(modelID)
	return nil
}

func decodeNative(ctx *C.struct_llama_context, tokens []C.llama_token) error {
	if len(tokens) == 0 {
		return nil
	}
	batch := C.llama_batch_get_one(&tokens[0], C.int32_t(len(tokens)))
	rc := C.llama_decode(ctx, batch)
	if rc != 0 {
		return ErrGenerationFailed
	}
	return nil
}

func tokenizeNative(vocab *C.struct_llama_vocab, prompt string) ([]C.llama_token, error) {
	cPrompt := C.CString(prompt)
	defer C.free(unsafe.Pointer(cPrompt))

	textLen := C.int32_t(len(prompt))
	nMax := C.int32_t(256)
	for {
		buf := make([]C.llama_token, int(nMax))
		n := C.llama_tokenize(vocab, cPrompt, textLen, &buf[0], nMax, C.bool(false), C.bool(false))
		if n >= 0 {
			return buf[:int(n)], nil
		}
		if int32(n) == -2147483648 {
			return nil, ErrGenerationFailed
		}
		needed := C.int32_t(-n)
		if needed <= nMax {
			return nil, ErrGenerationFailed
		}
		nMax = needed
	}
}

func tokenToPieceNative(vocab *C.struct_llama_vocab, token C.llama_token) (string, error) {
	buf := make([]byte, 64)
	n := C.llama_token_to_piece(
		vocab,
		token,
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.int32_t(len(buf)),
		0,
		C.bool(false),
	)
	if n < 0 {
		size := -int(n)
		buf = make([]byte, size)
		n = C.llama_token_to_piece(
			vocab,
			token,
			(*C.char)(unsafe.Pointer(&buf[0])),
			C.int32_t(len(buf)),
			0,
			C.bool(false),
		)
		if n < 0 {
			return "", ErrGenerationFailed
		}
	}
	return string(buf[:n]), nil
}

func newNativeSampler(cfg LlamaConfig) *C.struct_llama_sampler {
	sparams := C.llama_sampler_chain_default_params()
	chain := C.llama_sampler_chain_init(sparams)
	if chain == nil {
		return nil
	}
	if cfg.GreedySampler {
		greedy := C.llama_sampler_init_greedy()
		C.llama_sampler_chain_add(chain, greedy)
		return chain
	}
	temp := C.llama_sampler_init_temp(C.float(cfg.Temperature))
	C.llama_sampler_chain_add(chain, temp)
	seed := cfg.Seed
	if seed == 0 {
		seed = 0xFFFFFFFF
	}
	dist := C.llama_sampler_init_dist(C.uint32_t(seed))
	C.llama_sampler_chain_add(chain, dist)
	return chain
}

// generateStats is defined in stats.go (no build constraints) so it is
// available to both CGO and Windows dynamic DLL backends.


func generateNative(modelID string, prompt string, cfg LlamaConfig) (string, uint32, error) {
	out, n, _, err := generateNativeTimed(modelID, prompt, cfg)
	return out, n, err
}

func generateNativeTimed(modelID string, prompt string, cfg LlamaConfig) (string, uint32, generateStats, error) {
	nativeMu.Lock()
	defer nativeMu.Unlock()

	modelPtr, ok := nativeModels[modelID]
	if !ok {
		return "", 0, generateStats{}, fmt.Errorf("native model not loaded")
	}
	ctxPtr, ok := nativeContexts[modelID]
	if !ok {
		return "", 0, generateStats{}, fmt.Errorf("native context not loaded")
	}

	model := (*C.struct_llama_model)(modelPtr)
	ctx := (*C.struct_llama_context)(ctxPtr)

	mem := C.llama_get_memory(ctx)
	C.llama_memory_clear(mem, C.bool(true))

	vocab := C.llama_model_get_vocab(model)
	tokens, err := tokenizeNative(vocab, prompt)
	if err != nil {
		return "", 0, generateStats{}, err
	}

	batchSize := int(effectiveBatchSize(cfg))
	t0 := C.llama_time_us()

	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		if err := decodeNative(ctx, tokens[i:end]); err != nil {
			return "", 0, generateStats{}, err
		}
	}

	sampler := newNativeSampler(cfg)
	if sampler == nil {
		return "", 0, generateStats{}, ErrGenerationFailed
	}
	defer C.llama_sampler_free(sampler)

	maxTokens := cfg.MaxTokens
	if maxTokens == 0 {
		maxTokens = 128
	}

	var stats generateStats
	var output strings.Builder
	var produced uint32
	firstToken := true

	for produced < maxTokens {
		token := C.llama_sampler_sample(sampler, ctx, -1)
		if C.llama_vocab_is_eog(vocab, token) {
			break
		}

		if firstToken {
			stats.TTFTMicros = int64(C.llama_time_us() - t0)
			firstToken = false
		}

		piece, err := tokenToPieceNative(vocab, token)
		if err != nil {
			return "", 0, generateStats{}, err
		}
		output.WriteString(piece)
		produced++

		var single [1]C.llama_token
		single[0] = token
		if err := decodeNative(ctx, single[:]); err != nil {
			return "", 0, generateStats{}, err
		}
	}

	if !firstToken {
		stats.DecodeMicros = int64(C.llama_time_us()-t0) - stats.TTFTMicros
	}

	return output.String(), produced, stats, nil
}
