# Version Progress

Spec-level completion. Update when tasks complete or acceptance criteria are met.

| Version | Progress | Status | Target month | Blocking tasks |
|---------|----------|--------|--------------|----------------|
| V0.1 Kernel | 100% | Shipped | 1 | — |
| V0.2 Registry | 100% | Shipped | 1 | — |
| V0.3 Hardware | 100% | Shipped | 1 | — |
| V0.4 Runtime | 100% | Shipped | 2 | — |
| V0.5 Model Registry | 100% | Shipped | 2 | — |
| V0.6 Inference | 85% | In Progress | 2 | GPU offload (CUDA) |
| V0.7 Scheduler | 0% | Not Started | 3 | MLFQ (010 FIFO skeleton shipped) |
| V0.8 Agents | 0% | Not Started | 4 | 014, integration |
| V1.0 | 0% | Not Started | 6 | full integration |

**Month 1 complete:** v0.1–v0.3 shipped together at git tag `v0.3` (single foundation release). Event core (003–004) delivered alongside registry foundation.

**v0.4 shipped:** Runtime shell (008); stub loader (009) and FIFO scheduler skeleton (010) completed same sprint but outside v0.4 spec scope.

**v0.5 shipped:** tag `v0.5` — persistent model registry (006-registry-persistence).

**v0.6 in progress:** 009a–009c complete (runtime ↔ LlamaProvider wired). Remaining: GPU offload (CUDA).

**Note:** Tool registry (v0.2 spec) deferred — no Month 1 task.

**Monthly plans:** [roadmap/months/](../roadmap/months/)

---
**Layer:** planning
**Related:** [../../docs/specs/](../../docs/specs/), [../releases/](../releases/)
**Last updated:** 2026-06-10
