#!/bin/sh
# Adapter contract (mirrors codex-adapter.sh): runs one headless opencode agent
# inside $SPECTACULAR_EVAL_WORKSPACE, writes a strict AgentResult JSON object to
# $SPECTACULAR_EVAL_RESULT and a JSONL event trace to $SPECTACULAR_EVAL_TRACE.
# Raw opencode events are kept verbatim; the adapter appends normalized events
# (usage totals, tool markers, spectacular.eval.* semantic events) so
# ParseTraceMetrics observes tokens, tool calls, and observations.
set -eu

: "${SPECTACULAR_EVAL_WORKSPACE:?missing workspace}"
: "${SPECTACULAR_EVAL_PROMPT:?missing prompt}"
: "${SPECTACULAR_EVAL_RESULT:?missing result path}"
: "${SPECTACULAR_EVAL_TRACE:?missing trace path}"
: "${SPECTACULAR_EVAL_SCHEMA:?missing output schema}"
: "${SPECTACULAR_EVAL_MODEL:?missing model}"

eval_cli_mode="${SPECTACULAR_EVAL_CLI_MODE:-usable}"
eval_opencode=$(command -v opencode) || {
  echo "opencode executable not found" >&2
  exit 127
}
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
eval_result=$SPECTACULAR_EVAL_RESULT
eval_trace=$SPECTACULAR_EVAL_TRACE
eval_schema=$SPECTACULAR_EVAL_SCHEMA
eval_model=$SPECTACULAR_EVAL_MODEL
unset SPECTACULAR_EVAL_WORKSPACE SPECTACULAR_EVAL_RESULT \
  SPECTACULAR_EVAL_TRACE SPECTACULAR_EVAL_MODEL \
  SPECTACULAR_EVAL_CASE SPECTACULAR_EVAL_KIND SPECTACULAR_EVAL_CLI_MODE

eval_raw=$(mktemp "${TMPDIR:-/tmp}/spectacular-opencode-trace.XXXXXX")
trap 'rm -f "$eval_raw"' EXIT HUP INT TERM

eval_contract="------------------------------------------------------------
MANDATORY OUTPUT CONTRACT
1. First, use your read tool on the JSON Schema file at exactly this path: $eval_schema
2. Perform the task described below inside the current working directory.
3. Your FINAL assistant message must be EXACTLY ONE JSON object that is valid against that schema, reporting what you actually did. No markdown fences, no prose before or after the object.
Result fields: role, phase, status, summary, next_action, owner_gate, owner_questions, references_loaded, files_read, commands_run, safety_notes.
------------------------------------------------------------
"

{
  printf '%s\n' "$eval_contract"
  cat "$SPECTACULAR_EVAL_PROMPT"
} > "$eval_raw.prompt"

"$eval_opencode" run \
  --dir "$eval_workspace" \
  --auto \
  --model "$eval_model" \
  --format json \
  -- \
  "$(cat "$eval_raw.prompt")" > "$eval_raw"

cp "$eval_raw" "$eval_trace"

# Normalized token totals (opencode emits per-step counters).
jq -rs '
  map(select(.part.type == "step-finish") | .part.tokens)
  | if length > 0 then
      {
        type: "spectacular.eval.usage",
        input_tokens: (map(.input // 0) | add),
        cached_input_tokens: (map(.cache.read // 0) | add),
        output_tokens: (map((.output // 0) + (.reasoning // 0)) | add)
      }
    else empty end
' "$eval_raw" >> "$eval_trace"

# Semantic observation events derived from tool parts.
jq -rs '
  [.[] | select(.part.type == "tool") | .part]
  | map(
      if .tool == "read" and ((.state.input.filePath // "") | length > 0) then
        [{type: "tool_call"},
         {type: "spectacular.eval.file_read", path: .state.input.filePath}]
      elif .tool == "bash" and ((.state.input.command // "") | length > 0) then
        [{type: "tool_call"},
         {type: "spectacular.eval.command", command: .state.input.command}]
      else [{type: "tool_call"}]
      end)
  | flatten[]
' "$eval_raw" >> "$eval_trace"

jq -rs '
  [.[] | select(.part.type == "text")]
  | if length == 0 then "" else last | .part.text end
' "$eval_raw" | python3 -c '
import sys, json
target = sys.argv[1]
text = sys.stdin.read()
if not text.strip():
    sys.exit(4)
start, end = text.find("{"), text.rfind("}")
if start < 0 or end <= start:
    sys.exit(5)
obj = text[start:end + 1]
json.loads(obj)
with open(target, "w") as handle:
    handle.write(obj + "\n")
' "$eval_result"
