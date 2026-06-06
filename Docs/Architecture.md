# TinyBrain Architecture

TinyBrain works like a local AI runtime, not like a single chatbot. Its job is to route work, manage models, and keep resource usage under control on constrained hardware.

## High-level flow

1. A user submits a task.
2. TinyBrain classifies the task and decides which agent should handle it.
3. The scheduler assigns priority and decides when the task should run.
4. The runtime loads the required model only when needed.
5. The agent produces output.
6. TinyBrain stores state, updates metrics, and frees resources when the task is done.

## Main layers

### 1. Interface layer
This is the entry point for users. It accepts requests, returns responses, and exposes system status.

### 2. Router layer
This layer decides what kind of work the request represents. For example, it can route a task toward planning, coding, browsing, or reasoning.

### 3. Scheduler layer
This is the control center. It decides order, priority, fairness, and preemption. Its purpose is to keep the system responsive even when multiple tasks compete for limited hardware.

### 4. Runtime layer
This layer manages model lifecycle. It loads models, unloads them, swaps them when needed, and keeps the active footprint small.

### 5. Agent layer
Agents are specialized workers. Each agent is responsible for a specific type of task rather than trying to do everything.

### 6. Memory layer
TinyBrain treats memory as a managed resource. It keeps active state available, preserves important context, and releases what is no longer needed.

### 7. Telemetry layer
This layer tracks what the system is doing: resource usage, task flow, queue depth, model activity, and performance behavior.

## Core design principles

TinyBrain is built around a few ideas:

- local first
- hardware aware
- dynamic instead of static
- specialized instead of monolithic
- resource disciplined instead of wasteful

## How the pieces fit together

TinyBrain does not run everything at once. It chooses the smallest useful set of components for the task, activates them when needed, and shuts them down when they are idle. That allows the system to behave more like an operating system for AI workloads than a fixed model wrapper.

## What this architecture is trying to solve

The architecture exists to make local AI feel practical on real machines. It aims to reduce wasted memory, improve responsiveness, and make multiple specialized models cooperate under one control plane.

In short: TinyBrain is designed to coordinate AI work the way an OS coordinates processes.