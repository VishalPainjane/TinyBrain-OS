# Task 010 — Scheduler

## Status

Complete

## Goal

Implement FIFO scheduler skeleton with interface ready for MLFQ migration.

## Context

Accepted decision: start FIFO, migrate to MLFQ at v0.7. Scheduler schedules VRAM/processes — never calls inference directly.

## Requirements

- Scheduler interface per [docs/contracts/scheduler.md](../docs/contracts/scheduler.md)
- Queue interface with FIFO implementation
- Enqueue, Schedule (dequeue + request runtime load), Preempt (shell), Boost (no-op until MLFQ)
- Integration with process table (read/update state)

## Files

- `internal/scheduler/scheduler.go`
- `internal/scheduler/queue.go`
- `internal/scheduler/scheduler_test.go`

## Acceptance Criteria

- [x] Enqueue adds process to queue
- [x] Schedule dequeues highest-priority (FIFO for v1) process and returns selected process
- [x] Preempt marks process PREEMPTED in process table
- [x] Scheduling policy only — no runtime, loader, or registry imports
- [x] Tests demonstrate two-process queue ordering

## Out Of Scope

- MLFQ Q0–Q3 (v0.7)
- Token quantum, boost/aging
- Real runtime load calls (mock runtime in tests)

## Related

- Spec: [docs/specs/v0.7-scheduler.md](../docs/specs/v0.7-scheduler.md)
- Migration: [planning/architecture-evolution/migration-paths.md](../planning/architecture-evolution/migration-paths.md)

---
**Layer:** task
