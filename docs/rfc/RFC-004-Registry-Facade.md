# RFC-004: Registry Facade

## Status

Proposal

## Problem

Month 1 ships separate `AgentRegistry` and `ModelRegistry` implementations. The registry contract documents a unified `Registry` interface with agent, model, and tool operations. Runtime and scheduler consumers will need a single entry point for capability discovery.

## Proposal

Introduce a thin `Registry` facade that composes agent, model, and (future) tool registries behind the contract interface. No change to underlying in-memory stores.

## Alternatives

- Keep split registries and inject both into runtime (more wiring, violates contract ergonomics)
- Monolithic registry struct (larger single type, harder to test in isolation)

## Open Questions

- When should tool registry land — v0.2 completion or v0.5?
- Does facade live in `internal/registry/facade.go` or `internal/registry/registry.go`?

## Not Scheduled Until

Task explicitly assigned (likely v0.4 runtime or v0.5 persistence work).

Spec version: v0.4+

---
**Layer:** planning
**Related:** [../contracts/registry.md](../contracts/registry.md), [../../planning/releases/v0.2.md](../../planning/releases/v0.2.md)
