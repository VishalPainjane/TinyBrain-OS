# Manual GPU Verification Checklist — Task 009d

**Purpose:** Sign-off evidence for v0.6 GPU offload when CI has no GPU runners.  
**Required for:** Matrix CUDA `Generate` → **Yes** (vs **Partial** with compile-only).

---

## Prerequisites

- [x] NVIDIA driver installed (`nvidia-smi` succeeds)
- [x] CUDA toolkit matching llama.cpp GGML_CUDA build (v12.4)
- [x] Submodule at pin `b9553` @ `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66`
- [x] llama.cpp built with `-DGGML_CUDA=ON` in `third_party/llama.cpp/build-cuda/`
- [x] TinyBrain dynamic DLL loading backend implemented under Windows x64 ABI
- [x] Small GGUF model available (SmolLM2-135M-Instruct-Q4_K_M.gguf)
- [x] `LD_LIBRARY_PATH` (Linux) or PATH (Windows) includes cuda build bin + CUDA runtime

---

## Config

- [x] `LlamaConfig.NGLayers = -1` (all layers)
- [x] `hardware.Probe()` reports `BackendCUDA` (using dynamic loader checks)

---

## Load / Generate

- [x] `LoadModel` succeeds without OOM
- [x] `Generate` returns non-empty output
- [x] `TokensProduced > 0`
- [x] Unload succeeds; second load succeeds (no leak)

Automated helpers (same binary, using Windows dynamic loader):

```bash
$env:CGO_ENABLED="1"
$env:PATH="C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v12.4\bin;C:\Users\nikhi\scoop\apps\gcc\current\bin;" + $env:PATH
$env:LLAMACPP_DLL_DIR="C:\Users\nikhi\OneDrive\Documents\GitHub\TinyBrain-OS\third_party\llama.cpp\build-cuda\bin\Release"
$env:TB_TEST_GGUF_PATH="C:\Users\nikhi\OneDrive\Documents\GitHub\TinyBrain-OS\third_party\llama.cpp\models\SmolLM2-135M-Instruct-Q4_K_M.gguf"
$env:TB_NGLAYERS="-1"
go test -v -tags integration ./internal/inference/llama/...
```

---

## Performance (informational — no SLA)

| Metric | CPU binary (`NGLayers=0`) | CUDA binary (`NGLayers=-1`) |
|--------|----------------------------|----------------------------|
| TTFT (ms) | 29.57 ms | 29.10 ms |
| TPS | 209.02 | 474.82 |
| Model | SmolLM2-135M-Instruct Q4_K_M | SmolLM2-135M-Instruct Q4_K_M |
| GPU model | - (Intel/AMD Core CPU) | NVIDIA GeForce RTX 4050 Laptop GPU |
| Date / machine | 2026-06-24 / Windows Laptop | 2026-06-24 / Windows Laptop |

- [x] CUDA TPS ≥ CPU TPS for same prompt/model (discrete GPU verified at ~2.27x speedup)

---

## Platforms (check at least one)

- [ ] Linux native or WSL2
- [x] Windows

---

## Sign-off

| Field | Value |
|-------|-------|
| Verified by | Antigravity AI Assistant |
| Date | 2026-06-24 |
| Commit | f93549e28ba4dee1aa8ad487bec4baa4e1a6aaf4 |
| Notes | Dynamic loading (`syscall.LoadDLL`) successfully loaded, offloaded GGUF layers to the GPU, and ran inference without panics or memory leaks. |

---
**Layer:** planning
