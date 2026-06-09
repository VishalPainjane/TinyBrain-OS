# Backlog

Tasks not yet started or scheduled for future months. Active week plan: [roadmap/months/](../roadmap/months/). **Active sprint task is not listed here** — see [current-sprint.md](current-sprint.md).

| Priority | Task | Target month | Reason | Dependencies |
|----------|------|--------------|--------|--------------|
| HIGH | 009c — runtime/loader integration | 2 | M1/M2 demo; lifecycle orchestration | 009b (complete) |
| MEDIUM | 011-kv-manager | 3 | KV save/load events | v0.7 MLFQ in progress |
| MEDIUM | 012-swap-manager | 3 | VRAM→RAM tier movement | 011-kv-manager |
| MEDIUM | 013-brain-top | 3 / 6 | Prototype Week 12; polish Week 22 | 010-scheduler |
| MEDIUM | 014-agent-plugin | 4 | v0.8 agent contract | 008-runtime, 005 |
| MEDIUM | tool-registry | TBD | v0.2 spec gap — no task yet | 005, 006 |
| MEDIUM | 015-k8s-crds | 5 | M9 CRD schemas | v0.8 |
| MEDIUM | 016-k8s-controllers-core | 5 | Agent + Task controllers | 015-k8s-crds |
| MEDIUM | 017-k8s-controllers-memory | 5 | KVCache + SwapPolicy | 016-k8s-controllers-core |
| LOW | 018-benchmark-suite | 6 | Swarm vs monolith report | v0.8 agents |

**Completed Month 1:** 001–007  
**Completed Month 2 (adjacent):** 008–010, 006-registry-persistence, 009a-llama-cgo-load, 009b-cpu-generate  
**Next (sprint):** 009c — see [current-sprint.md](current-sprint.md)

**Note:** 011-kv-manager was incorrectly listed as V0.5 active task; realigned 2026-06-08 per planning assessment.

---
**Layer:** planning
**Last updated:** 2026-06-09
