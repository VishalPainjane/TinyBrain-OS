# Architecture Review — Task 009b (CPU Generate)

**Status:** Complete — implementation and integration verified (2026-06-09)  
**Date:** 2026-06-09  
**Scope:** `internal/inference/llama/` — real CPU `Generate` via CGO. No runtime, loader, registry, or scheduler changes.

| Gate | Status |
|------|--------|
| Architecture review | Approved |
| API compatibility (b9553 @ 9e3b928) | Approved |
| Code implementation | **Complete** |
| Linux CGO verification | **Complete** |
| Real GGUF integration | **Complete** |

**Pre-read:** [009a-architecture-review.md](009a-architecture-review.md), [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md), [inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md), [009a-build-tags.md](009a-build-tags.md), [009a-registry-resolver.md](009a-registry-resolver.md), [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md)

**Submodule pin (verified 2026-06-09):**

| Field | Value |
|-------|-------|
| Tag | `b9553` |
| SHA | `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66` |
| Header source | `third_party/llama.cpp/include/llama.h` |

---

## Gates Cleared

| Gate | Document | Status |
|------|----------|--------|
| A — Lifecycle | [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md) | Closed |
| B — Resolver | [009a-registry-resolver.md](009a-registry-resolver.md) | Unchanged |
| C — Dependency | [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md) | Unchanged |
| D — Build tags | [009a-build-tags.md](009a-build-tags.md) | Unchanged |
| E — API pin | This document § API Compatibility Report | Closed |

---

## Approved API Corrections (Binding)

These corrections are mandatory for 009b implementation. Verified against `llama.h` at SHA `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`.

| # | Rule |
|---|------|
| 1 | Use `llama_memory_clear(llama_get_memory(ctx), true)` — **not** `llama_kv_cache_clear` (absent at pin) |
| 2 | Obtain vocab via `llama_model_get_vocab(model)` before `llama_tokenize` / `llama_token_to_piece` |
| 3 | EOS check via `llama_vocab_is_eog(vocab, token)` — **not** deprecated `llama_token_is_eog` |
| 4 | Sampler chain APIs only: `llama_sampler_chain_init` → `llama_sampler_chain_add` → `llama_sampler_sample` |
| 5 | Create `llama_context` via `llama_init_from_model` during `LoadModel` (not lazy on first Generate) |
| 6 | On unload: `llama_free(ctx)` **before** `llama_model_free(model)` |
| 7 | Store contexts in `nativeContexts` map alongside `nativeModels` |

**Implementation constraints (unchanged):**

- No runtime, scheduler, loader, registry, contract, or architecture changes
- No CUDA execution, no runtime wiring, no 009c work

---

## 009a Assumption Verification

| 009a assumption | Status | Evidence |
|-----------------|--------|----------|
| `ModelResolver` port; no registry in `provider.go` / `load_*.go` | **Valid** | `resolver.go`, `registry_resolver.go` |
| Submodule pin `b9553` @ `9e3b928…` | **Valid** | `git -C third_party/llama.cpp rev-parse HEAD` |
| CPU load/unload via CGO mmap; `NGLayers=0` | **Valid** | `bindings_cgo.go`, `load_cpu.go` |
| `Generate` stub returns `ErrNotImplemented` | **Valid** | `port_stubs.go` |
| Core packages unchanged; runtime unwired to `LlamaProvider` | **Valid** | `runtime.go` delegates to port only |
| `CGO_ENABLED=0 go test ./...` green | **Valid** | CI `test` job |
| `inference-cgo` CI job green | **Valid** | `.github/workflows/ci.yml` |
| Build tags mutually exclusive | **Valid** | `load_cpu.go`, `provider_stub.go` |
| INV-008 satisfied | **Valid** | `tests/import_boundary_test.go` |
| ACTIVE = ready for Generate with valid context | **Gap (009b)** | 009a loads `llama_model` only; no `llama_context` |

**Deviations carried forward (not blockers):**

1. `bindings_cgo.go` uses `cgo && !cuda && !rocm && !metal && !vulkan` (not bare `cgo`).
2. Port stubs merged into `port_stubs.go` (not separate `generate_stub.go`).
3. Native handles in package-level `nativeModels` map; 009b adds `nativeContexts` map.

