# TinyBrain OS — AGENTS.md

## Project Identity

TinyBrain OS is a hardware-aware AI runtime kernel for dynamically orchestrating arbitrary AI agents under resource constraints on local hardware. Agents are plugins; the runtime, scheduler, and process model are the core — not fixed agent roles.

## Architecture Rules

- Agents are plugins — never hardcode agent types in core
- Hardware determines model selection (ADR-001)
- Scheduler never depends on the inference engine
- Runtime never depends on the UI
- Cloud is optional; local-first by default
- Structured JSON IPC — no natural language between components
- Follow [docs/constitution.md](docs/constitution.md) and [docs/architecture/invariants.md](docs/architecture/invariants.md)

## Required Reading Order

1. [planning/README.md](planning/README.md)
2. [planning/execution/current-sprint.md](planning/execution/current-sprint.md)
3. [docs/current.md](docs/current.md)
4. [docs/glossary.md](docs/glossary.md)
5. [docs/constitution.md](docs/constitution.md)
6. [docs/architecture/invariants.md](docs/architecture/invariants.md)
7. Relevant [docs/adr/](docs/adr/) and [docs/contracts/](docs/contracts/)
8. Active [tasks/NNN-*.md](tasks/)

## Coding Rules

- Go 1.22+, module: `github.com/VishalPainjane/TinyBrain-OS`
- Package layout: `internal/process`, `internal/events`, `internal/registry`, `internal/hardware`, `internal/runtime`, `internal/loader`, `internal/scheduler`
- Each package: `interface.go`, `types.go`, service implementation, tests
- No inference imports in scheduler (INV-001)
- No scheduler imports in runtime (INV-002)
- No llama.cpp outside inference adapter package (INV-008)
- Respect forbidden packages in [docs/current.md](docs/current.md)
- New docs must use [docs/templates/](docs/templates/)

## Testing Rules

- Write tests for new behavior; table-driven where appropriate
- Add regression tests when fixing bugs
- Test architectural boundaries, not just happy paths
- Run `go test ./...` before marking tasks complete
- Verify scheduler, runtime, and UI remain decoupled

## Important Locations

| Path | Purpose |
|------|---------|
| [planning/](planning/) | Sprint execution, metrics, releases, retrospectives |
| [docs/constitution.md](docs/constitution.md) | Architectural law |
| [docs/glossary.md](docs/glossary.md) | Term definitions |
| [docs/architecture/](docs/architecture/) | Subsystem descriptions |
| [docs/architecture/invariants.md](docs/architecture/invariants.md) | System truths |
| [docs/contracts/](docs/contracts/) | Interfaces and ownership |
| [docs/specs/](docs/specs/) | Version specifications |
| [docs/adr/](docs/adr/) | Architecture decisions |
| [docs/templates/](docs/templates/) | Document scaffolds |
| [docs/current.md](docs/current.md) | Current version, task, forbidden work |
| [tasks/](tasks/) | Implementation tasks |
| [detail.md](detail.md) | Master spec reference |

## Working Rule

If a new idea conflicts with the constitution or an ADR, revise the idea before implementing it. Log rejected approaches in [planning/decisions/rejected.md](planning/decisions/rejected.md).

## Repository Governance

Behavioral rules live in: `.cursor/rules/`

Development history lives in: `planning/`

Architecture decisions live in: `docs/adr/`

Future proposals live in: `docs/rfc/`

Agents must consult these locations before making changes.
