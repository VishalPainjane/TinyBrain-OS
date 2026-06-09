# Architecture Review — Task 009c (Runtime / Loader / Inference Integration)

**Status:** Approved for implementation (2026-06-09)  
**Date:** 2026-06-09  
**Scope:** Wire `ModelRuntime` → `loader.Loader` → `InferenceProvider` (`LlamaProvider`). Consolidate `ModelResolver` ownership in runtime (Option B). Close 009a lifecycle debt. No scheduler, agent, KV, GPU, contract, or architecture doc changes.

| Gate | Status |
|------|--------|
| Architecture review | **Approved** |
| Resolver design (Option B) | **Approved** |
| Code implementation | Not started |

**Pre-read:** [009a-architecture-review.md](009a-architecture-review.md), [009b-architecture-review.md](009b-architecture-review.md), [009a-registry-resolver.md](009a-registry-resolver.md), [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md), [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md), [v0.6-inference.md](../../docs/specs/v0.6-inference.md)

**Gates inherited (unchanged from 009a/009b):** A Lifecycle, C Dependency, D Build tags, E API pin b9553. Gate B Resolver — **superseded for ownership location** by Option B below; behavioral rules unchanged (inference must not import registry).

---

## Resolver Design (Approved — Option B)

**Decision:** Single shared resolver contract owned by `internal/runtime/`. Inference consumes `runtime.ModelResolver`; duplicate inference-owned resolver types and registry adapter are removed.

| Artifact | Owner (009c) |
|----------|----------------|
| `ModelSpec` | `internal/runtime/resolver.go` |
| `ModelResolver` | `internal/runtime/resolver.go` |
| `RegistryResolver` | `internal/runtime/registry_resolver.go` (sole registry import for model resolution in kernel path) |
| `LlamaProvider.resolver` field | `runtime.ModelResolver` |
| `ModelDefinition` metadata | `internal/registry/` (INV-004, unchanged) |

**Removed from inference tree:**

- `internal/inference/llama/resolver.go` — delete
- `internal/inference/llama/registry_resolver.go` — delete

**Wiring at composition root:**

```text
ModelRegistry (or static / file / remote map)
        │
        ▼
runtime.NewRegistryResolver(reg)  ──► runtime.ModelResolver (single instance)
        │
        ├─► NewIntegratedModelRuntime(..., resolver, ...)  → loader.Load(modelID, spec.Path)
        └─► llama.NewLlamaProvider(resolver, cfg)          → provider.LoadModel(modelID)
```

**009a-registry-resolver.md note:** Port behavior and ADR-004 goals are preserved. Ownership moves from inference package to runtime package; inference `provider.go` / `load_*.go` remain registry-free.

---

## 1. Exact Files to Create

| File | Purpose |
|------|---------|
| `tasks/009c-runtime-integration.md` | Task spec with acceptance criteria |
| `internal/runtime/resolver.go` | `ModelSpec`, `ModelResolver` (moved from inference) |
| `internal/runtime/registry_resolver.go` | `RegistryResolver` adapter (moved from inference) |
| `internal/runtime/registry_resolver_test.go` | Resolver adapter tests (migrated from `provider_test.go`) |
| `internal/runtime/runtime_loader_test.go` | Loader-orchestration unit tests (`StubProvider` + `StubLoader`, `CGO_ENABLED=0`) |
| `internal/runtime/runtime_integration_test.go` | `//go:build cgo && integration` — `ModelRuntime` + `LlamaProvider` + shared resolver + real GGUF |

**Not created:** `cmd/`, `internal/app/`, `ModelPathResolver`, parallel inference resolver, new top-level packages.

---

## 2. Exact Files to Modify

