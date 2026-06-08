# Architecture Review — Task 006-registry-persistence (V0.5)

**Status:** Approved — implementation in progress  
**Date:** 2026-06-08  
**Approved:** 2026-06-08  
**Scope:** Persistent model metadata only — no KV manager, no inference, no agent persistence

## Summary

V0.5 adds a file-backed `ModelStore` adapter behind the existing `ModelRegistry` API. Callers (runtime at v0.6, loader) continue to use `RegisterModel` / `GetModel` / `ListModels`. In-memory implementation remains for tests. This follows [migration-paths.md](../architecture-evolution/migration-paths.md) Path 3.

## Problem Statement

Task 006 shipped a mutex-protected `map[string]ModelDefinition`. Definitions are lost on process exit. v0.6 inference needs durable paths and memory budgets across restarts. Execution docs incorrectly targeted task 011 (KV manager), which belongs to Month 3 memory subsystems after v0.7 MLFQ.

## Architectural Invariants

| Invariant | Impact on this task |
|-----------|---------------------|
| INV-004 (registry owns capability definitions) | Persistence stays in `internal/registry/` only |
| Registry must not execute or schedule | No runtime/scheduler imports added to registry |
| Adapter swap, not rewrite | Public method signatures unchanged |
| No scope creep | Agents, tools, KV — out of scope |
| Dependency approval required | New store dependency needs explicit sign-off |

## Current State

```text
internal/registry/
  models.go      — ModelDefinition + in-memory ModelRegistry
  models_test.go — register/get/list/duplicate/not-found
  agents.go      — separate in-memory AgentRegistry (unchanged)
```

`ModelRegistry` is a concrete struct, not yet behind an interface. Migration Path 3 step 1 (interface) was planned but not implemented — **this task completes steps 1 and 3**.

## Proposed Design

### Layering

```text
┌─────────────────────────────────────┐
│  Callers (future: runtime, loader)  │
└─────────────────┬───────────────────┘
                  │ RegisterModel / GetModel / ListModels
┌─────────────────▼───────────────────┐
│  ModelRegistry (facade, optional)   │
│  or ModelStore interface            │
└─────────────────┬───────────────────┘
        ┌─────────┴─────────┐
        ▼                   ▼
 InMemoryStore         BboltStore
 (tests, dev)          (go.etcd.io/bbolt file)
        │                   │
        └─────────┬─────────┘
                  ▼
           models.yaml seed (bootstrap only)
```

### Interface shape (documentation target)

```go
// ModelStore persists model definitions. See docs/contracts/registry.md.
type ModelStore interface {
    RegisterModel(def ModelDefinition) error
    GetModel(id string) (ModelDefinition, error)
    ListModels() []ModelDefinition
}
```

`ModelRegistry` either embeds `ModelStore` or becomes a type alias — implementation detail left to coding phase. Tests must prove behavioral parity with task 006.

### Persistence backend recommendation

| Option | Pros | Cons |
|--------|------|------|
| **bbolt** (`go.etcd.io/bbolt`) | Matches migration-path doc; key-value fits ID lookup; single file; pure Go | No ad-hoc queries; bucket design required |
| **SQLite** (`modernc.org/sqlite`) | Familiar SQL; easy schema introspection; pure Go driver available | Heavier for three CRUD ops; migration boilerplate |

**Recommendation:** **bbolt** for V0.5 MVP.

- Access pattern is exactly `map[id]ModelDefinition` — no joins or filters
- Aligns with "BoltDB or SQLite — TBD" ordering in specs and migration path
- Minimizes schema migration surface before v0.6 inference work
- Single `.db` file parallels future KV block storage mental model (separate concern, same locality)

**Fallback:** If team prefers SQL tooling for ops/debugging, `modernc.org/sqlite` is acceptable — record decision in `accepted.md` before coding.

### Serialization (approved)

`BboltStore` writes `ModelDefinition` values directly into a bbolt bucket keyed by `def.ID` (`[]byte`). Values are `encoding/gob`-encoded inside `BboltStore` — no separate JSON-per-ID persistence layer or interim file store.

