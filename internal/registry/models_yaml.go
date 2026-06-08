package registry

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// modelsYAMLFile is the on-disk seed format for model definitions.
type modelsYAMLFile struct {
	Models []modelYAMLEntry `yaml:"models"`
}

type modelYAMLEntry struct {
	ID           string   `yaml:"id"`
	Path         string   `yaml:"path"`
	SizeBytes    uint64   `yaml:"size_bytes"`
	MemoryBudget uint64   `yaml:"memory_budget"`
	Capabilities []string `yaml:"capabilities"`
	Quantization string   `yaml:"quantization"`
}

// LoadModelsYAML reads path and registers each model into store.
// Duplicate IDs in the file or store return ErrDuplicateID.
func LoadModelsYAML(path string, store ModelStore) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read models yaml %q: %w", path, err)
	}

	var file modelsYAMLFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse models yaml %q: %w", path, err)
	}

	seen := make(map[string]struct{}, len(file.Models))
	for _, entry := range file.Models {
		if entry.ID == "" {
			return fmt.Errorf("models yaml %q: model entry missing id", path)
		}
		if _, dup := seen[entry.ID]; dup {
			return fmt.Errorf("models yaml %q: duplicate id %q", path, entry.ID)
		}
		seen[entry.ID] = struct{}{}

		def := ModelDefinition{
			ID:           entry.ID,
			Path:         entry.Path,
			SizeBytes:    entry.SizeBytes,
			MemoryBudget: entry.MemoryBudget,
			Capabilities: entry.Capabilities,
			Quantization: entry.Quantization,
		}
		if err := store.RegisterModel(def); err != nil {
			return fmt.Errorf("register model %q from %q: %w", entry.ID, path, err)
		}
	}

	return nil
}
