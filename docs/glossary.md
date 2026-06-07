# TinyBrain OS Glossary

Canonical definitions for TinyBrain OS terms. Use these definitions consistently across all documentation and code comments. If a new term is introduced, add it here before using it elsewhere.

See also: [architecture/invariants.md](architecture/invariants.md)

---

## Agent

**Definition:** A configurable plugin that represents a capability in the system. An agent is defined by name, model profile, tool set, resource profile, and priority — not by a hardcoded type in core code. Example fleet entries (planner, coder) are sample configurations, not architectural types.

**Owner:** Registry + plugin layer (`internal/registry/`, future `internal/agents/`)

**Related:** [Capability](#capability), [Process](#process), [Registry](#registry), [Tool](#tool), [architecture/agents.md](architecture/agents.md)

---

## Process

**Definition:** An operating-system-style execution unit representing a running or scheduled agent instance. Every process has a unique identifier, lifecycle state, priority, resource usage, and optional KV cache reference. The scheduler operates on processes, not on agent definitions directly.

**Owner:** Kernel (`internal/process/`)

**Related:** [ProcessState](#processstate), [Task](#task), [Preemption](#preemption), [Hibernation](#hibernation), [architecture/kernel.md](architecture/kernel.md)

---

## ProcessState

**Definition:** The lifecycle state of a process. Valid states: NEW, READY, RUNNING, WAITING, PREEMPTED, HIBERNATED, TERMINATED. Every process is always in exactly one state.

**Owner:** Kernel (`internal/process/`)

**Related:** [Process](#process), [Preemption](#preemption), [Hibernation](#hibernation), [contracts/process.md](contracts/process.md)

---

## Task

**Definition:** A unit of work submitted to TinyBrain by a user or API. Tasks are classified, routed to an agent capability, scheduled, and executed. A task may spawn one or more processes during its lifecycle.

**Owner:** Scheduler input / API layer (future `internal/api/`)

**Related:** [Process](#process), [Agent](#agent), [Scheduler](#scheduler)

---

## Scheduler

**Definition:** The control component that decides execution order, priority, preemption, and fairness. The scheduler schedules VRAM and process execution — it never talks to inference engines directly. It delegates model operations to the runtime.

**Owner:** `internal/scheduler/` (not yet implemented)

**Related:** [Runtime](#runtime), [Preemption](#preemption), [Process](#process), [architecture/scheduler.md](architecture/scheduler.md)

---

## Runtime

**Definition:** The subsystem that manages model lifecycle: load, unload, warm, prefetch, generate, and context save/restore. The runtime is the only component that invokes inference providers. It does not schedule work or define agent capabilities.

**Owner:** `internal/runtime/` (not yet implemented)

**Related:** [Provider](#provider), [Model](#model), [KV Cache](#kv-cache), [architecture/runtime.md](architecture/runtime.md)

---

## Registry

**Definition:** The system of record for all configurable capabilities: agent definitions, model definitions, and tool definitions. The registry describes what exists; the scheduler and runtime consume those descriptions. Core code never hardcodes specialized agent roles.

**Owner:** `internal/registry/` (not yet implemented)

**Related:** [Agent](#agent), [Model](#model), [Tool](#tool), [Capability](#capability), [architecture/registry.md](architecture/registry.md)

---

## Capability

**Definition:** An abstract skill or function that an agent can perform. Capabilities are declared in agent definitions and discovered via the registry. A capability is not a hardcoded class (e.g., "Planner" is a capability name in a config, not a Go type).

**Owner:** Registry

**Related:** [Agent](#agent), [Registry](#registry), [Tool](#tool)

---

## Tool

**Definition:** An external action or ability that agents may request through the runtime (search, filesystem, terminal, git, etc.). Agents never execute tools directly; a separate tool execution layer handles requests on their behalf.

**Owner:** Tool execution layer (future `internal/tools/`)

**Related:** [Agent](#agent), [Registry](#registry), [Process](#process)

---

## Model

**Definition:** A loadable inference artifact, typically GGUF format with Q4_K_M quantization. Models are described in the model registry with path, size, memory requirements, and capability metadata. Model selection is determined by hardware profile, not fixed per agent role.

**Owner:** Model registry (`internal/registry/`, future `internal/loader/`)

**Related:** [Provider](#provider), [Hardware Profile](#hardware-profile), [Runtime](#runtime)

---

## Provider

**Definition:** An inference adapter implementing the `InferenceProvider` interface. Local GGUF via llama.cpp is the primary provider; optional cloud providers (OpenAI, Anthropic, vLLM) implement the same interface. The scheduler and agents never reference providers directly.

**Owner:** Inference adapter layer (future `internal/inference/`)

**Related:** [Runtime](#runtime), [Model](#model), ADR-004

---

## Hardware Profile

**Definition:** A classification of the local execution environment based on detected RAM, VRAM, CPU, and inference backend. Profiles (Tiny, Standard, Workstation) guide model selection, scheduling aggressiveness, and resource planning. Not just labels — active inputs to runtime decisions.

**Owner:** `internal/hardware/` (not yet implemented)

**Related:** [Model](#model), [Registry](#registry), [architecture/hardware.md](architecture/hardware.md)

---

## KV Cache

**Definition:** The key-value attention state produced during model inference. In TinyBrain, KV cache is a first-class persistence primitive — not an implementation detail. KV caches can be saved, compressed, moved across memory tiers, and restored to avoid recomputing prompts.

**Owner:** Memory / KV manager (future `internal/kv/`, `internal/memory/`)

**Related:** [Hibernation](#hibernation), [Runtime](#runtime), [architecture/memory.md](architecture/memory.md)

---

## Hibernation

**Definition:** A process state (HIBERNATED) where execution is suspended, model weights are unloaded, but KV cache is preserved in a lower memory tier (RAM or NVMe). Enables context preservation without holding VRAM. Also called context hibernation.

**Owner:** Memory layer + process state machine

**Related:** [KV Cache](#kv-cache), [ProcessState](#processstate), [Process](#process)

---

## Preemption

**Definition:** Interrupting a running process to allow higher-priority work to execute. The preempted process transitions to PREEMPTED state; its KV cache may be preserved. Analogous to OS process preemption. The scheduler owns preemption decisions.

**Owner:** Scheduler (`internal/scheduler/`)

**Related:** [ProcessState](#processstate), [Scheduler](#scheduler), [Process](#process)

---
**Layer:** glossary
**Source:** detail.md Parts 1–5, implemented code
**Last verified against code:** 2026-06-07
**Related:** [constitution.md](constitution.md), [architecture/invariants.md](architecture/invariants.md)
