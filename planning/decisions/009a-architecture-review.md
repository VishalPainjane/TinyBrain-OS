# Architecture Review — Task 009a (Phase 1)

**Status:** Approved for implementation (gates closed 2026-06-08)  
**Date:** 2026-06-08  
**Scope:** `internal/inference/llama/` skeleton; resolver-based load/unload; mmap GGUF via CGO. No `Generate`. Core packages unchanged.

**Pre-read:** [postmortems/v0.5.md](../postmortems/v0.5.md), [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md), [009a-registry-resolver.md](009a-registry-resolver.md), [009a-build-tags.md](009a-build-tags.md), [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md)

## Gates Cleared

| Gate | Document | Status |
|------|----------|--------|
| A — Lifecycle | [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md) | Closed |
| B — Resolver | [009a-registry-resolver.md](009a-registry-resolver.md) | Closed |
| C — Dependency | [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md) | Closed |
| D — Build tags | [009a-build-tags.md](009a-build-tags.md) | Closed |

## Cross-Platform Constraints (binding)

Per [cross-platform.md](../../docs/architecture/cross-platform.md) and [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md):

1. `filepath` only in shared orchestration.
2. `provider.go`, `config.go`, `errors.go`, resolver types — OS-agnostic.
3. GPU/OS specifics in tagged `load_*.go` only.
4. Core packages — **zero diffs**.
5. CPU baseline; CUDA file optional at compile time.
6. **No `internal/registry` import** in `provider.go` or `load_*.go`.

## Files to Create

| File | Build constraint |
|------|------------------|
| `internal/inference/llama/doc.go` | — |
| `internal/inference/llama/config.go` | — |
| `internal/inference/llama/provider.go` | — |
| `internal/inference/llama/errors.go` | — |
| `internal/inference/llama/context.go` | — |
| `internal/inference/llama/resolver.go` | — (`ModelResolver`, `ModelSpec` types) |
| `internal/inference/llama/registry_resolver.go` | — (sole file importing `registry`) |
| `internal/inference/llama/provider_stub.go` | `!cgo` |
| `internal/inference/llama/bindings_cgo.go` | `cgo` (no GPU LDFLAGS) |
| `internal/inference/llama/load_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` |
| `internal/inference/llama/load_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` |
| `internal/inference/llama/generate_stub.go` | — |
| `internal/inference/llama/context_stub.go` | — |
| `internal/inference/llama/provider_test.go` | — |
| `internal/inference/llama/load_integration_test.go` | `cgo && integration` |
| `third_party/llama.cpp/` | submodule |
| `.gitmodules` | if submodule |

## Files to Modify

| File | Change |
|------|--------|
| `README.md` | Linux, Windows, macOS build sections |
| `.gitignore` | Submodule build artifacts |
| `.github/workflows/ci.yml` | Add `inference-cgo` job at 009a merge |

## Package Boundaries

- `internal/inference/llama` → `runtime`, `hardware`, inference resolver types
- `registry_resolver.go` → only inference file importing `registry`
- `internal/runtime` → must not import `inference`

## Acceptance Criteria

- [ ] [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md) followed
- [ ] `ModelResolver` injection; no registry import in provider/load files
- [ ] Build tags per [009a-build-tags.md](009a-build-tags.md)
- [ ] INV-008 satisfied
- [ ] Core packages unchanged
- [ ] `CGO_ENABLED=0 go test ./...` green
- [ ] CPU load on Linux with integration tag + `TB_TEST_GGUF_PATH`
- [ ] No `Generate` implementation
- [ ] [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md) updated on merge

---
**Layer:** planning
