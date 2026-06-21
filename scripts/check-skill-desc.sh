#!/usr/bin/env bash
# check-skill-desc.sh — measure a SKILL.md `description` block against Codex's limit.
#
# Codex skips loading any skill whose frontmatter `description` exceeds 1024 chars.
# Claude Code's limit is 1536, so an over-long description loads fine in Claude Code
# while silently breaking Codex. This guard catches it before a release ships.
#
# WHAT IS MEASURED: the `description` field ALONE — not description + when_to_use.
# Proven by the v1.17.2 incident: trimming description 1146 → 986 cleared Codex's
# error even though description + when_to_use stayed at 1253 (> 1024). See
# .spectacular/decisions/D7.md.
#
# Usage:
#   measure_skill_desc <skill.md>        # echoes char count, returns 0
#   skill_desc_severity <count>          # echoes error|warning|pass
#   bash check-skill-desc.sh <skill.md>...   # standalone: prints a line per file,
#                                            # exits 1 if any file is over the limit.
#
# Thresholds: error > 1024, warning > 1000 (last ~24 chars before the cap), else pass.
# 1000 sits above the deliberate ~986 steady state so a healthy file is a clean
# pass; the warning fires only when a description is genuinely about to break Codex.

SKILL_DESC_ERROR=1024
SKILL_DESC_WARN=1000

# Extract the `description` value (single-line or YAML literal block `|`) from a
# SKILL.md and echo its character count. Multibyte-safe (awk length() counts
# characters, not bytes). Echoes 0 if no description / file unreadable.
measure_skill_desc() {
  local file="$1"
  [[ -r "$file" ]] || { echo 0; return; }
  awk '
    /^---$/ { fm = !fm; if (!fm && seen) exit; next }
    # literal-block form:  description: |
    fm && /^description:[[:space:]]*\|[[:space:]]*$/ { in_desc=1; seen=1; next }
    # single-line form:    description: some text   /   description: "some text"
    fm && /^description:[[:space:]]*[^|[:space:]]/ {
      val = $0
      sub(/^description:[[:space:]]*/, "", val)
      gsub(/^"|"$/, "", val)
      total += length(val)
      seen=1
      exit
    }
    # continuation lines of a literal block: indented; stop at next top-level key
    in_desc {
      if ($0 ~ /^[^[:space:]]/ || $0 ~ /^---$/) { exit }
      line = $0
      sub(/^[[:space:]]+/, "", line)   # strip the block indent
      # join with a newline between wrapped lines (matches how the value renders)
      if (started) total += 1          # +1 for the joining newline
      total += length(line)
      started=1
    }
    END { print total + 0 }
  ' "$file"
}

# Map a char count to a severity word.
skill_desc_severity() {
  local n="$1"
  if (( n > SKILL_DESC_ERROR )); then echo "error"
  elif (( n > SKILL_DESC_WARN )); then echo "warning"
  else echo "pass"; fi
}

# Standalone entrypoint — only runs when executed directly, not when sourced.
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  rc=0
  for f in "$@"; do
    n="$(measure_skill_desc "$f")"
    sev="$(skill_desc_severity "$n")"
    case "$sev" in
      error)   printf '  ✗ %-50s description %s chars (over %s — Codex will skip this skill)\n' "$f" "$n" "$SKILL_DESC_ERROR"; rc=1 ;;
      warning) printf '  ! %-50s description %s chars (%s under the %s cap — trim soon)\n' "$f" "$n" "$((SKILL_DESC_ERROR - n))" "$SKILL_DESC_ERROR" ;;
      pass)    printf '  ✓ %-50s description %s chars\n' "$f" "$n" ;;
    esac
  done
  exit "$rc"
fi
