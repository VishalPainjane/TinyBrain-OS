# Architecture Review — Real Inference CI Gate (v3)

**Status:** Approved for implementation  
**Date:** 2026-06-09 (v3 — operational readiness pass incorporated)  
**Scope:** Workflow, CI assets (`testdata/ci/`), documentation only. **No Go code changes.**

**Supersedes:** v2 optional JSON skip guard → v3 **mandatory** pipefail-safe execution + minimum pass-count enforcement (M1–M2). v2 split jobs, tag + SHA verification unchanged.

**Goal:** Green CI on `main` proves real GGUF inference on Linux CPU — adapter and runtime E2E — with **no silent test skips** and **no masked `go test` failures**.

---

## v3 Delta (operational readiness — M1–M6)

| ID | Change from v2 | Binding for implementation |
|----|----------------|----------------------------|
| **M1** | JSON skip guard was "recommended" | **Mandatory:** all integration test steps use pipefail-safe execution; `go test` exit code must not be masked by `\| tee` or JSON piping |
| **M2** | Skip grep only | **Mandatory:** fail if pass count below job threshold, skip count > 0, or test run count = 0 |
| **M3** | Cache key bump implied | **Documented procedure:** cmake flag changes require `LLAMA_BUILD_CACHE_KEY` version suffix bump in `testdata/ci/README.md` |
| **M4** | Checksum invalidates cache | **Documented procedure:** HuggingFace URL and SHA256 updated together; checksum is source of truth |
| **M5** | "After first green run" one line | **Explicit rollout sequence** for branch protection (§8.2) |
| **M6** | "Acceptable for v1" brief note | **Documented cold-cache expectations:** duplicate bootstrap, duration estimates, rationale (§4.5, `testdata/ci/README.md`) |

---

## 1. Exact Files to Create

| File | Purpose |
|------|---------|
| `testdata/ci/README.md` | Model identity, source URL, license, size limit; **M3** llama build cache bump procedure; **M4** GGUF URL+SHA256 coupled refresh; **M6** cold-cache expectations; what green CI proves / does not prove |
| `testdata/ci/smollm2-135m-instruct-q4_k_m.gguf.sha256` | Single-line SHA256 hex (source of truth read by workflow) |
| `testdata/ci/llama.cpp.pin` | Submodule pin record: upstream tag `b9553`, commit SHA `9e3b928…` (human + CI reference) |
| `tasks/stab-real-inference-ci-gate.md` | Task spec with acceptance criteria |
| `planning/decisions/real-inference-ci-gate-ci-proves.md` | *(optional excerpt in README instead)* — can live as §7 of this doc only |

**Not created:** Go files, CUDA assets, GGUF binary in repo (`.gitignore` keeps `*.gguf` out).

---

## 2. Exact Files to Modify

| File | Change |
|------|--------|
| `.github/workflows/ci.yml` | Add `inference-integration` + `inference-integration-runtime` jobs; remove runtime integration step from `inference-cgo`; add workflow `env`; shared bootstrap steps |
| `README.md` | CI jobs table; local integration mirrors `testdata/ci/`; link to what green CI proves |
| `docs/architecture/inference-backend-matrix.md` | CPU CI coverage: unit **Yes**, integration **Yes** (merge-blocking); document two integration jobs |
| `planning/decisions/accepted.md` | Row: Real Inference CI Gate — deterministic GGUF + dual integration jobs |
| `planning/decisions/009a-llama-cpp-dependency.md` | CI table: replace optional `inference-integration` with required jobs; remove `vars.TB_TEST_GGUF_PATH` merge dependency |
| `planning/decisions/real-inference-ci-gate-architecture-review.md` | This document (v3) |
| `planning/risks/technical-risks.md` | Close / update “E2E not enforced on merge” risk |

**Explicitly zero-diff:**

| Path | Reason |
|------|--------|
| `internal/inference/**` | No inference implementation changes |
| `internal/runtime/**` | No runtime changes |
| `internal/scheduler/**`, `internal/loader/**`, etc. | Out of scope |
| CUDA / GPU workflow steps | Forbidden |

