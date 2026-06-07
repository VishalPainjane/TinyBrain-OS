# Hardware

TinyBrain adapts behavior to detected hardware rather than assuming a fixed environment.

## Responsibilities

- Detect RAM, VRAM, CPU, inference backend (CUDA, Metal, ROCm)
- Classify into hardware profiles
- Feed profile into registry model selection and scheduler aggressiveness

## Hardware Profiles

| Profile | RAM | GPU | Target |
|---------|-----|-----|--------|
| Tiny | 8 GB | CPU only | Most constrained |
| Standard | 16 GB | RTX 3050 (4 GB VRAM) | **Baseline design target** |
| Workstation | 64 GB+ | High-end GPU | High-capability systems |

Profiles are inputs to runtime, scheduler, and registry — not decorative labels.

## Boot Sequence (future)

1. **Phase 1 — Hardware Detection:** collect RAM, VRAM, CPU, backend
2. **Phase 2 — Capability Classification:** assign profile tier (0–4)
3. **Phase 3 — Agent Mapping:** registry selects best available models per agent capability

Same codebase, different model fleets on different machines — no user configuration required.

## Baseline Constraint

Primary design envelope: **16 GB RAM, RTX 3050, 4 GB VRAM** (consumer laptop).

## Inputs

System probes at startup; optional user override (future).

## Outputs

Hardware profile; capability tier; recommended model fleet.

## Dependencies (allowed)

OS APIs for memory/GPU detection.

## Dependencies (forbidden)

Hardcoded model assignments; scheduler logic.

## Future Plans

Tier 0–4 granularity; backend-specific optimizations; benchmark-driven profile tuning.

## Non-Goals

Model inference; agent execution.

## Related ADRs

ADR-001

---
**Layer:** architecture
**Source:** detail.md Part 1, Day 6 hardware-profiles.md
**Related:** [registry.md](registry.md), [research/hardware-profiles-notes.md](../research/hardware-profiles-notes.md)
