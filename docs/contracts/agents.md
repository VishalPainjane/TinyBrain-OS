# Contract: Agents

## Owner

Package: `internal/agents/` — plugin contract, executor, sample reference plugin

## Consumers

- Runtime integration (future task pipeline)
- Registry (agent definitions resolve to plugin factories)

## Agent Interface

```go
type Agent interface {
    ID() string
    Execute(ctx ExecuteContext, req TaskRequest) (TaskResult, error)
}
```

## RuntimeAPI

Agents call inference only through `RuntimeAPI.Generate` (INV-003). No `InferenceProvider` or `internal/inference` imports.

## TaskRequest / TaskResult

Structured JSON IPC — not natural language between components.

| Field | TaskRequest | TaskResult |
|-------|-------------|------------|
| TaskID | string | string |
| PID | string | — |
| Input | string | — |
| AgentID | — | string |
| Output | — | JSON string |

## Executor

`Executor.Run` publishes `AgentStarted` and `AgentStopped` on the event bus while delegating to `Agent.Execute`.

## Responsibilities

- Define generic plugin contract (ADR-002)
- Execute tasks via runtime API
- Emit agent lifecycle events

## Must NOT Own

- Registry writes
- Scheduling
- Direct model/inference access
- Hardcoded agent types (Planner, Coder, etc.)

## Related

- [architecture/agents.md](../architecture/agents.md)
- ADR-002, INV-003
- [tasks/014-agent-plugin.md](../../tasks/014-agent-plugin.md)

---
**Layer:** contract
