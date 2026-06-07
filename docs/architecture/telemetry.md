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
| agent_load_time | Model load duration |
| agent_unload_time | Model unload duration |
| swap_time | Tier movement latency |
| queue_wait_time | Time in queue before execution |
| vram_usage / ram_usage | Resource utilization |
| tokens_per_second | Inference throughput |
| ttft | Time to first token |

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
