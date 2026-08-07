---
type: concept-piece
id: PZL-128
status: captured
domain: context-retrieval
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-010, PZL-121]
overlaps_with: [PZL-004, PZL-005, PZL-123]
conflicts_with: []
tags: [progressive-disclosure, runtime, measurement, cards]
updated: 2026-08-07
---

# Frequency-weighted runtime cards

## Core message

Measure which instructions are used together, then package frequent runtime paths into small
cards while leaving rare rules behind explicit retrieval boundaries.

## Value

Turns “make the skill leaner” into a measurable retrieval-design problem.

## Assumptions

- Runtime loading can be observed or approximated consistently.
- A cold-start budget and task-success threshold are declared before optimization.

## Evidence and collisions

Issue #11 identifies progressive-loading work but lacks its optimization objective. It converges
with PZL-010 and must share a briefing schema with #12 rather than create a second projection.

## Trade-offs and recommendation

Fewer loaded tokens may increase lookup hops. Establish a benchmark corpus and jointly measure
tokens, first useful action, errors, and retrieval misses before changing boundaries.

## Decision

Pending.
