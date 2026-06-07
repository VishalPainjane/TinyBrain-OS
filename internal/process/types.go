package process

import "time"

// Process is a record in the process table.
// See docs/contracts/process.md.
type Process struct {
	PID            string
	AgentRef       string
	State          ProcessState
	Priority       int
	MemoryUsage    uint64
	VRAMUsage      uint64
	KVCacheID      string
	LastExecution  time.Time
	TokensProduced int
	TaskID         string
}
