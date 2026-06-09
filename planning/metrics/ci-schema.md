# CI Run Metrics — Artifact Schema

Schema for the **`ci-run-record-{run_id}`** workflow artifact: one JSON object per **`main` push** CI run, produced by `ci-metrics-collect`. Schema version: **1**.

**Source of truth:** GitHub Actions artifacts and step summaries — not an in-repo file.

## Top-level fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schema_version` | int | yes | Always `1` for this schema |
| `run_id` | string | yes | GitHub Actions run id |
| `run_attempt` | int | yes | Run attempt number |
| `workflow` | string | yes | `"CI"` |
| `event` | string | yes | `push` (v1 scope) |
| `ref` | string | yes | e.g. `refs/heads/main` |
| `head_sha` | string | yes | Commit SHA (short or full) |
| `conclusion` | string | yes | `success` \| `failure` \| `cancelled` \| `skipped` |
| `created_at` | string | yes | ISO 8601 UTC |
| `updated_at` | string | yes | ISO 8601 UTC |
| `workflow_duration_s` | int | no | Wall time for full workflow |
| `jobs` | array | yes | Per-job objects (see below) |

## Job object fields

| Field | Type | Jobs | Description |
|-------|------|------|-------------|
| `job_name` | string | all | `test`, `inference-cgo`, `inference-integration`, `inference-integration-runtime` |
| `conclusion` | string | all | Job conclusion |
| `duration_s` | int | all | GitHub-reported job duration |
| `setup_go_cache` | string | all | `hit` \| `miss` \| `unknown` |
| `go_test_duration_s` | int | test, inference-cgo | Primary test step |
| `llama_build_cache` | string | integration, inference-cgo | `hit` \| `miss` (integration only meaningful for cache) |
| `llama_build_duration_s` | int | integration, inference-cgo | Seconds; 0 on integration cache hit |
| `gguf_cache` | string | integration | `hit` \| `miss` |
| `gguf_download_duration_s` | int | integration | 0 when file already present |
| `gguf_verify_duration_s` | int | integration | sha256 verify step |
| `integration_test_duration_s` | int | integration | M1+M2 block wall time |
| `cgo_test_duration_s` | int | inference-cgo | CGO unit tests |
| `m2_pass_count` | int | integration | From M2 summary |
| `m2_run_count` | int | integration | From M2 summary |
| `m2_skip_count` | int | integration | From M2 summary |
| `m2_fail_count` | int | integration | From M2 summary |
| `failure_phase` | string | all | `none` or phase enum (see review §9) |
| `failure_class` | string | all | Short enum: `SIGILL`, `pass_count_low`, `sha256_mismatch`, etc.; empty on success |
| `llama_cache_key` | string | integration | Env snapshot |
| `gguf_cache_key` | string | integration | Env snapshot |

## Enums

**`failure_phase`:** `none`, `bootstrap`, `llama_build`, `gguf_download`, `checksum`, `integration_test`, `m2_guard`, `other`

**Cache fields:** `hit`, `miss`

## Example record

Artifact: `ci-run-record-27204998799` (file: `ci-run-record.json`).

```json
{"schema_version":1,"run_id":"27203437159","run_attempt":1,"workflow":"CI","event":"push","ref":"refs/heads/main","head_sha":"e293b98","conclusion":"success","created_at":"2026-06-09T12:00:00Z","updated_at":"2026-06-09T12:08:00Z","workflow_duration_s":480,"jobs":[{"job_name":"test","conclusion":"success","duration_s":45,"setup_go_cache":"hit","go_test_duration_s":38,"failure_phase":"none","failure_class":""},{"job_name":"inference-cgo","conclusion":"success","duration_s":420,"setup_go_cache":"hit","llama_build_duration_s":360,"cgo_test_duration_s":25,"failure_phase":"none","failure_class":""},{"job_name":"inference-integration","conclusion":"success","duration_s":180,"setup_go_cache":"hit","llama_build_cache":"hit","llama_build_duration_s":0,"gguf_cache":"hit","gguf_download_duration_s":0,"gguf_verify_duration_s":2,"integration_test_duration_s":45,"m2_pass_count":13,"m2_run_count":13,"m2_skip_count":0,"m2_fail_count":0,"failure_phase":"none","failure_class":"","llama_cache_key":"llama-cpu-b9553-ubuntu-portable-v2","gguf_cache_key":"smollm2-135m-instruct-q4_k_m-v1"},{"job_name":"inference-integration-runtime","conclusion":"success","duration_s":175,"setup_go_cache":"hit","llama_build_cache":"hit","llama_build_duration_s":0,"gguf_cache":"hit","gguf_download_duration_s":0,"gguf_verify_duration_s":2,"integration_test_duration_s":40,"m2_pass_count":1,"m2_run_count":1,"m2_skip_count":0,"m2_fail_count":0,"failure_phase":"none","failure_class":"","llama_cache_key":"llama-cpu-b9553-ubuntu-portable-v2","gguf_cache_key":"smollm2-135m-instruct-q4_k_m-v1"}]}
```

## Downloading run history

```bash
gh run list --branch main --workflow CI --limit 7
gh run download <run_id> -n ci-run-record-<run_id> -D /tmp/ci-metrics
```

Per-job detail is also available as artifacts `ci-metrics-{job_name}` and in each job's GitHub Actions step summary.

## Privacy and allowlist

- Do **not** record: `CI_GGUF_URL`, runner temp paths, Hugging Face tokens, env secrets.
- Model identity: filename + checksum file hash only (via cache key suffix).

---
**Layer:** planning  
**Related:** [ci-baseline.md](ci-baseline.md), [../decisions/ci-observability-architecture-review.md](../decisions/ci-observability-architecture-review.md)
