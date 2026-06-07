# ADR-005: Local-First

## Status

Accepted

## Context

Many AI systems assume cloud APIs and large datacenter hardware. TinyBrain targets consumer hardware (16 GB RAM, RTX 3050) and offline-capable operation. Cloud dependency would violate the project thesis and exclude the primary use case.

## Decision

TinyBrain is local-first, edge-first, and offline-capable by default.

- Core operation requires no OpenAI, Anthropic, Gemini, or external API.
- Cloud providers are optional adapters behind `InferenceProvider` — not required dependencies.
- Primary model format: GGUF via llama.cpp on local hardware.
- No Postgres, Redis, or Kafka in MVP — BoltDB/SQLite only when persistence is needed.

## Consequences

### Positive

- Runs on constrained hardware without network.
- Privacy-preserving by default.
- Clear differentiation from cloud-first AI wrappers.
- Simpler MVP deployment.

### Negative

- Model quality bounded by local hardware.
- Cloud fallback requires explicit adapter work.
- Cannot assume unlimited compute for demos.

### Trade-off

Optimize for local practicality; add cloud as optional enhancement, not foundation.

---
**Layer:** decision
**Related:** [../constitution.md](../constitution.md), [../../planning/assumptions.md](../../planning/assumptions.md)
