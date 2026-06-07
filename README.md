# TinyBrain OS

## Vision

TinyBrain OS is a hardware-aware AI runtime kernel for dynamically orchestrating specialized agents on constrained local hardware. Instead of one monolithic model handling every task, TinyBrain treats agents as scheduled processes — loaded when needed, swapped when idle, and coordinated under strict VRAM and RAM budgets.

Local AI fails on consumer hardware when treated as a single-model problem. TinyBrain exists to make it practical: dynamic capability allocation, not static model-centric workflows.

See also: [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)

## Architecture Diagram

```text
User
  ↓
API
  ↓
Router
  ↓
Scheduler
  ↓
Runtime
  ↓
InferenceProvider
  ↓
Models (GGUF / llama.cpp)
```

Core principles:

- Local-first, hardware-aware
- Agents are plugins (registry-defined, not hardcoded)
- Structured JSON IPC between components
- Scheduler independent from inference engine
- Dynamic model swapping under resource budget

Full architecture: [docs/architecture/overview.md](docs/architecture/overview.md)

## Quick Start

**Current state:** v0.4 Runtime shipped at tag `v0.4`. Month 2 active — v0.5 model registry next.

```bash
go test ./...
```

Future commands (not yet implemented):

```bash
tinybrain run "your task"
brain-top
```

## Current Version

**V0.5 Model Registry** — [docs/current.md](docs/current.md) | Sprint: [planning/execution/current-sprint.md](planning/execution/current-sprint.md)

## Roadmap

| Version | Goal | Status |
|---------|------|--------|
| v0.1 | Kernel — process model + table | Shipped |
| v0.2 | Registry — agent/model definitions + events | Shipped |
| v0.3 | Hardware profiler | Shipped |
| v0.4 | Runtime shell + stub provider | Shipped |
| v0.5 | Persistent model registry | In Progress |
| v0.6 | llama.cpp inference | Planned |
| v0.7 | MLFQ scheduler | Planned |
| v0.8 | Plugin agents | Planned |
| v1.0 | Integrated runtime + brain-top | Planned |

Details: [docs/specs/](docs/specs/) and [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)

---

*TinyBrain OS is a systems project — a local operating system for AI agents, not a chatbot wrapper.*
