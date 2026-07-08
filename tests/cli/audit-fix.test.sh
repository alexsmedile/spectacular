#!/usr/bin/env bash
# tests/cli/audit-fix.test.sh — spectacular audit + fix soft-DB collections (v1.25.0)
#
# Covers the real logic: auto-numbering (A<N>/F<N>), the verified gate warning,
# from-audit validation, and the --into-fix graduation (cause copy + disposition).

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_eq()   { [[ "$1" == "$2" ]] && pass || fail "$3: got '$1' want '$2'"; }
assert_file() { [[ -f "$1" ]] && pass || fail "expected file: $1"; }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }

new_ws() { local d="$1"; rm -rf "$d"; mkdir -p "$d/.spectacular"; printf 'name: t\n' > "$d/.spectacular/config.yaml"; }

# ── auto-numbering: A1, A2, A3 in sequence; F likewise ────────────────────────
scenario_autonumber() {
  echo "Scenario autonumber: sequential A<N>/F<N> ids"
  local d="/tmp/spec-af-num"; new_ws "$d"
  ( cd "$d"
    "$CLI" audit new "one"   >/dev/null
    "$CLI" audit new "two"   >/dev/null
    "$CLI" audit new "three" >/dev/null )
  assert_file "$d/.spectacular/audits/A1-one.md"
  assert_file "$d/.spectacular/audits/A2-two.md"
  assert_file "$d/.spectacular/audits/A3-three.md"
  rm -rf "$d"
}

# ── fix verified gate: no --verified-by → verified: null + nonzero-safe exit ───
scenario_verified_gate() {
  echo "Scenario verified-gate: unverified fix logs with verified: null"
  local d="/tmp/spec-af-gate"; new_ws "$d"
  local code
  ( cd "$d" && "$CLI" fix new "unproven" >/dev/null 2>&1 ) && code=0 || code=$?
  assert_exit "$code" 0 "fix new without --verified-by still exits 0"
  assert_contains "$d/.spectacular/fixes/F1-unproven.md" "verified: null"
  ( cd "$d" && "$CLI" fix new "proven" --verified-by "tests/x.sh" >/dev/null 2>&1 )
  grep -q "verified: null" "$d/.spectacular/fixes/F2-proven.md" && fail "F2 should be verified" || pass
  rm -rf "$d"
}

# ── from-audit validation: pointing at a missing audit errors ─────────────────
scenario_from_audit_validation() {
  echo "Scenario from-audit: missing audit ref is rejected"
  local d="/tmp/spec-af-fa"; new_ws "$d"
  local code
  ( cd "$d" && "$CLI" fix new "x" --from-audit A99 --verified-by "y" >/dev/null 2>&1 ) && code=0 || code=$?
  assert_exit "$code" 1 "fix new --from-audit A99 (nonexistent) exits 1"
  rm -rf "$d"
}

# ── skeleton: entries carry problem → intended → cause → fix → criteria ───────
scenario_skeleton() {
  echo "Scenario skeleton: fix entry carries the full bug-fixing skeleton"
  local d="/tmp/spec-af-skel"; new_ws "$d"
  ( cd "$d" && "$CLI" fix new "skel" --problem P --intended I --cause C --fix X \
      --criteria K --verified-by V --signature "sig words" >/dev/null 2>&1 )
  local file="$d/.spectacular/fixes/F1-skel.md"
  assert_contains "$file" "## Problem"
  assert_contains "$file" "## Intended behavior"
  assert_contains "$file" "## Success criteria"
  assert_contains "$file" "## Signature"
  assert_contains "$file" "signature: sig words"    # retrievable via frontmatter
  rm -rf "$d"
}

# ── into-fix graduation: closes audit + scaffolds F, copies ALL matching slots ─
scenario_into_fix() {
  echo "Scenario into-fix: audit resolve --into-fix graduates + copies slots"
  local d="/tmp/spec-af-into"; new_ws "$d"
  ( cd "$d" && "$CLI" audit new "bug title" --problem "the problem" --intended "the intent" >/dev/null )
  # fill the audit's root cause so the copy-forward is exercised
  sed -i.bak 's|_(the actual cause, once found — or "not yet found")_|the real cause|' "$d/.spectacular/audits/A1-bug-title.md"
  rm -f "$d/.spectacular/audits/A1-bug-title.md.bak"
  ( cd "$d" && "$CLI" audit resolve A1 --disposition "user phrasing" --into-fix --verified-by "repro" >/dev/null )
  assert_file "$d/.spectacular/fixes/F1-bug-title.md"
  assert_contains "$d/.spectacular/audits/A1-bug-title.md" "status: resolved"
  # disposition owned by --into-fix (no doubling of the user's phrasing)
  assert_contains "$d/.spectacular/audits/A1-bug-title.md" "disposition: became fix F1"
  assert_contains "$d/.spectacular/fixes/F1-bug-title.md"  "from_audit: A1"
  # slots copied forward from the audit
  assert_contains "$d/.spectacular/fixes/F1-bug-title.md"  "the problem"
  assert_contains "$d/.spectacular/fixes/F1-bug-title.md"  "the intent"
  assert_contains "$d/.spectacular/fixes/F1-bug-title.md"  "the real cause"
  rm -rf "$d"
}

scenario_autonumber
scenario_verified_gate
scenario_from_audit_validation
scenario_skeleton
scenario_into_fix

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
