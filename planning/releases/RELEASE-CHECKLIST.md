# Release Checklist

One-page gate for shipping a version tag. Solo developer use only.

Full planning sync (optional extras): [roadmap/update-checklist.md](../roadmap/update-checklist.md) Tier A + B.

---

## Before tagging

### Verification

- [ ] All merge-blocking CI jobs green on `main`:
  - `test`
  - `inference-cgo`
  - `inference-integration`
  - `inference-integration-runtime`
- [ ] Local: `CGO_ENABLED=0 go test ./...`
- [ ] If inference/runtime changed: local integration tier (see [CONTRIBUTING.md](../../CONTRIBUTING.md))
- [ ] If CUDA changed: [009d manual GPU checklist](../decisions/009d-manual-gpu-checklist.md) signed

### Documentation

- [ ] [planning/releases/vX.Y.md](.) — Status → Shipped; features, lessons, known issues
- [ ] [CHANGELOG.md](../../CHANGELOG.md) — version section with date and user-facing changes
- [ ] [docs/current.md](../../docs/current.md) — version, sprint, forbidden work
- [ ] [README.md](../../README.md) — current version row if user-facing scope changed
- [ ] [planning/metrics/progress.md](../metrics/progress.md) — version at 100%
- [ ] [planning/metrics/repo-health.md](../metrics/repo-health.md) — last verified tag, known gaps

### Git

- [ ] Final task in [planning/execution/completed.md](../execution/completed.md) with commit hash
- [ ] Tag: `vX.Y` (e.g. `v0.7`)
- [ ] Tag points at commit that passed CI above

---

## After tagging

- [ ] [planning/postmortems/](../postmortems/) — file from [template.md](../postmortems/template.md)
- [ ] [planning/execution/current-sprint.md](../execution/current-sprint.md) — next sprint goal
- [ ] Record manual gaps explicitly (e.g. CUDA partial) in release notes and repo-health

---

## What green CI does not prove

- GPU/CUDA runtime on real hardware (manual checklist)
- Metal, ROCm, Vulkan backends
- Multi-agent orchestration or brain-top
- Performance vs [benchmark targets](../benchmarks/target-metrics.md)

---

**Layer:** planning  
**Last updated:** 2026-06-11
