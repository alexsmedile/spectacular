#!/usr/bin/env bash
# tests/cli/pack.test.sh — convention-pack lifecycle + init/doctor wiring tests.
#
# Covers VERIFY scenarios for convention-pack-application:
#   - pack list / install / remove / show across scopes
#   - init consumes config.yaml convention_pack: in scaffold mode
#   - doctor conventions area: suggest / scaffold / enforce mode behavior
#   - mechanical fix for pack-driven gitignore drift

set -u

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
LOCAL_SKILL="$REPO_ROOT/skills/spectacular"

# Use a fake HOME so pack install tests don't touch the user's real ~/.spectacular/
FAKE_HOME="/tmp/spectacular-pack-test-home"

# ── helpers ───────────────────────────────────────────────────────────────────
fail_count=0
pass_count=0

assert_file_exists() {
  local path="$1"
  if [[ -f "$path" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected file: $path"; fail_count=$((fail_count + 1)); fi
}

assert_dir_exists() {
  local path="$1"
  if [[ -d "$path" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected dir: $path"; fail_count=$((fail_count + 1)); fi
}

assert_absent() {
  local path="$1"
  if [[ ! -e "$path" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected absent: $path"; fail_count=$((fail_count + 1)); fi
}

assert_file_contains() {
  local path="$1" pattern="$2"
  if [[ -f "$path" ]] && grep -qF -- "$pattern" "$path"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected '$path' to contain '$pattern'"
    fail_count=$((fail_count + 1))
  fi
}

assert_stdout_contains() {
  local output="$1" pattern="$2"
  if echo "$output" | grep -qF -- "$pattern"; then
    pass_count=$((pass_count + 1))
  else
    echo "    ✗ expected stdout to contain '$pattern'"
    echo "      got: $(echo "$output" | head -3)"
    fail_count=$((fail_count + 1))
  fi
}

assert_exit_code() {
  local expected="$1" actual="$2"
  if [[ "$expected" == "$actual" ]]; then pass_count=$((pass_count + 1))
  else echo "    ✗ expected exit $expected, got $actual"; fail_count=$((fail_count + 1)); fi
}

seed_workspace() {
  local dir="$1"
  rm -rf "$dir"
  mkdir -p "$dir/.agents/skills"
  ln -s "$LOCAL_SKILL" "$dir/.agents/skills/spectacular"
  mkdir -p "$dir/.spectacular"
  cat > "$dir/.spectacular/skills.lock" <<EOF
spectacular:
  ref: local-dev
  sha: local
  url: file://${LOCAL_SKILL}
EOF
}

reset_fake_home() {
  rm -rf "$FAKE_HOME"
  mkdir -p "$FAKE_HOME"
}

# ── scenarios ─────────────────────────────────────────────────────────────────

scenario_1_pack_list_shows_bundled_and_appstore() {
  echo "Scenario 1: pack list shows bundled + app-store scopes"
  reset_fake_home
  local out
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack list 2>&1)

  assert_stdout_contains "$out" "[bundled]"
  assert_stdout_contains "$out" "minimal"
  assert_stdout_contains "$out" "[app-store]"
  assert_stdout_contains "$out" "alex-default"
}

scenario_2_pack_install_minimal() {
  echo "Scenario 2: pack install minimal copies to user scope"
  reset_fake_home

  local out exit_code
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack install minimal 2>&1)
  exit_code=$?

  assert_exit_code 0 "$exit_code"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/minimal/pack.md"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/minimal/templates/.gitignore"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/minimal/references/why-minimal.md"
  assert_stdout_contains "$out" "✓ installed minimal"
  assert_stdout_contains "$out" "convention_pack:"
}

scenario_3_pack_install_already_exists_fails() {
  echo "Scenario 3: re-installing same pack errors"
  # depends on scenario 2 having installed minimal
  local out exit_code
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack install minimal 2>&1)
  exit_code=$?

  assert_exit_code 1 "$exit_code"
  assert_stdout_contains "$out" "already installed"
}

scenario_4_pack_show_minimal() {
  echo "Scenario 4: pack show prints frontmatter"
  local out
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack show minimal 2>&1)

  assert_stdout_contains "$out" "pack:   minimal"
  assert_stdout_contains "$out" "scope:"
  assert_stdout_contains "$out" "pack: minimal"
  assert_stdout_contains "$out" "version: 1.0"
}

scenario_5_pack_remove_user_scope() {
  echo "Scenario 5: pack remove deletes user-scope pack"
  local out exit_code
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack remove minimal 2>&1)
  exit_code=$?

  assert_exit_code 0 "$exit_code"
  assert_absent "$FAKE_HOME/.spectacular/packs/minimal"
  assert_stdout_contains "$out" "removed user-scope pack"
}

scenario_6_pack_remove_bundled_refuses_without_force() {
  echo "Scenario 6: pack remove refuses bundled scope without --force"
  reset_fake_home
  local out exit_code
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack remove minimal 2>&1)
  exit_code=$?

  # bundled minimal still exists; should be refused
  assert_exit_code 1 "$exit_code"
  assert_stdout_contains "$out" "Use --force"
}

