# ADR-006: Windows GPU Backend via Runtime DLL Loading

## Status

Accepted — 2026-06-23

## Context

TinyBrain OS requires CUDA GPU inference on Windows (RTX hardware, ADR-001).
llama.cpp's CUDA code is compiled by `nvcc`, which on Windows mandates MSVC
(`cl.exe`) as the host C++ compiler.  Go's CGO toolchain on Windows uses
MinGW-GCC.  These two ABI families **cannot be statically linked** into the
same binary:

| Requirement | Toolchain |
|-------------|-----------|
| Go CGO | MinGW-GCC (produces `.a` archives) |
| nvcc host compiler on Windows | MSVC (produces `.lib` / `.dll`) |

Attempts to link MSVC-built `.lib` files against MinGW-GCC CGO fail with
`undefined reference to GOMP_parallel` and similar ABI boundary errors.

Static CGO works on Linux CUDA because `nvcc` accepts GCC as host compiler,
producing GNU-compatible archives.

## Decision

On Windows, TinyBrain will load the llama.cpp inference library **at runtime**
using `syscall.LoadDLL` / `syscall.NewProc` rather than via static CGO linkage.

### Consequences

**Positive**

1. **ABI decoupling** — The Go scheduler and runtime kernel are never compiled
   against MSVC symbols.  Only `syscall.Proc.Call()` crosses the boundary.

2. **Fault isolation** — If `ggml-cuda.dll` faults or is absent, the DLL
   loader catches the error at `LoadDLL` time and gracefully degrades to the
   CPU backend.  The Go scheduler is never taken down by a CUDA driver crash.

3. **Forward-compatible pattern** — Metal (macOS), ROCm (Linux/Windows), and
   Vulkan can all follow the same dynamic probe pattern:
   `backend_macos_metal.go`, `backend_linux_rocm.go`, `backend_windows_vulkan.go`.
   No additional build tags or static link complexity is required.

4. **Distribution readiness** — End-users ship `llama.dll` + `ggml-cuda.dll`
   alongside the TinyBrain binary.  No MinGW runtime, no Visual C++ Redistributable
   dependency on the Go side.  The MSVC runtime is contained inside the DLL.

5. **Single binary** — The Windows binary probes for `ggml-cuda.dll` at startup
   and selects the GPU or CPU backend automatically.  No separate `-tags cuda`
   build is needed on Windows.

**Negative / Trade-offs**

- C struct layouts in `backend_windows_dynamic.go` must be kept in sync with
  llama.cpp's `llama.h` / `ggml.h` ABI.  They are pinned to build tag b9553
  (`9e3b928`).  A submodule update requires re-verification of struct offsets.

- `syscall.Proc.Call()` does not provide Go's native type safety.  Marshalling
  errors are runtime panics.  The `decodeBatch` helper uses a fixed-size stack
  buffer to receive the `llama_batch` struct returned by value.

- CI cannot automatically test the GPU path without a GPU runner.
  Manual verification via `TB_CUDA_INTEGRATION=1` is required (see
  `009d-manual-gpu-checklist.md`).

## Implementation

| File | Purpose |
|------|---------|
| `internal/inference/llama/backend.go` | Private `backend` interface |
| `internal/inference/llama/backend_cgo.go` | CGO backend (CPU; Linux/macOS CUDA) |
| `internal/inference/llama/backend_windows_dynamic.go` | Windows runtime DLL loader |
| `internal/inference/llama/provider_windows.go` | Probes for `ggml-cuda.dll`; selects backend |
| `internal/inference/llama/provider_non_windows.go` | Always returns `cgoBackend` on non-Windows |
| `internal/inference/llama/generate_stub.go` | `stubBackend` for `!cgo` builds |

### Probe algorithm

```
startup → provider_windows.go:selectBackend()
  → NewWindowsDynamicBackend(dllDir)
      → stat ggml-cuda.dll          (ErrDLLNotFound → fallback)
      → LoadDLL(ggml-base.dll)
      → LoadDLL(ggml-cuda.dll)      (optional; non-fatal)
      → LoadDLL(llama.dll)
      → resolveProcs()              (all exported symbols)
  → on error: return &cgoBackend{} (CPU fallback)
```

### DLL search order

1. `LLAMACPP_DLL_DIR` environment variable (CI override)
2. Directory of the running executable

## Alternatives Rejected

| Alternative | Reason rejected |
|-------------|----------------|
| Use MinGW GCC as nvcc host on Windows | `nvcc` on Windows requires `cl.exe`; not supported |
| Import `.lib` via `dlltool` `.a` wrappers | Produces duplicate symbol errors when both MSVC and MinGW runtimes are linked |
| WSL2 passthrough | Requires WSL installed; not a standalone Windows experience |
| Ship pre-built Go plugin (`.so`) | Go plugins are Linux-only |

Rejected approaches are logged in [planning/decisions/rejected.md](../../planning/decisions/rejected.md).

---
**Layer:** decision  
**Related:** [ADR-001](ADR-001-Hardware-Aware-Runtime.md), [inference-backend-matrix.md](../architecture/inference-backend-matrix.md), [cross-platform.md](../architecture/cross-platform.md)  
**Supercedes:** None — extends 009d CUDA adapter with Windows DLL strategy  
**Last updated:** 2026-06-23
