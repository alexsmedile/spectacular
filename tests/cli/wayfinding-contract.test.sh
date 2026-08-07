#!/usr/bin/env bash
# UUIDv7 Wayfinding graph contract.
set -u
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"; CLI="$ROOT/cli/spectacular"; d=/tmp/spectacular-wayfinding-contract
p=0; f=0
pass(){ p=$((p+1)); }; fail(){ echo "    ✗ $1"; f=$((f+1)); }
id(){ awk '/^id:/{print $2;exit}' "$1"; }
rm -rf "$d"; mkdir -p "$d/.spectacular"
printf 'project:\n  name: wayfinding-test\n' > "$d/.spectacular/config.yaml"

echo "Scenario UUIDv7 graph: slugs create records and IDs persist dependencies"
(cd "$d" && "$CLI" question new tenant --question "Which tenant model?" >/dev/null)
qid=$(id "$d/.spectacular/questions/tenant.md")
(cd "$d" && "$CLI" research new isolation --summary "Evaluate isolation" --blocked-by "$qid" >/dev/null)
rid=$(id "$d/.spectacular/research/isolation.md")
(cd "$d" && "$CLI" spike new scale --summary "Test scale" --blocked-by "$qid" >/dev/null)
sid=$(id "$d/.spectacular/spikes/scale.md")
[[ "$qid" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-7 ]] && pass || fail "question has UUIDv7"
grep -qF -- "- $qid" "$d/.spectacular/research/isolation.md" && pass || fail "dependency persists UUIDv7"
out=$(cd "$d" && "$CLI" wayfind status)
[[ "$out" == *"$qid"* && "$out" == *"$rid"* ]] && pass || fail "status projects UUIDv7 graph"
(cd "$d" && "$CLI" wayfind resolve "$qid" --answer "Database per tenant" >/dev/null)
[[ -f "$d/.spectacular/archive/questions/tenant.md" ]] && pass || fail "resolved UUID question archives"
out=$(cd "$d" && "$CLI" wayfind next)
[[ "$out" == *"$sid"* || "$out" == *"$rid"* ]] && pass || fail "a formerly blocked UUIDv7 node becomes actionable"

echo "Scenario compatibility: a migrated numeric alias resolves but new records stay UUIDv7"
mkdir -p "$d/.spectacular/ideas"
printf '%s\n' '---' 'id: IDEA-001' 'status: parked' 'summary: legacy' '---' '# Legacy' > "$d/.spectacular/ideas/legacy.md"
(cd "$d" && "$CLI" id migrate --uuidv7 --apply --yes >/dev/null)
iid=$(id "$d/.spectacular/ideas/legacy.md")
[[ "$(cd "$d" && "$CLI" id resolve IDEA-001 --context idea)" == "$iid" ]] && pass || fail "legacy alias resolves"
(cd "$d" && "$CLI" idea new fresh >/dev/null)
[[ ! -f "$d/.spectacular/ideas/IDEA-002-fresh.md" && -f "$d/.spectacular/ideas/fresh.md" ]] && pass || fail "new idea uses slug path only"

echo "Results: $p passed, $f failed"; rm -rf "$d"; [[ $f -eq 0 ]]
