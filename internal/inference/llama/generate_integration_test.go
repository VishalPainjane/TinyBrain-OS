//go:build cgo && integration

package llama

import (
	"errors"
	"os"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func integrationProvider(t *testing.T) (*LlamaProvider, string) {
	t.Helper()

	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	modelID := "test-model"
	cfg := DefaultConfig()
	cfg.GreedySampler = true

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

func TestLlamaProvider_Generate_integration(t *testing.T) {
	p, modelID := integrationProvider(t)
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
	t.Logf("Generate output (%d tokens): %q", resp.TokensProduced, resp.Output)
}

func TestLlamaProvider_Generate_afterUnload(t *testing.T) {
	p, modelID := integrationProvider(t)

	if err := p.UnloadModel(modelID); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}

	_, err := p.Generate(runtime.GenerateRequest{ModelID: modelID, Prompt: "Hello"})
	if !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("Generate() error = %v, want ErrModelNotLoaded", err)
	}
}

func TestLlamaProvider_Generate_TTFT_TPS(t *testing.T) {
	p, modelID := integrationProvider(t)
	defer func() { _ = p.UnloadModel(modelID) }()

	output, tokens, stats, err := p.generateBackendTimed(modelID, "Hello")
	if err != nil {
		t.Fatalf("generateBackendTimed() error = %v", err)
	}
	if output == "" || tokens == 0 {
		t.Fatalf("output=%q tokens=%d, want non-empty", output, tokens)
	}
	if stats.TTFTMicros <= 0 {
		t.Fatalf("TTFTMicros = %d, want > 0", stats.TTFTMicros)
	}
	if stats.DecodeMicros < 0 {
		t.Fatalf("DecodeMicros = %d, want >= 0", stats.DecodeMicros)
	}

	tps := float64(tokens) / (float64(stats.DecodeMicros) / 1_000_000)
	t.Logf("TTFT=%dus decode=%dus tokens=%d TPS=%.2f", stats.TTFTMicros, stats.DecodeMicros, tokens, tps)
}

func TestLlamaProvider_LoadGenerateUnload_cycles(t *testing.T) {
	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	modelID := "test-model"
	cfg := DefaultConfig()
	cfg.GreedySampler = true

	for i := 0; i < 3; i++ {
		p := NewLlamaProvider(staticResolver{
			specs: map[string]runtime.ModelSpec{
				modelID: {ID: modelID, Path: path},
			},
		}, cfg)

		if err := p.LoadModel(modelID); err != nil {
			t.Fatalf("cycle %d LoadModel() error = %v", i, err)
		}

		_, err := p.Generate(runtime.GenerateRequest{ModelID: modelID, Prompt: "Hi"})
		if err != nil {
			t.Fatalf("cycle %d Generate() error = %v", i, err)
		}

		if err := p.UnloadModel(modelID); err != nil {
			t.Fatalf("cycle %d UnloadModel() error = %v", i, err)
		}
	}

	p := NewLlamaProvider(staticResolver{}, cfg)
	if err := p.UnloadModel(modelID); !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("UnloadModel() after cycles error = %v, want ErrModelNotLoaded", err)
	}
}
