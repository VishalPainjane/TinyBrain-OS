# Inference Backend Capability Matrix

**Source of truth** for inference backend support across operating systems and hardware.  
**Rule:** Update this document **before** implementing or enabling any backend capability. Do not assume feature parity between backends — track parity explicitly.

**Engine:** llama.cpp via `internal/inference/llama/` (INV-008)  
**Port:** `runtime.InferenceProvider` ([contracts/runtime.md](../contracts/runtime.md))  
**Related:** [cross-platform.md](cross-platform.md), [inference-lifecycle.md](inference-lifecycle.md), [hardware.md](hardware.md), [009a-build-tags.md](../../planning/decisions/009a-build-tags.md)

## Status Legend

| Symbol | Meaning |
|--------|---------|
| **Yes** | Implemented and verified |
| **Partial** | Implemented with documented limitations |
| **Planned** | Scoped; not yet implemented |
| **Stub** | Fake implementation for tests only |
| **No** | Not supported; no current plan |
| **N/A** | Not applicable to this backend |

## Capability Matrix

Capabilities apply to the **llama.cpp adapter** unless noted. Feature parity across backends is **not** assumed.

| Capability | CPU | CUDA | ROCm/HIP | Metal | Vulkan |
|------------|-----|------|----------|-------|--------|
| LoadModel | **Partial** (009a) | **Partial** (009d) | Planned | Planned | Planned |
| UnloadModel | **Partial** (009a) | **Partial** (009d) | Planned | Planned | Planned |
| Generate | **Partial** (009b) | **Partial** (009d) | Planned | Planned | Planned |
| SaveContext | Stub (009a) | Stub (009a) | Stub | Stub | Stub |
| RestoreContext | Stub (009a) | Stub (009a) | Stub | Stub | Stub |
| Quantization (GGUF Q4_K_M) | Planned | Planned | Planned | Planned | Planned |
| Quantization (other GGUF) | Planned | Planned | Planned | Planned | Planned |
| Multi-GPU | No | No | No | No | No |
| mmap load | **Partial** (009a) | **Partial** (009d) | Planned | Planned | Planned |
| CI coverage | **Partial** (009a) | **Partial** (009d) | No | No | No |
| Warm / Prefetch / Evict | N/A (loader) | N/A | N/A | N/A | N/A |
| Concurrent models ACTIVE | No | No | No | No | No |
| Backend fallback to CPU | Planned | Planned | Planned | Planned | Planned |

### Test stub (not a production backend)

| Capability | `StubProvider` (`internal/runtime/`) |
|------------|--------------------------------------|
| LoadModel | **Yes** — all platforms |
| UnloadModel | **Yes** |
| Generate | **Stub** — canned response |
| SaveContext | **Stub** — empty blob |
| RestoreContext | **Stub** |
| Quantization | **N/A** |
| Multi-GPU | **N/A** |
| CI coverage | **Yes** — `CGO_ENABLED=0`, default CI |

`StubProvider` exists for port testing per ADR-004. It is **not** listed in the GPU backend matrix above and must not be treated as llama.cpp parity.

---

## CPU

### Supported operating systems

| OS | Status |
|----|--------|
| Linux | **Partial** (009a) — `inference-cgo` CI green |
| Windows | **Partial** (009a) — manual dev verification |
| macOS | Planned — manual verification |

### Supported hardware

- x86_64 and arm64 CPUs
- No discrete GPU required
- Baseline for all profiles ([hardware.md](hardware.md))

### Known limitations

- Lower throughput than GPU backends
- Large models constrained by RAM on Tiny profile
- No GPU layer offload (`NGLayers = 0`)

### Build requirements

- `CGO_ENABLED=1`
- C/C++ toolchain (gcc/clang on Linux/macOS; MSVC on Windows)
- llama.cpp submodule built **without** GPU flags for CPU-only binary
- Build tag: [009a-build-tags.md](../../planning/decisions/009a-build-tags.md) — `cgo && !cuda && !rocm && !metal && !vulkan`

### Test coverage status

