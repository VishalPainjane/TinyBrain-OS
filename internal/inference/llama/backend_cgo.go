//go:build cgo

package llama

// cgoBackend routes all inference calls through static CGO bindings.
// On Linux/macOS this is the only backend.  On Windows it is used for the CPU
// path and for any platform where ggml-cuda.dll is not present.
//
// The actual native calls live in bindings_common.go (shared across CPU and
// CUDA static builds) and bindings_cpu.go / bindings_cuda.go (link flags).
type cgoBackend struct{}

func (b *cgoBackend) loadModel(path, modelID string, cfg LlamaConfig) error {
	return loadNativeModel(path, modelID, cfg)
}

func (b *cgoBackend) unloadModel(modelID string) error {
	return unloadNativeModel(modelID)
}

func (b *cgoBackend) generate(modelID, prompt string, cfg LlamaConfig) (string, uint32, generateStats, error) {
	return generateNativeTimed(modelID, prompt, cfg)
}
