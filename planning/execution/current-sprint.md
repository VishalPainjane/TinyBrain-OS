# Current Sprint

## Sprint

**Name:** Runtime and Inference
**Goal:** Ship v0.6 llama.cpp adapter (Month 2 Weeks 7–8) — **complete (untagged)**
**Version:** V0.6

## Current Task

(none — v0.6 inference scope complete)

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
- [x] 009d-gpu-offload-cuda

## Next

- Cut tag `v0.6` when release notes and manual CUDA sign-off (optional for tag) are accepted
- Month 3 / v0.7 scheduler (MLFQ) — see [master-roadmap.md](../roadmap/master-roadmap.md)

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
- Release: [planning/releases/v0.6.md](../releases/v0.6.md) (untagged)
- Prior release: tag `v0.5` — [planning/releases/v0.5.md](../releases/v0.5.md)

---
**Layer:** planning
**Last updated:** 2026-06-10
