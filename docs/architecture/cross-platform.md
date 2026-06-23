# Cross-Platform Architecture

TinyBrain OS must remain cross-platform. Development may occur on a single OS at a time; **architecture and package decisions must assume future support on Windows, Linux, and macOS**, and on CPU-only, NVIDIA CUDA, AMD ROCm/HIP, Apple Metal, and Vulkan-compatible GPUs.

Current implementation may target one platform or backend when required. Shortcuts are allowed only when portability risks are **documented** with a **migration path** — never silently accepted.

## Platforms

| OS | Status | Notes |
|----|--------|-------|
| Windows | Active development | CGO requires MSVC toolchain |
| Linux | Supported (CI default) | Primary server and Kubernetes target |
| macOS | Planned | Metal backend via inference adapter |

## Inference Backends

| Backend | Hardware | Windows isolation | Linux/macOS isolation |
|---------|----------|-------------------|-----------------------|
| CPU | All platforms | Static CGO (MinGW-GCC) | Static CGO (GCC/Clang) |
| CUDA | NVIDIA | **Runtime DLL loading** (`syscall.LoadDLL`) — [ADR-006](../adr/ADR-006-Windows-GPU-Dynamic-DLL-Backend.md) | Static CGO `-tags cuda` |
| ROCm/HIP | AMD | Runtime DLL loading (future — follows ADR-006 pattern) | Static CGO `-tags rocm` (future) |
| Metal | Apple Silicon / macOS GPU | N/A | CGO `-tags metal` (future) |
| Vulkan | Cross-vendor GPU | Runtime DLL loading (future) | Static CGO `-tags vulkan` (future) |

Backend selection flows from [hardware probe](hardware.md) → `hardware.Backend` → `internal/inference/llama/provider_windows.go:selectBackend()` (Windows) or `provider_non_windows.go:selectBackend()`. No backend logic in core packages.

**Capability truth:** [inference-backend-matrix.md](inference-backend-matrix.md) — update before implementing any backend; feature parity is not assumed.

## Rules

### 1. No hardcoded OS paths

- Use `filepath` / `path/filepath` for path operations.
- Registry `ModelDefinition.Path` is host-local; never assume `/` or `\` or drive letters in core packages.
- Config paths come from env, flags, or registry — not compile-time OS constants in `internal/process`, `scheduler`, `runtime`, `registry`, or `loader`.

### 2. No hardcoded platform assumptions

- Do not assume CUDA is present, that GPU exists, or that mmap behaves identically — probe and configure.
- Do not assume little-endian, single GPU, or little-core count without detection.
- `runtime.GOOS` / `GOARCH` may appear only in OS-specific probe or inference adapter files — not in core scheduling or process logic.

### 3. Prefer portable boundaries

- **Interfaces** at package edges (`InferenceProvider`, `Loader`, `ModelStore`, `Prober`).
- **Adapters** for persistence, inference backends, and OS probes.
- **Build tags** for CGO, CUDA, Metal, ROCm, Vulkan, and integration tests.
- **Configuration** for paths, layer counts, and backend choice — not `#ifdef` scattered through core code.

### 4. Platform-specific code isolation

Platform-specific code lives only in:

| Mechanism | Example |
|-----------|---------|
| Build tags | `//go:build cgo`, `//go:build cuda`, `//go:build windows` |
| Runtime DLL loading (Windows GPU) | `backend_windows_dynamic.go` via `syscall.LoadDLL` |
| Backend probe files | `provider_windows.go`, `provider_non_windows.go` |
| Adapter packages | `internal/inference/llama/` |
| OS-specific files | `probe_windows.go`, `probe_unix.go` |

> [!NOTE]
> On Windows, GPU backends (CUDA, future ROCm/Vulkan) use **runtime DLL loading** rather
> than static CGO because MSVC-compiled DLLs are ABI-incompatible with MinGW-GCC CGO.
> This is established in [ADR-006](../adr/ADR-006-Windows-GPU-Dynamic-DLL-Backend.md).
> Future GPU backends on Windows MUST follow this pattern.

