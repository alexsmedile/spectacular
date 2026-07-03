#!/usr/bin/env bash
# tests/cli/decide.test.sh — spectacular decide + summary count, both modes
#
# Regression guards for two bugs found 2026-07-03:
#   Bug 1  flat-mode decisions counted 0 in `summary` (only index-mode counted)
#   Bug 2  flat-mode `decide` returned exit 1 despite a successful append
#          (trailing `[[ -n "$session" ]] &&` short-circuits with no open session)

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }

assert_exit()   { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_output_contains(){ echo "$1" | grep -qF -- "$2" && pass || fail "output should contain: $2"; }

# ── flat mode: no decisions/ folder; ADRs are prose blocks in DECISIONS.md ────
scenario_flat() {
  echo "Scenario flat: decide exits 0 + summary counts prose blocks"
  local dir="/tmp/spectacular-decide-flat"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular"
  printf 'name: flat\n' > "$dir/.spectacular/config.yaml"

  # Bug 2: no open session → must still exit 0
  local code; (cd "$dir" && "$CLI" decide "First flat call" --consequences x >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "flat decide (no session) exits 0"

  (cd "$dir" && "$CLI" decide "Second flat call" --consequences y >/dev/null 2>&1)

  # Bug 1: summary must count the two prose blocks, not 0
  local out; out=$(cd "$dir" && "$CLI" summary 2>&1)
  assert_output_contains "$out" "Decisions:  2"
  rm -rf "$dir"
}

# ── index mode: decisions/ folder; ADRs are D<N>.md files ─────────────────────
scenario_index() {
  echo "Scenario index: per-file ADRs counted; unaffected by the flat fix"
  local dir="/tmp/spectacular-decide-index"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular/decisions"
  printf 'name: idx\n' > "$dir/.spectacular/config.yaml"
  printf -- '---\nmode: index\n---\n# Decisions\n' > "$dir/.spectacular/DECISIONS.md"

  (cd "$dir" && "$CLI" decide "Index call one" --consequences a >/dev/null 2>&1)
  (cd "$dir" && "$CLI" decide "Index call two" --consequences b >/dev/null 2>&1)

  local out; out=$(cd "$dir" && "$CLI" summary 2>&1)
  assert_output_contains "$out" "Decisions:  2"
  [[ -f "$dir/.spectacular/decisions/D1.md" && -f "$dir/.spectacular/decisions/D2.md" ]] && pass || fail "D1/D2 files written"
  rm -rf "$dir"
}

scenario_flat
scenario_index

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