| Test type | Status |
|-----------|--------|
| Unit (`CGO_ENABLED=0` stub path) | **Yes** |
| CGO compile (Linux CI) | **Yes** — `inference-cgo` job |
| Integration (real GGUF) | **Yes** (STAB-001) — merge-blocking `inference-integration` + `inference-integration-runtime`; tag `integration`, deterministic GGUF + SHA256; 5 llama + 1 runtime test |
| Windows automated | No — manual |
| macOS automated | No — manual |

---

## CUDA (NVIDIA)

### Supported operating systems

| OS | Status |
|----|--------|
| Linux | **Partial** (009d) — `-tags cuda` static CGO build; manual GPU per checklist |
| Windows | **Partial** (ADR-006) — runtime DLL loading (`syscall.LoadDLL`); no MSVC/CGO conflict |
| macOS | **No** — NVIDIA does not support CUDA on macOS |

### Supported hardware

- NVIDIA GPUs with compatible driver + CUDA toolkit
- `hardware.BackendCUDA` from probe
- Single GPU only (device index `0` default)
- Verified: RTX 4050 Laptop GPU (SM 8.9, Ada Lovelace, CUDA 12.4)

### Known limitations

- Not available on macOS
- Windows: `backend_windows_dynamic.go` struct layouts are pinned to llama.cpp b9553 (`9e3b928`); must be re-verified on submodule update
- Multi-GPU not supported (v0.6+)
- CI GPU runners not available in default GitHub Actions
- Windows DLL probe is silent on CPU fallback — callers should check `hardware.Probe()` to detect expected GPU

### Build requirements

**Linux (static CGO):**
- `CGO_ENABLED=1`, GCC, llama.cpp built with `-DGGML_CUDA=ON` (MinGW Makefiles)
- Build tag: `cuda` — routes through `bindings_cuda.go` → `bindings_common.go`

**Windows (runtime DLL):**
- NVIDIA CUDA Toolkit 12.x + MSVC 2022 (for `nvcc` host)
- llama.cpp built with `-DGGML_CUDA=ON` using Visual Studio generator → produces `ggml-cuda.dll`, `ggml-base.dll`, `llama.dll`
- **No** build tag required — single binary probes for DLLs at startup via `provider_windows.go:selectBackend()`
- DLL directory: alongside binary, or override with `LLAMACPP_DLL_DIR` env var
- See [ADR-006](../adr/ADR-006-Windows-GPU-Dynamic-DLL-Backend.md)

### Test coverage status

| Test type | Status |
|-----------|--------|
| Unit (`NGLayers`, `ConfigFromProbe`) | **Yes** — `CGO_ENABLED=0` config tests |
| CGO CPU fallback (CUDA tag absent) | **Yes** — CPU binary forces `n_gpu_layers=0` |
| Windows DLL probe (no DLL) | **Yes** — `provider_windows.go` falls back to `cgoBackend` |
| Integration with GPU (Linux) | **Partial** — `cuda_integration_test.go`; `TB_CUDA_INTEGRATION=1` manual |
| Integration with GPU (Windows DLL) | **Yes** — manual check-off signed off; all integration tests pass on Windows CUDA GPU |
| CI | **Partial** — CPU `inference-cgo` unchanged; no GPU runner |

---

## ROCm / HIP (AMD)

### Supported operating systems

| OS | Status |
|----|--------|
| Linux | Planned |
| Windows | Partial — ROCm Windows support limited vs Linux |
| macOS | **No** |

### Supported hardware

- AMD GPUs with ROCm-compatible stack
- `hardware.BackendROCm` from probe (detection may be partial at first)

### Known limitations

- ROCm availability narrower than CUDA
- llama.cpp HIP build flags differ from CUDA
- Probe detection for AMD less mature than NVIDIA on Windows
- Multi-GPU not supported

### Build requirements

- ROCm SDK on Linux
- llama.cpp built with HIP/ROCm flags
- Build tag: `rocm` — `load_rocm.go` (not created in 009a)

### Test coverage status

| Test type | Status |
|-----------|--------|
| All | **No** — documented deferral past 009a |

---

## Metal (Apple)

### Supported operating systems

