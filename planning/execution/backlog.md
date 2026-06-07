# Backlog

Tasks not yet started or scheduled for future months. Active week plan: [roadmap/months/](../roadmap/months/).

| Priority | Task | Target month | Reason | Dependencies |
|----------|------|--------------|--------|--------------|
| MEDIUM | 011-kv-manager | 3 | KV save/load events | v0.4 runtime |
| MEDIUM | 012-swap-manager | 3 | VRAM→RAM tier movement | 011-kv-manager |
| MEDIUM | 013-brain-top | 3 / 6 | Prototype Week 12; polish Week 22 | 010-scheduler |
| MEDIUM | 014-agent-plugin | 4 | v0.8 agent contract | 008-runtime, 005 |
| MEDIUM | tool-registry | TBD | v0.2 spec gap — no task yet | 005, 006 |
| MEDIUM | 015-k8s-crds | 5 | M9 CRD schemas | v0.8 |
| MEDIUM | 016-k8s-controllers-core | 5 | Agent + Task controllers | 015-k8s-crds |
| MEDIUM | 017-k8s-controllers-memory | 5 | KVCache + SwapPolicy | 016-k8s-controllers-core |
| LOW | 018-benchmark-suite | 6 | Swarm vs monolith report | v0.8 agents |

**Completed Month 1:** 001–007  
**Completed Month 2 (in progress):** 008–010

---
**Layer:** planning
**Last updated:** 2026-06-07
