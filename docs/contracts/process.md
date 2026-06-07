# Contract: Process

## Owner

Package: `internal/process/`

## Consumers

- `internal/scheduler/` — reads process table, requests state transitions
- `internal/runtime/` — updates resource usage fields
- Telemetry — read-only queries

## ProcessState

```go
// Documentation only — matches internal/process/state.go
type ProcessState int

const (
    New ProcessState = iota
    Ready
    Running
    Waiting
    Preempted
    Hibernated
    Terminated
)
```

Valid transitions (scheduler owns transition logic):

```text
NEW → READY → RUNNING ↔ WAITING
RUNNING → PREEMPTED → READY | HIBERNATED
HIBERNATED → READY
* → TERMINATED
```

## Process Record

| Field | Type | Description |
|-------|------|-------------|
| PID | string | Unique process identifier |
| AgentRef | string | Registry agent definition ID |
| State | ProcessState | Current lifecycle state |
| Priority | int | Scheduling priority |
| MemoryUsage | uint64 | RAM bytes |
| VRAMUsage | uint64 | VRAM bytes |
| KVCacheID | string | Associated KV block (optional) |
| LastExecution | time | Last active timestamp |
| TokensProduced | int | Token count this session |
| TaskID | string | Associated task |

## Process Table Operations

| Operation | Description |
|-----------|-------------|
| Create | Insert new process in NEW state |
| Get | Lookup by PID — O(1) target |
| List | Return all processes |
| UpdateState | Transition state (scheduler only) |
| Delete | Remove terminated process |

## Responsibilities

- Define and enforce ProcessState enum
- Store process records
- Provide O(1) lookup by PID

## Must NOT Own

- Scheduling policy or transition decisions
- Model loading
- Agent execution logic

## Related

- [architecture/kernel.md](../architecture/kernel.md)
- [tasks/001-process-state.md](../../tasks/001-process-state.md)
- [tasks/002-process-table.md](../../tasks/002-process-table.md)

---
**Layer:** contract
**Last verified against code:** 2026-06-07 (Month 1 foundation / tag v0.3)
