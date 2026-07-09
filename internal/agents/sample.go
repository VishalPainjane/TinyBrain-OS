package agents

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	err := ctx.Runtime.LoadModel(p.modelID)
	if err != nil && !errors.Is(err, runtime.ErrModelAlreadyLoaded) {
		return TaskResult{}, fmt.Errorf("agent %s failed to load model %s: %w", p.id, p.modelID, err)
	}

	messages := []runtime.ChatMessage{
		{Role: "system", Content: "You are a helpful AI assistant."},
	}
	var history []runtime.ChatMessage
	if err := json.Unmarshal([]byte(req.Input), &history); err == nil && len(history) > 0 {
		messages = append(messages, history...)
	} else {
		messages = append(messages, runtime.ChatMessage{Role: "user", Content: req.Input})
	}

	prompt, _, err := ctx.Runtime.FormatChat(p.modelID, messages, runtime.FormatChatOpts{
		AddGenerationPrompt: true,
	})
	if err != nil {
		// Fallback to raw input if formatting fails or is unsupported
		prompt = req.Input
	}

	resp, err := ctx.Runtime.Generate(runtime.GenerateRequest{
		ModelID:   p.modelID,
		Prompt:    prompt,
		MaxTokens: 4096,
	})
	if err != nil {
		return TaskResult{}, err
	}

	outText := resp.Output
	outText = strings.ReplaceAll(outText, "<|im_end|>", "")
	outText = strings.TrimSpace(outText)

	payload := map[string]any{
		"agent_id":        p.id,
		"task_id":         req.TaskID,
		"model_id":        resp.ModelID,
		"tokens_produced": resp.TokensProduced,
		"text":            outText,
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
