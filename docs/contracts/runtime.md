# Contract: Runtime

## Owner

Package: `internal/runtime/` (not yet implemented)

## Consumers

- Scheduler — load/unload/generate commands
- Agent plugins — generate requests (via runtime API, not direct)

## ModelRuntime Interface

```go
// Documentation only
type ModelRuntime interface {
    LoadModel(modelID string) error
    UnloadModel(modelID string) error
    Generate(req GenerateRequest) (GenerateResponse, error)
    SaveContext(id string) error
    RestoreContext(id string) error
}
```

## InferenceProvider Interface

```go
// Documentation only — implemented by adapters, called by runtime
type InferenceProvider interface {
    LoadModel(modelID string) error
    UnloadModel(modelID string) error
    Generate(req GenerateRequest) (GenerateResponse, error)
    SaveContext(id string) error
    RestoreContext(id string) error
}
```

## GenerateRequest / GenerateResponse

Structured JSON — not natural language control messages.

## Responsibilities

- Model lifecycle: load, unload, warm, prefetch, evict
- Delegate inference to InferenceProvider
- Context save/restore coordination with KV manager
- Runtime metrics reporting

## Must NOT Own

- Queue management or scheduling policy
- Agent business logic
- Registry writes
- Direct llama.cpp imports outside adapter package (INV-008)

## Related

- [architecture/runtime.md](../architecture/runtime.md)
- ADR-004
- INV-001, INV-002, INV-008
- [tasks/008-runtime.md](../../tasks/008-runtime.md)
- [tasks/009-model-loader.md](../../tasks/009-model-loader.md)

---
**Layer:** contract
