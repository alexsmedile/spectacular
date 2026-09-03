#!/bin/sh
set -eu

test -f "$SPECTACULAR_EVAL_SCHEMA"

printf '%s\n' '{"event":"fake","message":"orient.md loaded"}' > "$SPECTACULAR_EVAL_TRACE"
printf '%s\n' '{
  "role": "Orchestrator",
  "phase": "orient",
  "status": "done",
  "summary": "M1 O1 R1 recovered",
  "next_action": "inspect the current Objective",
  "owner_gate": "",
  "owner_questions": [],
  "references_loaded": ["references/orient.md"],
  "files_read": [".spectacular/PROJECT.md"],
  "commands_run": [],
  "safety_notes": []
}' > "$SPECTACULAR_EVAL_RESULT"
