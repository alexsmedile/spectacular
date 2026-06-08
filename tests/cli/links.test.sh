#!/usr/bin/env bash
# tests/cli/links.test.sh — smoke tests for `spectacular links` CLI (v1.16.0+).
#
# Two example scenarios demonstrating the link graph:
#   Example A: depends-on + inverse required-by
#   Example B: blocks + inverse blocked-by + archived dep = satisfied
#
# Each scenario builds an isolated workspace under /tmp/, exercises
# `spectacular links`, and asserts on output. Inverse edges are computed at
# read time; no test modifies B's file to add blocked-by/required-by.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_contains() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — missing: $pattern"
    fail_count=$((fail_count + 1))
  fi
}

assert_lacks() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -qF -- "$pattern"; then
    echo "    ✗ $desc — should not contain: $pattern"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi
}

assert_json_contains() {
  local out="$1" pattern="$2" desc="$3"
  if echo "$out" | grep -q "$pattern"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ $desc — JSON missing: $pattern"
    fail_count=$((fail_count + 1))
  fi
}

mk_workspace() {
  local dir
  dir=$(mktemp -d)
  mkdir -p "$dir/.spectacular/requests"
  cat > "$dir/.spectacular/config.yaml" <<'YAML'
project:
  name: link-graph-test
  owner: test
YAML
  echo "$dir"
}

mk_plan() {
  local dir="$1" slug="$2"
  mkdir -p "$dir/.spectacular/requests/$slug"
  local plan="$dir/.spectacular/requests/$slug/PLAN.md"
  printf -- '---\nstatus: active\npriority: medium\nowner: test\nupdated: 2026-06-08\nsummary: "%s fixture"\nrelated: []\n---\n' "$slug" > "$plan"
}

mk_plan_with_deps() {
  local dir="$1" slug="$2" dep_field="$3" targets="$4"
  mkdir -p "$dir/.spectacular/requests/$slug"
  local plan="$dir/.spectacular/requests/$slug/PLAN.md"
  printf -- '---\n' > "$plan"
  printf 'status: active\npriority: medium\nowner: test\nupdated: 2026-06-08\n' >> "$plan"
  printf 'summary: "%s fixture"\nrelated: []\n' "$slug" >> "$plan"
  printf '%s:\n' "$dep_field" >> "$plan"
  IFS=',' read -ra tarr <<< "$targets"
  for t in "${tarr[@]}"; do
    printf '  - %s\n' "$t" >> "$plan"
  done
  printf -- '---\n' >> "$plan"
}

mk_archive_plan() {
  local dir="$1" slug="$2"
  mkdir -p "$dir/.spectacular/archive/$slug"
  cat > "$dir/.spectacular/archive/$slug/PLAN.md" <<EOF
---
status: verified
priority: medium
owner: test
updated: 2025-01-01
summary: "$slug (archived)"
related: []
EOF
}

# ── Example A: depends-on + inverse required-by ────────────────────────────────
echo "── Example A: depends-on / required-by ──"
W=$(mk_workspace)
mk_plan           "$W" "auth-backend"
mk_plan_with_deps "$W" "user-profile" "depends-on" "auth-backend"

out=$(cd "$W" && bash "$CLI" links 2>&1)

assert_contains "$out" "user-profile" "A: user-profile appears in graph"
assert_contains "$out" "depends-on:  auth-backend" "A: forward depends-on edge shown"
assert_contains "$out" "auth-backend" "A: auth-backend appears in graph"
assert_contains "$out" "required-by: user-profile" "A: inverse required-by computed"
assert_lacks    "$out" "not found"    "A: no dangling warning"

# JSON output
json=$(cd "$W" && bash "$CLI" links --json 2>&1)
assert_json_contains "$json" '"depends_on":\[{"slug":"auth-backend"' "A: JSON depends_on populated"
assert_json_contains "$json" '"required_by":\["user-profile"\]'      "A: JSON required_by populated"

rm -rf "$W"

# ── Example B: blocks + blocked-by + archived dep = satisfied ─────────────────
echo "── Example B: blocks / blocked-by + archived dep ──"
W=$(mk_workspace)
mk_plan           "$W" "ui-redesign"
mk_plan_with_deps "$W" "design-system" "blocks" "ui-redesign"
mk_archive_plan   "$W" "old-components"
mk_plan_with_deps "$W" "widget-library" "depends-on" "old-components"

out=$(cd "$W" && bash "$CLI" links 2>&1)

assert_contains "$out" "design-system" "B: design-system appears"
assert_contains "$out" "blocks:      ui-redesign" "B: forward blocks edge shown"
assert_contains "$out" "blocked-by:  design-system" "B: inverse blocked-by computed"
assert_contains "$out" "widget-library" "B: widget-library appears"
assert_contains "$out" "✓ (shipped)" "B: archived dep shown as satisfied"
assert_lacks    "$out" "⚠ (not found)" "B: no dangling warning for archived dep"

# Single-slug view
single=$(cd "$W" && bash "$CLI" links ui-redesign 2>&1)
assert_contains "$single" "blocked-by:  design-system" "B: single-slug shows blocked-by"
assert_lacks    "$single" "widget-library" "B: single-slug scoped to ui-redesign only"

# doctor links: dangling slug
W2=$(mk_workspace)
mk_plan_with_deps "$W2" "alpha" "depends-on" "nonexistent-slug"
doctor_out=$(cd "$W2" && bash "$CLI" doctor links 2>&1)
assert_contains "$doctor_out" "nonexistent-slug" "B: doctor flags dangling depends-on slug"

rm -rf "$W" "$W2"

# ── Summary ────────────────────────────────────────────────────────────────────
echo ""
echo "Results: $pass_count passed, $fail_count failed"
if [[ $fail_count -gt 0 ]]; then
  exit 1
else
  exit 0
fi
