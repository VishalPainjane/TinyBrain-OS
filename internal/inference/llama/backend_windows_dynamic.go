//go:build windows && !cgo_cuda_static

// Package llama — Windows GPU dynamic backend.
//
// # Architecture
//
// On Windows, NVIDIA CUDA libraries are compiled with MSVC (cl.exe), but Go's
// CGO toolchain on Windows uses MinGW-GCC.  These two ABI families cannot be
// statically linked into the same binary.
//
// Instead, this file implements the backend interface using syscall.LoadDLL /
// syscall.NewProc to late-bind the MSVC-compiled llama.dll at runtime.  The
// call convention is the standard Windows C calling convention (cdecl / ms-abi)
// exported by llama.cpp's public C API (llama.h).
//
// # Backend selection (runtime probe)
//
// NewWindowsDynamicBackend probes for ggml-cuda.dll alongside the binary.  If
// the DLL is not found the function returns (nil, ErrDLLNotFound) and the
// caller (provider_windows.go) falls back to the static CGO CPU backend.
//
// # Struct layout
//
// llama.cpp's public C API passes structs by value.  We replicate the exact
// memory layout here so that uintptr-level marshalling is correct.  The layouts
// are pinned to llama.cpp commit 9e3b928 (build tag: b9553).  If the submodule
// is updated these structs must be re-verified.
//
// # Future backends
//
// This file establishes the pattern for Metal (macOS), ROCm (Linux/Windows),
// and Vulkan.  Each backend builds a backend implementation that probes for its
// shared library and exposes the same interface.
package llama

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	osruntime "github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

// ErrDLLNotFound is returned when ggml-cuda.dll or llama.dll cannot be located.
var ErrDLLNotFound = errors.New("llama CUDA DLL not found")

// ErrDLLSymbol is returned when a required exported symbol is missing.
var ErrDLLSymbol = errors.New("llama DLL missing required export")

// llamaDLLBackend implements backend via runtime-loaded Windows DLLs.
// All fields are set once during NewWindowsDynamicBackend and then read-only.
type llamaDLLBackend struct {
	llamaDLL    *syscall.DLL
	ggmlDLL     *syscall.DLL // ggml-base.dll — loaded first so llama.dll resolves it
	ggmlCPUDLL  *syscall.DLL // ggml-cpu.dll
	ggmlCUDADLL *syscall.DLL // ggml-cuda.dll (optional)
	ggmlCoreDLL *syscall.DLL // ggml.dll

	// model lifecycle
	procModelDefaultParams *syscall.Proc // llama_model_default_params
	procModelLoadFromFile  *syscall.Proc // llama_model_load_from_file
	procModelFree          *syscall.Proc // llama_model_free
	procModelGetVocab      *syscall.Proc // llama_model_get_vocab

	// context lifecycle
	procContextDefaultParams *syscall.Proc // llama_context_default_params
	procInitFromModel        *syscall.Proc // llama_init_from_model
	procFree                 *syscall.Proc // llama_free
	procGetMemory            *syscall.Proc // llama_get_memory
	procMemoryClear          *syscall.Proc // llama_memory_clear

	// backend init
	procBackendInit *syscall.Proc // llama_backend_init
	procTimeUs      *syscall.Proc // llama_time_us

	// tokenisation
	procTokenize     *syscall.Proc // llama_tokenize
	procTokenToPiece *syscall.Proc // llama_token_to_piece
	procVocabIsEog   *syscall.Proc // llama_vocab_is_eog

	// batched decode
	procBatchGetOne *syscall.Proc // llama_batch_get_one
	procDecode      *syscall.Proc // llama_decode

	// sampler chain
	procSamplerChainDefaultParams *syscall.Proc // llama_sampler_chain_default_params
	procSamplerChainInit          *syscall.Proc // llama_sampler_chain_init
	procSamplerChainAdd           *syscall.Proc // llama_sampler_chain_add
	procSamplerInitGreedy         *syscall.Proc // llama_sampler_init_greedy
	procSamplerInitTemp           *syscall.Proc // llama_sampler_init_temp
	procSamplerInitDist           *syscall.Proc // llama_sampler_init_dist
	procSamplerInitGrammar        *syscall.Proc // llama_sampler_init_grammar
	procSamplerSample             *syscall.Proc // llama_sampler_sample
	procSamplerFree               *syscall.Proc // llama_sampler_free

	procStateGetSize *syscall.Proc // llama_state_get_size
	procStateGetData *syscall.Proc // llama_state_get_data
	procStateSetData *syscall.Proc // llama_state_set_data

	// templates & metadata
	procChatApplyTemplate *syscall.Proc // llama_chat_apply_template
	procModelMetaValStr   *syscall.Proc // llama_model_meta_val_str

	// state
	backendOnce sync.Once
	mu          sync.Mutex
	models      map[string]uintptr // modelID → llama_model*
	contexts    map[string]uintptr // modelID → llama_context*
	stateBlobs  map[string][]byte  // ctxID → []byte
}

