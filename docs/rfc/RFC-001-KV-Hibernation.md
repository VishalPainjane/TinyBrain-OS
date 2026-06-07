# RFC-001: KV Hibernation

## Status

Proposal

## Problem

On Standard profile (4 GB VRAM), only one small model fits comfortably with active KV cache. Pausing agents without KV preservation forces expensive prompt recomputation. Most agent frameworks ignore KV cache lifecycle.

## Proposal

Implement full KV hibernation pipeline:

1. On PREEMPTED/HIBERNATED: serialize KV cache from VRAM
2. Compress FP16 → Q4 for RAM/SSD storage (~72% target reduction)
3. Track blocks in KV Block Pool with location metadata
4. On resume: restore KV to VRAM before generation continues
5. Tier movement: VRAM → RAM → NVMe based on idle time and pressure

## Alternatives

- Discard KV on pause (simple, bad performance)
- Keep all KV in VRAM (doesn't scale)
- External Redis for KV (rejected — local-first)

## Open Questions

- llama.cpp KV export API stability
- Compression done in Go vs native library
- Block size granularity

## Not Scheduled Until

v1.0 (basic save/load); full compression pipeline post-v1.0

---
**Layer:** planning
**Related:** [../research/kv-cache-notes.md](../research/kv-cache-notes.md), [../architecture/memory.md](../architecture/memory.md)
