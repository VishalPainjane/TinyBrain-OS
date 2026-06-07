# KV Cache Notes

**Non-normative — do not treat as architecture law.**

## What Is KV Cache

During autoregressive inference, attention key-value matrices grow with each token. Recomputing from scratch on resume is expensive. KV cache stores these matrices for reuse.

## Context Hibernation

When a process is paused/preempted, save KV cache instead of discarding. On resume, restore KV — avoid reprocessing the full prompt.

TinyBrain states: PREEMPTED (KV may remain in VRAM), HIBERNATED (KV moved to lower tier, weights unloaded).

## Compression Target

Proposed pipeline: FP16 → INT8 → Q4 for stored KV. Goal ~72% memory reduction (detail.md Part 4). Enables more suspended contexts on constrained hardware.

## Research References

LMCache, KVSwap-style architectures exploring KV offloading for throughput gains. See arXiv literature on hierarchical KV cache systems.

## RFC

Full implementation proposal: [RFC-001-KV-Hibernation.md](../rfc/RFC-001-KV-Hibernation.md)

---
**Layer:** research
