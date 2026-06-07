# llama.cpp Notes

**Non-normative — do not treat as architecture law.**

## Why Primary Engine

- Native GGUF support
- mmap-backed model loading (lazy page load, not full RAM copy)
- KV cache handling built-in
- GPU offloading (CUDA, Metal, Vulkan)
- Active ecosystem, local-first alignment (ADR-005)

## Quantization

Default: **Q4_K_M** — balance of quality and memory savings on consumer hardware.

Format: **GGUF** — mmap compatible, llama.cpp native.

## Integration Approach

Never import llama.cpp outside adapter package (INV-008). Adapter implements `InferenceProvider` interface (ADR-004). Runtime calls adapter; scheduler unaware.

## CGO Considerations

- Windows toolchain requirements (MSVC, CGO_ENABLED)
- Build complexity — mitigate with stub provider until v0.6
- Risk tracked in [planning/risks/technical-risks.md](../../planning/risks/technical-risks.md)

## mmap

Weights loaded via mmap — pages faulted on demand. Critical for fitting models in constrained RAM/VRAM on Standard profile (16GB + 4GB VRAM).

---
**Layer:** research
