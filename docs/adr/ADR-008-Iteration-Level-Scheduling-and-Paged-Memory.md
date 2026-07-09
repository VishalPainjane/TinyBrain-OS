# ADR-008: Iteration-Level Scheduling and Paged Memory

## Status

Accepted

> **TinyBrain OS schedules sequences, not models; keeps weights resident; and manages prefill, decode, and KV cache as separate resource classes.**

## Context

The initial architectural concept of TinyBrain OS applied classical OS process models to LLM orchestration, treating an entire LLM inference operation (including its loaded weights) as a schedulable "process" that could be context-switched in and out of memory. This approach created a massive I/O bottleneck and performance contradiction (e.g., spending 6 seconds for 128 tokens of generation, where almost all time was consumed by Time-To-First-Token/model loading overhead).

When scaling to concurrent agents, swapping whole gigabyte-scale models in and out of VRAM across the PCIe bus destroys true concurrency and violates low latency requirements. 

Modern high-throughput inference engines (like vLLM, TensorRT-LLM, ORCA, and Hugging Face TGI) have shown that to achieve performance, the system must decouple the static model weights from the mutable execution state (the KV cache). The true unit of an LLM "process" is the *sequence* (prompt + generated tokens) and its associated KV Cache, not the model itself.

## Decision

We will transition the TinyBrain OS architecture (v2) from Macro-Process Scheduling to Micro-Iteration Scheduling by adopting the following paradigms:

1. **Resident Model Workers (Data Plane)**
   Model weights will be treated as immutable execution hardware (a persistent worker / resident model service) loaded once into VRAM. Spawning, loading, and swapping whole models per request is prohibited unless explicitly handling admission limits.

2. **Sequence as a Process (Control Plane)**
   The core abstraction of an OS "Process" is redefined. A Process in TinyBrain OS is no longer the model executable; it is the mutable execution context (the KV Cache) and the decoding progress of an agent's request sequence.

3. **Iteration-Level Scheduling (In-Flight Batching)**
   The scheduler operates at the iteration granularity (output-token boundary), not the request granularity. We frame preemption at the token boundary, not as a CPU-style process switch. After each forward pass, the scheduler chooses which sequence advances next, checking the queue to eject completed sequences and admit new ones mid-flight, enabling stall-free continuous batching.

4. **Paged KV-Cache Memory Allocator (Virtual Memory)**
   To handle concurrent contexts without severe memory fragmentation, VRAM will be managed using fixed-size blocks/pages managed by a page table (the exact block size and reservation ratios are left to backend configuration). 
   If VRAM pressure occurs, the kernel can context-switch by preempting an active sequence, evicting/swapping its KV blocks to host RAM via PCIe, and restoring them later, without ever moving the model weights.

5. **Prefill and Decode as Separate First-Class Lanes**
   Prefill (processing the prompt) and decode (generating tokens) are treated as distinct workloads with different compute/memory profiles. Prompts will be processed using chunked prefills to avoid stalling ongoing decodes. Prefill/decode aggregation vs. disaggregation is a tradeoff depending on the SLO mix; our baseline is chunked prefill co-batched with decode steps, with full disaggregation as a future configuration.

### TinyBrain OS v2 Architectural Specification

To implement this, the runtime is decoupled into a clear Control Plane and Data Plane:

```text
                  +---------------------------------------+
                  |             CONTROL PLANE             |
                  |  [Request Ingest] -> [MLFQ Queue]     |
                  +-------------------+-------------------+
                                      | (Token Boundary Decisions)
                                      v
                  +---------------------------------------+
                  |              DATA PLANE               |
                  |  +---------------------------------+  |
                  |  | Persistent Inference Daemon     |  |
                  |  | (Weights Pinned in VRAM)        |  |
                  |  +---------------------------------+  |
                  |  | Paged KV-Cache Memory Allocator |  |
                  |  | [VRAM Blocks] <-> [Host RAM]    |  |
                  |  +---------------------------------+  |
                  +---------------------------------------+
```

## Consequences

### Positive
- Vastly increases concurrent agent capacity and total system throughput.
- Eliminates VRAM fragmentation via block allocation.
- True continuous batching and context switching become lightweight, taking milliseconds instead of seconds.
- Accurately tracks independent metrics: Time-To-First-Token (TTFT) and Inter-Token Latency (ITL).

### Negative
- Significantly increases the complexity of the Data Plane (CGO/CUDA backend) due to custom block allocators and iteration loops.
- Demands tight coupling between the Go scheduler's control plane and the backend's sequence execution state.

### Implementation Roadmap
1. Build one persistent resident model worker per GPU.
2. Implement prompt templating in the Router before dispatch.
3. Implement continuous batching for requests at the token/iteration level.
4. Optimize for the right metrics: explicitly track TTFT (Time-To-First-Token), ITL (Inter-Token Latency), sustained decode throughput, queue wait time, model load time, KV cache hit/eviction rate, GPU utilization, prompt length, generated length, batch size at each iteration, prefill/decode split, cancellation rate, OOM/retry count, and tail latency.
5. Introduce the KV-cache block manager with a clear Memory Pressure Policy (Admit, Queue, Spill, Reject).
6. (Later) Add chunked prefill strategies to separate prefill/decode lanes.

## Scope and Non-Goals

- **No Full Model Hot-Swap Per Request**: TinyBrain OS v2 does not promise full model hot-swapping for every inference request. 
- **Not CPU RAM**: We do not treat GPU memory like CPU RAM. The core promise is: resident model + sequence scheduling + KV-aware memory management.
- **Future Scale-Out**: Multi-GPU, multi-node parallelism, and deep prefill/decode disaggregation (passing KV cache across the network) are explicitly documented as "future scale-out modes", not part of the core v2 single-node mental model.

---
**Layer:** decision
**Related:** [../architecture/kernel.md](../architecture/kernel.md), [../architecture/scheduler.md](../architecture/scheduler.md), [../architecture/memory.md](../architecture/memory.md)
