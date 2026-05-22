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
  assert_dir_exists "$dir/.spectacular/current"

  assert_file_absent "$dir/.spectacular/PRINCIPLES.md"
  assert_file_absent "$dir/.spectacular/ARCHITECTURE.md"
  assert_file_absent "$dir/.spectacular/ROADMAP.md"
  assert_file_absent "$dir/.spectacular/STACK.md"
  assert_file_absent "$dir/.spectacular/DECISIONS.md"

  assert_file_contains "$dir/.spectacular/PRD.md" "kit: blank"
  assert_file_contains "$dir/.spectacular/config.yaml" "kit: blank"
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

# ── run ───────────────────────────────────────────────────────────────────────

scenario_1_bare_init
scenario_2_kit_coding
scenario_3_interactive
scenario_4_idempotent_rerun
scenario_5_with_flag
scenario_6_minimal_overrides_kit
scenario_7_interactive_abort
scenario_8_global_flag

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"

if [[ $fail_count -gt 0 ]]; then
  exit 1
fi
exit 0
