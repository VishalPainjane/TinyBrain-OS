# Version Progress

Spec-level completion. Update when tasks complete or acceptance criteria are met.

| Version | Progress | Status | Target month | Blocking tasks |
|---------|----------|--------|--------------|----------------|
| V0.1 Kernel | 100% | Shipped | 1 | — |
| V0.2 Registry | 100% | Shipped | 1 | — |
| V0.3 Hardware | 100% | Shipped | 1 | — |
| V0.4 Runtime | 100% | Shipped | 2 | — |
| V0.5 Model Registry | 100% | Shipped | 2 | — |
| V0.6 Inference | 100% | Shipped | 2 | — |
| V0.7 Scheduler | 100% | Shipped | 3 | — (tag `v0.7` @ `a0f90a7`) |
| V0.8 Agents | 100% | Shipped | 4 | — (tag `v0.8` @ `a40b0fd`) |
| V1.0 | 0% | Not Started | 6 | full integration |

**Month 1 complete:** v0.1–v0.3 shipped together at git tag `v0.3` (single foundation release). Event core (003–004) delivered alongside registry foundation.

**v0.4 shipped:** Runtime shell (008); stub loader (009) and FIFO scheduler skeleton (010) completed same sprint but outside v0.4 spec scope.

**v0.5 shipped:** tag `v0.5` — persistent model registry (006-registry-persistence).

**v0.6 shipped:** tag `v0.6` — llama.cpp inference (009a–009d). CUDA matrix status promoted to Full upon manual GPU checklist verification and sign-off on Windows.

**v0.7 shipped:** tag `v0.7` — MLFQ scheduler (010). Month 3 memory foundations (011–013) on `main` post-tag.

**v0.8 shipped:** tag `v0.8` — Agent plugin contract (014), fleet registry (015), event-driven pipeline (016), sequential agent workflow (017).

**Note:** Tool registry (v0.2 spec) deferred — no Month 1 task.

**Monthly plans:** [roadmap/months/](../roadmap/months/)

**CI health (STAB-002):** merge-blocking job timing and cache history in artifact `ci-run-record-{run_id}`; baselines in [ci-baseline.md](ci-baseline.md).

---
**Layer:** planning
**Related:** [../../docs/specs/](../../docs/specs/), [../releases/](../releases/)
**Last updated:** 2026-06-11
