# Architecture — Current State

Code-first snapshot. Update after every completed task.

**Last verified:** 2026-06-24

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
| MLFQ Scheduler (Q0-Q3, boost, preemption) | `internal/scheduler/` | Complete, tested |
| KV Manager (stub metadata, events) | `internal/kv/` | Complete, tested |
| Swap Manager (events, VRAM↔RAM shell) | `internal/swap/` | Complete, tested |
| brain-top TUI | `cmd/brain-top/` | Complete, production (023) |
| Agent plugin contract & sample executor | `internal/agents/` | Complete |
| Agent registry config (fleet.yaml) | `internal/registry/` | Complete |
| Event-driven workflow orchestration | `internal/router/` | Complete |
| llama.cpp adapter (CPU load/unload/generate) | `internal/inference/llama/` | Complete CPU path; wired via runtime (009c) |
| llama.cpp CUDA offload (dynamic DLL, static CGO, `NGLayers`) | `internal/inference/llama/` | **Complete** — Windows dynamic CUDA backend + static Linux CGO verified (manual GPU checklist signed off) |
| tinybrain CLI (doctor, probe, models, run, status, workflow) | `cmd/tinybrain/` | Complete |
| Kubernetes Operator (Agent, Task, memory CRDs & controllers) | `internal/k8s/`, `cmd/operator/` | Complete, tested (M9) |
| Scheduler ↔ runtime command wiring (via EventBus) | `internal/agents/`, `internal/scheduler/` | Complete, event-driven |
| Benchmark Suite | `cmd/benchmark/` | Complete (018) |

## In Progress

(none)

## Not Implemented

- Metal / ROCm / Vulkan inference backends (CUDA adapter verified & signed off via ADR-006 DLL loading; Metal/ROCm/Vulkan deferred)
- Tool registry
- Telemetry package

## Active Packages

`internal/process/`, `internal/events/`, `internal/registry/`, `internal/hardware/`, `internal/runtime/`, `internal/loader/`, `internal/scheduler/`, `internal/inference/llama/`, `internal/kv/`, `internal/swap/`, `internal/agents/`, `internal/router/`, `internal/k8s/`, `cmd/tinybrain/`, `cmd/brain-top/`, `cmd/operator/`

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
