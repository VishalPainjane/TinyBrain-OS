# Assumptions

Beliefs that may be wrong. Separate from [decisions/accepted.md](decisions/accepted.md) — decisions are choices made; assumptions are beliefs awaiting validation.

Status values: `Active` | `Needs validation` | `Invalidated` | `Retired`

| Assumption | Status | Validation |
|------------|--------|------------|
| Project targets local-first execution | Active | Confirmed by ADR-005 and constitution |
| Most users have ≤16 GB RAM | Active | Baseline Standard profile (16 GB + RTX 3050) |
| llama.cpp remains primary inference engine | Active | detail.md Part 1; revisit if vLLM matures |
| 4 GB VRAM is sufficient for 2-agent demo fleet | Needs validation | Benchmark at v0.8 on Standard profile |
| Go channels sufficient for event bus at MVP scale | Active | Revisit if event rate exceeds ~100/sec |
| Users will define agent fleets via registry | Active | Validated by ADR-002 design; demo at v0.8 |
| mmap-backed GGUF loading works reliably on Windows | Needs validation | Test at v0.6 inference integration |
| FIFO scheduler skeleton won't block MLFQ migration | Active | migration-paths.md documents upgrade path |
| Single developer can reach v1.0 in ~6 months | Active | Track via metrics/velocity.md |
| Windows VRAM detectable without nvidia-smi | Needs validation | v0.3 returns 0 VRAM when nvidia-smi absent; Standard profile requires GPU metrics |

When an assumption is invalidated: update status here, log in [retrospectives/](retrospectives/), consider new ADR.

---
**Layer:** planning
**Last updated:** 2026-06-07
