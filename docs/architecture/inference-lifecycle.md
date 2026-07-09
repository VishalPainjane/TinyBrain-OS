# Inference Lifecycle

Canonical model lifecycle across runtime, loader, and inference adapter.  
**Purpose:** Single source of truth for state names, ownership, transitions, and events — prevents duplicate state machines and platform-specific drift.

**Related:** [runtime.md](runtime.md), [inference-backend-matrix.md](inference-backend-matrix.md), [cross-platform.md](cross-platform.md), ADR-004

---

## 1. Canonical Model Lifecycle

Six states apply to **resident model weight** residency. These names are fixed across all packages. Note that in TinyBrain OS v2, the model weights are decoupled from the sequence execution context (KV Cache).

| State | Meaning |
|-------|---------|
| `NOT_LOADED` | Model not resident; no weights mapped; no llama context |
| `LOADING` | Load in progress (mmap, GPU init, context create) |
| `ACTIVE` | Ready for `Generate`; weights and inference context valid |
| `WARM` | Resident but idle; retained for fast reuse; not serving inference |
| `UNLOADING` | Teardown in progress (context free, unmap) |
| `UNLOADED` | Fully evicted; may transition to `NOT_LOADED` or `LOADING` on next load |

**Sequence Lifecycle:** process states (`NEW`, `READY`, `RUNNING`, `PREEMPTED`, `HIBERNATED`, …) in `internal/process/` apply to the individual request sequences (KV Caches), not the model weights. This is an independent subsystem managed by the sequence scheduler.

```text
NOT_LOADED ──Load──► LOADING ──success──► ACTIVE ◄──Warm──┐
     ▲                    │                  │            │
     │                    fail               Generate       │
     │                    ▼                  │            │
     └────────────── UNLOADED ◄── UNLOADING ◄┘            │
                            ▲              ▲              │
                            │              └── Evict ────┤
                            └──────── Prefetch/Warm path ─┘
                                      (WARM state)
```

---

## 2. State Ownership

Only one package **authoritatively sets** each transition. Others may read or mirror for tests until integration.

| Transition | Owner (target) | Owner (009a actual) | Notes |
|------------|----------------|---------------------|-------|
| → `LOADING` | Loader or inference adapter | **Inference adapter** | 009a: loader unwired |
| → `ACTIVE` | Inference adapter | **Inference adapter** | After successful llama load |
| → `WARM` | Loader | *Deferred* | `StubLoader.Warm` exists; not wired |
| → `UNLOADING` | Loader or inference adapter | **Inference adapter** | 009a unload path |
| → `UNLOADED` | Loader or inference adapter | **Inference adapter** | |
| → `NOT_LOADED` | Loader (registry of known models) | *Implicit* | Default for unknown model ID |

**Rule:** Inference adapter owns **engine context** lifetime. Loader owns **capacity / eviction policy** when wired. Runtime owns **orchestration and events** — it does not store lifecycle state in 009a.

**009a debt (documented):** `StubLoader` and `StubProvider` track parallel state ([v0.4 postmortem](../planning/postmortems/v0.4.md)). Integration task will make loader the policy owner and inference the execution owner.

---

## 3. Runtime Ownership

| Responsibility | Runtime (`internal/runtime/`) |
|----------------|-------------------------------|
| Public API | `LoadModel`, `UnloadModel`, `Generate`, `SaveContext`, `RestoreContext` |
| Delegation | Calls `InferenceProvider` methods |
| Lifecycle state storage | **No** in v0.4 shell — does not persist `ModelState` |
| Model ID → path resolution | **Target:** composition root or future runtime wiring via resolver; **not** in 009a runtime changes |
| Loader orchestration | **Deferred** — warm, prefetch, evict |
| Event emission | **Yes** — `ModelLoaded`, `ModelUnloaded` on successful provider load/unload |

Runtime is the **orchestration façade** agents and (future) scheduler use. It must remain free of llama.cpp, registry persistence details, and OS-specific logic.

---

## 4. Loader Ownership

| Responsibility | Loader (`internal/loader/`) |
|----------------|----------------------------|
| API | `Load`, `Unload`, `Warm`, `Prefetch`, `Evict`, `State` |
| State machine | Tracks `ModelState` per `modelID` |
| Path input | Receives `(modelID, path)` — does not resolve registry |
| mmap / engine | **No** — policy and state only in stub; real mmap in inference adapter (009a) |
| Eviction | `EvictionPolicy` shell (LRU) when over capacity |

**009a:** Loader unchanged. Inference adapter performs real mmap; loader integration is a **future integration task**.

---

## 5. Inference Ownership

