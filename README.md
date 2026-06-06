# TinyBrain OS

## What is TinyBrain?

TinyBrain OS is a local-first, resource-aware runtime for orchestrating specialized LLM agents on constrained hardware. It treats agents like operating-system processes and focuses on scheduling, memory management, model loading, and structured inter-agent communication rather than a conventional chatbot workflow.

## Why does it exist?

Most AI systems assume large cloud hardware or a single monolithic model. TinyBrain exists to explore a different idea: smaller specialized models can be orchestrated intelligently to behave like a lightweight operating system for AI work. The goal is to run useful agent workflows on consumer hardware while staying offline-capable and resource-efficient.

## High-level architecture

```
User
  ↓
Proxy Router
  ↓
Orchestrator
  ↓
Scheduler (MLFQ)
  ↓
Agent Runtime
  ↓
Local Models (GGUF / llama.cpp)
```

Core principles:
- Local-first and edge-first
- Dynamic model swapping instead of keeping everything loaded
- Structured JSON IPC between agents
- Preemption and priority-based scheduling
- KV cache persistence as a first-class runtime primitive
- Memory mapping and VRAM/RAM/NVMe tiering for constrained devices

## Current roadmap

**Day 1**
- Write the project README
- Lock the core thesis and scope

**MVP**
- Single model load/unload flow
- Basic model registry
- Minimal runtime wrapper around `llama.cpp`

**Next**
- Planner + Coder agent loop
- Scheduler with priority queues
- Process table and live state tracking

**Later**
- KV cache save/restore
- brain-top TUI
- Kubernetes operator and CRDs
- Benchmarks and observability

## Current status

TinyBrain OS is in the design and planning phase. The core architecture, agent roles, runtime direction, and scheduling philosophy are defined, but the implementation is still intentionally small for the first milestone. The immediate goal is to build a clean foundation before adding more agents, more tooling, or more infrastructure.

---

*TinyBrain OS is being built as a systems project, not a chatbot project.*
