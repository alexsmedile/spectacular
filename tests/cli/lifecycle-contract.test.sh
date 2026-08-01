#!/usr/bin/env bash
# Canonical lifecycle contract: evidence gates, corrections, and archive-first history.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
WS="/tmp/spectacular-lifecycle-contract"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }

rm -rf "$WS"; mkdir -p "$WS/.spectacular/specs"
printf 'project:\n  name: lifecycle-test\nlast_build: 0\n' > "$WS/.spectacular/config.yaml"
printf '%s\n' '---' 'version: 1.0' 'updated: 2026-08-01' 'summary: "index"' '---' '# Index' > "$WS/.spectacular/specs/index.md"

echo "Scenario specs: approval, implementation evidence, atomic supersession, archive"
(cd "$WS" && "$CLI" spec new original --summary "Original behavior" >/dev/null)
(cd "$WS" && "$CLI" spec approve s1 --evidence "user chose it" >/dev/null)
(cd "$WS" && "$CLI" request new original-build --from s1 >/dev/null)
PLAN1="$WS/.spectacular/requests/original-build/PLAN.md"
sed -i.bak 's/status: planned/status: verified/' "$PLAN1"; rm -f "$PLAN1.bak"
(cd "$WS" && "$CLI" docs-impact original-build --none --reason "No public surface" >/dev/null)
(cd "$WS" && "$CLI" spec implement s1 --verified-against "commit abc123" >/dev/null)
SPEC1="$WS/.spectacular/specs/SPC-001-original.md"
assert_contains "$SPEC1" "status: implemented"
assert_contains "$SPEC1" "verified_against: commit abc123"

(cd "$WS" && "$CLI" spec new replacement --summary "Replacement behavior" --supersedes s1 >/dev/null)
(cd "$WS" && "$CLI" spec approve s2 --evidence "user approved replacement" >/dev/null)
(cd "$WS" && "$CLI" request new replacement-build --from s2 >/dev/null)
PLAN2="$WS/.spectacular/requests/replacement-build/PLAN.md"
sed -i.bak 's/status: planned/status: verified/' "$PLAN2"; rm -f "$PLAN2.bak"
(cd "$WS" && "$CLI" docs-impact replacement-build --required --evidence "docs/workflow.md updated" >/dev/null)
(cd "$WS" && "$CLI" spec implement s2 --verified-against "build b99" >/dev/null)
SPEC2="$WS/.spectacular/specs/SPC-002-replacement.md"
assert_contains "$SPEC1" "status: superseded"
assert_contains "$SPEC1" "superseded_by: SPC-002"
assert_contains "$SPEC2" "status: implemented"

(cd "$WS" && "$CLI" spec deprecate s2 --reason "Behavior removed" >/dev/null)
out=$(cd "$WS" && "$CLI" spec archive s2 --reason "No longer relevant")
[[ "$out" == *"Dry-run only"* && -f "$SPEC2" ]] && pass || fail "spec archive defaults to dry-run"
(cd "$WS" && "$CLI" spec archive s2 --reason "No longer relevant" --apply --yes >/dev/null)
[[ -f "$WS/.spectacular/archive/specs/SPC-002-replacement.md" && ! -f "$SPEC2" ]] && pass || fail "terminal spec moved archive-first"

echo "Scenario knowledge: inconclusive fog, question provenance, memory correction, verified fixes"
(cd "$WS" && "$CLI" research new uncertain --summary "Test uncertainty" >/dev/null)
(cd "$WS" && "$CLI" research resolve r1 --result inconclusive --outcome "Need narrower test" --evidence "report.md" >/dev/null)
out=$(cd "$WS" && "$CLI" wayfind status)
[[ "$out" == *"RES-001     fog"* ]] && pass || fail "inconclusive completion remains fog"

(cd "$WS" && "$CLI" question new choice --question "Choose?" >/dev/null)
(cd "$WS" && "$CLI" question resolve q1 --answer "A" --resolved-by user --source "interview 2026-08-01" >/dev/null)
ARCHIVED_Q="$WS/.spectacular/archive/questions/QUE-001-choice.md"
[[ -f "$ARCHIVED_Q" && ! -f "$WS/.spectacular/questions/QUE-001-choice.md" ]] && pass || fail "resolved question archived"
assert_contains "$ARCHIVED_Q" "archived_from: resolved"
assert_contains "$ARCHIVED_Q" "resolved_by: user"
assert_contains "$ARCHIVED_Q" "resolution_source: interview 2026-08-01"

echo "Scenario retention: implemented and abandoned specs leave active context archive-first"
(cd "$WS" && "$CLI" spec new integrated --summary "Integrated behavior" >/dev/null)
(cd "$WS" && "$CLI" spec approve s3 --evidence "approved" >/dev/null)
(cd "$WS" && "$CLI" request new integrated-build --from s3 >/dev/null)
PLAN3="$WS/.spectacular/requests/integrated-build/PLAN.md"
sed -i.bak 's/status: planned/status: verified/' "$PLAN3"; rm -f "$PLAN3.bak"
(cd "$WS" && "$CLI" docs-impact integrated-build --none --reason "No public surface" >/dev/null)
(cd "$WS" && "$CLI" spec implement s3 --verified-against "commit def456" >/dev/null)
out=$(cd "$WS" && "$CLI" spec archive s3 --reason "Merged; code and tests are authoritative")
[[ "$out" == *"Dry-run only"* && -f "$WS/.spectacular/specs/SPC-003-integrated.md" ]] && pass || fail "implemented spec archive defaults to dry-run"
(cd "$WS" && "$CLI" spec archive s3 --reason "Merged; code and tests are authoritative" --apply --yes >/dev/null)
ARCHIVED_S3="$WS/.spectacular/archive/specs/SPC-003-integrated.md"
[[ -f "$ARCHIVED_S3" && ! -f "$WS/.spectacular/specs/SPC-003-integrated.md" ]] && pass || fail "implemented spec archived out of live set"
assert_contains "$ARCHIVED_S3" "archived_from: implemented"
assert_contains "$ARCHIVED_S3" "verified_against: commit def456"

