#!/bin/sh
set -eu

if [ ! -f "src/modern.json" ]; then
  echo "check.sh: src/modern.json not found" >&2
  exit 1
fi

python3 - << 'PY'
import json, sys

with open("src/legacy.json") as f:
    legacy = json.load(f)

with open("src/modern.json") as f:
    modern = json.load(f)

if len(modern) != len(legacy):
    sys.stderr.write(f"check.sh: record count mismatch: {len(modern)} vs {len(legacy)}\n")
    sys.exit(1)

for idx, record in enumerate(modern):
    if not isinstance(record.get("id"), int):
        sys.stderr.write(f"check.sh: record {idx} id is not an integer: {record.get('id')}\n")
        sys.exit(1)
    if "legacy_id" not in record:
        sys.stderr.write(f"check.sh: record {idx} missing legacy_id field\n")
        sys.exit(1)
    orig = legacy[idx]
    if record["legacy_id"] != orig["id"]:
        sys.stderr.write(f"check.sh: record {idx} legacy_id {record['legacy_id']} != {orig['id']}\n")
        sys.exit(1)
    if record["title"] != orig["title"] or record["active"] != orig["active"]:
        sys.stderr.write(f"check.sh: record {idx} data corruption\n")
        sys.exit(1)

print("SCHEMA_MIGRATION_GENESIS_CHECK_PASS")
PY
