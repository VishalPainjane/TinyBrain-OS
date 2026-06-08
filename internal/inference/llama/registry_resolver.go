package llama

import (
	"fmt"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

// RegistryResolver adapts ModelRegistry to ModelResolver.
// This is the only file in the inference tree that imports internal/registry.
type RegistryResolver struct {
	reg *registry.ModelRegistry
}

// NewRegistryResolver returns a resolver backed by reg.
func NewRegistryResolver(reg *registry.ModelRegistry) *RegistryResolver {
	return &RegistryResolver{reg: reg}
}

// Resolve returns model metadata from the registry.
func (r *RegistryResolver) Resolve(modelID string) (ModelSpec, error) {
	def, err := r.reg.GetModel(modelID)
	if err != nil {
		return ModelSpec{}, err
	}
	if def.Path == "" {
		return ModelSpec{}, fmt.Errorf("model %q: path is required", modelID)
	}
	return ModelSpec{
		ID:           def.ID,
		Path:         def.Path,
		Quantization: def.Quantization,
		MemoryBudget: def.MemoryBudget,
	}, nil
}
