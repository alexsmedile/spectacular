---
status: active
priority: medium
owner: alex
updated: 2026-08-05
build: b44
docs_impact: pending
summary: "Audit soft-DB read views and specify compact entity projections, escalation paths, and omission-regression tests"
related:
  - PRD.md
source_type: issue
source_ref: "alexsmedile/spectacular#12"
sensitivity: normal
activated_at: 2026-08-05
activated_by: alex
activated_against: 0de02e1491b82da3fffc0d2346aba2e363fe3599
---

# Plan — soft-db-projection-inventory

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

Define the compact, entity-specific read contract that lets an agent select and
inspect a high-use soft-DB record without batch-reading Markdown, while
preserving the record body as evidence. This serves the PRD success criterion
for a correct cold briefing from navigation signals and is a **design/test
inventory only**: no query engine, parser, cache, or production CLI change is
authorized here.
## Constraints

- Preserve the accepted outcome and named boundaries of Issue alexsmedile/spectacular#12.
- Build on the merged Issue #11 retrieval baseline: `spectacular status --brief
  --json` is the workspace-orientation command. Do not create a competing
  workspace briefing or reopen that default in this request.
- Preserve frontmatter as the queryable signal layer and Markdown bodies as
  named-record evidence. Direct named-file access remains valid after
  selection and for repair/debugging.
- Do not add a generic query language, arbitrary filters, caches, broad YAML
  parsing, or a generic heading/"first sections" slicer.
- Keep Bash 3.2 compatibility and deterministic ordering/output as required
  implementation constraints.
- Link remote evidence; do not copy Issue bodies or comments into the request.
## Understanding

<!--
  OPTIONAL authoring slot, but REQUIRED before `planned → active` by the
  `understand-before-change` policy (@Implementation). Fill it here for a
  typical request; escalate to a dedicated requests/<slug>/UNDERSTANDING.md
  (same three subheads) for large ones — the policy is satisfied by EITHER.
  Not one of the 7 required authoring slots; it gates implementation, not planning.
-->

### How it works now

Issue #11 now supplies the bounded workspace orientation: `status --brief
--json` returns blockers, health, ranked next action, and fleet, then routes to
the existing request overview/brief. Collection reads remain uneven. Requests
have list, overview, active implementation brief, and full-bundle views.
Decisions and memories have lists plus a generic frontmatter-and-heading-outline
skim. Questions, research, spikes, ideas, audits, and fixes have a list only;
callers must open entry Markdown directly for any record-level decision. Some
lists extract a body heading with `awk`, which bypasses the intended signal
layer.

### What changes

The proposed future contract is documented in `PROJECTION-INVENTORY.md`:
one bounded list projection for selection, one entity-specific named detail
projection for its primary decision, and a literal `--full` escalation to the
complete evidence record. It also names the test fixtures and negative tests
needed before any CLI work begins.

### What stays the same

Markdown files remain the source of detailed evidence; frontmatter remains the
indexed subset. Existing mutators, canonical IDs, aliases, direct-file repair,
and request `--brief` semantics remain compatible. The merged #11 orientation
route remains outside this request's collection-projection scope.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose per-entity projections over a shared heading slicer because each
  record type has a different decision-critical shape and headings are
  evidence, not a stable query protocol.
- Chose list → named detail → literal full evidence over mandatory CLI-only
  access because selection benefits from bounded projections while repair and
  debugging require direct named-file access.
- Deferred constrained filters, index rebuild, frontmatter-subset validation,
  and caching until the basic omission contract is proven in real collection
  reads.

## Milestones

- M1 — Align the design with #11's merged `status --brief --json` orientation
  contract and reserve the collection-only implementation boundary.
- M2 — Entity detail contracts and omission-regression fixtures are specified
  for the high-use collections.
- M3 — The smallest compatible CLI slice is implemented only after M1/M2,
  followed by documentation and a measurement-based decision on subsequent
  work.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Source: Issue alexsmedile/spectacular#12.
- Upstream baseline: merged Issue #11 retrieval work and its
  `benchmarks/retrieval-baseline.md` evidence. It defines workspace orientation
  but does not implement collection projections.
- Existing code, tests, documentation, and decisions define the design
  destination; escalate to spec-first if a new workspace-wide query or metadata
  schema contract proves necessary.
## Validation

- M1 — assert the design uses `status --brief --json` as the established
  workspace orientation and introduces no competing briefing command.
- M2 — review `PROJECTION-INVENTORY.md` against every named command and entry
  template; each compact view names a safe decision and exact escalation.
- M3 — future CLI tests prove every decision-critical field is present,
  escaping to `--full` preserves the body verbatim, output order is stable,
  and Bash 3.2 syntax remains valid.
## Deliverables

- `PROJECTION-INVENTORY.md`: current-command audit, proposed compact grammar,
  decision-risk analysis, test inventory, and #11 alignment.
- A sequenced `TASKS.md` that starts implementation only after Issue #11.
- No implementation, commit, push, or pull request in this design request.
