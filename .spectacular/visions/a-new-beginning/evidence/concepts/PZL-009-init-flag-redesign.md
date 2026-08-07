---
type: concept-piece
id: PZL-009
status: captured
domain: cli-compatibility
sources: [source-001, source-003, source-004, source-005]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-006, PZL-007, PZL-008]
overlaps_with: [PZL-006, PZL-008, PZL-058, PZL-060]
conflicts_with: []
tags: [cli, init, flags, compatibility, collections]
updated: 2026-08-07
---

# Init flag redesign

## Core message

Make minimal initialization the default, add an explicit full-scaffold mode, and
remove flags whose behavior no longer earns a separate surface.

## Value

Could make the common path obvious and reduce confusing combinations.

## Assumptions

- The new default scaffold has already been accepted.
- Existing automation can migrate without disproportionate cost.
- Reserved engineering collections provide insufficient value to retain.

## Evidence and collisions

`--minimal` is used in tests and existing documented flows. Reserved collection
flags encode identity/path intent even where workflows are not implemented.
Bundling default changes, `--full`, `--minimal` removal, and collection-flag
removal hides several independent compatibility decisions.
Source 005 adds a broader v2 grammar migration and a compatibility release,
making init changes one command family inside a project-wide migration contract.

## Trade-offs and recommendation

Cleaner common path versus a breaking CLI migration and lost explicit capability
intent. Split into separate decision packets after the scaffold and kit contracts
are chosen. Disputed as a bundle.

## Decision

Pending.
