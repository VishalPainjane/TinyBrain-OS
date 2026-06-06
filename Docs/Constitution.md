# TinyBrain Constitution

This document is the source of truth for TinyBrain’s architecture and design decisions.

## Core Rules

1. Hardware determines model selection.
2. Agents are plugins.
3. The scheduler never depends on the inference engine.
4. The runtime never depends on the UI.
5. Cloud is optional.
6. Local-first design.

## Decision Policy

Any future proposal, feature, or AI-generated suggestion must be checked against this constitution before being accepted.

If a suggestion conflicts with these rules, the suggestion is rejected or revised to fit the constitution.

## Design Priority

When trade-offs appear, TinyBrain should favor:

- local execution over cloud dependency
- hardware awareness over fixed assumptions
- modular agents over tightly coupled systems
- runtime stability over UI convenience
- clear boundaries over hidden coupling

## Purpose

This constitution keeps TinyBrain consistent as it grows. It prevents the project from drifting away from its local-first, hardware-aware runtime identity.
