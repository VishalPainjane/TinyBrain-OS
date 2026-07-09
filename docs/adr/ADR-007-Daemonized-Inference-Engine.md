# ADR-007: Daemonized Inference Engine

## Status

Accepted

## Context

The `tinybrain run` architecture artificially bottlenecks performance by performing a cold boot for every inference task. The 6-second execution time for a small 135M parameter model is dominated by Time-To-First-Token (TTFT) overhead, such as spawning the process, loading weights, allocating context, and CGO bridging, rather than true generation latency.

To build a high-throughput pipeline and accurately measure Tokens-Per-Second (TPS), the engine must isolate cold-boot latency from token generation latency.

## Decision

We will introduce a daemon execution model (`tinybrain daemon`) that initializes the kernel, scheduler, event bus, and runtime once, loading the inference engine into memory (utilizing `mmap` where applicable). 
Tasks will be submitted to the daemon via a new client command (`tinybrain submit`) using Structured JSON IPC over a local HTTP API.

1. **Daemon (`cmd/tinybrain/daemon.go`)**: Long-running process hosting the TinyBrain OS core. It exposes a local REST API for task submission.
2. **Client (`cmd/tinybrain/submit.go`)**: A thin client that submits JSON-structured tasks to the daemon and streams output to standard out.
3. **IPC Protocol**: Local HTTP. Chosen for simplicity, robust standard library support, cross-platform compatibility, and adherence to the "Structured JSON IPC" invariant.

## Consequences

### Positive

- Solves "The 21 TPS Problem" by entirely removing cold-boot overhead from the generation loop.
- Exposes true engine performance (hundreds of TPS on VRAM/DDR4).
- Paves the way for multi-task scheduling and concurrent agent execution against a single model in memory.

### Negative

- Increases CLI complexity (users must now manage a daemon process).
- Requires a local TCP port to be bound by the application.

### Trade-off

The performance gains of keeping the inference engine resident in memory far outweigh the added operational complexity of a daemon process.

---
**Layer:** decision
**Related:** [../architecture/kernel.md](../architecture/kernel.md), [../../tasks/022-system-integration.md](../../tasks/022-system-integration.md)
