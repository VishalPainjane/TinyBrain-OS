# Task 011 — KV Manager

## Status

Complete

## Goal

Implement KV cache save/load stub with KVStored and KVLoaded events.

## Context

Context hibernation requires a KV block manager. Month 3 Week 11. See [RFC-001](../../docs/rfc/RFC-001-KV-Hibernation.md).

## Requirements

- `internal/kv/` package with `Manager` interface and `StubManager`
- KV block pool metadata: ID, PID, size, tier (VRAM/RAM/NVMe), last access
- `Allocate`, `Save` (VRAM→RAM), `Load` (RAM→VRAM), `Get`, `Delete`
- Publish `KVStored` / `KVLoaded` on event bus
- No llama.cpp KV export (metadata-only stub)

## Files

- `internal/kv/types.go`
- `internal/kv/manager.go`
- `internal/kv/manager_test.go`

## Acceptance Criteria

- [x] Allocate registers block in VRAM
- [x] Save moves block to RAM and emits KVStored
- [x] Load moves block to VRAM and emits KVLoaded
- [x] Tests demonstrate event delivery on bus
- [x] No inference imports (INV-008)

## Out Of Scope

- Full FP16→Q4 compression pipeline (post-v1.0)
- Real llama.cpp KV export (deferred)
- NVMe tier transitions (stub tier type only)

## Related

- Month plan: [planning/roadmap/months/month-03.md](../planning/roadmap/months/month-03.md)
- Architecture: [docs/architecture/memory.md](../docs/architecture/memory.md)
- Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)

---
**Layer:** task
**Target month:** 3
