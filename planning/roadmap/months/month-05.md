# Month 5 — Platform Extension (Week 17–20)

**Status:** Complete
**Month goal:** Kubernetes operator with CRDs and reconciliation controllers; achieve M9.
**Versions targeted:** Post-v0.8 platform (not a semver bump — feeds v1.0 ecosystem)
**Milestones targeted:** M9 (Kubernetes operator)
**Exit demo:** `kubectl apply` Agent CR triggers model load via controller reconcile loop.

## Success criteria (month end)

- [x] CRDs defined: Agent, Task, KVCache, SwapPolicy
- [x] Agent + Task controllers reconcile
- [x] KVCache + SwapPolicy controllers reconcile
- [x] M9 demonstrated on local single-node cluster (kind/minikube)
- [x] Documented: v1.0 demo does not require K8s (operator is parallel track)

---

## Week 17 — CRD design

**Goal:** YAML schemas for four CRDs.

**Tasks:** [015-k8s-crds](../../tasks/015-k8s-crds.md)

**Version progress:** —

**Milestone:** M9 (partial)

**Deliverable:** CRD manifests + validation schema; align with [RFC-003](../../docs/rfc/RFC-003-Kubernetes-Operator.md).

**Demo:** `kubectl apply --dry-run=client` succeeds for sample Agent CR.

**Risks to watch:** Over-engineering CRD count — stick to four CRDs only.

**Forbidden this week:** Rewriting core runtime for K8s

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 18 — Agent and Task controllers

**Goal:** Reconcile loop skeleton for Agent and Task resources.

**Tasks:** [016-k8s-controllers-core](../../tasks/016-k8s-controllers-core.md)

**Version progress:** —

**Milestone:** M9 (partial)

**Deliverable:** Controller watches Agent CR → ensures model loaded state.

**Demo:** Apply Agent CR; controller logs reconcile; runtime LoadModel called via adapter.

**Risks to watch:** controller-runtime learning curve — time-box spike.

**Forbidden this week:** Changing scheduler contracts for K8s

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 19 — KVCache and SwapPolicy controllers

**Goal:** Declarative swap thresholds and KV tracking.

**Tasks:** [017-k8s-controllers-memory](../../tasks/017-k8s-controllers-memory.md)

**Version progress:** —

**Milestone:** M9 (partial)

**Deliverable:** SwapPolicy CR changes scheduler threshold; KVCache CR tracks block metadata.

**Demo:** Apply SwapPolicy with vramThreshold 80%; scheduler respects config.

**Risks to watch:** Coupling operator to MLFQ internals — use config adapter boundary.

**Forbidden this week:** —

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 20 — Operator end-to-end

**Goal:** Full M9 demo on local cluster.

**Tasks:** integration

**Version progress:** —

**Milestone:** M9

**Deliverable:** Helm or manifest bundle; README section for K8s deploy mode.

**Demo:** Deploy operator → apply Agent + Task CRs → workload runs.

**Risks to watch:** Scope creep into multi-node — single-node only for M9.

**Forbidden this week:** v1.0 integration (Month 6)

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

Log platform decisions in [accepted.md](../../decisions/accepted.md).

---

## Month-end review

**Planned vs actual:** Month 5 successfully concluded. We built the Kubernetes API types and Custom Resource Definitions (`Agent`, `Task`, `KVCache`, `SwapPolicy`), implemented the `controller-runtime` reconciler logic in Go, successfully completely isolated the logic to the new `cmd/operator` entrypoint to protect the core `internal/runtime`, and finally bundled the solution using Kustomize manifests. M9 is complete.

**Carry-forward to Month 6:** v1.0 integration, brain-top polish, benchmarks.

**Note:** K8s operator ships in Month 5 but v1.0 local demo (Month 6) remains non-K8s per [v1.0 spec](../../docs/specs/v1.0.md).

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../../docs/rfc/RFC-003-Kubernetes-Operator.md](../../docs/rfc/RFC-003-Kubernetes-Operator.md)
