package runtime

// ModelRuntime manages model lifecycle and delegates inference to a provider.
// See docs/contracts/runtime.md and ADR-004.
type ModelRuntime interface {
	LoadModel(modelID string) error
	UnloadModel(modelID string) error
	Generate(req GenerateRequest) (GenerateResponse, error)
	SaveContext(modelID, ctxID string) error
	RestoreContext(modelID, ctxID string) error
	FormatChat(modelID string, messages []ChatMessage, opts FormatChatOpts) (prompt string, tmpl string, err error)
	GetMetadata(modelID string) (ModelCapabilities, error)
}

// InferenceProvider is the inference port implemented by engine adapters.
// See docs/contracts/runtime.md and ADR-004.
type InferenceProvider interface {
	LoadModel(modelID string) error
	UnloadModel(modelID string) error
	Generate(req GenerateRequest) (GenerateResponse, error)
	SaveContext(modelID, ctxID string) error
	RestoreContext(modelID, ctxID string) error
	FormatChat(modelID string, messages []ChatMessage, opts FormatChatOpts) (prompt string, tmpl string, err error)
	GetMetadata(modelID string) (ModelCapabilities, error)
}
