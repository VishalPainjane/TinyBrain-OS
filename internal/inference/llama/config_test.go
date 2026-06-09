package llama

import (
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/hardware"
)

func TestDefaultConfig_NGLayersZero(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.NGLayers != 0 {
		t.Fatalf("DefaultConfig().NGLayers = %d, want 0", cfg.NGLayers)
	}
}

func TestConfigFromProbe_CUDA(t *testing.T) {
	cfg := ConfigFromProbe(hardware.ProbeResult{Backend: hardware.BackendCUDA})
	if cfg.NGLayers != -1 {
		t.Fatalf("ConfigFromProbe(CUDA).NGLayers = %d, want -1", cfg.NGLayers)
	}
}

func TestConfigFromProbe_CPU(t *testing.T) {
	cfg := ConfigFromProbe(hardware.ProbeResult{Backend: hardware.BackendCPU})
	if cfg.NGLayers != 0 {
		t.Fatalf("ConfigFromProbe(CPU).NGLayers = %d, want 0", cfg.NGLayers)
	}
}

func TestConfigFromProbe_Metal(t *testing.T) {
	cfg := ConfigFromProbe(hardware.ProbeResult{Backend: hardware.BackendMetal})
	if cfg.NGLayers != 0 {
		t.Fatalf("ConfigFromProbe(Metal).NGLayers = %d, want 0", cfg.NGLayers)
	}
}

func TestEffectiveNGLayers_Clamp(t *testing.T) {
	tests := []struct {
		name string
		in   int32
		want int32
	}{
		{"all layers", -1, -1},
		{"zero", 0, 0},
		{"positive", 32, 32},
		{"below min clamped", -2, -1},
		{"far below min clamped", -100, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EffectiveNGLayers(tt.in); got != tt.want {
				t.Fatalf("EffectiveNGLayers(%d) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
