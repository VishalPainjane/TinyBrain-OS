package events

import "time"

// Type identifies a core lifecycle event on the event bus.
// See ADR-003 and tasks/003-event-types.md.
type Type string

const (
	TypeTaskCreated         Type = "TaskCreated"
	TypeTaskAssigned        Type = "TaskAssigned"
	TypeProcessSpawned      Type = "ProcessSpawned"
	TypeProcessStateChanged Type = "ProcessStateChanged"
	TypeAgentStarted        Type = "AgentStarted"
	TypeAgentStopped        Type = "AgentStopped"
	TypeModelLoaded         Type = "ModelLoaded"
	TypeModelUnloaded       Type = "ModelUnloaded"
	TypeSwapStarted         Type = "SwapStarted"
	TypeSwapCompleted       Type = "SwapCompleted"
	TypeKVStored            Type = "KVStored"
	TypeKVLoaded            Type = "KVLoaded"
	TypeKVCompressed        Type = "KVCompressed"
	TypeKVDecompressed      Type = "KVDecompressed"
	TypeTaskCompleted       Type = "TaskCompleted"
)

// Event is a typed lifecycle message published on the event bus.
type Event struct {
	Type      Type
	Timestamp time.Time
	Payload   any
}

// TaskCreatedPayload is the payload for TypeTaskCreated.
type TaskCreatedPayload struct {
	TaskID  string
	Input   string
	AgentID string
}

// TaskAssignedPayload is the payload for TypeTaskAssigned.
type TaskAssignedPayload struct {
	TaskID  string
	AgentID string
}

// ProcessSpawnedPayload is the payload for TypeProcessSpawned.
type ProcessSpawnedPayload struct {
	PID      string
	AgentRef string
	TaskID   string
}

// ProcessStateChangedPayload is the payload for TypeProcessStateChanged.
type ProcessStateChangedPayload struct {
	PID       string
	OldState  string
	NewState  string
}

// AgentStartedPayload is the payload for TypeAgentStarted.
type AgentStartedPayload struct {
	AgentID string
	PID     string
}

// AgentStoppedPayload is the payload for TypeAgentStopped.
type AgentStoppedPayload struct {
	AgentID string
	PID     string
}

// ModelLoadedPayload is the payload for TypeModelLoaded.
type ModelLoadedPayload struct {
	ModelID string
}

// ModelUnloadedPayload is the payload for TypeModelUnloaded.
type ModelUnloadedPayload struct {
	ModelID string
}

// SwapStartedPayload is the payload for TypeSwapStarted.
type SwapStartedPayload struct {
	FromModelID string
	ToModelID   string
}

// SwapCompletedPayload is the payload for TypeSwapCompleted.
type SwapCompletedPayload struct {
	FromModelID string
	ToModelID   string
}

// KVStoredPayload is the payload for TypeKVStored.
type KVStoredPayload struct {
	KVCacheID string
	PID       string
}

// KVLoadedPayload is the payload for TypeKVLoaded.
type KVLoadedPayload struct {
	KVCacheID string
	PID       string
}

// KVCompressedPayload is the payload for TypeKVCompressed.
type KVCompressedPayload struct {
	KVCacheID        string
	CompressionRatio float64
	LatencyMs        int64
}

// KVDecompressedPayload is the payload for TypeKVDecompressed.
type KVDecompressedPayload struct {
	KVCacheID string
	LatencyMs int64
}

// TaskCompletedPayload is the payload for TypeTaskCompleted.
type TaskCompletedPayload struct {
	TaskID string
	Result string
}

// NewEvent creates an event with the given type, payload, and timestamp.
func NewEvent(eventType Type, payload any, at time.Time) Event {
	return Event{
		Type:      eventType,
		Timestamp: at,
		Payload:   payload,
	}
}

// AllTypes returns every core event type defined for v1.
func AllTypes() []Type {
	return []Type{
		TypeTaskCreated,
		TypeTaskAssigned,
		TypeProcessSpawned,
		TypeProcessStateChanged,
		TypeAgentStarted,
		TypeAgentStopped,
		TypeModelLoaded,
		TypeModelUnloaded,
		TypeSwapStarted,
		TypeSwapCompleted,
		TypeKVStored,
		TypeKVLoaded,
		TypeKVCompressed,
		TypeKVDecompressed,
		TypeTaskCompleted,
	}
}