// ─────────────────────────────────────────────────────────────────────────────
// Struct layouts (pinned to llama.cpp b9553 / 9e3b928)
// ─────────────────────────────────────────────────────────────────────────────
//
// These Go structs must match the C struct ABI exactly.  Fields are ordered as
// declared in llama.h and ggml.h.  Padding is explicit via _ fields.

// llamaModelParams mirrors llama_model_params (llama.h).
type llamaModelParams struct {
	devices             uintptr // ggml_backend_dev_t *
	tensorBuftOverrides uintptr // const struct llama_model_tensor_buft_override *
	nGPULayers          int32
	splitMode           int32 // llama_split_mode enum (int32)
	mainGPU             int32
	_                   int32 // pad to 8-byte align

	tensorSplit  uintptr // const float*
	progressCB   uintptr // llama_progress_callback
	progressData uintptr // void*

	kvOverrides uintptr // const llama_model_kv_override*

	vocabOnly     uint8 // bool
	useMMAP       uint8 // bool
	useDirectIO   uint8 // bool
	useMLOCK      uint8 // bool
	checkTensors  uint8 // bool
	useExtraBufts uint8 // bool
	noHost        uint8 // bool
	noAlloc       uint8 // bool
}

// llamaChatMessage mirrors llama_chat_message (llama.h).
type llamaChatMessage struct {
	role    uintptr // const char *
	content uintptr // const char *
}

// llamaContextParams mirrors llama_context_params (llama.h).
type llamaContextParams struct {
	nCtx            uint32
	nBatch          uint32
	nUBatch         uint32
	nSeqMax         uint32
	nRsSeq          uint32
	nOutputsMax     uint32
	nThreads        int32
	nThreadsBatch   int32
	ctxType         int32 // llama_context_type
	ropeScalingType int32 // llama_rope_scaling_type
	poolingType     int32 // llama_pooling_type
	attentionType   int32 // llama_attention_type
	flashAttnType   int32 // llama_flash_attn_type
	ropeFreqBase    float32
	ropeFreqScale   float32
	yarnExtFactor   float32
	yarnAttnFactor  float32
	yarnBetaFast    float32
	yarnBetaSlow    float32
	yarnOrigCtx     uint32
	defragThold     float32
	_               int32 // pad to 8-byte align

	cbEval            uintptr // ggml_backend_sched_eval_callback
	cbEvalUserData    uintptr // void *
	typeK             uint32  // ggml_type
	typeV             uint32  // ggml_type
	abortCallback     uintptr // ggml_abort_callback
	abortCallbackData uintptr // void *

	embeddings uint8 // bool
	offloadKQV uint8 // bool
	noPerf     uint8 // bool
	opOffload  uint8 // bool
	swaFull    uint8 // bool
	kvUnified  uint8 // bool
	_          [2]byte // pad to 8-byte align

	samplers  uintptr // struct llama_sampler_seq_config *
	nSamplers uintptr // size_t
	ctxOther  uintptr // struct llama_context *
}

