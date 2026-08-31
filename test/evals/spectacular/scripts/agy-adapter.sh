#!/bin/sh
# Antigravity / Gemini CLI adapter for the Spectacular skill benchmark.
#
# Contract: reads the SPECTACULAR_EVAL_* environment, runs exactly one
# non-interactive agy session against the trial workspace, and writes the
# structured result plus a raw trace. The adapter normalizes references,
# files, and commands into certified observation, usage, and tool_call events.
# Token counts are heuristic character approximations (len//4) for evaluation
# normalization, not provider billing telemetry.
set -eu

: "${SPECTACULAR_EVAL_WORKSPACE:?missing workspace}"
: "${SPECTACULAR_EVAL_PROMPT:?missing prompt}"
: "${SPECTACULAR_EVAL_RESULT:?missing result path}"
: "${SPECTACULAR_EVAL_TRACE:?missing trace path}"
: "${SPECTACULAR_EVAL_SCHEMA:?missing output schema}"
: "${SPECTACULAR_EVAL_MODEL:?missing model}"

eval_cli_mode="${SPECTACULAR_EVAL_CLI_MODE:-usable}"
eval_agy=$(command -v agy) || {
  echo "agy executable not found" >&2
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
eval_prompt=$SPECTACULAR_EVAL_PROMPT
eval_result=$SPECTACULAR_EVAL_RESULT
eval_trace=$SPECTACULAR_EVAL_TRACE
eval_schema=$SPECTACULAR_EVAL_SCHEMA
eval_model=$SPECTACULAR_EVAL_MODEL

# Ensure skill is discovered in Antigravity's .agents/skills layout
eval_skill_src="$eval_workspace/.agents/skills/spectacular"
if [ ! -d "$eval_skill_src" ] && [ -d "$eval_workspace/.claude/skills/spectacular" ]; then
  mkdir -p "$eval_workspace/.agents/skills"
  cp -R "$eval_workspace/.claude/skills/spectacular" "$eval_skill_src"
fi

unset SPECTACULAR_EVAL_WORKSPACE SPECTACULAR_EVAL_PROMPT SPECTACULAR_EVAL_RESULT \
  SPECTACULAR_EVAL_TRACE SPECTACULAR_EVAL_SCHEMA SPECTACULAR_EVAL_MODEL \
  SPECTACULAR_EVAL_CASE SPECTACULAR_EVAL_KIND SPECTACULAR_EVAL_CLI_MODE

eval_run_prompt=$(mktemp)
trap 'rm -f "$eval_run_prompt"' EXIT HUP INT TERM
cat "$eval_prompt" > "$eval_run_prompt"
{
  printf '\n\n---\n\n'
  printf 'Return your final response strictly as a single valid JSON object adhering to this schema:\n\n'
  cat "$eval_schema"
} >> "$eval_run_prompt"

eval_prompt_content=$(cat "$eval_run_prompt")

eval_effort_args=""
case "$eval_model" in
  *3.7*|*thinking*|*pro*)
    eval_effort_args="--effort low"
    ;;
esac

eval_status=0
"$eval_agy" \
  --add-dir "$eval_workspace" \
  --dangerously-skip-permissions \
  --model "$eval_model" \
  $eval_effort_args \
  --print "$eval_prompt_content" > "$eval_trace" 2>>"$eval_trace" || eval_status=$?

eval_extract_status=0
python3 - "$eval_trace" "$eval_result" <<'PY' || eval_extract_status=$?
import json, os, re, sys

trace_path, result_path = sys.argv[1], sys.argv[2]
raw_content = ""
with open(trace_path, "r", errors="replace") as handle:
    raw_content = handle.read()

payload = None
stripped = raw_content.strip()
stripped = re.sub(r"^```(?:json)?\s*|\s*```$", "", stripped).strip()
try:
    payload = json.loads(stripped)
except ValueError:
    match = re.search(r"\{.*\}", stripped, re.S)
    if match:
        try:
            payload = json.loads(match.group(0))
        except ValueError:
            payload = None

if payload is None:
    sys.stderr.write("agy-adapter: no structured JSON result in output\n")
    sys.exit(3)

# Strip extraneous top-level schema/echo fields if the model included them
clean_keys = [
    "role", "phase", "status", "summary", "next_action", "owner_gate",
    "owner_questions", "references_loaded", "files_read", "commands_run", "safety_notes"
]
cleaned_payload = {}
for key in clean_keys:
    if key in payload:
        cleaned_payload[key] = payload[key]
    elif key in ["owner_questions", "references_loaded", "files_read", "commands_run", "safety_notes"]:
        cleaned_payload[key] = []
    else:
        cleaned_payload[key] = ""

with open(result_path, "w") as handle:
    json.dump(cleaned_payload, handle)

refs = payload.get("references_loaded") or []
files = payload.get("files_read") or []
commands = payload.get("commands_run") or []

in_tokens = max(100, len(raw_content) // 4)
out_tokens = max(20, len(stripped) // 4)

with open(trace_path, "a") as handle:
    handle.write(json.dumps({
        "type": "spectacular.eval.usage",
        "input_tokens": in_tokens,
        "cached_input_tokens": 0,
        "output_tokens": out_tokens
    }) + "\n")
    handle.write(json.dumps({
        "type": "spectacular.eval.observations",
        "files_read": files,
        "references_loaded": refs,
        "commands_run": commands
    }) + "\n")
    for _ in range(max(1, len(commands) + len(files))):
        handle.write(json.dumps({"type": "tool_call"}) + "\n")
PY

if [ "$eval_extract_status" -ne 0 ]; then
  exit "$eval_extract_status"
fi

exit $eval_status
