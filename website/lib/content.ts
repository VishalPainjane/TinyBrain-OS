export const GITHUB_URL =
  "https://github.com/VishalPainjane/TinyBrain-OS";

export const REPO_DOCS_BASE =
  "https://github.com/VishalPainjane/TinyBrain-OS/blob/main";

export const features = [
  {
    title: "Hardware-Aware Scheduling",
    description:
      "Probes RAM, VRAM, and GPU at boot. Model fleets, offload depth, and scheduler aggressiveness adapt to what the machine actually has — never hardcoded.",
    accent: "red",
  },
  {
    title: "Token-Boundary Preemption",
    description:
      "MLFQ scheduling at the token boundary — not the request. Continuous batching with in-flight multiplexing, inspired by vLLM but built in Go and CUDA.",
    accent: "orange",
  },
  {
    title: "Paged KV Virtual Memory",
    description:
      "Attention state lives in fixed VRAM blocks with a page table. Evict KV blocks to RAM or NVMe without reloading multi-gigabyte model weights.",
    accent: "yellow",
  },
  {
    title: "Context Hibernation",
    description:
      "Preempted agents compress with Zstandard and move down the memory hierarchy. Resume conversations without recomputing the full prompt.",
    accent: "orange",
  },
  {
    title: "Agents as Plugins",
    description:
      "Planner, Coder, Shell — all registry entries. The kernel only knows process, priority, state, and resources. Linux doesn't know Chrome; TinyBrain doesn't know agents.",
    accent: "red",
  },
  {
    title: "OpenAI-Compatible API",
    description:
      "Drop-in /v1/chat/completions with SSE streaming. Any tool that speaks OpenAI can route through TinyBrain on local hardware.",
    accent: "yellow",
  },
];

export const lifecycleSteps = [
  {
    step: "01",
    title: "Task submission",
    detail: "HTTP API or CLI submits work. OpenAI-compatible gateway accepts chat completions.",
  },
  {
    step: "02",
    title: "Router",
    detail: "Classifies intent, selects agent capability from registry, applies ChatML templating.",
  },
  {
    step: "03",
    title: "Scheduler",
    detail: "MLFQ admission, VRAM budget check, priority assignment. Never touches inference code.",
  },
  {
    step: "04",
    title: "Process table",
    detail: "OS-style NEW → READY → RUNNING lifecycle. Each agent is a first-class process with PID.",
  },
  {
    step: "05",
    title: "Runtime",
    detail: "Load, warm, generate. Orchestrates InferenceProvider and loader without importing scheduler.",
  },
  {
    step: "06",
    title: "CUDA inference",
    detail: "PagedAttention, Flash-Softmax, GQA — bare-metal forward pass through llama.cpp or Windows DLL backend.",
  },
  {
    step: "07",
    title: "Token stream",
    detail: "Non-blocking channels → SSE → client. Telemetry feeds brain-top in real time.",
  },
];

export const memoryTiers = [
  {
    tier: "VRAM",
    label: "Active weights + KV blocks",
    target: "< 90% utilization",
    color: "#E63946",
  },
  {
    tier: "RAM",
    label: "Compressed KV + warm pages",
    target: "< 70% utilization",
    color: "#FF6B35",
  },
  {
    tier: "NVMe",
    label: "Cold archive + hibernated sessions",
    target: "Swap partition for AI",
    color: "#FFD60A",
  },
];

export const adrs = [
  { id: "ADR-001", title: "Hardware-Aware Runtime", file: "docs/adr/ADR-001-Hardware-Aware-Runtime.md" },
  { id: "ADR-002", title: "Agent Plugin System", file: "docs/adr/ADR-002-Agent-Plugin-System.md" },
  { id: "ADR-003", title: "Event-Driven Core", file: "docs/adr/ADR-003-Event-Driven-Core.md" },
  { id: "ADR-004", title: "Hexagonal Architecture", file: "docs/adr/ADR-004-Hexagonal-Architecture.md" },
  { id: "ADR-005", title: "Local-First", file: "docs/adr/ADR-005-Local-First.md" },
  { id: "ADR-006", title: "Windows GPU Dynamic DLL", file: "docs/adr/ADR-006-Windows-GPU-Dynamic-DLL-Backend.md" },
  { id: "ADR-007", title: "Daemonized Inference Engine", file: "docs/adr/ADR-007-Daemonized-Inference-Engine.md" },
  { id: "ADR-008", title: "Iteration-Level Scheduling", file: "docs/adr/ADR-008-Iteration-Level-Scheduling-and-Paged-Memory.md" },
];

export const researchNotes = [
  { title: "MLFQ Scheduling", file: "docs/research/mlfq-notes.md", topic: "OS scheduling theory" },
  { title: "vLLM Architecture", file: "docs/research/vllm-notes.md", topic: "Continuous batching" },
  { title: "KV Cache Strategies", file: "docs/research/kv-cache-notes.md", topic: "Paged attention & swap" },
  { title: "llama.cpp Integration", file: "docs/research/llama-cpp-notes.md", topic: "GGUF inference backend" },
  { title: "Hardware Profiles", file: "docs/research/hardware-profiles-notes.md", topic: "Tier classification" },
];

