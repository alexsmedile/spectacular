#!/bin/sh
# doctor.sh — report what this environment can and cannot do.
#
# Read-only fallback for hosts without the `spectacular` CLI. Answers one
# question: which mode is the session in, and what is unavailable in it.
set -eu

ok=0
warn=0

say()  { printf '  %-10s %s\n' "$1" "$2"; }
script_dir=$(CDPATH= cd "$(dirname "$0")" && pwd)
. "$script_dir/lib/cli-mode.sh"

printf 'SPECTACULAR ENVIRONMENT\n\n'

# 1. The mechanical layer.
spectacular_cli_probe "$script_dir/../generated/mechanical-interface.json"
if [ "$SPECTACULAR_CLI_MODE" = full ]; then
  say "cli" "$SPECTACULAR_CLI_DETAIL"
  ok=$((ok+1))
  cli=1
else
  say "cli" "$SPECTACULAR_CLI_DETAIL"
  warn=$((warn+1))
  cli=0
fi

# 2. The workspace.
dir=$(pwd); workspace=""
while [ "$dir" != "/" ]; do
  if [ -d "$dir/.spectacular" ]; then workspace="$dir/.spectacular"; break; fi
  dir=$(dirname "$dir")
done
if [ -n "$workspace" ]; then
  say "workspace" "$workspace"
  ok=$((ok+1))
else
  say "workspace" "none found — not a Spectacular workspace"
  warn=$((warn+1))
fi

# 3. Git, which every binding depends on.
if command -v git >/dev/null 2>&1 && git rev-parse --git-dir >/dev/null 2>&1; then
  branch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo '?')
  if [ -n "$(git status --porcelain 2>/dev/null)" ]; then
    say "git" "$branch (uncommitted changes)"
  else
    say "git" "$branch (clean)"
  fi
  case "$branch" in
    main|master) say "" "on the default branch — branch before activating a Mission" ;;
  esac
  ok=$((ok+1))
else
  say "git" "not a git repository — commit and tree bindings unavailable"
  warn=$((warn+1))
fi

# 4. Optional TS helpers.
if command -v node >/dev/null 2>&1; then
  say "node" "$(node --version 2>/dev/null) — TS helpers available"
else
  say "node" "absent — shell fallbacks only"
fi

printf '\nMODE\n'
if [ "$cli" -eq 1 ]; then
  printf '  full — read, draft, and governed execution\n'
else
  printf '  reduced — read, explain, and draft only\n\n'
  printf '  Unavailable without the CLI: mission start, objective promote/finish,\n'
  printf '  run start, review record, handoff record, mission complete, contract\n'
  printf '  amend. These produce fingerprints and atomic writes; nothing here\n'
  printf '  emulates them.\n\n'
  printf '  Install:\n'
  printf '    install/install.sh install --prefix <abs-prefix> --source <abs-release> --runtime claude\n'
fi

printf '\n%d ok, %d needing attention\n' "$ok" "$warn"
