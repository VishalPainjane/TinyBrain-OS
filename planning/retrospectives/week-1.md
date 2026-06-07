# Week 1 Retrospective — Kernel Foundation

## Completed

- ProcessState type with 7 lifecycle states (NEW, READY, RUNNING, WAITING, PREEMPTED, HIBERNATED, TERMINATED)
- Unit tests for state validity, string representation, and All()
- Full documentation system planned and initialized

## Problems

- Task 001 doc referenced fixed agent names (Planner, Coder) — conflicts with plugin architecture
- `docs/` casing standardized during Month 1 finalization

## Decisions

- Process owns state; scheduler will own transitions ([accepted.md](../decisions/accepted.md))
- Docs must use evolved architecture; fixed agents are examples only
- Code-first truth: documentation reflects implemented state

## Lessons

- Contracts must be written before task 002 to avoid rework
- `planning/` layer needed to separate sprint ops from architecture law
- Glossary and invariants needed to prevent AI term drift
- Templates enforce consistent document structure for future AI sessions

## Next Week

- Process table (002)
- Begin event types design (003)

---
**Layer:** planning
**Date:** 2026-06-07
