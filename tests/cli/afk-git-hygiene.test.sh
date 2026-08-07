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
  (cd "$d" && "$CLI" afk run start spike-test --goal "Test SPK-001" --allowed-actions "spikes" --apply --yes >/dev/null && git add . && git commit -qm afk-run)
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
  [[ ! -d "$d/.spectacular/archive/afk-branches" ]] && pass || fail "dry-run cleanup writes no workspace residue"
  (cd "$d" && "$CLI" afk cleanup "$branch" --disposition abandoned --outcome "Approach failed" --evidence "SPK-001 benchmark" --apply --yes >/dev/null)
  (cd "$d" && git show-ref --verify --quiet "refs/heads/$branch") && code=0 || code=$?
  assert_exit "$code" 1 "confirmed cleanup deletes local branch"
  (cd "$d" && git for-each-ref --format='%(refname)' refs/spectacular/archive | grep -q '^refs/spectacular/archive/') && pass || fail "confirmed cleanup preserves an archive ref"
  [[ -z "$(cd "$d" && git status --porcelain)" ]] && pass || fail "confirmed cleanup leaves no workspace residue"
  (cd "$d" && "$CLI" afk cleanup nowhere --disposition abandoned --outcome x --evidence y --remote >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "remote deletion remains closed"
  rm -rf "$d"
}

scenario_merged_remote_cleanup() {
  echo "Scenario merged cleanup: one confirmation deletes a matching remote branch and preserves mismatches"
  local d="/tmp/spectacular-afk-remote-cleanup" remote="/tmp/spectacular-afk-remote-cleanup.git" branch="feat/remote-cleanup" guard="feat/remote-mismatch" code out archived
  init_repo "$d"; rm -rf "$remote"; git init --bare -q "$remote"
  (cd "$d" && "$CLI" afk configure --enable --apply --yes >/dev/null && git add . && git commit -qm policy && git remote add origin "$remote" && git push -q -u origin main && git switch -q -c "$branch" && touch remote.txt && git add remote.txt && git commit -qm feature && git push -q -u origin "$branch" && git switch -q main && git merge -q "$branch" && git push -q origin main)
  out=$(cd "$d" && "$CLI" afk cleanup "$branch" --disposition merged --outcome "Merged" --evidence "PR-1")
  [[ "$out" == *"remote: delete origin/$branch"* ]] && pass || fail "merged cleanup previews remote deletion"
  (cd "$d" && "$CLI" afk cleanup "$branch" --disposition merged --outcome "Merged" --evidence "PR-1" --apply --yes >/dev/null)
  (cd "$d" && git show-ref --verify --quiet "refs/heads/$branch") && code=0 || code=$?
  assert_exit "$code" 1 "merged cleanup deletes local branch"
  [[ -z "$(git --git-dir="$remote" for-each-ref --format='%(refname)' "refs/heads/$branch")" ]] && pass || fail "merged cleanup deletes matching remote branch"
  (cd "$d" && git switch -q -c "$guard" && touch guard.txt && git add guard.txt && git commit -qm guard && git push -q -u origin "$guard" && git switch -q main && git merge -q "$guard" && git push -q origin main && git switch -q -c remote-advance "$guard" && touch remote-advance.txt && git add remote-advance.txt && git commit -qm remote-advance && git push -q origin HEAD:"$guard" && git switch -q main && git branch -D remote-advance >/dev/null)
  out=$(cd "$d" && "$CLI" afk cleanup "$guard" --disposition merged --outcome "Merged" --evidence "PR-2" 2>&1 || true)
  [[ "$out" == *"moved past the merged local tip"* ]] && pass || fail "mismatched remote branch blocks cleanup"
  (cd "$d" && git show-ref --verify --quiet "refs/heads/$guard") && pass || fail "mismatched remote branch remains local"
  rm -rf "$d" "$remote"
}

