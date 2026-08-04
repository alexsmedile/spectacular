#!/usr/bin/env bash
# tests/cli/traffic.test.sh — local durable-evidence traffic preflight.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
fail_count=0
pass_count=0

pass() { echo "    ✓ $1"; pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { echo "$1" | grep -qF -- "$2" && pass "$3" || fail "$3 — missing: $2"; }

mk_workspace() {
  local dir
  dir=$(mktemp -d)
  mkdir -p "$dir/.spectacular/requests"
  printf 'project:\n  name: traffic-test\n  owner: test\n' > "$dir/.spectacular/config.yaml"
  echo "$dir"
}

mk_plan() {
  local dir="$1" slug="$2" extra="${3:-}"
  mkdir -p "$dir/.spectacular/requests/$slug"
  printf -- '---\nstatus: active\npriority: medium\nowner: test\nupdated: 2026-08-04\nsummary: "%s fixture"\nrelated: []\n%s---\n# %s\n' "$slug" "$extra" "$slug" > "$dir/.spectacular/requests/$slug/PLAN.md"
}

echo "── parallel: complete, disjoint named boundaries ──"
W=$(mk_workspace)
mk_plan "$W" "api" $'traffic-boundaries:\n  - api:routes\n'
mk_plan "$W" "docs" $'traffic-boundaries:\n  - docs:commands\n'
out=$(cd "$W" && bash "$CLI" traffic preflight api --against docs 2>&1)
assert_contains "$out" "parallel" "disjoint complete boundaries are parallel"
assert_contains "$out" "assessed " "assessment is explicitly time-bound"
rm -rf "$W"

echo "── conditional: shared named boundary ──"
W=$(mk_workspace)
mk_plan "$W" "cli-a" $'traffic-boundaries:\n  - cli:dispatch\n'
mk_plan "$W" "cli-b" $'traffic-boundaries:\n  - cli:dispatch\n'
out=$(cd "$W" && bash "$CLI" traffic preflight cli-a --against cli-b 2>&1)
assert_contains "$out" "conditional" "shared boundary is conditional"
assert_contains "$out" "shared boundary: cli:dispatch" "conditional reason names durable boundary"
rm -rf "$W"

echo "── serialized: dependency and shared release constraint ──"
W=$(mk_workspace)
mk_plan "$W" "consumer" $'depends-on:\n  - migration\n'
mk_plan "$W" "migration"
out=$(cd "$W" && bash "$CLI" traffic preflight consumer --against migration 2>&1)
assert_contains "$out" "serialized" "declared dependency serializes"
rm -rf "$W"
W=$(mk_workspace)
mk_plan "$W" "release-a" $'release-constraints:\n  - release:v2\n'
mk_plan "$W" "release-b" $'release-constraints:\n  - release:v2\n'
out=$(cd "$W" && bash "$CLI" traffic preflight release-a --against release-b 2>&1)
assert_contains "$out" "serialized" "shared release constraint serializes"
assert_contains "$out" "shared release constraint: release:v2" "serialized reason names release constraint"
rm -rf "$W"

echo "── unknown: insufficient durable evidence ──"
W=$(mk_workspace)
mk_plan "$W" "one"
mk_plan "$W" "two"
out=$(cd "$W" && bash "$CLI" traffic preflight one --against two 2>&1)
assert_contains "$out" "unknown" "missing traffic declarations are unknown"
assert_contains "$out" "insufficient durable traffic evidence" "unknown explains insufficient evidence"
json=$(cd "$W" && bash "$CLI" traffic preflight one --against two --json 2>&1)
assert_contains "$json" '"state":"unknown"' "JSON preserves unknown state"
rm -rf "$W"

echo "── links: conflicts-with is durable graph evidence ──"
W=$(mk_workspace)
mk_plan "$W" "writer" $'conflicts-with:\n  - reader\n'
mk_plan "$W" "reader"
out=$(cd "$W" && bash "$CLI" traffic preflight writer --against reader 2>&1)
assert_contains "$out" "serialized" "declared conflict serializes"
links=$(cd "$W" && bash "$CLI" links writer 2>&1)
assert_contains "$links" "conflicts-with: reader" "links renders conflict relationship"
doctor=$(cd "$W" && bash "$CLI" doctor links 2>&1)
assert_contains "$doctor" "all related:/depends-on:/blocks:/conflicts-with: targets resolve" "doctor validates conflict target"
rm -rf "$W"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
