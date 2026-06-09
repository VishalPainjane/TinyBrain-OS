# Task 009c — Runtime / Loader / Inference Integration

## Status

Complete (pending merge commit)

## Goal

Wire `ModelRuntime` to `loader.Loader` and `LlamaProvider` via shared `runtime.ModelResolver`. Close 009a lifecycle debt (loader policy + inference execution orchestrated by runtime).

## Context

009a/009b delivered `LlamaProvider` with load/unload/generate but no runtime integration. Inference-lifecycle.md defines loader as policy owner and inference as execution owner; runtime orchestrates and emits events.

## Requirements

- `runtime.ModelResolver` + `runtime.ModelSpec` owned by `internal/runtime/`
- `runtime.RegistryResolver` — sole registry-backed resolver in kernel path
- `LlamaProvider` consumes `runtime.ModelResolver` (no inference-owned resolver types)
- `NewIntegratedModelRuntime(provider, loader, resolver, bus)` — resolve → loader.Load → provider.LoadModel → `ModelLoaded`
- Rollback: provider load failure → `loader.Unload`; no event
- Unload: provider.UnloadModel → loader.Unload → `ModelUnloaded`
- `NewModelRuntime(provider, bus)` loader-less path preserved
- `TestRuntimeDoesNotImportInference` boundary test
- E2E integration: `runtime_integration_test.go` with `RegistryResolver`, real `LlamaProvider`, `TB_TEST_GGUF_PATH`

## Files

- `internal/runtime/resolver.go`
- `internal/runtime/registry_resolver.go`
- `internal/runtime/registry_resolver_test.go`
- `internal/runtime/runtime.go` (orchestration)
- `internal/runtime/runtime_loader_test.go`
- `internal/runtime/runtime_integration_test.go`
- `internal/inference/llama/provider.go` (resolver type migration)
- `tests/import_boundary_test.go`
- `planning/decisions/009c-architecture-review.md`

## Acceptance Criteria

- [x] Shared `runtime.ModelResolver` used by runtime and `LlamaProvider`
- [x] Integrated load sets loader `State` = `ACTIVE`
- [x] Integrated unload sets loader `State` = `UNLOADED`
- [x] `ModelLoaded` / `ModelUnloaded` on integrated success only
- [x] Provider load failure rolls back loader; no lifecycle event
- [x] Loader-less `NewModelRuntime` unchanged
- [x] `internal/runtime` does not import `internal/inference`
- [x] No `internal/registry` import in inference production code
- [x] `CGO_ENABLED=0 go test ./...` passes
- [x] E2E integration test written (`cgo && integration`, `TB_TEST_GGUF_PATH`)
- [x] Scheduler zero inference imports (unchanged)

## Out Of Scope

- Scheduler → runtime commands
- CUDA / GPU backends
- KV manager, SaveContext / RestoreContext implementation
- Warm / Prefetch / Evict on `ModelRuntime`
- Loader package changes

## Related

- Spec: [docs/specs/v0.6-inference.md](../docs/specs/v0.6-inference.md)
- Review: [planning/decisions/009c-architecture-review.md](../planning/decisions/009c-architecture-review.md)
- Architecture: [docs/architecture/inference-lifecycle.md](../docs/architecture/inference-lifecycle.md)

---
**Layer:** task
