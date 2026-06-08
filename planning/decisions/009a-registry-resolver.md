# Decision — ModelResolver Port (009a Registry Decoupling)

**Status:** Accepted  
**Date:** 2026-06-08  
**Gate:** B — Registry coupling review  
**Related:** [inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md), ADR-004

---

## Problem

The original 009a plan allowed `internal/inference/llama` to import `internal/registry` and call `ModelRegistry.GetModel` inside `LoadModel(modelID)`.

| Issue | Impact |
|-------|--------|
| **Persistence coupling** | Inference adapter depends on bbolt, gob encoding, and `ModelRegistry` concrete API |
| **Deployment coupling** | Kubernetes sidecar or remote inference pod may not run local registry |
| **ADR-004 erosion** | Adapter should depend on narrow ports, not kernel subsystems |
| **Test coupling** | Inference unit tests must construct registry or bbolt fixtures |
| **Future remote providers** | Cloud adapter would ignore registry paths but inherit import graph |

`internal/registry` is the correct owner of model **metadata** (INV-004). Inference should consume **resolved load specifications**, not registry types.

---

## Proposed Solution

Introduce a **`ModelResolver` port** owned by the inference layer (documentation target in `internal/inference/` — implemented at 009a coding phase).

### Port shape (documentation only)

```go
// ModelSpec is the minimum metadata required to load a GGUF model.
type ModelSpec struct {
    ID            string
    Path          string   // host-local absolute or volume-relative
    Quantization  string
    MemoryBudget  uint64
}

// ModelResolver resolves modelID to load specification.
// Implemented outside inference; injected at construction.
type ModelResolver interface {
    Resolve(modelID string) (ModelSpec, error)
}
```

### Wiring

```text
Composition root (cmd / test harness)
    │
    ├── RegistryModelResolver  ── wraps ModelRegistry.GetModel
    ├── FileModelResolver    ── static map / models.yaml (tests)
    └── RemoteModelResolver  ── future: K8s ConfigMap / env (sidecar)
            │
            ▼
    NewLlamaProvider(resolver ModelResolver, hw HardwareProfile, cfg LlamaConfig)
            │
            ▼
    LoadModel(modelID) → resolver.Resolve(modelID) → load GGUF at spec.Path
```

### Package import rules (009a)

| Package | May import |
|---------|------------|
| `internal/inference/llama` | `internal/runtime`, `internal/hardware`, `internal/inference` (resolver types) |
| `internal/inference/llama` | **Must not** import `internal/registry`, `internal/process`, `internal/scheduler` |
| `internal/registry` | Unchanged — no inference imports |
| Test adapter (e.g. `provider_test.go` or `internal/inference/registry_resolver.go` with build tag) | May import `registry` to implement `ModelResolver` — **adapter file lives in inference package only if it is test-only or explicitly named adapter** |

**Recommended adapter location:** `internal/inference/registry_resolver.go` implements `ModelResolver` by delegating to `*registry.ModelRegistry`. This file is the **only** inference subpackage file that imports registry — keeps `provider.go` and `load_*.go` registry-free. Alternative: adapter in future `cmd/tinybrain/` or `internal/app/`.

`LlamaProvider` holds `ModelResolver` interface field — not `*registry.ModelRegistry`.

### ADR-004 preservation

- `InferenceProvider` port unchanged — still `LoadModel(modelID string)`.
- Runtime unchanged in 009a — still delegates to provider.
- Sidecar: swap `LlamaProvider` for `RemoteGRPCProvider` with same port; resolver returns endpoint metadata in extended `ModelSpec` (future field).
- Kubernetes: pod mounts GGUF at known paths; `FileModelResolver` or ConfigMap-backed resolver — no bbolt in inference image.

---

## Migration Impact

| Artifact | Change |
|----------|--------|
| 009a plan | Replace `registry` import in `provider.go` with `ModelResolver` injection |
| `internal/runtime` | **None** in 009a |
| `internal/registry` | **None** |
| Tests | Fake resolver with fixed `ModelSpec`; optional integration test uses `RegistryModelResolver` adapter |
| `StubProvider` | Unchanged — no resolver (accepts any modelID) |
| Future runtime wiring | Runtime may accept resolver and pass to provider factory — separate task |

**No breaking change** to public registry or runtime APIs.

---

## Architecture Impact

| Invariant | Effect |
|-----------|--------|
| INV-004 | Registry still owns definitions; resolver is a read adapter |
| INV-008 | llama.cpp stays in inference; registry adapter is pure Go |
| Cross-platform | Resolver returns `filepath`-clean paths; OS-agnostic provider |
| inference-lifecycle | Runtime does not resolve paths in 009a; resolver at composition boundary |
| Kubernetes | Enables sidecar + volume-only inference without registry DB |

---

## Approval

- [x] `ModelResolver` port approved for 009a
- [x] `RegistryModelResolver` adapter approved as sole registry touchpoint in inference tree
- [x] `provider.go` / `load_*.go` must not import `internal/registry`

---
**Layer:** planning
