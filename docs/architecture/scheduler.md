# Scheduler

The scheduler is the control center of TinyBrain OS. Its primary job is managing **VRAM budget** and process execution order — not scheduling agents as abstract concepts.

## Responsibilities

- Queue management (FIFO skeleton → MLFQ at v0.7)
- Priority assignment and preemption
- Fairness, aging, starvation prevention
- Delegating model operations to runtime (never calling inference directly)
- Swap heuristic decisions (idle > 10s before swap)

## MLFQ Design (v0.7 target)

| Queue | Priority | Token quantum |
|-------|----------|---------------|
| Q0 | Highest | 32 tokens |
| Q1 | | 64 tokens |
| Q2 | | 128 tokens |
| Q3 | Lowest | 256 tokens |

Scheduling loop runs per token generated. Quantum exceeded → demote to lower queue. Boost all queues periodically (every 30s or 500 tokens) to prevent starvation.

## Preemption

Higher-priority work can interrupt lower-priority running processes. Preempted process → PREEMPTED state; KV may be preserved.

## Swap Heuristic

Never swap immediately on pause. If idle > 10 seconds → swap to lower tier. Otherwise keep warm. Reduces thrashing.

## Inputs

Task submissions (via router); process table state; resource metrics (VRAM, RAM).

## Outputs

Schedule decisions; preemption commands; runtime load/unload requests; state transition events.

## Dependencies (allowed)

Process table, runtime (via interface), event bus.

## Dependencies (forbidden)

InferenceProvider, llama.cpp, registry writes, UI.

## Future Plans

Full MLFQ with token quanta; VRAM threshold policies; SwapPolicy CRD integration (Kubernetes).

## Non-Goals

Model loading implementation; agent execution logic; inference.

## Related Contracts

[scheduler.md](../contracts/scheduler.md)

## Related ADRs

ADR-003, accepted decision (FIFO first)

---
**Layer:** architecture
**Source:** detail.md Part 4
**Related:** [kernel.md](kernel.md), [runtime.md](runtime.md)
