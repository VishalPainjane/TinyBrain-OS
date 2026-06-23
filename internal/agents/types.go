package agents

// TaskRequest is structured input for one agent execution.
type TaskRequest struct {
	TaskID string
	PID    string
	Input  string
}

// TaskResult is structured JSON output from an agent execution.
type TaskResult struct {
	TaskID  string
	AgentID string
	Output  string
}
