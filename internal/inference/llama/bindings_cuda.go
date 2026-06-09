//go:build cgo && cuda && !rocm && !metal && !vulkan

package llama

/*
#cgo linux LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build-cuda/bin -lllama -lggml -lggml-cuda -lggml-cpu -lggml-base -lstdc++ -lm -lpthread -ldl -lgomp -lcudart -lcublas
#cgo windows LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build-cuda/bin -lllama -lggml -lggml-cuda -lggml-cpu -lggml-base -lstdc++ -lm -lcudart -lcublas
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// CUDA linkage uses third_party/llama.cpp/build-cuda/ (cmake -DGGML_CUDA=ON).
// Library names verified for llama.cpp pin b9553 @ 9e3b928: llama, ggml, ggml-cuda, ggml-cpu, ggml-base.

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
	params.n_gpu_layers = C.int32_t(EffectiveNGLayers(cfg.NGLayers))
	params.main_gpu = 0

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
