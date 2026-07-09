#!/usr/bin/env bash
# tests/cli/status.test.sh — smoke tests for `spectacular status` fleet view.
#
# Covers the deterministic surfaces added by status-fleet-view (b23):
#   - bare `status`          → aligned fleet table (frontmatter + body signals)
#   - `status <slug>`        → single request card
#   - `status --json`        → structured contract for agents
# Body signals assume the enforced canonical schema (## Goal, ### M milestones,
# flush-left checkboxes). Each scenario builds an isolated workspace under /tmp/.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_exit_code() {
  local expected="$1" actual="$2" desc="$3"
  if [[ "$actual" == "$expected" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ $desc — expected exit $expected, got $actual"; fail_count=$((fail_count + 1)); fi
}
assert_output_contains() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then pass_count=$((pass_count + 1))
  else echo "    ✗ $desc — output missing: $pattern"; fail_count=$((fail_count + 1)); fi
}
assert_output_lacks() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    echo "    ✗ $desc — output should not contain: $pattern"; fail_count=$((fail_count + 1))
  else pass_count=$((pass_count + 1)); fi
}

# Seed a minimal workspace (just enough for status to iterate requests/).
seed_ws() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.spectacular/requests" "$dir/.spectacular/archive"
  cat > "$dir/.spectacular/config.yaml" <<'EOF'
workspace_schema: "2.0"
EOF
}

# Write a canonical request. Args: dir slug status priority build updated summary goal
seed_request() {
  local dir="$1" slug="$2" status="$3" prio="$4" build="$5" upd="$6" summary="$7" goal="$8"
  mkdir -p "$dir/.spectacular/requests/$slug"
  cat > "$dir/.spectacular/requests/$slug/PLAN.md" <<EOF
---
status: $status
priority: $prio
owner: t
updated: $upd
build: $build
summary: "$summary"
related:
  - PRD.md
---
# Plan — $slug
## Goal
$goal
## Constraints
- c
## Milestones
- M1 — one
- M2 — two
## Tasks
See TASKS.md
## Dependencies
- none
## Validation
- M1 — run: true
## Deliverables
- thing
EOF
}

scenario_1_fleet_table() {
  echo "Scenario 1: bare status renders fleet table sorted active→planned"
  local dir="/tmp/status-test-1"
  seed_ws "$dir"
  seed_request "$dir" "aaa-planned" "planned" "medium" "b5" "2026-07-01" "planned one" "goal a"
  seed_request "$dir" "zzz-active"  "active"  "high"   "b7" "2026-07-05" "active one"  "goal z"
  # minimal TASKS for progress
  cat > "$dir/.spectacular/requests/zzz-active/TASKS.md" <<'EOF'
---
status: active
related:
  - PLAN.md
---
# Tasks
### M1 — one
- [x] a
- [ ] b
EOF

  local out code
  out=$(cd "$dir" && "$CLI" status 2>&1) && code=0 || code=$?

  assert_exit_code "0" "$code" "fleet table exits 0"
  assert_output_contains "$out" "aaa-planned" "lists planned slug"
  assert_output_contains "$out" "zzz-active" "lists active slug"
  assert_output_contains "$out" "b7" "shows build column"
  assert_output_contains "$out" "1/2" "shows top-level progress for active request"
  # active must sort before planned: zzz-active line precedes aaa-planned line
  local order
  order=$(echo "$out" | grep -nE 'zzz-active|aaa-planned' | head -2)
  if echo "$order" | head -1 | grep -q 'zzz-active'; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ active sorts before planned — got: $order"; fail_count=$((fail_count + 1))
  fi
  rm -rf "$dir"
}

scenario_2_body_signals() {
  echo "Scenario 2: progress ignores indented subtasks, counts [~] as deferred"
  local dir="/tmp/status-test-2"
  seed_ws "$dir"
  seed_request "$dir" "mixed" "active" "high" "b9" "2026-01-01" "mixed" "ship it"
  cat > "$dir/.spectacular/requests/mixed/TASKS.md" <<'EOF'
---
status: active
related:
  - PLAN.md
---
# Tasks
### M1 — one
- [x] done
- [~] deferred
  - [ ] nested subtask (must NOT count)
### M2 — two
- [ ] open one
- [ ] open two
EOF

  local out
  out=$(cd "$dir" && "$CLI" status mixed 2>&1)
  assert_output_contains "$out" "1/3 (+1 def)" "progress: 1 done / 3 total, 1 deferred, nested ignored"
  assert_output_contains "$out" "current: M2 — two" "current milestone advances past fully-done M1"
  assert_output_contains "$out" "stale" "active request updated >14d ago flags stale"
  assert_output_contains "$out" "ship it" "card shows goal line"
  rm -rf "$dir"
}

scenario_3_json_contract() {
  echo "Scenario 3: status --json emits valid structured output with body-signal fields"
  local dir="/tmp/status-test-3"
  seed_ws "$dir"
  seed_request "$dir" "req-a" "planned" "low" "b1" "2026-07-01" "the a request" "do a"
  cat > "$dir/.spectacular/requests/req-a/TASKS.md" <<'EOF'
---
status: planned
related:
  - PLAN.md
---
# Tasks
### M1 — one
- [ ] a
- [x] b
EOF

  local out code
  out=$(cd "$dir" && "$CLI" status --json 2>&1) && code=0 || code=$?
  assert_exit_code "0" "$code" "--json exits 0"

  # Validate JSON if a parser is available; otherwise assert on key presence.
  if command -v python3 >/dev/null 2>&1; then
    if echo "$out" | python3 -m json.tool >/dev/null 2>&1; then
      pass_count=$((pass_count + 1))
    else
      echo "    ✗ --json output is not valid JSON"; fail_count=$((fail_count + 1))
    fi
  fi
  assert_output_contains "$out" '"slug":"req-a"' "json has slug"
  assert_output_contains "$out" '"current_milestone":"M1 — one"' "json has current_milestone"
  assert_output_contains "$out" '"done":1' "json progress done count"
  assert_output_contains "$out" '"total":2' "json progress total count"
  rm -rf "$dir"
}

scenario_4_card_missing_slug() {
  echo "Scenario 4: status <unknown-slug> errors cleanly"
  local dir="/tmp/status-test-4"
  seed_ws "$dir"
  seed_request "$dir" "real" "planned" "medium" "b1" "2026-07-01" "real" "goal"
  local out code
  out=$(cd "$dir" && "$CLI" status no-such-slug 2>&1) && code=0 || code=$?
  assert_exit_code "1" "$code" "unknown slug exits non-zero"
  assert_output_contains "$out" "no request 'no-such-slug'" "clear not-found message"
  rm -rf "$dir"
}

echo "═══ status.test.sh ═══"
scenario_1_fleet_table
scenario_2_body_signals
scenario_3_json_contract
scenario_4_card_missing_slug

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