scenario_7_init_consumes_convention_pack() {
  echo "Scenario 7: init with convention_pack: scaffold appends pack gitignore"
  local dir="/tmp/spectacular-pack-test-7"
  seed_workspace "$dir"
  reset_fake_home

  # First init to scaffold config.yaml
  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt7 --kit blank --minimal >/dev/null 2>&1)

  # Inject convention_pack stanza
  cat >> "$dir/.spectacular/config.yaml" <<EOF

convention_pack:
  source: alex-default
  mode: scaffold
EOF

  # Re-run init (idempotent + applies pack)
  local out
  out=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt7 --kit blank --minimal 2>&1)

  assert_file_contains "$dir/.gitignore" "_archive/"
  assert_file_contains "$dir/.gitignore" "_backups/"
  assert_file_contains "$dir/.gitignore" ".env.local"
  assert_file_contains "$dir/.gitignore" ".spectacular.local/"

  rm -rf "$dir"
}

scenario_8_doctor_conventions_no_pack_declared() {
  echo "Scenario 8: doctor conventions reports skip when no pack declared"
  local dir="/tmp/spectacular-pack-test-8"
  seed_workspace "$dir"
  reset_fake_home

  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt8 --kit blank --minimal >/dev/null 2>&1)

  local out exit_code
  out=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" doctor conventions 2>&1)
  exit_code=$?

  assert_exit_code 0 "$exit_code"
  assert_stdout_contains "$out" "no convention_pack declared"

  rm -rf "$dir"
}

scenario_9_doctor_enforce_mode_flags_errors() {
  echo "Scenario 9: doctor conventions in enforce mode flags missing gitignore as errors"
  local dir="/tmp/spectacular-pack-test-9"
  seed_workspace "$dir"
  reset_fake_home

  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt9 --kit blank --minimal >/dev/null 2>&1)

  cat >> "$dir/.spectacular/config.yaml" <<EOF

convention_pack:
  source: alex-default
  mode: enforce
EOF
  # Truncate gitignore so most pack entries are missing
  echo ".spectacular.local/" > "$dir/.gitignore"

  local out exit_code
  out=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" doctor conventions 2>&1)
  exit_code=$?

  assert_exit_code 2 "$exit_code"
  assert_stdout_contains "$out" "requires gitignore entry"
  assert_stdout_contains "$out" "error"

  rm -rf "$dir"
}

scenario_10_doctor_fix_repairs_pack_drift() {
  echo "Scenario 10: doctor --fix appends missing pack gitignore entries"
  local dir="/tmp/spectacular-pack-test-10"
  seed_workspace "$dir"
  reset_fake_home

  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt10 --kit blank --minimal >/dev/null 2>&1)

  cat >> "$dir/.spectacular/config.yaml" <<EOF

convention_pack:
  source: alex-default
  mode: enforce
EOF
  echo ".spectacular.local/" > "$dir/.gitignore"

  local out exit_code
  out=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" doctor conventions --fix 2>&1)
  exit_code=$?

  # After --fix, should exit 0 with all entries appended
  assert_exit_code 0 "$exit_code"
  assert_file_contains "$dir/.gitignore" "_archive/"
  assert_file_contains "$dir/.gitignore" "_backups/"
  assert_file_contains "$dir/.gitignore" ".env.local"
  assert_stdout_contains "$out" "✓ fixed [conventions]"

  rm -rf "$dir"
}

