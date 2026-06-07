package registry

import (
	"errors"
	"fmt"
	"sync"
)

// ErrDuplicateID is returned when registering a definition with an existing ID.
var ErrDuplicateID = errors.New("duplicate ID")

// ErrNotFound is returned when a definition ID does not exist.
var ErrNotFound = errors.New("not found")

// ResourceProfile describes memory and priority limits for an agent.
// See docs/contracts/registry.md.
type ResourceProfile struct {
	MemoryLimit uint64
	MaxPriority int
}

// AgentDefinition describes a plugin agent entry in the registry.
// See docs/contracts/registry.md.
type AgentDefinition struct {
	ID              string
	Name            string
	ModelProfile    string
	Tools           []string
	ResourceProfile ResourceProfile
	Priority        int
}

// AgentRegistry stores and serves agent definitions in memory.
type AgentRegistry struct {
	mu     sync.RWMutex
	agents map[string]AgentDefinition
}

// NewAgentRegistry returns an empty in-memory agent registry.
func NewAgentRegistry() *AgentRegistry {
	return &AgentRegistry{
		agents: make(map[string]AgentDefinition),
	}
}

// RegisterAgent inserts an agent definition. Duplicate IDs return ErrDuplicateID.
func (r *AgentRegistry) RegisterAgent(def AgentDefinition) error {
	if def.ID == "" {
		return fmt.Errorf("agent ID is required")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.agents[def.ID]; exists {
		return ErrDuplicateID
	}

	r.agents[def.ID] = def
	return nil
}

// GetAgent returns the agent definition for id.
func (r *AgentRegistry) GetAgent(id string) (AgentDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.agents[id]
	if !ok {
		return AgentDefinition{}, ErrNotFound
	}
	return def, nil
}

// ListAgents returns all registered agent definitions.
func (r *AgentRegistry) ListAgents() []AgentDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]AgentDefinition, 0, len(r.agents))
	for _, def := range r.agents {
		out = append(out, def)
	}
	return out
}
