# Task 022 — System Integration

**Status:** Complete
**Owner:** System Agent
**Layer:** tasks

## Context
Month 6, Week 21 marks the integration of all prior tasks into a single v1.0 pipeline. The kernel, registry, hardware, runtime, inference, scheduler, and agents must communicate seamlessly to execute an end-to-end workload without manual component invocation.

## Requirements

1. Single binary path (`cmd/tinybrain run` or `workflow`) runs the full pipeline.
2. The core pipeline relies on `EventBus` to achieve loose coupling between the `Scheduler` and the `Runtime`. Explicit command dispatch is discouraged in favor of process state change events triggering the `Agent Executor`.
3. Forbidden list cleared for core packages (i.e. no K8s imports in the core).

## Implementation Plan

1. Verify `cmd/tinybrain/workflow.go` effectively wires the `ptab`, `sched`, `bus`, `exec`, `listener`, and `runtime`.
2. Update `planning/architecture-evolution/current-state.md` to reflect that the "Scheduler -> runtime command wiring" is completely implemented via the `EventListener` bridging `TypeProcessStateChanged` events.
3. Validate through unit and boundary tests that no inference/K8s packages leak into the core scheduler/runtime boundaries.

## Acceptance Criteria

- [ ] `go run cmd/tinybrain/main.go workflow --prompt "test"` succeeds end-to-end.
- [ ] `go test ./...` passes (with `CGO_ENABLED=0`).
- [ ] `tests/import_boundary_test.go` passes.
- [ ] Documentation updated to reflect completion of wiring.
