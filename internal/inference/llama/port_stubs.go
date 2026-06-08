package llama

import "github.com/VishalPainjane/TinyBrain-OS/internal/runtime"

// Generate is not implemented in 009a.
func (p *LlamaProvider) Generate(_ runtime.GenerateRequest) (runtime.GenerateResponse, error) {
	return runtime.GenerateResponse{}, ErrNotImplemented
}

// SaveContext is not implemented in 009a.
func (p *LlamaProvider) SaveContext(_ string) error {
	return ErrNotImplemented
}

// RestoreContext is not implemented in 009a.
func (p *LlamaProvider) RestoreContext(_ string) error {
	return ErrNotImplemented
}
