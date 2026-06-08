//go:build cgo && !cuda && !rocm && !metal && !vulkan

package llama

// loadBackend loads a GGUF model with the CPU backend (NGLayers=0).
func (p *LlamaProvider) loadBackend(path string, modelID string) error {
	return loadNativeModel(path, modelID, p.cfg.UseMMAP)
}

// unloadBackend frees the native model handle.
func (p *LlamaProvider) unloadBackend(modelID string) error {
	return unloadNativeModel(modelID)
}
