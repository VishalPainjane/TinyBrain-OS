# Current Sprint

## Sprint

**Name:** Runtime and Inference
**Goal:** Ship v0.6 llama.cpp adapter (Month 2 Weeks 7–8)
**Version:** V0.6

## Current Task

(none — 009c complete; GPU offload next)

## In Progress

(none)

## Blocked

(none)

## Done

- [x] 001-process-state
- [x] 002-process-table
- [x] 003-event-types
- [x] 004-event-bus
- [x] 005-agent-registry
- [x] 006-model-registry
- [x] 006-registry-persistence
- [x] 007-hardware-profiler
- [x] 008-runtime
- [x] 009-model-loader
- [x] 010-scheduler
- [x] 009a-llama-cgo-load
- [x] 009b-cpu-generate
- [x] 009c-runtime-integration

## Next

- **GPU offload (CUDA)** — remaining v0.6 scope
- Month 3: [011-kv-manager](../../tasks/011-kv-manager.md) (after v0.7 MLFQ)

## Definition Of Done

- All acceptance criteria in active task met
- `go test ./...` passes
- Task moved to [completed.md](completed.md) with commit hash
- [docs/current.md](../../docs/current.md) synced

## Forbidden Work

- `internal/agents/`
- Kubernetes
- [011-kv-manager](../../tasks/011-kv-manager.md) and [012-swap-manager](../../tasks/012-swap-manager.md) (Month 3)

## Active References

- Spec: [docs/specs/v0.6-inference.md](../../docs/specs/v0.6-inference.md)
- Contract: [docs/contracts/runtime.md](../../docs/contracts/runtime.md)
- Prior release: tag `v0.5` — [planning/releases/v0.5.md](../releases/v0.5.md)

---
**Layer:** planning
**Last updated:** 2026-06-10
