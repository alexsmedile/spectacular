#!/bin/sh
set -eu

: "${SPECTACULAR_EVAL_WORKSPACE:?missing workspace}"
: "${SPECTACULAR_EVAL_PROMPT:?missing prompt}"
: "${SPECTACULAR_EVAL_RESULT:?missing result path}"
: "${SPECTACULAR_EVAL_TRACE:?missing trace path}"
: "${SPECTACULAR_EVAL_SCHEMA:?missing output schema}"
: "${SPECTACULAR_EVAL_MODEL:?missing model}"

eval_cli_mode="${SPECTACULAR_EVAL_CLI_MODE:-usable}"
eval_codex=$(command -v codex || command -v "$HOME/.bun/bin/codex" || echo "")
if [ -z "$eval_codex" ] || [ ! -x "$eval_codex" ]; then
  echo "codex executable not found" >&2
  exit 127
fi
eval_original_path=$PATH

without_spectacular_dirs() {
  eval_filtered_path=
  eval_old_ifs=$IFS
  IFS=:
  for eval_path_dir in $eval_original_path; do
    [ -n "$eval_path_dir" ] || eval_path_dir=.
    if [ -x "$eval_path_dir/spectacular" ]; then
      continue
    fi
    if [ -z "$eval_filtered_path" ]; then
      eval_filtered_path=$eval_path_dir
    else
      eval_filtered_path=$eval_filtered_path:$eval_path_dir
    fi
  done
  IFS=$eval_old_ifs
  printf '%s\n' "$eval_filtered_path"
}

case "$eval_cli_mode" in
  usable) ;;
  absent)
    PATH=$(without_spectacular_dirs)
    export PATH
    ;;
  incompatible)
    eval_bin_dir=$(mktemp -d)
    trap 'rm -rf "$eval_bin_dir"' EXIT HUP INT TERM
    printf '%s\n' '#!/bin/sh' 'echo "spectacular 1.0.0-incompatible"' > "$eval_bin_dir/spectacular"
    chmod +x "$eval_bin_dir/spectacular"
    PATH="$eval_bin_dir:$eval_original_path"
    export PATH
    ;;
  *)
    echo "unsupported SPECTACULAR_EVAL_CLI_MODE: $eval_cli_mode" >&2
    exit 2
    ;;
esac

eval_workspace=$SPECTACULAR_EVAL_WORKSPACE
eval_prompt=$SPECTACULAR_EVAL_PROMPT
eval_result=$SPECTACULAR_EVAL_RESULT
eval_trace=$SPECTACULAR_EVAL_TRACE
eval_schema=$SPECTACULAR_EVAL_SCHEMA
eval_model=$SPECTACULAR_EVAL_MODEL
unset SPECTACULAR_EVAL_WORKSPACE SPECTACULAR_EVAL_PROMPT SPECTACULAR_EVAL_RESULT \
  SPECTACULAR_EVAL_TRACE SPECTACULAR_EVAL_SCHEMA SPECTACULAR_EVAL_MODEL \
  SPECTACULAR_EVAL_CASE SPECTACULAR_EVAL_KIND SPECTACULAR_EVAL_CLI_MODE

"$eval_codex" exec \
  --cd "$eval_workspace" \
  --ephemeral \
  --ignore-user-config \
  --skip-git-repo-check \
  --approve-for-me \
  --model "$eval_model" \
  --output-schema "$eval_schema" \
  --output-last-message "$eval_result" \
  --json \
  - < "$eval_prompt" > "$eval_trace"
