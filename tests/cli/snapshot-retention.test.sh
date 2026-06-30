#!/usr/bin/env bash
# tests/cli/snapshot-retention.test.sh — snapshot system (b16)
#
# Covers:
#   M1  allowlist: DESIGN.md snapshot-able, non-canonical refused
#   M2  version coupling: @v<version> filename for versioned docs, counter fallback
#   M3  tiered retention + prune (origin + periodic + recent; dry-run / apply)
#   M4  folder rename snapshots/ → _snapshots/ migration; gitignore toggle

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

pass() { pass_count=$((pass_count + 1)); }
fail() { echo "    ✗ $1"; fail_count=$((fail_count + 1)); }

assert_file_exists()  { [[ -f "$1" ]] && pass || fail "expected file: $1"; }
assert_file_absent()  { [[ ! -f "$1" ]] && pass || fail "expected absent: $1"; }
assert_file_contains(){ [[ -f "$1" ]] && grep -qF -- "$2" "$1" && pass || fail "'$1' should contain '$2'"; }
assert_exit()         { [[ "$1" -eq "$2" ]] && pass || fail "$3: exit $1 want $2"; }
assert_eq()           { [[ "$1" == "$2" ]] && pass || fail "$3: got '$1' want '$2'"; }
assert_output_contains(){ echo "$1" | grep -qF -- "$2" && pass || fail "output should contain: $2"; }

# A workspace with a snapshots config block (period off, keep K) and N
# pre-made snapshots for one doc, each stamped with a given updated: date.
seed_store() {
  local dir="$1" doc="$2" keep="$3" period="$4"; shift 4
  rm -rf "$dir"; mkdir -p "$dir/.spectacular/_snapshots/$doc"
  cat > "$dir/.spectacular/config.yaml" <<EOF
name: $(basename "$dir")
snapshots:
  keep: $keep
  period: $period
EOF
  local i=1
  for d in "$@"; do
    printf -- '---\nupdated: %s\n---\nbody %s\n' "$d" "$i" > "$dir/.spectacular/_snapshots/$doc/@v${i}.md"
    i=$((i + 1))
  done
}

