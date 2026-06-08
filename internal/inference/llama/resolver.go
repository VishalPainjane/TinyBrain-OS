package llama

// ModelSpec is the minimum metadata required to load a GGUF model.
// See planning/decisions/009a-registry-resolver.md.
type ModelSpec struct {
	ID           string
	Path         string
	Quantization string
	MemoryBudget uint64
}

// ModelResolver resolves modelID to a load specification.
// Registry and other backends implement this at the composition root.
type ModelResolver interface {
	Resolve(modelID string) (ModelSpec, error)
}
