#!/usr/bin/env bash
# catalog.sh — what does each skill reference doc do?
#
# Read-only. Scans skills/spectacular/references/*.md and prints a catalog
# derived from each file's own frontmatter (description / when_to_use),
# falling back to a rules-file `summary:` or the first heading.
#
# Mirrors the SKILL.md frontmatter convention: `description` is the
# what-it-does field; `when_to_use` is the load trigger. SKILL.md's
# "Reference loading" table remains the authoritative routing source —
# this script is the self-describing, drift-free companion view.
#
# Usage:
#   catalog.sh                 # table: file · description
#   catalog.sh --when          # table: file · description · when_to_use
#   catalog.sh --missing       # only files lacking description/summary frontmatter
#   catalog.sh --json          # machine-readable
#
# Resolves its own location via symlink so it works from the installed binary.

set -euo pipefail

SOURCE="${BASH_SOURCE[0]}"
while [[ -L "$SOURCE" ]]; do
  DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
  SOURCE="$(readlink "$SOURCE")"; [[ "$SOURCE" != /* ]] && SOURCE="$DIR/$SOURCE"
done
SCRIPT_DIR="$(cd -P "$(dirname "$SOURCE")" && pwd)"
REF_DIR="$(cd "$SCRIPT_DIR/../references" && pwd)"

MODE="table"
case "${1:-}" in
  --when)    MODE="when" ;;
  --missing) MODE="missing" ;;
  --json)    MODE="json" ;;
  ""|--help|-h)
    [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]] && { sed -n '2,20p' "$SOURCE"; exit 0; } ;;
esac

# Extract a frontmatter scalar (description / summary / when_to_use).
# Handles `key: value`, `key: "value"`, and `key: |` block scalars (first line).
fm_field() {
  local file="$1" key="$2"
  awk -v k="$key" '
    NR==1 && $0!="---" { exit }                 # no frontmatter
    NR==1 { infm=1; next }
    infm && $0=="---" { exit }
    infm {
      if ($0 ~ "^"k":[[:space:]]*\\|[[:space:]]*$") { block=1; next }
      if (block) { sub(/^[[:space:]]+/,""); print; exit }
      if ($0 ~ "^"k":[[:space:]]*") {
        sub("^"k":[[:space:]]*",""); gsub(/^"|"$/,""); print; exit
      }
    }
  ' "$file"
}

first_heading() {
  grep -m1 -E '^#[[:space:]]' "$1" 2>/dev/null | sed -E 's/^#+[[:space:]]*//'
}

# description, falling back to summary (rules files), then first heading.
describe() {
  local file="$1" d
  d="$(fm_field "$file" description)"
  [[ -z "$d" ]] && d="$(fm_field "$file" summary)"
  [[ -z "$d" ]] && d="$(first_heading "$file")"
  echo "$d"
}

shopt -s nullglob
files=("$REF_DIR"/*.md)

if [[ "$MODE" == "json" ]]; then
  printf '[\n'
  first=1
  for f in "${files[@]}"; do
    name="$(basename "$f")"
    desc="$(describe "$f")"; when="$(fm_field "$f" when_to_use)"
    [[ $first -eq 0 ]] && printf ',\n'; first=0
    printf '  {"file":"%s","description":"%s","when_to_use":"%s"}' \
      "$name" "${desc//\"/\\\"}" "${when//\"/\\\"}"
  done
  printf '\n]\n'
  exit 0
fi

if [[ "$MODE" == "missing" ]]; then
  echo "Reference docs with no description/summary frontmatter:"
  n=0
  for f in "${files[@]}"; do
    if [[ -z "$(fm_field "$f" description)" && -z "$(fm_field "$f" summary)" ]]; then
      printf "  %s\n" "$(basename "$f")"; n=$((n+1))
    fi
  done
  echo "($n of ${#files[@]} missing)"
  exit 0
fi

# table / when
printf "Spectacular skill — reference doc catalog (%d docs)\n\n" "${#files[@]}"
for f in "${files[@]}"; do
  name="$(basename "$f")"
  desc="$(describe "$f")"
  if [[ "$MODE" == "when" ]]; then
    when="$(fm_field "$f" when_to_use)"
    printf "• %-28s %s\n" "$name" "$desc"
    [[ -n "$when" ]] && printf "  %-28s ↳ when: %s\n" "" "$when"
  else
    printf "• %-28s %s\n" "$name" "$desc"
  fi
done
