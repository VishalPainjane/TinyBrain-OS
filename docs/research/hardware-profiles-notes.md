# Hardware Profiles Notes

**Non-normative — do not treat as architecture law.**

## Profile Tiers

| Profile | RAM | GPU | Design intent |
|---------|-----|-----|---------------|
| Tiny | 8 GB | CPU only | Minimum viable local AI |
| Standard | 16 GB | RTX 3050 4GB | **Primary design target** |
| Workstation | 64 GB+ | High-end | Aggressive orchestration |

## Capability Tiers (future granularity)

detail.md proposes Tier 0–4 based on boot probe:

- Tier 0: CPU only
- Tier 1: 8 GB RAM
- Tier 2: 16 GB + RTX 3050
- Tier 3: 32 GB + RTX 4070 class
- Tier 4: 64 GB+ workstation

## Model Fleet Examples (config only — not hardcoded)

**Standard profile:** smaller models, aggressive swapping, possibly shared model for multiple capabilities.

**Workstation profile:** larger models per capability, less swapping, concurrent warm models.

Same code, different registry configs selected by hardware probe.

## Benchmark Considerations

Cross-profile benchmarks are not directly comparable. Document profile in all benchmark results. Assumption tracked in [planning/assumptions.md](../../planning/assumptions.md).

---
**Layer:** research
**Related:** [../architecture/hardware.md](../architecture/hardware.md)
