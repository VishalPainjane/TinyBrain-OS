#!/usr/bin/env bash
# ci-emit-metrics.sh — STAB-002 CI observability helpers (source, do not execute directly).
set -euo pipefail

# tb_ci_now prints Unix epoch seconds.
tb_ci_now() {
  date +%s
}

# tb_ci_cache_label maps actions/cache cache-hit output to hit|miss|unknown.
tb_ci_cache_label() {
  case "${1:-}" in
    true) echo "hit" ;;
    false) echo "miss" ;;
    *) echo "unknown" ;;
  esac
}

# tb_ci_env_set appends a key=value pair to GITHUB_ENV.
tb_ci_env_set() {
  echo "${1}=${2}" >> "${GITHUB_ENV:?}"
}

# tb_ci_summary_header writes the job metrics table header to GITHUB_STEP_SUMMARY.
tb_ci_summary_header() {
  local job_name="${1:?job name required}"
  {
    echo "## CI metrics — ${job_name}"
    echo ""
    echo "| Metric | Value |"
    echo "|--------|-------|"
  } >> "${GITHUB_STEP_SUMMARY:?}"
}

# tb_ci_summary_row appends one row to the step summary table.
tb_ci_summary_row() {
  local label="$1"
  local value="$2"
  echo "| ${label} | ${value} |" >> "${GITHUB_STEP_SUMMARY:?}"
}

# tb_ci_failure_phase resolves failure_phase from step outcomes (env STEP_*).
tb_ci_failure_phase() {
  if [ "${JOB_CONCLUSION:-success}" = "success" ]; then
    echo "none"
    return
  fi
  if [ "${STEP_VERIFY_SHA:-}" = "failure" ] || [ "${STEP_VERIFY_TAG:-}" = "failure" ] || [ "${STEP_READ_SHA256:-}" = "failure" ]; then
    echo "bootstrap"
  elif [ "${STEP_BUILD_LLAMA:-}" = "failure" ]; then
    echo "llama_build"
  elif [ "${STEP_DOWNLOAD_GGUF:-}" = "failure" ]; then
    echo "gguf_download"
  elif [ "${STEP_VERIFY_GGUF:-}" = "failure" ] || [ "${STEP_ASSERT_GGUF:-}" = "failure" ]; then
    echo "checksum"
  elif [ "${STEP_INTEGRATION_TEST:-}" = "failure" ]; then
    if [ "${TB_M2_GUARD_FAILED:-}" = "1" ]; then
      echo "m2_guard"
    else
      echo "integration_test"
    fi
  elif [ "${STEP_GO_TEST:-}" = "failure" ] || [ "${STEP_CGO_TEST:-}" = "failure" ]; then
    echo "other"
  else
    echo "other"
  fi
}

# tb_ci_write_job_metrics_json writes ci-job-metrics.json for the collector job.
tb_ci_write_job_metrics_json() {
  local out="${RUNNER_TEMP:?}/ci-job-metrics.json"
  python3 <<'PY'
import json, os

def env(name, default=""):
    return os.environ.get(name, default)

def env_int(name, default=0):
    raw = env(name, str(default))
    try:
        return int(raw)
    except ValueError:
        return default

job = {
    "job_name": env("JOB_NAME"),
    "conclusion": env("JOB_CONCLUSION", "unknown"),
    "setup_go_cache": env("TB_SETUP_GO_CACHE", "unknown"),
    "failure_phase": env("TB_FAILURE_PHASE", "none"),
    "failure_class": env("TB_FAILURE_CLASS", ""),
}

optional_int = [
    ("go_test_duration_s", "TB_GO_TEST_DURATION_S"),
    ("llama_build_duration_s", "TB_LLAMA_BUILD_DURATION_S"),
    ("cgo_test_duration_s", "TB_CGO_TEST_DURATION_S"),
    ("gguf_download_duration_s", "TB_GGUF_DOWNLOAD_DURATION_S"),
    ("gguf_verify_duration_s", "TB_GGUF_VERIFY_DURATION_S"),
    ("integration_test_duration_s", "TB_INTEGRATION_TEST_DURATION_S"),
    ("m2_pass_count", "TB_M2_PASS_COUNT"),
    ("m2_run_count", "TB_M2_RUN_COUNT"),
    ("m2_skip_count", "TB_M2_SKIP_COUNT"),
    ("m2_fail_count", "TB_M2_FAIL_COUNT"),
]
for key, env_key in optional_int:
    if env(env_key):
        job[key] = env_int(env_key)

optional_str = [
    ("llama_build_cache", "TB_LLAMA_BUILD_CACHE"),
    ("gguf_cache", "TB_GGUF_CACHE"),
    ("llama_cache_key", "TB_LLAMA_CACHE_KEY"),
    ("gguf_cache_key", "TB_GGUF_CACHE_KEY"),
]
for key, env_key in optional_str:
    if env(env_key):
        job[key] = env(env_key)

path = os.path.join(os.environ["RUNNER_TEMP"], "ci-job-metrics.json")
with open(path, "w", encoding="utf-8") as f:
    json.dump(job, f, separators=(",", ":"))
print(f"wrote {path}")
PY
}

# tb_ci_emit_job_summary writes the step summary table and job metrics JSON.
tb_ci_emit_job_summary() {
  export TB_FAILURE_PHASE
  TB_FAILURE_PHASE="$(tb_ci_failure_phase)"
  tb_ci_summary_header "${JOB_NAME:?}"
  tb_ci_summary_row "conclusion" "${JOB_CONCLUSION:-unknown}"
  tb_ci_summary_row "failure_phase" "${TB_FAILURE_PHASE}"
  [ -n "${TB_SETUP_GO_CACHE:-}" ] && tb_ci_summary_row "setup-go cache" "${TB_SETUP_GO_CACHE}"
  [ -n "${TB_GO_TEST_DURATION_S:-}" ] && tb_ci_summary_row "go test (s)" "${TB_GO_TEST_DURATION_S}"
  [ -n "${TB_CGO_TEST_DURATION_S:-}" ] && tb_ci_summary_row "cgo test (s)" "${TB_CGO_TEST_DURATION_S}"
  [ -n "${TB_LLAMA_BUILD_CACHE:-}" ] && tb_ci_summary_row "llama build cache" "${TB_LLAMA_BUILD_CACHE}"
  [ -n "${TB_LLAMA_BUILD_DURATION_S:-}" ] && tb_ci_summary_row "llama build (s)" "${TB_LLAMA_BUILD_DURATION_S}"
  [ -n "${TB_GGUF_CACHE:-}" ] && tb_ci_summary_row "GGUF cache" "${TB_GGUF_CACHE}"
  [ -n "${TB_GGUF_DOWNLOAD_DURATION_S:-}" ] && tb_ci_summary_row "GGUF download (s)" "${TB_GGUF_DOWNLOAD_DURATION_S}"
  [ -n "${TB_GGUF_VERIFY_DURATION_S:-}" ] && tb_ci_summary_row "GGUF verify (s)" "${TB_GGUF_VERIFY_DURATION_S}"
  [ -n "${TB_INTEGRATION_TEST_DURATION_S:-}" ] && tb_ci_summary_row "integration test (s)" "${TB_INTEGRATION_TEST_DURATION_S}"
  [ -n "${TB_M2_PASS_COUNT:-}" ] && tb_ci_summary_row "M2 pass/run/skip/fail" "${TB_M2_PASS_COUNT}/${TB_M2_RUN_COUNT:-?}/${TB_M2_SKIP_COUNT:-?}/${TB_M2_FAIL_COUNT:-?}"
  tb_ci_write_job_metrics_json
}
