<div align="center">

# TinyBrain OS

**A hardware-aware AI runtime kernel for dynamically orchestrating<br>arbitrary AI agents under resource constraints on local hardware.**

[![Go Report Card](https://goreportcard.com/badge/github.com/VishalPainjane/TinyBrain-OS)](https://goreportcard.com/report/github.com/VishalPainjane/TinyBrain-OS)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Version](https://img.shields.io/badge/release-v1.1-blue.svg)](https://github.com/VishalPainjane/TinyBrain-OS/releases/tag/v1.1)
[![CI Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/VishalPainjane/TinyBrain-OS/actions)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

---

[Architecture](#architecture) | [Quick Start](#quick-start) | [Build from Source](#building-from-source) | [Deploy](#kubernetes-deployment) | [Roadmap](#development-roadmap) | [Contributing](#contributing)

</div>

---

## Why TinyBrain OS

Local AI on consumer hardware breaks down the moment you try to run more than one model. A single 7B model can consume an entire GPU's VRAM budget, leaving nothing for a second agent, a code-analysis tool, or a retrieval pipeline. The standard answer -- "buy more hardware" or "use the cloud" -- defeats the purpose of running locally.

TinyBrain OS solves this by treating AI agents the way an operating system treats processes: they are **scheduled**, **swapped**, and **coordinated** under strict VRAM and RAM budgets. Models are loaded when needed, hibernated when idle, and preempted when higher-priority work arrives. The system adapts to whatever hardware it finds -- from a laptop with 4 GB of RAM to a workstation with dual GPUs.

> **This is a systems project -- a local operating system for AI agents, not a chatbot wrapper.**

---

## Key Capabilities

<table>
<tr>
<td width="50%">

### Hardware-Aware Scheduling

The runtime detects available CPU, RAM, and VRAM at startup and continuously adapts. Model selection, GPU layer offloading, and scheduling aggressiveness are all derived from hardware profiles -- never hardcoded.

</td>
<td width="50%">

### Dynamic Model Swapping

Models and agents are loaded, unloaded, and swapped in and out of VRAM under strict memory budgets. A multi-level feedback queue (MLFQ) scheduler prioritizes active work and hibernates idle contexts.

</td>
</tr>
<tr>
<td width="50%">

### Process-Based Agent Model

Every agent runs as an independent process with its own lifecycle, priority, KV cache, and resource budget. Agents are plugins defined in a registry -- the core never hardcodes agent types like "Planner" or "Coder".

</td>
<td width="50%">

### Strict Architectural Boundaries

The scheduler never imports the inference engine. The runtime never imports the scheduler. Agents never call models directly. These boundaries are enforced by import rules, tested in CI, and documented as system invariants.

</td>
</tr>
<tr>
<td width="50%">

### KV Cache as First-Class Primitive

Attention state is saved, compressed (Zstandard), moved across memory tiers, and restored -- so that resuming a conversation does not require recomputing the entire prompt from scratch.

</td>
<td width="50%">

### OpenAI-Compatible API

A built-in HTTP server exposes `/v1/chat/completions` with Server-Sent Events streaming, making TinyBrain a drop-in replacement for any tool that speaks the OpenAI protocol.

</td>
</tr>
</table>

---

## Architecture

TinyBrain OS follows a layered, hexagonal architecture where each subsystem communicates through well-defined contracts. The core insight is borrowed from operating system design: **treat inference workloads as schedulable processes, not monolithic services.**

```
                            +---------------------+
                            |     API / Router     |  <-- HTTP, CLI
                            +----------+----------+
                                       |
                            +----------v----------+
                            |      Scheduler       |  MLFQ queues, preemption, aging
                            |  (never sees models) |
                            +----------+----------+
                                       |
                            +----------v----------+
                            |       Runtime        |  load, unload, warm, generate
                            |  (never sees sched)  |
                            +----------+----------+
                                       |
                     +-----------------+-----------------+
                     |                                   |
          +----------v----------+             +----------v----------+
          | InferenceProvider   |             |     KV / Memory     |
          | llama.cpp | cloud   |             | cache, swap, tiers  |
          +---------------------+             +---------------------+
                     |
          +----------v----------+
          |   Models (GGUF)     |
          |   Q4_K_M, Q8, F16  |
          +---------------------+
```

**Boundary Matrix** -- who may call whom:

| Component | May Call | Must NOT Call |
|---|---|---|
| API / Router | Scheduler, Registry | Runtime, InferenceProvider |
| Scheduler | Runtime, Process table, Event bus | InferenceProvider, llama.cpp |
| Runtime | InferenceProvider, Loader, KV Manager | Scheduler, UI |
| Registry | _(passive data store)_ | Runtime, Scheduler |
| Agents (plugins) | Runtime (via contract) | Models, Providers, Scheduler |
| UI / brain-top | Telemetry, read-only APIs | Scheduler, Runtime control |

Full architecture documentation: [docs/architecture/overview.md](docs/architecture/overview.md)
System invariants: [docs/architecture/invariants.md](docs/architecture/invariants.md)
Architecture Decision Records: [docs/adr/](docs/adr/)

---

## Quick Start

> **Prerequisites:** [Go 1.22+](https://go.dev/dl/) installed and available on your `PATH`.

### Build

```bash
# Clone the repository
git clone https://github.com/VishalPainjane/TinyBrain-OS.git
cd TinyBrain-OS

# Build the core daemon
go build -o tinybrain ./cmd/tinybrain

# Build the system monitor
go build -o brain-top ./cmd/brain-top
```

### Run

```bash
# Verify system health and detect hardware capabilities
./tinybrain doctor
./tinybrain probe --json

# Manage the model registry
./tinybrain models list
./tinybrain models pull --id <model-name> --url <model-url>

# List available agents
./tinybrain agents list

# Monitor system resources and agent processes in real time
./brain-top snapshot
```

### Configuration

All configuration is driven by environment variables. No config files are required for basic operation.

| Variable | Purpose | Default |
|---|---|---|
| `TB_MODELS_DB` | Path to the model registry database | `~/.tinybrain/models.db` |
| `TB_MODELS_SEED` | YAML seed file for initial model registry population | _(none)_ |
| `TB_NGLAYERS` | Number of GPU layers to offload (`-1` = all) | `0` |
| `TB_LLAMA_LIB_DIR` | Directory containing the llama.cpp shared library | _(auto-detect)_ |

---

## Building from Source

Inference requires CGO and a compiled llama.cpp shared library. The pure-Go path (scheduler, registry, events, process table) works without CGO.

### Toolchain Requirements

| Platform | Required Toolchain |
|---|---|
| Linux | `build-essential`, `cmake` 3.14+, `CGO_ENABLED=1` |
| Windows | MSVC Build Tools, `cmake` 3.14+, `CGO_ENABLED=1` |
| macOS | Xcode Command Line Tools, `cmake` 3.14+, `CGO_ENABLED=1` |

### Step-by-Step: llama.cpp Integration

```bash
# 1. Initialize the llama.cpp submodule
git submodule update --init --recursive third_party/llama.cpp

# 2. Configure CMake (CPU-only example)
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

# 3. Build the shared library
cmake --build third_party/llama.cpp/build --target llama -j

# 4. Run inference tests
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
CGO_ENABLED=1 go test ./internal/inference/llama/...
```

<details>
<summary><strong>CUDA GPU acceleration</strong></summary>

Replace the CMake configuration step with:

```bash
cmake -S third_party/llama.cpp -B third_party/llama.cpp/build \
  -DCMAKE_BUILD_TYPE=Release \
  -DGGML_CUDA=ON \
  -DLLAMA_BUILD_TESTS=OFF \
  -DLLAMA_BUILD_TOOLS=OFF \
  -DLLAMA_BUILD_EXAMPLES=OFF \
  -DLLAMA_BUILD_SERVER=OFF
```

Set `TB_NGLAYERS=-1` to offload all layers to GPU. On Windows, the runtime uses dynamic DLL loading (`syscall.LoadDLL`) per [ADR-006](docs/adr/ADR-006-Windows-GPU-Dynamic-DLL-Backend.md).

</details>

<details>
<summary><strong>Running without CGO</strong></summary>

The core subsystems -- scheduler, registry, event bus, process table, and hardware prober -- are pure Go and require no CGO:

```bash
CGO_ENABLED=0 go test ./...
```

This is the default CI verification tier and is sufficient for all non-inference development.

</details>

---

## Kubernetes Deployment

TinyBrain OS includes a Kubernetes operator for declarative agent and memory management on clusters.

```bash
# Deploy with raw manifests
kubectl apply -k deploy/k8s/

# Or deploy with the Helm chart
helm install tinybrain-operator deploy/helm/tinybrain-operator/
```

> **Note:** The Kubernetes operator is a parallel deployment track. The local daemon does not require Kubernetes and runs directly on the host machine.

---

## Project Structure

```
TinyBrain-OS/
  cmd/
    tinybrain/         CLI daemon and user-facing commands
    brain-top/         Real-time system monitor (TUI)
    operator/          Kubernetes operator entry point
    benchmark/         Performance benchmark harness
  internal/
    agents/            Agent plugin contracts and fleet management
    api/               HTTP server, OpenAI-compatible endpoints
    events/            Typed event bus for inter-component messaging
    hardware/          Hardware detection, GPU enumeration, profiling
    inference/         InferenceProvider adapters (llama.cpp, cloud)
    k8s/               Kubernetes CRD controllers and reconcilers
    kv/                KV cache persistence and compression
    loader/            GGUF model loader and memory-mapped I/O
    memory/            Memory tier management (VRAM, RAM, NVMe)
    process/           OS-style process table and lifecycle states
    registry/          Agent, model, and tool registry (bbolt-backed)
    router/            Task classification and agent routing
    runtime/           Model lifecycle orchestration
    scheduler/         MLFQ scheduler, preemption, priority aging
    swap/              VRAM-to-RAM context swap manager
    telemetry/         Metrics collection and export
    tools/             Tool execution layer for agent capabilities
  deploy/
    k8s/               Kubernetes manifests (kustomize)
    helm/              Helm chart for the TinyBrain operator
  docs/
    architecture/      Subsystem design documents
    adr/               Architecture Decision Records (ADR-001 through ADR-008)
    contracts/         Interface contracts and ownership boundaries
    specs/             Version specifications
    constitution.md    Architectural law -- the project's non-negotiable rules
    glossary.md        Canonical term definitions
  tests/               Cross-package integration and boundary tests
  benchmarks/          Performance benchmark scripts and report templates
  third_party/
    llama.cpp/         Inference backend (Git submodule)
```

---

## Development Roadmap

TinyBrain OS is built in monthly development sprints. Each phase ships a tagged release with a clear set of subsystems.

| Phase | Milestone | Release | Subsystems |
|:---:|---|:---:|---|
| Month 1 | Single-Node Foundation | `v0.3` | Kernel, Process Table, Event Bus, Hardware Prober |
| Month 2 | Inference Engine | `v0.6` | Runtime interface, llama.cpp integration, GGUF loader |
| Month 3 | OS Memory Model | `v0.7` | KV Cache Manager, CPU Swap Manager, MLFQ Scheduler |
| Month 4 | Agents as Processes | `v0.8` | Agent Plugin API, Fleet Registry, Event Pipeline |
| Month 5 | Control Plane | `v0.9` | Kubernetes CRDs, Fleet Operator, Network Bridge |
| Month 6 | v1.0 Release | `v1.0` | System Integration, brain-top, Benchmark Suite |
| **Month 7** | **Advanced Subsystems** | **`v1.1`** | **RadixAttention, KV Compression, Helm, OpenAI API** |

Detailed sprint plans: [planning/execution/current-sprint.md](planning/execution/current-sprint.md)
Master roadmap: [planning/roadmap/master-roadmap.md](planning/roadmap/master-roadmap.md)
Release history: [CHANGELOG.md](CHANGELOG.md)

---

## Design Principles

TinyBrain OS is governed by a [constitution](docs/constitution.md) -- a set of non-negotiable architectural rules that every change must satisfy.

| Principle | What It Means |
|---|---|
| **Hardware determines model selection** | The runtime adapts to detected capabilities. Fixed model sizes per agent are forbidden. |
| **Agents are plugins** | Agent capabilities are registry entries, not hardcoded Go types. |
| **Scheduler never sees inference** | The scheduler delegates all model operations to the runtime. It never imports llama.cpp. |
| **Runtime never sees the scheduler** | Bidirectional decoupling. Each subsystem evolves independently. |
| **Local-first by default** | Core operation requires no external API. Cloud providers are optional adapters. |
| **Structured JSON IPC only** | Components communicate via typed JSON messages. Natural-language IPC is forbidden. |

---

## Testing

TinyBrain OS uses a tiered testing strategy. Every commit must pass the fast tier; inference changes require the integration tier.

```bash
# Fast tier -- pure Go, no CGO, no GPU required
CGO_ENABLED=0 go test ./...

# Integration tier -- requires llama.cpp and a GGUF model
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
export TB_TEST_GGUF_PATH=/path/to/model.gguf
CGO_ENABLED=1 go test -tags integration -count=1 ./internal/inference/llama/...
CGO_ENABLED=1 go test -tags integration -count=1 ./internal/runtime/...

# CLI tier -- command output golden tests and fuzzing
go test ./cmd/tinybrain/...
go test -fuzz=Fuzz -fuzztime=30s ./cmd/tinybrain/...
```

Full policy: [docs/testing-policy.md](docs/testing-policy.md)

---

## Contributing

Contributions are welcome. Before submitting a pull request, please review the project guidelines:

| Document | Purpose |
|---|---|
| [CONTRIBUTING.md](CONTRIBUTING.md) | Development workflow, commit conventions, definition of done |
| [docs/testing-policy.md](docs/testing-policy.md) | Test tiers, coverage expectations, regression test naming |
| [docs/constitution.md](docs/constitution.md) | Architectural rules that every change must satisfy |
| [SECURITY.md](SECURITY.md) | Vulnerability reporting and secrets policy |
| [CHANGELOG.md](CHANGELOG.md) | Release history and notable changes |

---

## License

TinyBrain OS is released under the [MIT License](https://opensource.org/licenses/MIT).

---

<div align="center">
<br>

**[Documentation](docs/)** | **[Architecture Decisions](docs/adr/)** | **[Changelog](CHANGELOG.md)** | **[Security](SECURITY.md)**

<br>
<sub>Built with discipline by the TinyBrain OS contributors.</sub>

</div>
