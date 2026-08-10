#!/usr/bin/env bash
set -euo pipefail

die() { echo "refused: $*" >&2; exit 3; }

[[ $# -eq 3 ]] || die "usage: stage-local.sh <codex|claude> <empty-disposable-root> <spectacular-binary>"

runtime_name="$1"
destination_root="$2"
binary_path="$3"
script_dir="$(cd "$(dirname "$0")" && pwd)"
plugin_root="$(cd "$script_dir/.." && pwd)"

case "$runtime_name" in codex|claude) ;; *) die "runtime must be codex or claude" ;; esac
[[ -d "$destination_root" ]] || die "destination must be an existing disposable directory"
destination_root="$(cd "$destination_root" && pwd -P)"
home_root=""
if [[ -n "${HOME:-}" && -d "$HOME" ]]; then
  home_root="$(cd "$HOME" && pwd -P)"
fi
[[ "$destination_root" != "/" ]] || die "destination cannot be filesystem root"
[[ -z "$home_root" || "$destination_root" != "$home_root" ]] || die "destination cannot be the user home"
[[ -z "$(find "$destination_root" -mindepth 1 -maxdepth 1 -print -quit)" ]] || die "destination must be empty"

"$script_dir/preflight.sh" "$runtime_name" "$plugin_root" "$binary_path" >/dev/null

mkdir -p "$destination_root/bin" "$destination_root/plugins/spectacular"
cp "$binary_path" "$destination_root/bin/spectacular"
cp -R "$plugin_root/skills" "$destination_root/plugins/spectacular/skills"

if [[ "$runtime_name" == "codex" ]]; then
  mkdir -p "$destination_root/plugins/spectacular/.codex-plugin"
  cp "$plugin_root/.codex-plugin/plugin.json" "$destination_root/plugins/spectacular/.codex-plugin/plugin.json"
else
  mkdir -p "$destination_root/plugins/spectacular/.claude-plugin"
  cp "$plugin_root/.claude-plugin/plugin.json" "$destination_root/plugins/spectacular/.claude-plugin/plugin.json"
fi

printf 'runtime=%s\nstaged_root=%s\nplugin=%s\nbinary=%s\nresult=staged-disposable-only\n' \
  "$runtime_name" "$destination_root" "$destination_root/plugins/spectacular" "$destination_root/bin/spectacular"
