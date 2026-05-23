#!/usr/bin/env bash
# tests/cli/docs.test.sh — spectacular's docs surface in v1.2.0+ (DEPRECATED state)
#
# Public-facing docs work moved to the pageworks skill in v1.2.0. Spectacular's
# docs verbs still function for backward compatibility but emit a deprecation
# banner. `doctor docs` is slimmed to discovery-only. Removal lands in v2.0.0.
#
# Scenarios:
#   1. `spectacular docs init` still scaffolds docs.yaml + index.md (backward compat)
#   2. `spectacular docs init` emits deprecation banner on stderr
#   3. `doctor docs` skips silently when no docs/ folder exists
#   4. `doctor docs` reports docs/ folder + manifest presence (discovery-only)
#   5. `doctor docs` info-level pageworks install hint when pageworks not in PATH
#   6. `doctor docs` does NOT validate frontmatter / orphans / renderer-block (moved to pageworks)
#   7. `spectacular docs new` / review / status emit deprecation + skill-verb message
#   8. `spectacular docs --help` shows DEPRECATED banner
#   9. `spectacular docs export` still works (covered fully in docs-export.test.sh — this is a smoke check)

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_file_exists() {
  if [[ -f "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected file: $1"; fail_count=$((fail_count + 1)); fi
}
assert_contains() {
  if echo "$1" | grep -qF "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ output missing: $2"; fail_count=$((fail_count + 1)); fi
}
assert_not_contains() {
  if ! echo "$1" | grep -qF "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ output unexpectedly contains: $2"; fail_count=$((fail_count + 1)); fi
}
assert_exit() {
  local got="$1" want="$2" label="$3"
  if [[ "$got" -eq "$want" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ $label: exit $got want $want"; fail_count=$((fail_count + 1)); fi
}

run_cli() {
  local dir="$1"; shift
  (cd "$dir" && "$CLI" "$@" 2>&1)
}

scenario_1_init_still_works() {
  echo "Scenario 1: 'spectacular docs init' still scaffolds (backward compat)"
  local dir="/tmp/spectacular-docs-test-1"
  rm -rf "$dir" && mkdir -p "$dir"

  run_cli "$dir" docs init >/dev/null

  assert_file_exists "$dir/docs/docs.yaml"
  assert_file_exists "$dir/docs/index.md"
  assert_file_exists "$dir/docs/getting-started/install.md"

  rm -rf "$dir"
}

scenario_2_init_deprecation_banner() {
  echo "Scenario 2: 'spectacular docs init' emits deprecation banner"
  local dir="/tmp/spectacular-docs-test-2"
  rm -rf "$dir" && mkdir -p "$dir"

  local output
  output=$(run_cli "$dir" docs init 2>&1)

  assert_contains "$output" "deprecated in v1.2.0"
  assert_contains "$output" "pageworks"

  rm -rf "$dir"
}

scenario_3_doctor_docs_skips_when_absent() {
  echo "Scenario 3: doctor docs skips silently when no docs/ folder"
  local dir="/tmp/spectacular-docs-test-3"
  rm -rf "$dir" && mkdir -p "$dir"

  local output exit_code=0
  if output=$(run_cli "$dir" doctor docs 2>&1); then exit_code=0; else exit_code=$?; fi

  assert_exit "$exit_code" 0 "doctor docs (no docs/) exits 0"
  assert_not_contains "$output" "docs/ folder present"

  rm -rf "$dir"
}

scenario_4_doctor_docs_discovery_only() {
  echo "Scenario 4: doctor docs reports folder + manifest presence (discovery-only)"
  local dir="/tmp/spectacular-docs-test-4"
  rm -rf "$dir" && mkdir -p "$dir/docs"
  touch "$dir/docs/docs.yaml"

  local output
  output=$(run_cli "$dir" doctor docs 2>&1)

  assert_contains "$output" "docs/ folder present"
  assert_contains "$output" "docs.yaml manifest present"

  rm -rf "$dir"
}

scenario_5_doctor_pageworks_hint() {
  echo "Scenario 5: doctor docs info hint when pageworks not installed"
  local dir="/tmp/spectacular-docs-test-5"
  rm -rf "$dir" && mkdir -p "$dir/docs"
  touch "$dir/docs/docs.yaml"

  # Only assert the hint if pageworks isn't actually in PATH (CI/dev env varies)
  if ! command -v pageworks >/dev/null 2>&1; then
    local output
    output=$(run_cli "$dir" doctor docs 2>&1)
    assert_contains "$output" "pageworks not installed"
    assert_contains "$output" "github.com/alexsmedile/pageworks"
  else
    echo "    ⊘ pageworks installed in PATH — skipping hint assertion"
    pass_count=$((pass_count + 1))
  fi

  rm -rf "$dir"
}

scenario_6_doctor_no_deep_validation() {
  echo "Scenario 6: doctor docs does NOT validate frontmatter/orphans/renderers"
  local dir="/tmp/spectacular-docs-test-6"
  rm -rf "$dir" && mkdir -p "$dir/docs/getting-started"
  cat > "$dir/docs/docs.yaml" <<'YAML'
site:
  name: deep-test
sections:
  - id: getting-started
    title: Getting Started
    order: 1
    pages: [install]
renderers:
  bogus-renderer:
    foo: bar
YAML
  # install.md missing on disk + bogus renderer key — neither should be flagged
  # by spectacular's slim doctor (both belong to pageworks now).
  cat > "$dir/docs/getting-started/orphan.md" <<'EOF'
no-frontmatter-here
EOF

  local output
  output=$(run_cli "$dir" doctor docs 2>&1)

  assert_not_contains "$output" "missing required frontmatter"
  assert_not_contains "$output" "orphan"
  assert_not_contains "$output" "not a recognized renderer"
  assert_not_contains "$output" "declared in docs.yaml but file missing"

  rm -rf "$dir"
}

scenario_7_skill_verbs_deprecation() {
  echo "Scenario 7: skill verbs (new|review|status) emit deprecation + skill-verb message"
  local dir="/tmp/spectacular-docs-test-7"
  rm -rf "$dir" && mkdir -p "$dir"

  for verb in new review status; do
    local output exit_code=0
    if output=$(run_cli "$dir" docs "$verb" 2>&1); then exit_code=0; else exit_code=$?; fi
    if [[ $exit_code -ne 0 ]]; then pass_count=$((pass_count + 1))
    else echo "    ✗ docs $verb: expected non-zero exit (got $exit_code)"; fail_count=$((fail_count + 1)); fi
    assert_contains "$output" "deprecated in v1.2.0"
    assert_contains "$output" "skill verb"
  done

  rm -rf "$dir"
}

scenario_8_docs_help_deprecation() {
  echo "Scenario 8: 'spectacular docs --help' shows DEPRECATED banner"
  local output
  output=$("$CLI" docs --help 2>&1)
  assert_contains "$output" "DEPRECATED"
  assert_contains "$output" "pageworks"
}

scenario_9_export_still_works_smoke() {
  echo "Scenario 9: 'spectacular docs export mkdocs' still works (smoke check)"
  local dir="/tmp/spectacular-docs-test-9"
  rm -rf "$dir" && mkdir -p "$dir"

  run_cli "$dir" docs init >/dev/null
  run_cli "$dir" docs export mkdocs --no-workflow >/dev/null

  assert_file_exists "$dir/mkdocs.yml"

  rm -rf "$dir"
}

echo "=== docs.test.sh (v1.2.0 deprecation state) ==="
scenario_1_init_still_works
scenario_2_init_deprecation_banner
scenario_3_doctor_docs_skips_when_absent
scenario_4_doctor_docs_discovery_only
scenario_5_doctor_pageworks_hint
scenario_6_doctor_no_deep_validation
scenario_7_skill_verbs_deprecation
scenario_8_docs_help_deprecation
scenario_9_export_still_works_smoke

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
