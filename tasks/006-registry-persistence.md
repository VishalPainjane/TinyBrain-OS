# Task 006-registry-persistence — Model Registry Persistence

## Status

Complete

## Goal

Add durable storage for model definitions so metadata survives process restart, without changing the `ModelRegistry` public API established in task 006.

## Context

Task [006-model-registry](006-model-registry.md) shipped an in-memory `ModelRegistry` at v0.2/v0.3. [v0.5-model-registry.md](../docs/specs/v0.5-model-registry.md) requires persistence before v0.6 inference can rely on stable model paths and memory budgets. [migration-paths.md](../planning/architecture-evolution/migration-paths.md) Path 3 documents the adapter-swap approach.

Execution was incorrectly pointed at [011-kv-manager](011-kv-manager.md) (Month 3 memory subsystem). This task realigns V0.5 with the approved planning assessment.

## Requirements

1. **Persistence backend** — **bbolt** (`go.etcd.io/bbolt`); gob-encoded values in bucket keyed by model ID (approved 2026-06-08)
2. **API preservation** — `RegisterModel`, `GetModel`, `ListModels` signatures and semantics unchanged from task 006
3. **Durability** — Writes flush to disk; reopening the store returns the same definitions
4. **Seed format** — `models.yaml` file loads definitions at registry open (merge or bootstrap; document behavior)
5. **Error handling** — Duplicate ID, missing ID, empty ID, corrupt store — wrap with context; no silent swallowing
6. **Concurrency** — Thread-safe reads/writes matching in-memory registry behavior

## Files

| File | Purpose |
|------|---------|
| `internal/registry/models.go` | Extract shared types; keep or refactor in-memory constructor |
| `internal/registry/models_memory.go` | In-memory implementation (extracted from current `models.go`) |
| `internal/registry/models_store.go` | `ModelStore` interface |
| `internal/registry/models_memory.go` | `InMemoryStore` |
| `internal/registry/models_bolt.go` | `BboltStore` (gob values in bbolt bucket) |
| `internal/registry/models_yaml.go` | `models.yaml` parse and seed loader |
| `internal/registry/models_*_test.go` | Table-driven tests per existing `models_test.go` pattern |
| `testdata/models.yaml` | Fixture for seed-format tests |

`models.yaml` parsing uses `gopkg.in/yaml.v3` (seed loader only; not a persistence backend).

## Acceptance Criteria

- [x] `RegisterModel` / `GetModel` / `ListModels` behave identically to task 006 for in-memory and persistent implementations
- [x] Model definitions survive simulated process restart (close store, reopen, `ListModels` returns same entries)
- [x] `ListModels` returns path, `SizeBytes`, and `MemoryBudget` for every registered model
- [x] `models.yaml` seed file populates registry on first open; format documented in task or contract addendum
- [x] Duplicate ID returns `ErrDuplicateID`; missing ID returns `ErrNotFound`; empty ID returns error
- [x] `go test ./internal/registry/...` passes
- [x] No new imports in `internal/scheduler/`, `internal/runtime/`, or `internal/loader/` (registry-only change)
- [x] BoltDB vs SQLite decision recorded in [accepted.md](../planning/decisions/accepted.md)

## Out Of Scope

- Agent registry persistence (agents remain in-memory)
- Tool registry
- Hugging Face integration or automatic model download
- Hardware-profile filtering of list results (optional v0.3 enhancement — defer)
- KV manager ([011-kv-manager](011-kv-manager.md))
- llama.cpp / GGUF loading (v0.6)
- Unified `Registry` facade ([RFC-004](../docs/rfc/RFC-004-Registry-Facade.md))

## Implementation Plan

**Status:** Approved 2026-06-08.

### Step 0 — Decision gate ✓

- bbolt backend logged in `planning/decisions/accepted.md`
- Architecture compliance recorded in architecture review

### Step 1 — Interface extraction (Path 3, step 1)

- `ModelStore` interface: `RegisterModel`, `GetModel`, `ListModels`, `Close`
- `InMemoryStore` — extracted map + mutex from task 006
- `ModelRegistry` delegates to `ModelStore`; `NewModelRegistry()` uses `InMemoryStore`
- Existing tests pass unchanged

### Step 2 — BboltStore (Path 3, step 3)

- `BboltStore` opens `go.etcd.io/bbolt` file; bucket `models`; key = ID, value = gob(`ModelDefinition`)
- `NewBboltModelRegistry(dbPath, seedPath)` factory
- Table-driven tests: register → close → reopen → verify

### Step 3 — `models.yaml` seed

- Define minimal YAML schema:

```yaml
models:
  - id: tinyllama-q4
    path: /models/tinyllama-q4.gguf
    size_bytes: 637534208
    memory_budget: 2147483648
    quantization: Q4_K_M
    capabilities: [chat]
```

- `LoadModelsYAML(path string, store ModelStore) error` — skip or error on duplicate IDs (document: recommend error)
- Integration test: seed file → list → restart → list matches

### Step 4 — Factory and defaults

- `NewModelRegistry()` remains in-memory for unit tests
- `NewPersistentModelRegistry(dbPath string, seedPath string) (*ModelRegistry, error)` or equivalent — single entry point for runtime use at v0.6
- Document env var or config hook if needed (no secrets in repo)

### Step 5 — Verification and ship

- Run `go test ./...`
- Demo: register models → exit → restart → list (document in `planning/releases/v0.5.md`)
- Tier A docs: `completed.md`, `current-sprint.md`, `docs/current.md`
- Tier B on ship: tag `v0.5`, update `planning/releases/v0.5.md`

## Related

- Spec: [docs/specs/v0.5-model-registry.md](../docs/specs/v0.5-model-registry.md)
- Contract: [docs/contracts/registry.md](../docs/contracts/registry.md)
- Prior task: [006-model-registry.md](006-model-registry.md)
- Architecture review: [planning/decisions/006-registry-persistence-architecture-review.md](../planning/decisions/006-registry-persistence-architecture-review.md)
- Migration path: [planning/architecture-evolution/migration-paths.md](../planning/architecture-evolution/migration-paths.md) Path 3
- Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)
- Month plan: [planning/roadmap/months/month-02.md](../planning/roadmap/months/month-02.md) Week 6

---
**Layer:** task
**Target version:** v0.5
**Target month:** 2
