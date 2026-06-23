//go:build cgo && cuda && !rocm && !metal && !vulkan

package llama

// generateBackendTimed is kept for backward-compat with cuda_integration_test.go.
// On Windows the actual GPU execution routes through the dynamic DLL backend
// (backend_windows_dynamic.go) selected in provider_windows.go.
func (p *LlamaProvider) generateBackendTimed(modelID string, prompt string) (string, uint32, generateStats, error) {
	return p.b.generate(modelID, prompt, p.cfg)
}
