package registry_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func TestBboltModelRegistry_PersistenceSurvivesRestart(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "models.db")

	def := registry.ModelDefinition{
		ID:           "model-7b",
		Path:         "/models/7b.gguf",
		SizeBytes:    4_000_000_000,
		MemoryBudget: 6_000_000_000,
		Capabilities: []string{"chat"},
		Quantization: "Q4_K_M",
	}

	reg, err := registry.NewBboltModelRegistry(dbPath, "")
	if err != nil {
		t.Fatalf("NewBboltModelRegistry() error = %v", err)
	}
	if err := reg.RegisterModel(def); err != nil {
		t.Fatalf("RegisterModel() error = %v", err)
	}
	if err := reg.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reg2, err := registry.NewBboltModelRegistry(dbPath, "")
	if err != nil {
		t.Fatalf("reopen NewBboltModelRegistry() error = %v", err)
	}
	defer reg2.Close()

	got, err := reg2.GetModel("model-7b")
	if err != nil {
		t.Fatalf("GetModel() error = %v", err)
	}
	if got.Path != def.Path || got.SizeBytes != def.SizeBytes || got.MemoryBudget != def.MemoryBudget {
		t.Errorf("GetModel() = %+v, want %+v", got, def)
	}

	list := reg2.ListModels()
	if len(list) != 1 {
		t.Fatalf("ListModels() len = %d, want 1", len(list))
	}
}

func TestBboltModelRegistry_DuplicateID(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "models.db")

	reg, err := registry.NewBboltModelRegistry(dbPath, "")
	if err != nil {
		t.Fatalf("NewBboltModelRegistry() error = %v", err)
	}
	defer reg.Close()

	def := registry.ModelDefinition{ID: "dup", Path: "p"}
	if err := reg.RegisterModel(def); err != nil {
		t.Fatalf("first RegisterModel() error = %v", err)
	}
	if err := reg.RegisterModel(def); !errors.Is(err, registry.ErrDuplicateID) {
		t.Errorf("second RegisterModel() error = %v, want ErrDuplicateID", err)
	}
}

func TestBboltModelRegistry_SeedWhenEmpty(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "models.db")
	seedPath := filepath.Join("..", "..", "testdata", "models.yaml")

	reg, err := registry.NewBboltModelRegistry(dbPath, seedPath)
	if err != nil {
		t.Fatalf("NewBboltModelRegistry() error = %v", err)
	}
	defer reg.Close()

	got, err := reg.GetModel("tinyllama-q4")
	if err != nil {
		t.Fatalf("GetModel() error = %v", err)
	}
	if got.Path != "/models/tinyllama-q4.gguf" {
		t.Errorf("Path = %q, want /models/tinyllama-q4.gguf", got.Path)
	}
	if got.SizeBytes == 0 || got.MemoryBudget == 0 {
		t.Errorf("SizeBytes or MemoryBudget not set: %+v", got)
	}
}

func TestBboltModelRegistry_SkipSeedWhenNotEmpty(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "models.db")
	seedPath := filepath.Join("..", "..", "testdata", "models.yaml")

	reg, err := registry.NewBboltModelRegistry(dbPath, "")
	if err != nil {
		t.Fatalf("NewBboltModelRegistry() error = %v", err)
	}
	if err := reg.RegisterModel(registry.ModelDefinition{ID: "existing", Path: "p"}); err != nil {
		t.Fatalf("RegisterModel() error = %v", err)
	}
	if err := reg.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	reg2, err := registry.NewBboltModelRegistry(dbPath, seedPath)
	if err != nil {
		t.Fatalf("reopen with seed error = %v", err)
	}
	defer reg2.Close()

	if _, err := reg2.GetModel("tinyllama-q4"); !errors.Is(err, registry.ErrNotFound) {
		t.Errorf("GetModel(tinyllama-q4) error = %v, want ErrNotFound when store not empty", err)
	}
	if _, err := reg2.GetModel("existing"); err != nil {
		t.Errorf("GetModel(existing) error = %v", err)
	}
}

func TestLoadModelsYAML_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	badPath := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(badPath, []byte("models:\n  - id: []\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := registry.NewInMemoryStore()
	if err := registry.LoadModelsYAML(badPath, store); err == nil {
		t.Fatal("LoadModelsYAML() error = nil, want parse error")
	}
}
