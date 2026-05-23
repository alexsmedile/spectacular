#!/usr/bin/env bash
# tests/cli/mutator.test.sh — CLI mutator verbs (v0.8.0+)
#
# Tests the 5 verbs that replace skill-side manual file edits:
#   touch    — bump updated: field
#   new      — scaffold request
#   promote  — advance lifecycle state
#   snapshot — versioned copy + version bump
#   archive  — move to archive/ + rewrite inbound links

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_file_exists() {
  if [[ -f "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected file: $1"; fail_count=$((fail_count + 1)); fi
}
assert_dir_exists() {
  if [[ -d "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected dir: $1"; fail_count=$((fail_count + 1)); fi
}
assert_dir_absent() {
  if [[ ! -d "$1" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected absent dir: $1"; fail_count=$((fail_count + 1)); fi
}
assert_file_contains() {
  if [[ -f "$1" ]] && grep -qF -- "$2" "$1"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count + 1)); fi
}
assert_file_lacks() {
  if [[ -f "$1" ]] && ! grep -qF -- "$2" "$1"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected '$1' to NOT contain '$2'"; fail_count=$((fail_count + 1)); fi
}
assert_output_contains() {
  if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count + 1)); fi
}
assert_exit() {
  if [[ "$1" -eq "$2" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ $3: exit $1 want $2"; fail_count=$((fail_count + 1)); fi
}

seed_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir"
  (cd "$dir" && "$CLI" init --kit blank --name "$(basename "$dir")" >/dev/null 2>&1)
}

scenario_1_touch_basic() {
  echo "Scenario 1: touch bumps updated; idempotent; refuses non-frontmatter"
  local dir="/tmp/spectacular-mutator-1"
  seed_workspace "$dir"

  sed -i.bak 's/^updated:.*/updated: 2026-01-01/' "$dir/.spectacular/PRD.md"
  rm -f "$dir/.spectacular/PRD.md.bak"

  (cd "$dir" && "$CLI" touch .spectacular/PRD.md >/dev/null)
  assert_file_contains "$dir/.spectacular/PRD.md" "updated: $(date +%Y-%m-%d)"

  local out
  out=$(cd "$dir" && "$CLI" touch .spectacular/PRD.md 2>&1)
  assert_output_contains "$out" "already"

  echo "no frontmatter" > "$dir/plain.md"
  local code
  (cd "$dir" && "$CLI" touch plain.md >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "non-frontmatter file rejected"

  rm -rf "$dir"
}

scenario_2_new_basic() {
  echo "Scenario 2: new scaffolds PLAN+TASKS; refuses duplicates and invalid slugs"
  local dir="/tmp/spectacular-mutator-2"
  seed_workspace "$dir"

  (cd "$dir" && "$CLI" new feature-x --summary "test feature" >/dev/null)
  assert_file_exists "$dir/.spectacular/requests/feature-x/PLAN.md"
  assert_file_exists "$dir/.spectacular/requests/feature-x/TASKS.md"
  assert_file_contains "$dir/.spectacular/requests/feature-x/PLAN.md" 'status: planned'
  assert_file_contains "$dir/.spectacular/requests/feature-x/PLAN.md" 'priority: medium'
  assert_file_contains "$dir/.spectacular/requests/feature-x/PLAN.md" 'test feature'

  (cd "$dir" && "$CLI" new urgent-fix --status active --priority high --summary "urgent" >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/urgent-fix/PLAN.md" 'status: active'
  assert_file_contains "$dir/.spectacular/requests/urgent-fix/PLAN.md" 'priority: high'

  local code
  (cd "$dir" && "$CLI" new feature-x >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "duplicate slug refused"

  (cd "$dir" && "$CLI" new Bad_Slug >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "invalid slug refused"

  rm -rf "$dir"
}

scenario_3_promote_basic() {
  echo "Scenario 3: promote advances state; refuses backward without --force"
  local dir="/tmp/spectacular-mutator-3"
  seed_workspace "$dir"
  (cd "$dir" && "$CLI" new req1 --summary "one" >/dev/null)

  (cd "$dir" && "$CLI" promote req1 >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" 'status: active'
  assert_file_contains "$dir/.spectacular/requests/req1/TASKS.md" 'status: active'

  (cd "$dir" && "$CLI" promote req1 >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" 'status: review'

  (cd "$dir" && "$CLI" promote req1 --to verified >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" 'status: verified'

  local code
  (cd "$dir" && "$CLI" promote req1 --to active >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "backward transition refused without --force"

  (cd "$dir" && "$CLI" promote req1 --to active --force >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/req1/PLAN.md" 'status: active'

  (cd "$dir" && "$CLI" promote no-such-slug >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "unknown slug refused"

  rm -rf "$dir"
}

scenario_4_snapshot_basic() {
  echo "Scenario 4: snapshot creates @vN; bumps version; idempotent on no body change"
  local dir="/tmp/spectacular-mutator-4"
  seed_workspace "$dir"

  (cd "$dir" && "$CLI" snapshot .spectacular/PRD.md >/dev/null)
  assert_file_exists "$dir/.spectacular/PRD@v1.md"
  assert_file_contains "$dir/.spectacular/PRD.md" "version: 1.2"

  local out
  out=$(cd "$dir" && "$CLI" snapshot .spectacular/PRD.md 2>&1)
  assert_output_contains "$out" "no body changes"
  if [[ -f "$dir/.spectacular/PRD@v2.md" ]]; then
    echo "    ✗ second snapshot should not have created PRD@v2.md"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi

  echo "new content" >> "$dir/.spectacular/PRD.md"
  (cd "$dir" && "$CLI" snapshot .spectacular/PRD.md >/dev/null)
  assert_file_exists "$dir/.spectacular/PRD@v2.md"

  echo "---" > "$dir/random.md"
  echo "title: foo" >> "$dir/random.md"
  echo "---" >> "$dir/random.md"
  local code
  (cd "$dir" && "$CLI" snapshot random.md >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "non-canonical file refused"

  rm -rf "$dir"
}

scenario_5_archive_basic() {
  echo "Scenario 5: archive moves dir + rewrites inbound related: links"
  local dir="/tmp/spectacular-mutator-5"
  seed_workspace "$dir"

  (cd "$dir" && "$CLI" new req-a --summary "first" >/dev/null)
  (cd "$dir" && "$CLI" new req-b --summary "depends on a" >/dev/null)

  sed -i.bak 's|  - PRD.md|  - PRD.md\
  - ../req-a/PLAN.md|' "$dir/.spectacular/requests/req-b/PLAN.md"
  rm -f "$dir/.spectacular/requests/req-b/PLAN.md.bak"
  assert_file_contains "$dir/.spectacular/requests/req-b/PLAN.md" "../req-a/PLAN.md"

  local code
  (cd "$dir" && "$CLI" archive req-a >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "planned status refused without --force"

  (cd "$dir" && "$CLI" promote req-a --to verified --force >/dev/null)
  (cd "$dir" && "$CLI" archive req-a >/dev/null)
  assert_dir_absent "$dir/.spectacular/requests/req-a"
  assert_dir_exists "$dir/.spectacular/archive/req-a"
  assert_file_contains "$dir/.spectacular/archive/req-a/PLAN.md" 'status: archived'
  assert_file_contains "$dir/.spectacular/archive/req-a/PLAN.md" "archived: $(date +%Y-%m-%d)"

  assert_file_contains "$dir/.spectacular/requests/req-b/PLAN.md" "../../archive/req-a/PLAN.md"
  assert_file_lacks "$dir/.spectacular/requests/req-b/PLAN.md" "../req-a/PLAN.md"

  (cd "$dir" && "$CLI" archive no-such-slug >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "unknown slug refused"

  (cd "$dir" && "$CLI" archive req-a >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "already-archived slug refused"

  rm -rf "$dir"
}

scenario_6_promote_archive_combo() {
  echo "Scenario 6: promote --archive chains into archive when reaching verified"
  local dir="/tmp/spectacular-mutator-6"
  seed_workspace "$dir"
  (cd "$dir" && "$CLI" new combo --summary "test" >/dev/null)

  (cd "$dir" && "$CLI" promote combo --to verified --force --archive >/dev/null)
  assert_dir_absent "$dir/.spectacular/requests/combo"
  assert_dir_exists "$dir/.spectacular/archive/combo"
  assert_file_contains "$dir/.spectacular/archive/combo/PLAN.md" 'status: archived'

  rm -rf "$dir"
}

scenario_8_new_target_version() {
  echo "Scenario 8: new --target-version writes target_version: to PLAN frontmatter"
  local dir="/tmp/spectacular-mutator-8"
  seed_workspace "$dir"

  (cd "$dir" && "$CLI" new feat-tagged --target-version 0.9.0 --summary "tagged" >/dev/null)
  assert_file_contains "$dir/.spectacular/requests/feat-tagged/PLAN.md" "target_version: 0.9.0"

  # Without --target-version, no target_version field added
  (cd "$dir" && "$CLI" new feat-untagged --summary "untagged" >/dev/null)
  assert_file_lacks "$dir/.spectacular/requests/feat-untagged/PLAN.md" "target_version:"

  rm -rf "$dir"
}

scenario_7_help_flags() {
  echo "Scenario 7: each verb --help exits 0 with usage line"
  local out code
  for verb in touch new promote snapshot archive; do
    out=$("$CLI" "$verb" --help 2>&1) && code=0 || code=$?
    assert_exit "$code" 0 "$verb --help exits 0"
    assert_output_contains "$out" "Usage: spectacular $verb" "$verb help shows usage"
  done
}

scenario_9_doctor_precondition_archive() {
  echo "Scenario 9: archive refuses on doctor errors; --skip-doctor bypasses"
  local dir="/tmp/spectacular-mutator-9"
  seed_workspace "$dir"
  (cd "$dir" && "$CLI" new req-broken --summary "test" >/dev/null)
  (cd "$dir" && "$CLI" promote req-broken --to verified --force >/dev/null)

  # Break workspace: delete required AGENTS.md (workspace area error)
  rm "$dir/.spectacular/AGENTS.md"

  # Archive should refuse
  local out code
  out=$(cd "$dir" && "$CLI" archive req-broken 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "archive refused on doctor error"
  assert_output_contains "$out" "Refusing 'archive'"
  assert_output_contains "$out" "--skip-doctor"

  # Request should still be in requests/, not archive/
  assert_dir_exists "$dir/.spectacular/requests/req-broken"
  assert_dir_absent "$dir/.spectacular/archive/req-broken"

  # --skip-doctor should bypass
  (cd "$dir" && "$CLI" archive req-broken --skip-doctor >/dev/null)
  assert_dir_absent "$dir/.spectacular/requests/req-broken"
  assert_dir_exists "$dir/.spectacular/archive/req-broken"

  rm -rf "$dir"
}

scenario_10_doctor_precondition_clean_passes() {
  echo "Scenario 10: archive on clean workspace passes through precondition"
  local dir="/tmp/spectacular-mutator-10"
  seed_workspace "$dir"
  (cd "$dir" && "$CLI" new req-ok --summary "test" >/dev/null)
  (cd "$dir" && "$CLI" promote req-ok --to verified --force >/dev/null)

  # Clean workspace — archive should succeed without --skip-doctor
  local out code
  out=$(cd "$dir" && "$CLI" archive req-ok 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "archive succeeded on clean workspace"
  assert_dir_exists "$dir/.spectacular/archive/req-ok"

  rm -rf "$dir"
}

echo "=== mutator.test.sh ==="
scenario_1_touch_basic
scenario_2_new_basic
scenario_3_promote_basic
scenario_4_snapshot_basic
scenario_5_archive_basic
scenario_6_promote_archive_combo
scenario_7_help_flags
scenario_8_new_target_version
scenario_9_doctor_precondition_archive
scenario_10_doctor_precondition_clean_passes

echo ""
echo "Results: ${pass_count} passed, ${fail_count} failed"
exit $((fail_count > 0))
