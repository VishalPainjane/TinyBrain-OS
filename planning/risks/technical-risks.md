# Technical Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| llama.cpp CGO integration harder than expected | Medium | High | InferenceProvider abstraction first (v0.4); stub until v0.6 |
| MLFQ token-quantum scheduling causes swap thrashing | Medium | High | Swap heuristic (10s idle); migration-path FIFO→MLFQ |
| 4 GB VRAM budget too tight for demo fleet | Medium | Medium | Hardware profiles; dynamic model downgrade |
| Process table design blocks scheduler | Low | High | Contracts locked before task 002; migration-path doc |
| Scope creep into agents before runtime | High | High | Forbidden work in current-sprint; spec gates |
| Documentation drift from code | Medium | Medium | Code-first rule; completed.md + retrospectives |
| Windows CGO toolchain friction for llama.cpp | Medium | Medium | Document build requirements; test early at v0.6 |
| BoltDB vs SQLite undecided delays v0.5 | Low | Low | Mark TBD in v0.5 spec; in-memory registry first |

---
**Layer:** planning
**Last updated:** 2026-06-07