---

## 3. Workflow Design

### 3.1 Job topology

```text
push / pull_request → main
│
├─ test                          [merge blocker]  CGO_ENABLED=0  go test ./...
│
├─ inference-cgo                 [merge blocker]  CGO unit only (no -tags integration)
│     └─ remove: runtime integration step
│
├─ inference-integration         [merge blocker]  provider / llama integration (5 tests)
│     └─ parallel with runtime job (shared caches)
│
└─ inference-integration-runtime [merge blocker]  runtime E2E (1 test)
      └─ parallel with llama job (shared caches)
```

| Job | Package under test | `-tags` | Tests |
|-----|-------------------|---------|-------|
| `inference-integration` | `./internal/inference/llama/...` | `integration` | 5 (load + generate suite) |
| `inference-integration-runtime` | `./internal/runtime/...` | `integration` | 1 (`TestIntegratedRuntime_Llama_LoadGenerateUnload`) |

**Why two jobs (practical):**

- Failures isolate **adapter** vs **orchestration** in GitHub UI.
- Jobs run **in parallel** — total wall time ≈ one bootstrap + max(llama tests, runtime test).
- Shared Actions caches avoid duplicating download/build cost.

**Why not a setup job + artifacts:** Artifact handoff adds complexity; cache keys on submodule SHA + GGUF SHA are sufficient.

### 3.2 Workflow-level `env`

```yaml
env:
  GO_VERSION: "1.22"
  LLAMA_SUBMODULE_SHA: "9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66"
  LLAMA_SUBMODULE_TAG: "b9553"
  CI_GGUF_FILENAME: "smollm2-135m-instruct-q4_k_m.gguf"
  CI_GGUF_URL: "https://huggingface.co/HuggingFaceTB/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf"
  CI_GGUF_CACHE_KEY: "smollm2-135m-instruct-q4_k_m-v1"
  LLAMA_BUILD_CACHE_KEY: "llama-cpu-b9553-ubuntu-v1"
  CI_MIN_INTEGRATION_PASSES_LLAMA: "5"
  CI_MIN_INTEGRATION_PASSES_RUNTIME: "1"
```

`CI_GGUF_SHA256` is **read from** `testdata/ci/smollm2-135m-instruct-q4_k_m.gguf.sha256` in a bootstrap step (single source of truth in repo, not duplicated in YAML).

Per-job minimum pass thresholds (M2):

| Job | Env var | Minimum `"Action":"pass"` count |
|-----|---------|----------------------------------|
| `inference-integration` | `CI_MIN_INTEGRATION_PASSES_LLAMA` | **5** |
| `inference-integration-runtime` | `CI_MIN_INTEGRATION_PASSES_RUNTIME` | **1** |

### 3.3 Bootstrap sequence (identical in both integration jobs)

Each integration job runs these steps **in order**:

| Step | Name | On failure |
|------|------|------------|
| 1 | `actions/checkout@v4` with `submodules: recursive` | clone failure |
| 2 | Verify submodule commit SHA | exit 1 — pin drift |
| 3 | Verify upstream tag at HEAD | exit 1 — tag drift |
| 4 | Read GGUF SHA256 from `testdata/ci/*.sha256` | exit 1 — missing file |
| 5 | `actions/setup-go@v5` | — |
| 6 | Install `build-essential`, `cmake` | — |
| 7 | Cache restore: llama.cpp `build/` dir | miss → rebuild |
| 8 | Build llama.cpp CPU (cmake flags unchanged from today) | exit 1 |
| 9 | Assert `libllama.so` exists | exit 1 |
| 10 | Cache restore: GGUF model dir | miss → download |
| 11 | Download GGUF if missing (`curl -fsSL`) | exit 1 |
| 12 | `sha256sum -c` against repo checksum file | exit 1 — **checksum mismatch** |
| 13 | Assert model file exists + size ≤ 250 MB | exit 1 — **missing / oversized** |
| 14 | Export `MODEL_PATH`, `TB_TEST_GGUF_PATH`, `LD_LIBRARY_PATH` | — |
| 15 | Package-specific integration test execution (M1 + M2 — see §3.6) | exit 1 — test failure, masked exit, skip, or pass-count violation |

