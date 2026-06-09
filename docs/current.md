# Current

**Version:** V0.6 Inference
**Sprint:** Runtime and Inference
**Task:** (none — 009c complete; GPU offload next)

**009a status:** Merged — llama.cpp CPU load/unload (`b9553` @ `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`)

**009b status:** Complete (`8cdd9b1`, PR #4) — CPU `Generate` on `LlamaProvider`; Linux CGO + real GGUF integration verified

**009c status:** Merged (`ee186c1`, PR #5 @ `f20475f`) — `ModelRuntime` orchestrates resolver → loader → `LlamaProvider`; shared `runtime.ModelResolver`

**Forbidden:** agents, Kubernetes, KV manager (011 — Month 3)

→ Full sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)

**Prior release:** tag `v0.5` — [planning/releases/v0.5.md](../planning/releases/v0.5.md)
