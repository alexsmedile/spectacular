---
type: concept-piece
id: PZL-088
status: captured
domain: work-unit-model
sources: [source-007]
source_authority: synthesized-proposal
assessment: promising
evidence_status: unverified
disposition: pending
depends_on: [PZL-066, PZL-080]
overlaps_with: [PZL-074, PZL-089, PZL-099]
conflicts_with: []
tags: [mission, graph, delta, baseline, transaction]
updated: 2026-08-07
---

# Mission as versioned graph transaction

## Core message

A Mission names its baseline graph version, intended contract changes, touched
boundaries, completion boundary, autonomy, objectives, risks, assumptions, and evidence.

## Value

Makes impact and reconciliation explicit against a stable accepted baseline.

## Assumptions

- The graph has coherent version identity and atomic-enough reconciliation.
- Concurrent Missions can detect overlapping deltas safely.

## Evidence and collisions

Current contract UUIDs, spec provenance, request baselines, and Git SHAs provide pieces;
no system-graph version or transaction semantics exist.

## Trade-offs and recommendation

Precise deltas versus transaction/version complexity. Prove one single-Mission slice
before designing concurrency or global graph versions.

## Decision

Pending.
