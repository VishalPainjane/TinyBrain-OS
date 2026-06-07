package registry_test

import (
	"errors"
	"testing"

	"github.com/VishalPainjane/TinyBrain-OS/internal/registry"
)

func TestAgentRegistry_RegisterAndGet(t *testing.T) {
	reg := registry.NewAgentRegistry()
	def := registry.AgentDefinition{
		ID:           "agent-alpha",
		Name:         "Alpha Agent",
		ModelProfile: "model-small",
		Tools:        []string{"search"},
		ResourceProfile: registry.ResourceProfile{
			MemoryLimit: 1_000_000_000,
			MaxPriority: 10,
		},
		Priority: 5,
	}

	if err := reg.RegisterAgent(def); err != nil {
		t.Fatalf("RegisterAgent() error = %v", err)
	}

	got, err := reg.GetAgent("agent-alpha")
	if err != nil {
		t.Fatalf("GetAgent() error = %v", err)
	}
	if got.Name != def.Name {
		t.Errorf("Name = %q, want %q", got.Name, def.Name)
	}
}

func TestAgentRegistry_ListAgents(t *testing.T) {
	reg := registry.NewAgentRegistry()

	for _, id := range []string{"a1", "a2"} {
		if err := reg.RegisterAgent(registry.AgentDefinition{ID: id, Name: id}); err != nil {
			t.Fatalf("RegisterAgent(%s) error = %v", id, err)
		}
	}

	list := reg.ListAgents()
	if len(list) != 2 {
		t.Fatalf("ListAgents() len = %d, want 2", len(list))
	}
}

func TestAgentRegistry_DuplicateID(t *testing.T) {
	reg := registry.NewAgentRegistry()
	def := registry.AgentDefinition{ID: "dup", Name: "one"}

	if err := reg.RegisterAgent(def); err != nil {
		t.Fatalf("first RegisterAgent() error = %v", err)
	}
	if err := reg.RegisterAgent(def); !errors.Is(err, registry.ErrDuplicateID) {
		t.Errorf("second RegisterAgent() error = %v, want ErrDuplicateID", err)
	}
}

func TestAgentRegistry_GetNotFound(t *testing.T) {
	reg := registry.NewAgentRegistry()

	_, err := reg.GetAgent("missing")
	if !errors.Is(err, registry.ErrNotFound) {
		t.Errorf("GetAgent() error = %v, want ErrNotFound", err)
	}
}