| Responsibility | Inference adapter (`internal/inference/llama/`) |
|----------------|--------------------------------------------------|
| Port | Implements `runtime.InferenceProvider` |
| Engine state | llama model handle, context, backend device init |
| Load / unload | `LOADING` → `ACTIVE` / `UNLOADING` → `UNLOADED` internally |
| Generate | Transitions only while `ACTIVE` (009b) |
| SaveContext / RestoreContext | Stub until task 011 KV manager |
| Model metadata | Via **`ModelResolver` port** — not `internal/registry` import ([009a-registry-resolver.md](../../planning/decisions/009a-registry-resolver.md)) |
| Backend selection | From `hardware.Backend` + build tags — see [build tag matrix](../../planning/decisions/009a-build-tags.md) |

Inference **must not** import `internal/process`, `internal/scheduler`, or `internal/registry`.

---

## 6. Event Emission Points

| Event | Emitter | When | Payload |
|-------|---------|------|---------|
| `ModelLoaded` | **Runtime** | After `InferenceProvider.LoadModel` succeeds | `ModelID` |
| `ModelUnloaded` | **Runtime** | After `InferenceProvider.UnloadModel` succeeds | `ModelID` |
| `SwapStarted` | Runtime (future) | Model swap begins | `FromModelID`, `ToModelID` |
| `SwapCompleted` | Runtime (future) | Swap ends | TBD |
| `KVStored` / `KVLoaded` | KV manager (task 011) | Context tier moves | TBD |

**009a:** Only `ModelLoaded` / `ModelUnloaded` apply when runtime is wired to `LlamaProvider` in tests or demo. Inference adapter does **not** publish bus events directly.

---

## 7. Future Integration Boundaries

```text
Scheduler (future)
    │ load/unload commands (no inference imports)
    ▼
Runtime
    ├──► Loader (warm / prefetch / evict / State)
    ├──► ModelResolver ◄── Registry adapter (composition root)
    └──► InferenceProvider (llama / remote / sidecar)
              └──► llama.cpp / RPC (INV-008)
```

| Boundary | Rule |
|----------|------|
| Scheduler → runtime | Commands only; no `modelID` → path resolution in scheduler |
| Runtime → registry | Read via resolver at composition root; runtime package unchanged in 009a |
| Runtime → loader | Runtime calls loader before/after provider on integrated path (future) |
| Runtime → inference | `InferenceProvider` port only |
| Inference → registry | **Forbidden** — use `ModelResolver` |
| Loader → inference | **No direct import** — coordinated by runtime |

**Integration task (post-009b):** On `LoadModel(modelID)`, runtime resolves spec → loader `Load(modelID, path)` for state → provider `LoadModel(modelID)` for engine — or merged orchestration documented in task spec.

---

## 8. Kubernetes Visibility Expectations

| Concern | Expectation |
|---------|-------------|
| Control plane | Pure Go kernel (`process`, `scheduler`, `runtime` shell, `registry`, `events`) — no CGO in minimal image |
| Inference workload | Separate image or init container with llama.cpp + GGUF volume mount |
| Model paths | `ModelResolver` returns paths inside pod volume (`/models/...`) — not host registry DB |
| Lifecycle probes | **Future:** readiness when model `ACTIVE`; liveness independent of load state |
| State visibility | `ModelLoaded` / `ModelUnloaded` events for observability; optional metrics sidecar |
| GPU | Node selector + device plugin; CUDA/ROCm/Metal/Vulkan per image variant — one backend per image |
| Sidecar pattern | Remote `InferenceProvider` adapter; resolver returns RPC endpoint; scheduler/runtime unchanged (ADR-004) |
| Config | `models.yaml` seed or ConfigMap → registry or flat file resolver — not hot-reload in v0.5 |

**009a:** Document expectations only; no K8s manifests.

---

## 9. Package Dependency Summary

| Package | Lifecycle role | Imports inference? | Imports registry? |
|---------|----------------|------------------|-------------------|
| `process` | Process states only | No | No |
| `scheduler` | Future commands | No | No (read profiles later) |
| `runtime` | Orchestration, events | No (port interface) | Deferred |
| `registry` | Metadata store | No | — |
| `loader` | Policy / state machine | No | No |
| `inference/llama` | Engine execution | — | **No** (resolver port) |
| Composition root / tests | Wire resolver | Yes | Yes (adapter only) |

---

## 10. Non-Goals

- KV cache tier states (task 011)
- Process lifecycle
- Agent plugin lifecycle
- Multi-model concurrent `ACTIVE` (v0.6 out of scope)

---
**Layer:** architecture  
**Last updated:** 2026-06-08
