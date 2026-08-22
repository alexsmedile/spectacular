#!/bin/sh
# orient.sh — locate a Spectacular workspace and report what is live.
#
# Read-only fallback for hosts without the `spectacular` CLI. Touches flat
# frontmatter fields only; anything requiring nested YAML belongs in the CLI or
# a TS helper.
set -eu

fail() { printf '%s\n' "$1" >&2; exit 1; }
script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
. "$script_dir/lib/cli-mode.sh"

# Walk up for .spectacular/, like the CLI does.
dir=$(pwd)
workspace=""
while [ "$dir" != "/" ]; do
  if [ -d "$dir/.spectacular" ]; then workspace="$dir/.spectacular"; break; fi
  dir=$(dirname "$dir")
done
[ -n "$workspace" ] || fail "no .spectacular/ workspace found from $(pwd)"

root=$(dirname "$workspace")
printf 'workspace: %s\n' "$workspace"

# Root Anchor
if [ -f "$workspace/PROJECT.md" ]; then
  title=$(sed -n 's/^title:[[:space:]]*//p' "$workspace/PROJECT.md" | head -1)
  [ -n "$title" ] && printf 'project:   %s\n' "$title"
else
  printf 'project:   (no PROJECT.md — greenfield workspace)\n'
fi

# Exact schema and release compatibility decide which mode the caller is in.
spectacular_cli_probe "$script_dir/../generated/mechanical-interface.json"
if [ "$SPECTACULAR_CLI_MODE" = full ]; then
  printf 'cli:       %s\n' "$SPECTACULAR_CLI_DETAIL"
else
  printf 'cli:       %s; read-only mode\n' "$SPECTACULAR_CLI_DETAIL"
fi

printf '\nMISSIONS\n'
found=0
for d in "$workspace"/missions/*/; do
  [ -d "$d" ] || continue
  ref=$(basename "$d")
  record="$d$ref.md"
  [ -f "$record" ] || record=$(ls "$d"*.md 2>/dev/null | grep -v '/index.md$' | head -1)
  [ -n "$record" ] && [ -f "$record" ] || continue
  found=1
  status=$(sed -n 's/^status:[[:space:]]*//p' "$record" | head -1)
  title=$(sed -n 's/^title:[[:space:]]*//p' "$record" | head -1)
  printf '  %-12s %-12s %s\n' "${ref%%-*}" "${status:-unknown}" "$title"
done
[ "$found" -eq 1 ] || printf '  (none)\n'

# The live Mission is the one not completed or resolved.
printf '\nLIVE\n'
live=0
for d in "$workspace"/missions/*/; do
  [ -d "$d" ] || continue
  ref=$(basename "$d")
  record="$d$ref.md"
  [ -f "$record" ] || continue
  status=$(sed -n 's/^status:[[:space:]]*//p' "$record" | head -1)
  case "$status" in
    completed|resolved|archived) ;;
    *) printf '  %s (%s)\n' "${ref%%-*}" "${status:-unknown}"; live=1 ;;
  esac
done
[ "$live" -eq 1 ] || printf '  nothing live — no Mission is currently active\n'

printf '\nOPEN PROPOSALS\n'
open=0
for f in "$workspace"/proposals/P*.md; do
  [ -f "$f" ] || continue
  status=$(sed -n 's/^status:[[:space:]]*//p' "$f" | head -1)
  case "$status" in
    accepted|rejected|withdrawn) ;;
    *) printf '  %-6s %s\n' "$(basename "$f" .md | cut -d- -f1)" "$(sed -n 's/^title:[[:space:]]*//p' "$f" | head -1)"; open=1 ;;
  esac
done
[ "$open" -eq 1 ] || printf '  (none)\n'

printf '\nreported from files only; fingerprints and bindings were not verified\n'
