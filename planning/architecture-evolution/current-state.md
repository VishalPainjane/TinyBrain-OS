# Architecture — Current State

Code-first snapshot. Update after every completed task.

**Last verified:** 2026-06-09

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
| Runtime shell + StubProvider | `internal/runtime/` | Complete, tested |
| Model loader (stub lifecycle) | `internal/loader/` | Complete, tested |
| FIFO scheduler skeleton | `internal/scheduler/` | Complete, tested |
| llama.cpp adapter (CPU load/unload/generate) | `internal/inference/llama/` | Partial (009a+009b) — context at load; Generate verified; not wired to runtime |

## In Progress

(none)

## Not Implemented

- Runtime ↔ LlamaProvider integration (009c)
- CUDA / Metal / ROCm / Vulkan inference backends
- Agent plugins, tool registry
- Telemetry, brain-top TUI
- KV manager, swap manager
- Kubernetes operator

## Active Packages

`internal/process/`, `internal/events/`, `internal/registry/`, `internal/hardware/`, `internal/runtime/`, `internal/loader/`, `internal/scheduler/`, `internal/inference/llama/`

## Tests

`go test ./...` — passing (`CGO_ENABLED=0`, default CI)  
`inference-cgo` CI job — passing (Linux CGO unit tests)  
Integration (`-tags integration`, `TB_TEST_GGUF_PATH`) — verified Linux WSL (009b)  
`tests/import_boundary_test.go` — scheduler must not import inference

## Governance

Repository rules frozen in `.cursor/rules/` (7 files). CI via `.github/workflows/ci.yml` (`test` + `inference-cgo`).

---
**Layer:** planning
**Related:** [future-state.md](future-state.md), [../../docs/current.md](../../docs/current.md)
