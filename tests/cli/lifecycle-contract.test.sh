#!/usr/bin/env bash
# UUIDv7 contract lifecycle: merged contract ancestry, archive history, aliases.

set -u
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
WS="/tmp/spectacular-lifecycle-contract"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
id_of() { awk '/^id:/{print $2; exit}' "$1"; }

rm -rf "$WS"; mkdir -p "$WS/.spectacular/specs"
printf 'project:\n  name: lifecycle-test\nlast_build: 0\nforge:\n  shared_base: "main"\n' > "$WS/.spectacular/config.yaml"
printf '%s\n' '---' 'version: 1.0' 'updated: 2026-08-01' 'summary: "index"' '---' '# Index' > "$WS/.spectacular/specs/index.md"
(cd "$WS" && git init -q && git branch -M main && git config user.email test@example.com && git config user.name Test && git add . && git commit -qm init)

echo "Scenario contracts: a merged spec commit, not mutable approval metadata, authorizes execution"
(cd "$WS" && "$CLI" spec new original --summary "Original behavior" >/dev/null)
SPEC1="$WS/.spectacular/specs/original.md"; ID1=$(id_of "$SPEC1")
[[ "$ID1" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-7 ]] && pass || fail "spec receives UUIDv7"
(cd "$WS" && git add . && git commit -qm contract-original && git switch -q -c request/original)
(cd "$WS" && "$CLI" request new original-build --from "$ID1" >/dev/null)
PLAN1="$WS/.spectacular/requests/original-build/PLAN.md"
assert_contains "$PLAN1" "contract: $ID1"
grep -q '^source_spec:' "$PLAN1" && fail "request omits legacy provenance" || pass
sed -i.bak 's/status: planned/status: verified/' "$PLAN1"; rm -f "$PLAN1.bak"
(cd "$WS" && "$CLI" docs-impact original-build --none --reason "No public surface" >/dev/null)
(cd "$WS" && "$CLI" spec implement "$ID1" --verified-against "commit abc123" >/dev/null)
assert_contains "$SPEC1" "status: implemented"

echo "Scenario ancestry: unmerged or missing-base contracts cannot create requests"
(cd "$WS" && git switch -q main && "$CLI" spec new pending --summary "Pending behavior" >/dev/null)
PENDING="$WS/.spectacular/specs/pending.md"; PENDING_ID=$(id_of "$PENDING")
(cd "$WS" && git switch -q request/original)
code=0; (cd "$WS" && "$CLI" request new rejected-build --from "$PENDING_ID" >/dev/null 2>&1) || code=$?
assert_exit "$code" 1 "unmerged contract rejected"

echo "Scenario retention: slug paths archive without changing immutable IDs"
(cd "$WS" && "$CLI" spec deprecate "$ID1" --reason "Behavior removed" >/dev/null)
out=$(cd "$WS" && "$CLI" spec archive "$ID1" --reason "No longer relevant")
[[ "$out" == *"Dry-run only"* && -f "$SPEC1" ]] && pass || fail "spec archive is preview-first"
(cd "$WS" && "$CLI" spec archive "$ID1" --reason "No longer relevant" --apply --yes >/dev/null)
ARCHIVED_SPEC="$WS/.spectacular/archive/specs/original.md"
[[ -f "$ARCHIVED_SPEC" && ! -f "$SPEC1" ]] && pass || fail "slug-path spec archived"
assert_contains "$ARCHIVED_SPEC" "id: $ID1"

echo "Scenario knowledge: UUIDv7 entities resolve by slug and archive safely"
(cd "$WS" && "$CLI" research new uncertain --summary "Test uncertainty" >/dev/null)
RID=$(id_of "$WS/.spectacular/research/uncertain.md")
(cd "$WS" && "$CLI" research resolve "$RID" --result inconclusive --outcome "Need narrower test" --evidence report.md >/dev/null)
out=$(cd "$WS" && "$CLI" wayfind status)
[[ "$out" == *"$RID"* && "$out" == *"fog"* ]] && pass || fail "inconclusive UUID research remains fog"
(cd "$WS" && "$CLI" question new choice --question "Choose?" >/dev/null)
QID=$(id_of "$WS/.spectacular/questions/choice.md")
(cd "$WS" && "$CLI" question resolve choice --answer A --resolved-by user --source interview >/dev/null)
ARCHIVED_Q="$WS/.spectacular/archive/questions/choice.md"
[[ -f "$ARCHIVED_Q" && ! -f "$WS/.spectacular/questions/choice.md" ]] && pass || fail "slug question archives"
assert_contains "$ARCHIVED_Q" "id: $QID"

echo "Scenario migration: preview has no writes; apply preserves aliases and rewrites references"
mkdir -p "$WS/.spectacular/ideas"
printf '%s\n' '---' 'id: IDEA-001' 'status: parked' 'summary: "Legacy idea"' '---' '# Legacy idea' > "$WS/.spectacular/ideas/legacy-idea.md"
printf '%s\n' '# Reference' 'IDEA-001' 'ideas/legacy-idea.md' > "$WS/.spectacular/LEGACY-REFERENCE.md"
out=$(cd "$WS" && "$CLI" id migrate --uuidv7)
[[ "$out" == *"Dry-run only"* && -f "$WS/.spectacular/ideas/legacy-idea.md" ]] && pass || fail "UUID migration previews without writes"
(cd "$WS" && "$CLI" id migrate --uuidv7 --apply --yes >/dev/null)
MIGRATED="$WS/.spectacular/ideas/legacy-idea.md"; MID=$(id_of "$MIGRATED")
[[ "$MID" =~ ^[0-9a-f]{8}-[0-9a-f]{4}-7 ]] && pass || fail "legacy record obtains UUIDv7"
assert_contains "$MIGRATED" "aliases: IDEA-001"
out=$(cd "$WS" && "$CLI" id resolve IDEA-001 --context idea)
[[ "$out" == "$MID" ]] && pass || fail "legacy alias resolves to UUIDv7"
assert_contains "$WS/.spectacular/LEGACY-REFERENCE.md" "$MID"
find "$WS/.spectacular/archive/id-migrations" -name uuidv7-mapping.tsv -type f | grep -q . && pass || fail "migration mapping receipt archived"

echo ""
echo "Results: $pass_count passed, $fail_count failed"
rm -rf "$WS"
[[ $fail_count -eq 0 ]]