Forbidden: platform branches inside `internal/process`, `internal/scheduler`, `internal/runtime`, `internal/registry`, or `internal/loader`.

### 5. Core packages remain platform agnostic

These packages must compile and test with `CGO_ENABLED=0` and must not import inference or OS-specific CGO:

- `internal/process`
- `internal/scheduler`
- `internal/runtime`
- `internal/registry`
- `internal/loader`

Hardware detection is the exception: `internal/hardware` uses OS-specific **probe files** behind a stable `Prober` interface.

### 6. Inference backends isolated in `internal/inference/`

All llama.cpp, CUDA, Metal, ROCm, Vulkan, and future cloud SDK imports exist **only** under `internal/inference/` (INV-008). Subpackages per backend are encouraged as backends multiply (e.g. `internal/inference/llama/backends/`).

### 7. Platform impact before new dependencies

Before merging a platform-specific dependency, document in `planning/decisions/`:

| Impact area | Questions |
|-------------|-----------|
| Linux | CI, containers, glibc/musl, GPU drivers |
| macOS | Xcode, Metal SDK, code signing, notarization |
| Windows | MSVC, CUDA toolkit path, DLL deployment |
| Kubernetes | GPU node selectors, device plugins, image size, CGO in containers, multi-arch builds |

See [009a-llama-cpp-dependency.md](../../planning/decisions/009a-llama-cpp-dependency.md) for the llama.cpp template.

### 8. Evolve toward multiple GPU APIs without rewrites

Inference adapter layout must allow adding backends without changing:

- `internal/runtime` (`InferenceProvider` port)
- `internal/scheduler`
- `internal/registry`

Pattern:

```text
internal/inference/llama/
  provider.go          # implements runtime.InferenceProvider
  config.go            # backend-agnostic settings
  load_cpu.go          # !cuda && !metal && !rocm
  load_cuda.go         # cuda
  load_metal.go        # darwin && metal  (future)
  load_rocm.go         # rocm             (future)
  load_vulkan.go       # vulkan           (future)
```

Shared load/unload orchestration stays in `provider.go`; backend files supply device init and layer offload only.

### 9. Portability first, optimization second

- Ship CPU-only path that passes CI everywhere before enabling GPU build tags.
- Default `NGLayers = 0` when backend unknown or probe fails.
- Optional GPU paths are additive — never replace the portable baseline.

### 10. Document portability risks

If a short-term choice limits future platforms, record in `planning/decisions/` or `planning/risks/technical-risks.md`:

- What is limited
- Why the shortcut was taken
- Migration path (file, interface, or build-tag plan)
- Target version to resolve

## Kubernetes (future)

Kernel packages (`process`, `scheduler`, `runtime`, `registry`, `events`) remain pure Go — suitable for minimal container images without CGO.

Inference workloads may run:

- **In-process** (single binary with CGO + GPU base image), or
- **Sidecar / remote** `InferenceProvider` adapter (future RFC) — scheduler and runtime port unchanged per ADR-004.

GPU nodes require device plugins and larger images; that is an **deployment** concern, not a reason to import CUDA into `internal/scheduler`.

## Related

- [inference-lifecycle.md](inference-lifecycle.md) — model state ownership and event emission
- [inference-backend-matrix.md](inference-backend-matrix.md) — per-backend capabilities, OS, CI status
- [hardware.md](hardware.md) — probe and `Backend` enum
- [invariants.md](invariants.md) — INV-008 inference isolation
- [runtime.md](runtime.md) — `InferenceProvider` port
- ADR-001, ADR-004, ADR-005
- [migration-paths.md](../../planning/architecture-evolution/migration-paths.md) Path 2
- [technical-risks.md](../../planning/risks/technical-risks.md)

---
**Layer:** architecture
**Last updated:** 2026-06-08
