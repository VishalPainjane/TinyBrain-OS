package llama

// LlamaConfig holds CPU-only llama.cpp load settings for 009a.
type LlamaConfig struct {
	ContextSize uint32
	Threads     int
	UseMMAP     bool
}

// DefaultConfig returns conservative CPU defaults (NGLayers=0 implicit).
func DefaultConfig() LlamaConfig {
	return LlamaConfig{
		ContextSize: 512,
		Threads:     4,
		UseMMAP:     true,
	}
}