| File | Change |
|------|--------|
| `internal/runtime/runtime.go` | Optional `loader.Loader` + `ModelResolver`; orchestrate load/unload; preserve loader-less path |
| `internal/runtime/runtime_test.go` | Keep existing stub-only tests green |
| `internal/inference/llama/provider.go` | `resolver` field type → `runtime.ModelResolver`; `Resolve` uses `runtime.ModelSpec` |
| `internal/inference/llama/provider_test.go` | `staticResolver` returns `runtime.ModelSpec`; remove `TestRegistryResolver_*` (moved to runtime) |
| `internal/inference/llama/provider_stub_test.go` | Update `staticResolver` / `ModelSpec` references |
| `internal/inference/llama/load_integration_test.go` | `staticResolver` uses `runtime.ModelSpec` |
| `internal/inference/llama/generate_integration_test.go` | `staticResolver` uses `runtime.ModelSpec` |
| `tests/import_boundary_test.go` | Add `TestRuntimeDoesNotImportInference` |
| `README.md` | Integrated wiring demo with single shared resolver |
| `.github/workflows/ci.yml` | Optional: `runtime_integration_test.go` when `TB_TEST_GGUF_PATH` set |

**Delete:**

| File | Reason |
|------|--------|
| `internal/inference/llama/resolver.go` | Superseded by `internal/runtime/resolver.go` |
| `internal/inference/llama/registry_resolver.go` | Superseded by `internal/runtime/registry_resolver.go` |

**Post-merge planning sync (merge commit, not blocking review):**

- `planning/execution/completed.md`
- `planning/execution/current-sprint.md`
- `docs/current.md`
- `docs/specs/v0.6-inference.md`
- `planning/releases/v0.6.md`
- `planning/architecture-evolution/current-state.md`
- `planning/metrics/progress.md`

**Explicitly zero-diff:**

- `docs/contracts/*`
- `docs/architecture/*` (including `inference-lifecycle.md`)
- `internal/scheduler/`
- `internal/registry/` (consumed only via `runtime/registry_resolver.go`)
- `internal/loader/` (preferred — no loader package changes)
- `load_cuda.go` and GPU backend files

---

## 3. Runtime Ownership Boundaries

| Responsibility | Owner (`internal/runtime/`) | 009c behavior |
|----------------|---------------------------|---------------|
| Public API | `ModelRuntime` | Unchanged: `LoadModel`, `UnloadModel`, `Generate`, `SaveContext`, `RestoreContext` |
| `ModelResolver` / `ModelSpec` | **Runtime** | Shared port for loader path and provider injection |
| `RegistryResolver` | **Runtime** | Sole registry adapter for model resolution |
| Orchestration | `modelRuntime` | Resolver → loader → provider on load; provider → loader on unload |
| Lifecycle state storage | **No** | Loader owns policy `ModelState` |
| Event emission | **Yes** | `ModelLoaded` / `ModelUnloaded` after full successful orchestration |
| Inference execution | Delegate | `Generate` → provider only |
| Context save/restore | Delegate | Provider stubs unchanged (011 deferred) |
| llama.cpp / CGO | **Forbidden** | INV-008 |
| Scheduler | **Forbidden** | INV-002 |

**Constructors (backward compatible):**

```text
NewModelRuntime(provider, bus)
    → loader-less; provider-only load/unload (existing tests)

NewIntegratedModelRuntime(provider, loader, resolver, bus)
    → 009c orchestration; single shared resolver instance
```

---

## 4. Inference Ownership Boundaries

| Responsibility | Owner (`internal/inference/llama/`) | 009c change |
|----------------|-----------------------------------|-------------|
| `InferenceProvider` port | `LlamaProvider` | Unchanged surface |
| Engine state | CGO (`nativeModels`, `nativeContexts`) | Unchanged |
| mmap + context + Generate | `load_*.go`, `generate_*.go`, `bindings_cgo.go` | Unchanged |
| `ModelResolver` types | **Consumer only** | Field type `runtime.ModelResolver`; no local `ModelSpec` |
| Registry import | **Forbidden** | `registry_resolver.go` deleted; no registry import in inference tree |
| Bus events | **None** | Unchanged |

**Rule:** `LlamaProvider` receives the same `runtime.ModelResolver` instance as `ModelRuntime`. Provider still calls `resolver.Resolve` internally on `LoadModel` — runtime uses `spec.Path` for `loader.Load` first.

---

## 5. Loader Ownership Boundaries

| Responsibility | Owner (`internal/loader/`) | 009c change |
|----------------|---------------------------|-------------|
| Policy state machine | `StubLoader` | **No code change** (preferred) |
| API used by runtime | `Load`, `Unload`, `State` | Runtime supplies path from `resolver.Resolve` |
| mmap / engine | **No** | Inference adapter |
| Warm / Prefetch / Evict | `StubLoader` | **Out of scope** |

