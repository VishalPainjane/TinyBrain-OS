# Architecture — Migration Paths

How to evolve MVP code without throwaway rewrites. Each path documents today → future → steps.

## Path 1: FIFO → MLFQ Scheduler

**Today:** Single queue, FIFO dequeue (task 010 skeleton)
**Future:** MLFQ Q0–Q3 with token quanta, preemption, boost (v0.7)

| Step | Action |
|------|--------|
| 1 | Introduce `Queue` interface behind scheduler |
| 2 | Implement single FIFO queue using interface |
| 3 | Add priority levels (multiple queues) without preemption |
| 4 | Enable preemption + token quantum demotion |
| 5 | Add boost/aging anti-starvation every N tokens or 30s |

**Preserves:** Enqueue/Schedule/Preempt contract in [contracts/scheduler.md](../../docs/contracts/scheduler.md)

---

## Path 2: Stub Runtime → llama.cpp

**Today:** `ModelRuntime` interface with no-op implementation (v0.4)
**Future:** llama.cpp GGUF adapter (v0.6)

| Step | Action |
|------|--------|
| 1 | Define `InferenceProvider` interface in contract |
| 2 | Stub provider returns canned responses |
| 3 | llama.cpp adapter implements same interface |
| 4 | Swap provider via config — zero scheduler changes |

**Preserves:** Runtime contract; INV-008

---

## Path 3: In-Memory Registry → Persistent

**Today:** `ModelStore` with `InMemoryStore` and `BboltStore` (v0.5 task 006-registry-persistence)
**Future:** agent/tool persistence (deferred)

| Step | Action |
|------|--------|
| 1 | Registry behind interface |
| 2 | In-memory implementation |
| 3 | Add persistence adapter without API change |

**Preserves:** Registry contract; callers unchanged

---

## Path 4: Go Channels → NATS

**Today:** `chan Event` in-process bus (v0.4)
**Future:** NATS for distributed mode (post-v1.0 if needed)

| Step | Action |
|------|--------|
| 1 | Event bus behind interface |
| 2 | Channel implementation |
| 3 | NATS adapter for multi-node deployment |

**Preserves:** Event type definitions; publisher/subscriber API

---
**Layer:** planning
**Related:** [../../docs/architecture/scheduler.md](../../docs/architecture/scheduler.md), [accepted.md](../decisions/accepted.md)
