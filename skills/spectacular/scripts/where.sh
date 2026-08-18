#!/bin/sh
# where.sh — resolve a Spectacular ref to its record path.
#
# Read-only fallback for hosts without the `spectacular` CLI. Pure path
# resolution: it finds the file, it does not parse or validate it.
#
# Usage: where.sh <ref>          e.g. M12, P8, D11, M12/R1
set -eu

[ $# -eq 1 ] || { printf 'usage: where.sh <ref>\n' >&2; exit 2; }
ref=$1

dir=$(pwd)
workspace=""
while [ "$dir" != "/" ]; do
  if [ -d "$dir/.spectacular" ]; then workspace="$dir/.spectacular"; break; fi
  dir=$(dirname "$dir")
done
[ -n "$workspace" ] || { printf 'no .spectacular/ workspace found\n' >&2; exit 1; }

# A nested ref (M12/R1) resolves inside its Mission bundle.
case "$ref" in
  */*)
    mission=${ref%%/*}
    child=${ref#*/}
    for d in "$workspace"/missions/"$mission"-*/ "$workspace"/archive/missions/"$mission"-*/; do
      [ -d "$d" ] || continue
      for sub in runs objectives evidence reviews handoffs decisions gaps assessments checkpoints; do
        for c in "$d$sub"/"$child"-*/ "$d$sub"/"$child"-*.md; do
          if [ -d "$c" ]; then
            inner="$c$(basename "$c").md"
            [ -f "$inner" ] && { printf '%s\n' "$inner"; exit 0; }
          elif [ -f "$c" ]; then
            printf '%s\n' "$c"; exit 0
          fi
        done
      done
      # Checkpoints nest one level deeper, inside a Run bundle.
      for c in "$d"runs/*/checkpoints/"$child"-*.md; do
        [ -f "$c" ] && { printf '%s\n' "$c"; exit 0; }
      done
    done
    ;;
  *)
    # Top-level: Mission, Proposal, Decision, Contract — live then archived.
    for p in \
      "$workspace"/missions/"$ref"-*/"$ref"-*.md \
      "$workspace"/archive/missions/"$ref"-*/"$ref"-*.md \
      "$workspace"/proposals/"$ref"-*.md \
      "$workspace"/archive/proposals/"$ref"-*.md \
      "$workspace"/decisions/"$ref"-*.md \
      "$workspace"/contracts/"$ref"-*.md \
      "$workspace"/contracts/"$ref".md
    do
      [ -f "$p" ] && { printf '%s\n' "$p"; exit 0; }
    done
    ;;
esac

printf 'no record found for ref: %s\n' "$ref" >&2
exit 1
