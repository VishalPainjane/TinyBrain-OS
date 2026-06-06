# ADR-001: Hardware-Aware Runtime

## Status

Accepted

## Context

Most local AI systems assume a fixed hardware environment and a fixed set of models.

For example, a system may always assign a 2B model to planning, a 7B model to coding, and a larger model to reasoning. This approach works on specific hardware configurations but does not scale across different devices.

Users may run TinyBrain on a wide range of systems, from laptops with limited RAM and VRAM to high-end workstations with abundant resources. A fixed model strategy would either underutilize powerful hardware or fail on constrained hardware.

TinyBrain aims to be a general-purpose local AI runtime rather than a workflow built around specific models.

## Decision

TinyBrain will adapt model selection based on detected hardware capabilities instead of using fixed model sizes.

At startup, TinyBrain will evaluate available system resources such as CPU, RAM, VRAM, and inference backends.

Based on these capabilities, the runtime will select the most appropriate models for available agents and workloads.

Agents represent capabilities, while model assignments remain dynamic.

The runtime may choose different model fleets on different machines without requiring code changes.

## Consequences

### Positive

* Supports a wider range of hardware profiles.
* Improves resource utilization.
* Enables the same architecture to run on laptops and workstations.
* Decouples agent capabilities from specific model implementations.
* Makes TinyBrain more future-proof as new models become available.

### Negative

* Introduces additional runtime complexity.
* Requires hardware detection and capability profiling.
* Benchmarking becomes more difficult because configurations may differ between machines.

### Trade-off

TinyBrain sacrifices some simplicity in exchange for portability, flexibility, and better resource efficiency.

This decision aligns with the project's vision of being a hardware-aware AI runtime rather than a fixed multi-agent workflow.
