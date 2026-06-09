//go:build cgo && integration

package llama

import (
	"errors"
	"os"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func TestLlamaProvider_LoadUnload_integration(t *testing.T) {
	path := os.Getenv("TB_TEST_GGUF_PATH")
	if path == "" {
		t.Skip("TB_TEST_GGUF_PATH not set")
	}

	p := NewLlamaProvider(staticResolver{
		specs: map[string]runtime.ModelSpec{
			"test-model": {ID: "test-model", Path: path},
		},
	}, DefaultConfig())

	if err := p.LoadModel("test-model"); err != nil {
		t.Fatalf("LoadModel() error = %v", err)
	}
	if err := p.UnloadModel("test-model"); err != nil {
		t.Fatalf("UnloadModel() error = %v", err)
	}
	if err := p.UnloadModel("test-model"); !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("UnloadModel() twice error = %v, want ErrModelNotLoaded", err)
	}
}
