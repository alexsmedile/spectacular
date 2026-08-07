---
type: concept-piece
id: PZL-076
status: captured
domain: lifecycle-model
sources: [source-006, source-007]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-066, PZL-067, PZL-073, PZL-074]
overlaps_with: [PZL-018, PZL-039, PZL-097, PZL-099, PZL-100]
conflicts_with: []
tags: [mission, lifecycle, blocked, aborted, consolidation]
updated: 2026-08-07
---

# One compact mission lifecycle

## Core message

Use `draft → ready → running → proving → closed`, with blocked and aborted exits,
and absorb research, bugs, fixes, tasks, verification, and sessions into the mission.

## Value

Offers one visible state machine instead of many partially overlapping lifecycles.

## Assumptions

- Absorbed record types do not require independent identity, provenance, or reuse.
- Blocked and aborted semantics are orthogonal or fully specified.

## Evidence and collisions

Current types differ in authority and evidence; Source 004's discovery collapse is
already disputed. The proposed states also replace the existing five-state request lifecycle.
Source 007 revises Mission states and adds a separate run lifecycle, complicating the
claim that the surrounding lifecycle becomes singular.

## Trade-offs and recommendation

Excellent legibility versus semantic loss and a breaking migration. Compare
behaviors and query needs before collapsing names or storage.

## Decision

Pending.