**Target ownership (closes 009a debt):**

| Transition | Owner after 009c |
|------------|------------------|
| → `LOADING` → `ACTIVE` (policy) | Loader (`Load`) |
| → `ACTIVE` (engine) | Inference (`LoadModel` after loader succeeds) |
| → `UNLOADING` → `UNLOADED` (policy) | Loader (`Unload` after provider succeeds) |
| Engine teardown | Inference (`UnloadModel` first) |

---

## 6. Event Flow

```text
LoadModel(modelID)
  │
  ├─ [integrated] resolver.Resolve(modelID) → spec.Path
  ├─ [integrated] loader.Load(modelID, spec.Path)
  ├─ provider.LoadModel(modelID)
  │     └─ on failure: loader.Unload(modelID) compensating rollback
  └─ bus.Publish(ModelLoaded { ModelID })

UnloadModel(modelID)
  │
  ├─ provider.UnloadModel(modelID)
  ├─ [integrated] loader.Unload(modelID)
  └─ bus.Publish(ModelUnloaded { ModelID })

Generate(req)
  └─ provider.Generate(req)    // no lifecycle events
```

**Failure = no event** (extends existing `TestModelRuntime_FailurePath_NoLifecycleEvents`).

**Not emitted in 009c:** `SwapStarted`, `SwapCompleted`, `KVStored`, `KVLoaded`.

---

## 7. Lifecycle Transitions

### Integrated load

```text
[Resolver]  modelID → ModelSpec.Path
[Loader]    NOT_LOADED ──Load(path)──► LOADING ──► ACTIVE
[Inference] ──LoadModel(modelID)──► engine ACTIVE
[Runtime]   emits ModelLoaded when both succeed
```

### Integrated unload

```text
[Inference] ──UnloadModel──► engine freed
[Loader]    ACTIVE ──Unload──► UNLOADING ──► UNLOADED
[Runtime]   emits ModelUnloaded when both succeed
```

### Generate

```text
[Loader]    stays ACTIVE
[Inference] memory_clear → prefill → decode
[Runtime]   delegate only
```

### Loader-less mode (backward compat)

```text
Provider-only load/unload/events — identical to pre-009c runtime.go
```

### Rollback rules

1. Provider load failure → `loader.Unload(modelID)`; no `ModelLoaded` event.
2. Unload order: **provider first**, loader second.
3. Integrated tests use `NewStubLoader()` with **capacity 0** (no silent eviction).

---

## 8. ModelRuntime Integration Plan

### Phase 1 — Resolver migration (Option B)

1. Add `internal/runtime/resolver.go` with `ModelSpec` + `ModelResolver`.
2. Move `RegistryResolver` to `internal/runtime/registry_resolver.go`.
3. Migrate `TestRegistryResolver_*` to `registry_resolver_test.go`.
4. Update `LlamaProvider` to use `runtime.ModelResolver`.
5. Delete `inference/llama/resolver.go` and `inference/llama/registry_resolver.go`.
6. Update all inference tests (`staticResolver`, integration tests) to `runtime.ModelSpec`.

### Phase 2 — Orchestration

7. Add `NewIntegratedModelRuntime(provider, loader, resolver, bus)`.
8. Keep `NewModelRuntime(provider, bus)` as loader-less wrapper.
9. Implement load/unload sequence in `runtime.go` (§6–§7).

### Phase 3 — Composition (tests / README)

```text
reg := registry.NewModelRegistry(...) // or test registry
resolver := runtime.NewRegistryResolver(reg)
loader := loader.NewStubLoader()
provider := llama.NewLlamaProvider(resolver, cfg)
rt := runtime.NewIntegratedModelRuntime(provider, loader, resolver, bus)
```

### Phase 4 — Verification

- Unit: stub stack + resolver migration tests.
- Integration: E2E via `ModelRuntime` with shared resolver + `TB_TEST_GGUF_PATH`.
- Boundary: `runtime` package must not import `internal/inference`.

---

