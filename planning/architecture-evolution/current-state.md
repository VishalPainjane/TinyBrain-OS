# Architecture — Current State

Code-first snapshot. Update after every completed task.

**Last verified:** 2026-06-07

## Implemented

| Component | Package | Status |
|-----------|---------|--------|
| ProcessState (7 states) | `internal/process/state.go` | Complete, tested |
| Process table (CRUD, O(1) lookup) | `internal/process/table.go` | Complete, tested |
| Event types (13 core events) | `internal/events/types.go` | Complete, tested |
| Event bus (channel pub/sub) | `internal/events/bus.go` | Complete, tested |
| Agent registry | `internal/registry/agents.go` | Complete, tested |
| Model registry | `internal/registry/models.go` | Complete, tested |
| Hardware probe + profiles | `internal/hardware/` | Complete, tested |

## In Progress

| Component | Task | Status |
|-----------|------|--------|
| Runtime shell | 008-runtime | Not started (Month 2) |

## Not Implemented

- Model loader, persistence
- Scheduler
- Inference provider / llama.cpp adapter
- Agent plugins, tool registry
- Telemetry, brain-top TUI
- KV manager, swap manager
- Kubernetes operator

## Active Packages

`internal/process/`, `internal/events/`, `internal/registry/`, `internal/hardware/`

## Tests

`go test ./...` — passing (process, events, registry, hardware)

## Governance

Repository rules frozen in `.cursor/rules/` (7 files). CI via `.github/workflows/ci.yml`.

---
**Layer:** planning
**Related:** [future-state.md](future-state.md), [../../docs/current.md](../../docs/current.md)
