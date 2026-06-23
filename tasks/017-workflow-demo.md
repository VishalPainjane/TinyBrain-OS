# Task 017 — Workflow Demo

## Status

In Progress

## Goal

Create a sequential two-agent workflow demo and ship v0.8.

## Context

Month 4, Week 16. This is the exit demo for v0.8. We need to demonstrate that two agents can run sequentially (e.g. Planner -> Coder) by utilizing the event bus, achieving M7 and M3 milestones.

## Requirements

- Add `Result` to `TaskCompletedPayload`.
- Update `EventListener` to populate `Result`.
- Create `tinybrain workflow` command to orchestrate the two-agent sequence.
- Document peak VRAM and TTFT in assumptions.md.
- Ship v0.8.

## Files

- `internal/events/types.go`
- `internal/agents/listener.go`
- `cmd/tinybrain/workflow.go`
- `cmd/tinybrain/main.go`
- `planning/assumptions.md`

## Acceptance Criteria

- [ ] `TaskCompletedPayload` contains agent's output.
- [ ] `tinybrain workflow` successfully routes a task to `sample-alpha`, parses its JSON, and routes a second task to `sample-beta`.
- [ ] Metrics (VRAM/TTFT) documented in `assumptions.md`.
- [ ] `v0.8` is fully shipped.

## Out Of Scope

- Complex branching workflows (only sequential required).
- Actual complex reasoning logic.

## Related

- Spec: [v0.8-agents.md](../docs/specs/v0.8-agents.md)
- Month plan: [month-04.md](../planning/roadmap/months/month-04.md)

---
**Layer:** task
**Target month:** 4
