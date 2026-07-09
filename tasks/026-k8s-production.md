# 026 Kubernetes Production Deploy

**Status:** In Progress
**Owner:** TBD
**Tags:** k8s, production, devops
**Layer:** deploy

## Context

The TinyBrain operator currently relies on plain Kustomize manifests in `deploy/k8s`. For a proper v1.1 production rollout, the cluster administrator needs the ability to customize replicas, resource limits, image tags, RBAC rules, and observability ports (metrics) seamlessly. A Helm Chart provides the industry-standard way to package and distribute these configurations for Kubernetes.

## Acceptance Criteria

- [ ] A new Helm Chart is created at `deploy/helm/tinybrain-operator`.
- [ ] `values.yaml` exposes standard production variables: `replicaCount`, `image`, `resources`, `leaderElection`, `metrics.port`, etc.
- [ ] The operator Deployment includes `livenessProbe` and `readinessProbe` targeting the manager's health endpoints (`:8081`).
- [ ] A Kubernetes `Service` is created to expose the operator's metrics port (`:8080`).
- [ ] CRDs (`agent`, `task`, `kvcache`, `swappolicy`) are bundled into the chart's `crds/` directory.
- [ ] RBAC ClusterRoles and RoleBindings are properly templated.
- [ ] The Helm Chart passes `helm lint` and `helm template` successfully.

## Technical Implementation Notes

- We will maintain the legacy `deploy/k8s/` structure alongside the new Helm chart during the transition period.
- Ensure the controller-runtime flags (`--metrics-bind-address`, `--health-probe-bind-address`, `--leader-elect`) match the templates in the deployment.

## Related Docs

- [Month 7 Roadmap](../planning/roadmap/months/month-07.md)
- [v2 Architecture](../docs/architecture/v2-architecture.md)
