//go:build cgo && !cuda && !rocm && !metal && !vulkan

package llama

/*
#cgo CFLAGS: -I${SRCDIR}/../../../third_party/llama.cpp/include -I${SRCDIR}/../../../third_party/llama.cpp/ggml/include
#cgo linux LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lstdc++ -lm -lpthread -ldl -lgomp
#cgo darwin LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lc++ -lm -lpthread
#cgo windows LDFLAGS: -L${SRCDIR}/../../../third_party/llama.cpp/build/bin -lllama -lggml -lggml-cpu -lggml-base -lstdc++ -lm
#include "llama.h"
#include <stdlib.h>
*/
import "C"

import (
	"sync"
	"unsafe"
)

var (
	backendOnce  sync.Once
	nativeMu     sync.Mutex
	nativeModels = make(map[string]unsafe.Pointer)
)

func initBackend() {
	backendOnce.Do(func() {
		C.llama_backend_init()
	})
}

func loadNativeModel(path string, modelID string, useMMAP bool) error {
	initBackend()

	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	params := C.llama_model_default_params()
	if useMMAP {
		params.use_mmap = C.bool(true)
	} else {
		params.use_mmap = C.bool(false)
	}
	params.n_gpu_layers = 0

	model := C.llama_model_load_from_file(cPath, params)
	if model == nil {
		return ErrPathInaccessible
	}

	nativeMu.Lock()
	defer nativeMu.Unlock()
	if old, ok := nativeModels[modelID]; ok {
		C.llama_model_free((*C.struct_llama_model)(old))
	}
	nativeModels[modelID] = unsafe.Pointer(model)
	return nil
}

func unloadNativeModel(modelID string) error {
	nativeMu.Lock()
	defer nativeMu.Unlock()

	model, ok := nativeModels[modelID]
	if !ok {
		return nil
	}
	C.llama_model_free((*C.struct_llama_model)(model))
	delete(nativeModels, modelID)
	return nil
}
