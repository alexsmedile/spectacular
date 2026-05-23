#!/usr/bin/env bash
# tests/cli/migrate.test.sh — workspace migration framework (v0.6.1+)
#
# Stage 1 ships with two backfilled migrations:
#   v0.4 → v0.5  rename current/ → specs/ (preserve contents, flat or subfolder)
#   v0.5 → v0.6  ensure specs/ exists as always-set
#
# Scenarios:
#   1. Migrate on up-to-date workspace is a clean no-op
#   2. v0.4-shape (Octopus-style: current/ + flat .md files, no workspace_schema)
#      → migrate renames to specs/, preserves flat layout, bumps schema to 0.6
#   3. Idempotence — running migrate twice; second is no-op
#   4. --dry-run lists planned migrations without writing
#   5. Broken state: both current/ AND specs/ present → migrate refuses
#   6. migrate outside a workspace exits non-zero
#   7. migrate --help shows usage

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_file_exists() {
  if [[ -f "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected file: $1"; fail_count=$((fail_count + 1)); fi
}
assert_dir_exists() {
  if [[ -d "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected dir: $1"; fail_count=$((fail_count + 1)); fi
}
assert_dir_absent() {
  if [[ ! -d "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected absent dir: $1"; fail_count=$((fail_count + 1)); fi
}
assert_file_contains() {
  if [[ -f "$1" ]] && grep -qF "$2" "$1"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count + 1)); fi
}
assert_output_contains() {
  if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count + 1)); fi
}
assert_exit() {
  if [[ "$1" -eq "$2" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ $3: exit $1 want $2"; fail_count=$((fail_count + 1)); fi
}

seed_v04_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.spectacular/current" "$dir/.spectacular/requests"
  # Octopus-shape contents: flat SCHEMA-*.md files
  cat > "$dir/.spectacular/current/SCHEMA-TASK.md" <<EOF
---
name: task
---
# Task schema
EOF
  cat > "$dir/.spectacular/current/AXIS-MODEL.md" <<EOF
---
name: axis
---
EOF
  # Required scaffold files (so doctor doesn't complain about other things)
  cat > "$dir/.spectacular/config.yaml" <<EOF
project:
  name: $(basename "$dir")
EOF
  echo "---" > "$dir/.spectacular/PRD.md"
  echo "---" > "$dir/.spectacular/SPEC.md"
}

scenario_1_up_to_date_noop() {
  echo "Scenario 1: migrate on up-to-date workspace is a clean no-op"
  local dir="/tmp/spectacular-migrate-test-1"
  rm -rf "$dir"; mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --kit blank >/dev/null 2>&1)

  local out code
  out=$(cd "$dir" && "$CLI" migrate 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "migrate on fresh workspace exits 0"
  assert_output_contains "$out" "Workspace is up to date"

  rm -rf "$dir"
}

scenario_2_v04_full_migration() {
  echo "Scenario 2: v0.4-shape (Octopus) migrates → specs/ with flat layout preserved"
  local dir="/tmp/spectacular-migrate-test-2"
  seed_v04_workspace "$dir"

  local out code
  out=$(cd "$dir" && "$CLI" migrate 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "migrate v0.4 → v0.6 exits 0"
  assert_output_contains "$out" "renamed .spectacular/current/ → .spectacular/specs/"
  assert_output_contains "$out" "workspace_schema is now 0.6"

  # State assertions
  assert_dir_absent "$dir/.spectacular/current"
  assert_dir_exists "$dir/.spectacular/specs"
  assert_file_exists "$dir/.spectacular/specs/SCHEMA-TASK.md"
  assert_file_exists "$dir/.spectacular/specs/AXIS-MODEL.md"
  assert_file_contains "$dir/.spectacular/config.yaml" 'workspace_schema: "0.6"'

  rm -rf "$dir"
}

scenario_3_idempotence() {
  echo "Scenario 3: running migrate twice — second is no-op"
  local dir="/tmp/spectacular-migrate-test-3"
  seed_v04_workspace "$dir"

  (cd "$dir" && "$CLI" migrate >/dev/null 2>&1)
  # Second run
  local out code
  out=$(cd "$dir" && "$CLI" migrate 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "second migrate exits 0"
  assert_output_contains "$out" "up to date"

  rm -rf "$dir"
}

scenario_4_dry_run() {
  echo "Scenario 4: --dry-run lists planned migrations without writing"
  local dir="/tmp/spectacular-migrate-test-4"
  seed_v04_workspace "$dir"

  local out code
  out=$(cd "$dir" && "$CLI" migrate --dry-run 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "dry-run exits 0"
  assert_output_contains "$out" "Dry-run"
  assert_output_contains "$out" "would rename"

  # State should be UNCHANGED
  assert_dir_exists "$dir/.spectacular/current"
  assert_dir_absent "$dir/.spectacular/specs"
  # workspace_schema should NOT have been added
  if grep -q "^workspace_schema:" "$dir/.spectacular/config.yaml"; then
    echo "    ✗ dry-run wrote workspace_schema to config.yaml"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi

  rm -rf "$dir"
}

scenario_5_broken_state_refuses() {
  echo "Scenario 5: both current/ AND specs/ present → migrate refuses with error"
  local dir="/tmp/spectacular-migrate-test-5"
  seed_v04_workspace "$dir"
  # Force broken state: also create specs/
  mkdir -p "$dir/.spectacular/specs"
  touch "$dir/.spectacular/specs/somefile.md"

  local out code
  out=$(cd "$dir" && "$CLI" migrate 2>&1) && code=0 || code=$?
  if [[ "$code" -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit, got 0"; fail_count=$((fail_count + 1)); fi
  assert_output_contains "$out" "both .spectacular/current/ AND .spectacular/specs/"

  rm -rf "$dir"
}

scenario_6_outside_workspace() {
  echo "Scenario 6: migrate outside a workspace exits non-zero"
  local dir="/tmp/spectacular-migrate-test-6"
  rm -rf "$dir"; mkdir -p "$dir"

  local code
  (cd "$dir" && "$CLI" migrate >/dev/null 2>&1) && code=0 || code=$?
  if [[ "$code" -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit, got 0"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_8_list_flag() {
  echo "Scenario 8: --list shows registered migrations from the registry"
  local out code
  out=$("$CLI" migrate --list 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "--list exits 0"
  assert_output_contains "$out" "Registered migrations"
  assert_output_contains "$out" "v04-to-v05"
  assert_output_contains "$out" "v05-to-v06"
  assert_output_contains "$out" "0.4 → 0.5"
  assert_output_contains "$out" "0.5 → 0.6"
}

scenario_9_to_flag() {
  echo "Scenario 9: --to <ver> stops at intermediate schema version"
  local dir="/tmp/spectacular-migrate-test-9"
  seed_v04_workspace "$dir"

  # Apply only v0.4 → v0.5
  local out code
  out=$(cd "$dir" && "$CLI" migrate --to 0.5 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "--to 0.5 exits 0"
  assert_output_contains "$out" "0.4 → 0.5"
  assert_file_contains "$dir/.spectacular/config.yaml" 'workspace_schema: "0.5"'
  # specs/ should now exist (rename happened) but no .gitkeep (v05→v06 didn't run)
  assert_dir_exists "$dir/.spectacular/specs"

  rm -rf "$dir"
}

scenario_10_downgrade_refused() {
  echo "Scenario 10: --to <older> refuses downgrade"
  local dir="/tmp/spectacular-migrate-test-10"
  rm -rf "$dir"; mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --kit blank >/dev/null 2>&1)

  local out code
  out=$(cd "$dir" && "$CLI" migrate --to 0.4 2>&1) && code=0 || code=$?
  if [[ "$code" -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit on downgrade attempt"; fail_count=$((fail_count + 1)); fi
  assert_output_contains "$out" "Downgrade not supported"

  rm -rf "$dir"
}

scenario_7_help_flag() {
  echo "Scenario 7: migrate --help shows usage"
  local out code
  out=$("$CLI" migrate --help 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "--help exits 0"
  assert_output_contains "$out" "Usage: spectacular migrate"
  assert_output_contains "$out" "--dry-run"
  assert_output_contains "$out" "--to"
  assert_output_contains "$out" "--from"
  assert_output_contains "$out" "--list"
}

echo "=== migrate.test.sh ==="
scenario_1_up_to_date_noop
scenario_2_v04_full_migration
scenario_3_idempotence
scenario_4_dry_run
scenario_5_broken_state_refuses
scenario_6_outside_workspace
scenario_7_help_flag
scenario_8_list_flag
scenario_9_to_flag
scenario_10_downgrade_refused

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