### 3.4 Submodule SHA verification (step 2)

```bash
cd third_party/llama.cpp
actual=$(git rev-parse HEAD)
if [ "$actual" != "$LLAMA_SUBMODULE_SHA" ]; then
  echo "Submodule SHA mismatch: got $actual want $LLAMA_SUBMODULE_SHA"
  exit 1
fi
```

### 3.5 Upstream tag verification (step 3)

```bash
cd third_party/llama.cpp
git fetch --tags --force origin 2>/dev/null || true
tag=$(git describe --tags --exact-match HEAD 2>/dev/null || true)
if [ "$tag" != "$LLAMA_SUBMODULE_TAG" ]; then
  echo "Submodule tag mismatch: got '${tag:-none}' want '$LLAMA_SUBMODULE_TAG'"
  exit 1
fi
```

**Note:** If shallow submodule clone lacks tags, bootstrap must `git fetch --tags` in submodule (document in `testdata/ci/README.md` troubleshooting). Tag `b9553` must point at SHA `9e3b928…`.

### 3.6 Integration test execution — anti-skip and pipefail (M1 + M2, mandatory)

Integration tests still contain `t.Skip` when env unset — **workflow must guarantee env is always set** before `go test`:

1. `MODEL_PATH` set in `$GITHUB_ENV` from deterministic cache path.
2. Step 13: `test -f "$MODEL_PATH"` — job fails before Go runs if missing.
3. Step 12: checksum fails before Go runs if tampered/wrong file.
4. Step 15: **mandatory** execution contract below.

#### M1 — Pipefail-safe execution (mandatory)

All integration test steps **must**:

- Enable `set -o pipefail` in the shell step (or equivalent) before invoking `go test`.
- **Never** treat JSON parsing as the sole failure signal when piping `go test` output.
- Capture and honor `go test` exit code explicitly when writing JSON to a file (write JSON to file first without masking exit code, or use `PIPESTATUS` after a pipe).

**Forbidden pattern (masks failures):**

```text
go test ... -json ./... | tee /tmp/out.json
# grep-only validation without checking go test exit code
```

**Required behavior:** If `go test` exits non-zero, step 15 fails **even if** JSON parsing would pass.

#### M2 — Minimum pass-count enforcement (mandatory)

After `go test -json` completes, step 15 **must fail** if **any** of:

| Condition | Threshold |
|-----------|-----------|
| Test pass count (`"Action":"pass"` for test events) | Below job minimum (llama: **5**, runtime: **1**) |
| Skip count (`"Action":"skip"`) | **> 0** |
| Test run count (`"Action":"run"` for tests) | **= 0** |

Implementation may parse JSON with `jq`, Python, or shell — workflow-only; no Go changes.

**Future tests:** Adding integration tests increases pass count; update `CI_MIN_INTEGRATION_PASSES_*` and this document when new merge-blocking tests are added.

#### Layered enforcement summary

| Layer | Catches |
|-------|---------|
| Preflight (steps 12–13) | Missing model, bad checksum, unset path |
| M1 pipefail | Compile failures, assertion failures, panics |
| M2 pass/skip/run counts | `t.Skip`, zero tests compiled/run, partial suite |

### 3.7 `inference-cgo` changes

**Keep:**

- checkout + submodules
- cmake build
- `go test ./internal/inference/llama/...` (no integration tag)

**Remove:**

```yaml
- name: Test runtime integration (CGO CPU)   # DELETE entire step
```

**Do not add** submodule tag verification to `inference-cgo` (pin verification is integration-job responsibility). Optional: add SHA verify to `inference-cgo` for early fail — **defer** to avoid scope creep.

### 3.8 Canonical CI model

