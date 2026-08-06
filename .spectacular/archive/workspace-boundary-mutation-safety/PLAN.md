---
status: archived
priority: medium
owner: alex
updated: 2026-08-06
build: b45
docs_impact: none
summary: "Implement SPC-006: Define additive workspace-boundary and schema-mutation safety hardening under schema 2.0"
related:
  - PRD.md
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
source_spec: SPC-006
source_type: spec
source_ref: SPC-006
source_spec_version: 1.0
source_spec_digest: "sha256:b22f392ba601c251aa77cdeeb0ec05f536fbd6ac30616df7fc9d9d4da3544818"
scaffolded_against: 5f57afa040bbe1bed576479d63d0853ef62ae98a
activated_at: 2026-08-06
activated_by: alex
activated_against: 5f57afa040bbe1bed576479d63d0853ef62ae98a
docs_impact_reason: Internal CLI and skill-reference hardening only; no public docs surface changed.
archived: 2026-08-06
---

# Plan — workspace-boundary-mutation-safety

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

Make every mutating CLI path safe when a workspace declares a schema newer than
the running CLI, while bringing local-workspace guidance and migration preflight
behavior into line with DEC-022—without changing schema 2.0 or migrating data.
## Constraints

- Stay within the approved SPC-006 requirements; unresolved scope changes return to discovery.
- Preserve existing behavior outside this specification.
- Production code and invariant tests become implementation truth after verification.
## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

`config_workspace_schema()` reads the declared marker and `schema_lt()` supports
older/equal comparisons. `doctor` identifies a newer workspace, but
`status --against-latest` labels every non-equal schema as behind and ordinary
mutators have no shared newer-schema write guard. `migrate` can write through
its registry apply loop and has no tracked-local-path preflight. The init
workflow still says local state takes precedence, contrary to DEC-022.

### What changes

Add a centralized schema-relation and pre-write guard, correct status wording,
and add migration's filename-only tracked-local-path refusal. Update the local
state guidance and focused Bash fixtures for schema relations, no-write
refusal, and the local-path preflight.

### What stays the same

`CURRENT_SCHEMA` remains `2.0`; the registry chain, workspace layout,
private-data contents, GitHub configuration, and root-artifact cleanup remain
outside this request. Read-only commands and first-time initialization remain
available; `migrate --list`, `--help`, and `--dry-run` remain read-only
diagnostics on a newer workspace.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose a centralized pre-write gate over per-command checks because every
  current and future mutator must share the same fail-closed boundary.
- Chose to permit read-only migration diagnostics on newer workspaces because
  they expose the upgrade condition without creating a mutation risk.

## Milestones

- M1 — Implement and verify the approved SPC-006 requirements.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Approved source specification: SPC-006.
- Linked decisions and questions remain in the source specification.
## Validation

- M1 — run/assert: every approved requirement has implementation evidence and relevant tests pass.
## Deliverables

- Production implementation and tests for SPC-006.
- Verification evidence and documentation-impact assessment.
