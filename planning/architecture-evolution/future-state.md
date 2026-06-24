# Architecture — Future State

Target end-state for TinyBrain OS. Not a commitment timeline — see [../roadmap/master-roadmap.md](../roadmap/master-roadmap.md) and [../../docs/specs/](../../docs/specs/).

## Runtime Kernel

- Hardware-aware boot sequence: detect → classify → map agents to models
- Process table with full lifecycle state machine
- Event-driven decoupled core (no direct component-to-component calls)

## Resource Management

- Dynamic model loading, unloading, warming, prefetching
- MLFQ scheduler with VRAM-aware preemption and token-quantum time slices
- KV hibernation across VRAM → RAM → NVMe memory tiers
- KV cache compression pipeline
- LRU eviction when VRAM exceeds threshold

## Plugin System

- Agent, model, and tool definitions in registry
- User-defined agent fleets via configuration (not hardcoded roles)
- Tool execution layer separate from agents

## Inference

- `InferenceProvider` interface with local GGUF (llama.cpp) primary adapter
- Optional cloud adapters (OpenAI, Anthropic, vLLM) without scheduler changes

## Visibility

- brain-top TUI: live process states, VRAM/RAM usage, queue depths, swap monitor
- OpenTelemetry traces at V2

## Platform (post-v1.0)

- Kubernetes operator with Agent, Task, KVCache, SwapPolicy CRDs
- Docker Compose and Helm deployment modes

---
**Layer:** planning
**Related:** [current-state.md](current-state.md), [migration-paths.md](migration-paths.md)
