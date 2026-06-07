# Milestone M4 — Hardware Profiling

## Goal

Detect local hardware and classify into Tiny, Standard, or Workstation profiles that guide model selection (ADR-001).

## Actual Outcome

- `internal/hardware/` package ships probe + classification
- `ClassifyProfile` matches architecture tiers in unit tests
- `OSProber` returns live RAM on developer machine

## Demo

```bash
go test ./internal/hardware/...
```

## Issues

- Windows VRAM detection depends on nvidia-smi availability
- Tool registry not shipped (v0.2 partial)

## Lessons

- Injectable ProbeResult keeps tests deterministic while OSProber validates real detection
- CPU-only high-RAM machines correctly stay Tiny per architecture rules

---
**Layer:** planning
**Related:** [../../docs/architecture/hardware.md](../../docs/architecture/hardware.md)
