#!/usr/bin/env bash
# tests/cli/policy-output.test.sh — tiered policy gate output + directive field (b29 M1)
#
# Covers:
#   1. warn + directive    — row shows "— <directive>" + "P<n> — <title>" (no full paragraph)
#   2. warn, no directive  — falls back to principle title only
#   3. block + directive   — directive + FULL principle line
#   4. --full              — restores full paragraphs on warn rows too
#   5. --json              — carries a "directive" key
#   6. advance gate        — _policy_consult_transition inherits the tiering
#   7. repo POLICY.md      — every policy has an authored `- directive:` line

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"

fail_count=0
pass_count=0

assert_output_contains() { if echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count+1)); else echo "    ✗ expected output to contain: $2"; fail_count=$((fail_count+1)); fi; }
assert_output_lacks()    { if ! echo "$1" | grep -qF -- "$2"; then pass_count=$((pass_count+1)); else echo "    ✗ expected output to NOT contain: $2"; fail_count=$((fail_count+1)); fi; }
assert_eq() { if [[ "$1" == "$2" ]]; then pass_count=$((pass_count+1)); else echo "    ✗ $3: got '$1' want '$2'"; fail_count=$((fail_count+1)); fi; }

seed() {
  local dir="$1"
  rm -rf "$dir"; mkdir -p "$dir"
  (cd "$dir" && git init -q 2>/dev/null)
  (cd "$dir" && "$CLI" init --kit blank --name "$(basename "$dir")" >/dev/null 2>&1)

  cat > "$dir/.spectacular/PRINCIPLES.md" <<'EOF'
# Principles

## 7. Seven title
Seven body line that only full output shows.

## 10. Ten title
Ten body line that only full output shows.

## 11. Eleven title
Eleven body line that only full output shows.
EOF

  cat > "$dir/.spectacular/POLICY.md" <<'EOF'
---
version: 1.0
updated: 2026-07-12
summary: "fixture"
---

# Policies

## @Implementation

### with-directive
- principle: 11
- severity: warn
- check: fixture condition
- directive: Always test the tiered row.

Warn prose that never prints at the gate.

### no-directive
- principle: 10
- severity: warn

Fallback prose that never prints at the gate.

### block-directive
- principle: 7
- severity: block
- check: fixture block condition
- directive: Blocked row directive sentence.

Block prose that never prints at the gate.
EOF
}

# ── Scenario 1–4: hook form tiering ──────────────────────────────────────────
scenario_hook_form() {
  echo "Scenario 1-4: hook form — tiered rows, fallback, block=full, --full"
  local dir="/tmp/spectacular-policy-out-1"; seed "$dir"
  local out
  out=$(cd "$dir" && "$CLI" policy @Implementation 2>&1)

  # 1. warn + directive: directive line + title trailer, NOT the full body
  assert_output_contains "$out" "— Always test the tiered row."
  assert_output_contains "$out" "→ P11 — Eleven title"
  assert_output_lacks    "$out" "Eleven body line"

  # 2. warn without directive: title fallback only
  assert_output_contains "$out" "→ P10 — Ten title"
  assert_output_lacks    "$out" "Ten body line"

  # 3. block + directive: directive + FULL principle line
  assert_output_contains "$out" "— Blocked row directive sentence."
  assert_output_contains "$out" "→ P7. Seven title — Seven body line that only full output shows."

  # policy prose never leaks into the gate
  assert_output_lacks "$out" "never prints at the gate"

  # 4. --full restores full paragraphs on warn rows
  out=$(cd "$dir" && "$CLI" policy @Implementation --full 2>&1)
  assert_output_contains "$out" "→ P11. Eleven title — Eleven body line that only full output shows."
  assert_output_contains "$out" "— Always test the tiered row."
}

# ── Scenario 5: --json carries directive ─────────────────────────────────────
scenario_json() {
  echo "Scenario 5: --json carries directive key"
  local dir="/tmp/spectacular-policy-out-1"   # reuse seed from scenario 1
  [[ -d "$dir/.spectacular" ]] || seed "$dir"
  local out
  out=$(cd "$dir" && "$CLI" policy --json 2>&1)
  assert_output_contains "$out" '"id":"with-directive"'
  assert_output_contains "$out" '"directive":"Always test the tiered row."'
  assert_output_contains "$out" '"directive":""'   # no-directive policy → empty string, key present
}

# ── Scenario 6: advance gate inherits the tiering ────────────────────────────
scenario_advance() {
  echo "Scenario 6: advance planned→active embeds the tiered gate"
  local dir="/tmp/spectacular-policy-out-2"; seed "$dir"
  (cd "$dir" && "$CLI" new r --summary t >/dev/null 2>&1)
  local out
  out=$(cd "$dir" && "$CLI" advance r 2>&1)
  assert_output_contains "$out" "policy gate @Implementation"
  assert_output_contains "$out" "— Always test the tiered row."
  assert_output_contains "$out" "→ P11 — Eleven title"
  assert_output_lacks    "$out" "Eleven body line"
}

# ── Scenario 7: repo POLICY.md fully authored ────────────────────────────────
scenario_repo_authored() {
  echo "Scenario 7: every repo policy has a - directive: line"
  local pol="$REPO_ROOT/.spectacular/POLICY.md"
  local n_pol n_dir
  n_pol=$(grep -c '^### ' "$pol")
  n_dir=$(grep -c '^- directive:' "$pol")
  assert_eq "$n_dir" "$n_pol" "directive count matches policy count"
}

scenario_hook_form
scenario_json
scenario_advance
scenario_repo_authored

rm -rf /tmp/spectacular-policy-out-1 /tmp/spectacular-policy-out-2

echo ""
echo "policy-output: $pass_count passed, $fail_count failed"
[[ $fail_count -eq 0 ]] || exit 1
exit 0
