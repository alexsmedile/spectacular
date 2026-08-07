#!/usr/bin/env bash
# Request namespace, compiled context, approved-spec handoff, and verify gate.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
WS="/tmp/spectacular-request-workflow"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_has() { [[ "$1" == *"$2"* ]] && pass || fail "$3: missing '$2'"; }

rm -rf "$WS"; mkdir -p "$WS"
(cd "$WS" && "$CLI" init --minimal --name request-test --skill-scope none >/dev/null)
(cd "$WS" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)

echo "Scenario approved specification creates one planned request bundle"
(cd "$WS" && "$CLI" spec new auth --summary "Authentication" >/dev/null)
SPEC="$WS/.spectacular/specs/SPC-001-auth.md"
SPEC="$WS/.spectacular/specs/auth.md"
awk '/^- _\(fill in\)_/ { print "- Accept signed login tokens."; print "- Reject expired login tokens."; next } { print }' "$SPEC" > "$SPEC.tmp"
mv "$SPEC.tmp" "$SPEC"
(cd "$WS" && git add . && git commit -qm spec)
(cd "$WS" && git switch -qc request/auth)
(cd "$WS" && "$CLI" request new auth-build --from auth >/dev/null)
PLAN="$WS/.spectacular/requests/auth-build/PLAN.md"
TASKS="$WS/.spectacular/requests/auth-build/TASKS.md"
[[ -f "$PLAN" && -f "$TASKS" ]] && pass || fail "PLAN and TASKS generated"
grep -Eq '^contract: [0-9a-f-]+$' "$PLAN" && pass || fail "immutable contract recorded"
grep -qF -- '- [ ] Accept signed login tokens.' "$TASKS" && pass || fail "requirements preserve approved order"
code=0; (cd "$WS" && "$CLI" request new duplicate --from s1 >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "duplicate source request refused"

echo "Scenario active brief is compiled and milestone aliases converge"
code=0; (cd "$WS" && "$CLI" request auth-build --brief >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "planned request has no implementation brief"
(cd "$WS" && "$CLI" request advance auth-build >/dev/null)
out1=$(cd "$WS" && "$CLI" request auth-build --brief --milestone M1)
out2=$(cd "$WS" && "$CLI" request auth-build --brief -m 1)
out3=$(cd "$WS" && "$CLI" request auth-build --brief -m1)
[[ "$out1" == "$out2" && "$out2" == "$out3" ]] && pass || fail "milestone aliases return identical brief"
assert_has "$out1" "Contract: " "brief contract provenance"
assert_has "$out1" "Execution boundary" "brief scope boundary"
grep -q '^activated_against:' "$PLAN" && pass || fail "activation git baseline recorded"
overview=$(cd "$WS" && "$CLI" request auth-build)
assert_has "$overview" "Task:" "overview current task"
json=$(cd "$WS" && "$CLI" request auth-build --brief -m1 --json)
assert_has "$json" '"schema":"spectacular.request.v1"' "brief JSON schema"
full=$(cd "$WS" && "$CLI" request auth-build --full)
[[ "$full" == *'<!-- PLAN.md -->'* && "$full" == *'<!-- TASKS.md -->'* && "$full" == *'<!-- SESSION.md -->'* ]] && pass || fail "full bundle stable core order"

echo "Scenario terminal act redirects and verify evidence owns the transition"
code=0; (cd "$WS" && "$CLI" spec act s1 >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "terminal spec act redirects"
code=0; (cd "$WS" && "$CLI" act s1 >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "short act redirects"
(cd "$WS" && "$CLI" request advance auth-build >/dev/null)
(cd "$WS" && "$CLI" docs-impact auth-build --none --reason "No separate docs impact" >/dev/null)
code=0; (cd "$WS" && "$CLI" request advance auth-build >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "verified transition refuses missing evidence"
printf '%s\n' '---' 'result: pass' '---' '# Verify log' > "$WS/.spectacular/requests/auth-build/VERIFY-LOG.md"
(cd "$WS" && "$CLI" request advance auth-build >/dev/null)
grep -q '^status: verified' "$PLAN" && pass || fail "passing evidence permits verified"

echo "Scenario decisions persist reusable tags"
(cd "$WS" && "$CLI" decide "Use tagged decisions" --tags cli,docs >/dev/null)
DECISION_FILE="$WS/.spectacular/decisions/use-tagged-decisions.md"
grep -q '^tags: \[cli, docs\]' "$DECISION_FILE" && pass || fail "decision tags persisted"
grep -Eq '^id: [0-9a-f-]+$' "$DECISION_FILE" && pass || fail "decision uses UUIDv7 identity"
for n in $(seq 2 59); do
  id=$(printf 'DEC-%03d' "$n")
  printf '%s\n' '---' "id: $id" 'type: decision' 'status: verified' 'tags: [test]' '---' '' "# $id — Test decision $n" '' '**Consequences:**' 'Test rationale.' > "$WS/.spectacular/decisions/$id-test.md"
done
(cd "$WS" && "$CLI" decide "Trigger tiered compaction" --tags indexing >/dev/null)
INDEX="$WS/.spectacular/decisions/index.md"
grep -qF -- 'Trigger tiered compaction' "$INDEX" && pass || fail "UUIDv7 decision remains indexed"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
rm -rf "$WS"
[[ $fail_count -eq 0 ]]
