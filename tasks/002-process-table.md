# Task 002 — Process Table

## Status

Complete

## Goal

Implement the process table with CRUD operations and O(1) PID lookup.

## Context

The kernel needs a process table analogous to Linux — tracking every running or scheduled agent process. brain-top and the scheduler will read from this table.

## Requirements

- Process struct per [docs/contracts/process.md](../docs/contracts/process.md)
- ProcessTable with Create, Get, List, UpdateState, Delete
- O(1) lookup by PID (map-backed)
- Thread-safe if concurrent access anticipated (mutex)

## Files

- `internal/process/table.go`
- `internal/process/table_test.go`
- `internal/process/types.go` (Process struct if separated)

## Acceptance Criteria

- [x] Create inserts process in NEW state
- [x] Get returns process by PID; error if not found
- [x] List returns all processes
- [x] UpdateState changes state; rejects invalid transitions optionally (or defer to scheduler)
- [x] Delete removes terminated process
- [x] Lookup is O(1)
- [x] Unit tests cover all operations

## Out Of Scope

- Scheduler transition logic
- Runtime integration
- Persistence across restarts

## Related

- Spec: [docs/specs/v0.1-kernel.md](../docs/specs/v0.1-kernel.md)
- Contract: [docs/contracts/process.md](../docs/contracts/process.md)
- Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)

---
**Layer:** task
