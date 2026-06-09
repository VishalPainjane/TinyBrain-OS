# Task 009d — GPU Offload (CUDA)

## Status

**Complete** — merged `ab06c60` (PR #9, 2026-06-10). CUDA matrix **Partial** until manual GPU checklist signed.

## Goal

Enable NVIDIA CUDA layer offload in the llama.cpp adapter (`-tags cuda`) without changing runtime, scheduler, loader, or registry packages. Close the remaining v0.6 spec item: GPU offload on CUDA-capable hardware.

## Context

009a–009c delivered CPU load/unload/generate wired to `ModelRuntime`. `load_cuda.go` was planned in 009a but never implemented; `bindings_cgo.go` hardcodes `n_gpu_layers = 0` and is tagged CPU-only. Hardware probe already detects `BackendCUDA` via `nvidia-smi`. NGLayers policy belongs in inference config + composition root, not core packages.

## Requirements

- Implement `load_cuda.go` + CUDA bindings per [009a-build-tags.md](../planning/decisions/009a-build-tags.md)
- Add `NGLayers` to `LlamaConfig`; apply at `llama_model_load_from_file` only in inference adapter
- Refactor CGO bindings: shared helpers vs backend-specific `#cgo LDFLAGS`
- llama.cpp built with `-DGGML_CUDA=ON` for CUDA binary; CPU CI path unchanged
- Optional compile-only CUDA CI job (no GPU runner)
- Manual GPU integration verification checklist (Linux/Windows)
- Update [inference-backend-matrix.md](../docs/architecture/inference-backend-matrix.md) on merge

## Files

See [planning/decisions/009d-architecture-review.md](../planning/decisions/009d-architecture-review.md) §1–§2.

## Acceptance Criteria

- [x] `CGO_ENABLED=0 go test ./...` passes (unchanged merge blocker)
- [x] `inference-cgo` CPU job passes (unchanged) — PR #9
- [ ] `go build -tags cuda ./internal/inference/llama/...` succeeds on Linux with CUDA toolkit + GGML_CUDA build — manual dev machine
- [x] CUDA binary sets `params.n_gpu_layers` from `LlamaConfig.NGLayers` (not hardcoded 0) — `bindings_cuda.go`
- [x] CPU binary (`!cuda` tag) still sets `NGLayers=0` regardless of config value — `bindings_cpu.go`
- [ ] Generate returns real tokens on CUDA-capable dev machine with `NGLayers > 0` (manual)
- [ ] TPS improvement vs CPU on same model/GPU documented in merge notes (no CI SLA)
- [x] `internal/runtime`, `internal/scheduler`, `internal/loader`, `internal/registry` — zero production diffs
- [x] Scheduler zero inference imports (unchanged boundary test)
- [x] INV-008 satisfied — CUDA CGO only under `internal/inference/llama/`
- [x] Matrix CUDA rows updated to **Partial** per verification level (manual GPU → **Yes**)

## Out Of Scope

- Runtime, scheduler, loader, registry package changes
- Metal, ROCm, Vulkan backends
- Multi-GPU, tensor split, device index > 0
- Bundling CUDA runtime DLLs/so
- CI GPU integration tests (no GPU runners)
- KV manager, SaveContext / RestoreContext
- Automatic probe → NGLayers wiring inside runtime (composition-root helper only)
- llama.cpp submodule pin bump (unless CUDA build blocked at b9553 — then pin bump is in-scope with ADR note)

## Related

- Review: [planning/decisions/009d-architecture-review.md](../planning/decisions/009d-architecture-review.md)
- Spec: [docs/specs/v0.6-inference.md](../docs/specs/v0.6-inference.md)
- Build tags: [planning/decisions/009a-build-tags.md](../planning/decisions/009a-build-tags.md)
- Dependency gate: [planning/decisions/009a-llama-cpp-dependency.md](../planning/decisions/009a-llama-cpp-dependency.md)
- Matrix: [docs/architecture/inference-backend-matrix.md](../docs/architecture/inference-backend-matrix.md)
- Cross-platform: [docs/architecture/cross-platform.md](../docs/architecture/cross-platform.md)

---
**Layer:** task
