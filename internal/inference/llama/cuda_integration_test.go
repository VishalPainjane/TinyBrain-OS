//go:build cgo && cuda && !rocm && !metal && !vulkan && integration

package llama

import (
	"os"
	"strconv"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func cudaIntegrationProvider(t *testing.T) (*LlamaProvider, string) {
	t.Helper()

	if os.Getenv("TB_CUDA_INTEGRATION") != "1" {
		t.Skip("TB_CUDA_INTEGRATION not set to 1")
	}

	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	modelID := "test-model"
	cfg := DefaultConfig()
	cfg.GreedySampler = true
	cfg.NGLayers = cudaIntegrationNGLayers()

	p := NewLlamaProvider(staticResolver{
		specs: map[string]runtime.ModelSpec{
			modelID: {ID: modelID, Path: path},
		},
	}, cfg)

	if err := p.LoadModel(modelID); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	return p, modelID
}

func cudaIntegrationNGLayers() int32 {
	if v := os.Getenv("TB_NGLAYERS"); v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return -1
		}
		return int32(n)
	}
	return -1
}

func TestLlamaProvider_LoadGenerate_CUDA(t *testing.T) {
	p, modelID := cudaIntegrationProvider(t)
	defer func() { _ = p.UnloadModel(modelID) }()

	resp, err := p.Generate(runtime.GenerateRequest{
		ModelID: modelID,
		Prompt:  "Hello",
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if resp.Output == "" {
		t.Fatal("Generate() Output is empty")
	}
	if resp.TokensProduced <= 0 {
		t.Fatalf("TokensProduced = %d, want > 0", resp.TokensProduced)
	}
	t.Logf("CUDA Generate (%d tokens, NGLayers=%d): %q",
		resp.TokensProduced, cudaIntegrationNGLayers(), resp.Output)
}

func TestLlamaProvider_CUDA_TPS(t *testing.T) {
	p, modelID := cudaIntegrationProvider(t)
	defer func() { _ = p.UnloadModel(modelID) }()

	output, tokens, stats, err := p.generateBackendTimed(modelID, "Hello")
	if err != nil {
		t.Fatalf("generateBackendTimed() error = %v", err)
	}
	if output == "" || tokens == 0 {
		t.Fatalf("output=%q tokens=%d, want non-empty", output, tokens)
	}

	decodeSecs := float64(stats.DecodeMicros) / 1_000_000
	var tps float64
	if decodeSecs > 0 {
		tps = float64(tokens) / decodeSecs
	}
	t.Logf("CUDA TTFT=%dus decode=%dus tokens=%d TPS=%.2f NGLayers=%d",
		stats.TTFTMicros, stats.DecodeMicros, tokens, tps, cudaIntegrationNGLayers())
}
