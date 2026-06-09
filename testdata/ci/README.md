# CI test fixtures — Real Inference CI Gate (STAB-001)

Deterministic assets for merge-blocking integration jobs `inference-integration` and `inference-integration-runtime`. See [real-inference-ci-gate-architecture-review.md](../../planning/decisions/real-inference-ci-gate-architecture-review.md) (v3).

## Canonical CI model

| Field | Value |
|-------|-------|
| Model | SmolLM2-135M-Instruct Q4_K_M |
| Quantization | Q4_K_M |
| Approx size | ~105 MB (< 250 MB limit) |
| License | Apache-2.0 (upstream HuggingFaceTB/SmolLM2-135M-Instruct) |
| Filename | `smollm2-135m-instruct-q4_k_m.gguf` (runner cache name) |
| Source file on Hugging Face | `SmolLM2-135M-Instruct-Q4_K_M.gguf` |

### Download URL

CI downloads from:

`https://huggingface.co/unsloth/SmolLM2-135M-Instruct-GGUF/resolve/main/SmolLM2-135M-Instruct-Q4_K_M.gguf`

**Note:** The v3 architecture review named `HuggingFaceTB/SmolLM2-135M-Instruct-GGUF`; that repository is not publicly resolvable (401/404) as of 2026-06-09. The unsloth mirror hosts the same filename, is derived from HuggingFaceTB/SmolLM2-135M-Instruct, and matches the SHA256 committed here. Any URL change must follow **M4** below.

The GGUF binary is **not** stored in git (see root `.gitignore` `*.gguf`). Actions cache + runner temp only.

## SHA256 source of truth

File: `smollm2-135m-instruct-q4_k_m.gguf.sha256` — single-line lowercase hex.

Workflow step 12 runs `sha256sum -c` on every integration job, including cache hits.

### M4 — GGUF URL + SHA256 coupled refresh

1. Download candidate GGUF from the target Hugging Face repo (or mirror).
2. Compute SHA256 locally (`sha256sum` / `Get-FileHash -Algorithm SHA256`).
3. Write hex to `smollm2-135m-instruct-q4_k_m.gguf.sha256`.
4. Update `CI_GGUF_URL` in `.github/workflows/ci.yml` **only if** the URL changed.
5. Bump `CI_GGUF_CACHE_KEY` version suffix if filename or checksum changes.
6. Commit `.sha256` and URL/env changes in the **same PR**.
7. **Never** update URL alone without verifying and updating checksum — `/resolve/main/` can drift.

## llama.cpp pin

Human reference: `llama.cpp.pin`. CI enforces:

- Submodule commit: `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`
- Upstream tag at HEAD: `b9553`

### M3 — llama build cache bump

When **any** cmake argument changes in `.github/workflows/ci.yml` (GPU flags, build targets, Release/Debug, etc.):

1. Bump `LLAMA_BUILD_CACHE_KEY` version suffix (e.g. `llama-cpu-b9553-ubuntu-v1` → `v2`).
2. Record the bump in the PR description.
3. Submodule SHA bump alone does **not** require a manual key bump — SHA is part of the cache key.

Before adding Windows/macOS jobs, include `${{ runner.os }}` in the cache key prefix.

## M6 — Cold-cache expectations (v1)

Both integration jobs run **in parallel** with identical bootstrap.

| Phase | Cold (per job) | Warm (typical) |
|-------|----------------|----------------|
| Checkout + pin verify | ~30–60s | ~30–60s |
| llama.cpp compile | ~3–8 min | cache hit, seconds |
| GGUF download + SHA256 | ~60–120s | cache hit, seconds |
| Integration tests | ~30–90s | ~30–90s |
| **Total** | ~5–12 min | ~2–4 min |

On cold cache, each parallel job may download ~105 MB and compile llama.cpp once. Actions cache amortizes cost on `main`. Fork PRs have read-only cache — expect cold behavior.

## What green CI proves

When `test`, `inference-cgo`, `inference-integration`, and `inference-integration-runtime` are all green on `main`:

- Pure Go packages pass with stub inference (`test`).
- llama.cpp CPU adapter compiles and passes CGO unit tests (`inference-cgo`).
- Submodule pinned at tag `b9553` / commit `9e3b928…`.
- Checksum-verified GGUF loads and generates on Ubuntu CPU (5 llama + 1 runtime integration tests).
- No silent `t.Skip` in integration jobs (pass/skip/run counts enforced).

Does **not** prove: Windows/macOS CGO, GPU/CUDA, latency SLAs, or Hugging Face availability without cache.

## Local integration (mirrors CI)

```bash
git submodule update --init --recursive third_party/llama.cpp
# build llama.cpp (see README.md)
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
export TB_TEST_GGUF_PATH=/path/to/smollm2-135m-instruct-q4_k_m.gguf
CGO_ENABLED=1 go test -tags integration ./internal/inference/llama/...
CGO_ENABLED=1 go test -tags integration ./internal/runtime/...
```

Verify local file against committed checksum:

```bash
sha256sum -c testdata/ci/smollm2-135m-instruct-q4_k_m.gguf.sha256
# (create line: "<hash>  /path/to/file" first)
```

## Troubleshooting

| Symptom | Check |
|---------|--------|
| Submodule tag mismatch | Run `git fetch --tags` inside `third_party/llama.cpp` |
| SHA256 failure after URL change | Re-run M4 coupled refresh |
| Stale llama build after cmake flag change | Bump `LLAMA_BUILD_CACHE_KEY` (M3) |
| All tests skipped | `TB_TEST_GGUF_PATH` unset — workflow bug; M2 should fail |
