# Task 007 — Hardware Profiler

## Status

Complete

## Goal

Detect system hardware and classify into TinyBrain hardware profiles.

## Context

ADR-001 requires hardware-aware model selection. Boot sequence Phase 1–2: detect resources, assign profile.

## Requirements

- Detect total RAM, available VRAM (if GPU present)
- Detect CPU info
- Detect inference backend: CUDA, Metal, ROCm, or CPU-only
- Classify into profile: Tiny, Standard, or Workstation
- Expose HardwareProfile struct for registry/runtime consumption

## Files

- `internal/hardware/probe.go`
- `internal/hardware/profile.go`
- `internal/hardware/probe_test.go`

## Acceptance Criteria

- [x] Probe returns RAM bytes on current machine
- [x] Profile classification logic matches [docs/architecture/hardware.md](../docs/architecture/hardware.md)
- [x] CPU-only machine classified as Tiny
- [x] Tests use injected/mock probe values

## Out Of Scope

- Model assignment (registry responsibility)
- Scheduler integration

## Related

- Spec: [docs/specs/v0.3-hardware.md](../docs/specs/v0.3-hardware.md)
- ADR-001

---
**Layer:** task
