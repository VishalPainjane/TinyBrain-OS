package agents

import (
	"encoding/json"
	"fmt"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
)

// SamplePlugin is a generic reference plugin — not a hardcoded core agent type.
// It demonstrates structured JSON output via the runtime Generate API.
type SamplePlugin struct {
	id      string
	modelID string
}

// NewSamplePlugin returns a plugin bound to registry agent id and model profile.
func NewSamplePlugin(id, modelID string) *SamplePlugin {
	return &SamplePlugin{id: id, modelID: modelID}
}

// ID implements Agent.
func (p *SamplePlugin) ID() string {
	return p.id
}

// Execute implements Agent by calling runtime.Generate and marshaling JSON output.
func (p *SamplePlugin) Execute(ctx ExecuteContext, req TaskRequest) (TaskResult, error) {
	resp, err := ctx.Runtime.Generate(runtime.GenerateRequest{
		ModelID: p.modelID,
		Prompt:  req.Input,
	})
	if err != nil {
		return TaskResult{}, err
	}

	payload := map[string]any{
		"agent_id":        p.id,
		"task_id":         req.TaskID,
		"model_id":        resp.ModelID,
		"tokens_produced": resp.TokensProduced,
		"text":            resp.Output,
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return TaskResult{}, fmt.Errorf("marshal agent output: %w", err)
	}

	return TaskResult{
		TaskID:  req.TaskID,
		AgentID: p.id,
		Output:  string(raw),
	}, nil
}
