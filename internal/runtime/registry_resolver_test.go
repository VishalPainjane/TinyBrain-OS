package runtime_test

import (
	"errors"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

func TestRegistryResolver_Resolve(t *testing.T) {
	t.Parallel()

	reg := registry.NewModelRegistry()
	_ = reg.RegisterModel(registry.ModelDefinition{
		ID:   "tiny",
		Path: "/models/tiny.gguf",
	})

	r := runtime.NewRegistryResolver(reg)
	spec, err := r.Resolve("tiny")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if spec.Path != "/models/tiny.gguf" {
		t.Fatalf("Path = %q, want /models/tiny.gguf", spec.Path)
	}

	_, err = r.Resolve("missing")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("Resolve(missing) error = %v, want ErrNotFound", err)
	}
}

func TestRegistryResolver_emptyPath(t *testing.T) {
	t.Parallel()

	reg := registry.NewModelRegistry()
	_ = reg.RegisterModel(registry.ModelDefinition{ID: "empty"})

	r := runtime.NewRegistryResolver(reg)
	_, err := r.Resolve("empty")
	if err == nil {
		t.Fatal("Resolve() error = nil, want path required error")
	}
}
