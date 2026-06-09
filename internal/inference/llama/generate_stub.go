//go:build !cgo

package llama

// generateBackend is unavailable when CGO is disabled.
func (p *LlamaProvider) generateBackend(modelID string, prompt string) (string, uint32, error) {
	return "", 0, ErrCGODisabled
}
