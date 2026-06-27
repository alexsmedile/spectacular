#!/usr/bin/env bash
# tests/cli/undo.test.sh — `spectacular undo` reverse gear (v1.22.0+)
#
# Covers the four undo milestones:
#   M1 advance undo      — status back one step (PLAN + TASKS), breadcrumb cleared
#   M2 archive undo      — dir back, status restored, archived: dropped, links reversed
#   M3 idea-promote undo — idea restored + status reset; request dir left by default
#   M4 guardrails        — nothing-to-undo, --dry-run mutates nothing, stale breadcrumb refused

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_dir_exists()  { if [[ -d "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected dir: $1"; fail_count=$((fail_count+1)); fi; }
assert_dir_absent()  { if [[ ! -d "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected absent dir: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_exists() { if [[ -f "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected file: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_absent() { if [[ ! -f "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected absent file: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_contains() { if [[ -f "$1" ]] && grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_file_lacks()    { if [[ -f "$1" ]] && ! grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to NOT contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_output_contains() { if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count+1)); else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count+1)); fi; }
assert_exit() { if [[ "$1" -eq "$2" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ $3: exit $1 want $2"; fail_count=$((fail_count+1)); fi; }

seed() {
  local dir="$1"
  rm -rf "$dir"; mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --kit blank --name "$(basename "$dir")" >/dev/null 2>&1)
}

scenario_1_undo_advance() {
  echo "Scenario 1 (M1): undo reverses advance; breadcrumb cleared"
  local dir="/tmp/spectacular-undo-1"
  seed "$dir"
  (cd "$dir" && "$CLI" new req1 --summary "t" >/dev/null)
  (cd "$dir" && "$CLI" advance req1 >/dev/null)        # planned → active
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" "status: active"
  assert_file_exists "$dir/.spectacular/.last-mutation"
  (cd "$dir" && "$CLI" undo >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" "status: planned"
  assert_file_contains "$dir/.spectacular/requests/req1/TASKS.md" "status: planned"
  assert_file_absent "$dir/.spectacular/.last-mutation"
  rm -rf "$dir"
}

scenario_2_undo_archive() {
  echo "Scenario 2 (M2): undo reverses archive — dir, status, archived:, links"
  local dir="/tmp/spectacular-undo-2"
  seed "$dir"
  (cd "$dir" && git init -q 2>/dev/null)
  (cd "$dir" && "$CLI" new req-a --summary "a" >/dev/null)
  (cd "$dir" && "$CLI" new req-b --summary "b" >/dev/null)
  printf '\nrelated:\n  - ../req-a/PLAN.md\n' >> "$dir/.spectacular/requests/req-b/PLAN.md"
  (cd "$dir" && "$CLI" advance req-a --to verified --force >/dev/null)
  (cd "$dir" && "$CLI" archive req-a --skip-doctor >/dev/null 2>&1)
  assert_dir_exists "$dir/.spectacular/archive/req-a"
  assert_dir_absent "$dir/.spectacular/requests/req-a"
  assert_file_contains "$dir/.spectacular/requests/req-b/PLAN.md" "../../archive/req-a/"
  # undo
  (cd "$dir" && "$CLI" undo >/dev/null 2>&1)
  assert_dir_exists "$dir/.spectacular/requests/req-a"
  assert_dir_absent "$dir/.spectacular/archive/req-a"
  assert_file_contains "$dir/.spectacular/requests/req-a/PLAN.md" "status: verified"
  assert_file_lacks "$dir/.spectacular/requests/req-a/PLAN.md" "archived:"
  assert_file_contains "$dir/.spectacular/requests/req-b/PLAN.md" "../req-a/PLAN.md"
  assert_file_lacks "$dir/.spectacular/requests/req-b/PLAN.md" "../../archive/req-a/"
  rm -rf "$dir"
}

scenario_3_undo_idea_promote() {
  echo "Scenario 3 (M3): undo reverses idea promote; request dir left by default"
  local dir="/tmp/spectacular-undo-3"
  seed "$dir"
  (cd "$dir" && "$CLI" idea new myidea >/dev/null 2>&1)
  (cd "$dir" && "$CLI" idea promote myidea >/dev/null 2>&1)
  assert_file_absent "$dir/.spectacular/ideas/myidea.md"
  assert_file_exists "$dir/.spectacular/archive/ideas/myidea.md"
  assert_dir_exists "$dir/.spectacular/requests/myidea"
  # undo, answer N to the removal prompt (no /dev/tty in CI → defaults to n)
  (cd "$dir" && echo "n" | "$CLI" undo >/dev/null 2>&1)
  assert_file_exists "$dir/.spectacular/ideas/myidea.md"
  assert_file_contains "$dir/.spectacular/ideas/myidea.md" "status: parked"
  assert_file_lacks "$dir/.spectacular/ideas/myidea.md" "promoted_to:"
  assert_dir_exists "$dir/.spectacular/requests/myidea"   # left in place (default)
  rm -rf "$dir"
}

scenario_4_guardrails() {
  echo "Scenario 4 (M4): nothing-to-undo, --dry-run, stale refusal"
  local dir="/tmp/spectacular-undo-4" out code
  seed "$dir"
  # Nothing to undo
  out=$(cd "$dir" && "$CLI" undo 2>&1); code=$?
  assert_exit "$code" 0 "nothing-to-undo exits 0"
  assert_output_contains "$out" "Nothing to undo"
  # --dry-run mutates nothing
  (cd "$dir" && "$CLI" new req1 --summary "t" >/dev/null)
  (cd "$dir" && "$CLI" advance req1 >/dev/null)
  out=$(cd "$dir" && "$CLI" undo --dry-run 2>&1)
  assert_output_contains "$out" "dry-run"
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" "status: active"  # unchanged
  assert_file_exists "$dir/.spectacular/.last-mutation"                            # breadcrumb kept
  # Stale breadcrumb: edit the file after the mutation → undo refuses
  sleep 2
  echo "" >> "$dir/.spectacular/requests/req1/PLAN.md"
  out=$(cd "$dir" && "$CLI" undo 2>&1); code=$?
  assert_exit "$code" 1 "stale breadcrumb refused"
  assert_output_contains "$out" "stale"
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" "status: active"  # still active
  # Bad argument rejected
  out=$(cd "$dir" && "$CLI" undo bogus 2>&1); code=$?
  assert_exit "$code" 1 "undo rejects unexpected args"
  rm -rf "$dir"
}

scenario_1_undo_advance
scenario_2_undo_archive
scenario_3_undo_idea_promote
scenario_4_guardrails

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
