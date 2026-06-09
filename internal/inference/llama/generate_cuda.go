//go:build cgo && cuda && !rocm && !metal && !vulkan

package llama

// generateBackend runs CUDA inference for a loaded model.
func (p *LlamaProvider) generateBackend(modelID string, prompt string) (string, uint32, error) {
	return generateNative(modelID, prompt, p.cfg)
}

// generateBackendTimed runs CUDA inference and returns TTFT/decode timing in microseconds.
func (p *LlamaProvider) generateBackendTimed(modelID string, prompt string) (string, uint32, generateStats, error) {
	return generateNativeTimed(modelID, prompt, p.cfg)
}
