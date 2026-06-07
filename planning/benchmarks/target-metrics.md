# Target Metrics

Performance targets (aspirational until benchmark suite exists). Distinct from [../metrics/](../metrics/) which tracks delivery velocity and version progress.

| Metric | Target | Measured at |
|--------|--------|-------------|
| Kernel startup (no models loaded) | < 2 s | v0.1 |
| Agent registration | < 100 ms | v0.2 |
| Process lookup by PID | O(1) | v0.1 task 002 |
| Runtime memory (kernel only) | < 200 MB | v0.1 |
| Model swap latency (warm) | < 5 s | v0.4 |
| TTFT (single 2B model, Standard profile) | < 3 s | v0.6 |
| TPS (Standard profile) | > 20 tok/s | v0.6 |
| Peak VRAM (2-agent demo) | < 3.6 GB (90% of 4 GB) | v0.8 |
| Swap latency (KV tier move) | < 500 ms | v1.0 |
| Queue wait time (P0 task) | < 100 ms | v0.7 |

---
**Layer:** planning
**Related:** [../../docs/specs/v1.0.md](../../docs/specs/v1.0.md)
