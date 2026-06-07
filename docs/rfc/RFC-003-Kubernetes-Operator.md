# RFC-003: Kubernetes Operator

## Status

Proposal

## Problem

TinyBrain needs enterprise-grade deployment story and declarative resource management for portfolio/recruiting credibility. Running as a single binary limits platform integration.

## Proposal

TinyBrain Kubernetes Operator with CRDs:

| CRD | Purpose |
|-----|---------|
| Agent | Declarative agent plugin config |
| Task | Work unit with assigned agent |
| KVCache | Track KV location, size, owner |
| SwapPolicy | Declarative VRAM/RAM thresholds for scheduler |

Controllers reconcile desired state: Agent Running → model loaded, resources assigned.

## Alternatives

- Docker Compose only (simpler, less impressive)
- Helm charts without operator (no reconciliation loop)
- 20+ CRDs (over-engineered — rejected)

## Open Questions

- Operator SDK vs controller-runtime directly
- Single-node vs multi-node cluster assumptions
- CRD versioning strategy

## Not Scheduled Until

Post-v1.0 (Month 5 in master roadmap)

---
**Layer:** planning
**Related:** [../../planning/roadmap/master-roadmap.md](../../planning/roadmap/master-roadmap.md)
