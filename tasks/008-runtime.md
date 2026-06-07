# Task 008 — Runtime

## Status

Not Started

## Goal

Implement ModelRuntime shell with StubProvider for testing.

## Context

Runtime is the moat — agents are replaceable. Scheduler delegates all model operations here. ADR-004 defines InferenceProvider boundary.

## Requirements

- ModelRuntime interface per [docs/contracts/runtime.md](../docs/contracts/runtime.md)
- StubProvider implementing InferenceProvider (canned responses)
- LoadModel, UnloadModel, Generate, SaveContext, RestoreContext
- Emit ModelLoaded/ModelUnloaded events via event bus

## Files

- `internal/runtime/runtime.go`
- `internal/runtime/interface.go`
- `internal/runtime/stub_provider.go`
- `internal/runtime/runtime_test.go`

## Acceptance Criteria

- [ ] Load/unload stub model succeeds
- [ ] Generate returns structured stub response
- [ ] No llama.cpp imports in runtime package
- [ ] Scheduler package can import runtime interface without inference imports

## Out Of Scope

- llama.cpp adapter (v0.6)
- KV manager integration
- Scheduler

## Related

- Spec: [docs/specs/v0.4-runtime.md](../docs/specs/v0.4-runtime.md)
- ADR-004

---
**Layer:** task
