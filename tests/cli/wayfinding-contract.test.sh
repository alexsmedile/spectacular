#!/usr/bin/env bash
# Canonical ID resolver + questions soft-DB substrate.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_eq() { [[ "$1" == "$2" ]] && pass || fail "$3: got '$1' want '$2'"; }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }

new_ws() {
  local d="$1"
  rm -rf "$d"
  mkdir -p "$d/.spectacular"
  printf 'project:\n  name: wayfinding-test\n' > "$d/.spectacular/config.yaml"
}

scenario_aliases() {
  echo "Scenario aliases: canonical, shorthand, and contextual numbers normalize"
  assert_eq "$($CLI id resolve d1)" "DEC-001" "d1"
  assert_eq "$($CLI id resolve res2)" "RES-002" "res2"
  assert_eq "$($CLI id resolve IDE-42)" "IDEA-042" "legacy IDE alias"
  assert_eq "$($CLI id resolve spec7)" "SPC-007" "reserved SPEC alias"
  assert_eq "$($CLI id resolve fnd1)" "FND-001" "finding-safe alias"
  assert_eq "$($CLI id resolve fix1)" "F1" "explicit fix alias resolves legacy fix"
  assert_eq "$($CLI id resolve bug1)" "BUG-001" "bug-safe alias"
  assert_eq "$($CLI id resolve b1)" "b1" "roadmap build owns b1"
  assert_eq "$($CLI id resolve 3 --context question)" "QUE-003" "contextual number"

  local code
  "$CLI" id resolve 3 >/dev/null 2>&1 && code=0 || code=$?
  assert_exit "$code" 1 "ambiguous naked number"
  "$CLI" id resolve f1 >/dev/null 2>&1 && code=0 || code=$?
  assert_exit "$code" 1 "ambiguous f1 is refused"
}

scenario_questions() {
  echo "Scenario questions: allocate, canonicalize dependencies, list, and resolve"
  local d="/tmp/spectacular-wayfinding-questions"
  new_ws "$d"
  (cd "$d" && "$CLI" question new tenant-isolation \
    --question "Which tenant isolation model?" \
    --context "Blocks the current spec." \
    --blocked-by res1,spk2 >/dev/null)
  (cd "$d" && "$CLI" question new billing-owner --question "Who owns billing?" >/dev/null)

  local q1="$d/.spectacular/questions/QUE-001-tenant-isolation.md"
  local q2="$d/.spectacular/questions/QUE-002-billing-owner.md"
  [[ -f "$q1" && -f "$q2" ]] && pass || fail "QUE-001/QUE-002 files written"
  assert_contains "$q1" "id: QUE-001"
  assert_contains "$q1" "  - RES-001"
  assert_contains "$q1" "  - SPK-002"
  assert_contains "$q1" "requires_user_input: true"

  local out
  out=$(cd "$d" && "$CLI" question list)
  [[ "$out" == *"QUE-001"* && "$out" == *"QUE-002"* ]] && pass || fail "question list shows both IDs"

  (cd "$d" && "$CLI" question resolve q1 --answer "Use database-per-tenant." >/dev/null)
  local archived_q1="$d/.spectacular/archive/questions/QUE-001-tenant-isolation.md"
  [[ ! -f "$q1" && -f "$archived_q1" ]] && pass || fail "resolved question leaves active collection and archives"
  assert_contains "$archived_q1" "status: archived"
  assert_contains "$archived_q1" "archived_from: resolved"
  assert_contains "$archived_q1" "Use database-per-tenant."
  rm -rf "$d"
}

