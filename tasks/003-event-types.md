# Task 003 — Event Types

## Status

Complete

## Goal

Define typed event definitions for the event-driven core.

## Context

ADR-003 requires decoupled components communicating via events. Typed events prevent stringly-typed bugs and enable metrics mapping.

## Requirements

- Event type enum or string constants for core events
- Event struct with Type, Timestamp, Payload (typed per event)
- Core events: TaskCreated, TaskAssigned, ProcessSpawned, ProcessStateChanged, AgentStarted, AgentStopped, ModelLoaded, ModelUnloaded, SwapStarted, SwapCompleted, KVStored, KVLoaded, TaskCompleted

## Files

- `internal/events/types.go`
- `internal/events/types_test.go`

## Acceptance Criteria

- [x] All core event types defined
- [x] Event struct with type-safe payload variants or generic Payload field
- [x] Tests verify event creation and type identification

## Out Of Scope

- Event bus implementation (task 004)
- NATS adapter

## Related

- ADR-003
- Spec: v0.2+ (events required before registry)
- Architecture: [docs/architecture/kernel.md](../docs/architecture/kernel.md)

---
**Layer:** task
