# Testing Policy

**Layer:** docs  
**Audience:** Solo developer and AI agents working in this repository.

---

## Policy

Every important behavior must have a direct test that asserts outputs and errors, not merely successful compilation.

Every architectural boundary must have a negative test that fails CI if forbidden imports or couplings appear (`tests/import_boundary_test.go`).

Every bug fix must include a named regression test (`TestRegression_<ShortName>`).

Parsers, configuration loaders, and CLI inputs must be fuzzed or exhaustively table-tested so malformed input never panics.

Integration tests with real GGUF models are required for inference and runtime changes. Required CI jobs block merges if integration tests skip or fall below documented pass thresholds.

A release is trusted only after all merge-blocking CI jobs pass, [RELEASE-CHECKLIST.md](../planning/releases/RELEASE-CHECKLIST.md) is complete, and any known manual gaps (e.g. CUDA GPU verification) are explicitly recorded—not assumed.

**False confidence is not acceptable.** Fast stub-path green (`CGO_ENABLED=0 go test ./...`) is necessary but never sufficient for shipping inference or runtime changes.

---

## Test types

| Type | When required | Location |
|------|---------------|----------|
| Unit | Every exported behavior | `*_test.go` next to code |
| Integration | Inference, runtime, loader | `-tags integration`, real GGUF |
| Boundary | Import invariants | `tests/import_boundary_test.go` |
| Regression | Every bug fix | `TestRegression_*` in owning package |
| Fuzz | CLI env/config (`cmd/tinybrain`), YAML seeds (`internal/registry`) | `Fuzz*`; run locally with `-fuzztime=30s` |
| Golden | CLI models list, normalized probe JSON, doctor header | `cmd/tinybrain/testdata/golden/`; `TB_UPDATE_GOLDEN=1` to update |
| Benchmark | Optional until task 018 | Local only; not release-blocking |

---

## CI jobs (merge-blocking)

| Job | Proves |
|-----|--------|
| `test` | Full tree, stubs, boundary tests |
| `inference-cgo` | llama adapter with real `libllama.so` |
| `inference-integration` | Real GGUF; ≥5 passes; 0 skips |
| `inference-integration-runtime` | Runtime E2E; ≥1 pass; 0 skips |

Details: [README.md](../README.md), [planning/metrics/repo-health.md](../planning/metrics/repo-health.md).

---

## Related

- [CONTRIBUTING.md](../CONTRIBUTING.md) — workflow
- [docs/architecture/invariants.md](architecture/invariants.md) — boundaries
- [testdata/ci/README.md](../testdata/ci/README.md) — integration fixture
