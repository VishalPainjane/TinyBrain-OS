# Month 2 — Runtime and Inference (Week 5–8)

**Status:** Planned
**Month goal:** Runtime abstraction with stub provider, then real GGUF inference — prove M1 and M2.
**Versions targeted:** v0.4, v0.5, v0.6
**Milestones targeted:** M1 (single model runtime), M2 (model switching), M5 (runtime + inference)
**Exit demo:** Load GGUF model, generate one prompt, swap to second model on Standard profile.

## Success criteria (month end)

- [x] v0.4 Runtime shipped (stub InferenceProvider)
- [ ] v0.5 Model Registry shipped (persistent metadata)
- [ ] v0.6 Inference shipped (llama.cpp adapter)
- [ ] M1 and M2 demonstrated with real GGUF Q4_K_M models
- [x] INV-001 and INV-008 verified (scheduler still has zero inference imports)
- [ ] Assumptions updated: mmap on Windows, 4GB VRAM budget

---

## Week 5 — Runtime shell

**Goal:** ModelRuntime + StubProvider; ship v0.4.

**Tasks:** [008-runtime](../../tasks/008-runtime.md)

**Version progress:** v0.4 → 100%

**Milestone:** M5 (partial — shell only)

**Deliverable:** Load/unload/swap two stub models; Generate returns canned JSON.

**Demo:** Test swaps model A → model B via runtime interface without llama.cpp.

**Risks to watch:** Interface churn — lock [contracts/runtime.md](../../docs/contracts/runtime.md) before coding.

**Forbidden this week:** scheduler, llama.cpp CGO, agent plugins

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Week 6 — Model loader and persistence

**Goal:** Loader lifecycle + persistent model registry; ship v0.5.

**Tasks:** [009-model-loader](../../tasks/009-model-loader.md), extend [006-model-registry](../../tasks/006-model-registry.md)

**Version progress:** v0.5 → 100%

**Milestone:** —

**Deliverable:** NOT_LOADED → ACTIVE → UNLOADED lifecycle; model defs survive restart (BoltDB or SQLite — TBD).

**Demo:** Register models, restart process, list still returns definitions.

**Risks to watch:** BoltDB vs SQLite undecided — log decision in accepted.md or mark TBD in spec.

**Forbidden this week:** scheduler, agent plugins

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

---

## Week 7 — llama.cpp integration start

**Goal:** CGO adapter skeleton; mmap load one GGUF (no generate yet).

**Tasks:** extend [009-model-loader](../../tasks/009-model-loader.md)

**Version progress:** v0.6 started (~50%)

**Milestone:** M5 (in progress)

**Deliverable:** LlamaCppProvider struct; one model loads into memory via mmap.

**Demo:** Load single small GGUF; log load time and memory footprint.

**Risks to watch:** CGO toolchain on Windows; CGO build docs in README if needed.

**Forbidden this week:** scheduler, agent plugins

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes)

---

## Week 8 — Real inference

**Goal:** Full generate path; ship v0.6; achieve M1 + M2.

**Tasks:** extend [009-model-loader](../../tasks/009-model-loader.md)

**Version progress:** v0.6 → 100%

**Milestone:** M1, M2, M5

**Deliverable:** Generate one prompt; swap between two GGUF models under resource budget.

**Demo:** Load model A → generate → unload → load model B → generate; log TTFT.

**Risks to watch:** 4GB VRAM insufficient for two models — validate assumption in assumptions.md.

**Forbidden this week:** scheduler (Month 3), agent plugins (Month 4)

**Update when week ends:** [Tier A](../update-checklist.md#tier-a--after-every-task-completes) + [Tier B on ship](../update-checklist.md#tier-b--after-every-version-ships-v0x-complete)

Validate assumptions in [assumptions.md](../../assumptions.md).

---

## Month-end review

**Planned vs actual:** _(fill when Month 2 closes)_

**Carry-forward to Month 3:** FIFO scheduler, MLFQ, KV manager, swap manager, brain-top prototype.

**Update when month ends:** [Tier C](../update-checklist.md#tier-c--after-every-month-closes)

---
**Layer:** planning
**Related:** [../../docs/specs/v0.6-inference.md](../../docs/specs/v0.6-inference.md)
