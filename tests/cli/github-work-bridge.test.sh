#!/usr/bin/env bash
# GitHub work routing, Issue/goal provenance, PR handoff, and reconciliation.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
WS="/tmp/spectacular-github-work-bridge"
REMOTE="/tmp/spectacular-github-work-bridge-remote.git"
BIN="/tmp/spectacular-github-work-bridge-bin"
GH_LOG="/tmp/spectacular-github-work-bridge-gh.log"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_has() { [[ "$1" == *"$2"* ]] && pass || fail "$3: missing '$2'"; }

rm -rf "$WS" "$REMOTE" "$BIN"; rm -f "$GH_LOG"; mkdir -p "$WS" "$BIN"
(cd "$WS" && "$CLI" init --minimal --name github-bridge-test --skill-scope none >/dev/null)
(cd "$WS" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)

echo "Scenario Issue and goal sources create lean requests without inventing specs"
(cd "$WS" && "$CLI" request new cache-fix --from-issue https://github.com/acme/widget/issues/12 --summary "Prevent stale cache reads" --sensitivity normal >/dev/null)
ISSUE_PLAN="$WS/.spectacular/requests/cache-fix/PLAN.md"
grep -q '^source_type: issue' "$ISSUE_PLAN" && pass || fail "Issue source type recorded"
grep -q '^source_ref: "acme/widget#12"' "$ISSUE_PLAN" && pass || fail "Issue URL normalized to canonical identity"
[[ ! -d "$WS/.spectacular/specs/SPC-001-cache-fix.md" ]] && pass || fail "Issue request does not manufacture spec"
code=0; (cd "$WS" && "$CLI" request new duplicate --from-issue acme/widget#12 --summary duplicate --sensitivity normal >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "duplicate Issue owner refused"
(cd "$WS" && "$CLI" request new release-check --from-goal release-readiness --summary "Prove release readiness" --sensitivity normal >/dev/null)
grep -q '^source_type: goal' "$WS/.spectacular/requests/release-check/PLAN.md" && pass || fail "goal source type recorded"
grep -q '^source_ref: release-readiness' "$WS/.spectacular/requests/release-check/PLAN.md" && pass || fail "goal source reference recorded"
code=0; (cd "$WS" && "$CLI" request new invalid --from-issue '#7' --summary invalid --sensitivity normal >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "ambiguous naked Issue number refused"
code=0; (cd "$WS" && "$CLI" request new unclassified --from-issue acme/widget#13 --summary unsafe >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "Issue request requires explicit sensitivity"

echo "Scenario integration manifest uses correct Issue relationship and opens a draft PR"
sed -i.bak \
  -e 's/<The current behavior\/structure this request touches\.>/Existing cache behavior is defined by code and tests./' \
  -e 's/<The specific surfaces this request modifies\.>/The accepted cache correction and regression test./' \
  -e 's/<The boundary — what this change deliberately leaves alone\.>/No product or API behavior beyond the accepted correction./' \
  "$ISSUE_PLAN"; rm -f "$ISSUE_PLAN.bak"
(cd "$WS" && "$CLI" request advance cache-fix >/dev/null && git add . && git commit -qm request)
git init -q --bare "$REMOTE"
(cd "$WS" && git remote add origin "$REMOTE" && git switch -q -c codex/cache-fix && printf 'change\n' > implementation.txt && git add implementation.txt && git commit -qm implementation && git push -qu origin codex/cache-fix)
out=$(cd "$WS" && "$CLI" github pr open cache-fix --name "Cache fix" --summary "Prevent stale reads" --validation "Regression test passes")
assert_has "$out" "## In plain language" "manifest leads with plain-language context"
assert_has "$out" "## Why this matters" "manifest explains why the request matters"
assert_has "$out" "## Decision requested" "manifest states the review decision"
assert_has "$out" "## What this PR changes" "manifest separates the change list"
assert_has "$out" "## What this PR does not change" "manifest makes the boundary explicit"
assert_has "$out" "## Review focus" "manifest tells reviewers what to assess"
assert_has "$out" "## Follow-up" "manifest explains the next boundary"
assert_has "$out" "<summary>Technical details</summary>" "manifest collapses technical provenance"
assert_has "$out" "Fixes acme/widget#12" "on-merge Issue closes from manifest"
assert_has "$out" "Spectacular request: \`cache-fix\`" "request provenance in manifest"
assert_has "$out" "open a draft PR" "dry run declares draft boundary"
out=$(cd "$WS" && "$CLI" github pr open cache-fix --name "Cache fix" --resolution on_release)
assert_has "$out" "Refs acme/widget#12" "release-gated Issue uses non-closing reference"
assert_has "$out" "must not close it on merge" "release gate explained"

HEAD_SHA=$(cd "$WS" && git rev-parse HEAD)
printf '#!/bin/sh\nprintf "%%s\\n" "$*" >> "%s"\nif [ "$1 $2" = "pr create" ]; then printf "https://github.com/acme/widget/pull/9\\n"; exit 0; fi\nif [ "$1 $2" = "pr view" ]; then printf "%%s\\n" "%s"; exit 0; fi\nif [ "$1 $2" = "pr checks" ]; then printf "no required checks reported\\n"; exit 1; fi\nif [ "$1 $2" = "issue view" ]; then printf "CLOSED\\n"; exit 0; fi\nexit 0\n' "$GH_LOG" "$HEAD_SHA" > "$BIN/gh"; chmod +x "$BIN/gh"
(cd "$WS" && PATH="$BIN:$PATH" "$CLI" github pr open cache-fix --name "Cache fix" --apply --yes >/dev/null)
grep -q -- '--draft' "$GH_LOG" && pass || fail "PR creation is draft"
grep -q '^github_pr: "https://github.com/acme/widget/pull/9"' "$ISSUE_PLAN" && pass || fail "PR URL recorded on request"

echo "Scenario ready handoff requires verified current-head evidence and never merges"
sed -i.bak 's/status: active/status: verified/' "$ISSUE_PLAN"; rm -f "$ISSUE_PLAN.bak"
TASKS="$WS/.spectacular/requests/cache-fix/TASKS.md"
sed -i.bak 's/status: active/status: verified/' "$TASKS"; rm -f "$TASKS.bak"
printf '%s\n' '---' 'result: pass' '---' '# Verify log' '' "**Outcome:** verified" '' "against: $HEAD_SHA" > "$WS/.spectacular/requests/cache-fix/VERIFY-LOG.md"
(cd "$WS" && PATH="$BIN:$PATH" "$CLI" github pr ready cache-fix --apply --yes >/dev/null)
grep -q 'pr ready https://github.com/acme/widget/pull/9' "$GH_LOG" && pass || fail "ready command called"
grep -q 'pr merge' "$GH_LOG" && fail "bridge must never merge" || pass
(cd "$WS" && git add .spectacular/requests/cache-fix && git commit -qm lifecycle-metadata && git push -q)
READY_HEAD=$(cd "$WS" && git rev-parse HEAD)
sed -i.bak "s/$HEAD_SHA/$READY_HEAD/" "$BIN/gh"; rm -f "$BIN/gh.bak"
json=$(cd "$WS" && PATH="$BIN:$PATH" "$CLI" github reconcile cache-fix --json)
assert_has "$json" '"status":"clean"' "request-only lifecycle descendant preserves verification coverage"

echo "Scenario reconciliation is read-only and reports a closed Issue with live work"
RELEASE_PLAN="$WS/.spectacular/requests/release-check/PLAN.md"
sed -i.bak 's/status: planned/status: active/; s/source_type: goal/source_type: issue/; s/source_ref: "release-readiness"/source_ref: "acme\/widget#44"/' "$RELEASE_PLAN"; rm -f "$RELEASE_PLAN.bak"
json=$(cd "$WS" && PATH="$BIN:$PATH" "$CLI" github reconcile release-check --json)
assert_has "$json" '"status":"discrepancies"' "reconcile reports discrepancy state"
assert_has "$json" "source Issue is closed" "closed Issue/live request discrepancy"
grep -qE 'issue (close|edit)|pr merge' "$GH_LOG" && fail "reconcile must not mutate GitHub" || pass

echo ""
echo "Results: $pass_count passed, $fail_count failed"
rm -rf "$WS" "$REMOTE" "$BIN"; rm -f "$GH_LOG"
[[ $fail_count -eq 0 ]]
