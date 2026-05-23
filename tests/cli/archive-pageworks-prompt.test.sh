#!/usr/bin/env bash
# tests/cli/archive-pageworks-prompt.test.sh
#
# Verifies the v1.2.0 archive-time pageworks-handoff prompt:
#   1. Prompt fires when docs/ exists AND the archived PLAN references SPEC.md/specs/
#   2. Prompt suppressed by --no-docs-prompt flag
#   3. Prompt suppressed by 'docs.prompt_on_archive: false' in .spectacular/config.yaml
#   4. Prompt silent when no docs/ folder exists
#   5. Prompt silent when archived request doesn't reference SPEC/specs/ARCHITECTURE/PRD

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

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

setup_workspace() {
  local dir="$1"
  rm -rf "$dir" && mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --name archive-prompt-test >/dev/null 2>&1)
}

setup_request_with_spec_ref() {
  local dir="$1" slug="$2"
  mkdir -p "$dir/.spectacular/requests/$slug"
  cat > "$dir/.spectacular/requests/$slug/PLAN.md" <<EOF
---
status: verified
priority: medium
owner: test
updated: 2026-05-23
summary: "Touches SPEC.md and specs/ — should trigger pageworks prompt"
---

# Plan
References SPEC.md and specs/foo/ for the capability tracking.
EOF
}

setup_request_no_spec_ref() {
  local dir="$1" slug="$2"
  mkdir -p "$dir/.spectacular/requests/$slug"
  cat > "$dir/.spectacular/requests/$slug/PLAN.md" <<EOF
---
status: verified
priority: medium
owner: test
updated: 2026-05-23
summary: "No spec references at all"
---

# Plan
This plan is purely about internal refactor; no public surface changes.
EOF
}

scenario_1_prompt_fires() {
  echo "Scenario 1: prompt fires when docs/ exists + PLAN references SPEC"
  local dir="/tmp/spectacular-archive-prompt-1"
  setup_workspace "$dir"
  mkdir "$dir/docs"
  touch "$dir/docs/docs.yaml"
  setup_request_with_spec_ref "$dir" "spec-touching"

  local output
  output=$(run_cli "$dir" archive spec-touching --skip-doctor)

  assert_contains "$output" "Public docs/ may need updates"
  assert_contains "$output" "pageworks"
}

scenario_2_no_docs_prompt_flag() {
  echo "Scenario 2: --no-docs-prompt suppresses the prompt"
  local dir="/tmp/spectacular-archive-prompt-2"
  setup_workspace "$dir"
  mkdir "$dir/docs"
  touch "$dir/docs/docs.yaml"
  setup_request_with_spec_ref "$dir" "spec-touching-2"

  local output
  output=$(run_cli "$dir" archive spec-touching-2 --skip-doctor --no-docs-prompt)

  assert_not_contains "$output" "Public docs/ may need updates"
}

scenario_3_config_suppression() {
  echo "Scenario 3: 'docs.prompt_on_archive: false' in config suppresses prompt"
  local dir="/tmp/spectacular-archive-prompt-3"
  setup_workspace "$dir"
  mkdir "$dir/docs"
  touch "$dir/docs/docs.yaml"
  cat >> "$dir/.spectacular/config.yaml" <<'YAML'

docs:
  prompt_on_archive: false
YAML
  setup_request_with_spec_ref "$dir" "spec-touching-3"

  local output
  output=$(run_cli "$dir" archive spec-touching-3 --skip-doctor)

  assert_not_contains "$output" "Public docs/ may need updates"
}

scenario_4_silent_when_no_docs() {
  echo "Scenario 4: prompt silent when no docs/ folder"
  local dir="/tmp/spectacular-archive-prompt-4"
  setup_workspace "$dir"
  # No docs/ folder
  setup_request_with_spec_ref "$dir" "spec-touching-4"

  local output
  output=$(run_cli "$dir" archive spec-touching-4 --skip-doctor)

  assert_not_contains "$output" "Public docs/ may need updates"
}

scenario_5_silent_when_no_spec_reference() {
  echo "Scenario 5: prompt silent when PLAN doesn't reference SPEC/specs/ARCHITECTURE/PRD"
  local dir="/tmp/spectacular-archive-prompt-5"
  setup_workspace "$dir"
  mkdir "$dir/docs"
  touch "$dir/docs/docs.yaml"
  setup_request_no_spec_ref "$dir" "no-spec-ref"

  local output
  output=$(run_cli "$dir" archive no-spec-ref --skip-doctor)

  assert_not_contains "$output" "Public docs/ may need updates"
}

echo "=== archive-pageworks-prompt.test.sh ==="
scenario_1_prompt_fires
scenario_2_no_docs_prompt_flag
scenario_3_config_suppression
scenario_4_silent_when_no_docs
scenario_5_silent_when_no_spec_reference

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
