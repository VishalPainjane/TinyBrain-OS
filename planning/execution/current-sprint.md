# Current Sprint

## Sprint

**Name:** Month 3 — KV and Swap Foundations
**Goal:** KV manager stub, swap manager, brain-top prototype
**Version:** pre-v1.0 memory subsystems (v0.7 scheduler shipped)

## Current Task

[013-brain-top](../../tasks/013-brain-top.md) — next

## In Progress

(none)

## Blocked

(none)

## Done (prior sprint)

- [x] 011-kv-manager — merged PR #13 @ `1a1acf0`
- [x] v0.7 MLFQ scheduler (010) — merged PR #12 @ `a0f90a7`
- [x] v0.6 inference (009a–009d) — tag `v0.6`

## Next

- 013-brain-top (prototype)

## Definition Of Done

- Swap manager records VRAM→RAM tier movement with lifecycle events
- `CGO_ENABLED=0 go test ./...` passes
- Boundary tests pass
- Task moved to [completed.md](completed.md) with commit hash
- [docs/current.md](../../docs/current.md) synced

## Forbidden Work

- `internal/agents/`
- Kubernetes
- Web dashboard
- API/Router layer

## Active References

- Task: [tasks/012-swap-manager.md](../../tasks/012-swap-manager.md)
- Architecture: [docs/architecture/memory.md](../../docs/architecture/memory.md)

---
**Layer:** planning
**Last updated:** 2026-06-11
