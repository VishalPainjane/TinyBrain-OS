# TinyBrain OS — Planning

## What is this folder?

The `planning/` folder is the execution workspace for TinyBrain OS. It tracks sprints, backlog, decisions, risks, metrics, releases, and retrospectives. It is separate from `docs/`, which holds architectural law, descriptions, contracts, and specifications.

**Rule:** Planning files describe *what we are doing and when*. Docs describe *what the system is and why*.

## How to use it

1. Read [execution/current-sprint.md](execution/current-sprint.md) for today's work.
2. Read [../docs/current.md](../docs/current.md) for hard boundaries (forbidden packages).
3. Read the active task in [../tasks/](../tasks/).
4. Read relevant [../docs/contracts/](../docs/contracts/) before writing code.
5. When a task completes: update [execution/completed.md](execution/completed.md), [metrics/velocity.md](metrics/velocity.md), [metrics/progress.md](metrics/progress.md), and [architecture-evolution/current-state.md](architecture-evolution/current-state.md).

## Reading order (strict)

Every AI agent starts here, then reads in this order:

1. [execution/current-sprint.md](execution/current-sprint.md)
2. [../docs/current.md](../docs/current.md)
3. [../docs/glossary.md](../docs/glossary.md)
4. [../docs/constitution.md](../docs/constitution.md)
5. [../docs/architecture/invariants.md](../docs/architecture/invariants.md)
6. Relevant [../docs/adr/](../docs/adr/) and [../docs/contracts/](../docs/contracts/)
7. Active [../tasks/NNN-*.md](../tasks/)
8. [architecture-evolution/current-state.md](architecture-evolution/current-state.md)
9. [assumptions.md](assumptions.md) when evaluating design choices

## Which files are authoritative?

| Question | Authoritative file |
|----------|-------------------|
| What can never be violated? | [../docs/constitution.md](../docs/constitution.md) |
| What must always remain true? | [../docs/architecture/invariants.md](../docs/architecture/invariants.md) |
| What does this term mean? | [../docs/glossary.md](../docs/glossary.md) |
| Why was this decided formally? | [../docs/adr/](../docs/adr/) |
| How does a subsystem work? | [../docs/architecture/](../docs/architecture/) |
| What interfaces exist? | [../docs/contracts/](../docs/contracts/) |
| What version are we building? | [../docs/specs/](../docs/specs/) |
| What am I implementing today? | [../tasks/](../tasks/) |
| What is the current sprint? | [execution/current-sprint.md](execution/current-sprint.md) |
| What packages are forbidden? | [../docs/current.md](../docs/current.md) |
| What tactical choice was made? | [decisions/accepted.md](decisions/accepted.md) |
| What idea was rejected? | [decisions/rejected.md](decisions/rejected.md) |
| What do we assume but haven't proven? | [assumptions.md](assumptions.md) |
| What exists in code today? | [architecture-evolution/current-state.md](architecture-evolution/current-state.md) |
| What is the long-term direction? | [roadmap/master-roadmap.md](roadmap/master-roadmap.md) |
| What happens this month or week? | [roadmap/months/month-0N.md](roadmap/months/) |
| What do I update when something finishes? | [roadmap/update-checklist.md](roadmap/update-checklist.md) |
| How fast are we delivering? | [metrics/velocity.md](metrics/velocity.md) |
| How complete is each version? | [metrics/progress.md](metrics/progress.md) |

## File Ownership Model

Every file belongs to exactly one layer. Do not mix layers.

| Layer | Location | Defines |
|-------|----------|---------|
| LAW | `docs/constitution.md` | What can never be violated |
| DECISIONS | `docs/adr/` | Why something exists |
| ARCHITECTURE | `docs/architecture/` | How a subsystem works |
| CONTRACTS | `docs/contracts/` | Interfaces, responsibilities, ownership |
| SPECS | `docs/specs/` | What version is being built |
| TASKS | `tasks/` | What gets implemented today |
| CURRENT | `docs/current.md` | Current version, sprint, task — nothing else |
| PLANNING | `planning/` | Execution orchestration, history, risks |
| GLOSSARY | `docs/glossary.md` | Canonical term definitions |
| INVARIANTS | `docs/architecture/invariants.md` | System truths |
| TEMPLATES | `docs/templates/` | Structural scaffolds for new docs |

## Templates

New ADRs, contracts, specs, tasks, and RFCs must copy from [../docs/templates/](../docs/templates/). Never invent new document structure.

---
**Layer:** planning
**Last updated:** 2026-06-07
