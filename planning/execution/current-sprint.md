# Current Sprint

## Sprint

**Name:** V0.7 MLFQ Scheduler
**Goal:** MLFQ queues, token quantum, preemption, boost/aging — ship v0.7 scheduler
**Version:** v0.7

## Current Task

[010-scheduler](../../tasks/010-scheduler.md) — MLFQ phase

## In Progress

(none — v0.7 scheduler core complete; tag pending)

## Blocked

(none)

## Done (prior sprint)

- [x] v0.6 inference (009a–009d) — tag `v0.6`
- [x] stab-003-tinybrain-cli — PR #10
- [x] repo-hygiene-testing — PR #10

## Next

- 011-kv-manager (after v0.7 scheduler core)
- 012-swap-manager
- brain-top prototype (013) — after scheduler visibility

## Definition Of Done

- All acceptance criteria in [v0.7-scheduler.md](../../docs/specs/v0.7-scheduler.md) met
- `CGO_ENABLED=0 go test ./...` passes
- Boundary tests pass (`tests/import_boundary_test.go`)
- Task moved to [completed.md](completed.md) with commit hash
- [docs/current.md](../../docs/current.md) synced

## Forbidden Work

- `internal/agents/`
- Kubernetes
- Web dashboard
- API/Router layer
- llama.cpp / inference adapter changes (unless bugfix)

## Active References

- Spec: [docs/specs/v0.7-scheduler.md](../../docs/specs/v0.7-scheduler.md)
- Contract: [docs/contracts/scheduler.md](../../docs/contracts/scheduler.md)
- Task: [tasks/010-scheduler.md](../../tasks/010-scheduler.md)

---
**Layer:** planning
**Last updated:** 2026-06-11
