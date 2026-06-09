package llama

// LlamaConfig holds CPU-only llama.cpp load and generate settings.
type LlamaConfig struct {
	ContextSize   uint32
	Threads       int
	UseMMAP       bool
	MaxTokens     uint32
	Temperature   float32
	BatchSize     uint32
	Seed          uint32
	GreedySampler bool
}

// DefaultConfig returns conservative CPU defaults (NGLayers=0 implicit).
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
	}
}
