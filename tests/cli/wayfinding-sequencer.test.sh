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
  local qid rid sid cid
  qid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/questions/product-choice.md")
  (cd "$d" && "$CLI" research new market-data --summary "Research market" --blocked-by "$qid" >/dev/null)
  rid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/research/market-data.md")
  (cd "$d" && "$CLI" spike new feasibility --summary "Test feasibility" --blocked-by "$qid" >/dev/null)
  sid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/spikes/feasibility.md")
  (cd "$d" && "$CLI" spec new launch --summary "Launch" --target-version v1.0.0-execution >/dev/null)
  cid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/specs/launch.md")
  sed -i.bak "s/related: \[\]/blocked_by:\n  - $rid\n  - $sid\nrelated: []/" "$d/.spectacular/specs/launch.md"; rm -f "$d/.spectacular/specs/launch.md.bak"

  order=$(cd "$d" && "$CLI" wayfind order)
  [[ "$order" == *"$qid"*$'\tfrontier'* && "$order" == *"$rid"*$'\tfog'* && "$order" == *"$sid"*$'\tfog'* && "$order" == *"$cid"*$'\tfog'* ]] && pass || fail "strict dependency-first order: $order"
  out=$(cd "$d" && "$CLI" wayfind next)
  [[ "$out" == *"$qid"* ]] && pass || fail "human question wins equal-priority uncertainty tie"

  sed -i.bak "s/blocked_by: \[\]/blocked_by:\n  - $cid/" "$d/.spectacular/questions/product-choice.md"; rm -f "$d/.spectacular/questions/product-choice.md.bak"
  (cd "$d" && "$CLI" wayfind status >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "cyclic graph refuses sequencing"
  sed -i.bak "s/blocked_by:/blocked_by: []/; /  - $cid/d" "$d/.spectacular/questions/product-choice.md"; rm -f "$d/.spectacular/questions/product-choice.md.bak"

  (cd "$d" && "$CLI" question resolve "$qid" --answer "Enterprise" >/dev/null)
  out=$(cd "$d" && "$CLI" wayfind next)
  [[ "$out" == *"$sid"* ]] && pass || fail "spike outranks research at equal priority"
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
  [[ "$before" == "$after" && -f "$d/.spectacular/ideas/later-optimization.md" ]] && pass || fail "park route creates idea without scope mutation"

  (cd "$d" && "$CLI" question new pricing --question "Which price?" >/dev/null)
  local qid bid
  qid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/questions/pricing.md")
  (cd "$d" && "$CLI" wayfind route icebox "$qid" --reason "Later milestone" >/dev/null)
  assert_contains "$d/.spectacular/questions/pricing.md" "status: deferred"
  (cd "$d" && "$CLI" wayfind resume "$qid" >/dev/null)
  (cd "$d" && "$CLI" wayfind resolve "$qid" --answer "Per seat" >/dev/null)
  assert_contains "$d/.spectacular/archive/questions/pricing.md" "archived_from: resolved"
  out=$(cd "$d" && "$CLI" wayfind route "find your way to" "$qid")
  [[ "$out" == *"Find your way to $qid"* ]] && pass || fail "find route resolves target"

  (cd "$d" && "$CLI" spec new billing --summary "Billing" >/dev/null)
  bid=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/specs/billing.md")
  (cd "$d" && "$CLI" wayfind route "act on goal" "$bid" --request-slug billing-work >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "act route preserves spec confirmation gate"
  out=$(cd "$d" && "$CLI" wayfind route "act on goal" "$bid" --request-slug billing-work 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "unmerged act route remains blocked"
  [[ "$out" == *"$bid"* ]] && pass || fail "act route names the contract"
  rm -rf "$d"
}

scenario_cross_layer_doctor() {
  echo "Scenario coherence: inferred edges and version inversions warn without mutation"
  local d="/tmp/spectacular-wayfinding-sequencer-doctor" out code before after
  new_ws "$d"
  (cd "$d" && "$CLI" research new future-api --summary "Future API" --target-version v1.1.0-discovery >/dev/null)
  (cd "$d" && "$CLI" research new evidence --summary "Evidence" --target-version v1.0.0-discovery >/dev/null)
  (cd "$d" && "$CLI" spec new consumer --summary "Consumer" --target-version v1.0.0-execution >/dev/null)
  local future_id evidence_id consumer_id consumer
  future_id=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/research/future-api.md")
  evidence_id=$(awk '/^id:/{print $2; exit}' "$d/.spectacular/research/evidence.md")
  consumer="$d/.spectacular/specs/consumer.md"; consumer_id=$(awk '/^id:/{print $2; exit}' "$consumer")
  sed -i.bak "s/related: \[\]/blocked_by:\n  - $future_id\nrelated: []/" "$consumer"; rm -f "$consumer.bak"
  printf '\n## Dependency note\nThis feature depends on %s.\n' "$evidence_id" >> "$consumer"
  before=$(cksum "$consumer")
  out=$(cd "$d" && "$CLI" doctor wayfinding 2>&1); code=$?
  after=$(cksum "$consumer")
  assert_exit "$code" 1 "coherence findings are warnings"
  [[ "$out" == *"target-version inversion"* && "$out" == *"$consumer_id"* ]] && pass || fail "version inversion surfaced: $out"
  [[ "$out" == *"inferred dependency missing from frontmatter"* && "$out" == *"$evidence_id"* ]] && pass || fail "prose dependency surfaced: $out"
  [[ "$before" == "$after" ]] && pass || fail "doctor does not reslot or write inferred edges"
  rm -rf "$d"
}

scenario_topology_and_ranking
scenario_resolve_and_routes
scenario_cross_layer_doctor

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
