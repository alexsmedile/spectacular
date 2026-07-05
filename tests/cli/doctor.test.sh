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
  cat > "$dir/.spectacular/POLICY.md" <<EOF
---
version: 1.0
updated: 2026-05-22
summary: "Operating policies — the practice layer paired with PRINCIPLES.md"
---
# Test — Operating Policies

## @Implementation

### understand-before-change
- principle: 7
- severity: block
- check: PLAN.md has a filled \`## Understanding\` section, OR a \`UNDERSTANDING.md\` exists

A request must not move planned → active until understanding is written down.
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
  # Create v1.0 and v1.2 of PRD in the v1.5.0+ layout; doctor flags missing v1.1.
  # (Root-level <DOC>@v<N>.md is now a legacy-migration warning, not gap detection.)
  mkdir -p "$dir/.spectacular/snapshots/PRD"
  cp "$dir/.spectacular/PRD.md" "$dir/.spectacular/snapshots/PRD/@v1.0.md"
  cp "$dir/.spectacular/PRD.md" "$dir/.spectacular/snapshots/PRD/@v1.2.md"

  local out
  out=$(cd "$dir" && "$CLI" doctor snapshots 2>&1)
  local code=$?

  assert_exit_code "1" "$code" "snapshot gap exits 1 (warning)"
  assert_output_contains "$out" "PRD/@v1.1.md" "gap finding mentions missing version"
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

scenario_8b_multi_area() {
  echo "Scenario 8b: multiple scoped areas run together (regression for status.md substrate trigger)"
  local dir="/tmp/doctor-test-8b"
  seed_clean "$dir"
  # Break frontmatter on PRD so frontmatter area would fire
  echo "no frontmatter here" > "$dir/.spectacular/PRD.md"

  local out
  out=$(cd "$dir" && "$CLI" doctor workspace frontmatter kits 2>&1)

  # All three areas should have run — output should mention each
  assert_output_contains "$out" "workspace" "workspace area ran"
  assert_output_contains "$out" "frontmatter" "frontmatter area ran"
  assert_output_contains "$out" "kits" "kits area ran"
  assert_output_contains "$out" "missing frontmatter delimiter" "frontmatter error surfaced when paired with other areas"

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

scenario_16_roadmap_icebox_rename() {
  echo "Scenario 16: doctor workspace flags pre-v0.7.2 Bucket-list ROADMAP; silent on Icebox"
  local dir="/tmp/doctor-test-16"
  seed_clean "$dir"

  # ROADMAP with old "Bucket list" heading (otherwise structured-shape so check 15 stays silent)
  cat > "$dir/.spectacular/ROADMAP.md" <<'EOF'
---
version: 2.0
updated: 2026-05-23
summary: "test"
related: []
---

# Test — Roadmap

## v1 — first

**Tier:** full
**Status:** active
**Phase:** mvp

**Scope (in):**
- a

## Bucket list

- something
EOF

  local out_bucket
  out_bucket=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_contains "$out_bucket" "Bucket list" \
    "pre-v0.7.2 Bucket-list ROADMAP triggers info"
  assert_output_contains "$out_bucket" "Icebox" \
    "info line names the convergent idiom 'Icebox'"

  # After rename → silent on the Bucket/Icebox check
  sed -i.bak 's/^## Bucket list/## Icebox/' "$dir/.spectacular/ROADMAP.md"
  rm -f "$dir/.spectacular/ROADMAP.md.bak"
  local out_icebox
  out_icebox=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_lacks "$out_icebox" "Bucket list" \
    "renamed ROADMAP no longer triggers Bucket-list info"

  rm -rf "$dir"
}

scenario_15_roadmap_shape_detection() {
  echo "Scenario 15: doctor workspace flags pre-v0.7.1 freeform ROADMAP; silent on structured"
  local dir="/tmp/doctor-test-15"
  seed_clean "$dir"

  # Old-shape ROADMAP: has ## v headings but no **Phase:** markers
  cat > "$dir/.spectacular/ROADMAP.md" <<'EOF'
---
version: 1.0
updated: 2026-01-01
summary: "old"
related: []
---

# Test — Roadmap

## v1 — old

**Status:** planned

- bullet
EOF

  local out_old
  out_old=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_contains "$out_old" "pre-v0.7.1 freeform shape" \
    "old-shape ROADMAP triggers info line"

  # New-shape ROADMAP: has **Phase:** markers
  cat > "$dir/.spectacular/ROADMAP.md" <<'EOF'
---
version: 2.0
updated: 2026-05-23
summary: "new"
related: []
---

# Test — Roadmap

## v1 — new

**Tier:** full
**Status:** planned
**Phase:** mvp

**Scope (in):**
- thing
EOF

  local out_new
  out_new=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_lacks "$out_new" "pre-v0.7.1 freeform shape" \
    "new-shape ROADMAP triggers no info line"

  rm -rf "$dir"
}

scenario_13_migration_chain_validates() {
  echo "Scenario 13: kits area validates migration chain — clean chain passes"
  local dir="/tmp/doctor-test-13"
  seed_clean "$dir"

  local out
  out=$(cd "$dir" && "$CLI" doctor kits 2>&1)
  assert_output_contains "$out" "migration chain validates" "clean registry passes"
}

scenario_14_migration_chain_gap() {
  echo "Scenario 14: kits area detects chain gap in migration registry"
  local dir="/tmp/doctor-test-14"
  seed_clean "$dir"

  # Inject a broken migration: from: "0.9" (no predecessor)
  local target_skill="$dir/.agents/skills/spectacular"
  # The symlink points to the real skill — we need a real copy to mutate
  rm "$target_skill"
  mkdir -p "$target_skill/references/migrations"
  cp "$LOCAL_SKILL/references/migrations/v04-to-v05.md" "$target_skill/references/migrations/"
  cp "$LOCAL_SKILL/references/migrations/v05-to-v06.md" "$target_skill/references/migrations/"
  cat > "$target_skill/references/migrations/v09-to-v10.md" <<'EOF'
---
id: v09-to-v10
from: "0.9"
to: "1.0"
description: "Bogus migration with no predecessor"
mechanical: true
reversible: false
apply-fn: migration_apply_v09_to_v10
affects: []
---
EOF
  # Required: also a templates/prd/kits/ dir for check_kits not to bail early
  mkdir -p "$target_skill/templates/prd/kits"
  cp "$LOCAL_SKILL/templates/prd/kits/blank.md" "$target_skill/templates/prd/kits/" 2>/dev/null || true

  local out code
  out=$(cd "$dir" && "$CLI" doctor kits 2>&1) && code=0 || code=$?
  assert_output_contains "$out" "chain gap" "gap detected"
  # apply-fn not defined → also an error
  assert_output_contains "$out" "not defined in cli/spectacular" "missing apply-fn detected"

  rm -rf "$dir"
}

scenario_12_v06_scaffold_suggestion() {
  echo "Scenario 12: v0.6+ scaffold suggestion surfaces missing PRINCIPLES/ARCHITECTURE/ROADMAP as one info line"
  local dir="/tmp/doctor-test-12"
  seed_clean "$dir"

  # blank kit → all three missing
  local out_all
  out_all=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_contains "$out_all" "v0.6+ conventional files missing: PRINCIPLES, ARCHITECTURE, ROADMAP" \
    "all three missing → single info line lists all three"
  assert_output_contains "$out_all" "spectacular init --with principles,architecture,roadmap" \
    "info line suggests init --with command"

  # Only ROADMAP missing → list just ROADMAP
  echo "---" > "$dir/.spectacular/PRINCIPLES.md"
  echo "---" > "$dir/.spectacular/ARCHITECTURE.md"
  local out_one
  out_one=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_contains "$out_one" "v0.6+ conventional files missing: ROADMAP" \
    "only ROADMAP missing → info line lists just ROADMAP"
  assert_output_lacks "$out_one" "PRINCIPLES, ARCHITECTURE" \
    "info line does not list files that exist"

  # All three present → silent
  echo "---" > "$dir/.spectacular/ROADMAP.md"
  local out_none
  out_none=$(cd "$dir" && "$CLI" doctor workspace 2>&1)
  assert_output_lacks "$out_none" "v0.6+ conventional files missing" \
    "all three present → no v0.6 info line"

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

# Helper: write a SKILL.md whose `description` literal block totals ~N chars.
_write_skill_with_desc() {
  local path="$1" n="$2"
  {
    echo "---"
    echo "name: spectacular"
    echo "description: |"
    local remaining="$n"
    local first=1
    while (( remaining > 0 )); do
      local take=$(( remaining < 90 ? remaining : 90 ))
      printf '  '; printf 'x%.0s' $(seq 1 "$take"); printf '\n'
      remaining=$(( remaining - take ))
      (( remaining > 0 )) && remaining=$(( remaining - 1 ))   # joining newline
      first=0
    done
    echo "version: 1.0.0"
    echo "---"
    echo "# Spectacular"
  } > "$path"
}

scenario_17_description_length() {
  echo "Scenario 17: SKILL.md description vs Codex 1024-char cap (error >1024, warning >1000)"
  local dir="/tmp/doctor-test-17"
  seed_clean "$dir"
  # Replace the symlinked skill with a real dir we can mutate, preserving the
  # bundled templates the kits area expects (copy from the live skill).
  rm "$dir/.agents/skills/spectacular"
  cp -R "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"
  local skill_md="$dir/.agents/skills/spectacular/SKILL.md"

  # over the cap → error
  _write_skill_with_desc "$skill_md" 1100
  local out_over
  out_over=$(cd "$dir" && "$CLI" doctor skill 2>&1)
  assert_output_contains "$out_over" "over Codex's 1024 cap" "over-limit description flagged as error"

  # in the warning band → warning, not error
  _write_skill_with_desc "$skill_md" 1010
  local out_warn
  out_warn=$(cd "$dir" && "$CLI" doctor skill 2>&1)
  assert_output_contains "$out_warn" "trim soon" "near-limit description flagged as warning"

  # comfortably under → no description finding among errors/warnings
  _write_skill_with_desc "$skill_md" 800
  local out_ok
  out_ok=$(cd "$dir" && "$CLI" doctor skill 2>&1)
  assert_output_contains "$out_ok" "description length ok" "short description passes"

  rm -rf "$dir"
}

scenario_18_debug_spine_validation() {
  echo "Scenario 18: debug/ trace spines — enum + invariant validation (v1.26.0+)"
  local dir="/tmp/doctor-test-18"
  seed_clean "$dir"
  mkdir -p "$dir/.spectacular/debug/clean-job" \
           "$dir/.spectacular/debug/bad-status" \
           "$dir/.spectacular/debug/bad-invariant"

  # clean spine → passes
  echo '{"slug":"clean-job","status":"investigating","symptom_class":"runtime_error"}' \
    > "$dir/.spectacular/debug/clean-job/job.json"
  # off-enum status (a 'reason' value leaked into the status slot — the real orchestrator bug)
  echo '{"slug":"bad-status","status":"needs-more-context"}' \
    > "$dir/.spectacular/debug/bad-status/job.json"
  # invariant: wont-fix must log no fix
  echo '{"slug":"bad-invariant","status":"resolved"}' \
    > "$dir/.spectacular/debug/bad-invariant/job.json"
  echo '{"disposition":"wont-fix","logged_fixes":["F9"]}' \
    > "$dir/.spectacular/debug/bad-invariant/outcome.json"

  local out code
  out=$(cd "$dir" && "$CLI" doctor debug 2>&1) && code=0 || code=$?

  assert_exit_code "1" "$code" "debug drift exits 1 (warnings)"
  assert_output_contains "$out" "job status 'needs-more-context' is off-enum" "off-enum status flagged"
  assert_output_contains "$out" "logged_fixes is non-empty" "wont-fix invariant flagged"

  # clean-only workspace passes with exit 0
  rm -rf "$dir/.spectacular/debug/bad-status" "$dir/.spectacular/debug/bad-invariant"
  local out_ok code_ok
  out_ok=$(cd "$dir" && "$CLI" doctor debug 2>&1) && code_ok=0 || code_ok=$?
  assert_exit_code "0" "$code_ok" "clean debug spine exits 0"
  assert_output_contains "$out_ok" "conform to schema enums" "clean spine reports pass"

  rm -rf "$dir"
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
scenario_8b_multi_area
scenario_9_json_output
scenario_10_mechanical_gitignore
scenario_11_help_flag
scenario_12_v06_scaffold_suggestion
scenario_13_migration_chain_validates
scenario_14_migration_chain_gap
scenario_15_roadmap_shape_detection
scenario_16_roadmap_icebox_rename
scenario_17_description_length
scenario_18_debug_spine_validation

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"

if [[ $fail_count -gt 0 ]]; then
  exit 1
fi
exit 0
