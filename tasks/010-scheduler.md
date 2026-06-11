# Task 010 — Scheduler

## Status

Complete — MLFQ phase (v0.7)

## Goal

Implement FIFO scheduler skeleton with interface ready for MLFQ migration.

### MLFQ phase (v0.7 — active)

Per [docs/specs/v0.7-scheduler.md](../docs/specs/v0.7-scheduler.md):

- [x] MLFQ Q0–Q3 queues (`MLFQQueue`, `MLFQScheduler`)
- [x] Token quantum demotion (`RecordToken`, `TokenQuantum`)
- [x] Preemption of lower-priority processes (`Schedule` + `preemptAndRequeue`)
- [x] Boost/aging anti-starvation (`Boost`, auto-boost every 500 tokens / 30s)
- [x] Swap idle heuristic (`ShouldSwap`, 10s threshold)

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
- `internal/scheduler/mlfq.go`
- `internal/scheduler/scheduler_test.go`

## Acceptance Criteria

- [x] Enqueue adds process to queue
- [x] Schedule dequeues highest-priority (FIFO for v1) process and returns selected process
- [x] Preempt marks process PREEMPTED in process table
- [x] Scheduling policy only — no runtime, loader, or registry imports
- [x] Tests demonstrate two-process queue ordering

## MLFQ Acceptance Criteria

- [x] `MLFQQueue` dequeues Q0 before Q1–Q3
- [x] High-priority enqueue preempts lower-priority runner (`TestMLFQScheduler_Preemption`)
- [x] Token quantum demotes process level (`TestMLFQScheduler_TokenQuantumDemotion`)
- [x] Boost resets queued processes to Q0 (`TestMLFQScheduler_Boost`, `TestMLFQScheduler_AutoBoostViaTokens`)
- [x] Per-level queue depths for telemetry (`QueueDepths`)
- [x] Swap idle heuristic — no swap before 10s idle (`ShouldSwap`)
- [x] No runtime, loader, or inference imports (boundary tests)

## Out Of Scope

- Real runtime load calls (task 011+ integration)
- Swap execution (task 012 swap manager)

## Related

- Spec: [docs/specs/v0.7-scheduler.md](../docs/specs/v0.7-scheduler.md)
- Migration: [planning/architecture-evolution/migration-paths.md](../planning/architecture-evolution/migration-paths.md)

---
**Layer:** task
