# Architecture Review — Task 009d (GPU Offload — CUDA)

**Status:** Approved for implementation (2026-06-09)  
**Date:** 2026-06-09  
**Scope:** `internal/inference/llama/` — CUDA layer offload via `-tags cuda`. Refactor CGO bindings for backend-specific LDFLAGS. No runtime, scheduler, loader, registry, contract, or lifecycle doc changes.

| Gate | Status |
|------|--------|
| Architecture review | **Approved** |
| API compatibility (b9553 @ 9e3b928) | **Approved** (see §5) |
| Build tag matrix (009a) | **Approved** — unchanged mutual exclusivity |
| Code implementation | **Not started** |

**Pre-read:** [009a-architecture-review.md](009a-architecture-review.md), [009b-architecture-review.md](009b-architecture-review.md), [009c-architecture-review.md](009c-architecture-review.md), [009a-build-tags.md](009a-build-tags.md), [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md), [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md), [cross-platform.md](../../docs/architecture/cross-platform.md), [v0.6-inference.md](../../docs/specs/v0.6-inference.md)

**Submodule pin (unchanged from 009a–009c):**

| Field | Value |
|-------|-------|
| Tag | `b9553` |
| SHA | `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66` |
| Header source | `third_party/llama.cpp/include/llama.h` |

**Gates inherited (unchanged):** A Lifecycle, B Resolver (Option B in runtime), C Dependency, D Build tags, E API pin b9553.

---

## Executive Summary

009d completes v0.6 by adding a **compile-time optional** CUDA backend isolated in `internal/inference/llama/`. The CPU path remains the CI merge blocker. GPU offload is controlled exclusively through `LlamaConfig.NGLayers`, applied to `llama_model_params.n_gpu_layers` at load time in the CUDA-tagged binding file. Hardware probe (`hardware.BackendCUDA`) informs composition-root config; core packages do not import CUDA or set layer counts.

**Critical refactor:** Today `bindings_cgo.go` is tagged `cgo && !cuda && !rocm && !metal && !vulkan` and hardcodes `n_gpu_layers = 0`. 009d splits shared CGO helpers from backend-specific LDFLAGS and load parameterization.

---

## 1. Exact Files to Create

| File | Build constraint | Purpose |
|------|------------------|---------|
| `tasks/009d-gpu-offload-cuda.md` | — | Task spec |
| `internal/inference/llama/load_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA `loadBackend` / `unloadBackend` |
| `internal/inference/llama/generate_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA `generateBackend` (delegates to shared native generate) |
| `internal/inference/llama/bindings_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA `#cgo LDFLAGS` + `loadNativeModel` with `cfg.NGLayers` |
| `internal/inference/llama/config_test.go` | — | Table-driven `NGLayers` defaults and clamp rules |
| `internal/inference/llama/config_probe.go` | — | Optional `ConfigFromProbe(hardware.ProbeResult) LlamaConfig` helper |
| `internal/inference/llama/load_cuda_test.go` | `cuda` or build-tag-agnostic config tests | Config mapping tests (no GPU required) |
| `planning/decisions/009d-manual-gpu-checklist.md` | — | Manual verification steps for merge sign-off |

**Not created:** `cmd/`, new top-level packages, runtime wiring, probe changes, `load_rocm.go`, `load_metal.go`, CI GPU integration job.

---

## 2. Exact Files to Modify

| File | Change |
|------|--------|
| `internal/inference/llama/config.go` | Add `NGLayers int32`; document semantics; update `DefaultConfig()` (`NGLayers: 0`) |
| `internal/inference/llama/bindings_cgo.go` | **Rename/split** → see §4 refactor plan. CPU file keeps `n_gpu_layers = 0`; extract shared generate/tokenize/decode helpers |
| `internal/inference/llama/doc.go` | Reference 009d review; note `-tags cuda` |
| `internal/inference/llama/provider.go` | Comment only: provider is backend-agnostic; backend selected by build tag |
| `README.md` | CUDA cmake flags, `-tags cuda` build, `TB_NGLAYERS`, DLL/PATH notes (Windows/Linux) |
| `.github/workflows/ci.yml` | Add optional `inference-cuda-compile` job (compile-only, no GPU) |
| `.gitignore` | Ensure CUDA cmake artifacts under `third_party/llama.cpp/build/` remain ignored |

