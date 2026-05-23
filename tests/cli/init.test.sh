#!/usr/bin/env bash
# tests/cli/init.test.sh — covers the 6 VERIFY.md scenarios for smart-init.
#
# Each scenario runs the CLI against an isolated /tmp dir + asserts file presence.
# Skill install step (curl + tarball fetch) is bypassed by pre-creating a stub
# .agents/skills/spectacular/ symlink to this repo's local skill — keeps tests
# offline + fast.

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
LOCAL_SKILL="$REPO_ROOT/skills/spectacular"

# ── helpers ───────────────────────────────────────────────────────────────────
fail_count=0
pass_count=0

assert_file_exists() {
  local path="$1"
  if [[ -f "$path" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected file: $path"
    fail_count=$((fail_count + 1))
  fi
}

assert_file_absent() {
  local path="$1"
  if [[ ! -e "$path" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected absent: $path"
    fail_count=$((fail_count + 1))
  fi
}

assert_dir_exists() {
  local path="$1"
  if [[ -d "$path" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected dir: $path"
    fail_count=$((fail_count + 1))
  fi
}

assert_file_contains() {
  local path="$1" pattern="$2"
  if [[ -f "$path" ]] && grep -qF "$pattern" "$path"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected '$path' to contain '$pattern'"
    fail_count=$((fail_count + 1))
  fi
}

# Bypass skill download: pre-create the local skill symlink so install_skill()
# detects the existing directory and skips fetching from GitHub.
seed_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.agents/skills"
  ln -s "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"
  # Pre-write a skills.lock so install_skill is short-circuited
  mkdir -p "$dir/.spectacular"
  cat > "$dir/.spectacular/skills.lock" <<EOF
spectacular:
  ref: local-dev
  sha: local
  url: file://${LOCAL_SKILL}
EOF
}

run_cli() {
  local dir="$1"; shift
  (cd "$dir" && "$CLI" init "$@" 2>&1)
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_bare_init() {
  echo "Scenario 1: bare init produces only the always-set"
  local dir="/tmp/spectacular-test-1"
  seed_workspace "$dir"

  run_cli "$dir" >/dev/null

  assert_file_exists "$dir/.spectacular/PRD.md"
  assert_file_exists "$dir/.spectacular/config.yaml"
  assert_file_exists "$dir/.spectacular/AGENTS.md"
  assert_dir_exists "$dir/.spectacular/requests"
  assert_dir_exists "$dir/.spectacular/specs"
  assert_file_exists "$dir/.spectacular/SPEC.md"

  assert_file_absent "$dir/.spectacular/PRINCIPLES.md"
  assert_file_absent "$dir/.spectacular/ARCHITECTURE.md"
  assert_file_absent "$dir/.spectacular/ROADMAP.md"
  assert_file_absent "$dir/.spectacular/STACK.md"
  assert_file_absent "$dir/.spectacular/DECISIONS.md"

  assert_file_contains "$dir/.spectacular/PRD.md" "kit: blank"
  assert_file_contains "$dir/.spectacular/config.yaml" "kit: blank"
  assert_file_contains "$dir/.spectacular/config.yaml" 'workspace_schema: "0.6"'
  assert_file_contains "$dir/.gitignore" ".spectacular.local/"

  rm -rf "$dir"
}

scenario_2_kit_coding() {
  echo "Scenario 2: --kit coding adds STACK + ARCHITECTURE"
  local dir="/tmp/spectacular-test-2"
  seed_workspace "$dir"

  run_cli "$dir" --kit coding >/dev/null

  assert_file_exists "$dir/.spectacular/PRD.md"
  assert_file_exists "$dir/.spectacular/config.yaml"
  assert_file_exists "$dir/.spectacular/AGENTS.md"

  assert_file_exists "$dir/.spectacular/STACK.md"
  assert_file_exists "$dir/.spectacular/ARCHITECTURE.md"

  assert_file_absent "$dir/.spectacular/PRINCIPLES.md"
  assert_file_absent "$dir/.spectacular/ROADMAP.md"
  assert_file_absent "$dir/.spectacular/DECISIONS.md"

  assert_file_contains "$dir/.spectacular/PRD.md" "kit: coding"

  rm -rf "$dir"
}

scenario_3_interactive() {
  echo "Scenario 3: -i interactive walks kit + suggested prompts"
  local dir="/tmp/spectacular-test-3"
  seed_workspace "$dir"

  # Interactive prompts (in order):
  #   name (default), summary (empty), agents-file (default), scope (default),
  #   kit=2 (coding), then per-suggested-doc y/n, then final "Proceed?" confirm.
  # Coding kit suggested = principles, roadmap, decisions → answer y / n / y, then y to confirm
  (cd "$dir" && printf '\n\n\n\n2\ny\nn\ny\ny\n' | "$CLI" init -i 2>&1) >/dev/null || true

  assert_file_exists "$dir/.spectacular/STACK.md"
  assert_file_exists "$dir/.spectacular/ARCHITECTURE.md"
  assert_file_exists "$dir/.spectacular/PRINCIPLES.md"
  assert_file_absent "$dir/.spectacular/ROADMAP.md"
  assert_file_exists "$dir/.spectacular/DECISIONS.md"

  rm -rf "$dir"
}

scenario_4_idempotent_rerun() {
  echo "Scenario 4: idempotent re-run skips existing files"
  local dir="/tmp/spectacular-test-4"
  seed_workspace "$dir"

  run_cli "$dir" --kit coding >/dev/null

  echo "USER ADDED CONTENT" >> "$dir/.spectacular/PRD.md"
  local before
  before=$(cat "$dir/.spectacular/PRD.md")

  local out
  out=$(run_cli "$dir" --kit coding)

  local after
  after=$(cat "$dir/.spectacular/PRD.md")
  if [[ "$before" == "$after" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ PRD.md was modified on re-run"
    fail_count=$((fail_count + 1))
  fi

  if echo "$out" | grep -q "PRD.md"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ stdout did not mention PRD.md"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$dir"
}

scenario_5_with_flag() {
  echo "Scenario 5: --with principles,roadmap adds exactly those two"
  local dir="/tmp/spectacular-test-5"
  seed_workspace "$dir"

  run_cli "$dir" --with principles,roadmap >/dev/null

  assert_file_exists "$dir/.spectacular/PRINCIPLES.md"
  assert_file_exists "$dir/.spectacular/ROADMAP.md"

  assert_file_absent "$dir/.spectacular/STACK.md"
  assert_file_absent "$dir/.spectacular/ARCHITECTURE.md"
  assert_file_absent "$dir/.spectacular/DECISIONS.md"

  # Unknown doc ID errors cleanly
  seed_workspace "$dir"
  if run_cli "$dir" --with FOOBAR >/dev/null 2>&1; then
    echo "    ✗ --with FOOBAR should error"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi

  rm -rf "$dir"
}

scenario_6_minimal_overrides_kit() {
  echo "Scenario 6: --kit coding --minimal scaffolds always-set only"
  local dir="/tmp/spectacular-test-6"
  seed_workspace "$dir"

  run_cli "$dir" --kit coding --minimal >/dev/null

  assert_file_exists "$dir/.spectacular/PRD.md"
  assert_file_exists "$dir/.spectacular/config.yaml"
  assert_file_exists "$dir/.spectacular/AGENTS.md"

  assert_file_absent "$dir/.spectacular/STACK.md"
  assert_file_absent "$dir/.spectacular/ARCHITECTURE.md"

  assert_file_contains "$dir/.spectacular/PRD.md" "kit: coding"

  rm -rf "$dir"
}

scenario_7_interactive_abort() {
  echo "Scenario 7: -i with 'n' at Proceed? aborts without writing"
  local dir="/tmp/spectacular-test-7"
  seed_workspace "$dir"

  # Same prompt sequence as scenario 3, but final 'n' to abort
  (cd "$dir" && printf '\n\n\n\n2\ny\nn\ny\nn\n' | "$CLI" init -i 2>&1) >/dev/null || true

  # No canonical doc files should exist
  assert_file_absent "$dir/.spectacular/PRD.md"
  assert_file_absent "$dir/.spectacular/config.yaml"
  assert_file_absent "$dir/.spectacular/AGENTS.md"
  assert_file_absent "$dir/.spectacular/STACK.md"
  assert_file_absent "$dir/.spectacular/PRINCIPLES.md"

  rm -rf "$dir"
}

scenario_8_global_flag() {
  echo "Scenario 8: --global uses fake HOME (no real ~/.agents pollution)"
  local fake_home="/tmp/spectacular-test-8-home"
  rm -rf "$fake_home"
  mkdir -p "$fake_home/.agents/skills" "$fake_home/project/.spectacular"
  ln -s "$LOCAL_SKILL" "$fake_home/.agents/skills/spectacular"
  cat > "$fake_home/project/.spectacular/skills.lock" <<EOF
spectacular:
  ref: local-dev
EOF

  (cd "$fake_home/project" && HOME="$fake_home" "$CLI" init --global) >/dev/null

  # Workspace scaffolded in project dir
  assert_file_exists "$fake_home/project/.spectacular/PRD.md"
  assert_file_exists "$fake_home/project/.spectacular/config.yaml"
  assert_file_exists "$fake_home/project/.spectacular/AGENTS.md"

  # Skill resolves under fake HOME (symlink seeded by us)
  if [[ -L "$fake_home/.agents/skills/spectacular" ]]; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ skill symlink missing under \$HOME/.agents/"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$fake_home"
}

scenario_9_flag_eat() {
  echo "Scenario 9: value-flags reject another flag as their value"
  local dir="/tmp/spectacular-test-9"
  rm -rf "$dir"; mkdir -p "$dir"
  local out
  out=$(cd "$dir" && "$CLI" init --name --kit coding 2>&1) || true
  if echo "$out" | grep -q "requires a value"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected 'requires a value' error, got: $out"
    fail_count=$((fail_count + 1))
  fi
  out=$(cd "$dir" && "$CLI" doctor --format 2>&1) || true
  if echo "$out" | grep -q "format requires a value"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected --format guard, got: $out"
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$dir"
}

scenario_10_skill_verb_stubs() {
  echo "Scenario 10: remaining skill-only verbs print interactive-flow message"
  # As of v0.8.0, archive/new/snapshot/promote moved to CLI mutator verbs.
  # Only status + remember remain as skill stubs.
  local out code
  for verb in status remember; do
    out=$("$CLI" "$verb" 2>&1) && code=0 || code=$?
    if [[ "$code" -ne 0 ]] && echo "$out" | grep -q "interactive skill flow"; then
      pass_count=$((pass_count + 1))
    else
      echo "    ✗ verb '$verb' did not produce expected stub message"
      fail_count=$((fail_count + 1))
    fi
  done
}

scenario_11_status_against_latest() {
  echo "Scenario 11: status --against-latest reports workspace_schema verdict"
  local dir="/tmp/spectacular-test-11"
  seed_workspace "$dir"
  run_cli "$dir" --kit blank >/dev/null

  # Fresh init → up to date
  local out_fresh code_fresh
  out_fresh=$(cd "$dir" && "$CLI" status --against-latest 2>&1)
  code_fresh=$?
  if [[ "$code_fresh" -eq 0 ]] && echo "$out_fresh" | grep -q "up to date"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ fresh init: expected 'up to date', got: $out_fresh"
    fail_count=$((fail_count + 1))
  fi

  # Simulated older workspace (strip the field) → behind
  sed -i.bak '/^workspace_schema:/d' "$dir/.spectacular/config.yaml"
  rm -f "$dir/.spectacular/config.yaml.bak"
  local out_old code_old
  out_old=$(cd "$dir" && "$CLI" status --against-latest 2>&1)
  code_old=$?
  if [[ "$code_old" -eq 0 ]] && echo "$out_old" | grep -q "behind CLI" && echo "$out_old" | grep -q "spectacular migrate"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ stripped field: expected 'behind CLI' + migrate suggestion, got: $out_old"
    fail_count=$((fail_count + 1))
  fi

  # Outside a workspace → exit 1
  local out_outside code_outside
  out_outside=$(cd /tmp && "$CLI" status --against-latest 2>&1) && code_outside=0 || code_outside=$?
  if [[ "$code_outside" -eq 1 ]] && echo "$out_outside" | grep -q "Not inside"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ outside workspace: expected exit 1 + 'Not inside', got: $out_outside (exit $code_outside)"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$dir"
}

# ── run ───────────────────────────────────────────────────────────────────────

scenario_1_bare_init
scenario_2_kit_coding
scenario_3_interactive
scenario_4_idempotent_rerun
scenario_5_with_flag
scenario_6_minimal_overrides_kit
scenario_7_interactive_abort
scenario_8_global_flag
scenario_9_flag_eat
scenario_10_skill_verb_stubs
scenario_11_status_against_latest

scenario_12_status_since() {
  echo "Scenario 12: status --since filters by frontmatter updated: field"
  local dir="/tmp/spectacular-test-12"
  seed_workspace "$dir"
  run_cli "$dir" --kit blank >/dev/null

  # Create requests with controlled updated: dates
  (cd "$dir" && "$CLI" new req-old --summary "old" >/dev/null)
  (cd "$dir" && "$CLI" new req-recent --summary "recent" >/dev/null)
  (cd "$dir" && "$CLI" promote req-recent --to active >/dev/null)

  # Backdate req-old's PLAN to long ago
  sed -i.bak 's/^updated:.*/updated: 2025-01-01/' "$dir/.spectacular/requests/req-old/PLAN.md"
  rm -f "$dir/.spectacular/requests/req-old/PLAN.md.bak"

  # YYYY-MM-DD form should include only req-recent (today >= 2026-01-01 > 2025-01-01)
  local out_abs code_abs
  out_abs=$(cd "$dir" && "$CLI" status --since 2026-01-01 2>&1)
  code_abs=$?
  if [[ "$code_abs" -eq 0 ]] && echo "$out_abs" | grep -q "req-recent" && ! echo "$out_abs" | grep -q "req-old"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ --since 2026-01-01: expected req-recent only, got: $out_abs"
    fail_count=$((fail_count + 1))
  fi

  # Relative form 7d should also include recent but not 2025-vintage
  local out_rel
  out_rel=$(cd "$dir" && "$CLI" status --since 7d 2>&1)
  if echo "$out_rel" | grep -q "req-recent" && ! echo "$out_rel" | grep -q "req-old"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ --since 7d: expected req-recent only"
    fail_count=$((fail_count + 1))
  fi

  # --since= form
  local out_eq
  out_eq=$(cd "$dir" && "$CLI" status --since=2026-01-01 2>&1)
  if echo "$out_eq" | grep -q "req-recent"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ --since= form failed"
    fail_count=$((fail_count + 1))
  fi

  # Empty bucket: future date returns 0 requests + 0 docs
  local out_empty
  out_empty=$(cd "$dir" && "$CLI" status --since 9999-01-01 2>&1)
  if echo "$out_empty" | grep -q "Requests (0)" && echo "$out_empty" | grep -q "Canonical docs (0)"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ future date should yield empty buckets, got: $out_empty"
    fail_count=$((fail_count + 1))
  fi

  # Missing arg -> exit 2
  local out_no code_no
  out_no=$(cd "$dir" && "$CLI" status --since 2>&1) && code_no=0 || code_no=$?
  if [[ "$code_no" -eq 2 ]] && echo "$out_no" | grep -q "argument required"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ missing arg: expected exit 2 + 'argument required', got: $out_no (exit $code_no)"
    fail_count=$((fail_count + 1))
  fi

  # Bad format -> exit 2
  local out_bad code_bad
  out_bad=$(cd "$dir" && "$CLI" status --since "notadate" 2>&1) && code_bad=0 || code_bad=$?
  if [[ "$code_bad" -eq 2 ]] && echo "$out_bad" | grep -q "cannot parse"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ bad format: expected exit 2 + 'cannot parse', got: $out_bad (exit $code_bad)"
    fail_count=$((fail_count + 1))
  fi

  # Outside workspace -> exit 1
  local out_outside code_outside
  out_outside=$(cd /tmp && "$CLI" status --since 7d 2>&1) && code_outside=0 || code_outside=$?
  if [[ "$code_outside" -eq 1 ]] && echo "$out_outside" | grep -q "Not inside"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ outside workspace: expected exit 1, got: $out_outside (exit $code_outside)"
    fail_count=$((fail_count + 1))
  fi

  rm -rf "$dir"
}

scenario_12_status_since

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"

if [[ $fail_count -gt 0 ]]; then
  exit 1
fi
exit 0
