# Month 4 — Agent Layer (Week 13–16)

**Status:** Complete
**Month goal:** Plugin agents and event-driven multi-agent pipeline; ship v0.8.
**Versions targeted:** v0.8
**Milestones targeted:** M7 (plugin agents), M3 (full end-to-end registry + events)
**Exit demo:** Two sample agents from registry config execute a task via events; peak VRAM logged on Standard profile.

## Success criteria (month end)

- [x] v0.8 Agents shipped
- [x] Agent plugin contract implemented (no hardcoded types in core)
- [x] Sample fleet YAML drives two-agent workflow
- [x] M7 and M3 demonstrated
- [x] Peak VRAM for 2-agent demo recorded in benchmarks notes

---

## Week 13 — Agent plugin contract

**Goal:** Generic agent interface + one sample plugin; v0.8 started.

**Tasks:** [014-agent-plugin](../../tasks/014-agent-plugin.md)

**Version progress:** v0.8 started (~25%)

**Milestone:** M7 (partial)

**Deliverable:** Agent executes via runtime API; returns structured JSON.

**Demo:** Single sample plugin runs one task in test.

**Risks to watch:** Agent importing inference — enforce INV-003.

**Forbidden this week:** hardcoded Planner/Coder types in internal/

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 14 — Sample fleet config

**Goal:** YAML fleet with two sample agents (labels only, not Go types).

**Tasks:** extend [005-agent-registry](../../tasks/005-agent-registry.md)

**Version progress:** v0.8 ~50%

**Milestone:** M7 (partial)

**Deliverable:** Load fleet from config; registry resolves agent → model → tools.

**Demo:** Load sample fleet YAML; list shows two agent definitions.

**Risks to watch:** Config schema drift — document in contracts/registry.md if fields change.

**Forbidden this week:** Kubernetes, cloud providers

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 15 — Event-driven pipeline

**Goal:** Wire task → route → schedule → execute entirely via events.

**Tasks:** integration across 003–010, 014

**Version progress:** v0.8 ~75%

**Milestone:** M3 (full)

**Deliverable:** End-to-end flow with no direct component-to-component lifecycle calls.

**Demo:** Submit task; trace events TaskCreated → ProcessSpawned → AgentStarted → TaskCompleted.

**Risks to watch:** Debugging difficulty — ensure telemetry logs event chain.

**Forbidden this week:** K8s operator work

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 16 — Ship v0.8 and M7

**Goal:** Two-agent workflow demo; ship v0.8.

**Tasks:** integration polish

**Version progress:** v0.8 → 100%

**Milestone:** M7, M3

**Deliverable:** Two sample agents run sequential workflow; metrics capture peak VRAM.

**Demo:** Full pipeline on Standard profile hardware; document TTFT and VRAM in assumptions/benchmarks.

**Risks to watch:** 4GB VRAM assumption — validate or invalidate in assumptions.md.

**Forbidden this week:** —

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Month-end review

**Planned vs actual:** All deliverables for Month 4 completed successfully! The event-driven architecture successfully orchestrates agents dynamically loaded from a YAML fleet configuration. Peak VRAM requirements for a 2-agent workflow proved to be well under the 4GB Standard Profile limit (~3.1GB observed).

**Carry-forward to Month 5:** Kubernetes CRDs and controllers.

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../../docs/specs/v0.8-agents.md](../../docs/specs/v0.8-agents.md)
