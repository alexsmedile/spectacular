#!/usr/bin/env bash
# Static contract checks for approved execution continuing through checkpoints.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
SKILL="$REPO_ROOT/skills/spectacular/SKILL.md"
REQUEST="$REPO_ROOT/skills/spectacular/references/request-workflow.md"
ACTIVE="$REPO_ROOT/skills/spectacular/references/active-request.md"
BUILD="$REPO_ROOT/skills/spectacular/references/build-workflow.md"
AGENTS="$REPO_ROOT/.spectacular/AGENTS.md"
CLI="$REPO_ROOT/cli/spectacular"
BEHAVIOR="$REPO_ROOT/tests/agents/continuation/README.md"

pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }
assert_lacks() { ! grep -qF -- "$2" "$1" && pass || fail "$1 should not contain: $2"; }

echo "Scenario execution has explicit terminal states"
assert_contains "$SKILL" "Execution-turn state machine"
assert_contains "$SKILL" "RUNNING"
assert_contains "$SKILL" "BLOCKED"
assert_contains "$SKILL" "COMPLETE"
assert_contains "$SKILL" "response is forbidden."
assert_contains "$SKILL" "Never end a turn with"
assert_contains "$SKILL" "action instead."

echo "Scenario runtime plan and progress channel are portable"
assert_contains "$SKILL" 'Codex `update_plan` or Claude `TaskCreate`/`TaskUpdate`'
assert_contains "$SKILL" "intermediate/progress channel"
assert_contains "$AGENTS" '`update_plan` in Codex; `TaskCreate` / `TaskUpdate` in Claude'
assert_contains "$AGENTS" 'exactly one `in_progress`'

echo "Scenario failed checks remain execution work"
assert_contains "$SKILL" "A red check is work to diagnose and repair"
assert_contains "$ACTIVE" "A failed task check remains part of execution"
assert_contains "$BUILD" "A failure starts an in-scope diagnose"

echo "Scenario generated briefs carry the immediate completion contract"
assert_contains "$REQUEST" "execution-turn completion contract"
assert_contains "$CLI" "## Completion contract"
assert_contains "$CLI" "Remain RUNNING while any plan item or completion condition remains"
assert_contains "$CLI" "A failed check starts in-scope diagnosis and repair"

echo "Scenario adversarial behavior eval covers premature finals and red checks"
assert_contains "$BEHAVIOR" "No terminal response occurs while a plan item is pending or in progress."
assert_contains "$BEHAVIOR" "The first red check does not produce a terminal response."
assert_contains "$BEHAVIOR" "Automatic failure phrases in a terminal response"
assert_contains "$BEHAVIOR" '`COMPLETE` or `BLOCKED`'

echo "Scenario obsolete checkpoint handoff instruction stays removed"
assert_lacks "$SKILL" "Identify the single highest-priority next action and ask what the user wants to do."

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
