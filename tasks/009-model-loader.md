# Task 009 — Model Loader

## Status

Not Started

## Goal

Implement model loader responsible for load, unload, warm, prefetch, and evict operations.

## Context

Model loader sits behind runtime. Uses mmap for GGUF (future v0.6). Scheduler never touches loader directly — only via runtime.

## Requirements

- Loader interface: Load, Unload, Warm, Prefetch, Evict
- Track model lifecycle states: NOT_LOADED, LOADING, ACTIVE, WARM, UNLOADING, UNLOADED
- Stub implementation for v0.4; real mmap in v0.6
- LRU eviction hook when VRAM full (shell only until memory layer exists)

## Files

- `internal/loader/loader.go`
- `internal/loader/types.go`
- `internal/loader/loader_test.go`

## Acceptance Criteria

- [ ] Load transitions NOT_LOADED → ACTIVE
- [ ] Unload transitions ACTIVE → UNLOADED
- [ ] Duplicate load prevented
- [ ] Tests use stub files or mock paths

## Out Of Scope

- llama.cpp CGO (v0.6)
- KV cache management
- Scheduler-driven eviction policy

## Related

- Spec: [docs/specs/v0.4-runtime.md](../docs/specs/v0.4-runtime.md), [v0.6-inference.md](../docs/specs/v0.6-inference.md)
- Architecture: [docs/architecture/runtime.md](../docs/architecture/runtime.md)

---
**Layer:** task
