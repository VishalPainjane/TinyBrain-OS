# Build Tag Matrix — Inference Backends (009a)

**Status:** Accepted  
**Date:** 2026-06-08  
**Gate:** D — Build tag review  
**Rule:** Exactly **one** GPU/CPU backend implementation file set is active per binary build. Tags are mutually exclusive.

---

## Design Principles

1. **Positive backend tag** selects one implementation (`cuda`, `rocm`, `metal`, or `vulkan`).
2. **CPU is default** when no GPU backend tag is set: `cgo && !cuda && !rocm && !metal && !vulkan`.
3. **Negate all other backends** on every GPU file to prevent duplicate symbols.
4. **`bindings_cgo.go`** uses `cgo` only — no `#cgo` LDFLAGS for GPU.
5. **`integration`** is orthogonal — adds tests, not backend selection.

---

## File × Build Tag Matrix

| File | Build constraint | Backend |
|------|------------------|---------|
| `provider_stub.go` | `!cgo` | None (Go stub) |
| `bindings_cgo.go` | `cgo` | Shared CGO preamble (no GPU flags) |
| `provider.go`, `config.go`, `errors.go`, `context.go` | *(none)* | All builds |
| `generate_stub.go`, `context_stub.go` | *(none)* | All builds |
| `load_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` | **CPU** |
| `load_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | **CUDA** |
| `load_rocm.go` | `cgo && rocm && !cuda && !metal && !vulkan` | **ROCm** (future) |
| `load_metal.go` | `cgo && metal && !cuda && !rocm && !vulkan` | **Metal** (future) |
| `load_vulkan.go` | `cgo && vulkan && !cuda && !rocm && !metal` | **Vulkan** (future) |
| `load_integration_test.go` | `cgo && integration` | Test overlay (uses active backend binary) |

**Metal note:** Apple toolchains may add `darwin` to `load_metal.go` when implemented: `cgo && darwin && metal && !cuda && !rocm && !vulkan`.

---

## Mutual Exclusivity Verification

| Backend pair | Both compile together? | Mechanism |
|--------------|------------------------|-----------|
| CPU + CUDA | **No** | CPU requires `!cuda`; CUDA requires `cuda` |
| CPU + ROCm | **No** | `!rocm` vs `rocm` |
| CPU + Metal | **No** | `!metal` vs `metal` |
| CPU + Vulkan | **No** | `!vulkan` vs `vulkan` |
| CUDA + ROCm | **No** | Each negates the other |
| CUDA + Metal | **No** | Each negates the other |
| CUDA + Vulkan | **No** | Each negates the other |
| ROCm + Metal | **No** | Each negates the other |
| ROCm + Vulkan | **No** | Each negates the other |
| Metal + Vulkan | **No** | Each negates the other |

**Verdict:** All five backends are mutually exclusive per binary.

---

## Build Commands (documentation)

| Target | Command |
|--------|---------|
| CI default (stub) | `CGO_ENABLED=0 go test ./...` |
| Linux CPU | `CGO_ENABLED=1 go build ./internal/inference/llama/...` |
| Linux CUDA | `CGO_ENABLED=1 go build -tags cuda ./internal/inference/llama/...` |
| Linux ROCm | `CGO_ENABLED=1 go build -tags rocm ./internal/inference/llama/...` |
| macOS CPU | `CGO_ENABLED=1 go build ./internal/inference/llama/...` |
| macOS Metal | `CGO_ENABLED=1 go build -tags metal ./internal/inference/llama/...` |
| Vulkan | `CGO_ENABLED=1 go build -tags vulkan ./internal/inference/llama/...` |
| Integration test | `CGO_ENABLED=1 go test -tags 'integration' ./internal/inference/llama/...` |

---

## 009a Scope

| Backend | File in 009a | Build tag |
|---------|--------------|-----------|
| CPU | `load_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` |
| CUDA | `load_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` |
| ROCm | — | Deferred |
| Metal | — | Deferred |
| Vulkan | — | Deferred |

**Correction applied:** Original 009a review used `cgo && !cuda` only — **superseded** by this matrix.

---
**Layer:** planning  
**Related:** [009a-architecture-review.md](009a-architecture-review.md), [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md)