**Pre-implementation doc update (required before first CUDA code commit):**

| File | Change |
|------|--------|
| `docs/architecture/inference-backend-matrix.md` | CUDA rows: **Planned → In Progress** with build/CI notes |

**Post-merge planning sync (merge commit, not blocking review):**

| File | Change |
|------|--------|
| `docs/specs/v0.6-inference.md` | GPU offload checkbox |
| `planning/releases/v0.6.md` | Feature complete |
| `planning/execution/completed.md` | 009d entry |
| `planning/execution/current-sprint.md` | Sprint done / v0.6 tag prep |
| `docs/current.md` | v0.6 complete or next task |
| `planning/architecture-evolution/current-state.md` | CUDA partial status |
| `planning/decisions/accepted.md` | Row for 009d CUDA compile-only CI if confirmed at merge |

**Explicitly zero-diff (enforced):**

| Path | Reason |
|------|--------|
| `internal/runtime/` | User constraint; GPU is adapter-only |
| `internal/scheduler/` | INV-001 |
| `internal/loader/` | User constraint |
| `internal/registry/` | User constraint |
| `internal/process/`, `internal/events/` | Out of scope |
| `docs/contracts/*` | `InferenceProvider` unchanged |
| `docs/architecture/inference-lifecycle.md` | Lifecycle ownership unchanged |
| `internal/hardware/probe.go` | Probe already detects CUDA; no 009d change required |

---

## 3. CUDA Build-Tag Strategy

Per [009a-build-tags.md](009a-build-tags.md) — **no changes to mutual exclusivity rules.**

### Tag matrix (009d)

| File | Build constraint | Backend |
|------|------------------|---------|
| `provider_stub.go` | `!cgo` | Go stub |
| `bindings_common.go` *(new name)* | `cgo` | Shared CFLAGS + `#include`; **no LDFLAGS** |
| `bindings_cpu.go` *(from bindings_cgo.go)* | `cgo && !cuda && !rocm && !metal && !vulkan` | CPU LDFLAGS + CPU load (`n_gpu_layers=0`) |
| `bindings_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA LDFLAGS + CUDA load (`cfg.NGLayers`) |
| `load_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` | CPU backend dispatch |
| `load_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA backend dispatch |
| `generate_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` | CPU generate dispatch |
| `generate_cuda.go` | `cgo && cuda && !rocm && !metal && !vulkan` | CUDA generate dispatch |
| `provider.go`, `config.go`, `errors.go`, `context.go` | *(none)* | All builds |
| `*_integration_test.go` | `cgo && integration` | Uses **active backend binary** |

### Build commands

| Target | Command |
|--------|---------|
| CI default (stub) | `CGO_ENABLED=0 go test ./...` |
| Linux CPU (merge blocker) | `CGO_ENABLED=1 go test ./internal/inference/llama/...` |
| Linux CUDA (local / compile CI) | `CGO_ENABLED=1 go build -tags cuda ./internal/inference/llama/...` |
| Linux CUDA tests (compile) | `CGO_ENABLED=1 go test -tags cuda ./internal/inference/llama/...` *(unit/config only; GPU tests skipped)* |
| Integration (CPU binary) | `CGO_ENABLED=1 go test -tags integration ./...` |
| Integration (CUDA binary, manual GPU) | `CGO_ENABLED=1 go test -tags 'cuda integration' ./internal/inference/llama/...` |

### llama.cpp cmake (CUDA)

```bash
cmake -S third_party/llama.cpp -B third_party/llama.cpp/build-cuda \
  -DCMAKE_BUILD_TYPE=Release \
  -DLLAMA_BUILD_TESTS=OFF \
  -DLLAMA_BUILD_TOOLS=OFF \
  -DLLAMA_BUILD_EXAMPLES=OFF \
  -DLLAMA_BUILD_SERVER=OFF \
  -DLLAMA_BUILD_APP=OFF \
  -DLLAMA_BUILD_UI=OFF \
  -DLLAMA_BUILD_COMMON=OFF \
  -DLLAMA_CURL=OFF \
  -DGGML_CUDA=ON \
  -DGGML_METAL=OFF \
  -DGGML_VULKAN=OFF \
  -DGGML_HIP=OFF \
  -DGGML_CCACHE=OFF
cmake --build third_party/llama.cpp/build-cuda --target llama -j
```

