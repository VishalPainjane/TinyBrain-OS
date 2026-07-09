package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// CalculatorTool is a simple mock tool that adds two numbers.
type CalculatorTool struct{}

func (c *CalculatorTool) Name() string {
	return "calculator"
}

func (c *CalculatorTool) Description() string {
	return "Adds two numbers together. Provide arguments as a JSON string with 'a' and 'b'."
}

func (c *CalculatorTool) Schema() ToolSchema {
	return ToolSchema{
		Type: "object",
		Properties: map[string]Property{
			"a": {Type: "string", Description: "First number"},
			"b": {Type: "string", Description: "Second number"},
		},
		Required: []string{"a", "b"},
	}
}

func (c *CalculatorTool) Execute(ctx context.Context, input string) (string, error) {
	var args map[string]string
	if err := json.Unmarshal([]byte(input), &args); err != nil {
		// Fallback: If input is not JSON, try returning a dummy result
		return "Invalid arguments format. Expected JSON like {\"a\":\"1\", \"b\":\"2\"}", nil
	}

	a, _ := strconv.Atoi(args["a"])
	b, _ := strconv.Atoi(args["b"])
	return fmt.Sprintf("%d", a+b), nil
}
