#!/usr/bin/env bash
# tests/cli/docs.test.sh — covers v0.6.0 public-docs surface.
#
# Scenarios:
#   1. Fresh `docs init` scaffolds docs.yaml + index + 3 sections + placeholder pages
#   2. `docs init --minimal` scaffolds only docs.yaml + index.md
#   3. Repeat `docs init` is non-destructive
#   4. `doctor docs` skips silently when no docs/ exists
#   5. `doctor docs` passes on freshly-init'd tree
#   6. `doctor docs` errors when a declared page file is missing
#   7. `doctor docs` warns on orphan files (present on disk, not in docs.yaml)
#   8. `doctor docs` errors on missing required frontmatter
#   9. `doctor docs --fix` injects frontmatter stub for files missing the delimiter
#  10. Flat-tree `extras:` resolution works (top-level pages, no section folder)
#  11. `docs new` is a skill verb — CLI refuses cleanly
#  12. `docs --help` shows usage

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
assert_file_absent() {
  if [[ ! -e "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected absent: $1"; fail_count=$((fail_count + 1)); fi
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
assert_contains() {
  if echo "$1" | grep -qF "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ output missing: $2"; fail_count=$((fail_count + 1)); fi
}

run_cli() {
  local dir="$1"; shift
  (cd "$dir" && "$CLI" "$@" 2>&1)
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_init_default() {
  echo "Scenario 1: docs init scaffolds docs.yaml + index + 3 sections + placeholder pages"
  local dir="/tmp/spectacular-docs-test-1"
  rm -rf "$dir" && mkdir -p "$dir"

  run_cli "$dir" docs init >/dev/null

  assert_file_exists "$dir/docs/docs.yaml"
  assert_file_exists "$dir/docs/index.md"
  assert_dir_exists "$dir/docs/getting-started"
  assert_dir_exists "$dir/docs/guides"
  assert_dir_exists "$dir/docs/reference"
  assert_file_exists "$dir/docs/getting-started/install.md"
  assert_file_exists "$dir/docs/getting-started/quickstart.md"
  assert_file_exists "$dir/docs/getting-started/concepts.md"
  assert_file_contains "$dir/docs/docs.yaml" "sections:"
  assert_file_contains "$dir/docs/docs.yaml" "getting-started"
  assert_file_contains "$dir/docs/getting-started/install.md" "title: Install"

  rm -rf "$dir"
}

scenario_2_init_minimal() {
  echo "Scenario 2: docs init --minimal scaffolds only docs.yaml + index.md"
  local dir="/tmp/spectacular-docs-test-2"
  rm -rf "$dir" && mkdir -p "$dir"

  run_cli "$dir" docs init --minimal >/dev/null

  assert_file_exists "$dir/docs/docs.yaml"
  assert_file_exists "$dir/docs/index.md"
  assert_file_absent "$dir/docs/getting-started"
  assert_file_absent "$dir/docs/guides"

  rm -rf "$dir"
}

scenario_3_repeat_init_nondestructive() {
  echo "Scenario 3: repeat docs init is non-destructive"
  local dir="/tmp/spectacular-docs-test-3"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null

  echo "USER CONTENT" >> "$dir/docs/index.md"
  local sentinel
  sentinel="$(md5sum "$dir/docs/index.md" 2>/dev/null | cut -d' ' -f1)"
  [[ -z "$sentinel" ]] && sentinel="$(md5 -q "$dir/docs/index.md" 2>/dev/null)"

  run_cli "$dir" docs init >/dev/null

  local after
  after="$(md5sum "$dir/docs/index.md" 2>/dev/null | cut -d' ' -f1)"
  [[ -z "$after" ]] && after="$(md5 -q "$dir/docs/index.md" 2>/dev/null)"

  if [[ "$sentinel" == "$after" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ index.md was overwritten"; fail_count=$((fail_count + 1)); fi

  rm -rf "$dir"
}

scenario_4_doctor_docs_skips_when_absent() {
  echo "Scenario 4: doctor docs skips silently when no docs/ exists"
  local dir="/tmp/spectacular-docs-test-4"
  rm -rf "$dir" && mkdir -p "$dir"

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 0 "doctor docs without docs/ exits 0"
  if echo "$output" | grep -q "^docs"; then
    echo "    ✗ expected no 'docs' findings"; fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi

  rm -rf "$dir"
}

scenario_5_doctor_clean() {
  echo "Scenario 5: doctor docs passes on freshly-init'd tree"
  local dir="/tmp/spectacular-docs-test-5"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 0 "doctor docs on clean tree"
  assert_contains "$output" "docs.yaml present"
  assert_contains "$output" "declared page present"

  rm -rf "$dir"
}

scenario_6_doctor_missing_declared() {
  echo "Scenario 6: doctor errors when a declared page is missing"
  local dir="/tmp/spectacular-docs-test-6"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null
  rm "$dir/docs/getting-started/install.md"

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 2 "doctor exits 2 on missing declared page"
  assert_contains "$output" "declared in docs.yaml but file missing"

  rm -rf "$dir"
}

scenario_7_doctor_orphan_warn() {
  echo "Scenario 7: doctor warns on orphan files"
  local dir="/tmp/spectacular-docs-test-7"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null
  cat > "$dir/docs/guides/orphan.md" <<EOF
---
title: Orphan
description: not declared
section: guides
status: draft
updated: 2026-05-23
---
# Orphan
EOF

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 1 "doctor exits 1 (warning) on orphan"
  assert_contains "$output" "orphan"

  rm -rf "$dir"
}

scenario_8_doctor_missing_frontmatter() {
  echo "Scenario 8: doctor errors on missing required frontmatter"
  local dir="/tmp/spectacular-docs-test-8"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null
  cat > "$dir/docs/getting-started/install.md" <<EOF
# Install
No frontmatter at all.
EOF

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 2 "doctor exits 2 on missing frontmatter"
  assert_contains "$output" "missing frontmatter delimiter"

  rm -rf "$dir"
}

scenario_9_doctor_fix_injects_stub() {
  echo "Scenario 9: doctor docs --fix injects frontmatter stub"
  local dir="/tmp/spectacular-docs-test-9"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null
  cat > "$dir/docs/getting-started/install.md" <<EOF
# Install
No frontmatter at all.
EOF

  run_cli "$dir" doctor docs --fix >/dev/null 2>&1 || true

  if head -1 "$dir/docs/getting-started/install.md" | grep -q '^---$'; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ frontmatter stub not injected"; fail_count=$((fail_count + 1))
  fi
  assert_file_contains "$dir/docs/getting-started/install.md" "title: install"
  assert_file_contains "$dir/docs/getting-started/install.md" "section: getting-started"
  assert_file_contains "$dir/docs/getting-started/install.md" "status: draft"

  rm -rf "$dir"
}

scenario_10_extras_flat_tree() {
  echo "Scenario 10: flat-tree extras: resolution works"
  local dir="/tmp/spectacular-docs-test-10"
  rm -rf "$dir" && mkdir -p "$dir/docs"
  cat > "$dir/docs/docs.yaml" <<EOF
site:
  name: Test
sections: []
extras:
  - changelog
EOF
  cat > "$dir/docs/changelog.md" <<EOF
---
title: Changelog
description: Project changes
section: ""
status: stable
updated: 2026-05-23
---
# Changelog
EOF

  local output exit_code
  output=$(run_cli "$dir" doctor docs)
  exit_code=$?

  assert_exit "$exit_code" 0 "doctor docs with extras"
  assert_contains "$output" "declared extra present"

  rm -rf "$dir"
}

scenario_11_skill_verbs_refused_by_cli() {
  echo "Scenario 11: CLI refuses skill verbs cleanly"
  local dir="/tmp/spectacular-docs-test-11"
  rm -rf "$dir" && mkdir -p "$dir"

  local output exit_code=0
  if output=$(run_cli "$dir" docs new install 2>&1); then exit_code=0; else exit_code=$?; fi

  if [[ $exit_code -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit (got $exit_code)"; fail_count=$((fail_count + 1)); fi
  assert_contains "$output" "skill verb"

  rm -rf "$dir"
}

scenario_12_docs_help() {
  echo "Scenario 12: docs --help shows usage"
  local output
  output=$("$CLI" docs --help 2>&1)
  assert_contains "$output" "Usage: spectacular docs"
  assert_contains "$output" "init"
  assert_contains "$output" "review"
}

# ── run all ───────────────────────────────────────────────────────────────────
echo "=== docs.test.sh ==="
scenario_1_init_default
scenario_2_init_minimal
scenario_3_repeat_init_nondestructive
scenario_4_doctor_docs_skips_when_absent
scenario_5_doctor_clean
scenario_6_doctor_missing_declared
scenario_7_doctor_orphan_warn
scenario_8_doctor_missing_frontmatter
scenario_9_doctor_fix_injects_stub
scenario_10_extras_flat_tree
scenario_11_skill_verbs_refused_by_cli
scenario_12_docs_help

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
