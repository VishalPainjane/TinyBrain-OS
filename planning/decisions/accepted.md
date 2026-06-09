# Accepted Decisions

Tactical decisions already made. Promote to [docs/adr/](../../docs/adr/) when permanent architecture is confirmed.

| Decision | Reason | Date | Impact |
|----------|--------|------|--------|
| Agents are plugins, not hardcoded classes | Evolved architecture — user-defined fleets | 2026-06-07 | Core never imports agent types |
| Go as core language | Kubernetes ecosystem, controller-runtime, concurrency | 2026-06-07 | All kernel code in Go |
| Event bus v1 uses Go channels | MVP simplicity; few hundred events/sec | 2026-06-07 | No NATS/Kafka until scale requires |
| Scheduler skeleton starts FIFO | Working queue before MLFQ complexity | 2026-06-07 | MLFQ deferred to v0.7 spec |
| llama.cpp via InferenceProvider adapter | Hexagonal boundary (ADR-004) | 2026-06-07 | Scheduler never touches engine |
| Code-first truth in all docs | Documentation reflects implemented state | 2026-06-07 | No planned-state masquerading as done |
| Process owns state; scheduler owns transitions | Clear ownership split | 2026-06-07 | State in kernel; transitions in scheduler later |
| GGUF Q4_K_M as default quantization | Balance of quality and memory on consumer hardware | 2026-06-07 | Model registry defaults |
| Event bus v1 has no ordering guarantee | Async goroutine dispatch per subscriber; MVP scale | 2026-06-07 | Publishers must not assume FIFO delivery |
| Repository governance frozen in .cursor/rules/ | AGENTS.md + 7 rule files + CI as institutional memory | 2026-06-07 | All agents follow github-workflow.mdc |
| EventBus unsubscribe via returned func | Equivalent to Unsubscribe; avoids subscription ID registry in v1 | 2026-06-07 | `Subscribe(eventType, handler) (unsubscribe func())` |
| Month 1 git history as single foundation release | Honest audit trail; no fabricated per-task commits | 2026-06-07 | Tag `v0.3` represents v0.1+v0.2+v0.3 |
| v0.5 model registry persistence uses bbolt | Key-value access matches map semantics; pure Go; migration Path 3 | 2026-06-08 | `go.etcd.io/bbolt` in `internal/registry/` only |
| models.yaml seeds bbolt store only when empty | Bootstrap on first run; no hot reload in v0.5 | 2026-06-08 | `NewBboltModelRegistry(dbPath, seedPath)` |
| ModelStore interface with InMemoryStore + BboltStore | Adapter swap; gob encode in BboltStore (no separate JSON persistence layer) | 2026-06-08 | `ModelRegistry` delegates to `ModelStore` |
| Cross-platform architecture mandatory | Windows/Linux/macOS + CPU/CUDA/ROCm/Metal/Vulkan; core packages OS-agnostic; inference isolated in `internal/inference/` | 2026-06-08 | [docs/architecture/cross-platform.md](../../docs/architecture/cross-platform.md) |
| Inference backend capability matrix | Track OS, hardware, limitations, build, CI per backend; no assumed parity | 2026-06-08 | [docs/architecture/inference-backend-matrix.md](../../docs/architecture/inference-backend-matrix.md) |
| ModelResolver port for inference | Inference must not import registry; adapter at composition boundary | 2026-06-08 | [009a-registry-resolver.md](009a-registry-resolver.md) |
| Inference lifecycle canonical doc | Single state ownership across runtime, loader, inference | 2026-06-08 | [docs/architecture/inference-lifecycle.md](../../docs/architecture/inference-lifecycle.md) |
| Backend build tags mutually exclusive | One backend per binary; see 009a-build-tags.md | 2026-06-08 | [009a-build-tags.md](009a-build-tags.md) |
| llama.cpp git submodule | Pinned upstream tag `b9553` @ `9e3b928fd8c9d14dbf15a8768b9fdd7e5c721d66` in `third_party/llama.cpp`; scoped CI library-only build | 2026-06-08 | [009a-llama-cpp-dependency.md](009a-llama-cpp-dependency.md) |
| Real Inference CI Gate | Deterministic checksum-verified GGUF + dual merge-blocking integration jobs (`inference-integration`, `inference-integration-runtime`); no `vars.TB_TEST_GGUF_PATH` | 2026-06-09 | [real-inference-ci-gate-architecture-review.md](real-inference-ci-gate-architecture-review.md) |

---
**Layer:** planning
**Related:** [rejected.md](rejected.md), [../assumptions.md](../assumptions.md)
