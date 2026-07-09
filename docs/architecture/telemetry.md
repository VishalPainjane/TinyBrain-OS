# Telemetry

Observability layer for TinyBrain OS. Most student projects skip this — TinyBrain should not.

## Responsibilities

- Collect metrics: VRAM, RAM, SSD, TPS, TTFT, queue depth, swap latency
- Emit events as metrics (TaskCreated → tasks_total; SwapFinished → swap_duration_ms)
- Feed brain-top TUI and future dashboard
- Support OpenTelemetry traces at V2

## Key Metrics

| Metric | Description |
|--------|-------------|
| agent_load_time | Model load duration (Time to pin weights) |
| queue_wait_time | Time in queue before execution starts |
| ttft | Time to First Token (measures startup/prefill overhead) |
| itl | Inter-Token Latency (measures true decode generation latency) |
| sustained_decode_throughput | True decoding tokens/sec during continuous batching |
| tail_latency | P95/P99 latency for both TTFT and ITL |
| prompt_length | Number of tokens in the prefill prompt |
| generated_length | Number of tokens produced in the decode phase |
| batch_size | Number of active sequences in the batch at each iteration |
| prefill_decode_split | Ratio of GPU time spent on prefill vs decode |
| cancellation_rate | Rate of sequences aborted before completion |
| oom_retry_count | Count of Out-Of-Memory errors and sequence retries |
| kv_cache_hit_rate | Ratio of prompt tokens avoiding recomputation |
| kv_cache_eviction_rate | Rate of KV block eviction to RAM swap |
| gpu_utilization | Real GPU compute saturation vs memory bandwidth wait |
| vram_usage / ram_usage | Paged memory block resource utilization |
| swap_time | Tier movement latency (`cudaMemcpyAsync`) |

## brain-top TUI (future)

Signature visibility feature — htop/btop for AI agents:

- Panel 1: Agent process states (RUNNING, WAITING, HIBERNATED, …)
- Panel 2: GPU/RAM/SSD utilization bars
- Panel 3: MLFQ queue depths (P0–P3)
- Panel 4: Swap monitor (VRAM → RAM → SSD movements)

Framework target: Bubble Tea (Go) / Charm Bracelet ecosystem.

## Inputs

Events from all subsystems; periodic resource probes.

## Outputs

Metric streams; log entries; TUI/dashboard data.

## Dependencies (allowed)

Read-only access to process table, scheduler state, runtime metrics.

## Dependencies (forbidden)

Controlling scheduler or runtime; UI logic in core packages.

## Future Plans

OpenTelemetry integration; Prometheus export; task timeline traces.

## Non-Goals

Scheduling; inference; agent logic.

---
**Layer:** architecture
**Source:** detail.md Part 2–3
**Related:** [overview.md](overview.md), [planning/benchmarks/target-metrics.md](../../planning/benchmarks/target-metrics.md)
