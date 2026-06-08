package llama

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

type staticResolver struct {
	specs map[string]ModelSpec
	err   error
}

func (s staticResolver) Resolve(modelID string) (ModelSpec, error) {
	if s.err != nil {
		return ModelSpec{}, s.err
	}
	spec, ok := s.specs[modelID]
	if !ok {
		return ModelSpec{}, registry.ErrNotFound
	}
	return spec, nil
}

func TestLlamaProvider_LoadModel_validation(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())

	if err := p.LoadModel(""); err == nil {
		t.Fatal("LoadModel(\"\") error = nil, want error")
	}
}

func TestLlamaProvider_LoadModel_resolveError(t *testing.T) {
	p := NewLlamaProvider(staticResolver{err: registry.ErrNotFound}, DefaultConfig())

	err := p.LoadModel("missing")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Fatalf("LoadModel() error = %v, want ErrNotFound", err)
	}
}

func TestLlamaProvider_LoadModel_pathInaccessible(t *testing.T) {
	p := NewLlamaProvider(staticResolver{
		specs: map[string]ModelSpec{
			"m1": {ID: "m1", Path: filepath.Join(t.TempDir(), "nope.gguf")},
		},
	}, DefaultConfig())

	err := p.LoadModel("m1")
	if !errors.Is(err, ErrPathInaccessible) {
		t.Fatalf("LoadModel() error = %v, want ErrPathInaccessible", err)
	}
}

func TestLlamaProvider_LoadModel_duplicate(t *testing.T) {
	p := NewLlamaProvider(staticResolver{
		specs: map[string]ModelSpec{
			"m1": {ID: "m1", Path: "/unused"},
		},
	}, DefaultConfig())

	p.mu.Lock()
	p.models["m1"] = &modelSlot{path: "/unused"}
	p.mu.Unlock()

	err := p.LoadModel("m1")
	if !errors.Is(err, runtime.ErrModelAlreadyLoaded) {
		t.Fatalf("LoadModel() error = %v, want ErrModelAlreadyLoaded", err)
	}
}

func TestLlamaProvider_UnloadModel_notLoaded(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())

	err := p.UnloadModel("m1")
	if !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("UnloadModel() error = %v, want ErrModelNotLoaded", err)
	}
}

func TestLlamaProvider_portStubs(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())

	_, err := p.Generate(runtime.GenerateRequest{ModelID: "m1"})
	if !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("Generate() error = %v, want ErrNotImplemented", err)
	}
	if err := p.SaveContext("c1"); !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("SaveContext() error = %v, want ErrNotImplemented", err)
	}
	if err := p.RestoreContext("c1"); !errors.Is(err, ErrNotImplemented) {
		t.Fatalf("RestoreContext() error = %v, want ErrNotImplemented", err)
	}
}

func TestRegistryResolver_Resolve(t *testing.T) {
	reg := registry.NewModelRegistry()
	_ = reg.RegisterModel(registry.ModelDefinition{
		ID:   "tiny",
		Path: "/models/tiny.gguf",
	})

	r := NewRegistryResolver(reg)
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
	reg := registry.NewModelRegistry()
	_ = reg.RegisterModel(registry.ModelDefinition{ID: "empty"})

	r := NewRegistryResolver(reg)
	_, err := r.Resolve("empty")
	if err == nil {
		t.Fatal("Resolve() error = nil, want path required error")
	}
}
