package llama

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

type staticResolver struct {
	specs map[string]runtime.ModelSpec
	err   error
}

func (s staticResolver) Resolve(modelID string) (runtime.ModelSpec, error) {
	if s.err != nil {
		return runtime.ModelSpec{}, s.err
	}
	spec, ok := s.specs[modelID]
	if !ok {
		return runtime.ModelSpec{}, registry.ErrNotFound
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
		specs: map[string]runtime.ModelSpec{
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
		specs: map[string]runtime.ModelSpec{
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

func TestLlamaProvider_Generate_notLoaded(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())

	_, err := p.Generate(runtime.GenerateRequest{ModelID: "m1", Prompt: "hi"})
	if !errors.Is(err, runtime.ErrModelNotLoaded) {
		t.Fatalf("Generate() error = %v, want ErrModelNotLoaded", err)
	}
}

func TestLlamaProvider_Generate_emptyModelID(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())

	_, err := p.Generate(runtime.GenerateRequest{ModelID: "", Prompt: "hi"})
	if err == nil {
		t.Fatal("Generate() error = nil, want error")
	}
}

func TestLlamaProvider_portStubs(t *testing.T) {
	p := NewLlamaProvider(staticResolver{}, DefaultConfig())
	p.models["m1"] = &modelSlot{path: "dummy"}

	if err := p.SaveContext("m1", "c1"); !errors.Is(err, ErrNotImplemented) && !errors.Is(err, ErrCGODisabled) {
		t.Fatalf("SaveContext() error = %v, want ErrNotImplemented or ErrCGODisabled", err)
	}
	if err := p.RestoreContext("m1", "c1"); !errors.Is(err, ErrNotImplemented) && !errors.Is(err, ErrCGODisabled) {
		t.Fatalf("RestoreContext() error = %v, want ErrNotImplemented or ErrCGODisabled", err)
	}
}
