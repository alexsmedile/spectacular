---
type: concept-piece
id: PZL-097
status: captured
domain: delivery-lifecycle
sources: [source-007]
source_authority: synthesized-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-073, PZL-074]
overlaps_with: [PZL-018, PZL-066, PZL-098]
conflicts_with: []
tags: [local, pr, merged, deployed, observed, completion]
updated: 2026-08-07
---

# Explicit completion boundary

## Core message

Each Mission declares whether completion means local proof, draft PR, merge,
deployment, or deployed-and-observed behavior.

## Value

Prevents local green checks from being mistaken for the promised delivery outcome.

## Assumptions

- Boundaries apply across supported project types or have neutral equivalents.
- Spectacular can observe provider states without owning their mutations.

## Evidence and collisions

This reconciles Source 002's operations concern and Source 006's draft-PR MVP through
configuration, but risks merging request and release lifecycles currently kept separate.

## Trade-offs and recommendation

Honest completion versus provider coupling and long-lived Missions. Preserve the field
as a delivery contract; keep provider execution and release records separable.

## Decision

Pending.
