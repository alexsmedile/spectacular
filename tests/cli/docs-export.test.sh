#!/usr/bin/env bash
# tests/cli/docs-export.test.sh — covers v1.1.0 renderer adapters.
#
# Scenarios:
#   1. `docs export` with no renderer arg errors cleanly
#   2. `docs export mkdocs` writes mkdocs.yml + .github/workflows/docs.yml
#   3. `docs export docusaurus` writes docusaurus.config.js + sidebars.js + workflow
#   4. `--no-workflow` skips the workflow file
#   5. Idempotency: second run skips both targets (exit 0)
#   6. `--force` overwrites existing targets
#   7. `// spectacular: do-not-overwrite` pin survives `--force`
#   8. Unknown renderer (mintlify) returns actionable error
#   9. Missing docs.yaml fails clean
#  10. mkdocs nav drops empty sections (sections with no declared pages)
#  11. docusaurus sidebars drops empty sections (parity with mkdocs)
#  12. renderers: block in docs.yaml is consumed (theme, primary, organizationName)
#  13. doctor docs: valid renderers: block produces pass entries
#  14. doctor docs: unknown renderer key produces warning
#  15. doctor docs: renderers: as scalar produces error
#  16. doctor docs: renderers: as list produces error

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

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
assert_file_contains() {
  if [[ -f "$1" ]] && grep -qF "$2" "$1"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count + 1)); fi
}
assert_file_not_contains() {
  if [[ -f "$1" ]] && ! grep -qF "$2" "$1"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected '$1' to NOT contain '$2'"; fail_count=$((fail_count + 1)); fi
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
assert_not_contains() {
  if ! echo "$1" | grep -qF "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ output unexpectedly contains: $2"; fail_count=$((fail_count + 1)); fi
}

run_cli() {
  local dir="$1"; shift
  (cd "$dir" && "$CLI" "$@" 2>&1)
}

# Common setup: empty dir → docs init → ready for export.
setup_docs() {
  local dir="$1"
  rm -rf "$dir" && mkdir -p "$dir"
  run_cli "$dir" docs init >/dev/null
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_no_renderer() {
  echo "Scenario 1: docs export with no renderer errors cleanly"
  local dir="/tmp/spectacular-export-test-1"
  setup_docs "$dir"

  local output exit_code=0
  if output=$(run_cli "$dir" docs export 2>&1); then exit_code=0; else exit_code=$?; fi

  if [[ $exit_code -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit"; fail_count=$((fail_count + 1)); fi
  assert_contains "$output" "renderer"

  rm -rf "$dir"
}

scenario_2_mkdocs() {
  echo "Scenario 2: docs export mkdocs writes mkdocs.yml + workflow"
  local dir="/tmp/spectacular-export-test-2"
  setup_docs "$dir"

  run_cli "$dir" docs export mkdocs >/dev/null

  assert_file_exists "$dir/mkdocs.yml"
  assert_file_exists "$dir/.github/workflows/docs.yml"
  assert_file_contains "$dir/mkdocs.yml" "site_name:"
  assert_file_contains "$dir/mkdocs.yml" "name: material"
  assert_file_contains "$dir/mkdocs.yml" "nav:"
  assert_file_contains "$dir/.github/workflows/docs.yml" "MkDocs"
  assert_file_contains "$dir/.github/workflows/docs.yml" "mkdocs build"

  rm -rf "$dir"
}

scenario_3_docusaurus() {
  echo "Scenario 3: docs export docusaurus writes config + sidebars + workflow"
  local dir="/tmp/spectacular-export-test-3"
  setup_docs "$dir"

  run_cli "$dir" docs export docusaurus >/dev/null

  assert_file_exists "$dir/docusaurus.config.js"
  assert_file_exists "$dir/sidebars.js"
  assert_file_exists "$dir/.github/workflows/docs.yml"
  assert_file_contains "$dir/docusaurus.config.js" "module.exports"
  assert_file_contains "$dir/docusaurus.config.js" "presets:"
  assert_file_contains "$dir/sidebars.js" "type: 'category'"
  assert_file_contains "$dir/.github/workflows/docs.yml" "Docusaurus"
  assert_file_contains "$dir/.github/workflows/docs.yml" "docusaurus build"

  rm -rf "$dir"
}

scenario_4_no_workflow() {
  echo "Scenario 4: --no-workflow skips deploy file"
  local dir="/tmp/spectacular-export-test-4"
  setup_docs "$dir"

  run_cli "$dir" docs export mkdocs --no-workflow >/dev/null

  assert_file_exists "$dir/mkdocs.yml"
  assert_file_absent "$dir/.github/workflows/docs.yml"

  rm -rf "$dir"
}

scenario_5_idempotency() {
  echo "Scenario 5: second run skips with clear report (exit 0)"
  local dir="/tmp/spectacular-export-test-5"
  setup_docs "$dir"

  run_cli "$dir" docs export mkdocs >/dev/null
  local output exit_code=0
  if output=$(run_cli "$dir" docs export mkdocs 2>&1); then exit_code=0; else exit_code=$?; fi

  assert_exit "$exit_code" 0 "second run exits 0"
  assert_contains "$output" "skipped"
  assert_contains "$output" "use --force"

  rm -rf "$dir"
}

scenario_6_force() {
  echo "Scenario 6: --force overwrites existing files"
  local dir="/tmp/spectacular-export-test-6"
  setup_docs "$dir"

  run_cli "$dir" docs export mkdocs >/dev/null
  echo "# USER_INJECTED_MARKER_42" > "$dir/mkdocs.yml"

  run_cli "$dir" docs export mkdocs --force >/dev/null

  assert_file_contains "$dir/mkdocs.yml" "Generated by spectacular"
  assert_file_not_contains "$dir/mkdocs.yml" "USER_INJECTED_MARKER_42"

  rm -rf "$dir"
}

scenario_7_pin_respected() {
  echo "Scenario 7: do-not-overwrite pin survives --force"
  local dir="/tmp/spectacular-export-test-7"
  setup_docs "$dir"

  printf '# spectacular: do-not-overwrite\nsite_name: PINNED\n' > "$dir/mkdocs.yml"
  local output
  output=$(run_cli "$dir" docs export mkdocs --force 2>&1)

  assert_contains "$output" "pinned"
  assert_file_contains "$dir/mkdocs.yml" "PINNED"
  assert_file_not_contains "$dir/mkdocs.yml" "Generated by spectacular"

  rm -rf "$dir"
}

scenario_8_unknown_renderer() {
  echo "Scenario 8: unknown renderer errors with adapter-ref pointer"
  local dir="/tmp/spectacular-export-test-8"
  setup_docs "$dir"

  local output exit_code=0
  if output=$(run_cli "$dir" docs export mintlify 2>&1); then exit_code=0; else exit_code=$?; fi

  if [[ $exit_code -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit"; fail_count=$((fail_count + 1)); fi
  assert_contains "$output" "not shipped"
  assert_contains "$output" "docs-renderer-adapters"

  rm -rf "$dir"
}

scenario_9_missing_manifest() {
  echo "Scenario 9: missing docs.yaml fails clean"
  local dir="/tmp/spectacular-export-test-9"
  rm -rf "$dir" && mkdir -p "$dir"

  local output exit_code=0
  if output=$(run_cli "$dir" docs export mkdocs 2>&1); then exit_code=0; else exit_code=$?; fi

  if [[ $exit_code -ne 0 ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected non-zero exit"; fail_count=$((fail_count + 1)); fi
  assert_contains "$output" "docs/docs.yaml not found"
  assert_contains "$output" "docs init"

  rm -rf "$dir"
}

scenario_10_mkdocs_drops_empty_sections() {
  echo "Scenario 10: mkdocs nav drops empty sections"
  local dir="/tmp/spectacular-export-test-10"
  setup_docs "$dir"

  run_cli "$dir" docs export mkdocs --no-workflow >/dev/null

  # 'Getting Started' has 3 pages → must appear in nav
  assert_file_contains "$dir/mkdocs.yml" "Getting Started:"
  # 'Guides' and 'Reference' have empty pages: [] → must NOT appear
  assert_file_not_contains "$dir/mkdocs.yml" "Guides:"
  assert_file_not_contains "$dir/mkdocs.yml" "Reference:"

  rm -rf "$dir"
}

scenario_11_docusaurus_drops_empty_sections() {
  echo "Scenario 11: docusaurus sidebars drops empty sections"
  local dir="/tmp/spectacular-export-test-11"
  setup_docs "$dir"

  run_cli "$dir" docs export docusaurus --no-workflow >/dev/null

  assert_file_contains "$dir/sidebars.js" "label: 'Getting Started'"
  assert_file_not_contains "$dir/sidebars.js" "label: 'Guides'"
  assert_file_not_contains "$dir/sidebars.js" "label: 'Reference'"

  rm -rf "$dir"
}

scenario_12_renderers_block_consumed() {
  echo "Scenario 12: renderers: block in docs.yaml is consumed by adapters"
  local dir="/tmp/spectacular-export-test-12"
  setup_docs "$dir"

  cat >> "$dir/docs/docs.yaml" <<'YAML'

renderers:
  mkdocs:
    theme: material
    primary: indigo
    scheme: slate
    repo_url: https://github.com/example/repo
    edit_uri: edit/main/docs/
  docusaurus:
    preset: classic
    organizationName: example
    projectName: repo
YAML

  run_cli "$dir" docs export mkdocs --no-workflow >/dev/null
  run_cli "$dir" docs export docusaurus --no-workflow >/dev/null

  assert_file_contains "$dir/mkdocs.yml" "primary: indigo"
  assert_file_contains "$dir/mkdocs.yml" "scheme: slate"
  assert_file_contains "$dir/mkdocs.yml" "repo_url: https://github.com/example/repo"
  assert_file_contains "$dir/mkdocs.yml" "edit_uri: edit/main/docs/"
  assert_file_contains "$dir/docusaurus.config.js" "organizationName: 'example'"
  assert_file_contains "$dir/docusaurus.config.js" "projectName: 'repo'"

  rm -rf "$dir"
}

# Scenarios 13-16 (doctor renderers: validation) moved to pageworks's
# tests/cli/doctor.test.sh in v1.2.0. Spectacular's doctor docs is now
# discovery-only; renderer-block validation lives in pageworks doctor.

# ── run all ───────────────────────────────────────────────────────────────────
echo "=== docs-export.test.sh ==="
scenario_1_no_renderer
scenario_2_mkdocs
scenario_3_docusaurus
scenario_4_no_workflow
scenario_5_idempotency
scenario_6_force
scenario_7_pin_respected
scenario_8_unknown_renderer
scenario_9_missing_manifest
scenario_10_mkdocs_drops_empty_sections
scenario_11_docusaurus_drops_empty_sections
scenario_12_renderers_block_consumed

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
