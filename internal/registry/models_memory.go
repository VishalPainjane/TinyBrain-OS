package registry

import (
	"fmt"
	"sync"
)

// InMemoryStore holds model definitions in a mutex-protected map.
type InMemoryStore struct {
	mu     sync.RWMutex
	models map[string]ModelDefinition
}

// NewInMemoryStore returns an empty in-memory model store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		models: make(map[string]ModelDefinition),
	}
}

// RegisterModel inserts a model definition. Duplicate IDs return ErrDuplicateID.
func (s *InMemoryStore) RegisterModel(def ModelDefinition) error {
	if def.ID == "" {
		return fmt.Errorf("model ID is required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.models[def.ID]; exists {
		return ErrDuplicateID
	}

	s.models[def.ID] = def
	return nil
}

// GetModel returns the model definition for id.
func (s *InMemoryStore) GetModel(id string) (ModelDefinition, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	def, ok := s.models[id]
	if !ok {
		return ModelDefinition{}, ErrNotFound
	}
	return def, nil
}

// ListModels returns all registered model definitions.
func (s *InMemoryStore) ListModels() []ModelDefinition {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]ModelDefinition, 0, len(s.models))
	for _, def := range s.models {
		out = append(out, def)
	}
	return out
}

// Close releases in-memory resources. It is a no-op.
func (s *InMemoryStore) Close() error {
	return nil
}
