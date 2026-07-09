package agents

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/VishalPainjane/TinyBrain-OS/internal/runtime"
	"github.com/VishalPainjane/TinyBrain-OS/internal/tools"
)

// ReActAgent implements the explicit Plan -> Call -> Observe -> Reflect -> Writeback loop.
type ReActAgent struct {
	id       string
	modelID  string
	registry tools.ToolRegistry
}

// NewReActAgent creates a ReAct-based agent with access to a specific ToolRegistry.
func NewReActAgent(id, modelID string, registry tools.ToolRegistry) *ReActAgent {
	return &ReActAgent{
		id:       id,
		modelID:  modelID,
		registry: registry,
	}
}

// ID implements Agent.
func (p *ReActAgent) ID() string {
	return p.id
}

type ToolCall struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

// generateToolCallGrammar is a GBNF grammar that forces the model to output a valid JSON array of tool calls.
// In a full implementation, this is generated dynamically based on the tool schema.
const toolCallGrammar = `
root ::= "[" ws call ("," ws call)* "]" ws
call ::= "{" ws "\"name\"" ws ":" ws string "," ws "\"arguments\"" ws ":" ws arguments "}"
arguments ::= "{" ws (string ":" ws string ("," ws string ":" ws string)*)? "}"
string ::= "\"" [a-zA-Z0-9_]* "\"" ws
ws ::= [ \t\n]*
`

// Execute runs the full state machine loop: Plan -> Call -> Observe -> Reflect -> Writeback
func (p *ReActAgent) Execute(ctx ExecuteContext, req TaskRequest) (TaskResult, error) {
	err := ctx.Runtime.LoadModel(p.modelID)
	if err != nil && !errors.Is(err, runtime.ErrModelAlreadyLoaded) {
		return TaskResult{}, fmt.Errorf("agent %s failed to load model %s: %w", p.id, p.modelID, err)
	}

	messages := []runtime.ChatMessage{
		{Role: "system", Content: "You are a helpful AI assistant with access to tools. If you need to use a tool, output a JSON array of tool calls. Otherwise, answer the user directly."},
	}
	
	// Attempt to parse input as full chat history from daemon
	var history []runtime.ChatMessage
	if err := json.Unmarshal([]byte(req.Input), &history); err == nil && len(history) > 0 {
		messages = append(messages, history...)
	} else {
		// Fallback for raw CLI strings
		messages = append(messages, runtime.ChatMessage{Role: "user", Content: req.Input})
	}

	maxTurns := 5
	finalOutput := ""

	for turn := 0; turn < maxTurns; turn++ {
		// 1. Plan & Format
		systemPrompt := `You are an AI with tools.
Available Tools:`
		
		toolsList := p.registry.List()
		for _, t := range toolsList {
			schemaBytes, _ := json.Marshal(t.Schema())
			systemPrompt += fmt.Sprintf("\n- %s: %s\n  Schema: %s", t.Name(), t.Description(), string(schemaBytes))
		}

		systemPrompt += `

To use a tool, you MUST output EXACTLY a JSON array of objects.
Do not use markdown blocks. Output raw JSON.
Example:
[{"name": "calculator", "arguments": {"a": "15", "b": "20"}}]

If you have the final answer, output plain text, not JSON.`
		messages[0].Content = systemPrompt

		prompt, _, err := ctx.Runtime.FormatChat(p.modelID, messages, runtime.FormatChatOpts{
			AddGenerationPrompt: true,
		})
		if err != nil {
			prompt = req.Input // primitive fallback
		}

		// 2. Call
		resp, err := ctx.Runtime.Generate(runtime.GenerateRequest{
			ModelID:   p.modelID,
			Prompt:    prompt,
			MaxTokens: 4096,
			// Grammar: toolCallGrammar, // Enforce strict JSON array schema
		})
		if err != nil {
			return TaskResult{}, err
		}

		output := strings.TrimSpace(resp.Output)
		output = strings.ReplaceAll(output, "<|im_end|>", "")
		output = strings.TrimSpace(output)
		
		fmt.Printf("\n--- TURN %d ---\n%s\n----------------\n", turn, output)

		// 3. Extract Tool Calls using Regex
		var calls []ToolCall
		
		// Regex to find {"name": "tool_name", "arguments": {"key": "value"}}
		// It handles optional whitespace and extracts the name and arguments object.
		re := regexp.MustCompile(`(?s)\{\s*"name"\s*:\s*"([^"]+)"\s*,\s*"arguments"\s*:\s*(\{.*?\})\s*\}`)
		matches := re.FindAllStringSubmatch(output, -1)

		if len(matches) == 0 {
			// If it has JSON-like characters but failed the regex, reflect!
			if strings.Contains(output, "{") || strings.Contains(output, "[") {
				messages = append(messages, runtime.ChatMessage{Role: "assistant", Content: output})
				messages = append(messages, runtime.ChatMessage{Role: "system", Content: "You attempted to use a tool, but your JSON format was invalid. You MUST output EXACTLY: [{\"name\": \"tool_name\", \"arguments\": {\"key\": \"value\"}}]"})
				continue
			}

			// No tool calls found, treat as writeback
			finalOutput = output
			break
		}

		for _, match := range matches {
			toolName := match[1]
			argsRaw := match[2]

			var argsMap map[string]interface{}
			if err := json.Unmarshal([]byte(argsRaw), &argsMap); err != nil {
				// Syntax error in arguments: Reflect and try again
				messages = append(messages, runtime.ChatMessage{Role: "assistant", Content: output})
				messages = append(messages, runtime.ChatMessage{Role: "system", Content: fmt.Sprintf("Error parsing arguments for tool %s: %v", toolName, err)})
				continue
			}

			strMap := make(map[string]string)
			for k, v := range argsMap {
				strMap[k] = fmt.Sprintf("%v", v)
			}

			calls = append(calls, ToolCall{
				Name:      toolName,
				Arguments: strMap,
			})
		}

		// If we had syntax errors and continued, calls would be empty here on this iteration
		if len(calls) == 0 {
			continue
		}

		messages = append(messages, runtime.ChatMessage{Role: "assistant", Content: output})

		for _, call := range calls {
			argsBytes, _ := json.Marshal(call.Arguments)
			resultStr, execErr := p.registry.Execute(context.Background(), call.Name, string(argsBytes))
			if execErr != nil {
				resultStr = fmt.Sprintf("Error executing %s: %v", call.Name, execErr)
			}
			// Append Observation
			messages = append(messages, runtime.ChatMessage{
				Role:    "tool",
				Content: fmt.Sprintf("Observation from tool '%s': %s", call.Name, resultStr),
			})
		}
		// Loop back to Plan -> Call based on new Observations
	}

	if finalOutput == "" {
		finalOutput = "Error: Agent exceeded maximum turns without reaching a final answer."
	}

	payload := map[string]any{
		"agent_id": p.id,
		"task_id":  req.TaskID,
		"model_id": p.modelID,
		"text":     finalOutput,
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