export const externalResearch = [
  { name: "PagedAttention (vLLM)", url: "https://arxiv.org/abs/2309.06180", type: "Paper" },
  { name: "FlashAttention", url: "https://arxiv.org/abs/2205.14135", type: "Paper" },
  { name: "Operating Systems: Three Easy Pieces", url: "https://pages.cs.wisc.edu/~remzi/OSTEP/", type: "Book" },
  { name: "llama.cpp", url: "https://github.com/ggml-org/llama.cpp", type: "Project" },
  { name: "Kubernetes Operator Pattern", url: "https://kubernetes.io/docs/concepts/extend-kubernetes/operator/", type: "Docs" },
];

export const rfcs = [
  { id: "RFC-001", title: "KV Hibernation", file: "docs/rfc/RFC-001-KV-Hibernation.md" },
  { id: "RFC-002", title: "MLFQ Scheduler", file: "docs/rfc/RFC-002-MLFQ-Scheduler.md" },
  { id: "RFC-003", title: "Kubernetes Operator", file: "docs/rfc/RFC-003-Kubernetes-Operator.md" },
  { id: "RFC-004", title: "Registry Facade", file: "docs/rfc/RFC-004-Registry-Facade.md" },
  { id: "RFC-005", title: "Event Ordering", file: "docs/rfc/RFC-005-Event-Ordering.md" },
];

export const docLinks = [
  {
    title: "Quick Start",
    description: "Build tinybrain, run doctor, pull models, monitor with brain-top.",
    href: `${REPO_DOCS_BASE}/README.md#quick-start`,
    category: "Getting started",
  },
  {
    title: "Building from Source",
    description: "CGO, llama.cpp submodule, CUDA acceleration, Windows DLL backend.",
    href: `${REPO_DOCS_BASE}/README.md#building-from-source`,
    category: "Getting started",
  },
  {
    title: "Architecture Overview",
    description: "Hexagonal layout, boundary matrix, subsystem contracts.",
    href: `${REPO_DOCS_BASE}/docs/architecture/overview.md`,
    category: "Architecture",
  },
  {
    title: "System Invariants",
    description: "INV-001 through INV-008 — enforceable architectural laws.",
    href: `${REPO_DOCS_BASE}/docs/architecture/invariants.md`,
    category: "Architecture",
  },
  {
    title: "Constitution",
    description: "Non-negotiable rules every change must satisfy.",
    href: `${REPO_DOCS_BASE}/docs/constitution.md`,
    category: "Governance",
  },
  {
    title: "Testing Policy",
    description: "Fast, integration, and CLI tiers. Boundary import tests.",
    href: `${REPO_DOCS_BASE}/docs/testing-policy.md`,
    category: "Development",
  },
  {
    title: "Contributing",
    description: "Session workflow, verification tiers, definition of done.",
    href: `${REPO_DOCS_BASE}/CONTRIBUTING.md`,
    category: "Development",
  },
  {
    title: "Kubernetes Deployment",
    description: "Operator, CRDs, Helm chart for cloud-native control plane.",
    href: `${REPO_DOCS_BASE}/README.md#kubernetes-deployment`,
    category: "Deploy",
  },
  {
    title: "Project Overview",
    description: "Full technical deep-dive — 23 sections, end-to-end.",
    href: `${REPO_DOCS_BASE}/PROJECT_OVERVIEW.md`,
    category: "Reference",
  },
  {
    title: "Glossary",
    description: "Canonical term definitions used across the project.",
    href: `${REPO_DOCS_BASE}/docs/glossary.md`,
    category: "Reference",
  },
];

export const releases = [
  { version: "v0.3", label: "Foundation", detail: "Kernel, process table, event bus, hardware prober" },
  { version: "v0.6", label: "Inference", detail: "llama.cpp CPU + CUDA, GGUF loader, mmap" },
  { version: "v0.7", label: "Memory + Scheduler", detail: "MLFQ Q0–Q3, swap manager, brain-top prototype" },
  { version: "v0.8", label: "Agents", detail: "Plugin API, fleet registry, event pipeline" },
  { version: "v0.9", label: "Control plane", detail: "K8s CRDs, fleet operator, network bridge" },
  { version: "v1.0", label: "Integration", detail: "System integration, benchmarks, production brain-top" },
  { version: "v1.1", label: "Advanced", detail: "RadixAttention, KV compression, Helm, OpenAI API", current: true },
];

export const inspirationSites = [
  { name: "llama.app", url: "https://www.llama.app/", note: "Hardware grid, install-first hero" },
  { name: "Ollama", url: "https://ollama.com", note: "CLI minimalism, developer clarity" },
  { name: "Razorpay Sprint 26", url: "https://www.awwwards.com/sites/razorpay-sprint-26", note: "Scroll-driven product story" },
  { name: "Orange Giga Game", url: "https://www.awwwards.com/sites/orange-giga-game", note: "Bold orange brand palette" },
];

export const navLinks = [
  { href: "#vision", label: "Vision" },
  { href: "#architecture", label: "Architecture" },
  { href: "#lifecycle", label: "How it works" },
  { href: "#memory", label: "Memory" },
  { href: "#research", label: "Research" },
  { href: "/docs", label: "Docs" },
];
