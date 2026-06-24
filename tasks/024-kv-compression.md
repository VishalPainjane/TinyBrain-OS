# 024 KV Compression Pipeline

**Status:** Not Started
**Owner:** TBD
**Tags:** memory, inference
**Layer:** core

## Context

To prevent VRAM exhaustion during long-running swarm sessions or when heavily context-switching between agents, we need to efficiently offload KV cache state to RAM and NVMe. Currently, the `internal/swap` and `internal/kv` packages handle the orchestration of these memory movements. However, the data transferred is uncompressed, resulting in significant I/O overhead and rapid RAM saturation.

## Acceptance Criteria

- [ ] The `internal/kv` manager compresses memory blocks using Zstandard before they are written out to RAM or NVMe.
- [ ] The `internal/kv` manager decompresses memory blocks when restoring them to VRAM.
- [ ] CPU overhead for compression/decompression does not exceed 15% of total system resources, ensuring the scheduler is not starved.
- [ ] Telemetry events are emitted for `KVCompressed` and `KVDecompressed` with metrics on compression ratio and latency.
- [ ] All tests in `internal/kv` and `internal/swap` pass with the new compression logic enabled.
- [ ] The integration passes the `CGO_ENABLED=0 go test ./...` test suite.

## Technical Implementation Notes

- Consider using a fast, pure-Go implementation of Zstandard to avoid introducing new CGO dependencies unless performance strictly mandates it.
- Profile the compression stage. If it blocks the main event loop, we may need to dispatch compression to a background worker pool managed by the scheduler.

## Related Docs

- [Month 7 Roadmap](../planning/roadmap/months/month-07.md)
- [Architecture: Future State](../planning/architecture-evolution/future-state.md)
