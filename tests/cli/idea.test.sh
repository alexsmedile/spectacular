#!/usr/bin/env bash
# tests/cli/idea.test.sh — destination-aware local IDEA handoffs

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0
fail_count=0

pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_file() { [[ -f "$1" ]] && pass || fail "expected file: $1"; }
assert_absent() { [[ ! -e "$1" ]] && pass || fail "expected absent: $1"; }
assert_contains() { [[ -f "$1" ]] && grep -qF -- "$2" "$1" && pass || fail "expected '$1' to contain: $2"; }
assert_output() { printf '%s' "$1" | grep -qF -- "$2" && pass || fail "expected output: $2"; }

new_ws() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --kit coding --with roadmap --name idea-test >/dev/null)
}

scenario_request_and_legacy() {
  echo "Scenario request destination: explicit and legacy compatibility"
  local d="/tmp/spectacular-idea-request" out
  new_ws "$d"
  (cd "$d" && "$CLI" idea new explicit >/dev/null)
  (cd "$d" && "$CLI" idea promote explicit --to request >/dev/null)
  assert_file "$d/.spectacular/requests/explicit/PLAN.md"
  assert_file "$d/.spectacular/archive/ideas/explicit.md"

  (cd "$d" && "$CLI" idea new legacy >/dev/null)
  out=$(cd "$d" && "$CLI" idea promote legacy 2>&1)
  assert_output "$out" "without --to is deprecated"
  assert_file "$d/.spectacular/requests/legacy/PLAN.md"
  rm -rf "$d"
}

scenario_shared_is_local_only() {
  echo "Scenario shared destination: reference required, no request created"
  local d="/tmp/spectacular-idea-shared" out
  new_ws "$d"
  (cd "$d" && "$CLI" idea new discuss >/dev/null)
  out=$(cd "$d" && "$CLI" idea promote discuss --to shared 2>&1 || true)
  assert_output "$out" "requires '--ref <stable-reference>'"
  assert_file "$d/.spectacular/ideas/discuss.md"

  (cd "$d" && "$CLI" idea promote discuss --to shared --ref "https://tracker.example/item/42" >/dev/null)
  assert_file "$d/.spectacular/archive/ideas/discuss.md"
  assert_contains "$d/.spectacular/archive/ideas/discuss.md" "shared:https://tracker.example/item/42"
  assert_absent "$d/.spectacular/requests/discuss"
  rm -rf "$d"
}

scenario_roadmap_requires_icebox_choice() {
  echo "Scenario roadmap destination: explicit Icebox placement"
  local d="/tmp/spectacular-idea-roadmap" out
  new_ws "$d"
  (cd "$d" && "$CLI" idea new roadmap-entry --hypothesis "Better capture routing" >/dev/null)
  out=$(cd "$d" && "$CLI" idea promote roadmap-entry --to roadmap 2>&1 || true)
  assert_output "$out" "requires '--placement icebox'"
  assert_file "$d/.spectacular/ideas/roadmap-entry.md"

  (cd "$d" && "$CLI" idea promote roadmap-entry --to roadmap --placement icebox >/dev/null)
  assert_file "$d/.spectacular/archive/ideas/roadmap-entry.md"
  assert_contains "$d/.spectacular/roadmaps/index.md" "Better capture routing"
  assert_absent "$d/.spectacular/requests/roadmap-entry"
  rm -rf "$d"
}

scenario_request_only_flags_rejected() {
  echo "Scenario guardrails: destination flags do not leak across routes"
  local d="/tmp/spectacular-idea-flags" out
  new_ws "$d"
  (cd "$d" && "$CLI" idea new guard >/dev/null)
  out=$(cd "$d" && "$CLI" idea promote guard --to shared --request-slug no 2>&1 || true)
  assert_output "$out" "--request-slug is valid only with --to request"
  assert_file "$d/.spectacular/ideas/guard.md"
  rm -rf "$d"
}

scenario_request_and_legacy
scenario_shared_is_local_only
scenario_roadmap_requires_icebox_choice
scenario_request_only_flags_rejected

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
