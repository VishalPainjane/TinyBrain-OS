# Task 013 — brain-top TUI

## Status

Complete (Month 3 prototype)

## Goal

Terminal dashboard showing live process states, resources, queues, and swap activity.

## Context

Signature visibility feature (M8). Prototype Month 3 Week 12; production polish Month 6 Week 22.

## Requirements

- `cmd/brain-top/` read-only dashboard (stdlib only, no Bubble Tea yet)
- Process panel from `ProcessTable` via adapter
- MLFQ queue depth panel via `MLFQScheduler.QueueDepths`
- Resource panel from `hardware.ProbeAndClassify`
- `snapshot` and `watch` modes; no scheduler/runtime mutation

## Files

- `cmd/brain-top/main.go`
- `cmd/brain-top/render.go`
- `cmd/brain-top/render_test.go`

## Acceptance Criteria

- [x] Renders process STATE column from process table (test with RUNNING/WAITING)
- [x] Renders Q0–Q3 queue depths
- [x] Renders hardware resource summary
- [x] Read-only — no scheduler mutation from TUI
- [x] `go test ./cmd/brain-top/...` passes

## Out Of Scope

- Web dashboard
- Direct scheduler mutation from TUI
- Bubble Tea interactive UI (Month 6)
- Swap activity panel (no swap history store yet)
- Live kernel IPC attachment

## Related

- Month plan: [month-03.md](../planning/roadmap/months/month-03.md), [month-06.md](../planning/roadmap/months/month-06.md)
- Architecture: [telemetry.md](../docs/architecture/telemetry.md)

---
**Layer:** task
**Target month:** 3 (prototype), 6 (production)