// llamaSamplerChainParams mirrors llama_sampler_chain_params (llama.h).
type llamaSamplerChainParams struct {
	noPerf uint8
}

// ─────────────────────────────────────────────────────────────────────────────
// Constructor
// ─────────────────────────────────────────────────────────────────────────────

// NewWindowsDynamicBackend attempts to locate and load the MSVC-compiled
// llama.dll and its dependencies from dllDir.  dllDir is typically the
// directory alongside the running binary.
//
// Returns (nil, ErrDLLNotFound) if the DLLs are absent (caller should fall back
// to the CGO CPU backend).
func NewWindowsDynamicBackend(dllDir string) (*llamaDLLBackend, error) {
	llamaPath := filepath.Join(dllDir, "llama.dll")
	if _, err := os.Stat(llamaPath); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrDLLNotFound, llamaPath)
	}

	var ggmlDLL, ggmlCPUDLL, ggmlCUDADLL, ggmlCoreDLL, llamaDLL *syscall.DLL
	var err error

	rollback := func() {
		if llamaDLL != nil {
			llamaDLL.Release()
		}
		if ggmlCoreDLL != nil {
			ggmlCoreDLL.Release()
		}
		if ggmlCUDADLL != nil {
			ggmlCUDADLL.Release()
		}
		if ggmlCPUDLL != nil {
			ggmlCPUDLL.Release()
		}
		if ggmlDLL != nil {
			ggmlDLL.Release()
		}
	}

	// Load ggml-base.dll first so that llama.dll's import resolves
	ggmlBasePath := filepath.Join(dllDir, "ggml-base.dll")
	ggmlDLL, err = syscall.LoadDLL(ggmlBasePath)
	if err != nil {
		return nil, fmt.Errorf("load ggml-base.dll: %w", err)
	}

	// Load ggml-cpu.dll
	ggmlCPUPath := filepath.Join(dllDir, "ggml-cpu.dll")
	ggmlCPUDLL, err = syscall.LoadDLL(ggmlCPUPath)
	if err != nil {
		rollback()
		return nil, fmt.Errorf("load ggml-cpu.dll: %w", err)
	}

	// Load optional ggml-cuda.dll (may not exist on CPU-only deployments)
	cudaPath := filepath.Join(dllDir, "ggml-cuda.dll")
	if _, statErr := os.Stat(cudaPath); statErr == nil {
		// Load into process address space so llama.dll can resolve CUDA symbols.
		ggmlCUDADLL, err = syscall.LoadDLL(cudaPath)
		if err != nil {
			// Non-fatal: llama.dll will fall back to CPU backend inside llama.cpp.
			// Log but continue.
			_ = fmt.Sprintf("warning: ggml-cuda.dll found but failed to load: %v", err)
		}
	}

	// Load ggml.dll
	ggmlCorePath := filepath.Join(dllDir, "ggml.dll")
	ggmlCoreDLL, err = syscall.LoadDLL(ggmlCorePath)
	if err != nil {
		rollback()
		return nil, fmt.Errorf("load ggml.dll: %w", err)
	}

	// Load llama.dll
	llamaDLL, err = syscall.LoadDLL(llamaPath)
	if err != nil {
		rollback()
		return nil, fmt.Errorf("load llama.dll: %w", err)
	}

	b := &llamaDLLBackend{
		llamaDLL:    llamaDLL,
		ggmlDLL:     ggmlDLL,
		ggmlCPUDLL:  ggmlCPUDLL,
		ggmlCUDADLL: ggmlCUDADLL,
		ggmlCoreDLL: ggmlCoreDLL,
		models:      make(map[string]uintptr),
		contexts:    make(map[string]uintptr),
		stateBlobs:  make(map[string][]byte),
	}

	// Resolve all required procedure addresses.
	if err := b.resolveProcs(); err != nil {
		rollback()
		return nil, err
	}

	return b, nil
}

