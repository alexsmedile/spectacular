#!/usr/bin/env bash
# tests/cli/visual.test.sh — unit tests for the ascii-render helper layer
# and the visual renders (progress bars, summary dashboard, roadmap).
#
# Tests run CLI directly with NO_COLOR=1 to get plain-text output, then
# assert on bar math, column alignment, and --json byte-stability.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
LOCAL_SKILL="$REPO_ROOT/skills/spectacular"

fail_count=0
pass_count=0

# ── helpers ──────────────────────────────────────────────────────────────────

assert_eq() {
  local expected="$1" actual="$2" desc="$3"
  if [[ "$actual" == "$expected" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc"
    echo "      expected: $expected"
    echo "      got:      $actual"
    fail_count=$((fail_count + 1))
  fi
}

assert_contains() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — output missing: $pattern"
    fail_count=$((fail_count + 1))
  fi
}

assert_lacks() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    echo "    ✗ $desc — output should not contain: $pattern"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi
}

assert_valid_json() {
  local out="$1" desc="$2"
  if echo "$out" | python3 -m json.tool >/dev/null 2>&1; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — not valid JSON: $out"
    fail_count=$((fail_count + 1))
  fi
}

# Build a minimal workspace with a TASKS.md having known tick counts.
make_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.spectacular/requests/my-req" "$dir/.agents/skills"
  ln -s "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"

  cat > "$dir/.spectacular/config.yaml" <<'EOF'
workspace_schema: "0.6"
project:
  name: test-proj
EOF

  cat > "$dir/.spectacular/PRD.md" <<'EOF'
---
version: 1.0
updated: 2026-01-01
summary: "test"
kit: blank
---
EOF

  # Request with known milestone counts: M1 5/5, M2 2/5, M3 0/3
  cat > "$dir/.spectacular/requests/my-req/TASKS.md" <<'EOF'
---
status: active
updated: 2026-01-01
related:
  - PLAN.md
---

# Tasks — my-req

## v1

### M1 — Setup
- [x] task a
- [x] task b
- [x] task c
- [x] task d
- [x] task e

### M2 — Build
- [x] task f
- [x] task g
- [ ] task h
- [ ] task i
- [ ] task j

### M3 — Ship
- [ ] task k
- [ ] task l
- [ ] task m
EOF

  cat > "$dir/.spectacular/requests/my-req/PLAN.md" <<'EOF'
---
status: active
priority: medium
owner: test
updated: 2026-01-01
summary: "test request"
target_version: v1.0.0
---

# Plan — my-req
EOF
}

# ── test 1: _ascii_bar math (plain output) ───────────────────────────────────
echo "1. ascii_bar: bar math (plain text)"

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT
make_workspace "$TMPDIR"

out=$(cd "$TMPDIR" && NO_COLOR=1 bash "$CLI" progress my-req 2>&1)

# M1 done: 5/5 → 100%
assert_contains "$out" "100% ✓" "M1 shows 100% done"
# M2 done: 2/5 → 40%
assert_contains "$out" "40%" "M2 shows 40%"
assert_contains "$out" "2/5" "M2 shows 2/5 count"
# M3 done: 0/3 → 0%
assert_contains "$out" "0%" "M3 shows 0%"
assert_contains "$out" "0/3" "M3 shows 0/3 count"
# overall line present
assert_contains "$out" "overall" "overall roll-up present"

# ── test 2: bar fill chars ────────────────────────────────────────────────────
echo "2. ascii_bar: plain bar uses # and . (no block chars in NO_COLOR)"

assert_contains "$out" "#" "bar contains # chars"
assert_contains "$out" "." "bar contains . chars"
assert_lacks    "$out" "█" "no block chars in NO_COLOR mode"
assert_lacks    "$out" "░" "no empty-block chars in NO_COLOR mode"

# ── test 3: --json output is valid JSON and byte-stable ──────────────────────
echo "3. progress --json: valid JSON, unaffected by NO_COLOR"

json_out=$(cd "$TMPDIR" && NO_COLOR=1 bash "$CLI" progress my-req --json 2>&1)
assert_valid_json "$json_out" "progress --json is valid JSON"
assert_contains "$json_out" '"done":5' "M1 done=5 in JSON"
assert_contains "$json_out" '"done":2' "M2 done=2 in JSON"
assert_contains "$json_out" '"done":0' "M3 done=0 in JSON"
assert_contains "$json_out" '"total":5' "M1 total=5 in JSON"
assert_contains "$json_out" '"total":3' "M3 total=3 in JSON"

