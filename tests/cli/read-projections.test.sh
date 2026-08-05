#!/usr/bin/env bash
# tests/cli/read-projections.test.sh — compact entity details retain decision-critical fields

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_output() { printf '%s' "$1" | grep -qF -- "$2" && pass || fail "expected output: $2"; }
assert_absent() { printf '%s' "$1" | grep -qF -- "$2" && fail "unexpected output: $2" || pass; }

new_ws() {
  local d="$1"
  rm -rf "$d"; mkdir -p "$d/.spectacular"
  printf 'name: read-projections\n' > "$d/.spectacular/config.yaml"
}

scenario_details_are_entity_specific() {
  echo "Scenario projections: named details retain fields without body-scanning lists"
  local d="/tmp/spectacular-read-projections" out
  new_ws "$d"
  (
    cd "$d" || exit 1
    "$CLI" decide "Use named projections" --context "Selection needs bounded signals" --consequences "Bodies remain evidence" --tags cli,read >/dev/null
    "$CLI" remember "Never make a list read a long evidence-only paragraph." --tag cli >/dev/null
    "$CLI" question new release-owner --question "Who owns the release?" --context "LONG-CONTEXT-ONLY" --priority high --blocked-by DEC-001 >/dev/null
    "$CLI" research new provider-limit --summary "Check provider limit" --blocked-by QUE-001 >/dev/null
    "$CLI" research resolve RES-001 --result inconclusive --outcome "Limit remains unknown" --evidence "provider-doc-v2" >/dev/null
    "$CLI" spike new parser-prototype --summary "Prototype parser" --blocked-by RES-001 >/dev/null
    "$CLI" spike resolve SPK-001 --result supported --outcome "Parser works" --evidence "prototype-log" >/dev/null
    "$CLI" idea new quick-import --hypothesis "Import is viable" --origin "support call" --priority high >/dev/null
    "$CLI" audit new "Config regression" --severity high --problem "Config is dropped" --intended "Config persists" >/dev/null
    sed -i.bak 's|_(the actual cause, once found — or "not yet found")_|Wrong key|' .spectacular/audits/A1-config-regression.md
    sed -i.bak 's|_(the suggested fix — still a proposal at audit stage, not yet applied/verified)_|Use canonical key|' .spectacular/audits/A1-config-regression.md
    rm -f .spectacular/audits/A1-config-regression.md.bak
    "$CLI" fix new "Config regression" --cause "Wrong key" --fix "Use canonical key" --verified-by "tests/config.test.sh" --signature "config-key" >/dev/null
  )

  out=$(cd "$d" && "$CLI" decision DEC-001)
  assert_output "$out" "Use named projections"
  assert_output "$out" "Bodies remain evidence"
  assert_output "$out" "Full evidence: spectacular decision DEC-001 --full"
  out=$(cd "$d" && "$CLI" decision DEC-001 --json)
  assert_output "$out" '"schema":"spectacular.decision.v1"'
  assert_output "$out" '"decision":"Use named projections"'

  out=$(cd "$d" && "$CLI" memory M1)
  assert_output "$out" "Never make a list"
  assert_output "$out" "Status: active"

  out=$(cd "$d" && "$CLI" question QUE-001)
  assert_output "$out" "Requires user input: true"
  assert_output "$out" "Blocked by:"
  out=$(cd "$d" && "$CLI" question list)
  assert_output "$out" "Who owns the release?"
  assert_absent "$out" "LONG-CONTEXT-ONLY"

  out=$(cd "$d" && "$CLI" research RES-001 --json)
  assert_output "$out" '"result":"inconclusive"'
  assert_output "$out" '"evidence":"provider-doc-v2"'
  out=$(cd "$d" && "$CLI" spike SPK-001)
  assert_output "$out" "Execution requires approval: true"
  assert_output "$out" "Parser works"

  out=$(cd "$d" && "$CLI" idea IDEA-001)
  assert_output "$out" "Import is viable"
  assert_output "$out" "support call"
  out=$(cd "$d" && "$CLI" audit A1 --json)
  assert_output "$out" '"root_cause":"Wrong key"'
  assert_output "$out" '"proposed_fix":"Use canonical key"'
  out=$(cd "$d" && "$CLI" fix F1)
  assert_output "$out" "Signature: config-key"
  assert_output "$out" "Verified by"
  out=$(cd "$d" && "$CLI" fix list --since 1d)
  assert_output "$out" "F1"
  out=$(cd "$d" && "$CLI" fix F1 --full)
  assert_output "$out" "## Success criteria"
  rm -rf "$d"
}

scenario_details_are_entity_specific
echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
