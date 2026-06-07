package runtime

// ModelRuntime manages model lifecycle and delegates inference to a provider.
// See docs/contracts/runtime.md and ADR-004.
type ModelRuntime interface {
	LoadModel(modelID string) error
	UnloadModel(modelID string) error
	Generate(req GenerateRequest) (GenerateResponse, error)
	SaveContext(id string) error
	RestoreContext(id string) error
}

// InferenceProvider is the inference port implemented by engine adapters.
// See docs/contracts/runtime.md and ADR-004.
type InferenceProvider interface {
	LoadModel(modelID string) error
	UnloadModel(modelID string) error
	Generate(req GenerateRequest) (GenerateResponse, error)
	SaveContext(id string) error
	RestoreContext(id string) error
}
