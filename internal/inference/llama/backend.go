package llama

import "github.com/VishalPainjane/TinyBrain-OS/internal/runtime"

// backend is the internal interface that decouples provider.go from the
// concrete inference implementation.  Each platform selects exactly one
// implementation at compile-time (CPU static CGO) or at run-time (Windows GPU
// dynamic DLL).
//
// This interface is private to the llama package.  Callers use LlamaProvider
// which implements runtime.InferenceProvider.
type backend interface {
	// loadModel loads a GGUF file into the backend and associates it with modelID.
	loadModel(path, modelID string, cfg LlamaConfig) error

	// unloadModel frees all native resources for modelID.
	unloadModel(modelID string) error

	// generate runs inference and returns (text, tokensProduced, stats, error).
	generate(req runtime.GenerateRequest, cfg LlamaConfig) (string, uint32, generateStats, error)

	// formatChat applies the model's chat template to a list of structured messages.
	formatChat(modelID string, messages []runtime.ChatMessage, opts runtime.FormatChatOpts) (string, string, error)

	// getMetadata retrieves capabilities directly from the model metadata.
	getMetadata(modelID string) (runtime.ModelCapabilities, error)

	// saveContext persists the native KV cache state for the model.
	saveContext(modelID, ctxID string) error

	// restoreContext loads a previously persisted KV cache state.
	restoreContext(modelID, ctxID string) error
}