scenario_11_pack_install_from_local_path() {
  echo "Scenario 11: pack install --from <local path> installs custom pack"
  reset_fake_home
  # Use the alex-default app-store entry as a synthetic --from source
  local src="$REPO_ROOT/packs/alex-default"
  local out exit_code
  out=$(cd "$REPO_ROOT" && HOME="$FAKE_HOME" "$CLI" pack install alex-default --from "$src" 2>&1)
  exit_code=$?

  assert_exit_code 0 "$exit_code"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/alex-default/pack.md"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/alex-default/templates/.gitignore"
  assert_file_exists "$FAKE_HOME/.spectacular/packs/alex-default/references/why-alex-default.md"
}

scenario_12_pack_help() {
  echo "Scenario 12: pack --help shows usage"
  local out exit_code
  out=$("$CLI" pack --help 2>&1)
  exit_code=$?

  assert_exit_code 0 "$exit_code"
  assert_stdout_contains "$out" "Usage: spectacular pack"
  assert_stdout_contains "$out" "install <name>"
  assert_stdout_contains "$out" "remove <name>"
  assert_stdout_contains "$out" "Pack scopes"
}

scenario_13_unknown_mode_fallback() {
  echo "Scenario 13: unknown convention_pack.mode value falls back to suggest + emits info note"
  local dir="/tmp/spectacular-pack-test-13"
  seed_workspace "$dir"
  reset_fake_home

  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" init --name pt13 --kit blank --minimal >/dev/null 2>&1)
  (cd "$dir" && HOME="$FAKE_HOME" "$CLI" pack install minimal >/dev/null 2>&1)

  # Unknown mode 'strict'
  cat >> "$dir/.spectacular/config.yaml" <<EOF

convention_pack:
  source: minimal
  mode: strict
EOF

  local out
  out=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" doctor conventions 2>&1)

  # Info line surfaces the unknown value + the fallback
  assert_stdout_contains "$out" "unknown convention_pack.mode 'strict'"
  assert_stdout_contains "$out" "falling back to 'suggest'"
  # Pack still resolves and reports as suggest mode
  assert_stdout_contains "$out" "active pack 'minimal' (suggest)"

  # Valid mode 'enforce' → no fallback info line
  sed -i.bak 's/mode: strict/mode: enforce/' "$dir/.spectacular/config.yaml"
  rm -f "$dir/.spectacular/config.yaml.bak"
  local out_valid
  out_valid=$(cd "$dir" && HOME="$FAKE_HOME" "$CLI" doctor conventions 2>&1)
  if echo "$out_valid" | grep -qF "unknown convention_pack.mode"; then
    echo "    ✗ valid 'enforce' mode should not trigger unknown-mode info"
    fail_count=$((fail_count + 1))
  else
    pass_count=$((pass_count + 1))
  fi

  rm -rf "$dir"
}

# ── run ───────────────────────────────────────────────────────────────────────
scenario_1_pack_list_shows_bundled_and_appstore
scenario_2_pack_install_minimal
scenario_3_pack_install_already_exists_fails
scenario_4_pack_show_minimal
scenario_5_pack_remove_user_scope
scenario_6_pack_remove_bundled_refuses_without_force
scenario_7_init_consumes_convention_pack
scenario_8_doctor_conventions_no_pack_declared
scenario_9_doctor_enforce_mode_flags_errors
scenario_10_doctor_fix_repairs_pack_drift
scenario_11_pack_install_from_local_path
scenario_12_pack_help
scenario_13_unknown_mode_fallback

# Cleanup fake HOME
rm -rf "$FAKE_HOME"

echo ""
echo "  Asserts: $pass_count passed, $fail_count failed"
if [[ $fail_count -gt 0 ]]; then
  echo "  ✗ tests/cli/pack.test.sh"
  exit 1
else
  echo "  ✓ tests/cli/pack.test.sh"
  exit 0
fi
