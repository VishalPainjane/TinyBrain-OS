# TinyBrain OS v2: Bare-Metal LLM Inference Architecture

**Author:** Nikhil (and TinyBrain OS Engineering)  
**Model Profile:** TinyLlama-1.1B (1.1 Billion Parameters, FP16)  
**Core Stack:** Go (Control Plane) / C++ & CUDA (Data Plane)

---

## 1. TinyBrain OS

Traditional research frameworks prioritize developer productivity over deployment efficiency. TinyBrain OS moves the serving stack into compiled Go and CUDA components, minimizing runtime overhead in scheduling, memory management, and request orchestration while allowing compute-intensive operations to execute directly on the GPU.

The core directive of this project is to achieve Total AI Sovereignty: a compiled, standalone binary capable of high-speed local inference with zero dependency on massive Python libraries.

## 2. HTTP Layer

The HTTP server is wrapped in an OpenAI API-compatible interface (`/v1/chat/completions`) to achieve plug-and-play compatibility with UI ecosystems like LibreChat. It streams generated tokens directly using Server-Sent Events (SSE). It handles client disconnects gracefully by propagating context cancellation down to the engine tick loop.

## 3. Scheduler

The Go control plane implements an Asynchronous Continuous Batching MLFQ Scheduler. It operates exactly like an enterprise AI API (e.g., vLLM or Triton Inference Server). 

An infinite background "Engine Tick" loop runs independently of the HTTP server.
**Iteration-Level Preemption**: The scheduler mixes **Prefill** sequences (reading new prompts) with **Decode** sequences (generating words) in the exact same batched forward pass.

## 4. Tokenizer

Tokenization is handled upstream of the CGO boundary in pure Go to avoid serialization overhead. The tokenizer generates standard sentencepiece tokens to build the continuous batch and efficiently constructs ChatML templates without allocating multiple strings.

## 5. Request Queue

Requests are continuously pulled from the lock-free ingress channel into the Waiting, Running, and Swapped queues, governed by the scheduler logic. Hysteresis logic prevents VRAM thrashing by maintaining minimum residence quotas for executing sequences.

## 6. Continuous Batch Builder

The Control Plane batches sequences for the GPU. A persistent C-allocated arena (C.malloc) minimizes serialization overhead across the CGO boundary while avoiding garbage collection of inference buffers. Go writes token IDs, positional offsets, and block mappings directly into this raw memory block.

## 7. KV Cache Manager

**Dynamic Block Allocation**: A custom `BlockAllocator` checks out fixed-size KV cache blocks (e.g., block 511) as a sequence grows, and immediately returns them to the free pool the millisecond an EOS (End of Sequence) token is generated.

## 8. CUDA Executor

The C++ backend (`backend.cu`) computes the Transformer architecture through custom-written hardware shaders. 
Model weights are memory-mapped using `mmap`, allowing demand-paged loading into host memory before asynchronous transfer to GPU memory. This minimizes host-side copies while avoiding expensive serialization. 
Dense GEMMs are executed using cuBLAS and CUTLASS-backed kernels configured for FP16 inputs with FP32 accumulation, enabling efficient utilization of NVIDIA Tensor Cores.

## 9. Attention

- **PagedAttention**: The engine implements the industry standard Paged KV Cache to eliminate memory fragmentation.
- **Grouped-Query Attention (GQA)**: The engine reduces KV-cache memory traffic by sharing key/value projections across multiple query heads.
- **Flash-Softmax**: Inspired by FlashAttention, the attention kernel computes numerically stable softmax statistics online without materializing the full attention matrix, significantly reducing temporary memory requirements while eliminating large per-thread scratch buffers.

## 10. Sampler

The Go control plane implements a real-time sampling pipeline that manages Temperature, Top-K, Top-P, Repetition Penalty, and Stop Sequences. This allows the engine to autonomously navigate ChatML conversation trees with human-like creativity and perfect grammatical structure.

## 11. Response Stream

Outputs are streamed seamlessly back up the stack through non-blocking Go channels to the SSE handler, completing the inference lifecycle safely and asynchronously.

---

## 12. Engineering Postmortems

The path to this stable architecture required overcoming several hardcore hardware traps:

- **The 8KB Stack Death Trap**: The original attention kernel allocated a static `float scores[2048]` array, which silently overflowed the hardware thread stack limit, corrupting warp memory. This was patched via the FlashAttention single-pass math.
- **The CGO Misalignment**: A byte-offset miscalculation in the CGO arena passed garbage positional data (Positions: 1901674108) to the C++ engine, spinning the RoPE (Rotary Position Embeddings) into pure noise. Realignment of the pointers restored the AI's contextual vision.
- **The Thread-Launch Lobotomy**: The attention kernel was accidentally launched with `NUM_KV_HEADS` (4 threads) instead of `NUM_HEADS` (32 threads). 87.5% of the neural network was uninitialized memory. Fixing the CUDA dim3 block configuration re-awakened the full parameter set.

## 13. Industry Alignment

The resulting architecture closely mirrors design patterns adopted by modern open-source and industrial inference systems, including paged KV-cache management, continuous batching, and heterogeneous control/data plane separation. These concepts are widely used in systems such as vLLM and TensorRT-LLM.

## 14. The Roadmap Forward

Next steps focus on pushing performance to the absolute hardware limit:

- **Continuous Prefix Caching**: Hash-based, copy-on-write reference counting of physical VRAM blocks to share system prompts across concurrent users.
- **CUDA Graphs**: Capturing execution graphs to completely remove kernel launch overhead.
- **Quantization & Types**: Introducing W8A16 dequantization and FP8 compute support.
- **Advanced Parallelism**: Implementing Tensor Parallelism and Pipeline Parallelism for multi-GPU setups.
- **Asynchronous Transfers**: Using pinned memory and CUDA streams to overlap copy and compute.
- **Speculative Decoding**: For massive TPS improvements.
- **FlashAttention-2**: Upgrading the attention kernel to the latest standard.
