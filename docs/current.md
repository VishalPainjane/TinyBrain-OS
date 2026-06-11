# Current

**Version:** V0.6 Inference — shipped in docs; **git tag `v0.6` not created yet**

**Sprint:** Product CLI — **complete** (post-v0.6 shell; not part of v0.6 tag scope)

**Task:** (none — all local work pending commit)

**Uncommitted work (three tracks):**

1. v0.6 release closure — `planning/releases/v0.6.md`, retrospective, `progress.md`
2. `stab-003` — `cmd/tinybrain` CLI composition root
3. Repo hygiene — `CHANGELOG.md`, `CONTRIBUTING.md`, `SECURITY.md`, testing policy, boundary/fuzz/golden tests, `RELEASE-CHECKLIST.md`, `repo-health.md`

**009a–009d:** Merged on `main` history; see [planning/releases/v0.6.md](../planning/releases/v0.6.md)

**CUDA:** Matrix **Partial** until [009d manual GPU checklist](../planning/decisions/009d-manual-gpu-checklist.md) signed

**Forbidden:** `internal/agents/`, Kubernetes, web dashboard, API/Router layer — see [current-sprint.md](../planning/execution/current-sprint.md). KV manager (011) and swap manager (012) remain sprint-forbidden until v0.7+ planning unlocks them.

→ Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)  
→ Trust status: [planning/metrics/repo-health.md](../planning/metrics/repo-health.md)  
→ Release gate: [planning/releases/RELEASE-CHECKLIST.md](../planning/releases/RELEASE-CHECKLIST.md)

**Next after commit:** v0.7 MLFQ scheduler ([master-roadmap.md](../planning/roadmap/master-roadmap.md)) or manual GPU sign-off