// resolveProcs fetches all required exported symbols from the DLLs.
func (b *llamaDLLBackend) resolveProcs() error {
	type binding struct {
		dest *(*syscall.Proc)
		dll  *syscall.DLL
		name string
	}
	ll := b.llamaDLL

	bindings := []binding{
		{&b.procModelDefaultParams, ll, "llama_model_default_params"},
		{&b.procModelLoadFromFile, ll, "llama_model_load_from_file"},
		{&b.procModelFree, ll, "llama_model_free"},
		{&b.procModelGetVocab, ll, "llama_model_get_vocab"},
		{&b.procContextDefaultParams, ll, "llama_context_default_params"},
		{&b.procInitFromModel, ll, "llama_init_from_model"},
		{&b.procFree, ll, "llama_free"},
		{&b.procGetMemory, ll, "llama_get_memory"},
		{&b.procMemoryClear, ll, "llama_memory_clear"},
		{&b.procBackendInit, ll, "llama_backend_init"},
		{&b.procTimeUs, ll, "llama_time_us"},
		{&b.procTokenize, ll, "llama_tokenize"},
		{&b.procTokenToPiece, ll, "llama_token_to_piece"},
		{&b.procVocabIsEog, ll, "llama_vocab_is_eog"},
		{&b.procBatchGetOne, ll, "llama_batch_get_one"},
		{&b.procDecode, ll, "llama_decode"},
		{&b.procSamplerChainDefaultParams, ll, "llama_sampler_chain_default_params"},
		{&b.procSamplerChainInit, ll, "llama_sampler_chain_init"},
		{&b.procSamplerChainAdd, ll, "llama_sampler_chain_add"},
		{&b.procSamplerInitGreedy, ll, "llama_sampler_init_greedy"},
		{&b.procSamplerInitTemp, ll, "llama_sampler_init_temp"},
		{&b.procSamplerInitDist, ll, "llama_sampler_init_dist"},
		{&b.procSamplerSample, ll, "llama_sampler_sample"},
		{&b.procSamplerFree, ll, "llama_sampler_free"},
		{&b.procStateGetSize, ll, "llama_state_get_size"},
		{&b.procStateGetData, ll, "llama_state_get_data"},
		{&b.procStateSetData, ll, "llama_state_set_data"},
		{&b.procChatApplyTemplate, ll, "llama_chat_apply_template"},
		{&b.procModelMetaValStr, ll, "llama_model_meta_val_str"},
	}

	for _, bnd := range bindings {
		proc, err := bnd.dll.FindProc(bnd.name)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrDLLSymbol, bnd.name)
		}
		*bnd.dest = proc
	}

	// Optional bindings
	if proc, err := ll.FindProc("llama_sampler_init_grammar"); err == nil {
		b.procSamplerInitGrammar = proc
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// backend interface implementation
// ─────────────────────────────────────────────────────────────────────────────

func (b *llamaDLLBackend) loadModel(path, modelID string, cfg LlamaConfig) error {
	b.backendOnce.Do(func() { b.procBackendInit.Call() })

	params := b.defaultModelParams(cfg)
	cPath, err := syscall.BytePtrFromString(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// llama_model_load_from_file(path, params) → llama_model*
	modelPtr, _, _ := b.procModelLoadFromFile.Call(
		uintptr(unsafe.Pointer(cPath)),
		uintptr(unsafe.Pointer(&params)),
	)
	if modelPtr == 0 {
		return ErrPathInaccessible
	}

	ctxParams := b.defaultContextParams(cfg)
	ctxPtr, _, _ := b.procInitFromModel.Call(
		modelPtr,
		uintptr(unsafe.Pointer(&ctxParams)),
	)
	if ctxPtr == 0 {
		b.procModelFree.Call(modelPtr)
		return fmt.Errorf("%w: context init failed", ErrGenerationFailed)
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	b.freeHandles(modelID) // idempotent cleanup
	b.models[modelID] = modelPtr
	b.contexts[modelID] = ctxPtr
	return nil
}

func (b *llamaDLLBackend) unloadModel(modelID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.models[modelID]; !ok {
		return nil
	}
	b.freeHandles(modelID)
	return nil
}

// freeHandles releases native resources.  Caller must hold b.mu.
func (b *llamaDLLBackend) freeHandles(modelID string) {
	if ctx, ok := b.contexts[modelID]; ok {
		b.procFree.Call(ctx)
		delete(b.contexts, modelID)
	}
	if model, ok := b.models[modelID]; ok {
		b.procModelFree.Call(model)
		delete(b.models, modelID)
	}
}

func (b *llamaDLLBackend) generate(req osruntime.GenerateRequest, cfg LlamaConfig) (string, uint32, generateStats, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	modelID := req.ModelID
	prompt := req.Prompt

	modelPtr, ok := b.models[modelID]
	if !ok {
		return "", 0, generateStats{}, fmt.Errorf("dynamic: native model not loaded")
	}
	ctxPtr, ok := b.contexts[modelID]
	if !ok {
		return "", 0, generateStats{}, fmt.Errorf("dynamic: native context not loaded")
	}

	// Clear KV cache
	memPtr, _, _ := b.procGetMemory.Call(ctxPtr)
	b.procMemoryClear.Call(memPtr, 1) // true = reset

	// Tokenise
	vocabPtr, _, _ := b.procModelGetVocab.Call(modelPtr)
	tokens, err := b.tokenize(vocabPtr, prompt)
	if err != nil {
		return "", 0, generateStats{}, err
	}

	batchSize := int(effectiveBatchSizeVal(cfg))
	t0R, _, _ := b.procTimeUs.Call()
	t0 := int64(t0R)

	// Prefill batches
	for i := 0; i < len(tokens); i += batchSize {
		end := i + batchSize
		if end > len(tokens) {
			end = len(tokens)
		}
		if decErr := b.decodeBatch(ctxPtr, tokens[i:end]); decErr != nil {
			return "", 0, generateStats{}, decErr
		}
	}

	// Build sampler
	samplerPtr := b.buildSampler(cfg, modelPtr, req.Grammar)
	if samplerPtr == 0 {
		return "", 0, generateStats{}, ErrGenerationFailed
	}
	defer b.procSamplerFree.Call(samplerPtr)

	maxTokens := cfg.MaxTokens
	if req.MaxTokens > 0 {
		maxTokens = uint32(req.MaxTokens)
	}
	if maxTokens == 0 {
		maxTokens = 128
	}

	var stats generateStats
	var out strings.Builder
	var produced uint32
	firstToken := true

	for produced < maxTokens {
		tokenR, _, _ := b.procSamplerSample.Call(samplerPtr, ctxPtr, uintptr(^uint32(0))) // -1 as uint
		token := int32(tokenR)

		// Check EOG
		eogR, _, _ := b.procVocabIsEog.Call(vocabPtr, uintptr(uint32(token)))
		isEog := byte(eogR) != 0
		if isEog {
			break
		}

		if firstToken {
			tNowR, _, _ := b.procTimeUs.Call()
			stats.TTFTMicros = int64(tNowR) - t0
			firstToken = false
		}

		piece, pieceErr := b.tokenToPiece(vocabPtr, token)
		if pieceErr != nil {
			return "", 0, generateStats{}, pieceErr
		}
		out.WriteString(piece)
		produced++

		single := [1]int32{token}
		if decErr := b.decodeBatch(ctxPtr, single[:]); decErr != nil {
			return "", 0, generateStats{}, decErr
		}
	}

	if !firstToken {
		tNowR, _, _ := b.procTimeUs.Call()
		stats.DecodeMicros = int64(tNowR) - t0 - stats.TTFTMicros
	}

	// Suppress unused import linting
	_ = time.Now
	_ = runtime.GOOS

	return out.String(), produced, stats, nil
}

func (b *llamaDLLBackend) getMetadata(modelID string) (osruntime.ModelCapabilities, error) {
	b.mu.Lock()
	modelPtr, ok := b.models[modelID]
	b.mu.Unlock()

	if !ok {
		return osruntime.ModelCapabilities{}, fmt.Errorf("dynamic: native model not loaded")
	}

	key, _ := syscall.BytePtrFromString("tokenizer.chat_template")
	buf := make([]byte, 32768)
	nR, _, _ := b.procModelMetaValStr.Call(
		modelPtr,
		uintptr(unsafe.Pointer(key)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	
	n := int32(nR)
	template := ""
	if n > 0 && n <= int32(len(buf)) {
		if buf[n-1] == 0 {
			n--
		}
		template = string(buf[:n])
	}
	
	return osruntime.ModelCapabilities{
		ModelID:            modelID,
		ChatTemplate:       template,
		SupportsMultimodal: false,
		SupportsTools:      false,
		SupportsGrammar:    false, 
	}, nil
}

func (b *llamaDLLBackend) formatChat(modelID string, messages []osruntime.ChatMessage, opts osruntime.FormatChatOpts) (string, string, error) {
	caps, err := b.getMetadata(modelID)
	if err != nil {
		return "", "", err
	}

	tmplStr := opts.TemplateName
	if tmplStr == "" {
		tmplStr = caps.ChatTemplate
	}
	
	var cChat []llamaChatMessage
	var pinStr [](*byte)
	
	for _, m := range messages {
		cRole, _ := syscall.BytePtrFromString(m.Role)
		cContent, _ := syscall.BytePtrFromString(m.Content)
		pinStr = append(pinStr, cRole, cContent)
		
		cChat = append(cChat, llamaChatMessage{
			role:    uintptr(unsafe.Pointer(cRole)),
			content: uintptr(unsafe.Pointer(cContent)),
		})
	}

	addAss := uintptr(0)
	if opts.AddGenerationPrompt {
		addAss = 1
	}

	var chatPtr uintptr
	if len(cChat) > 0 {
		chatPtr = uintptr(unsafe.Pointer(&cChat[0]))
	}

	n := int32(-1)
	if tmplStr != "" {
		cTmpl, err := syscall.BytePtrFromString(tmplStr)
		if err == nil {
			buf := make([]byte, 65536) // 64KB
			nR, _, _ := b.procChatApplyTemplate.Call(
				uintptr(unsafe.Pointer(cTmpl)),
				chatPtr,
				uintptr(len(cChat)),
				addAss,
				uintptr(unsafe.Pointer(&buf[0])),
				uintptr(len(buf)),
			)
			n = int32(nR)
			if n >= 0 && n <= int32(len(buf)) {
				if n > 0 && buf[n-1] == 0 {
					n--
				}
				runtime.KeepAlive(pinStr)
				return string(buf[:n]), tmplStr, nil
			}
		}
	}
	
	runtime.KeepAlive(pinStr)

	// Fallback mechanism if minja fails or template is missing
	fallbackTmpl := ""
	for _, m := range messages {
		fallbackTmpl += fmt.Sprintf("<|im_start|>%s\n%s<|im_end|>\n", m.Role, m.Content)
	}
	if opts.AddGenerationPrompt {
		fallbackTmpl += "<|im_start|>assistant\n"
	}
	return fallbackTmpl, "fallback_chatml", nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal helpers
// ─────────────────────────────────────────────────────────────────────────────

func (b *llamaDLLBackend) defaultModelParams(cfg LlamaConfig) llamaModelParams {
	// Call the DLL to get the defaults struct, then patch fields we care about.
	var p llamaModelParams
	b.procModelDefaultParams.Call(uintptr(unsafe.Pointer(&p)))
	if cfg.UseMMAP {
		p.useMMAP = 1
	} else {
		p.useMMAP = 0
	}
	p.nGPULayers = EffectiveNGLayers(cfg.NGLayers)
	return p
}

func (b *llamaDLLBackend) defaultContextParams(cfg LlamaConfig) llamaContextParams {
	var p llamaContextParams
	b.procContextDefaultParams.Call(uintptr(unsafe.Pointer(&p)))
	p.nCtx = cfg.ContextSize
	p.nThreads = int32(cfg.Threads)
	p.nBatch = effectiveBatchSizeVal(cfg)
	return p
}

func effectiveBatchSizeVal(cfg LlamaConfig) uint32 {
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

// tokenize calls llama_tokenize via the DLL.  Returns a slice of int32 token IDs.
func (b *llamaDLLBackend) tokenize(vocabPtr uintptr, prompt string) ([]int32, error) {
	cPrompt, err := syscall.BytePtrFromString(prompt)
	if err != nil {
		return nil, fmt.Errorf("tokenize: invalid prompt: %w", err)
	}

	nMax := int32(256)
	for {
		buf := make([]int32, int(nMax))
		nR, _, _ := b.procTokenize.Call(
			vocabPtr,
			uintptr(unsafe.Pointer(cPrompt)),
			uintptr(len(prompt)),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(nMax),
			1, // add_special = true  (prepend BOS, matching llama.cpp CLI / HF tokenizers)
			0, // parse_special = false
		)
		n := int32(nR)
		if n >= 0 {
			return buf[:n], nil
		}
		if n == -2147483648 {
			return nil, ErrGenerationFailed
		}
		needed := -n
		if needed <= nMax {
			return nil, ErrGenerationFailed
		}
		nMax = needed
	}
}

// tokenToPiece converts a single token ID to its string piece.
func (b *llamaDLLBackend) tokenToPiece(vocabPtr uintptr, token int32) (string, error) {
	buf := make([]byte, 64)
	nR, _, _ := b.procTokenToPiece.Call(
		vocabPtr,
		uintptr(uint32(token)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
		0,  // lstrip = 0
		0,  // special = false
	)
	n := int32(nR)
	if n < 0 {
		size := int(-n)
		buf = make([]byte, size)
		nR, _, _ = b.procTokenToPiece.Call(
			vocabPtr,
			uintptr(uint32(token)),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
			0,
			0,
		)
		n = int32(nR)
		if n < 0 {
			return "", ErrGenerationFailed
		}
	}
	return string(buf[:n]), nil
}

func (b *llamaDLLBackend) saveContext(modelID, ctxID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	ctxPtr, ok := b.contexts[modelID]
	if !ok {
		return fmt.Errorf("model not loaded: %s", modelID)
	}

	size, _, _ := b.procStateGetSize.Call(ctxPtr)
	if size == 0 {
		return fmt.Errorf("state size is 0")
	}

	stateData := make([]byte, size)
	nBytes, _, _ := b.procStateGetData.Call(ctxPtr, uintptr(unsafe.Pointer(&stateData[0])), size)
	if nBytes != size {
		return fmt.Errorf("failed to get full state, expected %d got %d", size, nBytes)
	}

	b.stateBlobs[ctxID] = stateData
	return nil
}

func (b *llamaDLLBackend) restoreContext(modelID, ctxID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	stateData, ok := b.stateBlobs[ctxID]
	if !ok {
		return fmt.Errorf("context not found: %s", ctxID)
	}

	ctxPtr, ok := b.contexts[modelID]
	if !ok {
		return fmt.Errorf("model not loaded: %s", modelID)
	}

	size := uintptr(len(stateData))
	nBytes, _, _ := b.procStateSetData.Call(ctxPtr, uintptr(unsafe.Pointer(&stateData[0])), size)
	if nBytes != size {
		return fmt.Errorf("failed to set full state, expected %d got %d", size, nBytes)
	}

	return nil
}

// decodeBatch submits a batch of tokens through llama_decode.
func (b *llamaDLLBackend) decodeBatch(ctxPtr uintptr, tokens []int32) error {
	if len(tokens) == 0 {
		return nil
	}
	// llama_batch_get_one(tokens*, n_tokens) returns a llama_batch by value.
	// On Windows x64 ABI, a struct larger than 8 bytes is returned via a hidden
	// first pointer argument.  We allocate space for the batch struct and pass its address.
	// llama_batch is: token* data, float* embd, llama_pos* pos, int32_t* n_seq_id,
	// llama_seq_id** seq_id, int8_t* logits, int32_t n_tokens.
	// Total ≈ 7 pointers + 1 int32 = 60 bytes on x64.  Allocate 128 bytes for safety.
	var batchBuf [128]byte
	b.procBatchGetOne.Call(
		uintptr(unsafe.Pointer(&batchBuf[0])), // hidden return pointer
		uintptr(unsafe.Pointer(&tokens[0])),
		uintptr(len(tokens)),
	)
	rcR, _, _ := b.procDecode.Call(ctxPtr, uintptr(unsafe.Pointer(&batchBuf[0])))
	if int32(rcR) != 0 {
		return ErrGenerationFailed
	}
	return nil
}

// buildSampler creates a sampler chain according to LlamaConfig and Grammar constraints.
func (b *llamaDLLBackend) buildSampler(cfg LlamaConfig, modelPtr uintptr, grammarStr string) uintptr {
	// llama_sampler_chain_default_params returns a struct by value.
	// On Windows x64 ABI a struct return uses a hidden first-pointer argument,
	// identical to the llama_batch pattern in decodeBatch.  We allocate the
	// struct and pass its address as arg 0; the DLL writes the result there.
	var cp llamaSamplerChainParams
	b.procSamplerChainDefaultParams.Call(uintptr(unsafe.Pointer(&cp)))

	// llama_sampler_chain_init takes llamaSamplerChainParams by value (1 byte).
	// Pass as uintptr so the calling convention matches the C ABI.
	chainPtr, _, _ := b.procSamplerChainInit.Call(uintptr(cp.noPerf))
	if chainPtr == 0 {
		return 0
	}

	if grammarStr != "" && b.procSamplerInitGrammar != nil {
		cGrammar, _ := syscall.BytePtrFromString(grammarStr)
		cRoot, _ := syscall.BytePtrFromString("root")
		grammarPtr, _, _ := b.procSamplerInitGrammar.Call(
			modelPtr,
			uintptr(unsafe.Pointer(cGrammar)),
			uintptr(unsafe.Pointer(cRoot)),
		)
		if grammarPtr != 0 {
			b.procSamplerChainAdd.Call(chainPtr, grammarPtr)
		}
	}

	if cfg.GreedySampler {
		greedyPtr, _, _ := b.procSamplerInitGreedy.Call()
		b.procSamplerChainAdd.Call(chainPtr, greedyPtr)
		return chainPtr
	}

	tempBits := *(*uint32)(unsafe.Pointer(&cfg.Temperature)) // float32 → uint32 bits
	tempPtr, _, _ := b.procSamplerInitTemp.Call(uintptr(tempBits))
	b.procSamplerChainAdd.Call(chainPtr, tempPtr)

	seed := cfg.Seed
	if seed == 0 {
		seed = 0xFFFFFFFF
	}
	distPtr, _, _ := b.procSamplerInitDist.Call(uintptr(seed))
	b.procSamplerChainAdd.Call(chainPtr, distPtr)
	return chainPtr
}
