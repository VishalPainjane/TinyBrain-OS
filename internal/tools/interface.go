package tools

import "context"

// ToolSchema defines the JSON schema for tool arguments.
type ToolSchema struct {
	Type       string              `json:"type"` // always "object"
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property defines a field in the ToolSchema.
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// Tool represents a single executable tool, conforming to MCP standards.
type Tool interface {
	Name() string
	Description() string
	Schema() ToolSchema
	Execute(ctx context.Context, input string) (string, error)
}

// ToolRegistry manages a collection of tools and executes them.
type ToolRegistry interface {
	Register(tool Tool) error
	Get(name string) (Tool, bool)
	List() []Tool
	Execute(ctx context.Context, name, input string) (string, error)
}
