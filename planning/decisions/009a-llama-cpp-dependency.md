# Decision Gate — Task 009a llama.cpp Dependency

**Status:** Approved (planning gate closed 2026-06-08)  
**Date:** 2026-06-08  
**Task:** 009a — llama.cpp package skeleton + CGO load  
**Cross-platform policy:** [docs/architecture/cross-platform.md](../../docs/architecture/cross-platform.md)  
**Capability matrix:** [docs/architecture/inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md)  
**Lifecycle:** [docs/architecture/inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md)  
**Build tags:** [009a-build-tags.md](009a-build-tags.md)  
**Resolver:** [009a-registry-resolver.md](009a-registry-resolver.md)

---

## Dependency

| Item | Type | Location |
|------|------|----------|
| llama.cpp | C/C++ library (git submodule, pinned tag) | `third_party/llama.cpp/` |
| C/C++ toolchain | Build-time | OS package manager / MSVC |

No new Go module dependencies for 009a (direct CGO + submodule).

---

## Submodule Strategy

| Decision | Detail |
|----------|--------|
| Mechanism | Git submodule: `https://github.com/ggml-org/llama.cpp.git` → `third_party/llama.cpp` |
| Clone | `git submodule update --init --recursive third_party/llama.cpp` |
| Shallow | Optional `--depth 1` for CI clone speed |
| Build artifacts | `third_party/llama.cpp/build/` gitignored |
| Consumer | CGO `#cgo CFLAGS/LDFLAGS` point to submodule include + static/shared lib output |
| Update policy | Explicit PR only; never floating `master` in CI |
| Fork | Not used — upstream ggml-org/llama.cpp |

---

## Version Pinning Strategy

| Decision | Detail |
|----------|--------|
| Pin method | Submodule commit SHA tied to upstream **release tag** `bNNNN` |
| Record | Tag name + SHA in this file and first 009a implementation commit message |
| Initial pin | Select latest stable `b*` tag at implementation start; verify CPU build on Linux + dev Windows |
| Compatibility | Document chosen tag in README; breaking API changes require new pin PR |
| Security | Dependabot not applicable to submodule; manual CVE review on pin bumps |
| Go module | No `go get` for llama.cpp — avoids duplicate version sources |

**Implementation note:** Exact tag recorded when submodule is added (first 009a code commit).

---

## CI Strategy

| Job | Runner | `CGO_ENABLED` | Tags | Purpose | 009a |
|-----|--------|---------------|------|---------|------|
| `test` (existing) | `ubuntu-latest` | `0` | — | Full `go test ./...` stub path | **Required** — merge blocker |
| `inference-cgo` (new) | `ubuntu-latest` | `1` | — | `go test ./internal/inference/llama/...` CPU backend | **Required** at 009a merge |
| `inference-integration` | `ubuntu-latest` | `1` | `integration` | Checksum-verified SmolLM2-135M Q4_K_M; 5 llama tests | **Required** — merge blocker (STAB-001) |
| `inference-integration-runtime` | `ubuntu-latest` | `1` | `integration` | Same GGUF bootstrap; 1 runtime E2E test | **Required** — merge blocker (STAB-001) |
| Windows compile | `windows-latest` | `1` | — | `go build ./internal/inference/llama/...` | Optional 009a |
| macOS compile | `macos-latest` | `1` | — | CPU build only | Optional 009a |
| CUDA GPU | GPU runner | `1` | `cuda` | Manual / nightly | Not 009a |

**Linux CI packages (inference-cgo job):** `build-essential`, `cmake` (if building llama from source in CI).

**Default rule:** Merges must pass `CGO_ENABLED=0 go test ./...`. CGO job added in same PR as inference package or immediately after — documented here as 009a merge requirement.

---

## Windows Strategy

| Topic | Decision |
|-------|----------|
| Toolchain | MSVC (Visual Studio Build Tools), `CGO_ENABLED=1` |
| Paths | `filepath` only; registry paths may be `C:\` or UNC |
| Build | `load_cpu.go` backend; document cmake/msbuild steps in README |
| CUDA | User-installed toolkit + DLLs; `-tags cuda` separate binary; not bundled |
| mmap | Manual integration test on dev machine; record in assumptions at 009a ship |
| CI | Not blocking 009a; optional `windows-latest` compile job |

---

## Linux Strategy

| Topic | Decision |
|-------|----------|
| CI primary | `ubuntu-latest`, glibc |
| CPU build | Default 009a verification path |
| CUDA | `-tags cuda` local/dev; NVIDIA container for future K8s |
| ROCm | `-tags rocm` future; not 009a |
| musl/Alpine | **Unsupported** in 009a — document risk; no CI |
| Paths | No hardcoded `/usr/local` |

---

## macOS Strategy

| Topic | Decision |
|-------|----------|
| Toolchain | Xcode Command Line Tools |
| CPU build | `cgo` without GPU tags — 009a verification target (manual) |
| Metal | `-tags metal` future file `load_metal.go`; not 009a |
| CUDA | **N/A** on macOS |
| CI | Not blocking 009a; optional `macos-latest` compile |
| Distribution | Code signing deferred (future release) |

---

## CUDA Strategy

| Topic | Decision |
|-------|----------|
| 009a scope | `load_cuda.go` + `-tags cuda`; config maps `hardware.BackendCUDA` → `NGLayers` |
| Default | `NGLayers=0` when tag absent or init fails |
| Isolation | All `#cgo` CUDA LDFLAGS in `load_cuda.go` only ([009a-build-tags.md](009a-build-tags.md)) |
| CI | Not required 009a |
| Fallback | CPU path when CUDA build not used or GPU init fails |
| Multi-GPU | **No** — device index `0` only |

---

## 009a Backend Scope

| Backend | 009a | Build |
|---------|------|-------|
| CPU | **Yes** | `cgo && !cuda && !rocm && !metal && !vulkan` |
| CUDA | Config + file | `cgo && cuda && !rocm && !metal && !vulkan` |
| ROCm | Deferred | Future `load_rocm.go` |
| Metal | Deferred | Future `load_metal.go` |
| Vulkan | Deferred | Future `load_vulkan.go` |

---

## Portability Risks

| Risk | Migration path |
|------|----------------|
| CUDA LDFLAGS in shared bindings | Split per [009a-build-tags.md](009a-build-tags.md) |
| OS build steps differ | README Linux / Windows / macOS sections |
| mmap on Windows | Manual test + assumptions.md |
| Registry coupling | [009a-registry-resolver.md](009a-registry-resolver.md) |

---

## Approval Checklist

- [x] [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md) accurate for 009a CPU/CUDA rows
- [x] Submodule strategy approved (git submodule → `third_party/llama.cpp`)
- [x] Version pinning strategy approved (tagged `bNNNN` + SHA)
- [x] Linux CI strategy approved (`CGO_ENABLED=0` default + `inference-cgo` job at merge)
- [x] Windows manual verification plan accepted
- [x] macOS manual verification plan accepted
- [x] CUDA config-only in 009a accepted
- [x] Kubernetes impact acknowledged
- [x] CPU-first, GPU-additive approach accepted
- [x] [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md) created
- [x] [009a-registry-resolver.md](009a-registry-resolver.md) accepted
- [x] [009a-build-tags.md](009a-build-tags.md) accepted
- [ ] Row in [accepted.md](accepted.md) for llama.cpp submodule — **on 009a merge** (not pre-code)
- [ ] Exact submodule tag SHA recorded — **on first implementation commit**

---

## Approval Status

**Gate C: APPROVED** — planning and policy complete. Submodule SHA is recorded at implementation time.

---
**Layer:** planning