**Separate build directory** (`build-cuda` vs `build`) avoids clobbering CPU CI artifacts on dev machines. `#cgo LDFLAGS` in `bindings_cuda.go` must point at the directory used for CUDA libs (document both in README; implementation verifies exact `.so` names at pin).

### Verification before merge

```bash
# Must pass (unchanged)
CGO_ENABLED=0 go test ./...

# Must pass (unchanged CPU job)
CGO_ENABLED=1 go test ./internal/inference/llama/...

# Must pass (new, compile-only)
CGO_ENABLED=1 go test -tags cuda ./internal/inference/llama/... -run TestDoesNotRequireGPU

# Must NOT be merge blocker
CGO_ENABLED=1 go test -tags 'cuda integration' ./internal/inference/llama/...  # manual + GPU
```

---

## 4. CGO Refactor Plan (bindings split)

**Problem:** `bindings_cgo.go` combines CPU `#cgo LDFLAGS`, load with hardcoded `n_gpu_layers=0`, and shared generate/tokenize helpers — all behind the CPU-only build tag. CUDA cannot link without duplicating generate logic or violating [009a-build-tags.md](009a-build-tags.md) ("no GPU LDFLAGS in shared bindings").

**Approved split:**

```text
bindings_common.go     # tag: cgo
  - #cgo CFLAGS only (include paths)
  - #include llama.h
  - nativeMu, nativeModels, nativeContexts maps
  - initBackend() → llama_backend_init()
  - tokenizeNative, decodeNative, tokenToPieceNative
  - newNativeSampler, generateNative, generateNativeTimed
  - freeNativeHandles, unloadNativeModel
  - effectiveBatchSize

bindings_cpu.go        # tag: cgo && !cuda && !rocm && !metal && !vulkan
  - #cgo LDFLAGS: -lllama -lggml -lggml-cpu -lggml-base ... (current)
  - loadNativeModel(path, modelID, cfg) with params.n_gpu_layers = 0

bindings_cuda.go       # tag: cgo && cuda && !rocm && !metal && !vulkan
  - #cgo LDFLAGS: -lllama -lggml -lggml-cuda -lggml-base ... (+ CUDA deps)
  - loadNativeModel(path, modelID, cfg) with params.n_gpu_layers = cfg.NGLayers
  - optional: params.main_gpu = 0 (single-GPU v0.6)
```

**Generate path:** Unchanged API surface — `generateNative*` stays in `bindings_common.go` because decode/sampler APIs are backend-agnostic once the model is loaded on GPU.

**Delete/rename:** Remove monolithic `bindings_cgo.go` after split (git mv to `bindings_cpu.go` + extract common).

---

## 5. llama.cpp GPU API Compatibility Review (Pin b9553 @ 9e3b928)

Verified against upstream `llama.h` at SHA `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66` (same pin as 009b). **Do not use deprecated or unpinned APIs.**

### Summary

| Area | API | Pinned status | 009d usage |
|------|-----|---------------|------------|
| Backend init | `llama_backend_init()` | **Present** | Once via `initBackend()` — registers CUDA when compiled with `GGML_CUDA` |
| Model params | `llama_model_default_params()` | **Present** | Base for load |
| GPU layers | `llama_model_params.n_gpu_layers` (`int32_t`) | **Present** | Set from `LlamaConfig.NGLayers` in CUDA bindings only |
| mmap | `llama_model_params.use_mmap` | **Present** | Unchanged from CPU |
| Context create | `llama_init_from_model()` | **Present** | Unchanged (009b) |
| Memory clear | `llama_memory_clear(llama_get_memory(ctx), true)` | **Present** | Unchanged (009b) |
| Decode / sampler | `llama_decode`, `llama_sampler_*` | **Present** | Unchanged (009b) |
| Backend free | `llama_backend_free()` | **Present** | Optional on process exit; not required per load cycle |
| KV cache clear | `llama_kv_cache_clear` | **Absent** | **Do not use** (009b rule) |

