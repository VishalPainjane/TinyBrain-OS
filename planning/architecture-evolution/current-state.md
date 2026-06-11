# Architecture — Current State

Code-first snapshot. Update after every completed task.

**Last verified:** 2026-06-11

## Implemented

| Component | Package | Status |
|-----------|---------|--------|
| ProcessState (7 states) | `internal/process/state.go` | Complete, tested |
| Process table (CRUD, O(1) lookup) | `internal/process/table.go` | Complete, tested |
| Event types (13 core events) | `internal/events/types.go` | Complete, tested |
| Event bus (channel pub/sub) | `internal/events/bus.go` | Complete, tested |
| Agent registry | `internal/registry/agents.go` | Complete, tested |
| Model registry (in-memory + bbolt) | `internal/registry/models*.go` | Complete, tested |
| Hardware probe + profiles | `internal/hardware/` | Complete, tested |
| Runtime + integrated orchestration | `internal/runtime/` | Complete, tested (009c) |
| ModelResolver + RegistryResolver | `internal/runtime/` | Complete, tested (009c) |
| Model loader (stub lifecycle) | `internal/loader/` | Complete, tested |
| FIFO scheduler skeleton | `internal/scheduler/` | Complete, tested |
| llama.cpp adapter (CPU load/unload/generate) | `internal/inference/llama/` | Complete CPU path; wired via runtime (009c) |
| llama.cpp CUDA offload (`-tags cuda`, `NGLayers`) | `internal/inference/llama/` | **Partial** (009d merged `ab06c60`); manual GPU checklist open |
| tinybrain CLI (doctor, probe, models, run, status) | `cmd/tinybrain/` | Complete (stab-003) |

## In Progress

(none)

## Not Implemented

- Metal / ROCm / Vulkan inference backends (CUDA adapter shipped 009d; runtime GPU proof manual)
- Agent plugins, tool registry
- Telemetry package, brain-top TUI
- KV manager, swap manager
- Kubernetes operator
- Scheduler → runtime command wiring

## Active Packages

`internal/process/`, `internal/events/`, `internal/registry/`, `internal/hardware/`, `internal/runtime/`, `internal/loader/`, `internal/scheduler/`, `internal/inference/llama/`, `cmd/tinybrain/`

## Tests

`go test ./...` — passing (`CGO_ENABLED=0`, default CI)  
`inference-cgo` CI job — CGO unit tests (Linux)  
Integration (`-tags integration`, `TB_TEST_GGUF_PATH`) — inference + runtime E2E (`runtime_integration_test.go`)  
`tests/import_boundary_test.go` — scheduler and runtime must not import inference

## Governance

Repository rules frozen in `.cursor/rules/` (7 files). CI via `.github/workflows/ci.yml` (`test`, `inference-cgo`, `inference-integration`, `inference-integration-runtime`).

---
**Layer:** planning
**Related:** [future-state.md](future-state.md), [../../docs/current.md](../../docs/current.md)
