# Task 005 — Agent Registry

## Status

Complete

## Goal

Implement agent definition registration and lookup.

## Context

Agents are plugins (ADR-002). The registry stores agent definitions — not hardcoded agent types.

## Requirements

- Implement Registry interface from [docs/contracts/registry.md](../docs/contracts/registry.md)
- AgentDefinition: ID, Name, ModelProfile, Tools, ResourceProfile, Priority
- RegisterAgent, GetAgent, ListAgents
- In-memory map implementation

## Files

- `internal/registry/agents.go`
- `internal/registry/agents_test.go`

## Acceptance Criteria

- [x] Register and retrieve agent by ID
- [x] List returns all registered agents
- [x] Duplicate ID returns error
- [x] No hardcoded agent type names in package logic

## Out Of Scope

- Model registry (task 006)
- Persistence (v0.5)
- Hardware-aware filtering (v0.3)

## Related

- Spec: [docs/specs/v0.2-registry.md](../docs/specs/v0.2-registry.md)
- ADR-002

---
**Layer:** task
