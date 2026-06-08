//go:build !cgo

package llama

// loadBackend is a no-op when CGO is disabled.
func (p *LlamaProvider) loadBackend(_ string, _ string) error {
	return ErrCGODisabled
}

// unloadBackend is a no-op when CGO is disabled.
func (p *LlamaProvider) unloadBackend(_ string) error {
	return ErrCGODisabled
}
