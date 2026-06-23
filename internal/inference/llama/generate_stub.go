//go:build !cgo

package llama

// stubBackend is the !cgo implementation of the backend interface.
// It returns ErrCGODisabled for all operations.
type stubBackend struct{}

func (s *stubBackend) loadModel(_, _ string, _ LlamaConfig) error {
	return ErrCGODisabled
}

func (s *stubBackend) unloadModel(_ string) error {
	return ErrCGODisabled
}

func (s *stubBackend) generate(_, _ string, _ LlamaConfig) (string, uint32, generateStats, error) {
	return "", 0, generateStats{}, ErrCGODisabled
}

// selectBackend returns the stub when CGO is disabled.
func selectBackend() backend {
	return &stubBackend{}
}
