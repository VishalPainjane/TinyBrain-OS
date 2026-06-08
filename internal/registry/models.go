package registry

import "fmt"

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

// ModelRegistry stores and serves model definitions via a ModelStore backend.
type ModelRegistry struct {
	store ModelStore
}

// NewModelRegistry returns a registry backed by an in-memory store.
func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{store: NewInMemoryStore()}
}

// NewBboltModelRegistry opens a bbolt-backed registry at dbPath.
// When seedPath is non-empty and the store is empty, models are loaded from seedPath.
func NewBboltModelRegistry(dbPath, seedPath string) (*ModelRegistry, error) {
	store, err := NewBboltStore(dbPath)
	if err != nil {
		return nil, err
	}

	if seedPath != "" {
		empty, err := store.IsEmpty()
		if err != nil {
			_ = store.Close()
			return nil, fmt.Errorf("check bbolt store empty: %w", err)
		}
		if empty {
			if err := LoadModelsYAML(seedPath, store); err != nil {
				_ = store.Close()
				return nil, err
			}
		}
	}

	return &ModelRegistry{store: store}, nil
}

// RegisterModel inserts a model definition. Duplicate IDs return ErrDuplicateID.
func (r *ModelRegistry) RegisterModel(def ModelDefinition) error {
	return r.store.RegisterModel(def)
}

// GetModel returns the model definition for id.
func (r *ModelRegistry) GetModel(id string) (ModelDefinition, error) {
	return r.store.GetModel(id)
}

// ListModels returns all registered model definitions.
func (r *ModelRegistry) ListModels() []ModelDefinition {
	return r.store.ListModels()
}

// Close releases registry store resources.
func (r *ModelRegistry) Close() error {
	return r.store.Close()
}