# ── M1 + M2: allowlist + version coupling ─────────────────────────────────────
scenario_allowlist_and_coupling() {
  echo "Scenario A: DESIGN.md allowed; @v couples to version:; counter fallback"
  local dir="/tmp/spectacular-snap-A"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular"

  # versioned doc → @v<version>, live doc bumps
  printf -- '---\nversion: 1.3\nupdated: 2026-06-30\n---\n# PRD\nx\n' > "$dir/.spectacular/PRD.md"
  (cd "$dir" && "$CLI" snapshot .spectacular/PRD.md >/dev/null)
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v1.3.md"
  assert_file_contains "$dir/.spectacular/PRD.md" "version: 1.4"

  # DESIGN.md is snapshot-able (M1) and has no version → counter @v1, @v2 (M2 fallback)
  printf -- '---\ntitle: design\n---\n# D\na\n' > "$dir/.spectacular/DESIGN.md"
  (cd "$dir" && "$CLI" snapshot .spectacular/DESIGN.md >/dev/null)
  assert_file_exists "$dir/.spectacular/_snapshots/DESIGN/@v1.md"
  assert_file_contains "$dir/.spectacular/DESIGN.md" "title: design"   # no version injected
  printf -- '---\ntitle: design\n---\n# D\nb changed\n' > "$dir/.spectacular/DESIGN.md"
  (cd "$dir" && "$CLI" snapshot .spectacular/DESIGN.md >/dev/null)
  assert_file_exists "$dir/.spectacular/_snapshots/DESIGN/@v2.md"

  # non-canonical refused
  printf -- '---\ntitle: x\n---\n' > "$dir/.spectacular/NOTES.md"
  local code; (cd "$dir" && "$CLI" snapshot .spectacular/NOTES.md >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "non-canonical refused"

  rm -rf "$dir"
}

# ── M3: tiered retention + prune ──────────────────────────────────────────────
scenario_prune_off() {
  echo "Scenario B: period off, @v1..@v6 keep 3 → keep @v1 @v4 @v5 @v6"
  local dir="/tmp/spectacular-snap-B"
  seed_store "$dir" PRD 3 off 2026-01-01 2026-01-02 2026-01-03 2026-01-04 2026-01-05 2026-01-06

  # dry-run changes nothing
  local out; out=$(cd "$dir" && "$CLI" snapshot prune 2>&1)
  assert_output_contains "$out" "would prune"
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v2.md"   # still there after dry-run

  # apply → @v2 @v3 gone, rest kept (non-git workspace → .trash)
  (cd "$dir" && "$CLI" snapshot prune --apply >/dev/null)
  assert_file_exists  "$dir/.spectacular/_snapshots/PRD/@v1.md"
  assert_file_absent  "$dir/.spectacular/_snapshots/PRD/@v2.md"
  assert_file_absent  "$dir/.spectacular/_snapshots/PRD/@v3.md"
  assert_file_exists  "$dir/.spectacular/_snapshots/PRD/@v4.md"
  assert_file_exists  "$dir/.spectacular/_snapshots/PRD/@v5.md"
  assert_file_exists  "$dir/.spectacular/_snapshots/PRD/@v6.md"
  assert_file_exists  "$dir/.spectacular/.trash/_snapshots/PRD/@v2.md"   # recoverable

  # idempotent
  out=$(cd "$dir" && "$CLI" snapshot prune --apply 2>&1)
  assert_output_contains "$out" "nothing to prune"
  rm -rf "$dir"
}

scenario_prune_monthly() {
  echo "Scenario C: monthly, keep 3, heavy April churn → keep @v1 @v9 @v10 @v11"
  local dir="/tmp/spectacular-snap-C"
  seed_store "$dir" PRD 3 month \
    2026-01-05 2026-04-02 2026-04-03 2026-04-04 2026-04-05 2026-04-06 \
    2026-04-07 2026-04-08 2026-04-09 2026-04-10 2026-04-11
  (cd "$dir" && "$CLI" snapshot prune --apply >/dev/null)
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v1.md"    # origin
  assert_file_absent "$dir/.spectacular/_snapshots/PRD/@v5.md"    # mid-April, no tier
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v9.md"    # recent
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v11.md"   # newest = April periodic + recent
  rm -rf "$dir"
}

scenario_prune_git() {
  echo "Scenario D: tracked snapshots pruned via git rm (no .trash)"
  local dir="/tmp/spectacular-snap-D"
  seed_store "$dir" PRD 2 off 2026-01-01 2026-01-02 2026-01-03 2026-01-04 2026-01-05
  (cd "$dir" && git init -q && git config user.email t@t.co && git config user.name t && git add -A && git commit -qm init)
  (cd "$dir" && "$CLI" snapshot prune --apply >/dev/null)
  assert_file_absent "$dir/.spectacular/_snapshots/PRD/@v2.md"
  assert_file_absent "$dir/.spectacular/.trash/_snapshots/PRD/@v2.md"   # git path, not .trash
  # keep 2 + origin = @v1 @v4 @v5 (3 kept); prune @v2 @v3 = 2 deletions
  local st; st=$(cd "$dir" && git status --short | grep -c "^D")
  assert_eq "$st" "2" "2 deletions staged (keep @v1 @v4 @v5)"
  rm -rf "$dir"
}

# ── M4: folder migration + gitignore ──────────────────────────────────────────
scenario_folder_migration() {
  echo "Scenario E: snapshots/ → _snapshots/ via doctor --fix, lossless"
  local dir="/tmp/spectacular-snap-E"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular/snapshots/PRD"
  printf 'name: e\n' > "$dir/.spectacular/config.yaml"   # default folder = _snapshots
  printf -- '---\nupdated: 2026-06-01\n---\nx\n' > "$dir/.spectacular/snapshots/PRD/@v1.0.md"

  local out; out=$(cd "$dir" && "$CLI" doctor snapshots 2>&1)
  assert_output_contains "$out" "config wants _snapshots/"
  (cd "$dir" && "$CLI" doctor --fix snapshots >/dev/null 2>&1)
  assert_file_exists "$dir/.spectacular/_snapshots/PRD/@v1.0.md"
  assert_file_absent "$dir/.spectacular/snapshots/PRD/@v1.0.md"
  rm -rf "$dir"
}

scenario_gitignore_toggle() {
  echo "Scenario F: snapshots.gitignore true adds line, false removes it"
  local dir="/tmp/spectacular-snap-F"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular/_snapshots/PRD"
  printf 'name: f\nsnapshots:\n  gitignore: true\n' > "$dir/.spectacular/config.yaml"
  printf '.spectacular.local/\n' > "$dir/.gitignore"

  (cd "$dir" && "$CLI" doctor --fix snapshots >/dev/null 2>&1)
  assert_file_contains "$dir/.gitignore" ".spectacular/_snapshots/"
  assert_file_contains "$dir/.gitignore" ".spectacular.local/"   # baseline preserved

  printf 'name: f\nsnapshots:\n  gitignore: false\n' > "$dir/.spectacular/config.yaml"
  (cd "$dir" && "$CLI" doctor --fix snapshots >/dev/null 2>&1)
  if grep -q "_snapshots" "$dir/.gitignore" 2>/dev/null; then
    fail "store ignore line should be removed when gitignore=false"
  else
    pass
  fi
  rm -rf "$dir"
}

# ── gap check: skip when counter + version names are mixed (b16 transition) ───
scenario_mixed_scheme_no_false_gap() {
  echo "Scenario G: mixed @vN + @vX.Y names skip the gap check (no false positive)"
  local dir="/tmp/spectacular-snap-G"
  rm -rf "$dir"; mkdir -p "$dir/.spectacular/_snapshots/SPEC"
  printf 'name: g\n' > "$dir/.spectacular/config.yaml"
  # counter-named (pre-b16) + one version-named (post-b16) in the same dir
  for v in 1 2 3; do printf -- '---\nupdated: 2026-0%s-01\n---\nx\n' "$v" > "$dir/.spectacular/_snapshots/SPEC/@v${v}.md"; done
  printf -- '---\nupdated: 2026-06-01\n---\nx\n' > "$dir/.spectacular/_snapshots/SPEC/@v1.6.md"
  local out; out=$(cd "$dir" && "$CLI" doctor snapshots 2>&1)
  if echo "$out" | grep -qi "version-sequence gap"; then
    fail "mixed scheme should not report a version-sequence gap"
  else
    pass
  fi
  assert_output_contains "$out" "mixed counter + version names"
  rm -rf "$dir"
}

scenario_allowlist_and_coupling
scenario_prune_off
scenario_prune_monthly
scenario_prune_git
scenario_folder_migration
scenario_gitignore_toggle
scenario_mixed_scheme_no_false_gap

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]]
