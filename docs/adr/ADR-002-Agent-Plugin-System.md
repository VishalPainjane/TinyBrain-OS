# ADR-002: Agent Plugin System

## Status

Accepted

## Context

Early designs embedded fixed agent types (Planner, Browser, Coder, Reasoner) as core classes. This couples the runtime to a specific workflow, prevents user-defined agent fleets, and violates the OS analogy — Linux does not know Chrome or VS Code exist; it only knows processes, memory, and priority.

The evolved architecture treats agents as registry-defined plugins with configurable model profiles, tool sets, and resource profiles.

## Decision

Agents are plugins registered via the registry. Core code defines an agent **contract** (interface), not agent **types** (classes).

- Agent definitions live in configuration (YAML, registry API, CLI install).
- Names like "planner" or "coder" are capability labels in config — not Go types in `internal/`.
- Adding a new agent requires registry entry + plugin implementation — never core code changes.

## Consequences

### Positive

- User-defined agent fleets without forking core.
- Core stays small and stable as agents proliferate.
- Aligns with constitution rule: agents are plugins.
- Enables hardware-adaptive model assignment per agent capability.

### Negative

- More indirection than hardcoded agents for early demos.
- Requires registry and plugin loading infrastructure before first agent demo.

### Trade-off

Sacrifice demo simplicity for long-term modularity and systems credibility.

---
**Layer:** decision
**Related:** [../architecture/agents.md](../architecture/agents.md), [../architecture/registry.md](../architecture/registry.md), [../../planning/decisions/rejected.md](../../planning/decisions/rejected.md)
