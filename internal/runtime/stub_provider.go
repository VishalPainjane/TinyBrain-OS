package runtime

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrModelAlreadyLoaded = errors.New("model already loaded")
	ErrModelNotLoaded     = errors.New("model not loaded")
	ErrContextNotFound    = errors.New("context not found")
)

// StubProvider is a test InferenceProvider with in-memory state and canned responses.
type StubProvider struct {
	mu       sync.Mutex
	loaded   map[string]struct{}
	contexts map[string][]byte
}

// NewStubProvider returns a stub inference provider for testing.
func NewStubProvider() *StubProvider {
	return &StubProvider{
		loaded:   make(map[string]struct{}),
		contexts: make(map[string][]byte),
	}
}

// LoadModel marks modelID as loaded in memory.
func (p *StubProvider) LoadModel(modelID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.loaded[modelID]; ok {
		return fmt.Errorf("%w: %s", ErrModelAlreadyLoaded, modelID)
	}
	p.loaded[modelID] = struct{}{}
	return nil
}

// UnloadModel removes modelID from the in-memory loaded set.
func (p *StubProvider) UnloadModel(modelID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.loaded[modelID]; !ok {
		return fmt.Errorf("%w: %s", ErrModelNotLoaded, modelID)
	}
	delete(p.loaded, modelID)
	return nil
}

// Generate returns a canned structured response when the model is loaded.
func (p *StubProvider) Generate(req GenerateRequest) (GenerateResponse, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.loaded[req.ModelID]; !ok {
		return GenerateResponse{}, fmt.Errorf("%w: %s", ErrModelNotLoaded, req.ModelID)
	}

	return GenerateResponse{
		ModelID:        req.ModelID,
		Output:         fmt.Sprintf("stub response for %s", req.ModelID),
		TokensProduced: 1,
	}, nil
}

// SaveContext stores an empty placeholder blob for id.
func (p *StubProvider) SaveContext(modelID, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.contexts[id] = []byte("stub")
	return nil
}

// RestoreContext returns ErrContextNotFound if the id was never saved.
func (p *StubProvider) RestoreContext(modelID, id string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.contexts[id]; !ok {
		return ErrContextNotFound
	}
	return nil
}

func (p *StubProvider) FormatChat(modelID string, messages []ChatMessage, opts FormatChatOpts) (string, string, error) {
	return "", "", nil
}

func (p *StubProvider) GetMetadata(modelID string) (ModelCapabilities, error) {
	return ModelCapabilities{ModelID: modelID}, nil
}