## 9. Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| AC-1 | `ModelRuntime.LoadModel` → `Generate` → `UnloadModel` returns real tokens via `LlamaProvider` | `runtime_integration_test.go` |
| AC-2 | Integrated load sets loader `State` = `ACTIVE` | `runtime_loader_test.go` |
| AC-3 | Integrated unload sets loader `State` = `UNLOADED` | `runtime_loader_test.go` |
| AC-4 | `ModelLoaded` / `ModelUnloaded` on integrated success | Event tests |
| AC-5 | Provider load failure → no event; loader rolled back | `runtime_loader_test.go` |
| AC-6 | Loader-less `NewModelRuntime(stub, bus)` unchanged | Existing `runtime_test.go` |
| AC-7 | `runtime.ModelResolver` shared by runtime and `LlamaProvider` | Code review; single resolver in composition |
| AC-8 | No `internal/registry` import in `internal/inference/llama/` | Import audit |
| AC-9 | `internal/runtime` does not import `internal/inference` | `TestRuntimeDoesNotImportInference` |
| AC-10 | Scheduler zero inference imports | Existing boundary test |
| AC-11 | `CGO_ENABLED=0 go test ./...` passes | CI `test` job |
| AC-12 | `CGO_ENABLED=1 go test ./internal/inference/llama/...` passes | CI `inference-cgo` job |
| AC-13 | Integration suite passes with `TB_TEST_GGUF_PATH` | Manual / optional CI |
| AC-14 | v0.6: Runtime ↔ `LlamaProvider` wiring complete | Spec checkbox on merge |
| AC-15 | `SaveContext` / `RestoreContext` remain stubs | Unchanged tests |
| AC-16 | No scheduler / agent / KV / GPU changes | Scope review |

---

## 10. Architecture Risks

| # | Risk | Severity | Mitigation |
|---|------|----------|------------|
| R-1 | Dual state machines desync (loader ACTIVE, provider not loaded) | High | Compensating `loader.Unload` on provider failure; integration tests assert `loader.State` |
| R-3 | Unload partial failure (provider ok, loader fails) | Medium | Return error; document recovery; loader failures rare (in-memory) |
| R-4 | `StubLoader` capacity eviction without provider unload | High | Default capacity 0 in integrated tests; eviction orchestration deferred |
| R-5 | Concurrent load/unload/generate races | Medium | Accept provider mutex behavior; document |
| R-6 | Scope creep into Warm/Prefetch/Evict | Medium | No new `ModelRuntime` methods |
| R-7 | Runtime imports inference | High | `TestRuntimeDoesNotImportInference` |
| R-8 | Duplicate error types (`loader.Err*` vs `runtime.Err*`) | Low | Map at boundary; `errors.Is` in tests |
| R-15 | Resolver migration breaks 009a/009b tests | Medium | See §11 — full test matrix before merge |

**Removed:** ~~R-2 dual resolver drift~~ — eliminated by Option B single contract.

---

## 11. Migration Risks

| Area | Impact | Mitigation |
|------|--------|------------|
| `inference/llama/resolver.go` deleted | All `ModelSpec` / `ModelResolver` references break | Move types to runtime first; update imports; then delete |
| `inference/registry_resolver.go` deleted | `TestRegistryResolver_*` moves to runtime | `registry_resolver_test.go` |
| `staticResolver` in inference tests | Signature change to `runtime.ModelSpec` | Update `provider_test.go`, integration tests |
| `NewRegistryResolver` call sites | Package path `runtime.NewRegistryResolver` | Grep and update tests / README |
| `NewModelRuntime(provider, bus)` callers | None | Signature unchanged |
| Direct `LlamaProvider` tests | Import/runtime type updates only | `go test ./internal/inference/llama/...` |
| v0.6 release tagging | Blocked until 009c merge | Release notes on merge |

### R-15 mitigation (mandatory pre-merge)

```bash
CGO_ENABLED=0 go test ./...
CGO_ENABLED=1 go test ./internal/inference/llama/...
# integration (local / optional CI):
TB_TEST_GGUF_PATH=/path/to/model.gguf go test -tags integration ./internal/inference/llama/... ./internal/runtime/...
```

---

## 12. Cross-Platform Impact

| OS | Impact |
|----|--------|
| Linux | Primary CI — resolver migration unit tests + CGO inference + optional runtime integration |
| Windows | Manual — same wiring; CGO + `TB_TEST_GGUF_PATH` |
| macOS | Manual CPU |

