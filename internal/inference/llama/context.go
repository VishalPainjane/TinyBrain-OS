package llama

// modelSlot tracks per-model load state on the Go side.
// Native llama handles are managed in load_cpu.go (CGO) or absent when CGO is off.
type modelSlot struct {
	path string
}
