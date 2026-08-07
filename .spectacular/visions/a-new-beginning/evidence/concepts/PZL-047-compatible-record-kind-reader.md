---
type: concept-piece
id: PZL-047
status: captured
domain: record-contract
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-033, PZL-061]
conflicts_with: []
tags: [bug, wayfinding, frontmatter, compatibility, stabilization]
updated: 2026-08-07
---

# Compatible record-kind reader

## Core message

Read canonical `kind:` with a legacy `type:` fallback through one shared helper,
then remove direct field reads from Wayfinder ranking and doctor paths.

## Value

Restores deterministic ranking across mixed-generation records without forcing an
unsafe all-at-once migration.

## Assumptions

- `kind` is the accepted canonical field.
- Legacy `type` records remain possible during compatibility.

## Evidence and collisions

Current Wayfinder readers use `type`, UUIDv7-era records use `kind`, and the
sequencer test reproduced a spike-versus-research ordering failure.

## Trade-offs and recommendation

A compatibility helper adds a short migration tail but centralizes the root
contract. Treat as a pre-refactor stabilization fix after auditing every caller.

## Decision

Pending.
