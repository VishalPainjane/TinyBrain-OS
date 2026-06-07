# Task 004 — Event Bus

## Status

Complete

## Goal

Implement in-process pub/sub event bus using Go channels (v1).

## Context

Components must not call each other directly for lifecycle coordination. The bus enables ADR-003 and future NATS migration via interface swap.

## Requirements

- EventBus interface: Publish, Subscribe, Unsubscribe (implemented as `Subscribe() (unsubscribe func())` — see planning/decisions/accepted.md)
- Channel-based implementation
- Subscribers receive events without blocking publishers (buffered channels or goroutine dispatch)
- Behind interface for future NATS adapter

## Files

- `internal/events/bus.go`
- `internal/events/bus_test.go`

## Acceptance Criteria

- [x] Publish delivers event to all subscribers
- [x] Unsubscribe stops delivery
- [x] Multiple subscribers on same event type work
- [x] Tests use interface — not concrete channel type

## Out Of Scope

- NATS, Kafka
- Persistence, replay

## Related

- ADR-003
- Task: [003-event-types.md](003-event-types.md)
- Migration: [planning/architecture-evolution/migration-paths.md](../planning/architecture-evolution/migration-paths.md)

---
**Layer:** task