### `models.yaml` seed

- **Purpose:** Bootstrap empty store on first run; not a live sync source in V0.5
- **Behavior:** On `Open`, if store empty and seed path provided, load seed; if store non-empty, skip seed (or error if seed conflicts — recommend skip-with-log for MVP)
- **Location:** Configurable path; example in `testdata/models.yaml`
- **Out of scope:** Hot reload, multi-file includes, env substitution

### Concurrency

Match existing `sync.RWMutex` semantics:

- bbolt: single writer — `RegisterModel` uses `Update`; reads use `View`
- SQLite: enable WAL; one connection or connection pool with mutex — keep simple for V0.5

### Error and corruption handling

- Corrupt DB file → return wrapped error on `Open`; do not panic
- Partial write → rely on bbolt transactional `Update`
- Invalid YAML → fail fast at seed load with path in error message

## Boundaries — Must NOT Change

| Package | Reason |
|---------|--------|
| `internal/scheduler/` | No registry persistence coupling |
| `internal/runtime/` | Wiring deferred to v0.6 integration |
| `internal/loader/` | Same |
| `internal/events/` | No new event types required for V0.5 spec |
| Agent registry | Separate subsystem; future task if needed |

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Interface extraction breaks existing tests | Low | Extract in-memory first; tests green before persistent |
| Backend indecision blocks start | Medium | Step 0 decision gate; default bbolt if no objection |
| New dependency without approval | Medium | Log in accepted.md; justify vs stdlib (insufficient for durable KV store) |
| Scope creep into agent persistence | High | Explicit out-of-scope in task; code review gate |
| Accidental 011 KV work | High | Execution docs realigned; 011 remains locked Month 3 |

## Alternatives Considered

| Alternative | Verdict |
|-------------|---------|
| Extend `models.go` in place without interface | Rejected — blocks test isolation and Path 3 |
| JSON file only (no DB) | Rejected — no atomic writes; race on concurrent register |
| Persist agents and models together | Rejected — out of V0.5 scope |
| Start 011 KV manager instead | Rejected — wrong version; blocked on v0.7; see planning assessment |

## Approval Checklist

- [x] Backend choice approved — **bbolt** (`go.etcd.io/bbolt`)
- [x] `ModelStore` → `InMemoryStore` / `BboltStore` approved
- [x] `models.yaml` seeds only when store empty; no hot reload
- [x] Dependency addition approved
- [x] No runtime, scheduler, loader, or KV manager changes
- [x] Task [006-registry-persistence.md](../../tasks/006-registry-persistence.md) acceptance criteria accepted

## Architecture Compliance (pre-code)

| Check | Result |
|-------|--------|
| INV-004 — registry owns model definitions | Pass — all persistence in `internal/registry/` |
| INV-001 — scheduler no runtime imports | Pass — scheduler untouched |
| INV-002 — runtime no scheduler imports | Pass — runtime untouched |
| INV-008 — no llama.cpp outside inference | Pass — no inference packages touched |
| Contract API unchanged (`RegisterModel`, `GetModel`, `ListModels`) | Pass — `ModelRegistry` delegates to `ModelStore` |
| Forbidden packages (agents, K8s, KV) | Pass — out of scope |
| Dependency justification | Pass — stdlib has no durable embedded KV; bbolt is pure Go |
| Migration Path 3 | Pass — interface + in-memory + bbolt adapter |

## Post-V0.5 (not this task)

- v0.6: runtime reads persistent registry for GGUF paths
- v0.7: MLFQ scheduler (extends 010)
- Month 3: 011 KV manager, 012 swap manager — after v0.7
- Optional: agent persistence RFC if fleet config must survive restart

---
**Layer:** planning
**Related:** [../../tasks/006-registry-persistence.md](../../tasks/006-registry-persistence.md), [../../docs/specs/v0.5-model-registry.md](../../docs/specs/v0.5-model-registry.md)