| OS | Status |
|----|--------|
| macOS | Planned |
| Linux | **No** |
| Windows | **No** |

### Supported hardware

- Apple Silicon (M-series) — primary target
- Intel Mac with AMD GPU — secondary; verify with llama.cpp Metal support at pin tag
- `hardware.BackendMetal` from probe

### Known limitations

- macOS-only backend
- Xcode / Metal SDK required
- Code signing and notarization for distribution (future release concern)
- Unified memory model differs from discrete VRAM budgeting

### Build requirements

- macOS, Xcode CLT
- llama.cpp built with `LLAMA_METAL=1` (or equivalent at pin tag)
- Build tag: `metal` + `darwin` — `load_metal.go` (not created in 009a)

### Test coverage status

| Test type | Status |
|-----------|--------|
| All | **No** — manual macOS verification when implemented |

---

## Vulkan

### Supported operating systems

| OS | Status |
|----|--------|
| Linux | Planned |
| Windows | Planned |
| macOS | Partial — MoltenVK / driver dependent |

### Supported hardware

- Vulkan 1.x compatible GPUs (cross-vendor)
- Not yet represented in `hardware.Backend` enum — probe extension required ([hardware/profile.go](../../internal/hardware/profile.go))

### Known limitations

- llama.cpp Vulkan backend maturity varies by release
- Driver quality differs by vendor
- Probe must detect Vulkan-capable GPU before selecting backend
- Multi-GPU not supported

### Build requirements

- Vulkan SDK + compatible loader
- llama.cpp built with Vulkan support
- Build tag: `vulkan` — `load_vulkan.go` (not created in 009a)
- Add `hardware.BackendVulkan` when probe exists

### Test coverage status

| Test type | Status |
|-----------|--------|
| All | **No** — future backend |

---

## Parity Gaps (explicit)

Track intentional differences — do not imply equal support.

| Gap | Backends affected | Resolution target |
|-----|-------------------|-------------------|
| CUDA generate runtime proof (Windows) | CUDA/Windows | **Resolved** — Signed off in [009d-manual-gpu-checklist.md](../planning/decisions/009d-manual-gpu-checklist.md) |
| Windows DLL struct layout pinned to b9553 | CUDA/Windows DLL | Re-verify `backend_windows_dynamic.go` on submodule update |
| SaveContext / RestoreContext stub | All llama backends | Task 011 KV manager |
| CUDA absent on macOS | CUDA | Permanent — use Metal on macOS |
| Metal absent on Linux/Windows | Metal | Permanent — use CUDA/ROCm/Vulkan |
| No CI GPU integration tests | CUDA, ROCm, Metal, Vulkan | Manual GPU (009d); GPU CI runner deferred |
| No multi-GPU | All | Future RFC / workstation profile |
| Vulkan not in probe enum | Vulkan | Probe task + matrix update |
| StubProvider has Generate stub | Test only (`internal/runtime/`) | Permanent — runtime port tests; llama Generate is separate (009b) |
| ROCm/Metal/Vulkan dynamic loaders not yet written | ROCm, Metal, Vulkan | ADR-006 establishes pattern; implementations deferred |

---

## Maintenance Protocol

1. **Before** adding or enabling a backend file (`load_*.go`), update this matrix to **Planned** with OS/hardware/build rows filled in.
2. **On merge**, change capability cells to **Yes** or **Partial** and update test coverage.
3. **On deferral**, keep **Planned** or **No** with reason in Known limitations — never leave cells ambiguous.
4. Cross-check [cross-platform.md](cross-platform.md) and [009a-llama-cpp-dependency.md](../../planning/decisions/009a-llama-cpp-dependency.md) when CUDA/ROCm/Metal/Vulkan rows change.
5. Version ship postmortem should note matrix deltas ([postmortems/](../../planning/postmortems/)).

---
**Layer:** architecture  
**Last updated:** 2026-06-23  
**Matrix version:** 1.5 (ADR-006 — Windows CUDA now uses runtime DLL loading; single binary probes for `ggml-cuda.dll`; struct layouts pinned to b9553)
