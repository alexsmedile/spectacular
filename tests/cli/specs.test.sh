#!/usr/bin/env bash
# tests/cli/specs.test.sh — covers v0.5.0 spec-rename: SPEC.md + specs/ + legacy migration.
#
# Scenarios:
#   1. Fresh init scaffolds SPEC.md + specs/ (with .gitkeep), no legacy current/
#   2. Init with --kit coding still scaffolds SPEC.md alongside STACK + ARCHITECTURE
#   3. Doctor `specs` area passes on a clean v0.5.0 workspace
#   4. Doctor detects legacy current/ (no specs/) → warning, mechanical fix available
#   5. `doctor specs --fix` migrates current/ → specs/, preserving contents
#   6. Conflict: both current/ AND specs/ present → error, no auto-fix
#   7. Per-capability specs/<cap>/SPEC.md is recognized + validated
#   8. Repeat init is non-destructive (SPEC.md not overwritten)

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
LOCAL_SKILL="$REPO_ROOT/skills/spectacular"

fail_count=0
pass_count=0

assert_file_exists() {
  if [[ -f "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected file: $1"; fail_count=$((fail_count + 1)); fi
}
assert_file_absent() {
  if [[ ! -e "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected absent: $1"; fail_count=$((fail_count + 1)); fi
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
assert_exit() {
  local got="$1" want="$2" label="$3"
  if [[ "$got" -eq "$want" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ $label: exit $got want $want"; fail_count=$((fail_count + 1)); fi
}

seed_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.agents/skills" "$dir/.spectacular"
  ln -s "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"
  cat > "$dir/.spectacular/skills.lock" <<EOF
spectacular:
  ref: local-dev
  sha: local
  url: file://${LOCAL_SKILL}
EOF
}

# Build a "v0.4.x legacy" workspace (current/ instead of specs/+SPEC.md)
seed_legacy() {
  local dir="$1"
  seed_workspace "$dir"
  mkdir -p "$dir/.spectacular/requests" "$dir/.spectacular/current"
  cat > "$dir/.spectacular/PRD.md" <<EOF
---
version: 1.1
updated: 2026-05-22
summary: "legacy"
kit: blank
---
# Legacy
EOF
  cat > "$dir/.spectacular/config.yaml" <<EOF
project:
  name: legacy
  summary: ""
agents:
  file: AGENTS.md
EOF
  cat > "$dir/.spectacular/AGENTS.md" <<EOF
---
version: 1.0
summary: "agents"
---
# Agents
EOF
}

run_cli() {
  local dir="$1"; shift
  (cd "$dir" && "$CLI" "$@" 2>&1)
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_fresh_init_scaffolds_spec() {
  echo "Scenario 1: fresh init scaffolds SPEC.md + specs/ (no current/)"
  local dir="/tmp/spectacular-specs-test-1"
  seed_workspace "$dir"

  run_cli "$dir" init >/dev/null

  assert_file_exists "$dir/.spectacular/SPEC.md"
  assert_dir_exists "$dir/.spectacular/specs"
  assert_file_exists "$dir/.spectacular/specs/.gitkeep"
  assert_dir_absent "$dir/.spectacular/current"
  assert_file_contains "$dir/.spectacular/SPEC.md" "System Spec"

  rm -rf "$dir"
}

scenario_2_kit_coding_still_has_spec() {
  echo "Scenario 2: --kit coding scaffolds SPEC.md alongside STACK + ARCHITECTURE"
  local dir="/tmp/spectacular-specs-test-2"
  seed_workspace "$dir"

  run_cli "$dir" init --kit coding >/dev/null

  assert_file_exists "$dir/.spectacular/SPEC.md"
  assert_file_exists "$dir/.spectacular/STACK.md"
  assert_file_exists "$dir/.spectacular/ARCHITECTURE.md"
  assert_dir_exists "$dir/.spectacular/specs"

  rm -rf "$dir"
}

scenario_3_doctor_specs_clean() {
  echo "Scenario 3: doctor specs passes on v0.5.0 workspace"
  local dir="/tmp/spectacular-specs-test-3"
  seed_workspace "$dir"
  run_cli "$dir" init >/dev/null

  local output exit_code
  output=$(run_cli "$dir" doctor specs)
  exit_code=$?

  assert_exit "$exit_code" 0 "doctor specs on clean workspace"
  if echo "$output" | grep -q "specs/ directory present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected pass message"; fail_count=$((fail_count + 1)); fi
  if echo "$output" | grep -q "SPEC.md present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected SPEC.md pass"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_4_detects_legacy_current() {
  echo "Scenario 4: doctor detects legacy current/ → warning"
  local dir="/tmp/spectacular-specs-test-4"
  seed_legacy "$dir"

  local output exit_code
  output=$(run_cli "$dir" doctor specs)
  exit_code=$?

  assert_exit "$exit_code" 1 "doctor exits 1 (warning) on legacy current/"
  if echo "$output" | grep -q "legacy current/"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected legacy current/ warning"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_5_fix_migrates_current_to_specs() {
  echo "Scenario 5: doctor specs --fix renames current/ → specs/ preserving contents"
  local dir="/tmp/spectacular-specs-test-5"
  seed_legacy "$dir"
  mkdir -p "$dir/.spectacular/current/auth"
  cat > "$dir/.spectacular/current/auth/SPEC.md" <<EOF
---
status: stable
summary: "auth"
---
# Auth
EOF

  run_cli "$dir" doctor specs --fix >/dev/null 2>&1 || true

  assert_dir_exists "$dir/.spectacular/specs"
  assert_dir_absent "$dir/.spectacular/current"
  assert_file_exists "$dir/.spectacular/specs/auth/SPEC.md"
  assert_file_contains "$dir/.spectacular/specs/auth/SPEC.md" "# Auth"

  rm -rf "$dir"
}

scenario_6_conflict_both_dirs_no_autofix() {
  echo "Scenario 6: both current/ + specs/ present → error, no auto-fix"
  local dir="/tmp/spectacular-specs-test-6"
  seed_legacy "$dir"
  mkdir -p "$dir/.spectacular/specs"

  local output exit_code
  output=$(run_cli "$dir" doctor specs)
  exit_code=$?

  assert_exit "$exit_code" 2 "doctor exits 2 (error) on conflict"
  if echo "$output" | grep -q "both .spectacular/current"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected conflict message"; fail_count=$((fail_count + 1)); fi

  # --fix should refuse to act on conflict
  run_cli "$dir" doctor specs --fix >/dev/null 2>&1 || true
  assert_dir_exists "$dir/.spectacular/current"
  assert_dir_exists "$dir/.spectacular/specs"

  rm -rf "$dir"
}

scenario_7_per_capability_spec_validated() {
  echo "Scenario 7: per-capability specs/<cap>/SPEC.md is recognized + validated"
  local dir="/tmp/spectacular-specs-test-7"
  seed_workspace "$dir"
  run_cli "$dir" init >/dev/null

  mkdir -p "$dir/.spectacular/specs/billing"
  cat > "$dir/.spectacular/specs/billing/SPEC.md" <<EOF
---
status: stable
summary: "billing"
---
# Billing
EOF

  local output
  output=$(run_cli "$dir" doctor specs)
  if echo "$output" | grep -q "billing/SPEC.md.*capability spec present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected per-capability pass"; fail_count=$((fail_count + 1)); fi

  # Capability without SPEC.md → warning
  mkdir -p "$dir/.spectacular/specs/auth"
  output=$(run_cli "$dir" doctor specs)
  if echo "$output" | grep -q "auth.*missing SPEC.md"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected missing-SPEC warning"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_9_flat_contract_docs_valid() {
  echo "Scenario 9: flat contract docs (top-level .md in specs/) are valid alongside subfolders"
  local dir="/tmp/spectacular-specs-test-9"
  seed_workspace "$dir"
  run_cli "$dir" init >/dev/null

  # Octopus-shape: flat SCHEMA-*.md files
  cat > "$dir/.spectacular/specs/SCHEMA-TASK.md" <<'EOF'
---
name: schema-task
version: 1.0
---
# Task schema
EOF
  cat > "$dir/.spectacular/specs/AXIS-MODEL.md" <<'EOF'
---
name: axis-model
---
# Axis model
EOF
  # Bad: missing frontmatter — should warn
  cat > "$dir/.spectacular/specs/BROKEN.md" <<'EOF'
# no frontmatter
EOF
  # Mixed: also add a capability subfolder
  mkdir -p "$dir/.spectacular/specs/task-management"
  cat > "$dir/.spectacular/specs/task-management/SPEC.md" <<'EOF'
---
name: task-management
---
# Task management
EOF

  local output
  output=$(run_cli "$dir" doctor specs)

  if echo "$output" | grep -q "SCHEMA-TASK.md.*contract doc present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected SCHEMA-TASK.md to be a passing contract doc"; fail_count=$((fail_count + 1)); fi

  if echo "$output" | grep -q "AXIS-MODEL.md.*contract doc present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected AXIS-MODEL.md to be a passing contract doc"; fail_count=$((fail_count + 1)); fi

  if echo "$output" | grep -q "BROKEN.md.*missing frontmatter"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected BROKEN.md to warn about frontmatter"; fail_count=$((fail_count + 1)); fi

  if echo "$output" | grep -q "task-management/SPEC.md.*capability spec present"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected mixed-mode capability spec to validate"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_8_repeat_init_nondestructive() {
  echo "Scenario 8: repeat init is non-destructive (SPEC.md not overwritten)"
  local dir="/tmp/spectacular-specs-test-8"
  seed_workspace "$dir"
  run_cli "$dir" init >/dev/null

  echo "USER CONTENT" >> "$dir/.spectacular/SPEC.md"
  local sentinel="$(md5sum "$dir/.spectacular/SPEC.md" 2>/dev/null | cut -d' ' -f1)"
  [[ -z "$sentinel" ]] && sentinel="$(md5 -q "$dir/.spectacular/SPEC.md" 2>/dev/null)"

  run_cli "$dir" init >/dev/null

  local after="$(md5sum "$dir/.spectacular/SPEC.md" 2>/dev/null | cut -d' ' -f1)"
  [[ -z "$after" ]] && after="$(md5 -q "$dir/.spectacular/SPEC.md" 2>/dev/null)"

  if [[ "$sentinel" == "$after" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ SPEC.md was overwritten on re-init"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_10_spec_drift_vs_archive() {
  echo "Scenario 10: SPEC.md drift flagged when an archived request is newer"
  local dir="/tmp/spectacular-specs-test-10"
  seed_workspace "$dir"
  run_cli "$dir" init >/dev/null

  # Pin SPEC.md updated to an old date
  local spec="$dir/.spectacular/SPEC.md"
  perl -0pi -e 's/^updated:.*$/updated: 2026-01-01/m' "$spec"

  # Archive a request with a newer updated date
  mkdir -p "$dir/.spectacular/archive/late-feature"
  cat > "$dir/.spectacular/archive/late-feature/PLAN.md" <<'EOF'
---
status: verified
updated: 2026-03-01
summary: "shipped later than SPEC.md was touched"
---
# Late Feature
EOF

  local out; out="$(run_cli "$dir" doctor specs)"
  if echo "$out" | grep -q "may be stale.*late-feature"; then pass_count=$((pass_count + 1))
  else echo "    ✗ drift not flagged"; echo "$out" | grep -i spec; fail_count=$((fail_count + 1)); fi

  # Now bump SPEC.md newer than the archive → no drift warning
  perl -0pi -e 's/^updated:.*$/updated: 2026-04-01/m' "$spec"
  out="$(run_cli "$dir" doctor specs)"
  if echo "$out" | grep -q "may be stale"; then
    echo "    ✗ drift still flagged after SPEC.md bumped newer"; fail_count=$((fail_count + 1))
  else pass_count=$((pass_count + 1)); fi

  rm -rf "$dir"
}

# ── run all ───────────────────────────────────────────────────────────────────
echo "=== specs.test.sh ==="
scenario_1_fresh_init_scaffolds_spec
scenario_2_kit_coding_still_has_spec
scenario_3_doctor_specs_clean
scenario_4_detects_legacy_current
scenario_5_fix_migrates_current_to_specs
scenario_6_conflict_both_dirs_no_autofix
scenario_7_per_capability_spec_validated
scenario_8_repeat_init_nondestructive
scenario_9_flat_contract_docs_valid
scenario_10_spec_drift_vs_archive

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
