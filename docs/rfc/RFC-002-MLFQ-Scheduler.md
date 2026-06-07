# RFC-002: MLFQ Scheduler

## Status

Proposal

## Problem

FIFO scheduling is insufficient for interactive vs background task prioritization. TinyBrain needs OS-grade scheduling adapted for token-generation workloads under VRAM constraints.

## Proposal

Full MLFQ implementation per [research/mlfq-notes.md](../research/mlfq-notes.md):

- Q0–Q3 with token quanta (32/64/128/256)
- Preemption on higher-priority arrival
- Boost/aging every 30s or 500 tokens
- VRAM-aware admission control (reject/defer if insufficient VRAM)

## Alternatives

- Strict priority queue (starvation risk)
- Round-robin (no priority)
- Wall-clock time slices (poor fit for LLM token timing)

## Open Questions

- Optimal quantum values for 2B vs 7B models
- Preemption mid-token vs between tokens
- Interaction with swap heuristic

## Not Scheduled Until

v0.7 spec (FIFO skeleton in task 010 first)

---
**Layer:** planning
**Related:** [../architecture/scheduler.md](../architecture/scheduler.md)