| Field | Value |
|-------|-------|
| Model | SmolLM2-135M-Instruct Q4_K_M |
| Quantization | Q4_K_M |
| Approx size | ~80–100 MB (< 250 MB limit) |
| License | Apache-2.0 (document in `testdata/ci/README.md`) |
| Rationale | Same model used in v0.6 manual verification (spec + release notes) |
| Storage | Actions cache + runner temp only |

---

## 4. Cache Strategy

### 4.1 GGUF model cache

| Key | Value |
|-----|-------|
| Path | `${{ runner.temp }}/tb-ci-models/` |
| Cache key | `${{ CI_GGUF_CACHE_KEY }}-${{ hashFiles('testdata/ci/smollm2-135m-instruct-q4_k_m.gguf.sha256') }}` |
| Invalidation | Any change to `.sha256` file → new key → re-download + verify |
| Restore-fallback | Not used — miss triggers download |

### 4.2 llama.cpp build cache

| Key | Value |
|-----|-------|
| Path | `third_party/llama.cpp/build/` |
| Cache key | `${{ LLAMA_BUILD_CACHE_KEY }}-${{ env.LLAMA_SUBMODULE_SHA }}` |
| Invalidation | Submodule SHA change **or** `LLAMA_BUILD_CACHE_KEY` version suffix bump |
| Hit behavior | Skip cmake if `build/bin/libllama.so` exists post-restore |

#### M3 — llama build cache procedure (mandatory documentation)

**Problem:** Cache key includes submodule SHA but **not** cmake flag fingerprint. A cmake flag change with unchanged SHA could restore a stale incompatible `build/` tree.

**Procedure** (must appear in `testdata/ci/README.md`):

1. When **any** cmake argument changes in `.github/workflows/ci.yml` (GPU flags, build targets, Release/Debug, etc.), bump `LLAMA_BUILD_CACHE_KEY` version suffix (e.g. `llama-cpu-b9553-ubuntu-v1` → `v2`).
2. Record the bump in the PR description and task completion note.
3. Submodule SHA bump alone does **not** require manual key bump — SHA in cache key auto-invalidates.
4. Before adding Windows/macOS jobs (v0.7+), include `${{ runner.os }}` in key prefix — do not share Ubuntu cache across OSes.

### 4.3 GGUF cache and M4 refresh procedure

#### M4 — GGUF URL + SHA256 coupled refresh (mandatory documentation)

**Source of truth:** `testdata/ci/smollm2-135m-instruct-q4_k_m.gguf.sha256` in the repository.

**Procedure** (must appear in `testdata/ci/README.md`):

1. Download candidate GGUF from HuggingFace (or mirror).
2. Compute SHA256 locally; write hex to `.sha256` file.
3. Update `CI_GGUF_URL` in workflow env **only if** URL changed.
4. Bump `CI_GGUF_CACHE_KEY` version suffix if filename or checksum changes.
5. Commit `.sha256` (+ URL/env change if applicable) in the same PR.
6. **Never** update URL alone without verifying and updating checksum — `/resolve/main/` can drift.

Cache key includes `hashFiles('.sha256')` — checksum change auto-invalidates GGUF cache. Step 12 `sha256sum -c` runs on every job regardless of cache hit.

### 4.4 Cache miss behavior

| Cache | Miss cost | Acceptable? |
|-------|-----------|-------------|
| GGUF | ~60–120s download + verify | Yes — first run / fork |
| llama build | ~3–8 min compile | Yes — amortized across subsequent runs |

Both integration jobs share cache keys → warm runs avoid duplicate work; see §4.5 for cold-run behavior.

### 4.5 M6 — Cold-cache expectations (v1 accepted)

**Behavior:** `inference-integration` and `inference-integration-runtime` run **in parallel** with identical bootstrap. On **cold cache** (first workflow run, new checksum, new submodule SHA, fork PR without cache write):

