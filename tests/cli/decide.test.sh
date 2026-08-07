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
  echo "Scenario decide: decide exits 0, writes UUIDv7 slug files, index.md, and summary counts correctly"
  local dir="/tmp/spectacular-decide-okf"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular"
  printf 'project:\n  name: okf\n' > "$dir/.spectacular/config.yaml"

  # Run decide (no open session) -> must exit 0 and bootstrap decisions/
  local code; (cd "$dir" && "$CLI" decide "First decision" --consequences "enable x #cli #docs" --tags cli,docs >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "decide exits 0"

  (cd "$dir" && "$CLI" decide "Second decision" --consequences "enable y" >/dev/null 2>&1)

  # Check UUIDv7 files and index IDs.
  [[ -f "$dir/.spectacular/decisions/first-decision.md" && -f "$dir/.spectacular/decisions/second-decision.md" ]] && pass || fail "slug files written"
  grep -Eq -- '\*\*[0-9a-f]{8}-[0-9a-f]{4}-7[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}\*\*' "$dir/.spectacular/decisions/index.md" && pass || fail "UUIDv7 decision ID indexed"
  [[ "$(grep -o '#cli' "$dir/.spectacular/decisions/index.md" | wc -l | tr -d ' ')" == 1 ]] && pass || fail "decision index emits tags once"

  # Legacy D<N> headings retain a clean title when a later decision regenerates the index.
  printf '%s\n' '# D3 — Legacy title' '' '**Consequences:**' 'legacy rationale' > "$dir/.spectacular/decisions/D3-legacy-title.md"
  (cd "$dir" && "$CLI" decide "Fourth decision" >/dev/null 2>&1)
  grep -qF -- 'Legacy title — legacy rationale' "$dir/.spectacular/decisions/index.md" && pass || fail "legacy D heading is indexed without duplicated prefix"

  # Check that index.md is created
  [[ -f "$dir/.spectacular/decisions/index.md" ]] && pass || fail "decisions/index.md written"

  # Summary must count all decision files.
  local out; out=$(cd "$dir" && "$CLI" summary 2>&1)
  assert_output_contains "$out" "Decisions:  4"
  rm -rf "$dir"
}

scenario_decide

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
