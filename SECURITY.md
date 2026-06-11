# Security Policy

TinyBrain OS is solo-developed and local-first. This document states how security is handled in the repository.

## Secrets

- **Never** commit API keys, tokens, passwords, or private keys.
- **Never** commit `.env` files with real credentials. Use environment variables at runtime.
- Assume the repository may become public at any time.

See also: `.cursor/rules/security.mdc`.

## Dependencies

- Pin third-party code where practical (e.g. llama.cpp submodule SHA in CI).
- Review dependency changes before merge; note CVEs in the task or release notes when relevant.
- CGO and native libraries (`libllama.so`) are part of the attack surface — keep pins documented in `testdata/ci/`.

## Reporting vulnerabilities

There is no public bug-bounty program. If you discover a security issue:

1. Do **not** open a public issue with exploit details.
2. Contact the repository owner through GitHub private security advisories or direct message if you have a channel.

For solo development: log the issue in `planning/decisions/` or a private note, fix on a branch, and mention the fix in release notes without exposing exploit steps.

## Supported versions

Only the latest tagged release on `main` receives fixes. Older tags are not maintained unless explicitly noted in release docs.

---

**Last updated:** 2026-06-11
