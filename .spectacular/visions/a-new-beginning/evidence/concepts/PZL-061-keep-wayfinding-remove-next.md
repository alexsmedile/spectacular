---
type: concept-piece
id: PZL-061
status: captured
domain: wayfinding
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: supported
disposition: pending
depends_on: [PZL-033, PZL-047]
overlaps_with: [PZL-045, PZL-053]
conflicts_with: [PZL-045]
tags: [wayfinding, next, dependencies, consolidation]
updated: 2026-08-07
---

# Keep dependency wayfinding, remove generic next

## Core message

Retain `wayfind next` for dependency-aware sequencing of discovery and specs, but
remove the ambiguous separate top-level `next` projection.

## Value

Preserves a unique product question while eliminating the name collision that
makes it appear redundant with general status.

## Assumptions

- Corrected Wayfinder results improve action quality often enough to earn the surface.
- Top-level next has no distinct protected outcome.

## Evidence and collisions

Live comparison shows wayfind and status select different work, confirming unique
semantics. The kind/type defect currently undermines ranking reliability. Source
004 instead proposes deleting Wayfinder entirely.

## Trade-offs and recommendation

Unique dependency guidance versus maintaining a scheduler-like subsystem. Fix and
measure the behavior before choosing between retention, reduction, or removal.

## Decision

Pending.
