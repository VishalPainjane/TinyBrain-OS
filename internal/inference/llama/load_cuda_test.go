package llama

import "testing"

// TestDoesNotRequireGPU verifies config helpers used by the CUDA build path without a GPU.
func TestDoesNotRequireGPU(t *testing.T) {
	cfg := DefaultConfig()
	cfg.NGLayers = -1
	if got := EffectiveNGLayers(cfg.NGLayers); got != -1 {
		t.Fatalf("EffectiveNGLayers(-1) = %d, want -1", got)
	}
}
