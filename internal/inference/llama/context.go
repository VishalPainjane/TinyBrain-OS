package llama

// modelSlot tracks per-model load state on the Go side.
// Native llama model handles live in nativeModels; context handles in nativeContexts (CGO).
type modelSlot struct {
	path string
}