scenario_typed_records() {
  echo "Scenario typed records: canonical decisions, ideas, research, and spikes"
  local d="/tmp/spectacular-wayfinding-records"
  new_ws "$d"
  (cd "$d" && "$CLI" decide "Use the graph resolver" >/dev/null)
  (cd "$d" && "$CLI" idea new voice-command-bar >/dev/null)
  (cd "$d" && "$CLI" research new auth-options --summary "Compare auth options" --blocked-by q1 >/dev/null)
  (cd "$d" && "$CLI" spike new auth-scale --summary "Test auth scale" >/dev/null)

  [[ -f "$d/.spectacular/decisions/DEC-001-use-the-graph-resolver.md" ]] && pass || fail "canonical decision file"
  [[ -f "$d/.spectacular/ideas/IDEA-001-voice-command-bar.md" ]] && pass || fail "canonical idea file"
  assert_contains "$d/.spectacular/ideas/IDEA-001-voice-command-bar.md" "## Working plan"
  [[ -f "$d/.spectacular/research/RES-001-auth-options.md" ]] && pass || fail "canonical research file"
  [[ -f "$d/.spectacular/spikes/SPK-001-auth-scale.md" ]] && pass || fail "canonical spike file"
  assert_contains "$d/.spectacular/research/RES-001-auth-options.md" "  - QUE-001"
  assert_contains "$d/.spectacular/spikes/SPK-001-auth-scale.md" "execution_requires_approval: true"

  local code
  (cd "$d" && "$CLI" decide "Autonomous without evidence" --autonomous >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "autonomous ADR evidence gate"
  (cd "$d" && "$CLI" research resolve r1 --result supported --outcome "Option A wins" --evidence "bench.md" >/dev/null)
  (cd "$d" && "$CLI" afk run start choose-auth --goal "Choose reversible auth mechanism" --allowed-actions "research,decisions" --apply --yes >/dev/null)
  (cd "$d" && "$CLI" decide "Use option A" --autonomous --technical --in-scope --reversible --no-product-tradeoff --alternatives "Option B" --evidence "RES-001" >/dev/null)
  assert_contains "$d/.spectacular/decisions/DEC-002-use-option-a.md" "origin: autonomous-evidence"
  assert_contains "$d/.spectacular/decisions/DEC-002-use-option-a.md" "evidence: \"RES-001\""

  assert_contains "$d/.spectacular/research/RES-001-auth-options.md" "status: completed"
  assert_contains "$d/.spectacular/research/RES-001-auth-options.md" "result: supported"
  assert_contains "$d/.spectacular/research/RES-001-auth-options.md" "**Evidence:** bench.md"
  rm -rf "$d"
}

scenario_spec_lifecycle() {
  echo "Scenario specs: draft blocks action; approved spec seeds a request"
  local d="/tmp/spectacular-wayfinding-specs" code
  new_ws "$d"
  mkdir -p "$d/.spectacular/specs"
  printf '%s\n' '---' 'version: 1.0' 'updated: 2026-08-01' 'summary: "index"' '---' '# Index' > "$d/.spectacular/specs/index.md"
  printf '\nlast_build: 0\n' >> "$d/.spectacular/config.yaml"
  (cd "$d" && "$CLI" spec new onboarding --summary "User onboarding" --target-version v1.0.0-discovery >/dev/null)
  local spec="$d/.spectacular/specs/SPC-001-onboarding.md"
  [[ -f "$spec" ]] && pass || fail "SPC-001 spec written"
  assert_contains "$spec" "status: draft"

  (cd "$d" && "$CLI" request new onboarding-work --from s1 >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "draft spec action gate"
  (cd "$d" && "$CLI" spec approve s1 --evidence "user approved" --target-version v1.0.0-execution >/dev/null)
  assert_contains "$spec" "status: approved"
  assert_contains "$spec" "version: 1.0"
  [[ -f "$d/.spectacular/_snapshots/specs/SPC-001-onboarding/@v1.md" ]] && pass || fail "confirmation snapshot written"

  (cd "$d" && "$CLI" request new onboarding-work --from s1 >/dev/null)
  assert_contains "$d/.spectacular/requests/onboarding-work/PLAN.md" "source_spec: SPC-001"
  rm -rf "$d"
}

scenario_frontier() {
  echo "Scenario frontier: dependencies derive fog; session surfaces human blockers"
  local d="/tmp/spectacular-wayfinding-frontier" out
  new_ws "$d"
  (cd "$d" && "$CLI" question new product-choice --question "Which market?" >/dev/null)
  (cd "$d" && "$CLI" research new market-evidence --summary "Research the market" --blocked-by q1 >/dev/null)
  (cd "$d" && "$CLI" spike new feasibility --summary "Test feasibility" >/dev/null)

  out=$(cd "$d" && "$CLI" wayfind status)
  [[ "$out" == *"QUE-001     frontier"* ]] && pass || fail "open unblocked question is frontier"
  [[ "$out" == *"RES-001     fog"* ]] && pass || fail "research blocked by open question is fog"
  [[ "$out" == *"SPK-001     frontier"* ]] && pass || fail "unblocked spike is frontier"

  out=$(cd "$d" && "$CLI" wayfind next)
  [[ "$out" == *"QUE-001"* && "$out" == *"requires user input"* ]] && pass || fail "human blocker wins next action"
  out=$(cd "$d" && "$CLI" session start --notes "frontier test")
  [[ "$out" == *"Questions requiring your input"* && "$out" == *"QUE-001"* ]] && pass || fail "session start surfaces question"

  (cd "$d" && "$CLI" question resolve q1 --answer "Enterprise first" >/dev/null)
  [[ ! -f "$d/.spectacular/questions/QUE-001-product-choice.md" && -f "$d/.spectacular/archive/questions/QUE-001-product-choice.md" ]] && pass || fail "resolved blocker archived outside active fog"
  out=$(cd "$d" && "$CLI" wayfind status)
  [[ "$out" == *"RES-001     frontier"* ]] && pass || fail "resolved dependency promotes research to frontier"
  (cd "$d" && "$CLI" wayfind defer r1 --reason "not relevant today" >/dev/null)
  out=$(cd "$d" && "$CLI" wayfind status)
  [[ "$out" == *"RES-001     deferred"* ]] && pass || fail "deferred loop leaves active frontier"
  rm -rf "$d"
}

scenario_doctor_dag() {
  echo "Scenario doctor: valid DAG passes; dangling links and cycles are reported"
  local d="/tmp/spectacular-wayfinding-doctor" out code
  new_ws "$d"
  (cd "$d" && "$CLI" question new first --question "First?" >/dev/null)
  (cd "$d" && "$CLI" research new second --summary "Second" --blocked-by q1 >/dev/null)
  (cd "$d" && "$CLI" doctor wayfinding >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "valid wayfinding DAG"

  sed -i.bak 's/QUE-001/QUE-999/' "$d/.spectacular/research/RES-001-second.md"; rm -f "$d/.spectacular/research/RES-001-second.md.bak"
  out=$(cd "$d" && "$CLI" doctor wayfinding 2>&1); code=$?
  assert_exit "$code" 1 "dangling dependency warning"
  [[ "$out" == *"dangling dependency"* && "$out" == *"QUE-999"* ]] && pass || fail "dangling ID named"

  sed -i.bak 's/QUE-999/QUE-001/' "$d/.spectacular/research/RES-001-second.md"; rm -f "$d/.spectacular/research/RES-001-second.md.bak"
  sed -i.bak 's/blocked_by: \[\]/blocked_by:\n  - RES-001/' "$d/.spectacular/questions/QUE-001-first.md"; rm -f "$d/.spectacular/questions/QUE-001-first.md.bak"
  out=$(cd "$d" && "$CLI" doctor wayfinding 2>&1); code=$?
  assert_exit "$code" 2 "dependency cycle error"
  [[ "$out" == *"dependency cycle"* ]] && pass || fail "cycle finding emitted"
  rm -rf "$d"
}

scenario_doctor_duplicate_spec_id() {
  echo "Scenario doctor: duplicate spec IDs are rejected"
  local d="/tmp/spectacular-wayfinding-duplicate-spec" out code
  new_ws "$d"
  mkdir -p "$d/.spectacular/specs"
  printf '%s\n' '---' 'id: SPC-001' 'type: specification' 'status: draft' '---' '# First' > "$d/.spectacular/specs/SPC-001-first.md"
  printf '%s\n' '---' 'id: SPC-001' 'type: specification' 'status: draft' '---' '# Second' > "$d/.spectacular/specs/SPC-002-second.md"
  out=$(cd "$d" && "$CLI" doctor wayfinding 2>&1); code=$?
  assert_exit "$code" 2 "duplicate spec ID is an error"
  [[ "$out" == *"duplicate canonical ID(s): SPC-001"* ]] && pass || fail "duplicate spec ID named"
  rm -rf "$d"
}

scenario_id_migration() {
  echo "Scenario migration: dry-run is inert; apply archives originals and rewrites IDs"
  local d="/tmp/spectacular-wayfinding-migrate" out code
  new_ws "$d"
  mkdir -p "$d/.spectacular/decisions" "$d/.spectacular/ideas" "$d/.spectacular/specs"
  printf '%s\n' '---' 'type: decisions' '---' '# Decisions' '- **D1** — Old' > "$d/.spectacular/decisions/index.md"
  printf '%s\n' '# D1 — Old' '**Decision:** yes' > "$d/.spectacular/decisions/D1-old.md"
  printf '%s\n' '---' 'type: idea' 'status: parked' 'updated: 2026-01-01' '---' '# Idea' > "$d/.spectacular/ideas/foo.md"
  printf '%s\n' '---' 'status: draft' 'updated: 2026-01-01' 'summary: "cap"' 'related: []' '---' '# Cap' > "$d/.spectacular/specs/cap.md"

  out=$(cd "$d" && "$CLI" id migrate)
  [[ "$out" == *"Dry-run only"* ]] && pass || fail "migration defaults to dry-run"
  [[ -f "$d/.spectacular/decisions/D1-old.md" && ! -f "$d/.spectacular/decisions/DEC-001-old.md" ]] && pass || fail "dry-run changes nothing"
  (cd "$d" && "$CLI" id migrate --apply >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "apply requires explicit yes"

  (cd "$d" && "$CLI" id migrate --apply --yes >/dev/null)
  [[ -f "$d/.spectacular/decisions/DEC-001-old.md" ]] && pass || fail "decision migrated"
  [[ -f "$d/.spectacular/ideas/IDEA-001-foo.md" ]] && pass || fail "idea migrated"
  [[ -f "$d/.spectacular/specs/SPC-001-cap.md" ]] && pass || fail "spec migrated"
  assert_contains "$d/.spectacular/decisions/index.md" "**DEC-001**"
  assert_contains "$d/.spectacular/specs/SPC-001-cap.md" "status: draft"
  local archived_count
  archived_count=$(find "$d/.spectacular/archive/id-migrations" -type f | wc -l | tr -d ' ')
  [[ "$archived_count" -ge 4 ]] && pass || fail "originals and mapping archived"
  rm -rf "$d"
}

scenario_aliases
scenario_questions
scenario_typed_records
scenario_spec_lifecycle
scenario_frontier
scenario_doctor_dag
scenario_doctor_duplicate_spec_id
scenario_id_migration

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
