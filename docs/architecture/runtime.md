# Runtime

The runtime manages model lifecycle and inference execution. It is the only component that invokes inference providers.

## Responsibilities

- Load, unload, warm, prefetch, and evict models
- Execute generate requests via InferenceProvider
- Save and restore execution context (KV cache coordination)
- Report runtime metrics

## Model Lifecycle States

| State | Meaning |
|-------|---------|
| NOT_LOADED | Model not in memory |
| LOADING | Weights being mapped/loaded |
| ACTIVE | Model loaded, ready for inference |
| WARM | Loaded but idle; candidate for retention |
| UNLOADING | Being removed from memory |
| UNLOADED | Fully removed |

Prevents duplicate loads and enables predictive prefetch (e.g., load coder while planner runs).

## API Surface (conceptual)

The runtime exposes: LoadModel, UnloadModel, Generate, SaveContext, RestoreContext. See [contracts/runtime.md](../contracts/runtime.md).

## Inputs

Scheduler commands; generate requests from agent plugins; registry model definitions.

## Outputs

Inference results; model state events; resource metrics.

## Dependencies (allowed)

InferenceProvider, model loader, KV manager, registry (read).

## Dependencies (forbidden)

Scheduler, UI, direct agent logic.

## Future Plans

Predictive loading based on agent transition patterns; LRU eviction when VRAM full; swap coordination with memory layer.

## Non-Goals

Queue management; task routing; agent business logic.

## Related Contracts

[runtime.md](../contracts/runtime.md)

## Related ADRs

ADR-004 (hexagonal architecture)

---
**Layer:** architecture
**Source:** detail.md Part 2, 4, 5
**Related:** [memory.md](memory.md), [invariants.md](invariants.md)
