package llama

// SaveContext is not implemented in 009b (deferred to task 011).
func (p *LlamaProvider) SaveContext(_ string) error {
	return ErrNotImplemented
}

// RestoreContext is not implemented in 009b (deferred to task 011).
func (p *LlamaProvider) RestoreContext(_ string) error {
	return ErrNotImplemented
}