(cd "$WS" && "$CLI" spec new integrated-v2 --summary "Replacement for archived implementation" --supersedes s3 >/dev/null)
(cd "$WS" && "$CLI" spec approve s4 --evidence "approved replacement" >/dev/null)
(cd "$WS" && "$CLI" request new integrated-v2-build --from s4 >/dev/null)
PLAN4="$WS/.spectacular/requests/integrated-v2-build/PLAN.md"
sed -i.bak 's/status: planned/status: verified/' "$PLAN4"; rm -f "$PLAN4.bak"
(cd "$WS" && "$CLI" docs-impact integrated-v2-build --none --reason "No public surface" >/dev/null)
(cd "$WS" && "$CLI" spec implement s4 --verified-against "commit fed654" >/dev/null)
assert_contains "$ARCHIVED_S3" "archived_from: implemented"
assert_contains "$WS/.spectacular/specs/SPC-004-integrated-v2.md" "status: implemented"

(cd "$WS" && "$CLI" spec new rejected --summary "Rejected path" >/dev/null)
(cd "$WS" && "$CLI" spec archive s5 --reason "Failed feasibility" --apply --yes >/dev/null)
assert_contains "$WS/.spectacular/archive/specs/SPC-005-rejected.md" "archived_from: draft"

(cd "$WS" && "$CLI" remember "The old fact" >/dev/null)
(cd "$WS" && "$CLI" remember "The corrected fact" --supersedes M1 >/dev/null)
assert_contains "$WS/.spectacular/memories/M1-the-old-fact.md" "status: superseded"
(cd "$WS" && "$CLI" memory retract M2 --reason "Source disproved" >/dev/null)
assert_contains "$WS/.spectacular/memories/M2-the-corrected-fact.md" "status: retracted"

code=0; (cd "$WS" && "$CLI" fix new "Unverified fix" >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "fix requires verification evidence"
(cd "$WS" && "$CLI" fix new "Verified fix" --verified-by "tests pass" >/dev/null)
[[ -f "$WS/.spectacular/fixes/F1-verified-fix.md" ]] && pass || fail "verified fix recorded"

echo "Scenario migration: singular collections stay readable and migrate archive-first"
mkdir -p "$WS/.spectacular/debug"
printf '%s\n' '---' 'status: resolved' '---' '# Legacy debug' > "$WS/.spectacular/debug/legacy.md"
out=$(cd "$WS" && "$CLI" collections migrate)
[[ "$out" == *"Dry-run only"* && -f "$WS/.spectacular/debug/legacy.md" ]] && pass || fail "collection migration previews without mutation"
(cd "$WS" && "$CLI" collections migrate --apply --yes >/dev/null)
[[ -f "$WS/.spectacular/debugs/legacy.md" ]] && pass || fail "legacy record moved to canonical plural collection"
archived_count=$(find "$WS/.spectacular/archive/collection-migrations" -path '*/debug/legacy.md' | wc -l | tr -d ' ')
[[ "$archived_count" -eq 1 ]] && pass || fail "legacy collection archived before movement"

echo "Scenario lifecycle migration: safe labels map only after explicit apply"
printf '%s\n' '---' 'id: SPC-009' 'type: specification' 'status: current' 'updated: 2026-01-01' 'summary: "legacy"' 'related: []' '---' '# Legacy spec' > "$WS/.spectacular/specs/SPC-009-legacy.md"
printf '%s\n' '---' 'type: memory' 'date: 2026-01-01' 'summary: "legacy memory"' '---' 'Legacy memory' > "$WS/.spectacular/memories/M9-legacy-memory.md"
out=$(cd "$WS" && "$CLI" lifecycle migrate)
[[ "$out" == *"SPC-009-legacy.md: current → approved"* && "$out" == *"M9-legacy-memory.md: missing → active"* ]] && pass || fail "lifecycle migration previews exact safe mappings"
assert_contains "$WS/.spectacular/specs/SPC-009-legacy.md" "status: current"
(cd "$WS" && "$CLI" lifecycle migrate --apply --yes >/dev/null)
assert_contains "$WS/.spectacular/specs/SPC-009-legacy.md" "status: approved"
assert_contains "$WS/.spectacular/memories/M9-legacy-memory.md" "status: active"
archived_count=$(find "$WS/.spectacular/archive/lifecycle-migrations" -path '*/specs/SPC-009-legacy.md' | wc -l | tr -d ' ')
[[ "$archived_count" -eq 1 ]] && pass || fail "lifecycle migration archives originals"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
rm -rf "$WS"
[[ $fail_count -eq 0 ]]
