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

---
**Layer:** planning
