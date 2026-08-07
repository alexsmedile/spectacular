#!/usr/bin/env bash
# Read-only Git commit review emitted by `spectacular session end`.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { [[ "$1" == *"$2"* ]] && pass || fail "$3: missing '$2'"; }
assert_eq() { [[ "$1" == "$2" ]] && pass || fail "$3: got '$1', want '$2'"; }

seed_workspace() {
  local d="$1"
  rm -rf "$d"; mkdir -p "$d"
  (cd "$d" && "$CLI" init --kit blank --name session-review >/dev/null)
  (cd "$d" && git init -q && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)
  (cd "$d" && "$CLI" spec new owned --summary "Ownership fixture" >/dev/null && git add . && git commit -qm contract && git switch -q -c request/owned)
  local contract_id; contract_id=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/specs/owned.md")
  (cd "$d" && "$CLI" request new owned --from "$contract_id" >/dev/null && git add . && git commit -qm request)
  (cd "$d" && "$CLI" session start --tag commit-review >/dev/null && git add . && git commit -qm session-start)
}

scenario_dirty_tree_is_read_only() {
  echo "Scenario dirty tree: staged, unstaged, untracked, and request paths are reported without Git mutation"
  local d="/tmp/spectacular-session-commit-review" out before_head after_head before_index after_index before_config after_config
  seed_workspace "$d"
  printf 'staged\n' > "$d/staged.md"
  (cd "$d" && git add staged.md)
  printf 'unstaged\n' > "$d/unstaged.md"
  printf '\nReview ownership fixture.\n' >> "$d/.spectacular/requests/owned/PLAN.md"
  printf 'untracked\n' > "$d/untracked.md"
  before_head=$(cd "$d" && git rev-parse HEAD)
  before_index=$(cd "$d" && git write-tree)
  before_config=$(cksum "$d/.git/config")

  out=$(cd "$d" && "$CLI" session end)
  after_head=$(cd "$d" && git rev-parse HEAD)
  after_index=$(cd "$d" && git write-tree)
  after_config=$(cksum "$d/.git/config")

  assert_contains "$out" "Git commit review (read-only)" "review header"
  assert_contains "$out" "Staged files: 1" "staged count"
  assert_contains "$out" "Unstaged files: 1" "unstaged count includes request edit"
  assert_contains "$out" "Untracked files: 2" "untracked count"
  assert_contains "$out" "staged.md" "changed staged path"
  assert_contains "$out" ".spectacular/requests/owned/PLAN.md" "request-owned path"
  assert_contains "$out" "Request-owned paths (direct folder match only)" "ownership disclaimer"
  assert_contains "$out" "Untracked files have no diffstat until staged." "untracked evidence limit"
  assert_contains "$out" "Suggested, human must verify:" "human verification label"
  assert_contains "$out" "cannot determine whether these changes are one commit or several" "no semantic grouping claim"
  assert_contains "$out" "git add -p" "human follow-up"
  assert_contains "$out" "did not stage, commit, amend, push, merge, reset, or stash" "non-mutation invariant"
  assert_eq "$after_head" "$before_head" "HEAD unchanged"
  assert_eq "$after_index" "$before_index" "Git index unchanged"
  assert_eq "$after_config" "$before_config" "Git config unchanged"
  rm -rf "$d"
}

scenario_clean_tree() {
  echo "Scenario clean tree: no proposal"
  local d="/tmp/spectacular-session-commit-clean" out
  seed_workspace "$d"
  out=$(cd "$d" && "$CLI" session end)
  assert_contains "$out" "Working tree: clean; no commit proposal." "clean tree result"
  rm -rf "$d"
}

scenario_dry_run_and_no_repository() {
  echo "Scenario dry-run/no repository: review is available without extra Git authority"
  local d="/tmp/spectacular-session-commit-dry" bare="/tmp/spectacular-session-commit-no-repo" out
  seed_workspace "$d"
  printf 'dirty\n' > "$d/dirty.md"
  out=$(cd "$d" && "$CLI" session end --dry-run)
  assert_contains "$out" "would close" "dry-run keeps close preview"
  assert_contains "$out" "Git commit review (read-only)" "dry-run review"
  rm -rf "$d"; mkdir -p "$bare"
  (cd "$bare" && "$CLI" init --kit blank --name no-repo >/dev/null && "$CLI" session start --tag no-repo >/dev/null)
  out=$(cd "$bare" && "$CLI" session end)
  assert_contains "$out" "Git commit review unavailable: not inside a Git repository." "no-repository result"
  rm -rf "$bare"
}

scenario_dirty_tree_is_read_only
scenario_clean_tree
scenario_dry_run_and_no_repository

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
