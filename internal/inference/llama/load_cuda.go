//go:build cgo && cuda && !rocm && !metal && !vulkan

package llama

// loadBackend loads a GGUF model and inference context with the CUDA backend.
func (p *LlamaProvider) loadBackend(path string, modelID string) error {
	return loadNativeModel(path, modelID, p.cfg)
}

// unloadBackend frees the native context and model handles.
func (p *LlamaProvider) unloadBackend(modelID string) error {
	return unloadNativeModel(modelID)
}
