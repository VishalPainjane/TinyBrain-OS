# Architecture Overview

TinyBrain OS is a local AI runtime that coordinates specialized agents under strict hardware budgets. It behaves like an operating system for AI workloads — not a single-model chatbot.

## High-Level Flow

```text
User
  ↓
API
  ↓
Router
  ↓
Scheduler
  ↓
Runtime
  ↓
InferenceProvider
  ↓
Models (GGUF)
```

1. User submits a task via API.
2. Router classifies the task and selects an agent capability from the registry.
3. Scheduler assigns priority, admits a process, and decides when it runs.
4. Runtime loads the required model, executes inference, manages context.
5. Agent plugin produces structured JSON output.
6. Telemetry records metrics; resources are released when done.

## TinyBrain OS v2 Architectural Specification

To support Iteration-Level Scheduling and Paged Memory, the architecture is decoupled into a clear Control Plane and Data Plane:

```text
                  +---------------------------------------+
                  |             CONTROL PLANE             |
                  |  [Request Ingest] -> [MLFQ Queue]     |
                  +-------------------+-------------------+
                                      | (Token Boundary Decisions)
                                      v
                  +---------------------------------------+
                  |              DATA PLANE               |
                  |  +---------------------------------+  |
                  |  | Persistent Inference Daemon     |  |
                  |  | (Weights Pinned in VRAM)        |  |
                  |  +---------------------------------+  |
                  |  | Paged KV-Cache Memory Allocator |  |
                  |  | [VRAM Blocks] <-> [Host RAM]    |  |
                  |  +---------------------------------+  |
                  +---------------------------------------+
```

### The Control Plane (Go / Engine Orchestration)
- **Admission Controller**: Receives `TaskCreated` events. Validates token budgets and registers sequence IDs.
- **MLFQ Iteration Scheduler**: Manages a pool of active sequences. Every iteration, it constructs an array of sequence pointers up to `MaxBatchSize` to pass to the engine backend.

### The Data Plane (CGO / CUDA Backend)
- **Resident Model Worker**: Loads the quantized model once during boot. Exposes an iterative execution function: `ExecuteIteration(Batch* batch)`.
- **Block Memory Manager**: Manages a fixed array of pre-allocated VRAM blocks. Tracks block lifespans, allocations, and coordinates asynchronous host-to-device memory transfers (`cudaMemcpyAsync`) during context swaps.

## Main Layers

| Layer | Responsibility |
|-------|----------------|
| Interface | Entry point — REST/gRPC API, status, streaming |
| Router | Task classification, applies model-specific chat templates before dispatch, agent selection |
| Scheduler | Order, priority, preemption, fairness, VRAM budget |
| Runtime | Model lifecycle, generate, context save/restore |
| Registry | Agent, model, tool definitions |
| Process (Kernel) | OS-style process table and lifecycle states |
| Memory | KV cache tiers, swap, hibernation |
| Telemetry | Metrics, logs, brain-top visibility |
| Agents | Plugin workers — configured, not hardcoded |

## Deployment Modes

| Mode | Command / method | Use case |
|------|------------------|----------|
| Developer | `tinybrain run` (future) | Local development |
| Docker | `docker compose up` (future) | Reproducible environment |
| Kubernetes | `helm install tinybrain` (future) | Platform deployment |

## Core Principles

- Local-first, hardware-aware
- Dynamic model swapping — not everything loaded at once
- Structured JSON IPC — not natural language between components
- Agents as plugins — registry-driven, not hardcoded roles
- Scheduler independent from inference engine

## Responsibilities

Describe the end-to-end system shape and layer boundaries.

## Inputs

User tasks, hardware environment, registry configurations.

## Outputs

Structured task results, telemetry, resource state.

## Dependencies (allowed)

All layers may read from registry and emit events.

## Dependencies (forbidden)

UI driving scheduler; scheduler calling inference directly.

## Future Plans

Kubernetes operator, brain-top TUI, benchmark suite, cloud provider adapters.

## Non-Goals

This document does not define interfaces (see contracts/) or implementation (see tasks/).

## Related Contracts

[process.md](../contracts/process.md), [registry.md](../contracts/registry.md), [runtime.md](../contracts/runtime.md), [scheduler.md](../contracts/scheduler.md)

## Related ADRs

ADR-001, ADR-002, ADR-003, ADR-004, ADR-005, ADR-008

---
**Layer:** architecture
**Source:** detail.md Part 2, System-map.md
**Related:** [kernel.md](kernel.md), [invariants.md](invariants.md)
