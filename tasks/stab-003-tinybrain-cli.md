# Task STAB-003 — tinybrain CLI

## Status

Complete

## Goal

Ship `cmd/tinybrain` as the composition root: terminal-first UX wiring existing kernel packages into runnable commands.

## Context

v0.6 inference works in tests but has no user-facing binary. Env vars, registry paths, and provider wiring belong at `cmd/` per 009c/009d ADRs.

## Requirements

- `tinybrain doctor` — CGO, llama build, registry DB, model path checks
- `tinybrain probe` — hardware profile; `--json` for scripting
- `tinybrain models list` — bbolt registry read
- `tinybrain run --model ID --prompt TEXT` — integrated load → generate → unload
- `tinybrain status` — read-only snapshot (profile, registry, CGO)
- Config via env: `TB_MODELS_DB`, `TB_MODELS_SEED`, `TB_NGLAYERS`, `TB_LLAMA_LIB_DIR`
- Default data dir: `~/.tinybrain/`

## Files

- `cmd/tinybrain/` — main + subcommands
- `tasks/stab-003-tinybrain-cli.md`
- `README.md` — CLI quick start

## Acceptance Criteria

- [x] `go build -o tinybrain ./cmd/tinybrain` succeeds with `CGO_ENABLED=0` (stub inference)
- [x] `tinybrain doctor` prints ok/warn/fail lines; non-zero exit on fail
- [x] `tinybrain probe` and `tinybrain probe --json` work
- [x] `tinybrain models list` reads registry (empty or seeded)
- [x] `tinybrain run` prints load/gen/unload phases; requires CGO + GGUF for real output
- [x] `tinybrain status` prints one-screen snapshot
- [x] `go test ./...` passes (`CGO_ENABLED=0`)
- [x] README documents CLI commands

## Out Of Scope

- brain-top TUI
- Scheduler wiring
- REST/gRPC API
- Telemetry package extraction
- Agent plugins

## Related

- Release: [planning/releases/v0.6.md](../planning/releases/v0.6.md)
- Architecture: [docs/architecture/telemetry.md](../docs/architecture/telemetry.md) (future brain-top)

---
**Layer:** task
