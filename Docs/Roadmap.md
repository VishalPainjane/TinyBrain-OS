# TinyBrain Roadmap

TinyBrain is built in stages so the runtime foundation comes before advanced agent features.

## Phase 1 — Foundation

Goal: establish the core identity of TinyBrain as a local, hardware-aware runtime.

Deliverables:
- `README.md`
- `vision.md`
- `architecture.md`
- `docs/constitution.md`
- `docs/adr/ADR-001-Hardware-Aware-Runtime.md`
- `AGENTS.md`

Success criteria:
- project purpose is clear
- core architectural rules are documented
- future decisions can be checked against the constitution

## Phase 2 — Runtime Core

Goal: build the smallest working runtime.

Deliverables:
- model load and unload flow
- basic task submission flow
- hardware profile detection
- simple agent routing
- runtime status reporting

Success criteria:
- TinyBrain can start, select a model, run a task, and release resources

## Phase 3 — Scheduler

Goal: make TinyBrain behave like a real runtime, not a fixed workflow.

Deliverables:
- priority-based scheduling
- queue management
- preemption support
- resource-aware task handling

Success criteria:
- competing tasks are handled predictably
- higher priority work can interrupt lower priority work

## Phase 4 — Agent System

Goal: turn agents into modular plugins.

Deliverables:
- plugin-based agent registration
- standardized agent contracts
- agent lifecycle handling
- task-to-agent routing

Success criteria:
- agents can be added without changing core scheduler logic

## Phase 5 — Telemetry and Visibility

Goal: make the runtime observable.

Deliverables:
- metrics collection
- task and model state reporting
- runtime logs
- lightweight terminal dashboard

Success criteria:
- users can inspect system behavior while it runs

## Phase 6 — Memory and Swap Strategy

Goal: reduce resource waste under constrained hardware.

Deliverables:
- memory-aware model handling
- context persistence strategy
- model swapping policy
- resource cleanup rules

Success criteria:
- TinyBrain remains usable under limited RAM and VRAM

## Phase 7 — Expansion

Goal: grow TinyBrain carefully without breaking the core design.

Deliverables:
- more agent types
- richer tool support
- optional cloud connectors
- benchmark suite
- improved developer UX

Success criteria:
- the system stays local-first and hardware-aware even as features grow

## Guiding Rule

Every new feature must support the core idea:

TinyBrain should adapt to hardware, not force hardware to adapt to TinyBrain.
