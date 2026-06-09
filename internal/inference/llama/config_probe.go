package llama

import "github.com/VishalPainjane/TinyBrain-OS/internal/hardware"

// ConfigFromProbe returns DefaultConfig with NGLayers set from the probed inference backend.
// BackendCUDA recommends -1 (all layers); all other backends use 0.
// The caller must build with -tags cuda for GPU offload to take effect.
func ConfigFromProbe(probe hardware.ProbeResult) LlamaConfig {
	cfg := DefaultConfig()
	switch probe.Backend {
	case hardware.BackendCUDA:
		cfg.NGLayers = -1
	default:
		cfg.NGLayers = 0
	}
	return cfg
}