---

## API Compatibility Report (Pin b9553 @ 9e3b928)

Verified against `third_party/llama.cpp/include/llama.h` at the pinned SHA. **Do not use upstream or deprecated APIs.**

### Summary

| Area | Planned API (pre-pin review) | Pinned API (b9553) | Verdict |
|------|------------------------------|--------------------|---------|
| Tokenizer | `llama_tokenize` | Present — requires `llama_vocab *` | **Adjust** — obtain vocab via `llama_model_get_vocab` |
| Sampler | `llama_sampler_*` (unspecified) | Chain + `llama_sampler_sample` | **Confirmed** — no legacy `llama_sample` |
| Context create | `llama_init_from_model` | Present | **Confirmed** |
| KV / memory clear | `llama_kv_cache_clear` | **Absent** | **Replace** with `llama_memory_clear` |
| Decode | `llama_decode` + batch | Present | **Confirmed** |
| Detokenize | `llama_token_to_piece` | Present | **Confirmed** |
| EOS | `llama_token_is_eog` | Deprecated | **Use** `llama_vocab_is_eog` |
| Model free | `llama_model_free` | Present (009a already uses) | **Confirmed** |
| Context free | `llama_free` | Present | **Confirmed** |

### 1. Tokenizer API

**Present at pin:**

```c
LLAMA_API const struct llama_vocab * llama_model_get_vocab(const struct llama_model * model);

LLAMA_API int32_t llama_tokenize(
    const struct llama_vocab * vocab,
    const char * text,
    int32_t text_len,
    llama_token * tokens,
    int32_t n_tokens_max,
    bool add_special,
    bool parse_special);

LLAMA_API int32_t llama_token_to_piece(
    const struct llama_vocab * vocab,
    llama_token token,
    char * buf,
    int32_t length,
    int32_t lstrip,
    bool special);

LLAMA_API int32_t llama_detokenize(
    const struct llama_vocab * vocab,
    const llama_token * tokens,
    int32_t n_tokens,
    char * text,
    int32_t text_len_max,
    bool remove_special,
    bool unparse_special);
```

**Rules for 009b:**

- Resolve vocab once per loaded model: `llama_model_get_vocab(model)`.
- `text_len` may be `-1` for null-terminated C strings (per common llama.cpp usage; verify at implementation).
- Negative return from `llama_tokenize` indicates required buffer size — allocate and retry.
- `add_special=true` for chat-style prompts on instruction models; default `true` for 009b single-prompt path.
- Prefer `llama_token_to_piece` per token in the decode loop (incremental output).
- **Do not use** deprecated `llama_token_*` helpers.

### 2. Sampler API

**Present at pin — chain-based sampler (no `llama_sample`):**

```c
LLAMA_API struct llama_sampler_chain_params llama_sampler_chain_default_params(void);
LLAMA_API struct llama_sampler * llama_sampler_chain_init(struct llama_sampler_chain_params params);
LLAMA_API void llama_sampler_chain_add(struct llama_sampler * chain, struct llama_sampler * smpl);
LLAMA_API llama_token llama_sampler_sample(struct llama_sampler * smpl, struct llama_context * ctx, int32_t idx);
LLAMA_API void llama_sampler_free(struct llama_sampler * smpl);

// Terminal samplers (pick one):
LLAMA_API struct llama_sampler * llama_sampler_init_greedy(void);
LLAMA_API struct llama_sampler * llama_sampler_init_dist(uint32_t seed);
// Optional chain members:
LLAMA_API struct llama_sampler * llama_sampler_init_temp(float t);
LLAMA_API struct llama_sampler * llama_sampler_init_top_p(float p, size_t min_keep);
```

**Official sample loop (from `llama.h` lines 1199–1228):**

1. Build chain at start: `chain_init` → `chain_add(temp)` → `chain_add(dist|greedy)`.
2. Each decode step: `llama_decode(ctx, batch)`.
3. Sample: `llama_sampler_sample(smpl, ctx, -1)` — `idx=-1` = last output.
4. `llama_sampler_sample` internally applies + accepts token.
5. Free chain after generation.

