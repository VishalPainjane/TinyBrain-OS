# Accepted Decisions

Tactical decisions already made. Promote to [docs/adr/](../../docs/adr/) when permanent architecture is confirmed.

| Decision | Reason | Date | Impact |
|----------|--------|------|--------|
| Agents are plugins, not hardcoded classes | Evolved architecture — user-defined fleets | 2026-06-07 | Core never imports agent types |
| Go as core language | Kubernetes ecosystem, controller-runtime, concurrency | 2026-06-07 | All kernel code in Go |
| Event bus v1 uses Go channels | MVP simplicity; few hundred events/sec | 2026-06-07 | No NATS/Kafka until scale requires |
| Scheduler skeleton starts FIFO | Working queue before MLFQ complexity | 2026-06-07 | MLFQ deferred to v0.7 spec |
| llama.cpp via InferenceProvider adapter | Hexagonal boundary (ADR-004) | 2026-06-07 | Scheduler never touches engine |
| Code-first truth in all docs | Documentation reflects implemented state | 2026-06-07 | No planned-state masquerading as done |
| Process owns state; scheduler owns transitions | Clear ownership split | 2026-06-07 | State in kernel; transitions in scheduler later |
| GGUF Q4_K_M as default quantization | Balance of quality and memory on consumer hardware | 2026-06-07 | Model registry defaults |
| Event bus v1 has no ordering guarantee | Async goroutine dispatch per subscriber; MVP scale | 2026-06-07 | Publishers must not assume FIFO delivery |
| Repository governance frozen in .cursor/rules/ | AGENTS.md + 7 rule files + CI as institutional memory | 2026-06-07 | All agents follow github-workflow.mdc |
| EventBus unsubscribe via returned func | Equivalent to Unsubscribe; avoids subscription ID registry in v1 | 2026-06-07 | `Subscribe(eventType, handler) (unsubscribe func())` |
| Month 1 git history as single foundation release | Honest audit trail; no fabricated per-task commits | 2026-06-07 | Tag `v0.3` represents v0.1+v0.2+v0.3 |

---
**Layer:** planning
**Related:** [rejected.md](rejected.md), [../assumptions.md](../assumptions.md)
