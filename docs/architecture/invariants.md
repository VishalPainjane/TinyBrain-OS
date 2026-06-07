# System Invariants

Things that must always remain true in TinyBrain OS. These are architectural constraints checkable in code review and future CI. Distinct from the [constitution](../constitution.md) (law) — invariants are specific, enforceable truths.

See [glossary.md](../glossary.md) for term definitions.

---

## INV-001: Scheduler never imports runtime

**Statement:** The `internal/scheduler/` package must not import `internal/runtime/` or any runtime subpackage.

**Enforcement:** Package import boundary; future CI lint rule.

**Related:** Constitution rule 3, [contracts/scheduler.md](../contracts/scheduler.md), [contracts/runtime.md](../contracts/runtime.md)

---

## INV-002: Runtime never imports scheduler

**Statement:** The `internal/runtime/` package must not import `internal/scheduler/` or any scheduler subpackage.

**Enforcement:** Package import boundary; future CI lint rule.

**Related:** Constitution rule 3, [architecture/runtime.md](runtime.md), [architecture/scheduler.md](scheduler.md)

---

## INV-003: Agents never directly invoke models

**Statement:** Agent plugins call the runtime API. They never call `InferenceProvider`, model loaders, or inference libraries directly.

**Enforcement:** Agent contract; code review; no inference imports in `internal/agents/`.

**Related:** Constitution rule 8, [architecture/agents.md](agents.md), [contracts/runtime.md](../contracts/runtime.md)

---

## INV-004: Registry owns discovery

**Statement:** Only the registry package registers, lists, and resolves agent, model, and tool definitions. Other packages read from the registry interface — they do not maintain parallel capability lists.

**Enforcement:** Single registry implementation; code review.

**Related:** [architecture/registry.md](registry.md), [contracts/registry.md](../contracts/registry.md)

---

## INV-005: Hardware determines model selection

**Statement:** Model assignment flows through hardware profile detection → capability classification → registry lookup. No hardcoded model paths or sizes in scheduler, runtime, or kernel code.

**Enforcement:** Code review; ADR-001 compliance check.

**Related:** ADR-001, [architecture/hardware.md](hardware.md), [glossary.md](../glossary.md#hardware-profile)

---

## INV-006: Structured JSON IPC only

**Statement:** Inter-component and inter-agent messages are structured JSON. Natural-language strings are not used as control messages between system components.

**Enforcement:** Event type definitions; schema validation at boundaries.

**Related:** Constitution rule 7, ADR-003

---

## INV-007: Core never hardcodes agent roles

**Statement:** No Go types, constants, or switch statements in `internal/` core packages encode fixed agent roles (Planner, Browser, Coder, Reasoner). These names may appear only in test fixtures, sample configs, and documentation examples.

**Enforcement:** Code review; grep audit for role names in `internal/` excluding tests and examples.

**Related:** ADR-002, [planning/decisions/rejected.md](../../planning/decisions/rejected.md)

---

## INV-008: Inference engine accessed only via Provider interface

**Statement:** llama.cpp, Ollama, vLLM, and cloud SDK imports exist only in adapter packages implementing `InferenceProvider`. No other package imports inference engine libraries.

**Enforcement:** Package import boundary; adapter isolation.

**Related:** ADR-004, [contracts/runtime.md](../contracts/runtime.md), [glossary.md](../glossary.md#provider)

---
**Layer:** invariants
**Source:** constitution, detail.md Part 5
**Last verified against code:** 2026-06-07
**Related:** [constitution.md](../constitution.md), [glossary.md](../glossary.md)