**009b recommendation:**

| Mode | Chain |
|------|-------|
| Integration test (deterministic) | `greedy` only |
| Default generation | `temp(cfg.Temperature)` → `dist(seed)` |

- Create sampler chain per `Generate` call; free before return.
- Call `llama_sampler_reset` only if reusing chain across calls (not needed per-call pattern).

### 3. Context Creation API

**Present at pin:**

```c
LLAMA_API struct llama_context_params llama_context_default_params(void);
LLAMA_API struct llama_context * llama_init_from_model(
    struct llama_model * model,
    struct llama_context_params params);
LLAMA_API void llama_free(struct llama_context * ctx);
```

**`llama_context_params` fields used in 009b:**

| Field | Source | Notes |
|-------|--------|-------|
| `n_ctx` | `LlamaConfig.ContextSize` | Query actual via `llama_n_ctx(ctx)` after create |
| `n_threads` | `LlamaConfig.Threads` | Also set via `llama_set_n_threads` if needed |
| `n_batch` | `min(n_ctx, 512)` or config | Must be ≥ prompt chunk size for prefill |
| `n_ubatch` | default or `= n_batch` | Physical batch cap |

**Deprecated — do not use:** `llama_new_context_with_model`.

**Lifecycle (009b):**

1. `llama_model_load_from_file` (009a, unchanged).
2. `llama_init_from_model(model, ctx_params)` — **new in 009b**, end of load path.
3. Store context pointer in `nativeContexts[modelID]`.
4. On unload: `llama_free(ctx)` then `llama_model_free(model)`.

### 4. KV / Memory Clear API

**`llama_kv_cache_clear` does not exist at pin b9553.**

**Replacement:**

```c
LLAMA_API llama_memory_t llama_get_memory(const struct llama_context * ctx);

LLAMA_API void llama_memory_clear(llama_memory_t mem, bool data);
```

**009b usage at start of each `Generate`:**

```c
llama_memory_clear(llama_get_memory(ctx), true);
```

- `data=true` clears data buffers and metadata (full reset for new prompt).
- Sequence ID defaults to `0` via `llama_batch_get_one`.

**Alternative (not needed 009b):** `llama_memory_seq_rm(mem, seq_id, p0, p1)` for partial clears.

### 5. Decode / Batch API

**Present at pin:**

```c
LLAMA_API struct llama_batch llama_batch_get_one(llama_token * tokens, int32_t n_tokens);
LLAMA_API int32_t llama_decode(struct llama_context * ctx, struct llama_batch batch);
```

**`llama_decode` return codes:**

| Value | Meaning |
|-------|---------|
| `0` | Success |
| `1` | No KV slot — reduce batch or increase context |
| `2` | Aborted |
| `-1` | Invalid batch |
| `< -1` | Fatal |

**Prefill strategy:**

- If `len(prompt_tokens) <= n_batch`: single `llama_batch_get_one(all_tokens)` + `llama_decode`.
- If longer: chunk into batches of `n_batch` tokens (loop).
- Positions tracked automatically when `batch.pos == NULL` (default for `llama_batch_get_one`).

**Generation loop:**

- Each new token: `llama_batch_get_one(&token, 1)` → `llama_decode` → `llama_sampler_sample(smpl, ctx, -1)`.
- Stop on `llama_vocab_is_eog(vocab, token)` or `MaxTokens` reached.

### 6. EOS Detection

```c
LLAMA_API bool llama_vocab_is_eog(const struct llama_vocab * vocab, llama_token token);
LLAMA_API llama_token llama_vocab_eos(const struct llama_vocab * vocab);
```

Use `llama_vocab_is_eog` in the decode loop.

### 7. Timing (TTFT / TPS)

```c
LLAMA_API int64_t llama_time_us(void);
```

Available for integration test metrics. No telemetry package in 009b.

---

## Required Adjustments to Pre-Pin 009b Review

