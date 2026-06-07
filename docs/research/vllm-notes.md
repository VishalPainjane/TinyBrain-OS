# vLLM Notes

**Non-normative — do not treat as architecture law.**

## What vLLM Offers

- High-throughput serving with PagedAttention
- Continuous batching
- Strong datacenter/multi-GPU story

## Why Not MVP Primary

- Heavier deployment than llama.cpp for single-user local runtime
- Less aligned with mmap-on-laptop use case
- TinyBrain targets 4GB VRAM consumer GPU, not A100 clusters

## Future Role

Potential secondary `InferenceProvider` adapter for Workstation profile users with multi-GPU setups. Same interface as llama.cpp adapter (ADR-004) — scheduler unchanged.

## When to Revisit

- Workstation profile benchmarks show llama.cpp throughput insufficient
- User demand for concurrent multi-user serving (out of current scope)

---
**Layer:** research
