package loader_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/loader"
)

func stubModelPath(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("stub-gguf"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func TestStubLoader_LoadTransition(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	var l loader.Loader = loader.NewStubLoader()

	state, err := l.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateNotLoaded {
		t.Fatalf("initial state = %s, want NOT_LOADED", state)
	}

	if err := l.Load("model-a", path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	state, err = l.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateActive {
		t.Fatalf("state after load = %s, want ACTIVE", state)
	}
}

func TestStubLoader_UnloadTransition(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	var l loader.Loader = loader.NewStubLoader()

	if err := l.Load("model-a", path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if err := l.Unload("model-a"); err != nil {
		t.Fatalf("Unload() error = %v", err)
	}

	state, err := l.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("state after unload = %s, want UNLOADED", state)
	}
}

func TestStubLoader_DuplicateLoadPrevented(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	var l loader.Loader = loader.NewStubLoader()

	if err := l.Load("model-a", path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if err := l.Load("model-a", path); !errors.Is(err, loader.ErrModelAlreadyLoaded) {
		t.Fatalf("duplicate Load() error = %v, want ErrModelAlreadyLoaded", err)
	}
}

func TestStubLoader_WarmAndPrefetch(t *testing.T) {
	t.Parallel()

	pathA := stubModelPath(t, "model-a.gguf")
	pathB := stubModelPath(t, "model-b.gguf")
	var l loader.Loader = loader.NewStubLoader()

	if err := l.Load("model-a", pathA); err != nil {
		t.Fatalf("Load(model-a) error = %v", err)
	}
	if err := l.Warm("model-a"); err != nil {
		t.Fatalf("Warm(model-a) error = %v", err)
	}

	state, err := l.State("model-a")
	if err != nil {
		t.Fatalf("State(model-a) error = %v", err)
	}
	if state != loader.StateWarm {
		t.Fatalf("state = %s, want WARM", state)
	}

	if err := l.Prefetch("model-b", pathB); err != nil {
		t.Fatalf("Prefetch(model-b) error = %v", err)
	}

	state, err = l.State("model-b")
	if err != nil {
		t.Fatalf("State(model-b) error = %v", err)
	}
	if state != loader.StateWarm {
		t.Fatalf("state = %s, want WARM", state)
	}
}

func TestStubLoader_Evict(t *testing.T) {
	t.Parallel()

	path := stubModelPath(t, "model-a.gguf")
	var l loader.Loader = loader.NewStubLoader()

	if err := l.Load("model-a", path); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if err := l.Evict("model-a"); err != nil {
		t.Fatalf("Evict() error = %v", err)
	}

	state, err := l.State("model-a")
	if err != nil {
		t.Fatalf("State() error = %v", err)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("state = %s, want UNLOADED", state)
	}
}

func TestStubLoader_LoadRejectsMissingPath(t *testing.T) {
	t.Parallel()

	var l loader.Loader = loader.NewStubLoader()
	err := l.Load("model-a", filepath.Join(t.TempDir(), "missing.gguf"))
	if !errors.Is(err, loader.ErrStubPathInaccessible) {
		t.Fatalf("Load() error = %v, want ErrStubPathInaccessible", err)
	}
}

func TestStubLoader_LRUEvictionShell(t *testing.T) {
	t.Parallel()

	pathA := stubModelPath(t, "model-a.gguf")
	pathB := stubModelPath(t, "model-b.gguf")
	pathC := stubModelPath(t, "model-c.gguf")

	policy := loader.NewLRUEvictionPolicy()
	l := loader.NewStubLoader(
		loader.WithCapacity(2),
		loader.WithEvictionPolicy(policy),
	)

	if err := l.Load("model-a", pathA); err != nil {
		t.Fatalf("Load(model-a) error = %v", err)
	}
	if err := l.Load("model-b", pathB); err != nil {
		t.Fatalf("Load(model-b) error = %v", err)
	}

	state, err := l.State("model-a")
	if err != nil {
		t.Fatalf("State(model-a) error = %v", err)
	}
	if state != loader.StateActive {
		t.Fatalf("model-a state = %s, want ACTIVE", state)
	}

	if err := l.Load("model-c", pathC); err != nil {
		t.Fatalf("Load(model-c) error = %v", err)
	}

	state, err = l.State("model-a")
	if err != nil {
		t.Fatalf("State(model-a) after eviction error = %v", err)
	}
	if state != loader.StateUnloaded {
		t.Fatalf("model-a state = %s, want UNLOADED after LRU eviction", state)
	}

	state, err = l.State("model-c")
	if err != nil {
		t.Fatalf("State(model-c) error = %v", err)
	}
	if state != loader.StateActive {
		t.Fatalf("model-c state = %s, want ACTIVE", state)
	}
}