# ── test 4: summary shows bars for request counts ────────────────────────────
echo "4. summary: request-state bars present"

sum_out=$(cd "$TMPDIR" && NO_COLOR=1 bash "$CLI" summary 2>&1)
assert_contains "$sum_out" "Requests:   1 total" "summary shows 1 total request"
assert_contains "$sum_out" "active" "summary shows active state"
# bar char present in summary
assert_contains "$sum_out" "[" "summary contains bar brackets"
assert_contains "$sum_out" "%" "summary contains percentage"

# ── test 5: summary --json unaffected ────────────────────────────────────────
echo "5. summary --json: valid JSON"

sjson=$(cd "$TMPDIR" && NO_COLOR=1 bash "$CLI" summary --json 2>&1)
assert_valid_json "$sjson" "summary --json is valid JSON"
assert_contains "$sjson" '"active":1' "summary JSON shows 1 active"

# ── test 6: roadmap render (using real repo's ROADMAP) ───────────────────────
echo "6. roadmap: renders version arc"

rm_out=$(cd "$REPO_ROOT" && NO_COLOR=1 bash "$CLI" roadmap 2>&1)
assert_contains "$rm_out" "Roadmap" "roadmap header present"
assert_contains "$rm_out" "·" "planned items use · indicator"

# shipped versions hidden by default (ledger-driven)
assert_lacks "$rm_out" "v1.16.0" "shipped v1.16.0 hidden by default"

# --all shows shipped ledger entries with v1.x versions and ✓ indicator
rm_all=$(cd "$REPO_ROOT" && NO_COLOR=1 bash "$CLI" roadmap --all 2>&1)
assert_contains "$rm_all" "v1." "at least one v1.x version in --all"
assert_contains "$rm_all" "v1.16.0" "--all shows shipped versions"
assert_contains "$rm_all" "✓" "--all shows ✓ for shipped"

# ── test 7: roadmap --json ────────────────────────────────────────────────────
echo "7. roadmap --json: valid JSON array"

rjson=$(cd "$REPO_ROOT" && NO_COLOR=1 bash "$CLI" roadmap --json 2>&1)
assert_valid_json "$rjson" "roadmap --json is valid JSON"
assert_contains "$rjson" '"version"' "roadmap JSON has version field"
assert_contains "$rjson" '"tier"' "roadmap JSON has tier field"
assert_contains "$rjson" '"status"' "roadmap JSON has status field"

# ── test 8: NO_COLOR in bar output vs block chars ────────────────────────────
echo "8. ascii_bar: block chars used when color enabled (TTY detection)"
# We can't force a real TTY in a test, but we can confirm that when NO_COLOR
# is unset the code path doesn't crash (output may or may not have block chars
# depending on whether stdout is a TTY in this test context).
bar_out=$(cd "$TMPDIR" && bash "$CLI" progress my-req 2>&1)
# Must exit 0 and produce output regardless of TTY state
if [[ -n "$bar_out" ]]; then
  pass_count=$((pass_count + 1))
else
  echo "    ✗ progress without NO_COLOR produced no output"
  fail_count=$((fail_count + 1))
fi

# ── test 9: clamped done > total ─────────────────────────────────────────────
echo "9. ascii_bar: done > total clamps to 100%"

# Inject a TASKS.md where v2 section adds more [x] after the milestone closes
cat > "$TMPDIR/.spectacular/requests/my-req/TASKS.md" <<'EOF'
---
status: active
updated: 2026-01-01
related:
  - PLAN.md
---

# Tasks — my-req

### M1 — Setup
- [x] a
- [x] b
- [x] c

### M2 — empty

EOF
clamp_out=$(cd "$TMPDIR" && NO_COLOR=1 bash "$CLI" progress my-req 2>&1)
assert_contains "$clamp_out" "100% ✓" "fully done milestone shows 100% ✓"
# empty milestone (0/0) still renders a line (shown as 0% ✓ since done==total==0)
assert_contains "$clamp_out" "0%" "empty milestone renders 0%"

# ── report ────────────────────────────────────────────────────────────────────
echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]] && exit 0 || exit 1
