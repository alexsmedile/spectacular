#!/usr/bin/env bash
set -euo pipefail

TARGET="${1:-}"
if [ -z "$TARGET" ]; then
  echo "Usage: bash skills/spectacular/scripts/count-tokens.sh <file-path|->" >&2
  exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 1. Try node runner first
if command -v node >/dev/null 2>&1 && [ -f "$SCRIPT_DIR/count-tokens.mjs" ]; then
  exec node "$SCRIPT_DIR/count-tokens.mjs" "$TARGET"
fi

# 2. Pure POSIX / Awk fallback
if [ "$TARGET" = "-" ]; then
  CONTENT="$(cat)"
  LABEL="stdin"
else
  if [ ! -f "$TARGET" ]; then
    echo "Error: file not found: $TARGET" >&2
    exit 1
  fi
  CONTENT="$(cat "$TARGET")"
  LABEL="$TARGET"
fi

LINES=$(printf "%s\n" "$CONTENT" | wc -l | tr -d ' ')
WORDS=$(printf "%s\n" "$CONTENT" | wc -w | tr -d ' ')
CHARS=$(printf "%s\n" "$CONTENT" | wc -m | tr -d ' ')
# Approximation: ~4 chars per token for o200k_base
TOKENS=$(( (CHARS + 3) / 4 ))

echo "Target: $LABEL"
echo "  Lines:      $LINES"
echo "  Words:      $WORDS"
echo "  Characters: $CHARS"
echo "  Tokens:     ~$TOKENS (POSIX estimate)"
