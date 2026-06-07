package loader

// ModelState is the lifecycle state of a loaded model.
// See docs/architecture/runtime.md.
type ModelState string

const (
	StateNotLoaded ModelState = "NOT_LOADED"
	StateLoading   ModelState = "LOADING"
	StateActive    ModelState = "ACTIVE"
	StateWarm      ModelState = "WARM"
	StateUnloading ModelState = "UNLOADING"
	StateUnloaded  ModelState = "UNLOADED"
)

// Valid reports whether s is a defined model lifecycle state.
func (s ModelState) Valid() bool {
	switch s {
	case StateNotLoaded, StateLoading, StateActive, StateWarm, StateUnloading, StateUnloaded:
		return true
	default:
		return false
	}
}

// Loader manages model load, unload, warm, prefetch, and evict operations.
// See tasks/009-model-loader.md.
type Loader interface {
	Load(modelID, path string) error
	Unload(modelID string) error
	Warm(modelID string) error
	Prefetch(modelID, path string) error
	Evict(modelID string) error
	State(modelID string) (ModelState, error)
}

// EvictionPolicy selects a model to evict when VRAM capacity is full.
// Shell hook until the memory layer exists.
type EvictionPolicy interface {
	SelectVictim(candidates []string) string
}
