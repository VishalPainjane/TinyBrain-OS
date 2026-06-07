# RFC-005: Event Ordering

## Status

Proposal

## Problem

Event bus v1 (`ChannelBus`) dispatches each publish to subscribers via goroutines with no ordering guarantee. Concurrent publishers can interleave deliveries. Telemetry and scheduler consumers may assume FIFO ordering incorrectly.

## Proposal

Document ordering semantics explicitly for v1 (none). For v2, evaluate:

- Per-event-type single consumer queue
- Sequence numbers on events
- NATS JetStream or similar if distributed

## Alternatives

- Synchronous publish (blocks publishers — rejected for v1)
- Global ordered log with replay (complexity — defer)

## Open Questions

- What event rate triggers migration from channels to NATS?
- Do any Month 2 subsystems require strict ordering?

## Not Scheduled Until

Measured need during runtime or scheduler integration (Month 2–3).

Spec version: TBD

---
**Layer:** planning
**Related:** [../adr/ADR-003-Event-Driven-Core.md](../adr/ADR-003-Event-Driven-Core.md), [../../planning/decisions/accepted.md](../../planning/decisions/accepted.md)
