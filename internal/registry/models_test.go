package registry_test

import (
	"errors"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func TestModelRegistry_RegisterAndGet(t *testing.T) {
	reg := registry.NewModelRegistry()
	def := registry.ModelDefinition{
		ID:           "model-7b",
		Path:         "/models/7b.gguf",
		SizeBytes:    4_000_000_000,
		MemoryBudget: 6_000_000_000,
		Capabilities: []string{"chat"},
		Quantization: "Q4_K_M",
	}

	if err := reg.RegisterModel(def); err != nil {
		t.Fatalf("RegisterModel() error = %v", err)
	}

	got, err := reg.GetModel("model-7b")
	if err != nil {
		t.Fatalf("GetModel() error = %v", err)
	}
	if got.Path != def.Path {
		t.Errorf("Path = %q, want %q", got.Path, def.Path)
	}
}

func TestModelRegistry_ListModels(t *testing.T) {
	reg := registry.NewModelRegistry()

	for _, id := range []string{"m1", "m2"} {
		if err := reg.RegisterModel(registry.ModelDefinition{ID: id, Path: id}); err != nil {
			t.Fatalf("RegisterModel(%s) error = %v", id, err)
		}
	}

	list := reg.ListModels()
	if len(list) != 2 {
		t.Fatalf("ListModels() len = %d, want 2", len(list))
	}
}

func TestModelRegistry_DuplicateID(t *testing.T) {
	reg := registry.NewModelRegistry()
	def := registry.ModelDefinition{ID: "dup", Path: "p"}

	if err := reg.RegisterModel(def); err != nil {
		t.Fatalf("first RegisterModel() error = %v", err)
	}
	if err := reg.RegisterModel(def); !errors.Is(err, registry.ErrDuplicateID) {
		t.Errorf("second RegisterModel() error = %v, want ErrDuplicateID", err)
	}
}

func TestModelRegistry_GetNotFound(t *testing.T) {
	reg := registry.NewModelRegistry()

	_, err := reg.GetModel("missing")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Errorf("GetModel() error = %v, want ErrNotFound", err)
	}
}
