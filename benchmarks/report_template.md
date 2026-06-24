# V1.0 Benchmark Report: Swarm vs Monolith

**Date:** YYYY-MM-DD
**Hardware Profile:** Standard (or actual profile used)
**CPU:** (Fill CPU)
**GPU:** (Fill GPU)
**Backend:** (CUDA / Metal / CPU)

## Goal
Compare the performance of a multi-agent swarm configuration against a single monolithic agent, validating the Month 6 roadmap target metrics.

## Workload
- **Prompt:** "Explain the theory of relativity in 5 paragraphs."
- **Swarm Model:** `tinyllama-q4` (2048 MB budget per worker)
- **Monolith Model:** `llama-3-8b-monolith` (8192 MB budget)

## Results

| Metric | Swarm (`tinyllama-q4`) | Monolith (`llama-3-8b`) | Target | Pass/Fail |
|--------|------------------------|-------------------------|--------|-----------|
| TTFT | ? s | ? s | < 3 s | ? |
| TPS | ? tok/s | ? tok/s | > 20 tok/s | ? |
| Peak RAM | ? MB | ? MB | N/A | - |
| Peak VRAM | ? MB | ? MB | < 3.6 GB | ? |

## Conclusion

(Summarize the findings. Note whether the swarm successfully achieved lower TTFT and higher TPS on the same hardware profile compared to the monolith, despite potential trade-offs in reasoning quality.)
