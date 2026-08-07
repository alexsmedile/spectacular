#!/usr/bin/env bash
# tests/cli/workspace.test.sh — read-only concurrent-work coordination baseline.
set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { echo "$1" | grep -qF -- "$2" && pass || fail "missing: $2"; }

W=$(mktemp -d)
mkdir -p "$W/.spectacular/requests/current"
printf 'project:\n  name: workspace-test\nforge:\n  shared_base: "main"\n' > "$W/.spectacular/config.yaml"
printf '%s\n' '---' 'status: active' 'updated: 2026-08-07' 'summary: "current"' 'branch: feat/current' '---' '# Plan' > "$W/.spectacular/requests/current/PLAN.md"
(cd "$W" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init && git switch -q -c feat/current)
printf 'owned request note\n' > "$W/.spectacular/requests/current/NOTE.md"
printf 'unrelated\n' > "$W/unrelated.txt"
before=$(cd "$W" && git status --porcelain)
out=$(cd "$W" && bash "$CLI" workspace preflight --request current)
after=$(cd "$W" && git status --porcelain)

echo "── read-only mixed worktree preflight ──"
assert_contains "$out" "[untracked] .spectacular/requests/current/NOTE.md"
assert_contains "$out" "belongs-to-known-work"
assert_contains "$out" "[untracked] unrelated.txt"
assert_contains "$out" "needs-preservation-branch"
[[ "$before" == "$after" ]] && pass || fail "preflight must not mutate Git state"
rm -rf "$W"

echo "── scoped preservation keeps staged unrelated work intact ──"
W=$(mktemp -d)
mkdir -p "$W/.spectacular"
printf 'project:\n  name: preserve-test\n' > "$W/.spectacular/config.yaml"
(cd "$W" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)
printf 'spec draft\n' > "$W/draft-spec.md"
printf 'staged unrelated\n' > "$W/unrelated.txt"
(cd "$W" && git add -- unrelated.txt)
before=$(cd "$W" && git status --porcelain)
(cd "$W" && bash "$CLI" workspace preserve draft-spec --paths draft-spec.md >/dev/null)
after_preview=$(cd "$W" && git status --porcelain)
[[ "$before" == "$after_preview" ]] && pass || fail "preservation preview must not mutate"
(cd "$W" && bash "$CLI" workspace preserve draft-spec --paths draft-spec.md --apply --yes >/dev/null)
branch=$(cd "$W" && git for-each-ref --format='%(refname:short)' 'refs/heads/preserve/draft-spec-*' | head -1)
[[ -n "$branch" ]] && pass || fail "preservation creates named branch"
[[ "$(cd "$W" && git show "$branch:draft-spec.md")" == "spec draft" ]] && pass || fail "preservation commit contains explicit path"
[[ "$(cd "$W" && git status --porcelain)" == *"A  unrelated.txt"* ]] && pass || fail "staged unrelated path remains staged"
rm -rf "$W"
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
