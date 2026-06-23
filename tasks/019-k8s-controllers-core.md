# Task 019 — Agent and Task Controllers

## Status

In Progress

## Goal

Reconcile loop skeleton for `Agent` and `Task` resources using `controller-runtime`.

## Context

Month 5, Week 18. Building off the CRDs defined in Task 018, we must now build the actual Kubernetes operator logic. We are intentionally avoiding heavy Kubebuilder scaffolding to keep the codebase lean and strictly adhere to our core `internal/runtime` separation.

## Requirements

- Add `sigs.k8s.io/controller-runtime` as a dependency.
- Define Go API structs in `internal/k8s/api/v1alpha1/` for Agent and Task.
- Implement `AgentReconciler` in `internal/k8s/controllers/agent_controller.go`. It should call an interface adapter that maps to `runtime.LoadModel`.
- Implement `TaskReconciler` in `internal/k8s/controllers/task_controller.go`.
- Scaffold the `cmd/operator/main.go` entrypoint to register the controllers and start the manager.

## Acceptance Criteria

- [ ] Go API definitions for `Agent` and `Task` exist.
- [ ] `AgentReconciler` and `TaskReconciler` implement the `reconcile.Reconciler` interface.
- [ ] `cmd/operator` compiles successfully.
- [ ] Core `internal/runtime` and `internal/scheduler` do not import Kubernetes packages.

## Out Of Scope

- KVCache and SwapPolicy controllers (Week 19).
- Deploying to a real cluster (can be mocked or unit tested for now).

---
**Layer:** task
**Target month:** 5
