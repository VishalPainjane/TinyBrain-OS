package runtime

// ModelSpec is the minimum metadata required to load a GGUF model.
// See planning/decisions/009a-registry-resolver.md and 009c-architecture-review.md.
type ModelSpec struct {
	ID           string
	Path         string
	Quantization string
	MemoryBudget uint64
}

// ModelResolver resolves modelID to a load specification.
// Implemented at the composition boundary; consumed by runtime and inference adapters.
type ModelResolver interface {
	Resolve(modelID string) (ModelSpec, error)
}
