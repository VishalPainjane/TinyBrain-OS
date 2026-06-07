# Kernel

The kernel is the foundation of TinyBrain OS. It provides the process abstraction — the OS analogue for agent execution.

## Responsibilities

- Define process lifecycle states (ProcessState)
- Maintain the process table (PID → process metadata)
- Track state, priority, memory usage, KV cache references
- Support boot phases: hardware detection → capability classification → agent mapping (future)

## Process Lifecycle States

| State | Meaning |
|-------|---------|
| NEW | Created, not yet admitted to scheduler |
| READY | Eligible to run, waiting for resources |
| RUNNING | Actively executing, holds resources |
| WAITING | Blocked on external dependency or tool |
| PREEMPTED | Interrupted, removed from active execution |
| HIBERNATED | Suspended, KV preserved, weights unloaded |
| TERMINATED | Completed or stopped permanently |

Every process is always in exactly one state.

## Process Table

Like a Linux process table. Tracks all active and recent processes:

- PID, agent capability reference, state, priority
- Memory usage (RAM, VRAM), KV cache ID
- Last execution time, tokens produced, associated task ID

brain-top (future) reads directly from this table.

## Event Catalog (Core)

Events emitted by kernel and consumed by scheduler/runtime:

- TaskCreated, TaskAssigned
- ProcessSpawned, ProcessStateChanged
- AgentStarted, AgentStopped
- SwapStarted, SwapCompleted
- KVStored, KVLoaded

## Inputs

Process creation requests from scheduler; state transition commands.

## Outputs

Process table queries; state change events.

## Dependencies (allowed)

Event bus (publish state changes).

## Dependencies (forbidden)

Inference engines, model loaders, UI.

## Future Plans

Boot sequence integration with hardware profiler; process groups; resource limits per process.

## Non-Goals

Scheduling policy (scheduler owns transitions logic); model loading (runtime owns).

## Related Contracts

[process.md](../contracts/process.md)

## Related ADRs

ADR-003 (event-driven core)

---
**Layer:** architecture
**Source:** detail.md Part 4–5, process-model.md
**Related:** [scheduler.md](scheduler.md), [glossary.md](../glossary.md)
