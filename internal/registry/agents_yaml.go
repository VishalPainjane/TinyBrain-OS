package registry

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// agentsYAMLFile is the on-disk seed format for agent definitions.
type agentsYAMLFile struct {
	Agents []agentYAMLEntry `yaml:"agents"`
}

type agentYAMLEntry struct {
	ID              string                   `yaml:"id"`
	Name            string                   `yaml:"name"`
	ModelProfile    string                   `yaml:"model_profile"`
	Tools           []string                 `yaml:"tools"`
	ResourceProfile resourceProfileYAMLEntry `yaml:"resource_profile"`
	Priority        int                      `yaml:"priority"`
}

type resourceProfileYAMLEntry struct {
	MemoryMB uint64 `yaml:"memory_mb"`
	Priority int    `yaml:"priority"`
}

// LoadAgentsYAML reads path and registers each agent into the registry.
// Duplicate IDs in the file or registry return ErrDuplicateID.
func LoadAgentsYAML(path string, r *AgentRegistry) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read agents yaml %q: %w", path, err)
	}

	var file agentsYAMLFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("parse agents yaml %q: %w", path, err)
	}

	seen := make(map[string]struct{}, len(file.Agents))
	for _, entry := range file.Agents {
		if entry.ID == "" {
			return fmt.Errorf("agents yaml %q: agent entry missing id", path)
		}
		if _, dup := seen[entry.ID]; dup {
			return fmt.Errorf("agents yaml %q: duplicate id %q", path, entry.ID)
		}
		seen[entry.ID] = struct{}{}

		def := AgentDefinition{
			ID:           entry.ID,
			Name:         entry.Name,
			ModelProfile: entry.ModelProfile,
			Tools:        entry.Tools,
			ResourceProfile: ResourceProfile{
				MemoryLimit: entry.ResourceProfile.MemoryMB * 1024 * 1024,
				MaxPriority: entry.ResourceProfile.Priority,
			},
			Priority: entry.Priority,
		}
		if err := r.RegisterAgent(def); err != nil {
			return fmt.Errorf("register agent %q from %q: %w", entry.ID, path, err)
		}
	}

	return nil
}
