---
status: verified
priority: medium
owner: alex
updated: 2026-08-06
build: b45
docs_impact: none
summary: "Implement SPC-005: Bounded AFK coordination convention with auditable delegation and protected Git boundaries"
related:
  - PRD.md
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
source_spec: SPC-005
source_type: spec
source_ref: SPC-005
source_spec_version: 1.0
source_spec_digest: "sha256:c679650931b4aa6ef557c29b34ba8cee339249e2eab8e5da74fe1f105eecd5d5"
scaffolded_against: 162df3b7430c8720ad695b46d5afdc0df057375c
activated_at: 2026-08-06
activated_by: alex
activated_against: 162df3b7430c8720ad695b46d5afdc0df057375c
docs_impact_reason: Changes are internal CLI and skill references; no Pageworks-owned public documentation was edited
---

# Plan — afk-orchestration-authority

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

Implement SPC-005 M1: a Bash 3.2-compatible, no-Git-mutation AFK run record
format with read-only inspection and doctor validation, so a human can review
declared authority and its event chronology before any future mutation-capable
AFK work is proposed.
## Constraints

- Stay within the approved SPC-005 requirements; unresolved scope changes return to discovery.
- Do not add commands that stage, commit, amend, push, open a PR, merge, reset,
  stash, delete branches, or mutate remote state.
- Preserve Issue #5's separate, read-only `session end` review.
- Keep legacy AFK run files readable and retain the existing AFK lifecycle.
- Use Bash 3.2 only and add no dependencies.
## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`afk run start` currently records only a goal, optional request, CSV allowed
actions, CSV HITL gates, and dates. `afk run status` prints those fields.
`doctor afk` validates status plus a small required-field set. Existing AFK
Git commands are independently protected by project opt-in and `--apply --yes`;
they do not have a per-run event ledger or integration-branch authority model.

### What changes

The AFK run record gains an explicit new-format authorization block and a
CLI-appended event log. New read-only inspection exposes the authority, event
chronology, gates, and missing evidence. `doctor afk` validates the new format
when present while treating missing new fields as legacy narrower authority.

### What stays the same

No AFK command gains Git-mutating authority. Existing branch creation, cleanup,
and draft-PR behavior stays unchanged; protected/default/release/deployment
merges and remote deletion remain outside AFK authority. Session-end commit
review remains independent and read-only.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose a versioned, optional run-record schema over rewriting legacy runs —
  because compatibility must fail closed without silently expanding authority.
- Chose structured CLI event append over direct free-form mutation — because it
  supplies a checkable chronology while keeping Markdown human-readable.

## Milestones

- M1 — Record, inspect, and validate AFK authority without Git mutation.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Approved source specification: SPC-005.
- Current AFK command/test conventions in `cli/spectacular` and
  `tests/cli/afk-git-hygiene.test.sh`.
## Validation

- M1 — focused AFK tests cover new-format start/event/status validation,
  malformed and legacy records, event ordering/evidence, and the no-Git-mutation
  guarantee; then run Bash syntax, pre-commit version guard, and diff check.
## Deliverables

- M1 implementation, focused AFK tests, and verification evidence.
- Documentation-impact assessment for the internal CLI/skill contract.
