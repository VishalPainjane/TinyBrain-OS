# Changelog

All notable user-facing changes to TinyBrain OS are documented here.

Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Version numbers match git tags (`v0.6` → section `0.6.0`).

Release process: [planning/releases/RELEASE-CHECKLIST.md](planning/releases/RELEASE-CHECKLIST.md).

---

## [Unreleased]

### Added

(none)

---

## [0.7.0] - 2026-06-23

Tag: `v0.7` — MLFQ scheduler and Memory subsystems release. Details: [planning/releases/v0.7.md](planning/releases/v0.7.md).

### Added

- MLFQ scheduler (`MLFQQueue`, `MLFQScheduler`) with Q0–Q3 queues
- Token quantum demotion via `RecordToken`
- Preemption of lower-priority running processes
- Boost/aging every 500 tokens or 30s
- KV manager stub (`KVStored` / `KVLoaded` events)
- Swap manager for VRAM→RAM tier moves (swap idle heuristic, 10s threshold)
- `brain-top` read-only dashboard prototype
- `QueueDepths()` for telemetry

### Notes

- FIFO scheduler preserved for tests and migration path
- Runtime load orchestration in `Schedule` remains deferred

---

## [0.6.0] - 2026-06-10

Tag: `v0.6` — Inference release. Details: [planning/releases/v0.6.md](planning/releases/v0.6.md).

### Added

- llama.cpp CPU adapter: GGUF load/unload, mmap, single-prompt `Generate` (009a–009b)
- `ModelRuntime` orchestration: resolver → loader → `LlamaProvider` (009c)
- CUDA GPU offload via `-tags cuda` and `NGLayers` / `TB_NGLAYERS` (009d)
- Real-GGUF CI integration jobs with checksum-verified SmolLM2-135M fixture

### Known limitations

- CUDA runtime proof requires manual GPU checklist (no GPU CI runner)
- Metal, ROCm, Vulkan backends deferred

---

## [0.5.0] - 2026-06-07

Tag: `v0.5` — Persistent model registry (bbolt). Details: [planning/releases/v0.5.md](planning/releases/v0.5.md).

---

## [0.4.0] - 2026-06-05

Tag: `v0.4` — Runtime shell + stub provider. Details: [planning/releases/v0.4.md](planning/releases/v0.4.md).

---

## [0.3.0] - 2026-06-01

Tag: `v0.3` — Foundation release: kernel (v0.1), registry/events (v0.2), hardware profiler (v0.3). Details: [planning/releases/v0.3.md](planning/releases/v0.3.md).

---

[Unreleased]: https://github.com/VishalPainjane/TinyBrain-OS/compare/v0.6...HEAD
[0.6.0]: https://github.com/VishalPainjane/TinyBrain-OS/releases/tag/v0.6
[0.5.0]: https://github.com/VishalPainjane/TinyBrain-OS/releases/tag/v0.5
[0.4.0]: https://github.com/VishalPainjane/TinyBrain-OS/releases/tag/v0.4
[0.3.0]: https://github.com/VishalPainjane/TinyBrain-OS/releases/tag/v0.3
