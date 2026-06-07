# TinyBrain OS Constitution

This document is the law of TinyBrain OS. Every feature, refactor, and AI-generated suggestion must comply with these rules before acceptance. See the [File Ownership Model](../planning/README.md#file-ownership-model) for how this document relates to other layers.

## Preamble

TinyBrain OS is a hardware-aware AI runtime kernel for dynamically orchestrating arbitrary AI agents under resource constraints on local hardware. It is a runtime system, not a chatbot wrapper.

## Non-Negotiable Rules

1. **Hardware determines model selection.** The runtime adapts model fleets to detected capabilities. Fixed model sizes per agent role are forbidden in core logic. *(Rationale: portability across laptops and workstations.)*

2. **Agents are plugins.** Agent capabilities are defined in the registry and loaded as configuration. Core code never hardcodes agent types such as Planner, Coder, Browser, or Reasoner. *(Rationale: user-defined fleets, not fixed workflows.)*

3. **The scheduler never depends on the inference engine.** The scheduler delegates all model operations to the runtime. It never imports llama.cpp, Ollama, vLLM, or any inference library. *(Rationale: engine swap without scheduler rewrite.)*

4. **The runtime never depends on the UI.** Dashboard, TUI, and API are consumers of runtime telemetry — not drivers of runtime behavior. *(Rationale: headless operation, testability.)*

5. **Cloud is optional.** Local-first, edge-first, offline-capable operation is the default. Cloud providers are adapters behind `InferenceProvider`, not required dependencies. *(Rationale: constrained hardware target.)*

6. **Local-first by default.** Core operation must not require OpenAI, Anthropic, Gemini, or any external API. *(Rationale: project thesis.)*

7. **Structured JSON IPC only.** Components communicate via structured JSON messages. Natural-language messages between agents are forbidden. *(Rationale: parseability, reduced agent drift.)*

8. **Registry describes capabilities; core consumes them.** The registry owns agent, model, and tool definitions. The scheduler and runtime read from the registry — they do not embed capability knowledge. *(Rationale: modularity.)*

## Boundary Matrix

| From | May call | Must NOT call |
|------|----------|---------------|
| API / Router | Scheduler, Registry | Runtime, InferenceProvider, llama.cpp |
| Scheduler | Runtime, Process table, Event bus | InferenceProvider, llama.cpp, Registry writes |
| Runtime | InferenceProvider, Loader, KV manager | Scheduler, UI |
| Registry | (none — passive data store) | Runtime, Scheduler |
| Agents (plugins) | Runtime (via contract) | Models, Providers, Scheduler |
| InferenceProvider | llama.cpp / cloud APIs | Scheduler, Process table |
| UI / brain-top | Telemetry, read-only APIs | Scheduler, Runtime control |

## Decision Policy

When documents conflict, resolve in this order:

1. Implemented code (code-first truth)
2. [current.md](current.md)
3. This constitution
4. [adr/](adr/)
5. [specs/](specs/)
6. [tasks/](../tasks/)
7. [planning/decisions/accepted.md](../planning/decisions/accepted.md)

If a proposal conflicts with this constitution, reject it or revise the proposal. If the constitution itself must change, create a new ADR and amend this document explicitly.

## Design Priority

When trade-offs appear, favor:

- local execution over cloud dependency
- hardware awareness over fixed assumptions
- modular agents over tightly coupled systems
- runtime stability over UI convenience
- clear boundaries over hidden coupling
- simple correct MVP over premature optimization

## Amendment Process

1. Identify the rule that needs change and document why in a new ADR.
2. Get ADR accepted before changing this file.
3. Update [architecture/invariants.md](architecture/invariants.md) if boundaries shift.
4. Update [glossary.md](glossary.md) if terms change.
5. Log the change in [planning/decisions/accepted.md](../planning/decisions/accepted.md).

---
**Layer:** law
**Source:** docs/constitution.md, detail.md Parts 1 and 5
**Last verified against code:** 2026-06-07
**Related:** [architecture/invariants.md](architecture/invariants.md), [glossary.md](glossary.md), [adr/](adr/)
