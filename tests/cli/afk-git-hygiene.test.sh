#!/usr/bin/env bash
# AFK Git policy, isolated branch, archive-first cleanup, and PR handoff.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }

init_repo() {
  local d="$1"
  rm -rf "$d"; mkdir -p "$d/.spectacular/specs" "$d/.spectacular/requests"
  printf 'project:\n  name: afk-test\nlast_build: 0\n' > "$d/.spectacular/config.yaml"
  printf '%s\n' '---' 'version: 1.0' 'updated: 2026-08-01' 'summary: "index"' '---' '# Index' > "$d/.spectacular/specs/index.md"
  (cd "$d" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)
}

scenario_policy_and_names() {
  echo "Scenario policy: status/configure are dry-run first; host prefixes compose with branch classes"
  local d="/tmp/spectacular-afk-policy" before after out code
  init_repo "$d"
  before=$(cksum "$d/.spectacular/config.yaml")
  out=$(cd "$d" && "$CLI" afk status)
  [[ "$out" == *"enabled:          false"* && "$out" == *"read-only"* ]] && pass || fail "default status explains disabled gate"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --allow-pr-create >/dev/null)
  after=$(cksum "$d/.spectacular/config.yaml")
  [[ "$before" == "$after" ]] && pass || fail "configure defaults to no mutation"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --allow-pr-create --apply >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "configure apply requires yes"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --allow-pr-create --apply --yes >/dev/null)
  assert_contains "$d/.spectacular/config.yaml" "enabled: true"
  assert_contains "$d/.spectacular/config.yaml" "branch_prefix: \"codex/\""
  [[ "$(cd "$d" && "$CLI" afk propose draft-spec auth)" == "codex/spec/draft-auth" ]] && pass || fail "draft branch name"
  [[ "$(cd "$d" && "$CLI" afk propose spike spk1)" == "codex/spike/prototype-spk-001" ]] && pass || fail "spike branch name"
  [[ "$(cd "$d" && "$CLI" afk propose fork idea1)" == "codex/fork/idea-idea-001" ]] && pass || fail "fork branch name"
  [[ "$(cd "$d" && "$CLI" afk propose feature billing --version 1.2.0)" == "codex/feat/v1.2.0-billing" ]] && pass || fail "feature branch name"
  rm -rf "$d"
}

scenario_isolation_and_cleanup() {
  echo "Scenario playground: clean preflight, provenance, dry-run cleanup, archive before delete"
  local d="/tmp/spectacular-afk-cleanup" out code branch="codex/spike/prototype-spk-001"
  init_repo "$d"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --apply --yes >/dev/null && git add . && git commit -qm policy)
  (cd "$d" && "$CLI" afk preflight spike >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "spike writes refused on primary"
  (cd "$d" && "$CLI" afk start spike spk1 >/dev/null)
  [[ "$(cd "$d" && git branch --show-current)" == main ]] && pass || fail "start defaults to proposal only"
  (cd "$d" && "$CLI" afk start spike spk1 --apply --yes >/dev/null)
  [[ "$(cd "$d" && git branch --show-current)" == "$branch" ]] && pass || fail "authorized start creates isolated branch"
  assert_contains "$d/.spectacular/afk/branches.md" "source: spk1"
  (cd "$d" && git add . && git commit -qm provenance && "$CLI" afk preflight spike >/dev/null)
  (cd "$d" && git switch -q main)
  (cd "$d" && "$CLI" afk cleanup "$branch" --disposition abandoned --outcome "Approach failed" --evidence "SPK-001 benchmark" >/dev/null)
  (cd "$d" && git show-ref --verify --quiet "refs/heads/$branch") && pass || fail "dry-run cleanup preserves branch"
  [[ ! -d "$d/.spectacular/archive/afk-branches" ]] && pass || fail "dry-run cleanup writes no archive"
  (cd "$d" && "$CLI" afk cleanup "$branch" --disposition abandoned --outcome "Approach failed" --evidence "SPK-001 benchmark" --apply --yes >/dev/null)
  (cd "$d" && git show-ref --verify --quiet "refs/heads/$branch") && code=0 || code=$?
  assert_exit "$code" 1 "confirmed cleanup deletes local branch"
  local archived; archived=$(find "$d/.spectacular/archive/afk-branches" -type f | head -1)
  assert_contains "$archived" "recoverable_from: git reflog"
  assert_contains "$archived" "**Evidence:** SPK-001 benchmark"
  (cd "$d" && "$CLI" afk cleanup nowhere --disposition abandoned --outcome x --evidence y --remote >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "remote deletion remains closed"
  rm -rf "$d"
}

scenario_pr_handoff() {
  echo "Scenario PR: verified/current/test gates; exact title; no merge"
  local d="/tmp/spectacular-afk-pr" bin="/tmp/spectacular-afk-bin" log="/tmp/spectacular-afk-gh.log" out code
  init_repo "$d"; rm -rf "$bin"; mkdir -p "$bin"; rm -f "$log"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --allow-pr-create --apply --yes >/dev/null)
  (cd "$d" && "$CLI" spec new billing --summary "Billing" >/dev/null && "$CLI" spec confirm s1 --evidence approved >/dev/null && "$CLI" spec act s1 --request-slug billing-work >/dev/null)
  sed -i.bak 's/status: planned/status: verified/' "$d/.spectacular/requests/billing-work/PLAN.md"; rm -f "$d/.spectacular/requests/billing-work/PLAN.md.bak"
  printf '%s\n' '---' 'updated: 2026-08-01' '---' '# Verify log' '' '**Outcome:** verified' > "$d/.spectacular/requests/billing-work/VERIFY-LOG.md"
  (cd "$d" && git add . && git commit -qm verified)
  (cd "$d" && "$CLI" afk start feature billing --version v1.0.0 --apply --yes >/dev/null && git add . && git commit -qm branch-provenance)

  (cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "fresh test proof required"
  out=$(cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed)
  [[ "$out" == *"[SpecTACular] Executed: v1.0.0 - Billing"* && "$out" == *"stop before merge"* ]] && pass || fail "dry-run prints exact title and merge boundary"
  (cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed --breaking --apply --yes >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "breaking change requires separate approval"

  printf '#!/bin/sh\nprintf "%%s\\n" "$*" > "%s"\nprintf "https://example.test/pr/1\\n"\n' "$log" > "$bin/gh"; chmod +x "$bin/gh"
  (cd "$d" && PATH="$bin:$PATH" "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed --apply --yes >/dev/null)
  assert_contains "$log" "pr create --title [SpecTACular] Executed: v1.0.0 - Billing"
  [[ "$(cd "$d" && git branch --show-current)" == "codex/feat/v1.0.0-billing" ]] && pass || fail "PR handoff does not merge or switch branch"
  rm -rf "$d" "$bin"; rm -f "$log"
}

scenario_policy_and_names
scenario_isolation_and_cleanup
scenario_pr_handoff

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
