//go:build !cgo

package llama

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLlamaProvider_LoadUnload_cgoDisabled(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stub.gguf")
	if err := os.WriteFile(path, []byte{0}, 0o600); err != nil {
		t.Fatal(err)
	}

	p := NewLlamaProvider(staticResolver{
		specs: map[string]ModelSpec{
			"m1": {ID: "m1", Path: path},
		},
	}, DefaultConfig())

	err := p.LoadModel("m1")
	if !errors.Is(err, ErrCGODisabled) {
		t.Fatalf("LoadModel() error = %v, want ErrCGODisabled when CGO off", err)
	}
}