### GPU-specific struct fields (b9553)

```c
struct llama_model_params {
    // ...
    int32_t n_gpu_layers; // number of layers to store in VRAM
    // ...
    int32_t main_gpu;     // single-GPU default 0 — set explicitly in 009d
    // ...
};
```

**Semantics (llama.cpp convention at pin):**

| `NGLayers` value | Behavior |
|------------------|----------|
| `0` | All layers on CPU (same as current 009b) |
| `1..N-1` | Offload that many layers to GPU |
| `-1` | Offload all layers to GPU (recommended default when `BackendCUDA` + cuda binary) |
| `> model layer count` | llama.cpp clamps — treat as success; document in tests |

**Not used in 009d (explicit deferral):**

| Field | Reason |
|-------|--------|
| `split_mode`, `tensor_split` | Multi-GPU out of scope |
| `main_gpu != 0` | Single GPU index 0 only |
| `llama_model_params.devices` | Device list API — defer to future probe task |
| Custom `ggml_backend_cuda_*` direct calls | Violates adapter boundary — use llama.h only |

### CUDA compile-time linkage

When `GGML_CUDA=ON`, cmake produces additional targets/libs (names vary by pin — verify locally):

- `libggml-cuda.so` (or `.dll`)
- CPU ggml may still exist; CUDA binary must link CUDA ggml backend

**Implementation gate:** First 009d commit must record exact `LDFLAGS` library list from `third_party/llama.cpp/build-cuda/bin` in `bindings_cuda.go` comments and README.

### API compatibility verdict

**Approved** — 009b generate/load APIs remain valid for CUDA; 009d adds only `n_gpu_layers` (+ optional `main_gpu`) on the existing load path. No new llama.h symbols required for Generate.

---

## 6. NGLayers Ownership

| Concern | Owner | 009d behavior |
|---------|-------|---------------|
| **Detection** — is CUDA available on host? | `internal/hardware` (`OSProber`, `nvidia-smi`) | **No change** — already sets `BackendCUDA` |
| **Configuration field** | `internal/inference/llama/config.go` (`LlamaConfig.NGLayers`) | **Add field** |
| **Default value** | `DefaultConfig()` | `NGLayers: 0` (safe CPU behavior) |
| **Policy mapping** — probe → layer count | **Composition root** (app/test/main), optional `ConfigFromProbe` helper in inference | When `BackendCUDA`: recommend `NGLayers: -1`; else `0` |
| **Application to engine** | `bindings_cuda.go` → `params.n_gpu_layers` | Only in `-tags cuda` binary |
| **Enforcement in CPU binary** | `bindings_cpu.go` | Always `n_gpu_layers = 0` — ignores config GPU values |
| **Runtime / loader / scheduler / registry** | **None** | Zero diffs — config passed via existing `NewLlamaProvider(resolver, cfg)` |

### Configuration sources (priority order)

1. Explicit `LlamaConfig` passed to `NewLlamaProvider` (highest)
2. Optional env at composition root only: `TB_NGLAYERS` (document in README; parsed in test helper or future `cmd/`, **not** in core packages)
3. `ConfigFromProbe(probe)` helper defaults:
   - `BackendCUDA` → `NGLayers: -1`
   - `BackendCPU` / `BackendMetal` / `BackendROCm` → `NGLayers: 0`
4. `DefaultConfig()` → `NGLayers: 0`

### Example composition (unchanged runtime API)

```go
profile, _ := hardware.ProbeAndClassify()
cfg := llama.DefaultConfig()
if profile.Probe.Backend == hardware.BackendCUDA {
    cfg.NGLayers = -1
}
provider := llama.NewLlamaProvider(resolver, cfg)
rt := runtime.NewIntegratedModelRuntime(provider, loader, resolver, bus)
```

**Rule:** Mismatch (CUDA probe + CPU binary) is safe — CPU binary forces `n_gpu_layers=0`. Opposite (CPU probe + CUDA binary + `NGLayers=-1`) may fail load or run slow — caller responsibility; document in README.

---

## 7. Configuration Changes

