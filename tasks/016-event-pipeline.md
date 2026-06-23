# Task 016 — Event-Driven Pipeline

## Status

In Progress

## Goal

Wire task routing, scheduling, and execution entirely via events without direct component-to-component lifecycle calls.

## Context

v0.8 and M3. Month 4 Week 15. The system must orchestrate work via the event bus to decouple the router, scheduler, and agent execution layers.

## Requirements

- Listen for `TaskCreated` and spawn processes via a router.
- `Scheduler` listens for `ProcessSpawned` to enqueue work, and publishes `ProcessStateChanged`.
- `Agent Executor` listens for `ProcessStateChanged` (Running) to start tasks.
- Command line test trace matches `TaskCreated → ProcessSpawned → AgentStarted → TaskCompleted`.

## Files

- `internal/router/router.go`
- `internal/scheduler/scheduler.go`
- `internal/agents/listener.go`
- `cmd/tinybrain/run.go`
- `tests/pipeline_integration_test.go`

## Acceptance Criteria

- [ ] `TaskCreated` successfully transitions into a spawned process.
- [ ] Process transitions to Running state entirely via scheduler event loops.
- [ ] Agent executes upon Running state via event subscription.
- [ ] End-to-end trace verifies no direct method calls between core lifecycle components.
- [ ] `go test ./...` passes.

## Out Of Scope

- Modifying inference or registry directly
- Kubernetes operator

## Related

- Spec: [v0.8-agents.md](../docs/specs/v0.8-agents.md)
- Month plan: [month-04.md](../planning/roadmap/months/month-04.md)

---
**Layer:** task
**Target month:** 4
