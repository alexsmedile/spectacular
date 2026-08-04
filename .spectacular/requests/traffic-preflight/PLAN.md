---
status: review
priority: medium
owner: alex
updated: 2026-08-04
build: b42
docs_impact: pending
summary: "Add a conservative local-first request traffic preflight using durable request evidence"
related:
  - PRD.md
source_type: issue
source_ref: "alexsmedile/spectacular#3"
sensitivity: normal
activated_at: 2026-08-04
activated_by: alex
activated_against: 0bbdee856645b23295ee2a60f5b6c38bc4590d9b
---

# Plan — traffic-preflight

<!--
  Canonical 7-slot PLAN template for a single request.
  Lives at: .spectacular/requests/<slug>/PLAN.md

  Rules:
  - PLAN is per-request. PRD is project-wide. Never put a PRD inside requests/.
  - This file's frontmatter `status:` is the single source of lifecycle state for the request.
  - The 7 required sections must appear IN ORDER, unnumbered:
      ## Goal, ## Constraints, ## Milestones, ## Tasks, ## Dependencies, ## Validation, ## Deliverables
    Extra sections (## Understanding, ## Decisions, request-specific) may appear
    BETWEEN them; doctor enforces the required set's presence + order, not a closed list.
  - All 7 required sections must be filled before this PLAN is considered usable.
  - Replace every <placeholder> with concrete content.
-->

## Goal

Classify a proposed request's local launch traffic as `parallel`, `conditional`, `serialized`, or `unknown` using only durable request declarations.
## Constraints

- Reuse SPC-003 and DEC-021 terminology; do not build a scheduler or a second coordination model.
- Never infer file-level safety; absent durable evidence is `unknown`.
- Keep GitHub branch/PR evidence optional and read-only; the local workspace remains sufficient.
- Do not mutate any request during assessment. Confirmed relationships are explicit PLAN frontmatter declarations.
## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`spectacular links` renders declared `depends-on`, `blocks`, and `related` request links, but does not classify launch compatibility.

### What changes

Add a read-only `spectacular traffic preflight` command, durable `conflicts-with`, `traffic-boundaries`, and `release-constraints` frontmatter conventions, and link-graph visibility for conflicts.

### What stays the same

Lifecycle ownership, request scheduling, and remote GitHub operations remain unchanged. Every assessment is recalculated and date-stamped from local PLAN evidence.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose a date-stamped read-only preflight over a stored verdict because current request evidence can change and DEC-021 requires recalculation.

## Milestones

- M1 — Define the conservative durable-evidence contract and command surface.
- M2 — Implement local classification and conflict-link visibility.
- M3 — Verify every traffic state and insufficient-evidence fallback.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Source: Issue alexsmedile/spectacular#3; SPC-003 and DEC-021 define the traffic model.
- Existing cross-request links command and isolated CLI test fixtures provide the local evidence path.
## Validation

- M1 — `bash tests/cli/traffic.test.sh` proves the documented input contract and command help.
- M2 — `bash tests/cli/traffic.test.sh` proves `parallel`, `conditional`, `serialized`, `unknown`, JSON, and conflict graph output.
- M3 — `bash -n cli/spectacular` and `bash tests/cli/{traffic,links,doctor}.test.sh` exit 0.
## Deliverables

- `spectacular traffic preflight <request> [--against <request>] [--json]`.
- Traffic evidence conventions in CLI help, command documentation, and request workflow.
- Focused regression coverage and validation record for Issue alexsmedile/spectacular#3.
