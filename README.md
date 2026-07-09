<div align="center">
  <h1>TinyBrain OS</h1>
  <p><em>A hardware-aware AI runtime kernel for dynamically orchestrating specialized agents on constrained local hardware.</em></p>
  
  [![Go Report Card](https://goreportcard.com/badge/github.com/VishalPainjane/TinyBrain-OS)](https://goreportcard.com/report/github.com/VishalPainjane/TinyBrain-OS)
  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
  [![Version](https://img.shields.io/badge/version-v1.0.0--dev-blue.svg)]()
  [![CI Status](https://img.shields.io/badge/build-passing-brightgreen)]()
</div>

<hr>

Instead of one monolithic model handling every task, **TinyBrain OS** treats AI agents as scheduled processes — loaded when needed, swapped when idle, and coordinated under strict VRAM and RAM budgets. Local AI fails on consumer hardware when treated as a single-model problem. TinyBrain exists to make it practical through dynamic capability allocation, not static model-centric workflows.

> [!IMPORTANT]
> *TinyBrain OS is a systems project — a local operating system for AI agents, not a chatbot wrapper.*

---

## Table of Contents

- [Key Features](#key-features)
- [Architecture Overview](#architecture-overview)
- [Quick Start](#quick-start)
  - [Installation](#installation)
  - [Basic Commands](#basic-commands)
- [Building from Source](#building-from-source)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Roadmap](#roadmap)
- [Contributing](#contributing)

---

## Key Features

- **Local-first & Hardware-aware:** Intelligently adapts to available CPU, RAM, and VRAM.
- **Dynamic Model Swapping:** Seamlessly swaps models and agents in/out of memory under strict budgets.
- **Process-based Agents:** Agents are treated as independent processes (plugins), loaded only when required.
- **Decoupled Architecture:** MLFQ Scheduler operates independently from the inference engine.
- **Structured JSON IPC:** Fast, robust communication between internal OS components.

---

## Architecture Overview

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

See the full architectural breakdown in [docs/architecture/overview.md](docs/architecture/overview.md).

---

## Quick Start

**Current state:** `v1.0` (Development) | **Current Sprint:** Month 7 (Advanced Subsystems). 
Month 4 memory and workflow orchestration foundations (KV, swap, brain-top, agents) are on `main`.

### Installation

> [!NOTE]
> For now, the project is built from source. Ensure you have Go 1.22+ installed.

```bash
# Build the core TinyBrain daemon
go build -o tinybrain ./cmd/tinybrain

# Build the system monitor (similar to 'top')
go build -o brain-top ./cmd/brain-top
```

### Basic Commands

Once built, you can run diagnostic checks and manage models/agents:

```bash
# Check system health and hardware capabilities
./tinybrain doctor
./tinybrain probe --json

# Manage models in the registry
./tinybrain models pull --id new-model --url <URL>
./tinybrain models list

# List available agents
./tinybrain agents list

# Monitor system resources and agent processes
./brain-top snapshot
```

#### Environment Variables

| Variable | Purpose |
|----------|---------|
| `TB_MODELS_DB` | Registry database path (default `~/.tinybrain/models.db`) |
| `TB_MODELS_SEED` | YAML seed when DB is empty (e.g. `testdata/models.yaml`) |
| `TB_NGLAYERS` | GPU layers for CUDA builds (`-1` = all) |
| `TB_LLAMA_LIB_DIR` | llama.cpp shared library directory |

---

## Building from Source

> [!CAUTION]
> Inference requires CGO and a built `llama.cpp` shared library. Be sure to configure CMake correctly for your specific GPU architecture to avoid fallback to CPU inference.

### Build Requirements

| OS | Toolchain |
|----|-----------|
| Linux | `build-essential`, `cmake`, `CGO_ENABLED=1` |
| Windows | MSVC Build Tools, `cmake`, `CGO_ENABLED=1` |
| macOS | Xcode CLT, `cmake`, `CGO_ENABLED=1` |

### Compiling with llama.cpp

```bash
# 1. Initialize submodule
git submodule update --init --recursive third_party/llama.cpp

# 2. Configure cmake (CPU backend example)
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

# 3. Build library
cmake --build third_party/llama.cpp/build --target llama -j

# 4. Run tests with library in path
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
CGO_ENABLED=1 go test ./internal/inference/llama/...
```

---

## Kubernetes Deployment

You can deploy the TinyBrain OS operator to a local Kubernetes cluster (like `kind` or `minikube`) to manage agents and memory thresholds declaratively.

```bash
kubectl apply -k deploy/k8s/
```
> [!WARNING]
> The Kubernetes operator is a parallel track. The local v1.0 demo does not require Kubernetes and runs directly on your host machine.

---

## Roadmap

| Phase | Goal | Status | Subsystems |
|---|---|---|---|
| [Month 1](planning/roadmap/months/month-01.md) | Single-Node Foundation | Done: Shipped `v0.3` | Kernel, Process Table, Event Bus, Hardware Prober |
| [Month 2](planning/roadmap/months/month-02.md) | The Inference Engine | Done: Shipped `v0.6` | Runtime interface, `llama.cpp` integration, GGUF loader |
| [Month 3](planning/roadmap/months/month-03.md) | OS Memory Model | Done: Shipped `v0.7` | KV Cache Manager, CPU Swap Manager, MLFQ Scheduler |
| [Month 4](planning/roadmap/months/month-04.md) | Agents as Processes | Done: Shipped `v0.8` | Agent Plugin API, Fleet Registry, Event Pipeline |
| [Month 5](planning/roadmap/months/month-05.md) | Control Plane | Done: Shipped `v0.9` | Kubernetes CRDs, Fleet Operator, Network Bridge |
| [Month 6](planning/roadmap/months/month-06.md) | V1.0 Release | Done: Shipped `v1.0` | System Integration, `brain-top`, Benchmark Suite |

See the master roadmap: [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)

---

## Contributing

We welcome contributions! Please review our guidelines before submitting PRs:

- **Solo workflow:** [CONTRIBUTING.md](CONTRIBUTING.md)
- **Testing policy:** [docs/testing-policy.md](docs/testing-policy.md)
- **Repo health:** [planning/metrics/repo-health.md](planning/metrics/repo-health.md)
- **Changelog:** [CHANGELOG.md](CHANGELOG.md)

---
<div align="center">
  <sub>Built with love by the TinyBrain OS Contributors</sub>
</div>
