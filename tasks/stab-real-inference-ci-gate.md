# Task STAB-001 — Real Inference CI Gate

## Status

In Progress — workflow + `testdata/ci/` + docs implemented; pending first green CI run + branch protection (AC-24)

## Goal

Implement [real-inference-ci-gate-architecture-review.md](../planning/decisions/real-inference-ci-gate-architecture-review.md) **v3** — green CI proves real GGUF inference on Linux CPU. Workflow + CI assets + docs only.

## Context

Current CI silently skips integration tests when `TB_TEST_GGUF_PATH` is unset. Stabilization sprint blocker before any CUDA or feature work. v3 adds mandatory pipefail (M1) and pass-count enforcement (M2) from operational readiness review.

## Requirements

- Merge-blocking jobs: `inference-integration`, `inference-integration-runtime`
- Keep `test`, `inference-cgo`; remove runtime integration step from `inference-cgo`
- Submodule SHA + tag `b9553` verification
- Deterministic GGUF provision + SHA256 gate (< 250 MB)
- **M1:** Pipefail-safe integration test execution
- **M2:** Minimum pass counts (llama: 5, runtime: 1); fail on skip > 0 or run count = 0
- **M3–M6:** Operational procedures in `testdata/ci/README.md`
- **M5:** Branch protection only after first green `main` run
- No Go / runtime / inference code changes

## Files

See architecture review v3 §1–§2.

## Acceptance Criteria

### Core (v2)

- [ ] AC-1: Merge-blocking `inference-integration` job
- [ ] AC-2: Separate `inference-integration-runtime` job
- [ ] AC-3: Keep `test` + `inference-cgo` purpose unchanged
- [ ] AC-4: Remove runtime integration from `inference-cgo`
- [ ] AC-5: Submodule SHA verified
- [ ] AC-6: Upstream tag `b9553` verified
- [ ] AC-7: CPU llama.cpp libraries built
- [ ] AC-8: Intentionally broken SHA256 fails CI
- [ ] AC-9: Deterministic model provision (no repo var dependency)
- [ ] AC-10: Missing model fails job
- [ ] AC-11: 5 llama integration tests pass
- [ ] AC-12: Runtime E2E test passes
- [ ] AC-14: Model < 250 MB
- [ ] AC-15: No Go feature changes
- [ ] AC-16: Docs updated (README, matrix, accepted, 009a, testdata/ci/README)
- [ ] AC-17: `CGO_ENABLED=0 go test ./...` green
- [ ] AC-18: `inference-cgo` green

### Operational (v3 — M1–M6)

- [ ] AC-13: No silent skip (skip count = 0, run count > 0)
- [ ] AC-19: Pipefail-safe execution; `go test` failure not masked by JSON/tee
- [ ] AC-20: Llama JSON pass count ≥ 5
- [ ] AC-21: Runtime JSON pass count ≥ 1
- [ ] AC-22: M3 llama build cache bump procedure in `testdata/ci/README.md`
- [ ] AC-23: M4 GGUF URL + SHA256 coupled refresh in `testdata/ci/README.md`
- [ ] AC-24: Branch protection configured per §8.2 after first green `main`
- [ ] AC-25: M6 cold-cache expectations documented in `testdata/ci/README.md`

## Implementation Order

1. Create `testdata/ci/` (checksum, pin, README with M3/M4/M6).
2. Update `.github/workflows/ci.yml` (M1, M2, integration jobs, trim `inference-cgo`).
3. Green on feature branch (four jobs).
4. Verify AC-8, AC-10, AC-19, AC-20/AC-21 (negative tests on branch).
5. Merge to `main`; first green run.
6. Branch protection rollout (§8.2).
7. Doc sync (README, matrix, accepted, 009a, technical-risks).
8. Complete task entry in `completed.md`.

## Out Of Scope

- CUDA, GPU, inference implementation, runtime changes
- Windows/macOS CI jobs
- Bootstrap job (deferred v1)

## Related

- Review: [planning/decisions/real-inference-ci-gate-architecture-review.md](../planning/decisions/real-inference-ci-gate-architecture-review.md) (v3)

---
**Layer:** task
