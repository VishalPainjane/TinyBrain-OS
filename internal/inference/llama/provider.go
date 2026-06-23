package llama

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

// Compile-time check that LlamaProvider implements runtime.InferenceProvider.
var _ runtime.InferenceProvider = (*LlamaProvider)(nil)

// LlamaProvider implements runtime.InferenceProvider using llama.cpp.
//
// The active hardware backend (CPU static CGO, or Windows GPU dynamic DLL) is
// selected once at construction via selectBackend().  All inference calls are
// delegated to the backend interface, keeping the provider and scheduler layers
// fully decoupled from hardware specifics.
//
// See ADR-006 for the rationale for runtime backend selection on Windows.
type LlamaProvider struct {
	mu       sync.Mutex
	b        backend // selected at construction; immutable after that
	resolver runtime.ModelResolver
	cfg      LlamaConfig
	models   map[string]*modelSlot
}

// NewLlamaProvider returns a provider backed by the best available inference
// engine on the current platform.
func NewLlamaProvider(resolver runtime.ModelResolver, cfg LlamaConfig) *LlamaProvider {
	return &LlamaProvider{
		b:        selectBackend(),
		resolver: resolver,
		cfg:      cfg,
		models:   make(map[string]*modelSlot),
	}
}

// LoadModel resolves modelID and loads the GGUF file via the selected backend.
func (p *LlamaProvider) LoadModel(modelID string) error {
	if modelID == "" {
		return fmt.Errorf("model ID is required")
	}

	p.mu.Lock()
	if _, ok := p.models[modelID]; ok {
		p.mu.Unlock()
		return fmt.Errorf("%w: %s", runtime.ErrModelAlreadyLoaded, modelID)
	}
	p.mu.Unlock()

	spec, err := p.resolver.Resolve(modelID)
	if err != nil {
		return err
	}

	cleanPath := filepath.Clean(spec.Path)
	if _, err := os.Stat(cleanPath); err != nil {
		return fmt.Errorf("%w: %s", ErrPathInaccessible, cleanPath)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.models[modelID]; ok {
		return fmt.Errorf("%w: %s", runtime.ErrModelAlreadyLoaded, modelID)
	}

	if err := p.b.loadModel(cleanPath, modelID, p.cfg); err != nil {
		return err
	}

	p.models[modelID] = &modelSlot{path: cleanPath}
	return nil
}

// UnloadModel frees the model and removes it from the loaded set.
func (p *LlamaProvider) UnloadModel(modelID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.models[modelID]; !ok {
		return fmt.Errorf("%w: %s", runtime.ErrModelNotLoaded, modelID)
	}

	if err := p.b.unloadModel(modelID); err != nil {
		return err
	}

	delete(p.models, modelID)
	return nil
}

// Generate runs single-prompt inference on a loaded model.
func (p *LlamaProvider) Generate(req runtime.GenerateRequest) (runtime.GenerateResponse, error) {
	if req.ModelID == "" {
		return runtime.GenerateResponse{}, fmt.Errorf("model ID is required")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.models[req.ModelID]; !ok {
		return runtime.GenerateResponse{}, fmt.Errorf("%w: %s", runtime.ErrModelNotLoaded, req.ModelID)
	}

	output, tokens, _, err := p.b.generate(req.ModelID, req.Prompt, p.cfg)
	if err != nil {
		return runtime.GenerateResponse{}, err
	}

	return runtime.GenerateResponse{
		ModelID:        req.ModelID,
		Output:         output,
		TokensProduced: int(tokens),
	}, nil
}
