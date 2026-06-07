# Registry

The registry is the system of record for all configurable capabilities in TinyBrain OS.

## Responsibilities

- Store and serve agent definitions (name, model profile, tools, resource profile, priority)
- Store and serve model definitions (path, size, memory requirements, capabilities)
- Store and serve tool definitions (external actions agents may request)
- Enable discovery — runtime and scheduler query capabilities, never embed them

## What the Registry Owns

### Agent Definitions

Describe what an agent is allowed to do, its name, and how the runtime should treat it. Agents are configuration entries — not hardcoded Go types.

### Model Definitions

Available models, capabilities, resource requirements. Selection adapts to hardware profile (ADR-001).

### Tool Definitions

External actions (search, filesystem, terminal, git) that agents request through the tool execution layer.

## Core Rule

The registry **describes** capabilities. The scheduler and runtime **consume** them. Core never hardcodes specialized agent roles.

## Future: Capability Manager

Evolve toward CLI-driven fleet management:

```text
tinybrain agent install reviewer
tinybrain agent install sql-agent
```

## Inputs

Registration requests (config files, API, CLI); hardware profile for model filtering.

## Outputs

Resolved agent, model, and tool definitions for scheduler and runtime.

## Dependencies (allowed)

Hardware profile (read-only, for model filtering).

## Dependencies (forbidden)

Executing agents, loading models, scheduling.

## Future Plans

Hugging Face integration; persistent storage (BoltDB/SQLite — TBD v0.5); validation of definition schemas.

## Non-Goals

Runtime execution; scheduling decisions; tool execution.

## Related Contracts

[registry.md](../contracts/registry.md)

## Related ADRs

ADR-001, ADR-002

---
**Layer:** architecture
**Source:** detail.md Day 5, registry.md
**Related:** [agents.md](agents.md), [hardware.md](hardware.md)
