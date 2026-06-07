# Contract: Registry

## Owner

Package: `internal/registry/` — agents and models implemented; tools deferred

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

## Implementation (v0.3)

Month 1 ships split in-memory registries: `AgentRegistry` and `ModelRegistry`. A unified `Registry` facade is proposed in [RFC-004-Registry-Facade.md](../rfc/RFC-004-Registry-Facade.md).

`ToolDefinition` registration is not yet implemented.

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

---
**Layer:** contract
