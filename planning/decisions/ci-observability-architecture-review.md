# Architecture Review — CI Observability (STAB-002)

**Status:** Approved for implementation  
**Date:** 2026-06-10  
**Scope:** Workflow instrumentation, metrics storage, documentation. **No Go code. No runtime/inference changes.**

**Builds on:** [real-inference-ci-gate-architecture-review.md](real-inference-ci-gate-architecture-review.md) (STAB-001 v3) — four merge-blocking jobs, M1–M2 guards, portable llama cache (`LLAMA_BUILD_CACHE_KEY: llama-cpu-b9553-ubuntu-portable-v2`), branch protection on `main`.

**Goal:** Measure and expose CI behavior so TinyBrain can see whether merge-blocking jobs stay **fast** and **reliable** over time — without weakening gates.

---

## 1. Exact Files to Create

| File | Purpose |
|------|---------|
| `planning/decisions/ci-observability-architecture-review.md` | This document (binding spec) |
| `tasks/stab-002-ci-observability.md` | Task spec, acceptance criteria, implementation order |
| `planning/metrics/ci-baseline.md` | SLO targets and interpretation guide (warm/cold bands from STAB-001 M6) |
| `planning/metrics/ci-runs.jsonl` | Append-only run history (one JSON object per line per `main` CI run) |
| `planning/metrics/ci-schema.md` | JSONL field definitions, enums, example row |
| `.github/scripts/ci-emit-metrics.sh` | Shared shell helpers: timers, cache-hit labels, `$GITHUB_STEP_SUMMARY` rows |

**Not created:** Go packages, external SaaS integrations, Datadog/Prometheus (deferred unless approved later).

---

## 2. Exact Files to Modify

| File | Change |
|------|--------|
| `.github/workflows/ci.yml` | Step timers, cache `id`/`outputs`, `$GITHUB_STEP_SUMMARY` tables, non-blocking `ci-metrics-collect` job (see §4–§9) |
| `testdata/ci/README.md` | CI observability section: metric names, cache semantics, how to read `ci-runs.jsonl`, baseline drift procedure |
| `README.md` | Link to `planning/metrics/ci-baseline.md` and observability summary |
| `planning/risks/technical-risks.md` | CI observability risks; extend STAB-001 mitigation row |
| `planning/metrics/progress.md` | One-line pointer to `ci-runs.jsonl` |
| `planning/decisions/accepted.md` | Row: STAB-002 CI observability adopted |
| `planning/execution/backlog.md` | STAB-002 entry when scheduled |

**Explicitly zero-diff:**

| Path | Reason |
|------|--------|
| `internal/**` | No inference/runtime/scheduler changes |
| `testdata/ci/*.sha256`, `testdata/ci/llama.cpp.pin` | Fixture identity unchanged |
| Branch protection check names | `test`, `inference-cgo`, `inference-integration`, `inference-integration-runtime` |

---

## 3. Metrics to Record

### 3.1 Per workflow run (`main` push; PR optional v2)

| Field | Source |
|-------|--------|
| `run_id`, `run_attempt`, `event`, `ref`, `head_sha`, `conclusion`, `created_at`, `updated_at` | GitHub context / `gh run view` |
| `workflow` | `"CI"` |

### 3.2 Per job

| Field | Jobs | Description |
|-------|------|-------------|
| `job_name`, `job_conclusion`, `job_duration_s` | all four | GitHub job timing |
| `setup_go_cache` | all four | `hit` \| `miss` \| `unknown` |
| `go_test_duration_s` | `test`, `inference-cgo` | Primary test step wall time |

### 3.3 Integration jobs (`inference-integration`, `inference-integration-runtime`)

