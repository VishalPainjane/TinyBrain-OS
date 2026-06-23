# Task 020 — Memory Controllers

## Status

Complete

## Goal

Declarative swap thresholds and KV tracking via Kubernetes Custom Resources.

## Context

Month 5, Week 19. Following the core Agent/Task controllers, we need to map the memory limits directly into Kubernetes. The `memory.tinybrain.io` API group encapsulates `KVCache` and `SwapPolicy` CRDs.

## Requirements

- Define Go API structs in `internal/k8s/api/memory/v1alpha1/` for `KVCache` and `SwapPolicy`.
- Register the `memory.tinybrain.io` GroupVersion.
- Implement `KVCacheReconciler` in `internal/k8s/controllers/kvcache_controller.go`.
- Implement `SwapPolicyReconciler` in `internal/k8s/controllers/swappolicy_controller.go`.
- Update the `cmd/operator/main.go` entrypoint to register the new scheme and controllers alongside the core controllers.

## Acceptance Criteria

- [x] Go API definitions for `KVCache` and `SwapPolicy` exist.
- [x] Both reconcilers implement the `reconcile.Reconciler` interface.
- [x] `cmd/operator` compiles successfully with the new memory group registered.
- [x] Core `internal/runtime` does not import Kubernetes packages.

## Out Of Scope

- Modifying the underlying local swap eviction logic in `internal/swap/` (this is purely the K8s interface wrapper).

---
**Layer:** task
**Target month:** 5