| # | Original plan | Adjustment (b9553) |
|---|---------------|-------------------|
| 1 | `llama_kv_cache_clear` before generate | Use `llama_memory_clear(llama_get_memory(ctx), true)` |
| 2 | Tokenize against model handle | Get `vocab = llama_model_get_vocab(model)` first |
| 3 | Legacy sampler assumptions | Chain API only: `chain_init` → `chain_add` → `sampler_sample` |
| 4 | Context created lazily on first Generate | Create at load end via `llama_init_from_model` (ACTIVE = context valid) |
| 5 | `llama_token_is_eog` | Use `llama_vocab_is_eog` |
| 6 | Single prefill assumption | Handle `llama_decode` return `1`; chunk if prompt > `n_batch` |
| 7 | Native storage | Add `nativeContexts map[string]unsafe.Pointer`; unload context before model |
| 8 | Integration test sampler | Use `llama_sampler_init_greedy()` for reproducibility |
| 9 | `bindings_cgo.go` scope | Extend with context create/destroy, memory clear, tokenize, decode, sampler helpers |
| 10 | Config `n_batch` | Add `BatchSize` to `LlamaConfig` or derive from `ContextSize` |

**Unchanged from approved review:**

- Scope limited to `internal/inference/llama/`
- No runtime / loader / registry / scheduler changes
- `ModelResolver` injection unchanged
- Build tags unchanged
- `SaveContext` / `RestoreContext` remain stubs
- Runtime wiring deferred to 009c

---

## Files to Create

| File | Build constraint | Purpose |
|------|------------------|---------|
| `internal/inference/llama/generate_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` | `generateBackend()` orchestration |
| `internal/inference/llama/generate_stub.go` | `!cgo` | `generateBackend()` → `ErrCGODisabled` |
| `internal/inference/llama/generate_integration_test.go` | `cgo && integration` | Real GGUF generate + TTFT/TPS |

---

## Files to Modify

| File | Change |
|------|--------|
| `internal/inference/llama/port_stubs.go` | Remove `Generate`; keep Save/Restore stubs |
| `internal/inference/llama/provider.go` | Add `Generate()` |
| `internal/inference/llama/bindings_cgo.go` | Context create/free; memory clear; tokenize; decode; sampler; extend unload order |
| `internal/inference/llama/context.go` | Optional Go-side metadata |
| `internal/inference/llama/config.go` | Add `MaxTokens`, `Temperature`, `BatchSize`, `Seed` |
| `internal/inference/llama/errors.go` | Add `ErrGenerationFailed` |
| `internal/inference/llama/provider_test.go` | Generate unit tests; update port stub test |
| `internal/inference/llama/load_integration_test.go` | Optional: reference generate test |
| `.github/workflows/ci.yml` | Optional integration generate when `TB_TEST_GGUF_PATH` set |
| `README.md` | Generate integration test docs |
| `docs/architecture/inference-backend-matrix.md` | CPU `Generate` → **Partial** on merge |
| `docs/specs/v0.6-inference.md` | Checkbox sync on merge |
| Planning sync | `completed.md`, `docs/current.md`, `architecture-evolution/current-state.md` |

**Not modified:** `internal/runtime/`, `internal/loader/`, `internal/registry/`, `internal/scheduler/`, `docs/contracts/*`.

---

## Generate() Lifecycle Design (Final)

```text
LoadModel(modelID)
  resolve → stat → loadNativeModel
    llama_backend_init (once)
    llama_model_load_from_file (mmap)
    llama_init_from_model → nativeContexts[modelID]   ← NEW
  models[modelID] = slot → ACTIVE

Generate(req)
  require models[req.ModelID] → else runtime.ErrModelNotLoaded
  llama_memory_clear(llama_get_memory(ctx), true)      ← NOT kv_cache_clear
  vocab = llama_model_get_vocab(model)
  tokens = llama_tokenize(vocab, prompt, ...)
  prefill: llama_batch_get_one + llama_decode (chunked if needed)
  sampler chain: temp → dist (or greedy in tests)
  loop: sample → token_to_piece → append
        batch_get_one(single) → decode
        until is_eog or MaxTokens
  return GenerateResponse

UnloadModel(modelID)
  llama_free(ctx) → llama_model_free(model) → delete maps
```

---

## Token Generation Flow (Final)

