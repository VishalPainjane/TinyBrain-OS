//go:build cgo && !cuda && !rocm && !metal && !vulkan

package llama

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lstdc++ -lm -lpthread -ldl -lgomp
#cgo darwin LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lc++ -lm -lpthread
#cgo windows LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lstdc++ -lm
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func loadNativeModel(path string, modelID string, cfg LlamaConfig) error {
	initBackend()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	params := C.llama_model_default_params()
	if cfg.UseMMAP {
		params.use_mmap = C.bool(true)
	} else {
		params.use_mmap = C.bool(false)
	}
	// CPU binary always forces CPU-only load regardless of cfg.NGLayers.
	params.n_gpu_layers = 0

	model := C.llama_model_load_from_file(cPath, params)
	if model == nil {
		return ErrPathInaccessible
	}

	ctxParams := C.llama_context_default_params()
	ctxParams.n_ctx = C.uint32_t(cfg.ContextSize)
	ctxParams.n_threads = C.int32_t(cfg.Threads)
	ctxParams.n_batch = C.uint32_t(effectiveBatchSize(cfg))

	ctx := C.llama_init_from_model(model, ctxParams)
	if ctx == nil {
		C.llama_model_free(model)
		return fmt.Errorf("%w: context init failed", ErrGenerationFailed)
	}

	nativeMu.Lock()
	defer nativeMu.Unlock()

	freeNativeHandles(modelID)
	nativeModels[modelID] = unsafe.Pointer(model)
	nativeContexts[modelID] = unsafe.Pointer(ctx)
	return nil
}
