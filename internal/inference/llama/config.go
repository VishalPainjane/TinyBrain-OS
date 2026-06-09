package llama

// LlamaConfig holds llama.cpp load and generate settings.
// NGLayers is applied only when the binary is built with -tags cuda; CPU builds force 0 at load.
// See docs/architecture/inference-backend-matrix.md.
type LlamaConfig struct {
	ContextSize   uint32
	Threads       int
	UseMMAP       bool
	MaxTokens     uint32
	Temperature   float32
	BatchSize     uint32
	Seed          uint32
	GreedySampler bool
	// NGLayers is the number of model layers to offload to GPU (-1 = all layers, 0 = CPU only).
	NGLayers int32
}

// DefaultConfig returns conservative defaults with NGLayers=0 (safe for CPU binaries).
func DefaultConfig() LlamaConfig {
	return LlamaConfig{
		ContextSize:   512,
		Threads:       4,
		UseMMAP:       true,
		MaxTokens:     128,
		Temperature:   0.8,
		BatchSize:     512,
		Seed:          0xFFFFFFFF,
		GreedySampler: false,
		NGLayers:      0,
	}
}

// EffectiveNGLayers clamps cfg.NGLayers for llama.cpp: values below -1 become -1 (all layers).
func EffectiveNGLayers(ng int32) int32 {
	if ng < -1 {
		return -1
	}
	return ng
}
