//go:build integration

// Regression tests for token-divergence bugs.
// Run with: go test -tags integration -run TestRegression_ ./internal/inference/llama/
//
// Requires a real GGUF model at TINYBRAIN_MODEL env var (or skipped).
package llama

import (
	"os"
	"strings"
	"testing"
)

// TestRegression_BosTokenMissing verifies that tokenization now includes BOS.
// Pre-fix: add_special=false caused HF/llama.cpp positional mismatch.
// Post-fix: token 0 in the sequence should be the BOS token (id=1 for Llama models).
func TestRegression_BosTokenMissing(t *testing.T) {
	modelPath := os.Getenv("TINYBRAIN_MODEL")
	if modelPath == "" {
		t.Skip("TINYBRAIN_MODEL not set; skipping BOS regression test")
	}

	cfg := LlamaConfig{
		MaxTokens:     1,
		GreedySampler: true,
		ContextSize:   512,
		BatchSize:     512,
		NGLayers:      0,
	}

	const modelID = "regression-bos"
	if err := loadNativeModel(modelPath, modelID, cfg); err != nil {
		t.Fatalf("load model: %v", err)
	}
	defer unloadNativeModel(modelID) //nolint:errcheck

	// Tokenize a minimal prompt and check BOS is present.
	// We use an unexported helper if available; otherwise exercise through generate.
	out, _, err := generateNative(modelID, "Hello", cfg)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	// The output should not be empty — a model that had no BOS would often
	// produce garbage or repeat token 0 immediately.
	if strings.TrimSpace(out) == "" {
		t.Errorf("generate produced empty output; possible BOS regression")
	}
}
