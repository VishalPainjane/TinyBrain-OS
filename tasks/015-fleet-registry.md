# Task 015 — Fleet Registry

## Status

In Progress

## Goal

Extend the agent registry with YAML loading capabilities to support a sample fleet of agents.

## Context

v0.8 and M7. Month 4 Week 14. Building on the agent plugin contract (Task 014), we need a way to define and load a fleet of agents (like Planner, Coder, etc.) dynamically rather than hardcoding them in the core runtime.

## Requirements

- Define `fleet.yaml` schema in `docs/contracts/registry.md`.
- Implement `LoadAgentsYAML(path string, r *AgentRegistry) error` in `internal/registry/agents_yaml.go`.
- Mirror the `LoadModelsYAML` design from v0.5.
- Parse `fleet.yaml` into the registry's `AgentDefinition` struct.
- Handle `ErrDuplicateID` gracefully.
- Provide a `testdata/fleet.yaml` with a sample fleet (e.g., `sample-alpha`, `sample-beta`).

## Files

- `tasks/015-fleet-registry.md`
- `docs/contracts/registry.md`
- `internal/registry/agents_yaml.go`
- `internal/registry/agents_yaml_test.go`
- `testdata/fleet.yaml`

## Acceptance Criteria

- [x] `fleet.yaml` schema documented
- [x] `LoadAgentsYAML` parses YAML into `AgentDefinition` structs
- [x] Duplicate IDs handled properly
- [x] `testdata/fleet.yaml` created with 2 sample agents
- [x] No hardcoded agent types introduced in core
- [ ] `go test ./...` passes

## Out Of Scope

- End-to-end event pipeline (Week 15)
- Actual execution of the fleet (Week 15+)

## Related

- Spec: [v0.8-agents.md](../docs/specs/v0.8-agents.md)
- Month plan: [month-04.md](../planning/roadmap/months/month-04.md)

---
**Layer:** task
**Target month:** 4
