package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

// runSubmit submits a task to the daemonized inference engine.
func runSubmit(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("submit", flag.ContinueOnError)
	fs.SetOutput(stderr)
	agentID := fs.String("agent", "qwen-coder", "agent ID to use")
	prompt := fs.String("prompt", "", "prompt text (required)")
	port := fs.Int("port", 8080, "daemon port")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *prompt == "" {
		fmt.Fprintln(stderr, "usage: tinybrain submit --agent ID --prompt TEXT")
		return 2
	}

	reqPayload := TaskRequest{
		AgentID: *agentID,
		Prompt:  *prompt,
	}

	body, err := json.Marshal(reqPayload)
	if err != nil {
		fmt.Fprintf(stderr, "error formatting request: %v\n", err)
		return 1
	}

	url := fmt.Sprintf("http://127.0.0.1:%d/v1/tasks", *port)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(stderr, "error connecting to daemon: %v\nIs the daemon running? Start it with 'tinybrain daemon'.\n", err)
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)
		fmt.Fprintf(stderr, "daemon returned error: %s\n", resp.Status)
		return 1
	}

	var result struct {
		TaskID  string  `json:"task_id"`
		Result  string  `json:"result"`
		Elapsed float64 `json:"elapsed"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Fprintf(stderr, "error reading response: %v\n", err)
		return 1
	}

	fmt.Fprintf(stdout, "\nOutput (Task: %s, Elapsed: %.2fs):\n%s\n", result.TaskID, result.Elapsed, result.Result)
	return 0
}