```text
GenerateRequest { ModelID, Prompt }
        │
        ▼
  provider.go: lock, loaded check
        │
        ▼
  generate_cpu.go
        │
        ├─ llama_memory_clear(get_memory(ctx), true)
        ├─ vocab = llama_model_get_vocab(model)
        ├─ llama_tokenize(vocab, prompt, add_special=true)
        ├─ llama_batch_get_one(prompt_tokens) → llama_decode  [prefill, TTFT ends]
        ├─ chain = sampler_chain_init → add(temp) → add(dist|greedy)
        └─ loop:
             token = llama_sampler_sample(chain, ctx, -1)
             piece = llama_token_to_piece(vocab, token)
             if llama_vocab_is_eog(vocab, token): break
             llama_batch_get_one(&token,1) → llama_decode
        │
        ▼
  GenerateResponse { ModelID, Output, TokensProduced }
```

---

## Context Ownership Boundaries

| Asset | Owner | Created | Destroyed |
|-------|-------|---------|-----------|
| `ModelDefinition` | Registry | Register | Registry close |
| `ModelSpec` | `ModelResolver` | `Resolve()` | N/A |
| `llama_model` | Inference CGO | `LoadModel` | `UnloadModel` (after context) |
| `llama_context` | Inference CGO | `LoadModel` (post-model) | `UnloadModel` (before model) |
| KV / memory | Inference CGO | First decode | `llama_memory_clear` each Generate |
| Sampler chain | Inference CGO | Each `Generate` | End of `Generate` |
| Loader `ModelState` | Loader | — | Unchanged until 009c |

---

## Runtime Ownership Review

| Responsibility | 009b impact |
|----------------|-------------|
| `Generate` delegation | **None** — runtime already delegates |
| Lifecycle events | **None** — load/unload only |
| Metrics | Inference integration tests only |
| llama.cpp imports | **None** in runtime (INV-008) |

---

## GGUF Integration Strategy

Unchanged from 009a load path. 009b adds context sizing from `LlamaConfig.ContextSize` and validates via integration test with `TB_TEST_GGUF_PATH`.

---

## CI Strategy

| Job | 009b change | Blocker |
|-----|-------------|---------|
| `test` (`CGO_ENABLED=0`) | New unit tests (not-loaded, validation) | Yes |
| `inference-cgo` | Generate validation tests | Yes |
| Integration (`-tags integration`) | Generate + TTFT/TPS when `TB_TEST_GGUF_PATH` set | No (recommended) |

No llama.cpp cmake flag changes required.

---

## Cross-Platform Impact

| OS | 009b |
|----|------|
| Linux | Primary CI |
| Windows | Manual CGO + integration |
| macOS | Manual CPU build |

CPU backend file only; GPU tags unchanged.

---

## Acceptance Criteria

| # | Criterion | Verification |
|---|-----------|--------------|
| AC-1 | Real tokens from loaded GGUF via `Generate` | Integration test |
| AC-2 | Unloaded model → `runtime.ErrModelNotLoaded` | Unit test |
| AC-3 | Non-empty `Output`, `TokensProduced > 0` | Integration test |
| AC-4 | TTFT/TPS measurable | Integration test timing |
| AC-5 | Save/Restore remain `ErrNotImplemented` | Port stub test |
| AC-6 | `CGO_ENABLED=0 go test ./...` passes | CI |
| AC-7 | `inference-cgo` passes | CI |
| AC-8 | Scheduler zero inference imports | Boundary test |
| AC-9 | No runtime/loader/registry/scheduler diffs | Scope check |
| AC-10 | APIs match b9553 pin only | Code review against this doc |
| AC-11 | Matrix CPU `Generate` → Partial on merge | Doc sync |

---

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| `llama_kv_cache_clear` assumed in early plans | High | Use `llama_memory_clear` (documented here) |
| Prompt longer than `n_batch` | Medium | Chunk prefill; handle decode return `1` |
| Context RAM on Tiny profile | Medium | Default `ContextSize=512`, conservative `MaxTokens` |
| Non-deterministic integration tests with `dist` | Low | Greedy sampler in integration tests |
| Windows mmap + generate | Medium | Manual validation |
| Global mutex during Generate | Low | Accept for 009b |
| CGO leak on error paths | Medium | defer frees; unload cycle test |

