# ADR-003: Event-Driven Core

## Status

Accepted

## Context

Direct coupling between components (e.g., Planner directly calling Coder) creates hidden dependencies, complicates testing, and prevents independent evolution of subsystems.

TinyBrain requires decoupled components that communicate through well-defined events, enabling the scheduler, runtime, and agents to evolve independently.

## Decision

All inter-component communication uses an event bus. Components publish and subscribe to typed events — they do not call each other directly for lifecycle coordination.

**Event bus v1:** Go channels (in-process).
**Event bus v2 (future):** NATS if distributed deployment requires it.

Core events include: TaskCreated, TaskAssigned, ProcessSpawned, ProcessStateChanged, AgentStarted, AgentStopped, ModelLoaded, ModelUnloaded, SwapStarted, SwapCompleted, KVStored, KVLoaded, TaskCompleted.

## Consequences

### Positive

- Components are independently testable.
- New subsystems can subscribe without modifying publishers.
- Enables metrics pipeline (every event → metric).
- Scales to distributed mode via bus adapter swap.

### Negative

- Harder to trace than direct calls (mitigated by telemetry).
- Event ordering and delivery semantics must be documented.

### Trade-off

Accept event indirection complexity for architectural decoupling.

---
**Layer:** decision
**Related:** [../architecture/kernel.md](../architecture/kernel.md), [../../planning/architecture-evolution/migration-paths.md](../../planning/architecture-evolution/migration-paths.md)
