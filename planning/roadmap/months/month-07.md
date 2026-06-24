# Month 7 — Advanced Subsystems (Week 25–28)

**Status:** Planned
**Month goal:** Implement KV Cache Compression pipeline, integrate cloud fallback providers, and deploy the full Kubernetes fleet operator in a production environment.
**Versions targeted:** v1.1
**Milestones targeted:** Post-v1.0 extensions
**Exit demo:** TinyBrain dynamically offloads inactive KV memory to NVMe using Zstandard compression, gracefully fails over to an OpenAI-compatible cloud endpoint when local VRAM is exhausted, and runs seamlessly across a multi-node Kubernetes cluster.

## Success criteria (month end)

- [ ] KV cache blocks are compressed before swapping to RAM/NVMe
- [ ] `CloudProvider` implements `InferenceProvider` for OpenAI/Anthropic
- [ ] Kubernetes memory controllers (KVCache, SwapPolicy) manage distributed state
- [ ] Local tests pass with >80% coverage on new subsystems

---

## Week 25 — KV Compression Pipeline

**Goal:** Protect VRAM exhaustion through transparent cache compression.

**Tasks:** [024-kv-compression](../../tasks/024-kv-compression.md)

**Version progress:** v1.1 started

**Deliverable:** Zstandard compression integrated into the `internal/swap` and `internal/kv` lifecycle.

**Risks to watch:** CPU overhead during compression blocking the scheduler.

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 26 — Cloud Provider Fallback

**Goal:** Standardized cloud bridging when local resources fail.

**Tasks:** [025-cloud-providers](../../tasks/025-cloud-providers.md)

**Deliverable:** `OpenAIProvider` adapting the `InferenceProvider` interface.

**Risks to watch:** Leaking cloud dependencies into the local-first scheduler core.

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 27 — Kubernetes Production Deploy

**Goal:** Shift the `cmd/operator` from local-dev to production-ready manifests.

**Tasks:** [026-k8s-production](../../tasks/026-k8s-production.md)

**Deliverable:** Helm charts or production Kustomize overlays supporting RBAC and distributed agent execution.

**Risks to watch:** Network routing complexities across Pods.

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 28 — v1.1 Release & Month 7 Wrap

**Goal:** Stabilize Advanced Subsystems and ship v1.1.

**Tasks:** Release polish, e2e testing.

**Version progress:** v1.1 → 100%

**Deliverable:** v1.1 release notes.

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---

## Month-end review

**Planned vs actual:** _(fill when Month 7 closes)_

**Post-v1.1:** TBD

---
**Layer:** planning
