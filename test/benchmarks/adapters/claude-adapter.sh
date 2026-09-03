#!/bin/sh
# Claude Code adapter for the Spectacular skill benchmark.
#
# Contract: reads the SPECTACULAR_EVAL_* environment, runs exactly one
# non-interactive Claude session against the trial workspace, and writes the
# structured result plus a JSONL trace.
#
# Artifact separation only. This is not an OS-enforced read sandbox, so runs
# using this adapter are `inconclusive` for strict isolation, same as Codex.
set -eu

: "${SPECTACULAR_EVAL_WORKSPACE:?missing workspace}"
: "${SPECTACULAR_EVAL_PROMPT:?missing prompt}"
: "${SPECTACULAR_EVAL_RESULT:?missing result path}"
: "${SPECTACULAR_EVAL_TRACE:?missing trace path}"
: "${SPECTACULAR_EVAL_SCHEMA:?missing output schema}"
: "${SPECTACULAR_EVAL_MODEL:?missing model}"

eval_cli_mode="${SPECTACULAR_EVAL_CLI_MODE:-usable}"
eval_claude=$(command -v claude || command -v "$HOME/.local/bin/claude" || echo "")
if [ -z "$eval_claude" ] || [ ! -x "$eval_claude" ]; then
  echo "claude executable not found" >&2
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

# The harness materializes the variant at .agents/skills/ for the Codex adapter.
# Claude Code discovers skills under .claude/skills/, so expose the same single
# materialized package there. A symlink keeps exactly one variant visible.
# A copy, not a symlink: the harness snapshots the workspace tree to detect
# unauthorized writes, and a symlinked directory breaks that walk.
eval_skill_src="$eval_workspace/.agents/skills/spectacular"
eval_skill_dst="$eval_workspace/.claude/skills/spectacular"
eval_skill_staged=0
if [ -d "$eval_skill_src" ] && [ ! -e "$eval_workspace/.claude" ]; then
  mkdir -p "$eval_workspace/.claude/skills"
  cp -R "$eval_skill_src" "$eval_skill_dst"
  eval_skill_staged=1
fi

cleanup_staged_skill() {
  # The harness snapshots the tree before and after this adapter runs, so the
  # staging directory must not survive into the post-run diff.
  if [ "$eval_skill_staged" = "1" ]; then
    rm -rf "$eval_workspace/.claude"
  fi
}

# The model must not see benchmark plumbing.
unset SPECTACULAR_EVAL_WORKSPACE SPECTACULAR_EVAL_PROMPT SPECTACULAR_EVAL_RESULT \
  SPECTACULAR_EVAL_TRACE SPECTACULAR_EVAL_SCHEMA SPECTACULAR_EVAL_MODEL \
  SPECTACULAR_EVAL_CASE SPECTACULAR_EVAL_KIND SPECTACULAR_EVAL_CLI_MODE

# Claude has no --output-schema flag; the shape is appended to the prompt and the
# final assistant message is extracted from the stream afterwards.
eval_run_prompt=$(mktemp)
trap 'rm -f "$eval_run_prompt"' EXIT HUP INT TERM
cat "$eval_prompt" > "$eval_run_prompt"
{
  printf '\n\n---\n\n'
  printf 'Return your final message as a single JSON object and nothing else: no prose before or after, no markdown code fence. It must validate against this JSON Schema:\n\n'
  cat "$eval_schema"
} >> "$eval_run_prompt"

eval_status=0
"$eval_claude" \
  --print \
  --model "$eval_model" \
  --output-format stream-json \
  --verbose \
  --permission-mode bypassPermissions \
  --add-dir "$eval_workspace" \
  < "$eval_run_prompt" > "$eval_trace" 2>>"$eval_trace" || eval_status=$?

# Extract the structured result and derive observations from tool telemetry
# rather than the model's own report of what it read.
eval_extract_status=0
python3 - "$eval_trace" "$eval_result" <<'PY' || eval_extract_status=$?
import json, os, re, sys

trace_path, result_path = sys.argv[1], sys.argv[2]
final_text = ""
files, refs, commands = [], [], []

READ_TOOLS = {"Read", "NotebookRead"}
SEARCH_TOOLS = {"Grep", "Glob"}


def record_tool(name, params):
    if not isinstance(params, dict):
        return
    if name in READ_TOOLS:
        path = params.get("file_path") or params.get("notebook_path")
        if path:
            files.append(path)
    elif name == "Bash":
        command = params.get("command")
        if command:
            commands.append(command)
    elif name in SEARCH_TOOLS:
        pattern = params.get("pattern", "")
        commands.append(f"{name} {pattern}".strip())
    elif name == "Skill":
        skill = params.get("skill")
        if skill:
            refs.append(skill)


with open(trace_path, "r", errors="replace") as handle:
    for line in handle:
        line = line.strip()
        if not line.startswith("{"):
            continue
        try:
            event = json.loads(line)
        except ValueError:
            continue
        message = event.get("message")
        if isinstance(message, dict):
            content = message.get("content")
            if isinstance(content, list):
                for block in content:
                    if not isinstance(block, dict):
                        continue
                    if block.get("type") == "tool_use":
                        record_tool(block.get("name", ""), block.get("input"))
                    elif block.get("type") == "text" and message.get("role") == "assistant":
                        final_text = block.get("text", "")
        if event.get("type") == "result" and isinstance(event.get("result"), str):
            final_text = event["result"]

# Any reference the model loaded shows up as a read under the skill package.
for path in files:
    if "skills/spectacular/references/" in path.replace(os.sep, "/"):
        refs.append(os.path.basename(path))

payload = None
if final_text:
    stripped = final_text.strip()
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
    sys.stderr.write("claude-adapter: no structured JSON result in final message\n")
    sys.exit(3)

with open(result_path, "w") as handle:
    json.dump(payload, handle)

# Adapter-authored observation event, appended to the trace as required by the
# harness. Sourced from tool telemetry above, never from payload self-report.
def unique(items):
    seen, output = set(), []
    for item in items:
        if item not in seen:
            seen.add(item)
            output.append(item)
    return output


with open(trace_path, "a") as handle:
    handle.write(json.dumps({
        "type": "spectacular.eval.observations",
        "files_read": unique(files),
        "references_loaded": unique(refs),
        "commands_run": unique(commands),
    }) + "\n")
PY

cleanup_staged_skill

if [ "$eval_extract_status" -ne 0 ]; then
  exit "$eval_extract_status"
fi

exit $eval_status
