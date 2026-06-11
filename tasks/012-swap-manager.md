# Task 012 — Swap Manager

## Status

Complete

## Goal

Move process state and KV blocks across VRAM → RAM memory tiers.

## Context

Memory tiering per [architecture/memory.md](../docs/architecture/memory.md). Month 3 Week 12.

## Requirements

- `internal/swap/` with `Manager` interface and `StubManager`
- `SwapOut`: KV VRAM→RAM via `kv.Save`, process → HIBERNATED, `SwapStarted`/`SwapCompleted`
- `SwapIn`: KV RAM→VRAM via `kv.Load`, process → READY
- Enforce scheduler idle heuristic on swap-out (`scheduler.ShouldSwap`)
- No scheduler policy changes

## Files

- `internal/swap/manager.go`
- `internal/swap/manager_test.go`

## Acceptance Criteria

- [x] SwapOut moves KV to RAM and hibernates process
- [x] SwapIn restores KV to VRAM and sets process READY
- [x] SwapStarted/SwapCompleted emitted with tier labels (VRAM/RAM)
- [x] Running and not-idle processes rejected on SwapOut
- [x] No inference imports (INV-008)

## Out Of Scope

- NVMe cold tier (post-v1.0)
- Scheduler policy changes
- Model weight unload (runtime)

## Related

- Month plan: [planning/roadmap/months/month-03.md](../planning/roadmap/months/month-03.md)
- Task: [011-kv-manager.md](011-kv-manager.md)

---
**Layer:** task
**Target month:** 3
