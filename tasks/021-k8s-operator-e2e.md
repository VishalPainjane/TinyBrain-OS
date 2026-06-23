# Task 021 — Operator End-to-End

## Status

Complete

## Goal

Full M9 demo deployment bundle.

## Context

Month 5, Week 20. We have built the CRDs and the Go operator. We need to bundle this up using Kustomize and Kubernetes manifests so users can easily deploy the TinyBrain Operator to their cluster.

## Requirements

- Provide `deploy/k8s/operator/rbac.yaml` with correct RBAC for the operator.
- Provide `deploy/k8s/operator/deployment.yaml` for the operator manager.
- Provide `deploy/k8s/kustomization.yaml` to bundle the CRDs and the operator.
- Update `README.md` with a "Kubernetes Deployment" section.

## Acceptance Criteria

- [x] `rbac.yaml` and `deployment.yaml` exist in `deploy/k8s/operator/`.
- [x] `kustomization.yaml` exists in `deploy/k8s/`.
- [x] `kubectl kustomize deploy/k8s/` successfully renders without errors.
- [x] `README.md` explicitly documents the new feature.

## Out Of Scope

- Helm charts (Kustomize is sufficient for M9).
- Multi-node clustering logic.

---
**Layer:** task
**Target month:** 5
