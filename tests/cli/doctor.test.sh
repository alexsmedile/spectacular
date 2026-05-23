#!/usr/bin/env bash
# tests/cli/doctor.test.sh — smoke tests for `spectacular doctor` CLI.
#
# Each scenario builds an isolated workspace under /tmp/, exercises doctor,
# asserts on findings + exit codes + fix behavior. Tests focus on CLI side
# only (detection + mechanical fixes). Judgment-fix flow is agent-side and
# lives in VERIFY.md scenarios.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
LOCAL_SKILL="$REPO_ROOT/skills/spectacular"

fail_count=0
pass_count=0

assert_exit_code() {
  local expected="$1" actual="$2" desc="$3"
  if [[ "$actual" == "$expected" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — expected exit $expected, got $actual"
    fail_count=$((fail_count + 1))
  fi
}

assert_output_contains() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — output missing: $pattern"
    fail_count=$((fail_count + 1))
  fi
}

assert_output_lacks() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    echo "    ✗ $desc — output should not contain: $pattern"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi
}

assert_file_exists() {
  local path="$1" desc="$2"
  if [[ -f "$path" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — expected file: $path"
    fail_count=$((fail_count + 1))
  fi
}

assert_file_contains() {
  local path="$1" pattern="$2" desc="$3"
  if [[ -f "$path" ]] && grep -qF "$pattern" "$path"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — '$path' should contain '$pattern'"
    fail_count=$((fail_count + 1))
  fi
}

# Build a clean test workspace with always-set + skill seed
seed_clean() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.spectacular/requests" "$dir/.spectacular/specs" "$dir/.agents/skills"
  touch "$dir/.spectacular/specs/.gitkeep"
  ln -s "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"
  cat > "$dir/.spectacular/PRD.md" <<EOF
---
version: 1.1
updated: 2026-05-22
summary: "test workspace"
kit: blank
---
# Test
EOF
  cat > "$dir/.spectacular/SPEC.md" <<EOF
---
version: 1.0
updated: 2026-05-22
summary: "Index of what this system actually is right now"
---
# Test — System Spec
EOF
  cat > "$dir/.spectacular/config.yaml" <<EOF
project:
  name: test
  summary: "x"
agents:
  file: AGENTS.md
EOF
  cat > "$dir/.spectacular/AGENTS.md" <<EOF
---
version: 1.0
updated: 2026-05-22
summary: "agents"
---
# Agents
EOF
  echo ".spectacular.local/" > "$dir/.gitignore"
  cat > "$dir/.spectacular/skills.lock" <<EOF
spectacular:
  ref: local-dev
EOF
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_clean_workspace_exits_zero() {
  echo "Scenario 1: clean workspace exits 0"
  local dir="/tmp/doctor-test-1"
  seed_clean "$dir"

  local out
  out=$(cd "$dir" && "$CLI" doctor 2>&1)
  local code=$?

  assert_exit_code "0" "$code" "clean workspace exits 0"
  assert_output_contains "$out" "0 error(s)" "no errors reported"
  assert_output_contains "$out" "0 warning(s)" "no warnings reported"

  rm -rf "$dir"
}

scenario_2_missing_always_set_file() {
  echo "Scenario 2: missing always-set file → error, exit 2"
  local dir="/tmp/doctor-test-2"
  seed_clean "$dir"
  rm "$dir/.spectacular/PRD.md"

  local out
  out=$(cd "$dir" && "$CLI" doctor 2>&1)
  local code=$?

  assert_exit_code "2" "$code" "missing PRD.md exits 2"
  assert_output_contains "$out" "PRD.md" "report mentions PRD.md"
  assert_output_contains "$out" "always-set file missing" "specific message"

  rm -rf "$dir"
}

scenario_3_missing_file_mechanical_fix() {
  echo "Scenario 3: --fix re-scaffolds missing always-set file"
  local dir="/tmp/doctor-test-3"
  seed_clean "$dir"
  rm "$dir/.spectacular/PRD.md"

  (cd "$dir" && "$CLI" doctor --fix >/dev/null 2>&1) || true

  assert_file_exists "$dir/.spectacular/PRD.md" "PRD.md re-created by --fix"
  assert_file_contains "$dir/.spectacular/PRD.md" "kit:" "re-created file has frontmatter"

  rm -rf "$dir"
}

scenario_4_malformed_frontmatter() {
  echo "Scenario 4: malformed frontmatter → flagged, NOT auto-fixed"
  local dir="/tmp/doctor-test-4"
  seed_clean "$dir"
  echo "no frontmatter at all just garbage" > "$dir/.spectacular/PRD.md"

  local out
  out=$(cd "$dir" && "$CLI" doctor 2>&1)
  local code=$?

  assert_exit_code "2" "$code" "malformed frontmatter exits 2"
  assert_output_contains "$out" "missing frontmatter delimiter" "specific frontmatter error"

  # Confirm --fix does NOT overwrite the malformed file (requires agent)
  local before
  before=$(cat "$dir/.spectacular/PRD.md")
  (cd "$dir" && "$CLI" doctor --fix >/dev/null 2>&1) || true
  local after
  after=$(cat "$dir/.spectacular/PRD.md")
  if [[ "$before" == "$after" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ malformed file should not be auto-fixed by --fix"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$dir"
}

scenario_5_snapshot_gap() {
  echo "Scenario 5: snapshot version gap → warning, exit 1"
  local dir="/tmp/doctor-test-5"
  seed_clean "$dir"
  # Create v1.0 and v1.2 of PRD; doctor should flag missing v1.1
  cp "$dir/.spectacular/PRD.md" "$dir/.spectacular/PRD@v1.0.md"
  cp "$dir/.spectacular/PRD.md" "$dir/.spectacular/PRD@v1.2.md"

  local out
  out=$(cd "$dir" && "$CLI" doctor snapshots 2>&1)
  local code=$?

  assert_exit_code "1" "$code" "snapshot gap exits 1 (warning)"
  assert_output_contains "$out" "PRD@v1.1.md" "gap finding mentions missing version"
  assert_output_contains "$out" "version-sequence gap" "specific message"

  rm -rf "$dir"
}

scenario_6_dangling_link() {
  echo "Scenario 6: broken related: link → warning"
  local dir="/tmp/doctor-test-6"
  seed_clean "$dir"
  # Add a related: pointing to nonexistent file
  mkdir -p "$dir/.spectacular/requests/foo"
  cat > "$dir/.spectacular/requests/foo/PLAN.md" <<EOF
---
status: planned
priority: low
owner: test
updated: 2026-05-22
summary: "test"
related:
  - ../nonexistent/PLAN.md
---
# Foo
EOF

  local out
  out=$(cd "$dir" && "$CLI" doctor links 2>&1)
  local code=$?

  assert_exit_code "1" "$code" "broken link exits 1"
  assert_output_contains "$out" "nonexistent/PLAN.md" "broken link surfaced"

  rm -rf "$dir"
}

scenario_7_lifecycle_active_without_session() {
  echo "Scenario 7: status: active without SESSION.md → warning"
  local dir="/tmp/doctor-test-7"
  seed_clean "$dir"
  mkdir -p "$dir/.spectacular/requests/bar"
  cat > "$dir/.spectacular/requests/bar/PLAN.md" <<EOF
---
status: active
priority: medium
owner: test
updated: 2026-05-22
summary: "test"
---
# Bar

## Validation
- something
EOF

  local out
  out=$(cd "$dir" && "$CLI" doctor lifecycle 2>&1)
  local code=$?

  assert_exit_code "1" "$code" "active without SESSION.md exits 1"
  assert_output_contains "$out" "without SESSION.md" "specific lifecycle message"

  rm -rf "$dir"
}

scenario_8_scoped_area() {
  echo "Scenario 8: scoped area run skips other areas"
  local dir="/tmp/doctor-test-8"
  seed_clean "$dir"
  rm "$dir/.gitignore"  # would normally flag as warning

  local out
  out=$(cd "$dir" && "$CLI" doctor frontmatter 2>&1)

  # Frontmatter area run should not mention .gitignore (which lives in workspace area)
  assert_output_lacks "$out" ".gitignore" "scoped area excludes other areas"

  rm -rf "$dir"
}

scenario_9_json_output() {
  echo "Scenario 9: --format json emits parseable JSON"
  local dir="/tmp/doctor-test-9"
  seed_clean "$dir"
  rm "$dir/.spectacular/PRD.md"  # ensure at least one finding

  local out
  out=$(cd "$dir" && "$CLI" doctor --format json 2>&1)

  assert_output_contains "$out" '"version":' "JSON has version field"
  assert_output_contains "$out" '"findings":' "JSON has findings array"
  assert_output_contains "$out" '"summary":' "JSON has summary object"
  assert_output_contains "$out" '"severity":' "findings have severity"
  assert_output_contains "$out" '"fix_type":' "findings have fix_type"

  # Validate with python json parser to confirm well-formed
  if echo "$out" | python3 -c 'import sys, json; json.load(sys.stdin)' 2>/dev/null; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ JSON output is not valid JSON"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$dir"
}

scenario_10_mechanical_gitignore() {
  echo "Scenario 10: --fix creates/appends .gitignore"
  local dir="/tmp/doctor-test-10"
  seed_clean "$dir"
  rm "$dir/.gitignore"

  (cd "$dir" && "$CLI" doctor --fix >/dev/null 2>&1) || true

  assert_file_exists "$dir/.gitignore" ".gitignore created by --fix"
  assert_file_contains "$dir/.gitignore" ".spectacular.local/" ".spectacular.local/ entry present"

  # Append case: existing file without entry
  rm -rf "$dir"
  seed_clean "$dir"
  echo "node_modules/" > "$dir/.gitignore"

  (cd "$dir" && "$CLI" doctor --fix >/dev/null 2>&1) || true

  assert_file_contains "$dir/.gitignore" "node_modules/" "existing entry preserved"
  assert_file_contains "$dir/.gitignore" ".spectacular.local/" "new entry appended"

  rm -rf "$dir"
}

scenario_11_help_flag() {
  echo "Scenario 11: doctor --help shows usage"
  local out
  out=$("$CLI" doctor --help 2>&1)
  local code=$?

  assert_exit_code "0" "$code" "--help exits 0"
  assert_output_contains "$out" "Usage: spectacular doctor" "help shows usage"
  assert_output_contains "$out" "--fix" "help mentions --fix"
  assert_output_contains "$out" "--format" "help mentions --format"
}

# ── run ───────────────────────────────────────────────────────────────────────

scenario_1_clean_workspace_exits_zero
scenario_2_missing_always_set_file
scenario_3_missing_file_mechanical_fix
scenario_4_malformed_frontmatter
scenario_5_snapshot_gap
scenario_6_dangling_link
scenario_7_lifecycle_active_without_session
scenario_8_scoped_area
scenario_9_json_output
scenario_10_mechanical_gitignore
scenario_11_help_flag

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"

if [[ $fail_count -gt 0 ]]; then
  exit 1
fi
exit 0