| Change | Location | Details |
|--------|----------|---------|
| `LlamaConfig.NGLayers int32` | `config.go` | Exported; comment references matrix |
| `DefaultConfig()` | `config.go` | `NGLayers: 0` |
| `ConfigFromProbe` *(optional)* | `config_probe.go` | Maps `hardware.ProbeResult.Backend` → recommended `NGLayers` |
| `TB_NGLAYERS` | README + manual checklist | Optional env for integration tests on GPU machines |
| No new registry fields | — | Quantization/path unchanged |
| No runtime config struct | — | Out of scope |

**Validation rules (unit tested):**

- Negative values only `-1` (all layers) or reject `< -1` with clamp to `-1`
- CPU build ignores positive `NGLayers` at load (hardcoded 0 in bindings)
- Document that `NGLayers` has no effect unless binary built with `-tags cuda`

---

## 8. CI Strategy (Without GPU Runners)

| Job | Runner | CGO | Tags | Purpose | Merge blocker? |
|-----|--------|-----|------|---------|----------------|
| `test` | `ubuntu-latest` | `0` | — | Full stub path | **Yes** |
| `inference-cgo` | `ubuntu-latest` | `1` | — | CPU llama build + tests | **Yes** |
| `inference-cuda-compile` *(new)* | `ubuntu-latest` + `nvidia/cuda:12.x-devel-ubuntu22.04` container OR toolkit install step | `1` | `cuda` | cmake `GGML_CUDA=ON`; `go test -tags cuda` config/unit tests; `go build -tags cuda` | **No** — advisory first PR; promote to required after green history |
| `inference-integration` | `ubuntu-latest` | `1` | `integration` | CPU E2E when `TB_TEST_GGUF_PATH` set | No |
| GPU integration | Dev machine | `1` | `cuda integration` | Manual checklist | No |

### `inference-cuda-compile` job design

```yaml
# Pseudocode — implementation in ci.yml
inference-cuda-compile:
  runs-on: ubuntu-latest
  container:
    image: nvidia/cuda:12.4.1-devel-ubuntu22.04
  steps:
    - checkout + submodules
    - install go 1.22, cmake, build-essential
    - cmake GGML_CUDA=ON → build-cuda/
    - go test -tags cuda ./internal/inference/llama/... -count=1
    - go build -tags cuda ./internal/inference/llama/...
```

**Constraints:**

- No `nvidia-smi`, no GPU device plugin — compile and unit tests only
- Skip tests requiring GPU with `t.Skip` unless `TB_CUDA_INTEGRATION=1`
- CPU `inference-cgo` job **unchanged** — still builds `-DGGML_CUDA=OFF`
- Cache: separate cache key suffix `-cuda` for `build-cuda/`

### What CI proves vs manual

| Claim | CI | Manual GPU |
|-------|-----|------------|
| CUDA code compiles | Yes | — |
| CPU path unbroken | Yes | — |
| `NGLayers` wired in CGO | Unit tests (mock/clamp) | — |
| Tokens from GPU offload | No | Yes |
| TPS > CPU | No | Yes (document) |

---

## 9. Windows Strategy

| Topic | Decision |
|-------|----------|
| Toolchain | MSVC + NVIDIA CUDA Toolkit (matching llama.cpp GGML_CUDA build) |
| Build | Separate `build-cuda` directory; `set CGO_ENABLED=1` && `go build -tags cuda ./internal/inference/llama/...` |
| Runtime DLLs | **Not bundled** — user must have CUDA runtime in PATH (`cudart64_*.dll`, etc.) |
| `#cgo LDFLAGS` | Windows block in `bindings_cuda.go` — `.lib` paths from cuda build dir |
| Probe | Existing `nvidia-smi` detection on Windows |
| CI | Not blocking — optional future `windows-latest` cuda **compile** job deferred |
| Verification | Manual checklist: load SmolLM with `NGLayers=-1`, compare TPS vs CPU |
| mmap | Same as 009b — manual validation |

---

## 10. Linux Strategy

| Topic | Decision |
|-------|----------|
| Primary target | Yes — first CUDA implementation and manual GPU verification |
| CPU CI | Unchanged `inference-cgo` on `ubuntu-latest` |
| CUDA compile CI | `inference-cuda-compile` with devel container |
| Local GPU dev | WSL2 or native Linux with NVIDIA driver + `-tags cuda` binary |
| `LD_LIBRARY_PATH` | Must include `third_party/llama.cpp/build-cuda/bin` and CUDA lib64 |
| K8s (future) | NVIDIA device plugin + CUDA base image — deployment concern only; no scheduler changes |
| musl/Alpine | **Unsupported** — same as 009a |

