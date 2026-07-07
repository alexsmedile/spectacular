#!/usr/bin/env bash
# tests/cli/archive-closure-gate.test.sh — the archive closure gate (v1.28.0+, b22)
#
# Covers:
#   1. tasks check    — open/reasonless-deferred box blocks; [x] and [~]+reason pass
#   2. verify check   — VERIFY.md without a ✅ walk row blocks; walked passes
#   3. spec check     — missing SPEC-DELTA.md blocks; NONE and a real delta pass
#   4. override       — --override <check> --reason records into archive_overrides:
#   5. no-op override  — overriding a check that already passes warns, doesn't error
#   6. --force split  — --force does NOT bypass closure checks (only the status gate)
#   7. grandfathered  — already-archived requests are never revalidated
#   8. delta integrity — doctor specs flags a MODIFIED delta quoting a missing bullet

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_dir_exists()  { if [[ -d "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected dir: $1"; fail_count=$((fail_count+1)); fi; }
assert_dir_absent()  { if [[ ! -d "$1" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ expected absent dir: $1"; fail_count=$((fail_count+1)); fi; }
assert_file_contains() { if [[ -f "$1" ]] && grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_file_lacks()    { if [[ -f "$1" ]] && ! grep -qF -- "$2" "$1"; then pass_count=$((pass_count+1)); else echo "    ✗ expected '$1' to NOT contain '$2'"; fail_count=$((fail_count+1)); fi; }
assert_output_contains() { if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count+1)); else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count+1)); fi; }
assert_exit() { if [[ "$1" -eq "$2" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ $3: exit $1 want $2"; fail_count=$((fail_count+1)); fi; }

seed() {
  local dir="$1"
  rm -rf "$dir"; mkdir -p "$dir"
  (cd "$dir" && git init -q 2>/dev/null)
  (cd "$dir" && "$CLI" init --kit blank --name "$(basename "$dir")" >/dev/null 2>&1)
}

# tick every TASKS box in a request
tick_tasks() {
  local t="$1/TASKS.md"
  [[ -f "$t" ]] && { sed -i.bak 's/^- \[ \]/- [x]/' "$t"; rm -f "$t.bak"; }
}
none_delta() { printf 'NONE — fixture, no spec impact\n' > "$1/SPEC-DELTA.md"; }

# ── Scenario 1: tasks check ───────────────────────────────────────────────────
scenario_1_tasks() {
  echo "Scenario 1: tasks check — open box blocks; [x] / [~]+reason pass"
  local dir="/tmp/spectacular-gate-1"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  none_delta "$R"

  # open boxes → block, naming the check
  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "open tasks blocks"
  assert_output_contains "$out" "✗ tasks"

  # tick all → pass
  tick_tasks "$R"
  (cd "$dir" && "$CLI" archive r --skip-doctor >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "ticked tasks pass"
  assert_dir_exists "$dir/.spectacular/archive/r"
  rm -rf "$dir"
}

# deferred-with-reason passes; reasonless [~] blocks
scenario_1b_deferred() {
  echo "Scenario 1b: [~] passes only with a ' — reason'"
  local dir="/tmp/spectacular-gate-1b"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  none_delta "$R"
  # make all boxes [~] WITHOUT a reason → block
  sed -i.bak 's/^- \[ \]/- [~]/' "$R/TASKS.md"; rm -f "$R/TASKS.md.bak"
  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "reasonless [~] blocks"
  assert_output_contains "$out" "deferred without reason"
  # add a reason to every [~] line → pass
  sed -i.bak 's/^- \[~\] \(.*\)$/- [~] \1 — deferred to next build/' "$R/TASKS.md"; rm -f "$R/TASKS.md.bak"
  (cd "$dir" && "$CLI" archive r --skip-doctor >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "[~]+reason passes"
  rm -rf "$dir"
}

# ── Scenario 2: verify check ──────────────────────────────────────────────────
scenario_2_verify() {
  echo "Scenario 2: verify check — VERIFY.md without a ✅ walk row blocks"
  local dir="/tmp/spectacular-gate-2"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  tick_tasks "$R"; none_delta "$R"
  printf -- '- [ ] some check\n' > "$R/VERIFY.md"   # exists, never walked (no log)

  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "unwalked VERIFY blocks"
  assert_output_contains "$out" "✗ verify"

  # add a VERIFY-LOG with a ✅ row → pass
  printf '# Verify Log\n| Check | Kind | Evidence | Result |\n|---|---|---|---|\n| M1 | run | ok | ✅ |\n' > "$R/VERIFY-LOG.md"
  (cd "$dir" && "$CLI" archive r --skip-doctor >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "walked VERIFY passes"
  rm -rf "$dir"
}

# ── Scenario 3: spec check + NONE delta ───────────────────────────────────────
scenario_3_spec() {
  echo "Scenario 3: spec check — missing SPEC-DELTA blocks; NONE passes"
  local dir="/tmp/spectacular-gate-3"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  tick_tasks "$R"

  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "missing SPEC-DELTA blocks"
  assert_output_contains "$out" "✗ spec"

  none_delta "$R"
  (cd "$dir" && "$CLI" archive r --skip-doctor >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "NONE delta passes"
  rm -rf "$dir"
}

# ── Scenario 4: override recording ────────────────────────────────────────────
scenario_4_override() {
  echo "Scenario 4: --override records into archive_overrides:"
  local dir="/tmp/spectacular-gate-4"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  # leave tasks open + no delta; override both
  local code
  (cd "$dir" && "$CLI" archive r --skip-doctor \
     --override tasks --reason "deferred to b23" \
     --override spec  --reason "no spec impact" >/dev/null 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "double override archives"
  local P="$dir/.spectacular/archive/r/PLAN.md"
  assert_file_contains "$P" "archive_overrides:"
  assert_file_contains "$P" "tasks — deferred to b23"
  assert_file_contains "$P" "spec — no spec impact"
  # exactly one archive_overrides: key (valid YAML, not duplicated)
  local n; n=$(grep -c '^archive_overrides:' "$P")
  assert_exit "$n" 1 "single archive_overrides key"

  # undo drops the whole archive_overrides block without orphaning items under related:
  (cd "$dir" && "$CLI" undo >/dev/null 2>&1)
  local RP="$dir/.spectacular/requests/r/PLAN.md"
  assert_file_lacks "$RP" "archive_overrides:"
  assert_file_lacks "$RP" "tasks — deferred to b23"
  assert_file_contains "$RP" "  - PRD.md"
  rm -rf "$dir"
}

# ── Scenario 5: no-op override warns, doesn't error ───────────────────────────
scenario_5_noop_override() {
  echo "Scenario 5: overriding a passing check warns but succeeds"
  local dir="/tmp/spectacular-gate-5"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  tick_tasks "$R"; none_delta "$R"   # both checks already pass
  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor --override tasks --reason "x" 2>&1) && code=0 || code=$?
  assert_exit "$code" 0 "no-op override still archives"
  assert_output_contains "$out" "already passes"
  # no archive_overrides stamped for a no-op
  assert_file_lacks "$dir/.spectacular/archive/r/PLAN.md" "archive_overrides:"
  rm -rf "$dir"
}

# ── Scenario 6: --force does NOT bypass closure checks ────────────────────────
scenario_6_force_split() {
  echo "Scenario 6: --force bypasses status gate only, not closure checks"
  local dir="/tmp/spectacular-gate-6"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)   # status: planned
  # planned + open tasks + no delta: --force clears status gate, closure still blocks
  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor --force 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "--force alone doesn't bypass closure gate"
  assert_output_contains "$out" "closure gate blocks"
  rm -rf "$dir"
}

# ── Scenario 7: grandfathered — already-archived never revalidated ────────────
scenario_7_grandfathered() {
  echo "Scenario 7: an already-archived request is refused, not revalidated"
  local dir="/tmp/spectacular-gate-7"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  (cd "$dir" && "$CLI" advance r --to review --force >/dev/null)
  tick_tasks "$R"; none_delta "$R"
  (cd "$dir" && "$CLI" archive r --skip-doctor >/dev/null 2>&1)
  assert_dir_exists "$dir/.spectacular/archive/r"
  # second archive of same slug → refused as already-archived (not a gate error)
  local out code
  out=$(cd "$dir" && "$CLI" archive r --skip-doctor 2>&1) && code=0 || code=$?
  assert_exit "$code" 1 "already-archived refused"
  assert_output_contains "$out" "already archived"
  rm -rf "$dir"
}

# ── Scenario 8: doctor specs flags a bad delta ────────────────────────────────
scenario_8_delta_integrity() {
  echo "Scenario 8: doctor specs flags a MODIFIED delta quoting a missing bullet"
  local dir="/tmp/spectacular-gate-8"; seed "$dir"
  local R="$dir/.spectacular/requests/r"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null)
  printf -- '- auth — password login\n' >> "$dir/.spectacular/SPEC.md"
  cat > "$R/SPEC-DELTA.md" <<'EOF'
### MODIFIED
- SPEC.md :: "auth — bullet that does not exist" -> "auth — sso"
EOF
  local out
  out=$(cd "$dir" && "$CLI" doctor specs 2>&1)
  assert_output_contains "$out" "not found in SPEC.md"
  # a valid delta (bullet exists) → no such warning
  cat > "$R/SPEC-DELTA.md" <<'EOF'
### MODIFIED
- SPEC.md :: "auth — password login" -> "auth — sso"
EOF
  out=$(cd "$dir" && "$CLI" doctor specs 2>&1)
  if echo "$out" | grep -qF "not found in SPEC.md"; then
    echo "    ✗ valid delta wrongly flagged"; fail_count=$((fail_count+1))
  else pass_count=$((pass_count+1)); fi
  rm -rf "$dir"
}

scenario_1_tasks
scenario_1b_deferred
scenario_2_verify
scenario_3_spec
scenario_4_override
scenario_5_noop_override
scenario_6_force_split
scenario_7_grandfathered
scenario_8_delta_integrity

echo ""
echo "Results: $pass_count passed, $fail_count failed"
[[ "$fail_count" -eq 0 ]]
