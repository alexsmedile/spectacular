#!/bin/sh
set -eu

: "${SPECTACULAR_EVAL_WORKSPACE:?missing workspace}"
: "${SPECTACULAR_EVAL_PROMPT:?missing prompt}"
: "${SPECTACULAR_EVAL_RESULT:?missing result path}"
: "${SPECTACULAR_EVAL_TRACE:?missing trace path}"
: "${SPECTACULAR_EVAL_SCHEMA:?missing output schema}"
: "${SPECTACULAR_EVAL_MODEL:?missing model}"

eval_cli_mode="${SPECTACULAR_EVAL_CLI_MODE:-usable}"
case "$eval_cli_mode" in
  usable) ;;
  absent)
    PATH="/Users/alex/.bun/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
    export PATH
    ;;
  incompatible)
    eval_bin_dir=$(mktemp -d)
    trap 'rm -rf "$eval_bin_dir"' EXIT HUP INT TERM
    printf '%s\n' '#!/bin/sh' 'echo "spectacular 1.0.0-incompatible"' > "$eval_bin_dir/spectacular"
    chmod +x "$eval_bin_dir/spectacular"
    PATH="$eval_bin_dir:/Users/alex/.bun/bin:/opt/homebrew/bin:/usr/bin:/bin:/usr/sbin:/sbin"
    export PATH
    ;;
  *)
    echo "unsupported SPECTACULAR_EVAL_CLI_MODE: $eval_cli_mode" >&2
    exit 2
    ;;
esac

codex exec \
  --cd "$SPECTACULAR_EVAL_WORKSPACE" \
  --ephemeral \
  --ignore-user-config \
  --skip-git-repo-check \
  --sandbox workspace-write \
  --approve-for-me \
  --model "$SPECTACULAR_EVAL_MODEL" \
  --output-schema "$SPECTACULAR_EVAL_SCHEMA" \
  --output-last-message "$SPECTACULAR_EVAL_RESULT" \
  --json \
  - < "$SPECTACULAR_EVAL_PROMPT" > "$SPECTACULAR_EVAL_TRACE"
