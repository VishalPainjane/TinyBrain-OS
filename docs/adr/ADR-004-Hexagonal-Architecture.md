# ADR-004: Hexagonal Architecture

## Status

Accepted

## Context

TinyBrain must support multiple inference backends (llama.cpp today, vLLM or custom CUDA engine tomorrow) and optional cloud providers (OpenAI, Anthropic) without rewriting the scheduler or runtime core.

Hardcoding llama.cpp into the runtime or scheduler would prevent engine swaps and cloud fallback.

## Decision

Introduce `InferenceProvider` as a port (interface). The runtime calls providers; providers implement engine-specific logic.

```go
// Documentation only — not compiled code
type InferenceProvider interface {
    LoadModel(modelID string) error
    UnloadModel(modelID string) error
    Generate(req GenerateRequest) (GenerateResponse, error)
    SaveContext(id string) error
    RestoreContext(id string) error
}
```

Implementations:
- `LlamaCppProvider` — local GGUF via llama.cpp (primary)
- `StubProvider` — no-op for testing (v0.4)
- Future: `OpenAIProvider`, `VLLMProvider`, etc.

The scheduler never references any provider or engine.

## Consequences

### Positive

- Engine swap without scheduler/runtime rewrite.
- Cloud optional via adapter — local-first preserved.
- Testable with stub provider before llama.cpp integration.
- Aligns with INV-008.

### Negative

- Interface design must be stable early; breaking changes are costly.
- Adapter layer adds indirection.

### Trade-off

Design the interface carefully upfront; gain long-term engine independence.

---
**Layer:** decision
**Related:** [../architecture/runtime.md](../architecture/runtime.md), [../contracts/runtime.md](../contracts/runtime.md), [../../planning/decisions/accepted.md](../../planning/decisions/accepted.md)
