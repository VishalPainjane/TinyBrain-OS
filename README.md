# TinyBrain OS

## Vision

TinyBrain OS is a hardware-aware AI runtime kernel for dynamically orchestrating specialized agents on constrained local hardware. Instead of one monolithic model handling every task, TinyBrain treats agents as scheduled processes — loaded when needed, swapped when idle, and coordinated under strict VRAM and RAM budgets.

Local AI fails on consumer hardware when treated as a single-model problem. TinyBrain exists to make it practical: dynamic capability allocation, not static model-centric workflows.

See also: [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)

## Architecture Diagram

```text
User
  ↓
API
  ↓
Router
  ↓
Scheduler
  ↓
Runtime
  ↓
InferenceProvider
  ↓
Models (GGUF / llama.cpp)
```

Core principles:

- Local-first, hardware-aware
- Agents are plugins (registry-defined, not hardcoded)
- Structured JSON IPC between components
- Scheduler independent from inference engine
- Dynamic model swapping under resource budget

Full architecture: [docs/architecture/overview.md](docs/architecture/overview.md)

## Quick Start

**Current state:** v0.5 shipped at tag `v0.5`. v0.6 in progress — 009a–009c complete (CPU inference wired to `ModelRuntime`); GPU offload remaining.

```bash
go test ./...
```

Default local and CI unit path uses `CGO_ENABLED=0` (stub inference). Merge-blocking CI also runs CGO unit and real-GGUF integration jobs — see [CI jobs](#ci-jobs-merge-blocking-on-main) below.

```bash
git submodule update --init --recursive third_party/llama.cpp
cmake -S third_party/llama.cpp -B third_party/llama.cpp/build \
  -DCMAKE_BUILD_TYPE=Release \
  -DLLAMA_BUILD_TESTS=OFF \
  -DLLAMA_BUILD_TOOLS=OFF \
  -DLLAMA_BUILD_EXAMPLES=OFF \
  -DLLAMA_BUILD_SERVER=OFF \
  -DLLAMA_BUILD_APP=OFF \
  -DLLAMA_BUILD_UI=OFF \
  -DLLAMA_BUILD_COMMON=OFF \
  -DLLAMA_CURL=OFF \
  -DGGML_CUDA=OFF \
  -DGGML_METAL=OFF \
  -DGGML_VULKAN=OFF \
  -DGGML_HIP=OFF \
  -DGGML_CCACHE=OFF
cmake --build third_party/llama.cpp/build --target llama -j
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
CGO_ENABLED=1 go test ./internal/inference/llama/...
```

Integration tests with a real GGUF file (mirrors CI — see [testdata/ci/README.md](testdata/ci/README.md)):

```bash
export TB_TEST_GGUF_PATH=/path/to/smollm2-135m-instruct-q4_k_m.gguf
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
CGO_ENABLED=1 go test -tags integration ./internal/inference/llama/...
CGO_ENABLED=1 go test -tags integration ./internal/runtime/...
```

### CI jobs (merge-blocking on `main`)

| Job | Purpose |
|-----|---------|
| `test` | `CGO_ENABLED=0 go test ./...` — stub inference path |
| `inference-cgo` | CGO unit tests for `./internal/inference/llama/...` (no integration tag) |
| `inference-integration` | Real GGUF — 5 llama adapter integration tests |
| `inference-integration-runtime` | Real GGUF — runtime E2E integration test |

Green CI on `main` requires all four jobs. Integration jobs download a checksum-verified SmolLM2-135M-Instruct Q4_K_M (~105 MB), verify llama.cpp pin `b9553`, and fail on silent skips. Branch protection on `main` requires all four check names.

**CI observability:** per-job timing and cache metrics in Actions step summaries; run history in [planning/metrics/ci-runs.jsonl](planning/metrics/ci-runs.jsonl) with baselines in [planning/metrics/ci-baseline.md](planning/metrics/ci-baseline.md). See [testdata/ci/README.md](testdata/ci/README.md).

### Integrated runtime wiring (009c)

Single shared `runtime.ModelResolver` instance for runtime and provider:

```go
reg := registry.NewModelRegistry()
_ = reg.RegisterModel(registry.ModelDefinition{ID: "tiny", Path: "/models/tiny.gguf"})

resolver := runtime.NewRegistryResolver(reg)
loader := loader.NewStubLoader()
provider := llama.NewLlamaProvider(resolver, llama.DefaultConfig())
bus := events.NewChannelBus(64)
rt := runtime.NewIntegratedModelRuntime(provider, loader, resolver, bus)

_ = rt.LoadModel("tiny")
resp, _ := rt.Generate(runtime.GenerateRequest{ModelID: "tiny", Prompt: "Hello"})
_ = rt.UnloadModel("tiny")
```

Loader-less path (tests, v0.4 compatibility): `runtime.NewModelRuntime(runtime.NewStubProvider(), bus)`.

### Build requirements (CPU backend)

| OS | Toolchain |
|----|-----------|
| Linux | `build-essential`, `cmake`, `CGO_ENABLED=1` |
| Windows | MSVC Build Tools, `cmake`, `CGO_ENABLED=1` |
| macOS | Xcode CLT, `cmake`, `CGO_ENABLED=1` |

Future commands (not yet implemented):

```bash
tinybrain run "your task"
brain-top
```

## Current Version

**V0.6 Inference (in progress)** — [docs/current.md](docs/current.md) | Sprint: [planning/execution/current-sprint.md](planning/execution/current-sprint.md)

## Roadmap

| Version | Goal | Status |
|---------|------|--------|
| v0.1 | Kernel — process model + table | Shipped |
| v0.2 | Registry — agent/model definitions + events | Shipped |
| v0.3 | Hardware profiler | Shipped |
| v0.4 | Runtime shell + stub provider | Shipped |
| v0.5 | Persistent model registry | Shipped |
| v0.6 | llama.cpp inference | In Progress (~85%) |
| v0.7 | MLFQ scheduler | Planned |
| v0.8 | Plugin agents | Planned |
| v1.0 | Integrated runtime + brain-top | Planned |

Details: [docs/specs/](docs/specs/) and [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)

---

*TinyBrain OS is a systems project — a local operating system for AI agents, not a chatbot wrapper.*
