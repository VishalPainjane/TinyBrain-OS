# Month 3 — Control Plane and Memory (Week 9–12)

**Status:** Planned
**Month goal:** Scheduler (FIFO → MLFQ), KV/swap foundations, brain-top prototype.
**Versions targeted:** v0.7 + memory subsystems (pre-v1.0)
**Milestones targeted:** M6 (scheduler), M8 (brain-top prototype, partial)
**Exit demo:** Preemption visible in process table; brain-top shows live process states; KV save/load emits events.

## Success criteria (month end)

- [ ] v0.7 Scheduler shipped (MLFQ + preemption)
- [ ] KV manager stub operational (task 011)
- [ ] Swap manager moves state VRAM → RAM (task 012)
- [ ] brain-top prototype reads process table (task 013)
- [ ] M6 and M8 (prototype) demonstrated

---

## Week 9 — FIFO scheduler skeleton

**Goal:** Working queue behind Scheduler interface; v0.7 started.

**Tasks:** [010-scheduler](../../tasks/010-scheduler.md) (FIFO phase only)

**Version progress:** v0.7 started (~25%)

**Milestone:** M6 (partial)

**Deliverable:** Enqueue, Schedule (FIFO dequeue), integration with process table.

**Demo:** Two processes enqueued; Schedule runs them in order; test proves ordering.

**Risks to watch:** Scheduler importing runtime concrete types — enforce interface only.

**Forbidden this week:** MLFQ preemption (Week 10), llama.cpp changes

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 10 — MLFQ and preemption

**Goal:** Q0–Q3 queues, token quanta, preemption; ship v0.7.

**Tasks:** [010-scheduler](../../tasks/010-scheduler.md) (MLFQ phase)

**Version progress:** v0.7 → 100%

**Milestone:** M6

**Deliverable:** Higher-priority process preempts lower; queue depths exposed for metrics.

**Demo:** Long-running low-priority task interrupted by high-priority submission.

**Risks to watch:** Swap thrashing — apply 10s idle heuristic from architecture doc.

**Forbidden this week:** agent plugins

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Week 11 — KV manager

**Goal:** Save/load KV stub; KVStored/KVLoaded events.

**Tasks:** [011-kv-manager](../../tasks/011-kv-manager.md)

**Version progress:** — (feeds v1.0)

**Milestone:** —

**Deliverable:** KV block allocation; save/load API shell; events on bus.

**Demo:** Save context ID → KVStored event; load → KVLoaded event in test.

**Risks to watch:** llama.cpp KV export API — may stub metadata only (see RFC-001).

**Forbidden this week:** full compression pipeline

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 12 — Swap manager and brain-top prototype

**Goal:** VRAM → RAM tier move; TUI shows process table.

**Tasks:** [012-swap-manager](../../tasks/012-swap-manager.md), [013-brain-top](../../tasks/013-brain-top.md) (prototype)

**Version progress:** —

**Milestone:** M8 (prototype)

**Deliverable:** Swap manager records tier movement; brain-top renders process states.

**Demo:** Run brain-top alongside scheduler test; see RUNNING/WAITING states update.

**Risks to watch:** TUI coupling to scheduler internals — read via telemetry interface only.

**Forbidden this week:** agent plugins (Month 4)

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Month-end review

**Planned vs actual:** _(fill when Month 3 closes)_

**Carry-forward to Month 4:** Agent plugin contract, sample fleet, v0.8 pipeline.

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../../docs/specs/v0.7-scheduler.md](../../docs/specs/v0.7-scheduler.md), [../../docs/rfc/RFC-001-KV-Hibernation.md](../../docs/rfc/RFC-001-KV-Hibernation.md)
