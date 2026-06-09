# Task STAB-002 — CI Observability

## Status

In Progress — workflow instrumentation + metrics scaffold; pending first `main` run with collector

## Goal

Implement [ci-observability-architecture-review.md](../planning/decisions/ci-observability-architecture-review.md) — measure and expose CI timing, cache behavior, and failure phases so merge-blocking jobs can be tracked for speed and reliability over time. Workflow + metrics + docs only.

## Context

STAB-001 established four merge-blocking jobs with M1–M2 guards and portable llama cache. GitHub Actions provides per-run job duration in the UI but no in-repo history. Cold-cache cost (M6) and cache invalidation events (e.g. `portable-v2` after SIGILL) need trend visibility before CUDA work expands CI surface.

## Requirements

- Per-job timing and cache hit/miss in `$GITHUB_STEP_SUMMARY`
- Track llama build, GGUF download, integration test duration separately
- Append `main` run metrics to `planning/metrics/ci-runs.jsonl`
- Non-blocking collector job; merge-blocking jobs unchanged
- Runtime integration M2 diagnostics parity with llama job
- No Go / runtime / inference code changes
- No new third-party dependencies

## Files

See architecture review §1–§2.

## Acceptance Criteria

### Observability emission

- [ ] AC-O1: Every merge-blocking job writes timing/cache subsection to `$GITHUB_STEP_SUMMARY`
- [ ] AC-O2: Integration jobs record cache + duration fields (llama build, GGUF, integration test)
- [ ] AC-O3: `inference-cgo` records `llama_build_duration_s` and `cgo_test_duration_s`
- [ ] AC-O4: `test` job records `go_test_duration_s`
- [ ] AC-O5: Runtime job M2 diagnostics parity with llama job
- [ ] AC-O6: Failure artifacts for both integration jobs
- [ ] AC-O7: Branch protection check names unchanged
- [ ] AC-O8: No new required third-party dependencies

### History and baselines

- [ ] AC-H1: `ci-runs.jsonl` receives one row per `main` CI run
- [ ] AC-H2: `ci-schema.md` documents every JSONL field
- [ ] AC-H3: `ci-baseline.md` defines warm/cold bands and review cadence
- [ ] AC-H4: Collector non-blocking

### Documentation

- [ ] AC-D1: `testdata/ci/README.md` observability section
- [ ] AC-D2: `README.md` links observability docs
- [ ] AC-D3: `accepted.md` and `technical-risks.md` updated

## Implementation Order

1. Metrics scaffold (`ci-schema.md`, `ci-baseline.md`, `ci-runs.jsonl`) — **done with review landing**
2. `.github/scripts/ci-emit-metrics.sh`
3. Instrument `.github/workflows/ci.yml` (timers, cache ids, summaries)
4. `ci-metrics-collect` job (non-blocking, `main` only)
5. Runtime M2 diagnostics + artifact upload
6. Doc sync
7. Feature branch: green run + failure validation
8. Merge; 5+ `main` runs vs baseline
9. Complete task in `completed.md`

## Out Of Scope

- CUDA / GPU CI jobs
- Shared bootstrap job (STAB-001 deferred)
- External observability SaaS
- Go code changes
- Changing M2 pass-count thresholds
- Making metrics collection merge-blocking

## Related

- Review: [planning/decisions/ci-observability-architecture-review.md](../planning/decisions/ci-observability-architecture-review.md)
- Prior: [stab-real-inference-ci-gate.md](stab-real-inference-ci-gate.md)
- Baseline: [planning/metrics/ci-baseline.md](../planning/metrics/ci-baseline.md)
- Schema: [planning/metrics/ci-schema.md](../planning/metrics/ci-schema.md)
- Sprint: [planning/execution/current-sprint.md](../planning/execution/current-sprint.md)

---
**Layer:** task
