#!/usr/bin/env python3
"""Merge per-job CI metrics artifacts into one run record JSON file (STAB-002)."""
from __future__ import annotations

import json
import os
import sys
from datetime import datetime, timezone
from pathlib import Path


def main() -> int:
    artifacts_dir = Path(os.environ.get("ARTIFACTS_DIR", "artifacts"))
    out_path = Path(os.environ["CI_RUN_RECORD_PATH"])
    run_id = os.environ["GITHUB_RUN_ID"]
    run_attempt = int(os.environ.get("GITHUB_RUN_ATTEMPT", "1"))
    head_sha = os.environ.get("GITHUB_SHA", "")[:12]
    event = os.environ.get("GITHUB_EVENT_NAME", "")
    ref = os.environ.get("GITHUB_REF", "")

    jobs: list[dict] = []
    for metrics_file in sorted(artifacts_dir.glob("ci-metrics-*/ci-job-metrics.json")):
        with metrics_file.open(encoding="utf-8") as f:
            jobs.append(json.load(f))

    if not jobs:
        print("no job metrics artifacts found", file=sys.stderr)
        return 1

    conclusions = {j.get("conclusion") for j in jobs}
    if "failure" in conclusions:
        conclusion = "failure"
    elif "cancelled" in conclusions:
        conclusion = "cancelled"
    elif conclusions == {"success"}:
        conclusion = "success"
    else:
        conclusion = "unknown"

    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    record = {
        "schema_version": 1,
        "run_id": run_id,
        "run_attempt": run_attempt,
        "workflow": "CI",
        "event": event,
        "ref": ref,
        "head_sha": head_sha,
        "conclusion": conclusion,
        "created_at": now,
        "updated_at": now,
        "jobs": jobs,
    }

    with out_path.open("w", encoding="utf-8") as f:
        json.dump(record, f, separators=(",", ":"))
        f.write("\n")

    print(f"wrote run record {run_id} to {out_path}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
