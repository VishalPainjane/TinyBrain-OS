# Manual GPU Verification Checklist — Task 009d

**Purpose:** Sign-off evidence for v0.6 GPU offload when CI has no GPU runners.  
**Required for:** Matrix CUDA `Generate` → **Yes** (vs **Partial** with compile-only).

---

## Prerequisites

- [ ] NVIDIA driver installed (`nvidia-smi` succeeds)
- [ ] CUDA toolkit matching llama.cpp GGML_CUDA build
- [ ] Submodule at pin `b9553` @ `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`
- [ ] llama.cpp built with `-DGGML_CUDA=ON` in `third_party/llama.cpp/build-cuda/`
- [ ] TinyBrain built with `CGO_ENABLED=1 go build -tags cuda ./internal/inference/llama/...`
- [ ] Small GGUF model available (e.g. SmolLM2-135M Q4_K_M)
- [ ] `LD_LIBRARY_PATH` (Linux) or PATH (Windows) includes cuda build bin + CUDA runtime

---

## Config

- [ ] `LlamaConfig.NGLayers = -1` (all layers) or explicit layer count
- [ ] `hardware.Probe()` reports `BackendCUDA` (optional cross-check)

---

## Load / Generate

- [ ] `LoadModel` succeeds without OOM
- [ ] `Generate` returns non-empty output
- [ ] `TokensProduced > 0`
- [ ] Unload succeeds; second load succeeds (no leak)

Optional automated helpers (same binary, requires GPU):

```bash
export TB_CUDA_INTEGRATION=1
export TB_TEST_GGUF_PATH=/path/to/model.gguf
export TB_NGLAYERS=-1   # optional; default -1
CGO_ENABLED=1 go test -tags 'cuda integration' ./internal/inference/llama/... \
  -run 'TestLlamaProvider_LoadGenerate_CUDA|TestLlamaProvider_CUDA_TPS' -v
```

---

## Performance (informational — no SLA)

| Metric | CPU binary (`NGLayers=0`) | CUDA binary (`NGLayers=-1`) |
|--------|----------------------------|----------------------------|
| TTFT (ms) | | |
| TPS | | |
| Model | | |
| GPU model | | |
| Date / machine | | |

- [ ] CUDA TPS ≥ CPU TPS for same prompt/model (expected on discrete GPU)

---

## Platforms (check at least one)

- [ ] Linux native or WSL2
- [ ] Windows (optional)

---

## Sign-off

| Field | Value |
|-------|-------|
| Verified by | |
| Date | |
| Commit | |
| Notes | |

---
**Layer:** planning
