#!/usr/bin/env bash
# Content-free progressive-loading baseline. It measures existing read views only.
# Bash 3.2 compatible; each run creates and removes its own fixture workspace.

set -eu

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
CLI="$REPO_ROOT/cli/spectacular"
FIXTURE="$REPO_ROOT/tests/benchmarks/fixtures/retrieval-baseline.tsv"
WS="$(mktemp -d "${TMPDIR:-/tmp}/spectacular-retrieval-baseline.XXXXXX")"
OUT="$WS/output"

cleanup() { rm -rf "$WS"; }
trap cleanup EXIT HUP INT TERM

seed_request() {
  local slug="$1" status="$2" milestone="$3"
  mkdir -p "$WS/.spectacular/requests/$slug"
  cat > "$WS/.spectacular/requests/$slug/PLAN.md" <<EOF
---
status: $status
priority: medium
owner: benchmark
updated: 2026-08-04
build: b0
summary: "Fixture request"
related:
  - PRD.md
---
# Plan — $slug

## Goal

Fixture goal.

## Constraints

- Fixture-only content.

## Milestones

- $milestone — Fixture milestone

## Tasks

See TASKS.md.

## Dependencies

- None.

## Validation

- $milestone — Fixture validation.

## Deliverables

- Fixture deliverable.
EOF
  cat > "$WS/.spectacular/requests/$slug/TASKS.md" <<EOF
---
status: $status
related:
  - PLAN.md
---
# Tasks — $slug

### $milestone — Fixture milestone
- [ ] Complete fixture task.
EOF
}

(cd "$WS" && "$CLI" init --minimal --name benchmark --skill-scope none >/dev/null)
seed_request active-work active M1
seed_request review-work review M1
mkdir -p "$WS/.spectacular/specs"
cat > "$WS/.spectacular/specs/SPC-001-fixture.md" <<'EOF'
---
id: SPC-001
type: specification
status: approved
target_version: "tbd"
updated: 2026-08-04
summary: "Fixture specification"
related: []
---
# SPC-001 — Fixture specification

## Intent

Fixture intent.

## Requirements

- Fixture requirement.
EOF
mkdir -p "$WS/.spectacular/fixes"
cat > "$WS/.spectacular/fixes/F001-fixture.md" <<'EOF'
---
id: F001
verified: 2026-08-04
signature: "fixture symptom"
---
# Fixture fix
EOF
cat > "$WS/.spectacular/requests/review-work/VERIFY.md" <<'EOF'
# Verification — review-work

- [ ] {assert} Fixture property holds.
EOF

run_flow() {
  local flow="$1"
  case "$flow" in
    unknown-workspace)
      (cd "$WS" && "$CLI" summary --json; "$CLI" requests --active --json; "$CLI" request active-work) >> "$OUT" ;;
    summary-candidate)
      (cd "$WS" && "$CLI" summary --json) >> "$OUT" ;;
    status-current)
      (cd "$WS" && "$CLI" status --json) >> "$OUT" ;;
    status-brief)
      (cd "$WS" && "$CLI" status --brief --json) >> "$OUT" ;;
    named-spec)
      (cd "$WS" && "$CLI" spec SPC-001 --json) >> "$OUT" ;;
    request-document-review)
      (cd "$WS" && "$CLI" request active-work --full) >> "$OUT" ;;
    named-request)
      (cd "$WS" && "$CLI" request active-work) >> "$OUT" ;;
    active-resume)
      (cd "$WS" && "$CLI" request active-work --brief) >> "$OUT" ;;
    simple-bug-fix)
      (cd "$WS" && "$CLI" fix list) >> "$OUT" ;;
    document-review)
      (cd "$WS" && "$CLI" show prd --section Goal) >> "$OUT"
      target_body=$(< "$WS/.spectacular/PRD.md") ;;
    verification)
      (cd "$WS" && "$CLI" request review-work) >> "$OUT"
      target_body=$(< "$WS/.spectacular/requests/review-work/VERIFY.md") ;;
    *) echo "unknown flow: $flow" >&2; exit 1 ;;
  esac
}

printf '| Flow | References | CLI calls | Output bytes | Full-body reads | Repeated reads | Next action exposed |\n'
printf '|---|---:|---:|---:|---:|---:|---|\n'
while IFS='|' read -r flow references calls full_reads repeated_reads next_action; do
  [[ -z "$flow" || "$flow" = \#* ]] && continue
  : > "$OUT"
  run_flow "$flow"
  bytes=$(wc -c < "$OUT" | tr -d ' ')
  if [[ "$references" == "-" ]]; then
    reference_count=0
  else
    reference_count=$(printf '%s' "$references" | awk -F, '{ print NF }')
  fi
  call_count=$(printf '%s' "$calls" | awk -F, '{ print NF }')
  printf '| %s | %s | %s | %s | %s | %s | %s |\n' \
    "$flow" "$reference_count" "$call_count" "$bytes" "$full_reads" "$repeated_reads" "$next_action"
done < "$FIXTURE"