| Field | Description |
|-------|-------------|
| `llama_build_cache` | `hit` \| `miss` (lib present before cmake) |
| `llama_build_duration_s` | 0 on cache hit; cmake+build wall time on miss |
| `gguf_cache` | `hit` \| `miss` (`actions/cache` output on `gguf-cache` step) |
| `gguf_download_duration_s` | 0 if file existed; curl wall time if downloaded |
| `gguf_verify_duration_s` | sha256sum step |
| `integration_test_duration_s` | M1+M2 `go test -json` block |
| `m2_pass_count`, `m2_run_count`, `m2_skip_count`, `m2_fail_count` | Parsed from existing summary / env |
| `failure_phase` | `none` \| `bootstrap` \| `llama_build` \| `gguf_download` \| `checksum` \| `integration_test` \| `m2_guard` \| `other` |
| `llama_cache_key`, `gguf_cache_key` | Env snapshot at run time (explains hit-rate drops after M3/M4 bumps) |

### 3.4 `inference-cgo` only

| Field | Description |
|-------|-------------|
| `llama_build_duration_s` | Full cmake build (no Actions cache today) |
| `cgo_test_duration_s` | CGO unit test step |

### 3.5 Reliability rollups (computed in collector or manual review)

- Job success rate (7-day / 30-day on `main`)
- Warm vs cold run total wall time (integration jobs)
- Cache hit ratio per cache key family
- Failure count by `failure_phase`
- Regression signal: job duration > baseline p95 for 3 consecutive runs

---

## 4. Where to Store CI Timing / History

### Primary store (repo-local, planning layer)

```
planning/metrics/ci-runs.jsonl
```

- One JSON object per line, one row per **`main` push** CI run.
- Append-only; never rewrite history.
- Schema: [ci-schema.md](../metrics/ci-schema.md).

### Secondary (ephemeral, GitHub-native)

- `$GITHUB_STEP_SUMMARY` markdown table per job (every run).
- Failure artifacts: `go-test-llama-json-{run_id}` (exists); add `go-test-runtime-json-{run_id}` and runtime M2 diagnostics parity.

### Collector job (non-merge-blocking)

**Option A (preferred):** Final job in same workflow:

```yaml
ci-metrics-collect:
  needs: [test, inference-cgo, inference-integration, inference-integration-runtime]
  if: always() && github.ref == 'refs/heads/main' && github.event_name == 'push'
  continue-on-error: true
```

**Option B:** Separate workflow `on: workflow_run` after `CI` completes on `main`.

Collector uses `gh run view` + step env files / logs to build JSONL row. Commits via bot PR to `planning/metrics/ci-runs.jsonl` (avoids bypassing branch protection with direct push).

### Human-readable snapshot

- [ci-baseline.md](../metrics/ci-baseline.md) — SLO bands from M6; updated manually when p95 shifts (not every run).

### Not in v1

External time-series DB, third-party dependencies.

---

## 5. How to Track Cache Hit / Miss Behavior

| Cache | Current state | Implementation change | Signal |
|-------|---------------|----------------------|--------|
| llama.cpp build | `Restore llama.cpp build cache` has **no** `id` | Add `id: llama-build-cache` | `steps.llama-build-cache.outputs.cache-hit` |
| llama.cpp build (confirm) | Build step logs `llama.cpp build cache hit` | Emit `TB_LLAMA_BUILD_CACHE=hit\|miss` to `$GITHUB_ENV` | File probe + skip cmake |
| GGUF model | `id: gguf-cache` exists | Read output in metrics step | `steps.gguf-cache.outputs.cache-hit` |
| setup-go | No explicit id | Add `id: setup-go` or parse cache log | `hit` \| `miss` \| `unknown` |

**Step summary example:**

```markdown
| Cache       | Result |
|-------------|--------|
| llama build | hit    |
| GGUF model  | miss   |
| setup-go    | hit    |
```

**JSONL:** `llama_build_cache`, `gguf_cache`, `setup_go_cache` as `"hit"|"miss"`.

**STAB-001 lesson:** Record `llama_cache_key` and `gguf_cache_key` in each row so hit-rate drops after M3/M4 bumps are explainable (e.g. `portable-v2` invalidation after SIGILL).

---

## 6. How to Track Model Download Time

In step `Download GGUF model if missing`:

1. If `[ -f "$MODEL_PATH" ]` → `GGUF_DOWNLOAD_S=0`, `GGUF_SOURCE=cache`.
2. Else: `START=$(date +%s)`, `curl -fsSL`, `END=$(date +%s)`, `GGUF_DOWNLOAD_S=$((END-START))`, `GGUF_SOURCE=download`.
3. Write to `$GITHUB_ENV`; append to step summary.

