# Current Sprint

## Sprint

**Name:** Runtime Foundation
**Goal:** Establish runtime shell and model loader primitives (Month 2)
**Version:** V0.4

## Current Task

008-runtime

## In Progress

- [ ] 008-runtime

## Blocked

(none)

## Done

- [x] 001-process-state
- [x] 002-process-table
- [x] 003-event-types
- [x] 004-event-bus
- [x] 005-agent-registry
- [x] 006-model-registry
- [x] 007-hardware-profiler

## Next

- 009-model-loader

## Definition Of Done

- All acceptance criteria in active task met
- `go test ./...` passes
- Task moved to [completed.md](completed.md) with commit hash
- [docs/current.md](../../docs/current.md) synced

## Forbidden Work

- `internal/scheduler/`
- llama.cpp / inference bindings (until task 009+)
- `internal/agents/`
- Kubernetes

## Active References

- Spec: [docs/specs/v0.4-runtime.md](../../docs/specs/v0.4-runtime.md)
- Contract: [docs/contracts/runtime.md](../../docs/contracts/runtime.md)
- Month 1 release: tag `v0.3`

---
**Layer:** planning
**Last updated:** 2026-06-07
