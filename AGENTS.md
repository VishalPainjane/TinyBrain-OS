# TinyBrain OS — AGENTS.md

## Project Overview

Project: TinyBrain OS

Purpose: a hardware-aware AI runtime kernel for running specialized agents locally on constrained hardware.

This repository is organized around a runtime-first design, not a prompt-wrapper design.

## Architecture Principles

- Agents are plugins.
- Hardware-aware scheduling comes first.
- The system is local-first by default.
- Cloud support is optional, not required.
- The scheduler must stay independent from the inference engine.
- The runtime must stay independent from the UI.
- Keep boundaries clean and coupling low.

## Development Rules

- Follow `docs/constitution.md` as the source of truth.
- Follow all ADRs in `docs/adr/`.
- Keep documents and code small, focused, and composable.
- Prefer interfaces and adapters over hardcoded behavior.
- Do not hardcode model choices into core logic.
- Do not let UI concerns leak into runtime or scheduler code.

## Testing Rules

- Write tests for new behavior.
- Add regression tests when fixing bugs.
- Test architectural boundaries, not just happy paths.
- Validate that scheduler, runtime, and UI remain decoupled.
- Verify that hardware-aware decisions behave as expected.

## Important Documents

- `docs/constitution.md`
- `docs/architecture.md`
- `docs/vision.md`
- `docs/adr/`

## Working Rule

If a new idea conflicts with the constitution or an ADR, revise the idea before implementing it.