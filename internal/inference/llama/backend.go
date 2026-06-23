package llama

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
	generate(modelID, prompt string, cfg LlamaConfig) (string, uint32, generateStats, error)
}
