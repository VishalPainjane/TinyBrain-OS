package tools

import (
	"context"
	"fmt"
	"sync"
)

type localRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry creates a new in-memory ToolRegistry.
func NewRegistry() ToolRegistry {
	return &localRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a new tool to the registry.
func (r *localRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[tool.Name()]; exists {
		return fmt.Errorf("tool %s already registered", tool.Name())
	}
	r.tools[tool.Name()] = tool
	return nil
}

// Get retrieves a tool by name.
func (r *localRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

// List returns all registered tools.
func (r *localRegistry) List() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}

// Execute looks up a tool and executes it with the given input payload.
func (r *localRegistry) Execute(ctx context.Context, name, input string) (string, error) {
	tool, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("tool %s not found", name)
	}
	return tool.Execute(ctx, input)
}
