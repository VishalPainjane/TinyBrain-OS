//go:build integration && windows && !cgo_cuda_static

package llama

import (
	"os"
	"testing"
)

// TestRegression_SamplerChainParamsABI verifies that buildSampler does not
// crash or return a null chain when the DLL backend is in use.
//
// Pre-fix: llama_sampler_chain_default_params called with no args, return value
// read as byte — corrupted params caused chainPtr=0 (silent generation failure).
func TestRegression_SamplerChainParamsABI(t *testing.T) {
	modelPath := os.Getenv("TINYBRAIN_MODEL")
	if modelPath == "" {
		t.Skip("TINYBRAIN_MODEL not set; skipping sampler ABI regression test")
	}

	dllDir := os.Getenv("TINYBRAIN_DLL_DIR")
	if dllDir == "" {
		t.Skip("TINYBRAIN_DLL_DIR not set; skipping DLL backend sampler ABI test")
	}

	b, err := NewWindowsDynamicBackend(dllDir)
	if err != nil {
		t.Skipf("DLL backend unavailable (%v); skipping", err)
	}

	cfg := LlamaConfig{
		MaxTokens:     5,
		GreedySampler: true,
		ContextSize:   512,
		BatchSize:     512,
	}

	const modelID = "regression-sampler-abi"
	if err := b.loadModel(modelPath, modelID, cfg); err != nil {
		t.Fatalf("load model: %v", err)
	}
	defer b.unloadModel(modelID) //nolint:errcheck

	modelPtr := b.models[modelID]
	samplerPtr := b.buildSampler(cfg, modelPtr, "")
	if samplerPtr == 0 {
		t.Fatal("buildSampler returned null; sampler chain params ABI is still broken")
	}
	b.procSamplerFree.Call(samplerPtr)
}
