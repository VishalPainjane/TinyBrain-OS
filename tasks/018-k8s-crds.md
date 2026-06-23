# Task 018 — Kubernetes CRDs

## Status

Complete

## Goal

Design the YAML schemas and validation rules for the four core Kubernetes Custom Resource Definitions (CRDs) required for the TinyBrain OS operator.

## Context

Month 5, Week 17. The platform is transitioning to support a Kubernetes operator model (M9 Milestone). Before writing controller reconciliation logic, we must strictly define the data contract through CRDs.

## Requirements

- Create `deploy/k8s/crds/agent.yaml`
- Create `deploy/k8s/crds/task.yaml`
- Create `deploy/k8s/crds/kvcache.yaml`
- Create `deploy/k8s/crds/swappolicy.yaml`
- Validate that the schemas map accurately to our existing Go types where appropriate, or abstract them safely for K8s users.
- Provide sample manifests in `deploy/k8s/samples/` for each CRD to ensure they apply correctly using `kubectl apply --dry-run=client`.

## Files

- `deploy/k8s/crds/*.yaml`
- `deploy/k8s/samples/*.yaml`

## Acceptance Criteria

- [x] Four CRD files exist with OpenAPI v3 validation schemas.
- [x] Four sample files exist.
- [x] `kubectl apply --dry-run=client -f deploy/k8s/samples/` passes without validation errors.

## Out Of Scope

- Controller Go logic (`Reconcile` functions).
- Kubebuilder scaffolding (we write raw YAML or minimal config for now to keep it lightweight).

## Related

- Spec: [RFC-003-Kubernetes-Operator.md](../docs/rfc/RFC-003-Kubernetes-Operator.md)
- Month plan: [month-05.md](../planning/roadmap/months/month-05.md)

---
**Layer:** task
**Target month:** 5
