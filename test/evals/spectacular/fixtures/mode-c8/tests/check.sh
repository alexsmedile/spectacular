#!/bin/sh
set -eu

RUNNER=""
if [ -f "src/main.py" ]; then
  RUNNER="python3 src/main.py"
elif [ -f "src/queue.py" ]; then
  RUNNER="python3 src/queue.py"
elif [ -f "src/main.js" ]; then
  RUNNER="node src/main.js"
elif [ -f "src/queue.js" ]; then
  RUNNER="node src/queue.js"
elif [ -f "src/main.go" ]; then
  RUNNER="go run ./src/..."
elif [ -f "src/runner.sh" ]; then
  RUNNER="sh src/runner.sh"
else
  for f in src/*; do
    if [ -x "$f" ] && [ ! -d "$f" ]; then
      RUNNER="$f"
      break
    fi
  done
fi

if [ -z "$RUNNER" ]; then
  echo "check.sh: no executable queue runner found in src/" >&2
  exit 1
fi

rm -f dlq.json results.log

# 1. Run queue simulation test in Python
python3 - << PY
import subprocess, json, sys, time, os

# Create sample input jobs
jobs = [
    {"id": "j_low_1", "priority": "low", "fail": False},
    {"id": "j_high_1", "priority": "high", "fail": False},
    {"id": "j_fail_1", "priority": "normal", "fail": True},
    {"id": "j_norm_1", "priority": "normal", "fail": False},
    {"id": "j_high_2", "priority": "high", "fail": False},
]

with open("jobs.json", "w") as f:
    json.dump(jobs, f)

# Execute runner
proc = subprocess.run("$RUNNER --jobs jobs.json", shell=True, capture_output=True, text=True)
if proc.returncode != 0:
    sys.stderr.write(f"Runner failed (exit {proc.returncode}): {proc.stderr}\n")
    sys.exit(1)

# Verify DLQ
if not os.path.exists("dlq.json"):
    sys.stderr.write("dlq.json not generated\n")
    sys.exit(1)

with open("dlq.json") as f:
    dlq = json.load(f)

# Assert failed job is in DLQ with 3 attempts
found_failed = False
for item in dlq:
    if item.get("id") == "j_fail_1":
        found_failed = True
        if item.get("attempts", 0) < 3:
            sys.stderr.write(f"Expected at least 3 attempts for j_fail_1, got {item.get('attempts')}\n")
            sys.exit(1)

if not found_failed:
    sys.stderr.write("j_fail_1 not found in dlq.json\n")
    sys.exit(1)

print("QUEUE_ENGINE_CHECK_PASS=OK")
PY

echo "PRIORITY_QUEUE_DLQ_GENESIS_CHECK_PASS"
