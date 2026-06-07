# Memory

TinyBrain treats memory as a managed, tiered resource — analogous to CPU cache hierarchy.

## Responsibilities

- KV cache allocation, movement, compression, restore, deletion
- Three-tier storage hierarchy
- LRU eviction when VRAM exceeds threshold
- Context hibernation coordination

## Memory Hierarchy

```text
VRAM  — active weights, active KV (< 90% utilization target)
  ↓
RAM   — compressed KV, warm models (< 70% utilization target)
  ↓
NVMe  — cold KV, archived sessions (swap partition for AI)
```

## KV Block Pool

Every context becomes a tracked block:

- ID, owning agent/process, size, location (VRAM/RAM/SSD), last access time

KV Block Manager responsibilities: allocate, move, compress, restore, delete — nothing else.

## KV Compression (future)

Target pipeline: FP16 → INT8 → Q4 for KV storage (~72% reduction goal).

## Context Hibernation

Process in HIBERNATED state: weights unloaded, KV preserved in lower tier. Avoids prompt recomputation on resume.

## Inputs

Runtime save/restore requests; scheduler swap decisions; process state changes.

## Outputs

KV location metadata; tier movement events; memory utilization metrics.

## Dependencies (allowed)

Runtime (context coordination), process table (owner tracking).

## Dependencies (forbidden)

Scheduling policy; inference execution.

## Future Plans

Full compression pipeline; predictive KV prefetch; KVCache Kubernetes CRD.

## Non-Goals

Model weight loading (runtime/loader); scheduling decisions.

## Related ADRs

RFC-001 (KV Hibernation proposal)

---
**Layer:** architecture
**Source:** detail.md Part 4
**Related:** [runtime.md](runtime.md), [glossary.md](../glossary.md)
