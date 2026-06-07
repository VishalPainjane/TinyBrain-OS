# MLFQ Notes

**Non-normative — do not treat as architecture law.**

Background research for TinyBrain scheduler design.

## OS MLFQ Background

Multi-Level Feedback Queue is a classic CPU scheduling algorithm. Interactive/short jobs stay in high queues; CPU-bound jobs demote to lower queues over time. Starvation prevented by periodic boost (all jobs moved to top queue).

References: standard OS textbooks; GeeksforGeeks MLFQ overview.

## TinyBrain Adaptation

Traditional OS uses time quanta in milliseconds. TinyBrain uses **token quanta** — scheduling decisions occur per generated token, not per wall-clock slice. This matches LLM inference where work unit is token generation.

| Queue | Token quantum |
|-------|---------------|
| Q0 | 32 |
| Q1 | 64 |
| Q2 | 128 |
| Q3 | 256 |

## Starvation Prevention

Boost all processes to Q0 every 30 seconds or every 500 generated tokens (TBD in v0.7 implementation).

## Migration

See [planning/architecture-evolution/migration-paths.md](../../planning/architecture-evolution/migration-paths.md) Path 1.

---
**Layer:** research