| Resource | Cold-run behavior | v1 decision |
|----------|-------------------|-------------|
| GGUF download | **Two** ~80–100 MB downloads (one per runner) | **Accept** |
| llama.cpp cmake build | **Two** full compiles (one per runner) | **Accept** |
| Cache upload | Both jobs may save same key; last write wins | **Accept** — content equivalent |
| Race conditions | No shared filesystem between jobs | **Safe** |

**Expected cold-start duration (order of magnitude):**

| Phase | Per job (cold) | Wall clock (parallel jobs) |
|-------|----------------|----------------------------|
| Checkout + pin verify | ~30–60s | ~30–60s |
| llama.cpp compile | ~3–8 min | ~3–8 min (overlap) |
| GGUF download + SHA256 | ~60–120s | ~60–120s (overlap) |
| Integration tests | ~30–90s | ~max(llama, runtime) |
| **Total (cold)** | ~5–12 min | ~5–12 min |

**Warm cache:** Subsequent runs on same checksum/SHA typically **~2–4 min** per integration job.

**Rationale for accepting duplicate bootstrap in v1:** Avoids bootstrap-job/artifact complexity; duplicate cost hits cold runs and fork PRs only; Actions cache amortizes on `main`. Revisit bootstrap job if cold-run time blocks team velocity (post-merge metric).

**Document in `testdata/ci/README.md`:** cold vs warm expectations, fork PR read-only cache note.

### 4.6 What is not cached

- Go module cache (handled by `setup-go` default)
- Submodule source (checkout each job)

---

## 5. Failure Modes

| # | Condition | Expected behavior | Silent skip? |
|---|-----------|-------------------|--------------|
| F-1 | `TB_TEST_GGUF_PATH` unset | N/A — workflow sets it | No — step 13 fails if path empty |
| F-2 | Model file missing after download | Step 11/13 exit 1 | No |
| F-3 | Wrong GGUF SHA256 | Step 12 `sha256sum -c` exit 1 | No |
| F-4 | Intentionally broken SHA in repo | Step 12 exit 1 on every integration job | No — **AC-8** |
| F-5 | Submodule SHA drift | Step 2 exit 1 | No |
| F-6 | Submodule tag ≠ `b9553` | Step 3 exit 1 | No |
| F-7 | cmake / link failure | Step 8/9 exit 1 | No |
| F-8 | HuggingFace download failure | Step 11 curl exit 1 | No |
| F-9 | Model > 250 MB | Step 13 size check exit 1 | No |
| F-10 | Generate returns empty / load fails | Step 15 go test exit 1 | No |
| F-11 | All tests skipped (env bug) | M2: skip count > 0 or run count = 0 | **Prevented** |
| F-12 | `inference-cgo` broken | Unit tests fail independently | N/A |
| F-13 | Cache serves corrupted file | SHA256 step catches before tests | No |
| F-14 | `go test` fails but pipe masks exit code | M1 pipefail / explicit exit check | **Prevented** |
| F-15 | Partial suite runs (3/5 pass) | M2 pass count below threshold | **Prevented** |
| F-16 | cmake flags change, SHA unchanged | M3 manual `LLAMA_BUILD_CACHE_KEY` bump | Ops procedure |
| F-17 | HF URL content drift, checksum stale | M4 coupled refresh | Ops procedure |

---

## 6. Acceptance Criteria Mapping