Runtime orchestration and resolver types are **pure Go** on all platforms. CGO variance remains in `internal/inference/llama/`.

---

## 13. CI Impact

| Job | 009c change | Blocker |
|-----|-------------|---------|
| `test` (`CGO_ENABLED=0`) | Resolver migration in runtime; `runtime_loader_test.go`; boundary test | Yes |
| `inference-cgo` | Inference compiles with `runtime.ModelResolver`; resolver tests in runtime package | Yes |
| Integration (`-tags integration`) | Optional runtime + inference E2E when `TB_TEST_GGUF_PATH` set | No (recommended) |

No llama.cpp cmake changes. No new CI jobs required.

---

## 14. Test Strategy

### Unit — `CGO_ENABLED=0` (always run)

| Test | Package | Assert |
|------|---------|--------|
| `TestRegistryResolver_Resolve` | `internal/runtime` | Migrated from inference |
| `TestRegistryResolver_emptyPath` | `internal/runtime` | Migrated from inference |
| `TestIntegratedRuntime_Load_setsLoaderActive` | `internal/runtime` | Loader ACTIVE after load |
| `TestIntegratedRuntime_Unload_setsLoaderUnloaded` | `internal/runtime` | Loader UNLOADED after unload |
| `TestIntegratedRuntime_LoadProviderFail_rollbackLoader` | `internal/runtime` | No event; loader rolled back |
| `TestIntegratedRuntime_Events` | `internal/runtime` | One loaded + one unloaded |
| `TestRuntimeDoesNotImportInference` | `tests` | No inference in runtime deps |
| All existing `LlamaProvider_*` validation tests | `internal/inference/llama` | Updated `staticResolver` types |
| All existing `runtime_test.go` stub tests | `internal/runtime` | Unchanged behavior |

### CGO — `inference-cgo` job

| Test | Assert |
|------|--------|
| Inference package unit + compile | Passes with `runtime.ModelResolver` |
| `registry_resolver` tests | Run under `internal/runtime` in `test` job |

### Integration — `cgo && integration`, `TB_TEST_GGUF_PATH`

| Test | Assert |
|------|--------|
| `TestIntegratedRuntime_Llama_LoadGenerateUnload` | Real tokens via `ModelRuntime` |
| `TestIntegratedRuntime_Llama_loaderState` | Loader ACTIVE / UNLOADED |
| Existing `load_integration_test.go` / `generate_integration_test.go` | Regression after resolver migration |

**Shared resolver in integration:** one `staticResolver` or `RegistryResolver` instance passed to both `NewIntegratedModelRuntime` and `NewLlamaProvider`.

### Explicitly not tested in 009c

- Scheduler → runtime commands
- Warm / Prefetch / Evict
- SaveContext / RestoreContext (stubs)
- CUDA / GPU
- Loader capacity eviction with provider sync
- Multi-model concurrent ACTIVE

---

## Implementation Sequence (Ordered)

1. `internal/runtime/resolver.go` + `registry_resolver.go` + `registry_resolver_test.go`
2. Update `LlamaProvider` + inference tests; delete inference resolver files
3. **Gate:** R-15 mitigation commands green
4. `NewIntegratedModelRuntime` + orchestration in `runtime.go`
5. `runtime_loader_test.go`
6. `runtime_integration_test.go`
7. `TestRuntimeDoesNotImportInference`
8. README + optional CI step
9. Planning doc sync on merge

---

## Scope Enforcement Checklist

- [ ] Single `runtime.ModelResolver` — no `ModelPathResolver`, no inference-owned resolver types
- [ ] `RegistryResolver` only in `internal/runtime/`
- [ ] No `internal/registry` import in `internal/inference/llama/`
- [ ] No `docs/contracts/` or `docs/architecture/` changes
- [ ] No scheduler, agent, KV, or GPU work
- [ ] No `ModelRuntime.Warm` / `Prefetch` / `Evict`
- [ ] R-15 mitigation run before merge

---
**Layer:** planning  
**Related:** [009b-architecture-review.md](009b-architecture-review.md), [009a-registry-resolver.md](009a-registry-resolver.md), [../../docs/specs/v0.6-inference.md](../../docs/specs/v0.6-inference.md)
