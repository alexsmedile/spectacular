#!/usr/bin/env bash
set -euo pipefail

die() { echo "refused: $*" >&2; exit 3; }

[[ $# -eq 3 ]] || die "usage: preflight.sh <codex|claude> <plugin-root> <spectacular-binary>"

runtime_name="$1"
plugin_root="$2"
binary_path="$3"
version_path="$plugin_root/VERSION"

case "$runtime_name" in
  codex) manifest_path="$plugin_root/.codex-plugin/plugin.json" ;;
  claude) manifest_path="$plugin_root/.claude-plugin/plugin.json" ;;
  *) die "runtime must be codex or claude" ;;
esac

[[ -f "$manifest_path" ]] || die "runtime manifest is missing"
[[ -f "$plugin_root/skills/spectacular/SKILL.md" ]] || die "canonical v2 Skill is missing"
[[ -f "$plugin_root/skills/spectacular/generated/mechanical-interface.json" ]] || die "generated mechanical catalog is missing"
[[ -f "$version_path" ]] || die "release version is missing"
[[ -x "$binary_path" ]] || die "local spectacular binary is missing or not executable"

release_version="$(sed -n '1p' "$version_path")"
[[ -n "$release_version" && "$release_version" != *[[:space:]]* ]] || die "release version is invalid"

grep -q '"name": "spectacular"' "$manifest_path" || die "runtime manifest identity is invalid"
grep -q '"version": "'"$release_version"'"' "$manifest_path" || die "runtime manifest version is invalid"
grep -q 'name: spectacular' "$plugin_root/skills/spectacular/SKILL.md" || die "Skill identity is invalid"
grep -q '^version: '"$release_version"'$' "$plugin_root/skills/spectacular/SKILL.md" || die "Skill version is invalid"
grep -q 'spectacular.command-catalog.v1' "$plugin_root/skills/spectacular/generated/mechanical-interface.json" || die "mechanical catalog schema is invalid"
grep -q '"release_version": "'"$release_version"'"' "$plugin_root/skills/spectacular/generated/mechanical-interface.json" || die "mechanical catalog version is invalid"
[[ "$($binary_path --version)" == "spectacular $release_version" ]] || die "binary version is invalid"

if grep -R -q 'spectacular status\|spectacular new\|requests/' "$plugin_root/skills/spectacular"; then
  die "v1 runtime language leaked into the v2 Skill bundle"
fi

printf 'runtime=%s\nplugin_root=%s\nbinary=%s\nresult=ready-for-disposable-staging\n' "$runtime_name" "$plugin_root" "$binary_path"
