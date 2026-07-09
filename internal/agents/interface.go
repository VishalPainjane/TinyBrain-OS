package agents

import "github.com/VishalPainjane/TinyBrain-OS/internal/runtime"

// Agent is the plugin contract for registry-defined workers.
// Core defines the interface only — not fixed agent types (ADR-002).
// See docs/architecture/agents.md and docs/contracts/agents.md.
type Agent interface {
	// ID returns the registry agent definition ID.
	ID() string
	// Execute runs one task using the runtime API (INV-003).
	Execute(ctx ExecuteContext, req TaskRequest) (TaskResult, error)
}

// RuntimeAPI is the narrow runtime port agents may call.
// Agents must not import inference or call InferenceProvider directly.
type RuntimeAPI interface {
	Generate(req runtime.GenerateRequest) (runtime.GenerateResponse, error)
	LoadModel(modelID string) error
	FormatChat(modelID string, messages []runtime.ChatMessage, opts runtime.FormatChatOpts) (string, string, error)
	GetMetadata(modelID string) (runtime.ModelCapabilities, error)
}

// ExecuteContext carries dependencies injected for a single execution.
type ExecuteContext struct {
	Runtime RuntimeAPI
}