| ID | Requirement | Verification method | Job |
|----|-------------|---------------------|-----|
| AC-1 | Dedicated merge-blocking `inference-integration` job | Branch protection + workflow YAML | `inference-integration` |
| AC-2 | Runtime integration separated when practical | Second job `inference-integration-runtime` | both integration jobs |
| AC-3 | Keep `test` + `inference-cgo` unchanged in purpose | Diff review | `test`, `inference-cgo` |
| AC-4 | Remove runtime integration from `inference-cgo` | Step deleted in diff | `inference-cgo` |
| AC-5 | Verify submodule SHA | Step 2 on integration jobs | integration jobs |
| AC-6 | Verify upstream tag `b9553` | Step 3 on integration jobs | integration jobs |
| AC-7 | Build CPU llama.cpp libraries | Steps 7–9 | integration jobs |
| AC-8 | **Intentionally broken SHA256 fails CI** | PR branch with wrong `.sha256` → red | integration jobs step 12 |
| AC-9 | Deterministic model provision | No `vars.TB_TEST_GGUF_PATH` | integration jobs |
| AC-10 | Missing model fails job | Delete file after download test / skip cache | step 13 |
| AC-11 | All 5 llama integration tests run | `-v` log: 5× PASS | `inference-integration` |
| AC-12 | Runtime E2E test runs | `-v` log: PASS | `inference-integration-runtime` |
| AC-13 | No silent skip | M2: skip count = 0, run count > 0 | integration jobs |
| AC-14 | Model < 250 MB | Step 13 size assert | integration jobs |
| AC-15 | No Go feature changes | `git diff` scope | — |
| AC-16 | Docs updated | README, matrix, accepted, 009a, testdata/ci/README | doc review |
| AC-17 | `CGO_ENABLED=0 go test ./...` still green | `test` job | `test` |
| AC-18 | CGO unit path still green | `inference-cgo` job | `inference-cgo` |
| AC-19 | **M1** Pipefail-safe integration execution | Deliberately failing test exits non-zero despite JSON tee | integration jobs |
| AC-20 | **M2** Llama pass count ≥ 5 | JSON parse or log inspection | `inference-integration` |
| AC-21 | **M2** Runtime pass count ≥ 1 | JSON parse or log inspection | `inference-integration-runtime` |
| AC-22 | **M3** Cache bump procedure documented | `testdata/ci/README.md` review | — |
| AC-23 | **M4** GGUF refresh procedure documented | `testdata/ci/README.md` review | — |
| AC-24 | **M5** Branch protection configured post-first-green | Settings audit after rollout §8.2 | — |
| AC-25 | **M6** Cold-cache expectations documented | `testdata/ci/README.md` review | — |

### AC-8 procedure (mandatory before merge)

1. On feature branch, temporarily set wrong hash in `testdata/ci/*.sha256`.
2. Push → both integration jobs must fail at checksum step.
3. Revert hash → jobs green.
4. Document result in PR / task completion note.

---

## 7. What Green CI Proves and Does Not Prove

### 7.1 Proves (when all four jobs green)

| Claim | Evidence |
|-------|----------|
| Pure Go kernel packages compile and pass tests with stubs | `test` job |
| llama.cpp CPU adapter compiles and passes CGO unit tests | `inference-cgo` |
| Submodule pinned at tag **`b9553`** commit **`9e3b928…`** | integration bootstrap steps 2–3 |
| CPU llama.cpp libraries build on Ubuntu with documented cmake flags | integration bootstrap steps 7–9 |
| Checksum-verified GGUF loads via mmap on Linux | integration tests |
| `LlamaProvider` produces non-empty tokens from real GGUF | 5 llama integration tests |
| `ModelRuntime` integrated path: resolver → loader → provider → events | runtime integration test |
| Scheduler still has zero inference imports | unchanged `test` job boundary tests |

### 7.2 Does NOT prove

| Claim | Why |
|-------|-----|
| Windows / macOS CGO or mmap | Linux-only runners |
| CUDA / GPU offload | CPU cmake; no `-tags cuda` |
| Large-model memory safety | SmolLM2-135M only; < 250 MB |
| TPS / TTFT SLA | Tests log timing; no threshold enforced |
| HuggingFace availability without cache | Cache mitigates; first-run needs network |
| Production model registry / persistence paths | Tests use temp registry + env path |
| Loader capacity / eviction correctness | Default capacity 0 in integration |
| SaveContext / RestoreContext | Stubs unchanged |
| Submodule tag fetch on all clone configurations | Document fetch workaround in README |

**Merge gate statement (for README + branch protection):**

> Green CI on `main` requires `test`, `inference-cgo`, `inference-integration`, and `inference-integration-runtime`. The two integration jobs prove real CPU GGUF inference at pin b9553 with checksum-verified SmolLM2-135M-Instruct Q4_K_M — adapter and runtime E2E — on Ubuntu. It does not prove cross-platform builds, GPU paths, or latency SLAs.

