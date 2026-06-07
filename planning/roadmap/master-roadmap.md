# Master Roadmap

CEO-level view. No technical details, packages, or interfaces — those live in [docs/specs/](../../docs/specs/) and [tasks/](../../tasks/).

## Current Version

**V0.4 Runtime** — Month 1 shipped at tag `v0.3`; Month 2 active.

## Target Version

**V1.0 Integrated Runtime** — Kernel, registry, hardware profiling, runtime, inference, scheduler, plugin agents, and live telemetry working together on consumer hardware.

## Major Milestones

| # | Milestone | Outcome | Status |
|---|-----------|---------|--------|
| M1 | Single model runtime | Prove load/unload one model under resource budget | Planned (Month 2) |
| M2 | Model switching | Dynamic swap between models | Planned (Month 2) |
| M3 | Registry + events | Plugin definitions and decoupled event core | Partial (events + registry shipped) |
| M4 | Hardware profiling | Adaptive behavior per detected profile | **Complete** |
| M5 | Runtime + inference | Real GGUF inference via provider adapter | Planned (Month 2) |
| M6 | Scheduler | Priority queues, preemption, VRAM-aware scheduling | Planned (Month 3) |
| M7 | Plugin agents | User-defined agent fleet via registry | Planned (Month 4) |
| M8 | brain-top TUI | Live runtime visibility (htop for AI) | Planned |
| M9 | Kubernetes operator | CRDs and controllers for Agent, Task, KVCache, SwapPolicy | Planned (Month 5) |
| M10 | Full swarm demo | Multi-agent workflow on Standard profile hardware | Planned (Month 6) |

## Estimated Timeline

| Phase | Focus | Target |
|-------|-------|--------|
| Month 1 | Kernel, events, registry, hardware profiling | **Complete** |
| Month 2 | Runtime shell, inference, model loader | Runtime + M1/M2 |
| Month 3 | Scheduler (MLFQ), KV/swap, brain-top prototype | Control plane |
| Month 4 | Plugin agents, sample fleet demo | Agent layer |
| Month 5 | Kubernetes operator, CRDs | Platform extension |
| Month 6 | Integration, benchmarks, polish | V1.0 release |

## Monthly Plans (detailed)

Week-by-week breakdown (Week 1–24). CEO summary stays here; execution detail lives in month files.

| Month | Weeks | Plan | Theme |
|-------|-------|------|-------|
| 1 | 1–4 | [month-01.md](months/month-01.md) | Foundation (**Complete**) |
| 2 | 5–8 | [month-02.md](months/month-02.md) | Runtime + inference |
| 3 | 9–12 | [month-03.md](months/month-03.md) | Control plane + memory |
| 4 | 13–16 | [month-04.md](months/month-04.md) | Agent layer |
| 5 | 17–20 | [month-05.md](months/month-05.md) | Kubernetes operator |
| 6 | 21–24 | [month-06.md](months/month-06.md) | V1.0 release |

**Update discipline:** [update-checklist.md](update-checklist.md) — Tier A (task), Tier B (version), Tier C (month).

## Alignment notes

| Topic | Detail |
|-------|--------|
| M1 / M2 | Achieved Week 8 (v0.6 inference), not Month 1 |
| M3 | Events + registry shipped Month 1; full pipeline Week 15 |
| M4 | Hardware profiling shipped Month 1 |
| brain-top | Prototype Week 12 (Month 3); production Week 22 (Month 6) |
| K8s operator | Month 5 (M9); v1.0 local demo does not require K8s |

## Dependencies

- Version chain: v0.1 → v0.2 → … → v1.0 (each spec depends on prior)
- External: Go 1.22+, llama.cpp at v0.6, GGUF models for inference demo
- Hardware baseline: 16 GB RAM, RTX 3050 (4 GB VRAM) for Standard profile validation

## Guiding Rule

Every milestone must support: **TinyBrain adapts to hardware, not the reverse.**

---
**Layer:** planning
**Related:** [../metrics/progress.md](../metrics/progress.md), [../../docs/specs/](../../docs/specs/)
