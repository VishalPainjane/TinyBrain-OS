# Month 1 — Foundation (Week 1–4)

**Status:** Complete
**Month goal:** Ship kernel, event core, registry, and hardware profiling — infrastructure for everything else.
**Versions targeted:** v0.1, v0.2, v0.3
**Milestones targeted:** M3 (partial — events), M4 (hardware profiling)
**Exit demo:** `go test ./...` (process, events, registry, hardware). Boot stdout probe deferred; registry demo via unit tests.

## Success criteria (month end)

- [x] v0.1 Kernel shipped (process table complete)
- [x] v0.2 Registry shipped (in-memory agent + model definitions)
- [x] v0.3 Hardware shipped (profile classification)
- [x] Event bus publishing core lifecycle events
- [x] All Month 1 tasks (002–007) in completed.md
- [x] `go test ./...` green for all implemented packages

---

## Week 1 — Finish kernel

**Goal:** Complete process table; ship v1.1.

**Tasks:** [002-process-table](../../tasks/002-process-table.md)

**Version progress:** v0.1 → 100%

**Milestone:** —

**Deliverable:** Process table with Create, Get, List, UpdateState, Delete; O(1) PID lookup.

**Demo:** `go test ./internal/process/...`

**Risks to watch:** Process table design blocking scheduler (mitigate via [contracts/process.md](../../docs/contracts/process.md)).

**Forbidden this week:** runtime, scheduler, llama.cpp, agents

**Update when week ends:** [Tier A checklist](../update-checklist.md#tier-a--after-every-task-completes)

On v0.1 ship also run: [Tier B](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Week 2 — Event core

**Goal:** Typed events and in-process event bus; start v0.2.

**Tasks:** [003-event-types](../../tasks/003-event-types.md), [004-event-bus](../../tasks/004-event-bus.md)

**Version progress:** v0.2 started

**Milestone:** M3 (partial — decoupled event core)

**Deliverable:** Core event types defined; channel-based pub/sub with tests.

**Demo:** Publish TaskCreated; subscriber receives typed event in test.

**Risks to watch:** Event ordering semantics undocumented — log in accepted.md if needed.

**Forbidden this week:** runtime, scheduler, llama.cpp, agents

**Update when week ends:** [Tier A checklist](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 3 — Registry

**Goal:** In-memory agent and model registry; ship v0.2.

**Tasks:** [005-agent-registry](../../tasks/005-agent-registry.md), [006-model-registry](../../tasks/006-model-registry.md)

**Version progress:** v0.2 → 100%

**Milestone:** M3 (registry half — definitions without fleet demo)

**Deliverable:** Register/list agents and models; no hardcoded agent types in core.

**Demo:** Register two sample definitions; list returns both.

**Risks to watch:** Scope creep into agent execution — registry is data only.

**Forbidden this week:** runtime, scheduler, llama.cpp, agent execution

**Update when week ends:** [Tier A checklist](../update-checklist.md#tier-a--after-every-task-completes)

On v0.2 ship also run: [Tier B](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Week 4 — Hardware probe

**Goal:** Detect and classify hardware profile; ship v0.3.

**Tasks:** [007-hardware-profiler](../../tasks/007-hardware-profiler.md)

**Version progress:** v0.3 → 100%

**Milestone:** M4

**Deliverable:** Boot probe assigns Tiny, Standard, or Workstation profile.

**Demo:** Run probe; stdout shows detected RAM, VRAM (if any), and profile name.

**Risks to watch:** Windows VRAM detection unreliable — mark assumption in assumptions.md if needed.

**Forbidden this week:** runtime, scheduler, llama.cpp, agents (opens in Month 2)

**Update when week ends:** [Tier A checklist](../update-checklist.md#tier-a--after-every-task-completes)

On v0.3 ship also run: [Tier B](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Month-end review

**Planned vs actual:** All four weeks delivered in working tree. Shipped to git as single foundation release at tag `v0.3`. Tool registry and boot stdout probe deferred. Repository governance frozen (7 Cursor rules + CI).

**Carry-forward to Month 2:** Runtime shell (008), model loader (009), persistence (v0.5), inference (v0.6).

**Update when month ends:** [Tier C checklist](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../master-roadmap.md](../master-roadmap.md), [../../docs/specs/v0.1-kernel.md](../../docs/specs/v0.1-kernel.md)
