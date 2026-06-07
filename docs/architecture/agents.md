# Agents

Agents are plugin workers — configured via the registry, not hardcoded in core.

## Responsibilities

- Execute tasks within their declared capability scope
- Request tools through the tool execution layer (never execute directly)
- Produce structured JSON output (not natural language for inter-agent communication)
- Operate within assigned memory quota and priority

## Agent Contract (conceptual)

Every agent plugin follows the same interface: receive context + task, return structured result. Specific agent names (planner, coder, reviewer) are **example fleet configurations** — not Go types in core.

## Example Fleet (sample config only)

```yaml
agents:
  planner:
    model: qwen-2b
    tools: [task_split]
  coder:
    model: deepseek-coder
    tools: [filesystem, git]
  reviewer:
    model: phi-3b
    tools: [static_analysis]
```

Users define their own fleets. TinyBrain core never imports `PlannerAgent` or similar.

## Tool Execution Layer

Separate from agents. Tools: search, filesystem, terminal, python, git. Agents emit tool requests as structured JSON; the tool layer executes and returns results.

## Structured Output

Agent outputs validated against JSON Schema / GBNF grammar for 100% parseability between components.

## Inputs

Task + context from runtime; tool results from tool layer.

## Outputs

Structured JSON results; tool requests; completion events.

## Dependencies (allowed)

Runtime API (via contract); registry (read own definition).

## Dependencies (forbidden)

Direct model access; scheduler; inference libraries.

## Future Plans

Agent install CLI; community agent packages; learned routing (v2+ router model).

## Non-Goals

Defining fixed agent roles in core; executing tools directly.

## Related Contracts

[registry.md](../contracts/registry.md), [runtime.md](../contracts/runtime.md)

## Related ADRs

ADR-002

---
**Layer:** architecture
**Source:** detail.md evolved architecture Part 5
**Related:** [registry.md](registry.md), [invariants.md](invariants.md)
