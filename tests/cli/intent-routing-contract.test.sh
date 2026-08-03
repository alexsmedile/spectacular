#!/usr/bin/env bash
# Static contract checks for the agentic intent-routing guard.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
SKILL="$REPO_ROOT/skills/spectacular/SKILL.md"
ROUTING="$REPO_ROOT/skills/spectacular/references/intent-routing.md"
REQUEST="$REPO_ROOT/skills/spectacular/references/new-request.md"
LIFECYCLE="$REPO_ROOT/skills/spectacular/references/spec-lifecycle.md"

pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }

echo "Scenario intent routing is explicit before a new SPC"
assert_contains "$SKILL" "First decision — does this need Spectacular?"
assert_contains "$SKILL" "Direct PR-shaped change"
assert_contains "$SKILL" "A Codex/harness plan is an ephemeral"
assert_contains "$SKILL" "“Implement this plan” means execute"
assert_contains "$SKILL" "“draft a spec” request, the intent receipt is mandatory"
assert_contains "$SKILL" "Never infer a new SPC from nearby docs impact or repository context."

echo "Scenario receipt preserves user intent over repository context"
assert_contains "$ROUTING" "Do not let repository context choose the work for the user."
assert_contains "$ROUTING" "## Intent receipt — required before a natural-language SPC draft"
assert_contains "$ROUTING" "Not doing: <closest plausible alternative>"
assert_contains "$ROUTING" "Draft this SPC? (yes / correct it)"
assert_contains "$ROUTING" "not permission to create"
assert_contains "$ROUTING" "## What each layer is for"
assert_contains "$ROUTING" "direct change +"
assert_contains "$ROUTING" "no need for a formal request"

echo "Scenario request and lifecycle flows share the guard"
assert_contains "$REQUEST" "First run [[intent-routing]]."
assert_contains "$LIFECYCLE" "[[intent-routing]] and show the user an intent receipt"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
