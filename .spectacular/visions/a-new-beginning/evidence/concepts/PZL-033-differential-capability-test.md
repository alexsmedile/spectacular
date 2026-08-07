---
type: concept-piece
id: PZL-033
status: captured
domain: capability-evaluation
sources: [source-003, source-004, source-005]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-026]
overlaps_with: [PZL-010, PZL-027, PZL-047, PZL-061]
conflicts_with: []
tags: [comparison, wayfinding, status, evidence]
updated: 2026-08-07
---

# Differential capability test

## Core message

Run allegedly overlapping capabilities against the same live state and compare
their decisions or outputs before deleting either.

## Value

Replaces intuition about duplication with concrete behavioral evidence.

## Assumptions

- The chosen scenario exercises each capability's intended question.
- Different output is evaluated for usefulness, not merely uniqueness.

## Evidence and collisions

The proposed test was run: `wayfind next` selected SPC-007 while `status --brief`
selected stance-layer. Wayfinding sequences discovery/spec dependencies; status
prioritizes request-fleet work. They are not currently equivalent. Source 005
independently argues for keeping the dependency-specific view, but a focused
sequencer test also reproduced a ranking defect, so usefulness should be measured
after correctness is restored.

## Trade-offs and recommendation

Evidence-backed pruning versus retaining unique but low-value behavior. Strong
evaluation method. Next decide whether the distinct question deserves its surface.

## Decision

Pending.