---

## Test Plan

### Unit (`CGO_ENABLED=0`)

- `Generate_notLoaded` → `ErrModelNotLoaded`
- `Generate_emptyModelID` → validation error
- Save/Restore stubs unchanged

### CGO unit (`inference-cgo`, no GGUF)

- Loaded-state tests deferred to integration (require real model)
- Port stub test updated (Generate no longer `ErrNotImplemented` when CGO+loaded)

### Integration (`cgo && integration`, `TB_TEST_GGUF_PATH`)

| Test | Assertions |
|------|------------|
| `Generate_integration` | Load → Generate → non-empty output |
| `Generate_afterUnload` | `ErrModelNotLoaded` |
| `Generate_TTFT_TPS` | Log timing; no SLA in 009b |
| `LoadGenerateUnload_x3` | No leak across cycles |

Sampler: `greedy` for deterministic output.

---

## Final Implementation Plan (Ordered)

### Phase 1 — Config and errors

1. Add `MaxTokens`, `Temperature`, `BatchSize`, `Seed` to `LlamaConfig` / `DefaultConfig()`.
2. Add `ErrGenerationFailed` to `errors.go`.

### Phase 2 — Native context lifecycle (bindings)

3. Add `nativeContexts` map alongside `nativeModels`.
4. After `llama_model_load_from_file`, call `llama_init_from_model` with `llama_context_default_params()` overrides.
5. Extend `unloadNativeModel`: `llama_free(ctx)` then `llama_model_free(model)`; delete both map entries.
6. Query `llama_n_ctx(ctx)` after create; fail load if context unusable.

### Phase 3 — CGO helpers (bindings_cgo.go)

7. `clearNativeMemory(ctx)` → `llama_memory_clear(llama_get_memory(ctx), true)`.
8. `tokenizePrompt(vocab, text)` → alloc buffer, `llama_tokenize`.
9. `decodeBatch(ctx, tokens)` → `llama_batch_get_one` + `llama_decode` with error mapping.
10. `createSampler(cfg, greedy bool)` → chain init/add/free wrapper.
11. `tokenToPiece(vocab, token)` → `llama_token_to_piece`.
12. `isEOG(vocab, token)` → `llama_vocab_is_eog`.

### Phase 4 — Generate orchestration

13. Implement `generateBackend` in `generate_cpu.go`:
    - memory clear
    - tokenize prompt
    - prefill (chunked)
    - sampler loop with EOS / MaxTokens
    - assemble output string + token count
14. Implement `generate_stub.go` for `!cgo`.
15. Add `Generate` to `provider.go`; remove from `port_stubs.go`.

### Phase 5 — Tests and CI

16. Update `provider_test.go`.
17. Add `generate_integration_test.go`.
18. Optional CI integration step.
19. Update README.

### Phase 6 — Doc sync (merge)

20. `inference-backend-matrix.md`, `v0.6-inference.md`, planning docs.

---

## Implementation Sequence Constraints

- **Do not** call `llama_kv_cache_clear` — it does not exist at pin.
- **Do not** use deprecated `llama_new_context_with_model`, `llama_free_model`, `llama_token_is_eog`.
- **Do not** import registry outside `registry_resolver.go`.
- **Do not** modify runtime, loader, registry, scheduler.
- Verify every CGO symbol against `llama.h` at `9e3b928` before merge.

---

## Exact File Diff Plan

### Create (3 files)

| File | Build tag | Lines (est.) | Purpose |
|------|-----------|--------------|---------|
| `internal/inference/llama/generate_cpu.go` | `cgo && !cuda && !rocm && !metal && !vulkan` | ~80 | `generateBackend(modelID, prompt, cfg)` — Go orchestration calling CGO helpers |
| `internal/inference/llama/generate_stub.go` | `!cgo` | ~10 | `generateBackend` → `ErrCGODisabled` |
| `internal/inference/llama/generate_integration_test.go` | `cgo && integration` | ~120 | Real GGUF generate, TTFT/TPS, cycle tests |

### Modify (8 files — inference package only)