**Do not** include sha256 verify in download time — separate field `gguf_verify_duration_s`.

Fork PRs: optional `is_fork: true` on PR runs so baselines are not polluted (v1 records `main` only).

---

## 7. How to Track llama Build Time

### Integration jobs (cached)

1. After cache restore, read `steps.llama-build-cache.outputs.cache-hit`.
2. In `Build llama.cpp (CPU, libraries only)`:
   - Cache hit (lib exists): `LLAMA_BUILD_S=0`, log `llama build: cache hit`.
   - Miss: time cmake configure + `cmake --build` → `LLAMA_BUILD_S`.

### `inference-cgo` (uncached)

Always time cmake+build; record `llama_build_duration_s` (establishes upper bound ~3–8 min per M6).

**Out of scope:** Sharing llama cache with `inference-cgo` (bootstrap redesign).

---

## 8. How to Track Integration Test Duration

Wall clock around existing M1+M2 block (guard logic unchanged):

```bash
TB_TEST_START=$(date +%s)
# existing go test -json ... + python M2
TB_TEST_END=$(date +%s)
INTEGRATION_TEST_S=$((TB_TEST_END - TB_TEST_START))
```

Extend stdout line:

```text
integration summary: pass=13 skip=0 run=13 duration_s=12
```

**Runtime job gap (STAB-002):** Add M2 diagnostics parity with llama job (pass/run/skip/fail counts, first 50 fail/output JSON lines, `go test exit code`, failure artifact upload). Thresholds unchanged.

**Optional diagnostic:** Sum `"Elapsed"` from `go test -json` pass events vs wall time.

---

## 9. How to Report Failures Clearly

| Layer | When | Content |
|-------|------|---------|
| Step log | Always | M2 diagnostics (llama; extend to runtime) |
| Step summary | `failure()` | Phase, counts, cache states, durations, `go test exit code` |
| Artifact | Integration failure | `go-test-llama-json-{run_id}`; add runtime equivalent |
| JSONL row | Collector on `main` | `conclusion=failure`, `failure_phase`, error class (`SIGILL`, `pass_count_low`, `sha256_mismatch`) — no secrets |
| `ci-baseline.md` | Manual triage | “If you see X, check Y” (portable build, M4, M3 bump) |

### Failure phase mapping

| Phase | Failing step (current workflow) |
|-------|--------------------------------|
| `bootstrap` | Verify llama.cpp submodule SHA / upstream tag |
| `llama_build` | Build llama.cpp (CPU, libraries only) |
| `gguf_download` | Download GGUF model if missing |
| `checksum` | Verify GGUF SHA256 / Assert GGUF model size |
| `integration_test` | `go test` exit ≠ 0 before M2 |
| `m2_guard` | Python M2 exits 1 |
| `other` | checkout, apt, setup-go |

Metrics collection must **not** be merge-blocking (`continue-on-error: true` on emit steps and collector).

---

## 10. Acceptance Criteria

### Observability emission (AC-O1–O8)

- [ ] **AC-O1:** Every merge-blocking job writes timing/cache subsection to `$GITHUB_STEP_SUMMARY`.
- [ ] **AC-O2:** Integration jobs record `llama_build_cache`, `gguf_cache`, `llama_build_duration_s`, `gguf_download_duration_s`, `integration_test_duration_s`.
- [ ] **AC-O3:** `inference-cgo` records `llama_build_duration_s` and `cgo_test_duration_s`.
- [ ] **AC-O4:** `test` job records `go_test_duration_s`.
- [ ] **AC-O5:** Runtime integration job has M2 failure diagnostics parity with llama job.
- [ ] **AC-O6:** Failure artifacts uploaded for both integration jobs on test-step failure.
- [ ] **AC-O7:** All four jobs remain merge-blocking; branch protection check names unchanged.
- [ ] **AC-O8:** No new required third-party dependencies.

### History and baselines (AC-H1–H4)

