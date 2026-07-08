#!/usr/bin/env bash
# tests/cli/decide.test.sh — spectacular decide + summary count

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }

assert_exit()   { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_output_contains(){ echo "$1" | grep -qF -- "$2" && pass || fail "output should contain: $2"; }

scenario_decide() {
  echo "Scenario decide: decide exits 0, writes D<N>-<slug>.md, index.md, and summary counts correctly"
  local dir="/tmp/spectacular-decide-okf"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular"
  printf 'project:\n  name: okf\n' > "$dir/.spectacular/config.yaml"

  # Run decide (no open session) -> must exit 0 and bootstrap decisions/
  local code; (cd "$dir" && "$CLI" decide "First decision" --consequences "enable x" >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "decide exits 0"

  (cd "$dir" && "$CLI" decide "Second decision" --consequences "enable y" >/dev/null 2>&1)

  # Check that D1/D2 slug-prefixed files are written
  [[ -f "$dir/.spectacular/decisions/D1-first-decision.md" && -f "$dir/.spectacular/decisions/D2-second-decision.md" ]] && pass || fail "D1/D2 files written"

  # Check that index.md is created
  [[ -f "$dir/.spectacular/decisions/index.md" ]] && pass || fail "decisions/index.md written"

  # Summary must count the two decisions
  local out; out=$(cd "$dir" && "$CLI" summary 2>&1)
  assert_output_contains "$out" "Decisions:  2"
  rm -rf "$dir"
}

scenario_decide

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
