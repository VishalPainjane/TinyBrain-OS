# Month 6 — V1.0 Release (Week 21–24)

**Status:** Planned
**Month goal:** Integrate all subsystems, polish brain-top, run benchmarks, ship v1.0 and M10.
**Versions targeted:** v1.0
**Milestones targeted:** M10 (full swarm demo)
**Exit demo:** `brain-top` live during multi-agent task on Standard profile; README quick start works; benchmark report published.

## Success criteria (month end)

- [ ] v1.0 Integrated Runtime shipped
- [ ] All v0.1–v0.8 acceptance criteria verified
- [ ] brain-top production panels: agents, queues, VRAM/RAM, swaps
- [ ] Benchmark: swarm vs monolith (TTFT, TPS, peak VRAM)
- [ ] M10 demonstrated
- [ ] README quick start commands accurate

---

## Week 21 — System integration

**Goal:** Wire kernel, registry, hardware, runtime, inference, scheduler, agents, telemetry.

**Tasks:** integration of all prior tasks

**Version progress:** v1.0 started (~25%)

**Milestone:** —

**Deliverable:** Single binary path runs full pipeline; forbidden list cleared for core packages.

**Demo:** End-to-end task without manual component invocation.

**Risks to watch:** Integration bugs at boundaries — run boundary tests per package.

**Forbidden this week:** New features — fix only

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 22 — brain-top production

**Goal:** Full TUI panels per architecture/telemetry.md.

**Tasks:** [013-brain-top](../../tasks/013-brain-top.md) (production polish)

**Version progress:** v1.0 ~50%

**Milestone:** M8 (complete)

**Deliverable:** Four panels: process states, resource bars, queue depths, swap monitor.

**Demo:** Run brain-top during two-agent workload; all panels update live.

**Risks to watch:** Terminal performance on Windows — test Bubble Tea on target OS.

**Forbidden this week:** K8s-only features in demo path

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 23 — Benchmark suite

**Goal:** Swarm vs monolith comparison report.

**Tasks:** [018-benchmark-suite](../../tasks/018-benchmark-suite.md)

**Version progress:** v1.0 ~75%

**Milestone:** —

**Deliverable:** Repeatable benchmark scripts; results vs [target-metrics.md](../../benchmarks/target-metrics.md).

**Demo:** Run benchmark; output TTFT, TPS, peak VRAM for swarm and monolith configs.

**Risks to watch:** Non-comparable hardware profiles — document profile in every result.

**Forbidden this week:** —

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 24 — Ship v1.0 and M10

**Goal:** Release integrated runtime; full swarm demo.

**Tasks:** release polish, docs sync

**Version progress:** v1.0 → 100%

**Milestone:** M10

**Deliverable:** v1.0 release notes in planning/releases/v1.0.md; README updated; all milestones M1–M8 verified, M9 optional parallel, M10 complete.

**Demo:** Standard profile: multi-agent workflow + brain-top + metrics logged.

**Risks to watch:** Last-minute scope creep — freeze features Monday of Week 24.

**Forbidden this week:** New ADRs without ship delay approval

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---

## Month-end review

**Planned vs actual:** _(fill when Month 6 closes)_

**Post-v1.0:** KV compression pipeline, cloud providers, full K8s production deploy — track in backlog and RFCs.

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../../docs/specs/v1.0.md](../../docs/specs/v1.0.md), [../releases/v1.0.md](../releases/v1.0.md)
