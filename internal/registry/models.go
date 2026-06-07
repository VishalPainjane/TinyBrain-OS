package registry

import (
	"fmt"
	"sync"
)

// ModelDefinition describes a model entry in the registry.
// See docs/contracts/registry.md.
type ModelDefinition struct {
	ID           string
	Path         string
	SizeBytes    uint64
	MemoryBudget uint64
	Capabilities []string
	Quantization string
}

// ModelRegistry stores and serves model definitions in memory.
type ModelRegistry struct {
	mu     sync.RWMutex
	models map[string]ModelDefinition
}

// NewModelRegistry returns an empty in-memory model registry.
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{
		models: make(map[string]ModelDefinition),
	}
}

// RegisterModel inserts a model definition. Duplicate IDs return ErrDuplicateID.
func (r *ModelRegistry) RegisterModel(def ModelDefinition) error {
	if def.ID == "" {
		return fmt.Errorf("model ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.models[def.ID]; exists {
		return ErrDuplicateID
	}

	r.models[def.ID] = def
	return nil
}

// GetModel returns the model definition for id.
func (r *ModelRegistry) GetModel(id string) (ModelDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.models[id]
	if !ok {
		return ModelDefinition{}, ErrNotFound
	}
	return def, nil
}

// ListModels returns all registered model definitions.
func (r *ModelRegistry) ListModels() []ModelDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]ModelDefinition, 0, len(r.models))
	for _, def := range r.models {
		out = append(out, def)
	}
	return out
}
