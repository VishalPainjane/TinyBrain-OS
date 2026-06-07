# Task 010 — Schedule Signature Contract Drift

**Status:** Open (deferred reconciliation)  
**Date:** 2026-06-07  
**Task:** [tasks/010-scheduler.md](../../tasks/010-scheduler.md)

## Discrepancy

| Layer | `Schedule` signature |
|-------|----------------------|
| [docs/contracts/scheduler.md](../../docs/contracts/scheduler.md) | `Schedule() error` |
| Task 010 implementation | `Schedule() (process.Process, error)` |

## Rationale

Task 010 implements **scheduling policy only** (queue, selection, state transitions). Runtime orchestration is deferred to a future integration task. Returning the selected process makes the scheduler usable without coupling to runtime.

## Resolution plan

Reconcile in a **future integration task** that wires scheduler output to runtime load/unload. Options:

1. Update contract to `Schedule() (Process, error)` and document runtime handoff.
2. Keep `Schedule() error` and add `SelectedProcess()` or event-based handoff.

Contract file is **not modified** in task 010.

---
**Layer:** planning
