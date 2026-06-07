# Roadmap Update Checklist

Canonical list of files to update when work finishes. Use the tier that matches what you just completed.

**End of week → Tier A** | **Version shipped → Tier B** | **Month closed → Tier C**

---

## Tier A — After every TASK completes

| # | File | Action |
|---|------|--------|
| 1 | Active [tasks/NNN-*.md](../../tasks/) | Mark Status: Complete; check all acceptance criteria |
| 2 | [execution/completed.md](../execution/completed.md) | Append: task ID, commit hash, outcome, files touched |
| 3 | [metrics/velocity.md](../metrics/velocity.md) | Fill Actual duration + Variance |
| 4 | [architecture-evolution/current-state.md](../architecture-evolution/current-state.md) | Move component from In Progress → Implemented |
| 5 | [docs/specs/v0.X-*.md](../../docs/specs/) | Check feature and acceptance criterion boxes |
| 6 | [metrics/progress.md](../metrics/progress.md) | Recalculate version completion % |
| 7 | [execution/current-sprint.md](../execution/current-sprint.md) | Move task to Done; set next Current Task |
| 8 | [docs/current.md](../../docs/current.md) | Sync version, sprint, task, forbidden packages (minimal only) |
| 9 | [execution/backlog.md](../execution/backlog.md) | Remove or deprioritize completed task |
| 10 | `go test ./...` | Must pass before marking task done |

---

## Tier B — After every VERSION ships (v0.X complete)

Complete **all Tier A** items for the final task of that version, then:

| # | File | Action |
|---|------|--------|
| 11 | [releases/v0.X.md](../releases/) | Status → Shipped; fill Completed Tasks, Demo, Lessons |
| 12 | [metrics/progress.md](../metrics/progress.md) | Set that version row to 100% |
| 13 | [retrospectives/week-N.md](../retrospectives/) | Log problems, decisions, lessons for the shipping week |
| 14 | [execution/current-sprint.md](../execution/current-sprint.md) | New sprint name, goal, forbidden packages for next version |
| 15 | [docs/current.md](../../docs/current.md) | Bump version; update forbidden list |
| 16 | [README.md](../../README.md) | Update Current Version and roadmap table row |
| 17 | [decisions/accepted.md](../decisions/accepted.md) | Log tactical decisions made during the version |
| 18 | [assumptions.md](../assumptions.md) | Mark assumptions Validated or Invalidated |
| 19 | [risks/technical-risks.md](../risks/technical-risks.md) | Close mitigated risks; add newly discovered risks |

---

## Tier C — After every MONTH closes

Complete **Tier B** for every version shipped that month, then:

| # | File | Action |
|---|------|--------|
| 20 | [roadmap/months/month-0N.md](months/) | Mark month Status: Complete; fill Planned vs Actual |
| 21 | [roadmap/master-roadmap.md](master-roadmap.md) | Update Current Version if changed |
| 22 | [metrics/progress.md](../metrics/progress.md) | Review all version rows against month targets |
| 23 | [retrospectives/week-N.md](../retrospectives/) | Write month-end engineering retrospective |
| 24 | [architecture-evolution/current-state.md](../architecture-evolution/current-state.md) | Full audit: compare doc vs implemented code |
| 25 | [architecture-evolution/future-state.md](../architecture-evolution/future-state.md) | Adjust if direction shifted during the month |
| 26 | [tasks/](../../tasks/) | Create or unlock tasks for next month if not yet present |

---

## Quick reference

| Event | Tier |
|-------|------|
| Finished a task | A |
| Shipped v0.1, v0.2, … v1.0 | B |
| Closed Month 1–6 | C |

---
**Layer:** planning
**Related:** [master-roadmap.md](master-roadmap.md), [months/](months/)
