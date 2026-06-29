#!/usr/bin/env bash
# tests/cli/roadmap-migrate.test.sh — `spectacular roadmap migrate` index mode (v1.23.0+, b18)
#
# Covers:
#   M1 dry-run        — reports moves, writes nothing
#   M2 migrate        — moves shipped-beyond-keep to roadmap/v*.md, keeps newest N inline,
#                       writes ## Shipped index, leaves planned/active/vision blocks alone
#   M3 idempotence    — re-run moves nothing
#   M4 doctor roadmap — index integrity (clean), orphan + stale detection, flat-mode nudge

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_file_exists()   { if [[ -f "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected file: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_absent()   { if [[ ! -f "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected absent file: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_contains() { if [[ -f "$1" ]] && grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_file_lacks()    { if [[ -f "$1" ]] && ! grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to NOT contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_output_contains() { if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count+1)); else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count+1)); fi; }

# Build a workspace with a ROADMAP carrying 5 shipped v-blocks + 1 planned + 1 vision.
seed() {
  local dir="$1"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular"
  cat > "$dir/.spectacular/ROADMAP.md" <<'EOF'
---
version: 1.0
updated: 2026-01-01
summary: "test roadmap"
---

# Roadmap

## Roadmap ledger

| build | slug | title | tier | target-version | status |
|-------|------|-------|------|----------------|--------|
| b1 | a | A | full | v1.0.0 | shipped |
| b2 | b | B | full | v1.1.0 | shipped |

---

## v1.0.0 — First

**Status:** shipped (2026-01-01)

Prose for v1.0.0.

## v1.1.0 — Second

**Status:** shipped (2026-01-02)

Prose for v1.1.0.

## v1.2.0 — Third

**Status:** shipped (2026-01-03)

Prose for v1.2.0.

## v1.3.0 — Fourth

**Status:** shipped (2026-01-04)

Prose for v1.3.0.

## v1.4.0 — Fifth

**Status:** shipped (2026-01-05)

Prose for v1.4.0.

## v1.5.0 — Planned next

**Status:** planned

Forward prose — must stay inline.

## v2.x — Vision

**Status:** someday

Direction only — must stay inline.

## Icebox

- some idea
EOF
}

scenario_1_dry_run() {
  echo "Scenario 1 (M1): dry-run reports moves, writes nothing"
  local dir="/tmp/spectacular-rmmig-1"; seed "$dir"
  local out; out=$(cd "$dir" && "$CLI" roadmap migrate --dry-run 2>&1)
  # 5 shipped, keep 3 → move 2 oldest (v1.0.0, v1.1.0)
  assert_output_contains "$out" "would move 2 of 5 shipped"
  assert_output_contains "$out" "roadmap/v1.0.0.md"
  assert_output_contains "$out" "roadmap/v1.1.0.md"
  assert_file_absent "$dir/.spectacular/roadmap/v1.0.0.md"   # nothing written
  assert_file_contains "$dir/.spectacular/ROADMAP.md" "## v1.0.0 — First"  # still inline
  rm -rf "$dir"
}

scenario_2_migrate() {
  echo "Scenario 2 (M2): migrate moves oldest shipped, keeps newest 3, indexes"
  local dir="/tmp/spectacular-rmmig-2"; seed "$dir"
  (cd "$dir" && "$CLI" roadmap migrate >/dev/null 2>&1)
  local rm="$dir/.spectacular/ROADMAP.md"
  # moved (oldest 2)
  assert_file_exists "$dir/.spectacular/roadmap/v1.0.0.md"
  assert_file_exists "$dir/.spectacular/roadmap/v1.1.0.md"
  assert_file_contains "$dir/.spectacular/roadmap/v1.0.0.md" "Prose for v1.0.0."
  assert_file_lacks "$rm" "## v1.0.0 — First"
  assert_file_lacks "$rm" "## v1.1.0 — Second"
  # kept inline (newest 3 shipped)
  assert_file_contains "$rm" "## v1.2.0 — Third"
  assert_file_contains "$rm" "## v1.4.0 — Fifth"
  # planned + vision never move
  assert_file_contains "$rm" "## v1.5.0 — Planned next"
  assert_file_contains "$rm" "## v2.x — Vision"
  # shipped index present, before Icebox
  assert_file_contains "$rm" "## Shipped"
  assert_file_contains "$rm" "v1.0.0 → roadmap/v1.0.0.md"
  assert_file_contains "$rm" "## Icebox"
  rm -rf "$dir"
}

scenario_3_idempotent() {
  echo "Scenario 3 (M3): re-run moves nothing"
  local dir="/tmp/spectacular-rmmig-3"; seed "$dir"
  (cd "$dir" && "$CLI" roadmap migrate >/dev/null 2>&1)
  local out; out=$(cd "$dir" && "$CLI" roadmap migrate 2>&1)
  assert_output_contains "$out" "nothing to migrate"
  rm -rf "$dir"
}

scenario_4_doctor() {
  echo "Scenario 4 (M4): doctor roadmap — clean, orphan, stale, flat nudge"
  local dir="/tmp/spectacular-rmmig-4"; seed "$dir"

  # flat-mode nudge: 5 shipped inline > keep 3
  local out; out=$(cd "$dir" && "$CLI" doctor roadmap 2>&1)
  assert_output_contains "$out" "shipped prose blocks inline"

  # after migrate → clean
  (cd "$dir" && "$CLI" roadmap migrate >/dev/null 2>&1)
  out=$(cd "$dir" && "$CLI" doctor roadmap 2>&1)
  assert_output_contains "$out" "have a corresponding file"

  # stale file (no index line)
  echo "x" > "$dir/.spectacular/roadmap/v0.0.1.md"
  out=$(cd "$dir" && "$CLI" doctor roadmap 2>&1)
  assert_output_contains "$out" "no Shipped index line: v0.0.1"

  # orphan index line (delete a real file)
  rm "$dir/.spectacular/roadmap/v1.0.0.md"
  out=$(cd "$dir" && "$CLI" doctor roadmap 2>&1)
  assert_output_contains "$out" "orphan Shipped index line"
  rm -rf "$dir"
}

scenario_1_dry_run
scenario_2_migrate
scenario_3_idempotent
scenario_4_doctor

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
