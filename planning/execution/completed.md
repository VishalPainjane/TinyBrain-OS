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
**Commit:** `8cdd9b1` (merged via PR #4, `8baa77e`)
**Outcome:** Success — CPU `Generate` on `LlamaProvider` via llama.cpp b9553; context at load; memory clear per call; Linux CGO + real GGUF integration verified (TTFT/TPS, 3× cycle)
**Files:** `internal/inference/llama/` (generate_cpu.go, generate_stub.go, generate_integration_test.go, bindings_cgo.go, provider.go, config.go, errors.go, port_stubs.go, provider_test.go, load_cpu.go, context.go, doc.go), `planning/decisions/009b-architecture-review.md`

**Note:** Runtime ↔ `LlamaProvider` wiring deferred to 009c. CUDA/GPU deferred. SaveContext/RestoreContext remain stubs (011).

## 009c-runtime-integration

**Completed:** 2026-06-09
**Commit:** `ee186c1` (merged via PR #5, `f20475f`)
**Outcome:** Success — `ModelRuntime` orchestrates `ModelResolver` → `loader.Load` → `LlamaProvider`; shared `runtime.ModelResolver` (Option B); rollback on provider failure; lifecycle events; loader-less constructor preserved; E2E `runtime_integration_test.go`; `TestRuntimeDoesNotImportInference`
**Files:** `internal/runtime/resolver.go`, `registry_resolver.go`, `registry_resolver_test.go`, `runtime.go`, `runtime_loader_test.go`, `runtime_integration_test.go`, `internal/inference/llama/provider.go` (+ test updates), `tests/import_boundary_test.go`, `tasks/009c-runtime-integration.md`, `planning/decisions/009c-architecture-review.md`

**Note:** CUDA/GPU deferred. SaveContext/RestoreContext remain stubs (011). CI E2E requires `TB_TEST_GGUF_PATH` on Linux.

## STAB-002-ci-observability

**Completed:** 2026-06-09
**Commit:** `30c0a39` (merge PR #8); initial instrumentation `0399025` (merge PR #7)
**Verified:** `main` CI run [27206160391](https://github.com/VishalPainjane/TinyBrain-OS/actions/runs/27206160391)
**Outcome:** Success — per-job Step Summary metrics; per-job `ci-metrics-*` artifacts; merged `ci-run-record-{run_id}` artifact on each `main` push; non-blocking `ci-metrics-collect` (no git writes to protected `main`); runtime M2 diagnostics parity; `ci-runs.jsonl` removed in favor of artifact-only history
**Files:** `.github/workflows/ci.yml`, `.github/scripts/ci-emit-metrics.sh`, `.github/scripts/ci-collect-run-metrics.py`, `planning/metrics/ci-{schema,baseline}.md`, `testdata/ci/README.md`, `tasks/stab-002-ci-observability.md`, `planning/decisions/ci-observability-architecture-review.md` (v2 delta)

**Note:** STAB-001 (real inference CI gate) shipped in PR #6 (`e293b98`). STAB-002 builds on four merge-blocking jobs without changing gate logic.

## 009d-gpu-offload-cuda

**Completed:** 2026-06-10
**Commit:** `ab06c60` (merge PR #9); `10dd8e5` feat, `4a6cd47` CGO preamble fix
**Outcome:** Success — CUDA GPU offload via `-tags cuda`; `LlamaConfig.NGLayers`; CGO split (`bindings_common` / `bindings_cpu` / `bindings_cuda`); `ConfigFromProbe`; guarded `cuda_integration_test.go`; CPU `inference-cgo` + integration jobs green on PR #9
**Files:** `internal/inference/llama/` (load_cuda.go, generate_cuda.go, bindings_*.go, config_probe.go, config_test.go, cuda_integration_test.go), `tasks/009d-gpu-offload-cuda.md`, `planning/decisions/009d-architecture-review.md`, `planning/decisions/009d-manual-gpu-checklist.md`

**Note:** Runtime/scheduler/loader/registry unchanged. CUDA matrix **Partial** until manual GPU checklist signed. No bundled CUDA runtime libs.

## stab-003-tinybrain-cli

**Completed:** 2026-06-11
**Commit:** `f6c8031`
**Outcome:** Success — `cmd/tinybrain` with doctor, probe, models list, run, status, version; fuzz + golden tests; README quick start
**Files:** `cmd/tinybrain/`, `tasks/stab-003-tinybrain-cli.md`, `README.md`

**Note:** `tinybrain run` requires `CGO_ENABLED=1` and built llama.cpp. Post-v0.6 product shell — tag `v0.6` applies to inference release docs commit, not necessarily this commit.

## repo-hygiene-testing

**Completed:** 2026-06-11
**Commit:** `f04feb1`
**Outcome:** Success — solo-dev hygiene docs, expanded boundary tests (INV-001/002/008), registry YAML fuzz, CLI golden tests
**Files:** `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`, `docs/testing-policy.md`, `planning/releases/RELEASE-CHECKLIST.md`, `planning/metrics/repo-health.md`, `tests/import_boundary_test.go`, `internal/registry/models_yaml_fuzz_test.go`, `cmd/tinybrain/*_fuzz_test.go`, `cmd/tinybrain/golden_test.go`

## 010-scheduler-mlfq

**Completed:** 2026-06-11
**Commit:** `a0f90a7` (merge PR #12; feat `55bbfa0`)
**Outcome:** Success — MLFQ Q0–Q3, token quantum demotion, preemption, boost/aging, swap idle heuristic; FIFO scheduler preserved
**Files:** `internal/scheduler/mlfq.go`, `internal/scheduler/queue.go`, `internal/scheduler/scheduler.go`, `internal/scheduler/scheduler_test.go`, `docs/specs/v0.7-scheduler.md`, `docs/contracts/scheduler.md`, `tasks/010-scheduler.md`

**Note:** Extends FIFO skeleton (`339666d`). Runtime load orchestration and swap execution remain deferred (011/012).

## 011-kv-manager

**Completed:** 2026-06-11
**Commit:** `1a1acf0` (merge PR #13; feat `8724988`)
**Outcome:** Success — Stub KV block pool, VRAM↔RAM save/load shell, KVStored/KVLoaded events on bus
**Files:** `internal/kv/types.go`, `internal/kv/manager.go`, `internal/kv/manager_test.go`, `tasks/011-kv-manager.md`

**Note:** Metadata-only stub; llama.cpp KV export and compression deferred per RFC-001.

## 012-swap-manager

**Completed:** 2026-06-11
**Commit:** `fc0a093` (merge PR #14; feat `10aad3f`)
**Outcome:** Success — VRAM↔RAM swap orchestration via kv.Save/Load, process HIBERNATED/READY, SwapStarted/SwapCompleted events
**Files:** `internal/swap/manager.go`, `internal/swap/manager_test.go`, `tasks/012-swap-manager.md`

**Note:** Uses `scheduler.ShouldSwap` for swap-out eligibility; no scheduler policy changes.

## 013-brain-top

**Completed:** 2026-06-11
**Commit:** `fc0a093` (merge PR #14; feat `df67c1a`)
**Outcome:** Success — read-only `cmd/brain-top` prototype (process panel, MLFQ depths, resources); snapshot/watch modes
**Files:** `cmd/brain-top/main.go`, `cmd/brain-top/render.go`, `cmd/brain-top/render_test.go`, `tasks/013-brain-top.md`

**Note:** Stdlib only; Bubble Tea and live kernel IPC deferred to Month 6.

## 014-agent-plugin

**Completed:** 2026-06-23
**Commit:** `5c6280e` (merge PR to main)
**Outcome:** Success — Agent plugin contract (Executor interface, Agent interface, SampleAgent implementation)
**Files:** `docs/contracts/agents.md`, `internal/agents/executor.go`, `internal/agents/interface.go`, `internal/agents/sample.go`, `tasks/014-agent-plugin.md`

## 015-fleet-registry

**Completed:** 2026-06-23
**Outcome:** Success — YAML loader for sample fleet configuration, supporting multiple agents
**Files:** `internal/registry/agents_yaml.go`, `internal/registry/agents_yaml_test.go`, `testdata/fleet.yaml`, `tasks/015-fleet-registry.md`

## 016-event-pipeline

**Completed:** 2026-06-23
**Outcome:** Success — Decoupled event-driven execution using router, scheduler coordinator, and agent listener
**Files:** `internal/router/router.go`, `internal/scheduler/coordinator.go`, `internal/agents/listener.go`, `tests/pipeline_integration_test.go`, `tasks/016-event-pipeline.md`

## 017-workflow-demo

**Completed:** 2026-06-23
**Commit:** `a40b0fd` (tag `v0.8`)
**Outcome:** Success — `tinybrain workflow` orchestrates sequential 2-agent execution via task completion events
**Files:** `cmd/tinybrain/workflow.go`, `cmd/tinybrain/main.go`, `internal/events/types.go`, `internal/agents/listener.go`, `tasks/017-workflow-demo.md`

---
**Layer:** planning