---

## 11. macOS Impact

| Topic | Impact |
|-------|--------|
| CUDA | **N/A** — NVIDIA does not support CUDA on macOS (permanent, matrix row **No**) |
| 009d code changes | None required for Darwin-specific files |
| CPU macOS build | Unchanged — `-tags cuda` not used on macOS dev |
| Future GPU on Mac | **Metal** backend (separate task) — not 009d |
| Probe | `BackendMetal` inferred on Darwin — `ConfigFromProbe` returns `NGLayers: 0` until Metal task |
| CI | No macOS CUDA job |

**Verdict:** 009d is transparent to macOS — CPU-only path remains the macOS inference story until a future Metal task.

---

## 12. Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| AC-1 | v0.6 spec: GPU offload on CUDA-capable hardware | Manual GPU checklist + merge notes |
| AC-2 | `load_cuda.go` exists with correct build tag | Code review vs [009a-build-tags.md](009a-build-tags.md) |
| AC-3 | `params.n_gpu_layers` from `LlamaConfig.NGLayers` in CUDA binary | Code review + unit test |
| AC-4 | CPU binary forces `n_gpu_layers=0` | Code review + CPU integration unchanged |
| AC-5 | Generate returns tokens on GPU (manual) | `009d-manual-gpu-checklist.md` |
| AC-6 | TPS improvement vs CPU documented | Manual log in checklist (no SLA) |
| AC-7 | `CGO_ENABLED=0 go test ./...` passes | CI `test` |
| AC-8 | CPU `inference-cgo` passes | CI unchanged |
| AC-9 | CUDA compile job green (or documented deferral with issue) | CI `inference-cuda-compile` |
| AC-10 | Zero diffs in runtime/scheduler/loader/registry | `git diff` scope |
| AC-11 | Scheduler zero inference imports | Boundary test |
| AC-12 | INV-008 — CUDA only in `internal/inference/llama/` | Import audit |
| AC-13 | Matrix CUDA rows updated | Doc sync on merge |
| AC-14 | README CUDA build section | Doc review |
| AC-15 | Rollback: CPU binary unaffected | Rebuild without `cuda` tag |

---

## 13. Risks

| ID | Risk | Severity | Mitigation |
|----|------|----------|------------|
| R-1 | Wrong `#cgo LDFLAGS` for CUDA libs at pin | High | Verify against cmake output; document in bindings comment |
| R-2 | CUDA toolkit / driver version skew | High | Pin CUDA devel image version in CI; README minimum driver |
| R-3 | `bindings_cgo.go` split breaks CPU path | High | CPU CI merge blocker; refactor in dedicated commit |
| R-4 | OOM on GPU with `NGLayers=-1` + large model | Medium | Document profile guidance; default `-1` only for small models in tests |
| R-5 | WSL2 GPU passthrough quirks | Medium | Manual checklist notes; CPU fallback |
| R-6 | Windows DLL hell | Medium | Document PATH; no bundling |
| R-7 | Accidental CUDA import in core packages | High | Import boundary tests unchanged |
| R-8 | Scope creep into runtime auto-config | Medium | Enforce zero-diff on forbidden packages |
| R-9 | Compile CI flaky (container size/time) | Low | `timeout-minutes: 60`; advisory first |
| R-10 | b9553 CUDA cmake flag name drift | Medium | Verify `GGML_CUDA` at pin; submodule bump only if blocked |
| R-11 | Integration tests fail on cuda-tagged CI without GPU | Medium | Skip GPU tests unless `TB_CUDA_INTEGRATION=1` |
| R-12 | Two cmake build dirs confuse developers | Low | README table: CPU=`build/`, CUDA=`build-cuda/` |

---

## 14. Rollback Plan

