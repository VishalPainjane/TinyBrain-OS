# Scheduler

The scheduler is the control center of TinyBrain OS. Its primary job is managing **VRAM budget** and **sequence execution order** (continuous batching) — not scheduling agents as abstract concepts, and not swapping entire models.

## Responsibilities

- Queue management (FIFO skeleton → MLFQ at v0.7)
- Admission Control (checking sequence limits before admitting to the active batch)
- Iteration-Level Scheduling (deciding which sequences fill the active batch per forward pass)
- Priority assignment and preemption
- Fairness, aging, starvation prevention
- Swap heuristic decisions (KV Cache eviction, never model eviction)

## MLFQ Design (v0.7 target)

| Queue | Priority | Token quantum |
|-------|----------|---------------|
| Q0 | Highest | 32 tokens |
| Q1 | | 64 tokens |
| Q2 | | 128 tokens |
| Q3 | Lowest | 256 tokens |

**Iteration-Level Scheduling (Continuous/In-Flight Batching):**
The scheduling loop runs at the iteration boundary (after every token generated). Requests do not run to completion uninterrupted. After each forward pass, the scheduler checks the queue, ejects completed sequences, and admits new sequences into the active batch mid-flight to avoid generation stragglers. Quantum exceeded → sequence paused, demote to lower queue. Boost all queues periodically to prevent starvation.

## Preemption

Higher-priority work can interrupt lower-priority active sequences. We frame preemption at the **output-token boundary**, not as a CPU-style process switch. When preempted, the sequence enters the PREEMPTED state, and the context switch is handled by evicting the lower-priority sequence's KV Cache blocks to host RAM (or spilling), allowing the high-priority sequence to allocate VRAM blocks.

## Swap Heuristic

Never swap KV cache immediately on pause. If a sequence is idle > 10 seconds → swap to lower tier (RAM). Otherwise keep warm. The model weights themselves are never swapped as part of preemption. Reduces thrashing.

## Inputs

Task submissions (via router); process table state; resource metrics (VRAM, RAM).

## Outputs

Schedule decisions; preemption commands; runtime load/unload requests; state transition events.

## Dependencies (allowed)

Process table, runtime (via interface), event bus.

## Dependencies (forbidden)

InferenceProvider, llama.cpp, registry writes, UI.

## Future Plans

Full MLFQ with token quanta; Continuous Batching implementation; Chunked Prefill (Prefill/Decode disaggregation) policies; VRAM threshold policies; SwapPolicy CRD integration.

## Non-Goals

Model loading implementation; agent execution logic; inference.

## Related Contracts

[scheduler.md](../contracts/scheduler.md)

## Related ADRs

ADR-003, accepted decision (FIFO first)
ADR-008 (Iteration-level scheduling and paged memory)

---
**Layer:** architecture
**Source:** detail.md Part 4
**Related:** [kernel.md](kernel.md), [runtime.md](runtime.md)
