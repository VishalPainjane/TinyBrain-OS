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

## Paged KV-Cache Memory Allocator (Virtual Memory)

To prevent severe VRAM fragmentation caused by contiguous memory allocation for unpredictable sequence lengths, the kernel implements a **Paged KV-Cache Manager** (inspired by PagedAttention). 

Every context sequence is divided into logical pages, mapped to fixed-size blocks/pages managed by a page table. The exact block size and reservation ratios are backend implementation details.

The Memory Manager maintains a block table (page table) that tracks:
- Logical Sequence ID to Physical Block mapping
- Block Owner (agent/process), location (VRAM/RAM/SSD), and last access time

Responsibilities: allocate blocks, free blocks, map virtual to physical addresses, and handle PCIe swap (Host-to-Device / Device-to-Host memory transfers for KV blocks).

## KV Compression (future)

Target pipeline: FP16 → INT8 → Q4 for KV storage (~72% reduction goal).

## Memory Pressure Policy

When VRAM blocks are exhausted, the kernel enforces a strict memory pressure policy to maintain goodput (not just raw throughput):
1. **Admit**: If sufficient VRAM blocks exist, new sequences are admitted to the active batch.
2. **Queue**: If VRAM is full, new sequences wait in the MLFQ queue (impacting `queue_wait_time` but preventing OOM).
3. **Spill (Preempt & Swap)**: The MLFQ scheduler preempts a lower-priority active sequence. The Memory Manager performs a context swap by copying the sequence's **KV cache blocks** from VRAM back to the host system RAM via `cudaMemcpy`. The **model weights are never moved or unloaded** during a sequence context switch.
4. **Reject**: If both VRAM and host RAM swap space are exhausted, or wait times exceed SLOs, new tasks are rejected with a 429/ResourceExhausted error.

Process in HIBERNATED state: KV preserved in lower tier (RAM or NVMe), allowing it to resume decoding later without recomputing the prefill prompt.

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
ADR-008 (Iteration-level scheduling and paged memory)

---
**Layer:** architecture
**Source:** detail.md Part 4
**Related:** [runtime.md](runtime.md), [glossary.md](../glossary.md)