- [ ] **AC-H1:** `planning/metrics/ci-runs.jsonl` receives one row per `main` CI run after merge.
- [ ] **AC-H2:** `planning/metrics/ci-schema.md` documents every JSONL field.
- [ ] **AC-H3:** `planning/metrics/ci-baseline.md` defines warm/cold duration bands and review cadence.
- [ ] **AC-H4:** Collector job is non-blocking; failed collector does not fail CI.

### Documentation (AC-D1–D3)

- [ ] **AC-D1:** `testdata/ci/README.md` observability section cross-links schema + baseline.
- [ ] **AC-D2:** `README.md` links CI observability docs.
- [ ] **AC-D3:** `planning/decisions/accepted.md` and `technical-risks.md` updated on task completion.

---

## 11. Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| Metrics collector commits bypass or fight branch protection | Medium | Medium | Bot PR; collector non-blocking |
| JSONL repo bloat | Low | Low | ~1 KB/run; archive quarterly to `planning/metrics/archive/ci-runs-YYYY.jsonl` |
| Timer noise from runner variance | High | Low | Track hit/miss separately; p95 baselines |
| Duplicate timing logic drifts across jobs | Medium | Medium | Single `.github/scripts/ci-emit-metrics.sh` |
| Leaking paths/URLs in metrics | Low | Medium | Schema allowlist; no `CI_GGUF_URL` in JSONL |
| Scope creep into bootstrap refactor | Medium | High | STAB-002 lock: no shared bootstrap job |
| PR runs pollute baseline | Medium | Low | v1: `main` only |
| False confidence from UI summaries alone | Medium | Medium | JSONL on `main` is trend source of truth |
| Metrics step breaks merge jobs | Low | High | `continue-on-error: true`; sidecar collector |

---

## Implementation Order

1. Land review + task + metrics scaffold (`ci-schema.md`, `ci-baseline.md`, empty `ci-runs.jsonl`).
2. Add `.github/scripts/ci-emit-metrics.sh` + workflow instrumentation.
3. Add non-blocking `ci-metrics-collect` job for `main`.
4. Extend runtime M2 diagnostics + artifact (reporting only).
5. Doc sync (`testdata/ci/README.md`, `README.md`, risks, accepted).
6. Verify on feature branch: one green run + one forced failure; validate summary + JSONL shape.
7. Merge; observe 5+ `main` runs against `ci-baseline.md`.

---

## Rollback Plan

1. Revert STAB-002 workflow commits; four-job topology unchanged.
2. Stop collector; retain `ci-runs.jsonl` as historical record.
3. Step summaries disappear; STAB-001 M2 guards and llama diagnostics remain.

---

## v2 Delta — Artifact-only history (STAB-002 cleanup)

**Date:** 2026-06-10  
**Reason:** Post-merge run [27204998799](https://github.com/VishalPainjane/TinyBrain-OS/actions/runs/27204998799) — `ci-metrics-collect` git push rejected (`GH006`) under protected `main`. Branch protection and bot bypass explicitly forbidden.

| v1 (superseded) | v2 (binding) |
|-----------------|--------------|
| Append to `planning/metrics/ci-runs.jsonl` on `main` | **Removed** — file deleted from repo |
| `git commit` + `git push` in collector | **Removed** |
| `permissions: contents: write` on collector | **Removed** |
| Backup artifact `ci-run-record-{run_id}` | **Primary** deliverable — single merged JSON per run |
| JSONL schema doc | **Artifact schema** doc (`ci-schema.md`) |

**Unchanged:** Step Summary metrics, per-job `ci-metrics-{job}` artifacts, merge-blocking four jobs, M2 guards, `ci-baseline.md`.

**AC-H1 superseded by AC-H1′:** each `main` push produces artifact `ci-run-record-{run_id}`.

---
**Layer:** planning  
**Related:** [real-inference-ci-gate-architecture-review.md](real-inference-ci-gate-architecture-review.md), [../metrics/ci-schema.md](../metrics/ci-schema.md), [../../tasks/stab-002-ci-observability.md](../../tasks/stab-002-ci-observability.md)
