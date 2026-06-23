# Task 014 — Agent Plugin

## Status

In Progress

## Goal

Implement generic agent plugin contract and one sample plugin (not a hardcoded core type).

## Context

v0.8 and M7. Month 4 Week 13. ADR-002.

## Requirements

- `internal/agents/` package with `Agent` interface and `RuntimeAPI` port (INV-003)
- `Executor` publishes `AgentStarted` / `AgentStopped` while running a plugin
- `SamplePlugin` demonstrates structured JSON output via `runtime.Generate`
- Contract documented in `docs/contracts/agents.md`

## Files

- `internal/agents/interface.go`
- `internal/agents/types.go`
- `internal/agents/executor.go`
- `internal/agents/sample.go`
- `internal/agents/executor_test.go`
- `docs/contracts/agents.md`
- `tests/import_boundary_test.go`

## Acceptance Criteria

- [x] `Agent` interface defined; no Planner/Coder types in core
- [x] Sample plugin runs one task in test via runtime API
- [x] Structured JSON output returned
- [x] Agent lifecycle events published on execution
- [x] `internal/agents` does not import `internal/inference`
- [x] `go test ./...` passes

## Out Of Scope

- Fixed Planner/Coder classes in core
- Sample fleet YAML (Week 14)
- End-to-end event pipeline (Week 15)

## Related

- Spec: [v0.8-agents.md](../docs/specs/v0.8-agents.md)
- Month plan: [month-04.md](../planning/roadmap/months/month-04.md)
- ADR-002

---
**Layer:** task
**Target month:** 4