---

## 8. Merge Gate (Exact)

### Required status checks on `main`

| Check | Job ID |
|-------|--------|
| Stub / unit path | `test` |
| CGO compile + unit | `inference-cgo` |
| Provider integration | `inference-integration` |
| Runtime integration | `inference-integration-runtime` |

### 8.2 M5 — Branch protection rollout sequence (mandatory)

**Do not** configure required checks until checks exist on `main`. Follow this order:

| Step | Action | Owner |
|------|--------|-------|
| 1 | Merge workflow + `testdata/ci/` PR to `main` | Implementer |
| 2 | Wait for **first successful** workflow run on `main` (all four jobs green) | CI |
| 3 | Open **Settings → Branches → Branch protection → `main`** | Maintainer |
| 4 | Enable **Require status checks to pass** | Maintainer |
| 5 | Search and select **exact check names** from step 2 run (see §8.3) | Maintainer |
| 6 | Enable **Require branches to be up to date** (if team policy allows) | Maintainer |
| 7 | Document configured check names in `README.md` CI section | Implementer |

**If check names are wrong:** Merges block until protection list matches actual job IDs. Re-run workflow on `main` and update settings — do not rename jobs after protection is enabled without updating settings.

### 8.3 Check naming edge cases

| Rule | Detail |
|------|--------|
| Job ID = check name | Do **not** set job-level `name:` override unless branch protection uses that exact string |
| Workflow `name: CI` | UI may display `CI / test` — use names **exactly as shown** in the PR checks dropdown from the first green `main` run |
| Expected IDs | `test`, `inference-cgo`, `inference-integration`, `inference-integration-runtime` |
| Fork PRs | All four jobs still run; cache read-only — not a branch protection issue |

---

## 9. Implementation Order

1. Create `testdata/ci/` — checksum file, `llama.cpp.pin`, **README with M3, M4, M6 procedures**; compute SHA256 locally from canonical GGUF.
2. Update `.github/workflows/ci.yml` — two integration jobs with **M1 pipefail + M2 pass-count enforcement**; trim runtime step from `inference-cgo`; set `LLAMA_BUILD_CACHE_KEY` with `v1` suffix.
3. Verify green on feature branch (all four jobs).
4. Verify **AC-8** (broken SHA256 → red at step 12).
5. Verify **AC-10** (missing model → red at step 13).
6. Verify **AC-19** (failing test not masked by JSON output).
7. Verify **AC-20/AC-21** (pass counts enforced; inject skip or reduce threshold on branch → red).
8. Merge to `main`; wait for first green `main` run.
9. Execute **§8.2 branch protection rollout** (AC-24).
10. Update README, matrix, accepted.md, 009a dependency doc, technical-risks.
11. Mark `tasks/stab-real-inference-ci-gate.md` complete with commit hash.

---

## 10. Rollback

| Level | Action |
|-------|--------|
| L1 | Remove integration jobs from required checks; keep advisory |
| L2 | Revert workflow + testdata/ci |
| L3 | Restore optional `vars.TB_TEST_GGUF_PATH` step in `inference-cgo` (not recommended) |

---

## 11. Final Recommendation

**APPROVED FOR IMPLEMENTATION**

v3 incorporates operational readiness findings M1–M6. No architecture redesign. Implementation may proceed on workflow + `testdata/ci/` + documentation only. No Go, runtime, inference, or CUDA changes.

| Gate | Status |
|------|--------|
| v2 architecture (split jobs, pin verify, deterministic GGUF) | Approved |
| Operational readiness (M1–M6) | Incorporated |
| Bootstrap job deferral | Accepted v1 |
| Implementation blockers | **None** |

---
**Layer:** planning  
**Related:** [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md), [tasks/stab-real-inference-ci-gate.md](../../tasks/stab-real-inference-ci-gate.md)
