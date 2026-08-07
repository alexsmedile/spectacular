---
type: concept-piece
id: PZL-057
status: captured
domain: dead-code
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-028]
conflicts_with: []
tags: [uuidv7, ids, dead-code, cleanup]
updated: 2026-08-07
---

# Remove dead canonical-ID allocator

## Core message

Delete `_next_canonical_id()` now that UUIDv7 identity allocation has replaced it,
and update any canonical prose that still implies the helper is active.

## Value

Removes misleading compatibility code and closes documentation drift.

## Assumptions

- No dynamic or external caller depends on the private helper.
- UUIDv7 migration is the accepted identity contract.

## Evidence and collisions

Only the function definition appears in executable code; a specification still
mentions the legacy allocation concept.

## Trade-offs and recommendation

Minimal compatibility risk versus preserving an unused fallback. Strong cleanup
candidate after one call-site and release-compatibility audit.

## Decision

Pending.
