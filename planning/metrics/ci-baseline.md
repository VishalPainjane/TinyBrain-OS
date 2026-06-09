# CI Duration Baselines

Reference bands for STAB-002 observability. Derived from STAB-001 M6 ([testdata/ci/README.md](../../testdata/ci/README.md)). Update when cache keys, runner images, or test counts change materially.

**Review cadence:** Weekly on `main` — compare last 7 runs in `ci-runs.jsonl` against bands below.

## Integration jobs (per job, parallel)

| Phase | Cold (cache miss) | Warm (cache hit) |
|-------|-------------------|------------------|
| Checkout + pin verify | 30–60 s | 30–60 s |
| llama.cpp compile | 3–8 min | &lt; 15 s |
| GGUF download + SHA256 | 60–120 s | &lt; 10 s |
| Integration tests (llama) | 30–90 s | 30–90 s |
| Integration tests (runtime) | 30–60 s | 30–60 s |
| **Job total (warm)** | — | **2–4 min** |
| **Job total (cold)** | — | **5–12 min** |

## `inference-cgo` (no llama Actions cache)

| Phase | Typical |
|-------|---------|
| llama.cpp compile | 3–8 min |
| CGO unit tests | 15–60 s |
| **Job total** | **4–10 min** |

## `test` (CGO disabled)

| Phase | Typical |
|-------|---------|
| `go test ./...` | 30–90 s |
| **Job total** | **1–2 min** |

## Workflow wall time (four jobs parallel)

| Condition | Typical |
|-----------|---------|
| Warm caches on `main` | 4–10 min (dominated by `inference-cgo` compile) |
| Cold integration caches | 5–12 min (integration jobs); `inference-cgo` still compiles |

## Reliability targets (v1 qualitative)

| Signal | Healthy | Investigate |
|--------|---------|-------------|
| `main` success rate (7d) | 100% | &lt; 100% |
| Integration cache hit rate (warm `main`) | &gt; 90% | &lt; 70% after stable week |
| `failure_phase=llama_build` + `SIGILL` | 0 | Any occurrence → check portable flags / cache key |
| `m2_skip_count` &gt; 0 | 0 | Any → STAB-001 regression |
| Job duration vs warm band | within band | &gt; p95 for 3 consecutive runs |

## Triage quick reference

| Symptom | Likely cause | Action |
|---------|--------------|--------|
| `SIGILL` in integration test output | Non-portable cached `libllama.so` | Verify `LLAMA_CMAKE_PORTABLE`; bump `LLAMA_BUILD_CACHE_KEY` |
| `pass_count` below threshold, `go test exit 0` | M2 guard | Inspect uploaded JSON artifact |
| Sudden GGUF download on every run | M4 checksum or cache key change | Verify `.sha256` + `CI_GGUF_CACHE_KEY` coupled |
| `inference-cgo` slow every run | Expected (no cache) | Do not compare to integration warm times |
| Collector stopped appending JSONL | Bot PR failure | Fix collector; history gap acceptable |

## Baseline update procedure

1. Run `jq` or manual scan of last 20 `main` rows in `ci-runs.jsonl`.
2. If warm p95 exceeds band by &gt; 20% for 2 weeks, widen band and note reason in this file.
3. If cold runs exceed 15 min routinely, open RFC for bootstrap job (out of STAB-002 scope).

---
**Layer:** planning  
**Related:** [ci-schema.md](ci-schema.md), [ci-runs.jsonl](ci-runs.jsonl)
