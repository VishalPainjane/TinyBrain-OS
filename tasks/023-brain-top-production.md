# Task 023 — brain-top Production Polish

**Status:** in-progress  
**Version:** v1.0  
**Sprint:** Month 6, Week 22  
**Depends on:** 013-brain-top (prototype), 022-system-integration  

## Objective

Harden the prototype `brain-top` TUI into a production-quality terminal dashboard with:

- ANSI colour and styled output
- Unicode box-drawing panel frames
- Four panels: process states, resource bars, queue depths, swap monitor
- Graceful signal handling (Ctrl+C)
- Cross-platform terminal clear
- Comprehensive test coverage

## Acceptance Criteria

- [ ] brain-top renders clean coloured output on Windows and Linux terminals
- [ ] Four panels display: Processes, Resources, Queues, Swap Monitor
- [ ] Signal handling works — Ctrl+C exits cleanly in watch mode
- [ ] Version bumped to 1.0.0
- [ ] Tests pass: `CGO_ENABLED=0 go test ./cmd/brain-top/... -v`

## Deliverables

- `cmd/brain-top/color.go` — ANSI colour/style system
- `cmd/brain-top/panel.go` — Box-drawing frame helpers
- `cmd/brain-top/signal.go` — Signal handling
- `cmd/brain-top/swap.go` — Swap monitor panel
- Modified `main.go` — Context-based watch loop, version bump
- Modified `render.go` — Visual overhaul with all four panels
- Test files for all new code

## References

- [docs/architecture/telemetry.md](../../docs/architecture/telemetry.md)
- [planning/roadmap/months/month-06.md](../../planning/roadmap/months/month-06.md) (Week 22)
- [tasks/013-brain-top.md](013-brain-top.md) (prototype)

---
**Layer:** tasks  
**Created:** 2026-06-24
