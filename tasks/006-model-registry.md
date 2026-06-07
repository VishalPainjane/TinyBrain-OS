# Task 006 — Model Registry

## Status

Complete

## Goal

Implement model definition registration and lookup.

## Context

Models are described by path, size, memory budget, and capabilities. Hardware profile will filter available models in v0.3+.

## Requirements

- ModelDefinition per contract
- RegisterModel, GetModel, ListModels
- In-memory implementation (persistence in v0.5)

## Files

- `internal/registry/models.go`
- `internal/registry/models_test.go`

## Acceptance Criteria

- [x] Register and retrieve model by ID
- [x] List returns all models with metadata
- [x] Duplicate ID returns error

## Out Of Scope

- BoltDB/SQLite persistence (v0.5)
- llama.cpp loading
- Hugging Face integration

## Related

- Spec: [docs/specs/v0.2-registry.md](../docs/specs/v0.2-registry.md), [v0.5-model-registry.md](../docs/specs/v0.5-model-registry.md)
- ADR-001

---
**Layer:** task
