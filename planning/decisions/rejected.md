# Rejected Decisions

Ideas explicitly rejected. Do not re-implement without a new ADR overriding the rejection.

| Decision | Status | Reason |
|----------|--------|--------|
| Fixed Planner/Browser/Coder/Reasoner classes in core | Rejected | Agents are plugins (ADR-002) |
| Postgres/Redis/Kafka in MVP | Rejected | Zero value, huge complexity (detail.md Part 3) |
| Scheduler talks directly to llama.cpp | Rejected | Violates constitution and INV-001/INV-008 |
| Hardcoded model sizes per agent role | Rejected | Hardware-aware runtime (ADR-001) |
| Natural-language IPC between agents | Rejected | Structured JSON only (constitution rule 7) |
| Node.js for orchestration core | Rejected | Go chosen for K8s ecosystem |
| Building agents before runtime/scheduler | Rejected | Runtime is the moat; agents are replaceable |
| Starting with MLFQ before FIFO skeleton | Rejected | Migration path FIFO→MLFQ |

---
**Layer:** planning
**Related:** [accepted.md](accepted.md), [../../docs/adr/](../../docs/adr/)
