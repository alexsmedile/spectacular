---
status: verified
priority: medium
owner: alex
updated: 2026-08-04
build: b43
docs_impact: required
summary: "Implement SPC-004: Define destination-based idea routing and reduce local IDEA storage"
related:
  - PRD.md
# Optional traffic evidence; omit until it is confirmed and complete.
# conflicts-with: [<request-slug>]
# traffic-boundaries: [<named, complete launch boundary>]
# release-constraints: [<shared release or migration constraint>]
source_spec: SPC-004
source_type: spec
source_ref: SPC-004
source_spec_version: 1.0
source_spec_digest: "sha256:9847fa947e9a665f0b7964eaff10d0fc053a1d691044bd4e98656758d28f6994"
scaffolded_against: a5b4d5fbf3073c8e4f87de24255d5959d1e3f0ea
activated_at: 2026-08-04
activated_by: alex
activated_against: a5b4d5fbf3073c8e4f87de24255d5959d1e3f0ea
docs_impact_evidence: Updated internal idea-routing references and architecture; public docs intentionally left unchanged
---

# Plan — idea-destination-routing

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

Make the place where an idea will be developed and shared its authoritative
home, while preserving local IDEA capture as a small compatibility and
private/offline surface.
## Constraints

- Stay within the approved SPC-004 requirements; unresolved scope changes return to discovery.
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

`idea new` always creates a committed `.spectacular/ideas/IDEA-NNN-*.md`.
`idea list` reads only that collection. `idea promote` has one hard-coded
destination: it scaffolds a request, marks the source `promoted`, and moves it
to `archive/ideas/`. The skill references already tell users not to mirror
Issues or Discussions, but the command grammar cannot express that route.

### What changes

The idea command gains a destination model. Request promotion remains the
compatible default; roadmap transfer is a local, explicit Icebox handoff;
shared transfer requires an existing stable reference and has no network path.
The skill and routing references make the selected destination authoritative
and explain that local IDEA storage is no longer the default shared backlog.

### What stays the same

Existing IDEA IDs, files, archives, `undo`, and Wayfinding remain compatible.
No GitHub API behavior, shared-destination adapter, automatic mirroring,
remote mutation, request auto-creation, or Vision auto-creation is added.

## Decisions

<!--
  Design calls made inside this request. Format: chose X over Y — because Z.
  Rejected alternatives stay listed; deleting them re-litigates them later.
  Project-wide calls go to DECISIONS.md via `spectacular decide` instead
  (see decisions-rules.md routing table). Empty is fine — no decisions yet.
-->

- Chose `request|roadmap|shared` over `issue|shared|roadmap|request` because
  Issue is a subtype of a shared destination, not a second local authority.
- Chose an existing-reference handoff for `shared` over remote filing because
  Spectacular cannot safely assume a remote provider, credentials, or approval.

## Milestones

- M1 — Add destination-aware local IDEA transitions and regression tests.
- M2 — Align skill/routing guidance and deprecation language with the new model.
- M3 — Verify compatibility, no-network behavior, and documentation boundaries.
## Tasks

<!-- Pointer. The executable checklist lives in TASKS.md, grouped by milestone. -->

See `TASKS.md`.

## Dependencies

- Approved source specification: SPC-004.
- Linked decisions and questions remain in the source specification.
## Validation

- M1 — run: focused idea/undo/Wayfinding CLI tests pass.
- M2 — assert: references agree that shared surfaces are authoritative and that
  roadmap transfer is an explicit Icebox handoff.
- M3 — run: bash syntax, version guard, and full relevant test suite pass.
## Deliverables

- Destination-aware IDEA CLI behavior and focused CLI regression coverage.
- Updated idea-routing skill/reference contract.
- Verification evidence and documentation-impact assessment.
