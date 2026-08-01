#!/usr/bin/env bash
# Strict DAG sequencing, metaphor routing, and cross-layer coherence.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
pass_count=0; fail_count=0
pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }
assert_exit() { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_contains() { grep -qF -- "$2" "$1" && pass || fail "$1 should contain: $2"; }

new_ws() {
  local d="$1"
  rm -rf "$d"; mkdir -p "$d/.spectacular/specs" "$d/.spectacular/requests"
  printf 'project:\n  name: sequencer-test\nlast_build: 0\n' > "$d/.spectacular/config.yaml"
  printf '%s\n' '---' 'version: 1.0' 'updated: 2026-08-01' 'summary: "index"' '---' '# Index' > "$d/.spectacular/specs/index.md"
}

scenario_topology_and_ranking() {
  echo "Scenario topology: dependencies are strict; priority then uncertainty select next"
  local d="/tmp/spectacular-wayfinding-sequencer-topology" out code order
  new_ws "$d"
  (cd "$d" && "$CLI" question new product-choice --question "Which market?" --priority medium >/dev/null)
  (cd "$d" && "$CLI" research new market-data --summary "Research market" --blocked-by q1 >/dev/null)
  (cd "$d" && "$CLI" spike new feasibility --summary "Test feasibility" --blocked-by q1 >/dev/null)
  (cd "$d" && "$CLI" spec new launch --summary "Launch" --target-version v1.0.0-execution >/dev/null)
  sed -i.bak 's/updated: 2026-08-01/blocked_by:\n  - RES-001\n  - SPK-001\nupdated: 2026-08-01/' "$d/.spectacular/specs/SPC-001-launch.md"; rm -f "$d/.spectacular/specs/SPC-001-launch.md.bak"

  order=$(cd "$d" && "$CLI" wayfind order)
  [[ "$order" == $'QUE-001\tfrontier\nRES-001\tfog\nSPK-001\tfog\nSPC-001\tfog' ]] && pass || fail "strict dependency-first order: $order"
  out=$(cd "$d" && "$CLI" wayfind next)
  [[ "$out" == *"QUE-001"* ]] && pass || fail "human question wins equal-priority uncertainty tie"
  (cd "$d" && "$CLI" question resolve q1 --answer "Enterprise" >/dev/null)
  out=$(cd "$d" && "$CLI" wayfind next)
  [[ "$out" == *"SPK-001"* ]] && pass || fail "spike outranks research at equal priority"

  sed -i.bak 's/blocked_by: \[\]/blocked_by:\n  - SPC-001/' "$d/.spectacular/questions/QUE-001-product-choice.md"; rm -f "$d/.spectacular/questions/QUE-001-product-choice.md.bak"
  (cd "$d" && "$CLI" wayfind status >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "cyclic graph refuses sequencing"
  rm -rf "$d"
}

scenario_resolve_and_routes() {
  echo "Scenario metaphors: routes reuse gated verbs and preserve scope"
  local d="/tmp/spectacular-wayfinding-sequencer-routes" out code before after
  new_ws "$d"
  mkdir -p "$d/.spectacular/requests/current"
  printf '%s\n' '---' 'status: active' '---' '# Tasks' '- [ ] Existing milestone scope' > "$d/.spectacular/requests/current/TASKS.md"
  before=$(cksum "$d/.spectacular/requests/current/TASKS.md")
  (cd "$d" && "$CLI" wayfind route "park this idea" later-optimization --hypothesis "Cache the index" --origin "execution discovery" >/dev/null)
  after=$(cksum "$d/.spectacular/requests/current/TASKS.md")
  [[ "$before" == "$after" && -f "$d/.spectacular/ideas/IDEA-001-later-optimization.md" ]] && pass || fail "park route creates idea without scope mutation"

  (cd "$d" && "$CLI" question new pricing --question "Which price?" >/dev/null)
  (cd "$d" && "$CLI" wayfind route icebox q1 --reason "Later milestone" >/dev/null)
  assert_contains "$d/.spectacular/questions/QUE-001-pricing.md" "status: deferred"
  (cd "$d" && "$CLI" wayfind resume q1 >/dev/null)
  (cd "$d" && "$CLI" wayfind resolve q1 --answer "Per seat" >/dev/null)
  assert_contains "$d/.spectacular/questions/QUE-001-pricing.md" "status: resolved"
  out=$(cd "$d" && "$CLI" wayfind route "find your way to" q1)
  [[ "$out" == *"Find your way to QUE-001"* ]] && pass || fail "find route resolves target"

  (cd "$d" && "$CLI" spec new billing --summary "Billing" >/dev/null)
  (cd "$d" && "$CLI" wayfind route "act on goal" s1 --request-slug billing-work >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "act route preserves spec confirmation gate"
  (cd "$d" && "$CLI" spec confirm s1 --evidence "approved" >/dev/null)
  (cd "$d" && "$CLI" wayfind route "act on goal" s1 --request-slug billing-work >/dev/null)
  assert_contains "$d/.spectacular/requests/billing-work/PLAN.md" "source_spec: SPC-001"
  rm -rf "$d"
}

scenario_cross_layer_doctor() {
  echo "Scenario coherence: inferred edges and version inversions warn without mutation"
  local d="/tmp/spectacular-wayfinding-sequencer-doctor" out code before after
  new_ws "$d"
  (cd "$d" && "$CLI" research new future-api --summary "Future API" --target-version v1.1.0-discovery >/dev/null)
  (cd "$d" && "$CLI" research new evidence --summary "Evidence" --target-version v1.0.0-discovery >/dev/null)
  (cd "$d" && "$CLI" spec new consumer --summary "Consumer" --target-version v1.0.0-execution >/dev/null)
  sed -i.bak 's/related: \[\]/blocked_by:\n  - RES-001\nrelated: []/' "$d/.spectacular/specs/SPC-001-consumer.md"; rm -f "$d/.spectacular/specs/SPC-001-consumer.md.bak"
  printf '\n## Dependency note\nThis feature depends on RES-002.\n' >> "$d/.spectacular/specs/SPC-001-consumer.md"
  before=$(cksum "$d/.spectacular/specs/SPC-001-consumer.md")
  out=$(cd "$d" && "$CLI" doctor wayfinding 2>&1); code=$?
  after=$(cksum "$d/.spectacular/specs/SPC-001-consumer.md")
  assert_exit "$code" 1 "coherence findings are warnings"
  [[ "$out" == *"target-version inversion"* && "$out" == *"SPC-001"* ]] && pass || fail "version inversion surfaced: $out"
  [[ "$out" == *"inferred dependency missing from frontmatter"* && "$out" == *"RES-002"* ]] && pass || fail "prose dependency surfaced: $out"
  [[ "$before" == "$after" ]] && pass || fail "doctor does not reslot or write inferred edges"
  rm -rf "$d"
}

scenario_topology_and_ranking
scenario_resolve_and_routes
scenario_cross_layer_doctor

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
