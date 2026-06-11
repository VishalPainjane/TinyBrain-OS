# Repository Health

Trust and verification status for solo development. Update at task completion and on every release.

**Live compass:** [docs/current.md](../../docs/current.md)  
**Version progress:** [progress.md](progress.md)

---

## Last verified

| Field | Value |
|-------|-------|
| Last shipped tag | `v0.6` |
| Sprint | Product CLI — complete (`stab-003`) |
| Active task | (none) |
| Uncommitted | v0.6 docs + CLI + hygiene (3 logical commits) |
| `go test ./...` (local) | Pass (`CGO_ENABLED=0`, 2026-06-11) |
| Boundary tests | `tests/import_boundary_test.go` — INV-001, INV-002, INV-008 |

---

## CI jobs (merge-blocking on `main`)

| Job | Tier | What it proves |
|-----|------|----------------|
| `test` | Fast | `CGO_ENABLED=0 go test ./...` — compile, stubs, boundaries |
| `inference-cgo` | CGO | llama adapter unit tests with `libllama.so` |
| `inference-integration` | Integration | Real GGUF; ≥5 passes; **0 skips** |
| `inference-integration-runtime` | Integration | Runtime orchestration; ≥1 pass; **0 skips** |

CI baselines: [ci-baseline.md](ci-baseline.md)  
Fixture docs: [testdata/ci/README.md](../../testdata/ci/README.md)

---

## Test tiers (local)

| Tier | Command | When |
|------|---------|------|
| Fast | `CGO_ENABLED=0 go test ./...` | Every commit |
| Integration | `CGO_ENABLED=1 go test -tags integration -count=1 ./internal/inference/llama/...` (+ runtime) | Inference/runtime changes |
| Manual GPU | 009d checklist | CUDA releases |

---

## Known gaps (not false confidence)

| Gap | Status | Mitigation |
|-----|--------|------------|
| CUDA GPU runtime | **Partial** — compile path exists; runtime unverified in CI | [009d manual checklist](../decisions/009d-manual-gpu-checklist.md) |
| Metal / ROCm / Vulkan | Not implemented | Deferred per v0.6 release |
| Benchmark suite | Not implemented | Task 018 / v0.8 |
| Event JSON fuzz | Not implemented | No wire unmarshaler yet |
| Fuzz tests | **Partial** — CLI env, registry YAML | `go test -fuzz=... -fuzztime=30s` |
| Golden CLI tests | **Done** — models list, probe JSON (normalized), doctor header | `TB_UPDATE_GOLDEN=1` to regenerate |
| E2E multi-agent | Not implemented | Forbidden until roadmap unlocks |

---

## What “green” means

**Proven:** Core packages compile and unit-test on stub path; import boundaries hold; llama CPU load/generate works on pinned GGUF in CI; runtime integration path works in CI.

**Not proven:** GPU inference in CI; cross-platform native builds beyond Ubuntu CI; performance; agent swarm behavior.

---

## Policy

[docs/testing-policy.md](../../docs/testing-policy.md) | [CONTRIBUTING.md](../../CONTRIBUTING.md)

---

**Layer:** planning  
**Last updated:** 2026-06-11
