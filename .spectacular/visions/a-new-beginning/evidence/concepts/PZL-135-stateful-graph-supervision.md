---
type: concept-piece
id: PZL-135
status: captured
domain: information-architecture
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-132, PZL-134, PZL-139]
overlaps_with: [PZL-120]
conflicts_with: []
tags: [mission, graph, state, supervision, visualization]
updated: 2026-08-07
---

# Stateful graph supervision projection

## Core message

Project live node state, blockers, evidence, and next eligible work onto a navigable graph once
the underlying Mission and node contracts are settled.

## Value

Offers a compressed supervision surface for long or branching programs.

## Assumptions

The displayed graph is derived from authoritative state transitions and stable node identities.

## Evidence and collisions

Issue #24 has a viable presentation direction through #18 but its data model depends on disputed
issues #20 and #31. A diagram alone must not imply executable orchestration.

## Trade-offs and recommendation

Separate renderer work from model binding. Validate the view over static fixtures first, then bind
only after the portfolio and node decisions are explicit.

## Decision

Pending.