scenario_pr_handoff() {
  echo "Scenario PR: verified/approved/test gates; reviewer-facing body; no merge"
  local d="/tmp/spectacular-afk-pr" bin="/tmp/spectacular-afk-bin" log="/tmp/spectacular-afk-gh.log" out code
  init_repo "$d"; rm -rf "$bin"; mkdir -p "$bin"; rm -f "$log"
  (cd "$d" && "$CLI" afk configure --enable --branch-prefix codex/ --allow-pr-create --apply --yes >/dev/null)
  (cd "$d" && "$CLI" spec new billing --summary "Billing" >/dev/null && git add . && git commit -qm billing-contract && git switch -q -c request/billing-work)
  local contract_id; contract_id=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/specs/billing.md")
  (cd "$d" && "$CLI" request new billing-work --from "$contract_id" >/dev/null)
  sed -i.bak 's/status: planned/status: verified/' "$d/.spectacular/requests/billing-work/PLAN.md"; rm -f "$d/.spectacular/requests/billing-work/PLAN.md.bak"
  printf '%s\n' '---' 'updated: 2026-08-01' '---' '# Verify log' '' '**Outcome:** verified' > "$d/.spectacular/requests/billing-work/VERIFY-LOG.md"
  (cd "$d" && git add . && git commit -qm verified)
  (cd "$d" && "$CLI" afk run start billing-run --goal "Implement approved billing spec" --request billing-work --allowed-actions "feature" --apply --yes >/dev/null && git add . && git commit -qm afk-run)
  (cd "$d" && "$CLI" afk start feature billing --version v1.0.0 --apply --yes >/dev/null && git add . && git commit -qm branch-provenance)

  (cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "fresh test proof required"
  out=$(cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed --summary "Add billing target validation" --summary "Document the retained target boundary" --validation "Contract fixtures pass")
  [[ "$out" == *"[Spectacular] Executed: v1.0.0 - Billing"* && "$out" == *"## In plain language"* && "$out" == *"## What this PR changes"* && "$out" == *"- Add billing target validation"* && "$out" == *"- Document the retained target boundary"* && "$out" == *"- Contract fixtures pass"* && "$out" == *"stop before merge"* ]] && pass || fail "dry-run prints a reviewer-facing title, body, and merge boundary"
  (cd "$d" && "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed --breaking --apply --yes >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "breaking change requires separate approval"

  printf '#!/bin/sh\nprintf "%%s\\n" "$*" > "%s"\nprintf "https://example.test/pr/1\\n"\n' "$log" > "$bin/gh"; chmod +x "$bin/gh"
  (cd "$d" && PATH="$bin:$PATH" "$CLI" afk pr billing-work --version v1.0.0 --name Billing --tests-passed --summary "Add billing target validation" --summary "Document the retained target boundary" --validation "Contract fixtures pass" --apply --yes >/dev/null)
  assert_contains "$log" "pr create --title [Spectacular] Executed: v1.0.0 - Billing"
  assert_contains "$log" "--draft"
  assert_contains "$log" "## In plain language"
  assert_contains "$log" "Add billing target validation"
  assert_contains "$log" "Document the retained target boundary"
  assert_contains "$log" "## Validation"
  assert_contains "$log" "Contract fixtures pass"
  assert_contains "$log" "Merge remains human-gated."
  assert_contains "$log" "Filed with [Spectacular](https://github.com/alexsmedile/spectacular)."
  grep -qF "Executed s1 through verified request" "$log" && fail "legacy opaque PR body removed" || pass
  [[ "$(cd "$d" && git branch --show-current)" == "codex/feat/v1.0.0-billing" ]] && pass || fail "PR handoff does not merge or switch branch"
  rm -rf "$d" "$bin"; rm -f "$log"
}

scenario_authority_record() {
  echo "Scenario authority record: opt-in record/event paths are auditable and make no Git mutation"
  local d="/tmp/spectacular-afk-authority" out code head branch
  init_repo "$d"
  head=$(cd "$d" && git rev-parse HEAD); branch=$(cd "$d" && git branch --show-current)
  (cd "$d" && "$CLI" afk run start afk-001 --goal "Validate AFK record" --authority-record --authorized-by alex --session S-001 --base main --goal-node N-001 --apply --yes >/dev/null)
  assert_contains "$d/.spectacular/afk/runs/afk-001.md" "authority_version: 1"
  assert_contains "$d/.spectacular/afk/runs/afk-001.md" "integration_branch: \"afk/afk-001/integration\""
  assert_contains "$d/.spectacular/afk/runs/afk-001.md" "### E-001 — authorization"
  [[ "$(cd "$d" && git rev-parse HEAD)" == "$head" && "$(cd "$d" && git branch --show-current)" == "$branch" ]] && pass || fail "authority record start makes no Git mutation"
  (cd "$d" && "$CLI" afk run event technical-choice --scope N-001 --rationale "Reuse existing helper" --evidence "cli/spectacular" >/dev/null)
  grep -q "E-002" "$d/.spectacular/afk/runs/afk-001.md" && fail "dry-run event must not append" || pass
  (cd "$d" && "$CLI" afk run event technical-choice --scope N-001 --rationale "Reuse existing helper" --evidence "cli/spectacular" --apply --yes >/dev/null)
  assert_contains "$d/.spectacular/afk/runs/afk-001.md" "### E-002 — technical-choice"
  out=$(cd "$d" && "$CLI" afk run status)
  [[ "$out" == *"authority: alex · session S-001"* && "$out" == *"events: 2"* ]] && pass || fail "status exposes authority and event count"
  (cd "$d" && "$CLI" doctor afk) >/dev/null && pass || fail "doctor accepts complete authority record"
  sed -i.bak 's/### E-002/### E-004/' "$d/.spectacular/afk/runs/afk-001.md"; rm -f "$d/.spectacular/afk/runs/afk-001.md.bak"
  out=$(cd "$d" && "$CLI" doctor afk)
  [[ "$out" == *"non-monotonic IDs"* ]] && pass || fail "doctor detects invalid event ordering"
  rm -rf "$d"
}

scenario_policy_and_names
scenario_isolation_and_cleanup
scenario_merged_remote_cleanup
scenario_pr_handoff
scenario_authority_record

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