| Level | Action | Result |
|-------|--------|--------|
| **L1 — Runtime** | Use CPU binary (build without `-tags cuda`) | Immediate CPU inference; ignore `NGLayers` |
| **L2 — Code revert** | Revert 009d merge commit | Restore monolithic `bindings_cgo.go` CPU-only |
| **L3 — CI** | Disable `inference-cuda-compile` job | CPU jobs continue; no merge impact |
| **L4 — Docs** | Revert matrix CUDA rows to **Planned** | Governance sync |
| **L5 — Release** | Ship v0.6 tag with CPU-only note if GPU blocked | Spec checkbox partial; defer GPU to v0.6.1 |

**Merge safety:** CPU paths are the merge blocker throughout — 009d cannot break default CI if refactor is correct.

---

## 15. Final Implementation Plan (Ordered)

### Phase 0 — Governance (pre-code)

1. Update matrix CUDA rows to **In Progress**.
2. Approve this review (gate closed).

### Phase 1 — Config

3. Add `NGLayers` to `LlamaConfig` + `config_test.go`.
4. Add optional `ConfigFromProbe` + tests.

### Phase 2 — Bindings refactor (CPU safety first)

5. Extract `bindings_common.go` from `bindings_cgo.go`.
6. Rename remainder to `bindings_cpu.go`; verify CPU CI green.
7. No behavior change in this phase.

### Phase 3 — CUDA backend

8. Build llama.cpp with `GGML_CUDA=ON`; record library names.
9. Add `bindings_cuda.go` with CUDA LDFLAGS + `loadNativeModel` using `cfg.NGLayers`.
10. Add `load_cuda.go`, `generate_cuda.go`.

### Phase 4 — CI + docs

11. Add `inference-cuda-compile` job (advisory).
12. Update README CUDA section.
13. Add `009d-manual-gpu-checklist.md`.

### Phase 5 — Manual GPU verification

14. Linux or Windows dev machine: `-tags cuda`, `NGLayers=-1`, SmolLM Q4_K_M.
15. Record TPS vs CPU in checklist.
16. Optional: `TB_CUDA_INTEGRATION=1` integration test guarded skip.

### Phase 6 — Merge sync

17. Matrix CUDA → **Partial** (manual GPU) or **Yes** (if checklist complete).
18. Update v0.6 spec, release notes, `docs/current.md`, sprint docs.

---

## 16. Test Plan

### Unit (`CGO_ENABLED=0`)

| Test | Assert |
|------|--------|
| `TestDefaultConfig_NGLayersZero` | Default is 0 |
| `TestConfigFromProbe_CUDA` | CUDA → `-1` |
| `TestConfigFromProbe_CPU` | CPU → `0` |

### CGO CPU (`inference-cgo`)

| Test | Assert |
|------|--------|
| Existing 009a–009c tests | Unchanged green |
| Regression: load still `n_gpu_layers=0` | Optional compile-time review |

### CGO CUDA compile (no GPU)

| Test | Assert |
|------|--------|
| `TestLlamaConfig_NGLayersClamp` | Invalid values clamped |
| Package compiles with `-tags cuda` | CI job |

### Manual GPU (`cuda integration`, `TB_TEST_GGUF_PATH`, `TB_CUDA_INTEGRATION=1`)

| Test | Assert |
|------|--------|
| `TestLlamaProvider_LoadGenerate_CUDA` | Non-empty output |
| `TestLlamaProvider_CUDA_TPS` | Log TPS > CPU baseline (informational) |

---

## 17. Architecture Boundaries (unchanged from 009c)

```text
hardware.Probe() → BackendCUDA (detection only)
        │
        ▼
composition root sets LlamaConfig.NGLayers
        │
        ▼
llama.NewLlamaProvider(resolver, cfg)   ← 009d: cfg carries NGLayers
        │
        ▼
load_cuda.go → bindings_cuda.go → n_gpu_layers
        │
        ▼
generate_cuda.go → bindings_common.go → Generate (unchanged API)
```

**Forbidden edges:** runtime → cuda, scheduler → inference, loader → NGLayers, registry → GPU config.

---

**Layer:** planning  
**Related:** [009a-build-tags.md](009a-build-tags.md), [009b-architecture-review.md](009b-architecture-review.md), [009c-architecture-review.md](009c-architecture-review.md), [../../tasks/009d-gpu-offload-cuda.md](../../tasks/009d-gpu-offload-cuda.md)
