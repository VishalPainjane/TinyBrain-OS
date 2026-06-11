# Current Sprint

## Sprint

**Name:** Month 3 — KV and Swap Foundations
**Goal:** KV manager stub, swap manager, brain-top prototype
**Version:** pre-v1.0 memory subsystems (v0.7 scheduler shipped)

## Current Task

[011-kv-manager](../../tasks/011-kv-manager.md)

## In Progress

- [ ] 012-swap-manager (next)

## Blocked

(none)

## Done (prior sprint)

- [x] v0.7 MLFQ scheduler (010) — merged PR #12 @ `a0f90a7`
- [x] v0.6 inference (009a–009d) — tag `v0.6`
- [x] stab-003-tinybrain-cli — PR #10
- [x] repo-hygiene-testing — PR #10

## Next

- 012-swap-manager
- 013-brain-top (prototype)

## Definition Of Done

- KV save/load stub with KVStored/KVLoaded events
- `CGO_ENABLED=0 go test ./...` passes
- Boundary tests pass (`tests/import_boundary_test.go`)
- Task moved to [completed.md](completed.md) with commit hash
- [docs/current.md](../../docs/current.md) synced

## Forbidden Work

- `internal/agents/`
- Kubernetes
- Web dashboard
- API/Router layer
- llama.cpp KV export changes (unless required for stub wiring)

## Active References

- Task: [tasks/011-kv-manager.md](../../tasks/011-kv-manager.md)
- RFC: [docs/rfc/RFC-001-KV-Hibernation.md](../../docs/rfc/RFC-001-KV-Hibernation.md)
- Architecture: [docs/architecture/memory.md](../../docs/architecture/memory.md)

---
**Layer:** planning
**Last updated:** 2026-06-11
