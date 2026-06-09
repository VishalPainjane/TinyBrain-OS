# Task STAB-002 — CI Observability

## Status

Complete — observability shipped (PR #7); artifact-only history cleanup in follow-up PR

## Goal

Implement [ci-observability-architecture-review.md](../planning/decisions/ci-observability-architecture-review.md) — measure and expose CI timing, cache behavior, and failure phases so merge-blocking jobs can be tracked for speed and reliability over time. Workflow + metrics + docs only.

## Context

STAB-001 established four merge-blocking jobs with M1–M2 guards and portable llama cache. GitHub Actions provides per-run job duration in the UI but no structured history. Cold-cache cost (M6) and cache invalidation events (e.g. `portable-v2` after SIGILL) need trend visibility before CUDA work expands CI surface.

Post-merge, in-repo `ci-runs.jsonl` git push failed under branch protection (`GH006`). Cleanup adopts artifact-only history per follow-up review.

## Requirements

- Per-job timing and cache hit/miss in `$GITHUB_STEP_SUMMARY`
- Track llama build, GGUF download, integration test duration separately
- Run-level `ci-run-record-{run_id}` artifact on each `main` push CI run
- Non-blocking collector job; merge-blocking jobs unchanged
- Runtime integration M2 diagnostics parity with llama job
- No Go / runtime / inference code changes
- No new third-party dependencies
- No git writes to `main` from CI

## Files

See architecture review §1–§2 and v2 delta.

## Acceptance Criteria

### Observability emission

- [x] AC-O1: Every merge-blocking job writes timing/cache subsection to `$GITHUB_STEP_SUMMARY`
- [x] AC-O2: Integration jobs record cache + duration fields (llama build, GGUF, integration test)
- [x] AC-O3: `inference-cgo` records `llama_build_duration_s` and `cgo_test_duration_s`
- [x] AC-O4: `test` job records `go_test_duration_s`
- [x] AC-O5: Runtime job M2 diagnostics parity with llama job
- [x] AC-O6: Failure artifacts for both integration jobs
- [x] AC-O7: Branch protection check names unchanged
- [x] AC-O8: No new required third-party dependencies

### History and baselines

- [x] AC-H1′: Each `main` push uploads artifact `ci-run-record-{run_id}` with merged JSON (schema v1)
- [x] AC-H2′: `ci-schema.md` documents artifact record fields
- [x] AC-H3: `ci-baseline.md` defines warm/cold bands and review cadence (artifact-based)
- [x] AC-H4′: Collector non-blocking; no git writes; job reaches green on successful merge

### Documentation

- [x] AC-D1: `testdata/ci/README.md` observability section
- [x] AC-D2: `README.md` links observability docs
- [x] AC-D3: `accepted.md` and `technical-risks.md` updated

### Cleanup (STAB-002 v2)

- [x] AC-C1: No `git commit` / `git push` in any workflow step
- [x] AC-C2: `ci-metrics-collect` has no `contents: write` permission
- [x] AC-C3: `planning/metrics/ci-runs.jsonl` removed from repository
- [ ] AC-C4: Post-cleanup `main` run: collector green; artifact downloadable with 4 jobs in `jobs[]` (verify after merge)
- [x] AC-C5: No live doc references `ci-runs.jsonl` as source of truth
- [x] AC-C6: Repository-wide search for `ci-runs.jsonl`; remaining refs classified

## Implementation Order

1. Metrics scaffold (`ci-schema.md`, `ci-baseline.md`) — done
2. `.github/scripts/ci-emit-metrics.sh` — done
3. Instrument `.github/workflows/ci.yml` — done (PR #7)
4. Runtime M2 diagnostics + artifact upload — done
5. Cleanup: artifact-only collector, delete jsonl, doc sync — this PR
6. Verify AC-C4 on first post-cleanup `main` run
7. Complete task entry in `completed.md`

## Out Of Scope

- CUDA / GPU CI jobs
- Shared bootstrap job (STAB-001 deferred)
- External observability SaaS
- Go code changes
- Changing M2 pass-count thresholds
- Making metrics collection merge-blocking
- Branch protection changes
- `github-actions[bot]` bypass

## Related

- Review: [planning/decisions/ci-observability-architecture-review.md](../planning/decisions/ci-observability-architecture-review.md)
- Prior: [stab-real-inference-ci-gate.md](stab-real-inference-ci-gate.md)
- Baseline: [planning/metrics/ci-baseline.md](../planning/metrics/ci-baseline.md)
- Schema: [planning/metrics/ci-schema.md](../planning/metrics/ci-schema.md)
- Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)

---
**Layer:** task
