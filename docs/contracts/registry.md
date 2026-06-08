# Contract: Registry

## Owner

Package: `internal/registry/` — agents (in-memory), models (in-memory + bbolt persistence); tools deferred

## Consumers

- Router — resolve agent for task type
- Scheduler — read agent resource profiles
- Runtime — resolve model definitions for load
- Tool layer — resolve tool definitions

## AgentDefinition

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique agent identifier |
| Name | string | Display name |
| ModelProfile | string | Reference to model definition |
| Tools | []string | Allowed tool IDs |
| ResourceProfile | ResourceProfile | Memory, priority limits |
| Priority | int | Default scheduling priority |

## ModelDefinition

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique model identifier |
| Path | string | GGUF file path |
| SizeBytes | uint64 | File size |
| MemoryBudget | uint64 | Estimated RAM+VRAM requirement |
| Capabilities | []string | Supported task types |
| Quantization | string | e.g., Q4_K_M |

## ToolDefinition

| Field | Type | Description |
|-------|------|-------------|
| ID | string | Unique tool identifier |
| Name | string | Display name |
| Handler | string | Tool execution handler reference |

## Implementation

### v0.3 — In-memory split registries

`AgentRegistry` and `ModelRegistry` (in-memory). A unified `Registry` facade is proposed in [RFC-004-Registry-Facade.md](../rfc/RFC-004-Registry-Facade.md).

`ToolDefinition` registration is not yet implemented.

### v0.5 — Model persistence

Model definitions may be stored in memory or in a bbolt file. Agent definitions remain in-memory only.

| Type | Constructors | Backend |
|------|--------------|---------|
| `ModelRegistry` | `NewModelRegistry()` | `InMemoryStore` |
| `ModelRegistry` | `NewBboltModelRegistry(dbPath, seedPath)` | `BboltStore` (`go.etcd.io/bbolt`) |
| `InMemoryStore` | `NewInMemoryStore()` | Mutex-protected map |
| `BboltStore` | `NewBboltStore(path)` | bbolt bucket `models`; key = model ID; value = gob-encoded `ModelDefinition` |

`LoadModelsYAML(path, store)` registers models from a seed file into any `ModelStore`. When used via `NewBboltModelRegistry`, seed loads **only if the store is empty**; no hot reload.

## Interface

```go
// Documentation only — unified facade; partial implementation today
type Registry interface {
    RegisterAgent(def AgentDefinition) error
    RegisterModel(def ModelDefinition) error
    RegisterTool(def ToolDefinition) error
    GetAgent(id string) (AgentDefinition, error)
    GetModel(id string) (ModelDefinition, error)
    GetTool(id string) (ToolDefinition, error)
    ListAgents() []AgentDefinition
    ListModels() []ModelDefinition
    ListTools() []ToolDefinition
}
```

### ModelStore (v0.5 — implemented)

```go
// ModelStore persists model definitions.
type ModelStore interface {
    RegisterModel(def ModelDefinition) error
    GetModel(id string) (ModelDefinition, error)
    ListModels() []ModelDefinition
    Close() error
}
```

Implementations: `InMemoryStore`, `BboltStore`.

### ModelRegistry (v0.5 — implemented)

```go
type ModelRegistry struct { /* delegates to ModelStore */ }

func NewModelRegistry() *ModelRegistry
func NewBboltModelRegistry(dbPath, seedPath string) (*ModelRegistry, error)

func (r *ModelRegistry) RegisterModel(def ModelDefinition) error
func (r *ModelRegistry) GetModel(id string) (ModelDefinition, error)
func (r *ModelRegistry) ListModels() []ModelDefinition
func (r *ModelRegistry) Close() error
```

`RegisterModel`, `GetModel`, and `ListModels` match v0.2/v0.3 semantics. `Close()` releases bbolt resources (no-op for in-memory).

### Errors

| Error | When |
|-------|------|
| `ErrDuplicateID` | `RegisterModel` with existing ID |
| `ErrNotFound` | `GetModel` with unknown ID |
| `fmt.Errorf("model ID is required")` | `RegisterModel` with empty ID |

## models.yaml seed format (v0.5)

Bootstrap file for empty bbolt stores. Field names map to `ModelDefinition`.

```yaml
models:
  - id: tinyllama-q4
    path: /models/tinyllama-q4.gguf
    size_bytes: 637534208
    memory_budget: 2147483648
    quantization: Q4_K_M
    capabilities:
      - chat
```

| YAML field | ModelDefinition field |
|------------|----------------------|
| `id` | `ID` (required) |
| `path` | `Path` |
| `size_bytes` | `SizeBytes` |
| `memory_budget` | `MemoryBudget` |
| `capabilities` | `Capabilities` |
| `quantization` | `Quantization` |

Duplicate IDs within the file return an error. See [testdata/models.yaml](../../testdata/models.yaml).

## Responsibilities

- Store and serve agent, model, tool definitions
- Single source of truth for capabilities (INV-004)

## Must NOT Own

- Runtime execution
- Scheduling
- Tool execution
- Hardware detection (reads profile, does not detect)

## Related

- [architecture/registry.md](../architecture/registry.md)
- ADR-002
- [tasks/005-agent-registry.md](../../tasks/005-agent-registry.md)
- [tasks/006-model-registry.md](../../tasks/006-model-registry.md)
- [tasks/006-registry-persistence.md](../../tasks/006-registry-persistence.md)
- [docs/specs/v0.5-model-registry.md](../specs/v0.5-model-registry.md)

---
**Layer:** contract
