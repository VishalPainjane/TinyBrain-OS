# Completed Tasks

Append-only history. Never delete entries.

## 001-process-state

**Completed:** 2026-06-07
**Commit:** ae9cd8315d8c6ff08827d2a44e85b68440eea668
**Outcome:** Success — ProcessState type with 7 states, unit tests pass
**Files:** internal/process/state.go, internal/process/state_test.go

## 002-process-table

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — Process table with CRUD, O(1) PID lookup, mutex-protected
**Files:** internal/process/types.go, internal/process/table.go, internal/process/table_test.go

## 003-event-types

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — 13 core event types with typed payloads
**Files:** internal/events/types.go, internal/events/types_test.go

## 004-event-bus

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — Channel-based EventBus with pub/sub tests via interface
**Files:** internal/events/bus.go, internal/events/bus_test.go

## 005-agent-registry

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — In-memory agent registry with register/get/list
**Files:** internal/registry/agents.go, internal/registry/agents_test.go

## 006-model-registry

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — In-memory model registry with register/get/list
**Files:** internal/registry/models.go, internal/registry/models_test.go

## 007-hardware-profiler

**Completed:** 2026-06-07
**Commit:** `b83ab88` (tag `v0.3`)
**Outcome:** Success — RAM/VRAM/CPU/backend probe and profile classification
**Files:** internal/hardware/profile.go, internal/hardware/probe.go, internal/hardware/probe_windows.go, internal/hardware/probe_unix.go, internal/hardware/probe_test.go

**Note:** Tasks 002–007 share the Month 1 foundation commit (tag `v0.3`). No per-task historical commits — honest single release.

## 008-runtime

**Completed:** 2026-06-07
**Commit:** `bb0b148`
**Outcome:** Success — ModelRuntime shell with StubProvider, lifecycle events, swap demo test
**Files:** internal/runtime/types.go, internal/runtime/interface.go, internal/runtime/runtime.go, internal/runtime/stub_provider.go, internal/runtime/runtime_test.go

## 009-model-loader

**Completed:** 2026-06-07
**Commit:** `514ff9f`
**Outcome:** Success — Stub loader with lifecycle states, warm/prefetch/evict, LRU eviction shell
**Files:** internal/loader/types.go, internal/loader/loader.go, internal/loader/loader_test.go

## 010-scheduler

**Completed:** 2026-06-07
**Commit:** `339666d`
**Outcome:** Success — FIFO scheduler skeleton, queue policy, state transitions (runtime orchestration deferred)
**Files:** internal/scheduler/queue.go, internal/scheduler/scheduler.go, internal/scheduler/scheduler_test.go

**Note:** Schedule signature drift recorded in [task-010-schedule-signature-drift.md](../decisions/task-010-schedule-signature-drift.md).

## 006-registry-persistence

**Completed:** 2026-06-08
**Commit:** `7de7c70` (tag `v0.5`)
**Outcome:** Success — bbolt `ModelStore`, `models.yaml` seed-on-empty, restart persistence tests pass
**Files:** internal/registry/models.go, internal/registry/models_store.go, internal/registry/models_memory.go, internal/registry/models_bolt.go, internal/registry/models_yaml.go, internal/registry/models_bolt_test.go, testdata/models.yaml

## 009a-llama-cgo-load

**Completed:** 2026-06-08
**Commit:** `f67c491` (merge PR #2)
**Outcome:** Success — llama.cpp CPU CGO LoadModel/UnloadModel via `LlamaProvider`; submodule pinned `b9553` @ `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`; `inference-cgo` CI green
**Files:** internal/inference/llama/, third_party/llama.cpp, .gitmodules, .github/workflows/ci.yml, tests/import_boundary_test.go, tests/deps.go, .gitignore, README.md (build docs)

**Note:** Generate, runtime wiring, and CUDA execution deferred to 009b+.

## 009b-cpu-generate

**Completed:** 2026-06-09
**Commit:** *(pending — commit after merge approval)*
**Outcome:** Success — CPU `Generate` on `LlamaProvider` via llama.cpp b9553; context at load; memory clear per call; Linux CGO + real GGUF integration verified (TTFT/TPS, 3× cycle)
**Files:** `internal/inference/llama/` (generate_cpu.go, generate_stub.go, generate_integration_test.go, bindings_cgo.go, provider.go, config.go, errors.go, port_stubs.go, provider_test.go, load_cpu.go, context.go, doc.go), `planning/decisions/009b-architecture-review.md`

**Note:** Runtime ↔ `LlamaProvider` wiring deferred to 009c. CUDA/GPU deferred. SaveContext/RestoreContext remain stubs (011).

---
**Layer:** planning
