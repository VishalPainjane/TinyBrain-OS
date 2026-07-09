//go:build !cgo

package llama

import "github.com/VishalPainjane/TinyBrain-OS/internal/runtime"

// stubBackend is the !cgo implementation of the backend interface.
// It returns ErrCGODisabled for all operations.
type stubBackend struct{}

func (s *stubBackend) loadModel(_, _ string, _ LlamaConfig) error {
	return ErrCGODisabled
}

func (s *stubBackend) unloadModel(_ string) error {
	return ErrCGODisabled
}

func (s *stubBackend) generate(_ runtime.GenerateRequest, _ LlamaConfig) (string, uint32, generateStats, error) {
	return "", 0, generateStats{}, ErrCGODisabled
}

func (s *stubBackend) formatChat(_ string, _ []runtime.ChatMessage, _ runtime.FormatChatOpts) (string, string, error) {
	return "", "", ErrCGODisabled
}

func (s *stubBackend) getMetadata(modelID string) (runtime.ModelCapabilities, error) {
	return runtime.ModelCapabilities{}, ErrCGODisabled
}

func (s *stubBackend) saveContext(_, _ string) error {
	return ErrCGODisabled
}

func (s *stubBackend) restoreContext(_, _ string) error {
	return ErrCGODisabled
}

// selectBackend returns the stub when CGO is disabled.
func selectBackend() backend {
	return &stubBackend{}
}
