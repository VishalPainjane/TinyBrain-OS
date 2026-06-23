//go:build cgo && !cuda && !rocm && !metal && !vulkan

package llama

// generateBackendTimed is kept for backward-compat with cuda_integration_test.go
// and any callers that need timing data directly.  New callers should use
// LlamaProvider.Generate which routes through the backend interface.
func (p *LlamaProvider) generateBackendTimed(modelID string, prompt string) (string, uint32, generateStats, error) {
	return p.b.generate(modelID, prompt, p.cfg)
}