| File | Diff summary |
|------|--------------|
| `internal/inference/llama/bindings_cgo.go` | Add `nativeContexts` map; `initNativeContext(model, modelID, cfg)` after model load; extend `unloadNativeModel` to `llama_free` then `llama_model_free`; add CGO helpers: `clearNativeMemory`, `tokenizeNative`, `decodeNativeBatch`, `sampleNativeToken`, `tokenToPieceNative`, `isNativeEOG`, `freeNativeSampler` |
| `internal/inference/llama/provider.go` | Add `Generate(req)` — lock, `models` check → `runtime.ErrModelNotLoaded`, delegate `generateBackend`, return `GenerateResponse` |
| `internal/inference/llama/port_stubs.go` | **Remove** `Generate` method; keep `SaveContext` / `RestoreContext` stubs |
| `internal/inference/llama/config.go` | Add `MaxTokens uint32`, `Temperature float32`, `BatchSize uint32`, `Seed uint32`; update `DefaultConfig()` |
| `internal/inference/llama/errors.go` | Add `ErrGenerationFailed` |
| `internal/inference/llama/context.go` | Optional: comment that native context lives in `nativeContexts` |
| `internal/inference/llama/provider_test.go` | Replace `TestLlamaProvider_portStubs` Generate assertion; add `TestLlamaProvider_Generate_notLoaded`, `TestLlamaProvider_Generate_emptyModelID`; keep Save/Restore stub tests |
| `internal/inference/llama/load_integration_test.go` | No required change (generate tests in separate file) |

### Modify (docs / CI — post-implementation, merge sync)

| File | Diff summary |
|------|--------------|
| `README.md` | Document `Generate` integration test invocation |
| `.github/workflows/ci.yml` | Optional: run `-tags integration` when `TB_TEST_GGUF_PATH` set |
| `docs/specs/v0.6-inference.md` | Checkboxes for generate + partial llama adapter |
| `docs/architecture/inference-backend-matrix.md` | CPU `Generate` → **Partial** |
| `planning/execution/completed.md` | 009b entry on ship |
| `docs/current.md` | Task status update |
| `planning/architecture-evolution/current-state.md` | Generate status |

### Explicitly zero-diff (enforced)

`internal/runtime/`, `internal/scheduler/`, `internal/loader/`, `internal/registry/`, `docs/contracts/`, `docs/architecture/inference-lifecycle.md`, `cmd/`, `load_cuda.go` (does not exist).

---

## Exact Tests to Add

### `provider_test.go` — always run (`CGO_ENABLED=0` OK)

| Test name | Setup | Action | Assert |
|-----------|-------|--------|--------|
| `TestLlamaProvider_Generate_notLoaded` | Empty provider | `Generate({ModelID:"m1", Prompt:"hi"})` | `errors.Is(err, runtime.ErrModelNotLoaded)` |
| `TestLlamaProvider_Generate_emptyModelID` | Empty provider | `Generate({ModelID:"", Prompt:"hi"})` | non-nil error |
| `TestLlamaProvider_portStubs` | Empty provider | `SaveContext` / `RestoreContext` | `ErrNotImplemented` (Generate case **removed**) |

### `generate_integration_test.go` — `cgo && integration`, requires `TB_TEST_GGUF_PATH`

| Test name | Setup | Action | Assert |
|-----------|-------|--------|--------|
| `TestLlamaProvider_Generate_integration` | Load model via `staticResolver` | `Generate({ModelID, Prompt:"Hello"})` | `err==nil`, `Output!=""`, `TokensProduced>0` |
| `TestLlamaProvider_Generate_afterUnload` | Load → Unload | `Generate` | `errors.Is(err, runtime.ErrModelNotLoaded)` |
| `TestLlamaProvider_Generate_TTFT_TPS` | Load model | `Generate` with `llama_time_us` timing around prefill vs decode | Log TTFT and TPS; no hard threshold |
| `TestLlamaProvider_LoadGenerateUnload_cycles` | Load | Generate → Unload × 3 | No panic; final unload → `ErrModelNotLoaded` |

**Sampler in integration tests:** `llama_sampler_init_greedy()` terminal sampler only (deterministic).

### Unchanged (regression)

