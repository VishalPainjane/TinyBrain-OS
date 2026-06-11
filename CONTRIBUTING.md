# Contributing — Solo Developer Workflow

TinyBrain OS is built by one developer. This guide is a **session workflow**, not community onboarding.

For architecture law, read [AGENTS.md](AGENTS.md). For live status, read [docs/current.md](docs/current.md).

---

## Start every session

1. [planning/execution/current-sprint.md](planning/execution/current-sprint.md)
2. [docs/current.md](docs/current.md)
3. Active task in [tasks/](tasks/)
4. Relevant [docs/contracts/](docs/contracts/) and [docs/specs/](docs/specs/)

Do not implement work outside the active task or forbidden list in `docs/current.md`.

---

## Implement one task

1. Read the task acceptance criteria.
2. List files you will touch; stay in scope.
3. Write tests with the behavior (see [docs/testing-policy.md](docs/testing-policy.md)).
4. Run verification (below).
5. One task → one commit on a feature branch (`feature/task-XXX-name`).
6. Sync docs per [planning/roadmap/update-checklist.md](planning/roadmap/update-checklist.md) Tier A.

---

## Verification tiers

### Fast (every commit)

```bash
CGO_ENABLED=0 go test ./...
```

Runs unit tests, stub inference path, and architecture fitness tests in `tests/`.

### Integration (inference, runtime, or loader changes)

```bash
# Build llama.cpp — see README and testdata/ci/README.md
export LD_LIBRARY_PATH=third_party/llama.cpp/build/bin
export TB_TEST_GGUF_PATH=/path/to/smollm2-135m-instruct-q4_k_m.gguf
CGO_ENABLED=1 go test -tags integration -count=1 ./internal/inference/llama/...
CGO_ENABLED=1 go test -tags integration -count=1 ./internal/runtime/...
```

### CLI (cmd/tinybrain changes)

```bash
go test ./cmd/tinybrain/...
go test -fuzz=Fuzz -fuzztime=30s ./cmd/tinybrain/...
```

Regenerate CLI golden files after intentional output changes:

```bash
TB_UPDATE_GOLDEN=1 go test ./cmd/tinybrain/ -run TestGolden
```

### Manual GPU (CUDA releases only)

[planning/decisions/009d-manual-gpu-checklist.md](planning/decisions/009d-manual-gpu-checklist.md)

---

## Definition of done (task)

- [ ] Acceptance criteria met
- [ ] Direct tests for new behavior
- [ ] Boundary/regression tests if imports or contracts changed
- [ ] `CGO_ENABLED=0 go test ./...` passes
- [ ] Integration tier run when touching inference/runtime
- [ ] Task marked complete; `docs/current.md` synced

---

## Definition of done (release)

Follow [planning/releases/RELEASE-CHECKLIST.md](planning/releases/RELEASE-CHECKLIST.md).

---

## Commit messages

Format: `feat(scope): description`, `fix(scope):`, `test(scope):`, `docs(scope):`, `chore(scope):`

One logical unit per commit. No `update`, `misc`, `final`.

---

## Regression tests

When fixing a bug or contract drift, add:

```text
TestRegression_<ShortName>
```

in the package where the bug lived. The test must fail without the fix.

---

## What not to do

- No architecture changes without ADR + contract updates
- No future-version features (see roadmap and forbidden work in `docs/current.md`)
- No secrets in code or logs
- No marking tasks done without tests
- No release without full CI green and release checklist

---

## Repository policy

Every important behavior must have a direct test. Every architectural boundary must have a negative test in CI. Every release is trusted only after merge-blocking CI and the release checklist pass. False confidence is not acceptable: `CGO_ENABLED=0` green does not prove inference works.

Full policy: [docs/testing-policy.md](docs/testing-policy.md)

---

## Health and trust

CI status, known gaps, and test tiers: [planning/metrics/repo-health.md](planning/metrics/repo-health.md)

Postmortems after versions: [planning/postmortems/](planning/postmortems/)
