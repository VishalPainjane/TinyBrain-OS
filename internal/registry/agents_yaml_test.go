package registry

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAgentsYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		content := `
agents:
  - id: sample-alpha
    name: "Sample Alpha"
    model_profile: "tinyllama-q4"
    tools:
      - "weather"
    resource_profile:
      memory_mb: 2048
      priority: 1
    priority: 10
  - id: sample-beta
    name: "Sample Beta"
    model_profile: "tinyllama-q8"
    tools: []
    resource_profile:
      memory_mb: 4096
      priority: 2
    priority: 5
`
		path := filepath.Join(t.TempDir(), "fleet.yaml")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test yaml: %v", err)
		}

		reg := NewAgentRegistry()
		err := LoadAgentsYAML(path, reg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		alpha, err := reg.GetAgent("sample-alpha")
		if err != nil {
			t.Fatalf("failed to get sample-alpha: %v", err)
		}
		if alpha.Name != "Sample Alpha" {
			t.Errorf("expected name 'Sample Alpha', got %q", alpha.Name)
		}
		if alpha.ResourceProfile.MemoryLimit != 2048*1024*1024 {
			t.Errorf("expected memory 2147483648, got %d", alpha.ResourceProfile.MemoryLimit)
		}

		beta, err := reg.GetAgent("sample-beta")
		if err != nil {
			t.Fatalf("failed to get sample-beta: %v", err)
		}
		if beta.Priority != 5 {
			t.Errorf("expected priority 5, got %d", beta.Priority)
		}
	})

	t.Run("duplicate_id_in_file", func(t *testing.T) {
		content := `
agents:
  - id: dupe
    name: "One"
  - id: dupe
    name: "Two"
`
		path := filepath.Join(t.TempDir(), "fleet.yaml")
		os.WriteFile(path, []byte(content), 0644)

		reg := NewAgentRegistry()
		err := LoadAgentsYAML(path, reg)
		if err == nil {
			t.Fatal("expected error on duplicate ID, got nil")
		}
	})

	t.Run("duplicate_id_in_registry", func(t *testing.T) {
		content := `
agents:
  - id: existing
    name: "Two"
`
		path := filepath.Join(t.TempDir(), "fleet.yaml")
		os.WriteFile(path, []byte(content), 0644)

		reg := NewAgentRegistry()
		reg.RegisterAgent(AgentDefinition{ID: "existing", Name: "One"})

		err := LoadAgentsYAML(path, reg)
		if err == nil {
			t.Fatal("expected error on duplicate ID, got nil")
		}
	})

	t.Run("missing_id", func(t *testing.T) {
		content := `
agents:
  - name: "No ID"
`
		path := filepath.Join(t.TempDir(), "fleet.yaml")
		os.WriteFile(path, []byte(content), 0644)

		reg := NewAgentRegistry()
		err := LoadAgentsYAML(path, reg)
		if err == nil {
			t.Fatal("expected error on missing ID, got nil")
		}
	})
}