| Test | Package |
|------|---------|
| `TestSchedulerDoesNotImportInference` | `tests/import_boundary_test.go` |
| All `runtime_test.go` tests with `StubProvider` | `internal/runtime/` |
| Existing 009a load validation tests | `provider_test.go` |
| `TestLlamaProvider_LoadUnload_integration` | `load_integration_test.go` |

---

## Acceptance Criteria Mapping

| ID | v0.6 spec / invariant | Criterion | Test / verification |
|----|----------------------|-----------|---------------------|
| AC-1 | v0.6: "One real generate call returns tokens from local GGUF" | `Generate` returns non-empty output from real GGUF | `TestLlamaProvider_Generate_integration` |
| AC-2 | Lifecycle: Generate only when ACTIVE | Unloaded model → error | `TestLlamaProvider_Generate_notLoaded`, `TestLlamaProvider_Generate_afterUnload` |
| AC-3 | Contract: `GenerateResponse.TokensProduced > 0` | Token count populated | `TestLlamaProvider_Generate_integration` |
| AC-4 | v0.6: "TTFT and TPS measurable" | Timing recorded in integration test | `TestLlamaProvider_Generate_TTFT_TPS` |
| AC-5 | 009b scope: Save/Restore deferred (011) | Stubs unchanged | `TestLlamaProvider_portStubs` |
| AC-6 | CI rule: `CGO_ENABLED=0` green | Full stub path passes | `go test ./...` CI `test` job |
| AC-7 | 009a CI: inference-cgo green | CGO compile + unit tests | CI `inference-cgo` job |
| AC-8 | INV-001: scheduler zero inference imports | Boundary test passes | `TestSchedulerDoesNotImportInference` |
| AC-9 | Implementation constraint | Zero diffs in forbidden packages | `git diff` scope review |
| AC-10 | API pin b9553 | Only approved symbols in CGO | Code review vs this doc § Approved API Corrections |
| AC-11 | Matrix maintenance | CPU `Generate` → Partial | `inference-backend-matrix.md` on merge |
| AC-12 | INV-008 | llama.cpp only in `internal/inference/` | Import audit |
| AC-13 | 009a: resolver port | No registry import in provider/load/generate files | `grep` / review |
| AC-14 | Input validation | Empty `ModelID` rejected | `TestLlamaProvider_Generate_emptyModelID` |

---

## Risk Checklist (Pre-Implementation)

Review before writing code; re-check before merge.

| # | Risk | Check before code | Check before merge |
|---|------|-------------------|-------------------|
| R-1 | `llama_kv_cache_clear` used by mistake | Grep CGO for `kv_cache` | Same |
| R-2 | Deprecated APIs (`llama_token_is_eog`, `llama_new_context_with_model`) | Grep CGO | Same |
| R-3 | Context not created at load → ACTIVE violation | `initNativeContext` in load path | Integration load+generate passes |
| R-4 | Wrong unload order (model before context) | Code review unload path | Unload integration test |
| R-5 | `nativeContexts` leak on load failure mid-path | Rollback on context create failure | `LoadGenerateUnload_cycles` |
| R-6 | Prompt > `n_batch` without chunking | Chunk logic in prefill | Long-prompt manual test (optional) |
| R-7 | `llama_decode` return `1` unhandled | Map to `ErrGenerationFailed` | Integration test |
| R-8 | CGO string/buffer leaks | `defer C.free` audit | Cycle test |
| R-9 | Mutex held during long generate blocks unload | Document; acceptable 009b | N/A |
| R-10 | Non-deterministic integration with `dist` sampler | Greedy in integration tests only | Test review |
| R-11 | Scope creep into runtime/009c | Zero-diff forbidden packages | `git diff --stat` |
| R-12 | CUDA / GPU code paths | No `load_cuda.go`; `n_gpu_layers=0` | Build without `cuda` tag |
| R-13 | Windows mmap + generate | N/A at code time | Manual if available |
| R-14 | `CGO_ENABLED=0` regression | Unit tests compile without CGO files | CI `test` job |

---
**Layer:** planning  
**Related:** [009a-architecture-review.md](009a-architecture-review.md), [../../docs/specs/v0.6-inference.md](../../docs/specs/v0.6-inference.md)
