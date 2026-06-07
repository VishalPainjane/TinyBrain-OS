# Task 001 — Process States

## Status

Complete

## Goal

Implement the lifecycle states of a TinyBrain process.

## Context

Every agent running inside TinyBrain is represented as a process. The scheduler will operate on process states. Agent names (planner, coder, etc.) are example fleet configurations — not core types.

## Requirements

Create a ProcessState type with states:

- NEW
- READY
- RUNNING
- WAITING
- PREEMPTED
- HIBERNATED
- TERMINATED

## Files

- `internal/process/state.go`
- `internal/process/state_test.go`

## Acceptance Criteria

- [x] ProcessState type exists
- [x] All 7 states are defined
- [x] States are strongly typed
- [x] Unit tests exist and pass

## Out Of Scope

- Scheduler
- Runtime
- Model loading
- Event bus
- Process table (task 002)

## Related

- Spec: [docs/specs/v0.1-kernel.md](../docs/specs/v0.1-kernel.md)
- Contract: [docs/contracts/process.md](../docs/contracts/process.md)
- Completed: [planning/execution/completed.md](../planning/execution/completed.md)

---
**Layer:** task
