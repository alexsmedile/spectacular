---
type: concept-piece
id: PZL-104
status: captured
domain: orchestration-pattern
sources: [source-008, source-009, source-010, source-012]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-103, PZL-114]
overlaps_with: [PZL-098]
conflicts_with: [PZL-071, PZL-079]
tags: [split, verify, merge, parallelism]
updated: 2026-08-07
---

# Split, work, check, merge

## Core message

For genuinely independent units, split the work, execute in bounded contexts, check
outputs against shared criteria, and merge only accepted results.

## Value

Makes joins and verification explicit while containing individual context growth.

## Assumptions

- Units are independent enough to avoid conflicting edits or duplicated discovery.
- A lead owns decomposition and semantic integration.

## Evidence and collisions

The pattern is plausible but Source 008 supplies no comparison. It conflicts with Source
006's single-agent MVP and its explicit deferral of parallel orchestration.

## Trade-offs and recommendation

Parallelism can reduce latency while increasing coordination, token, and merge costs.
Keep this as an earned pattern triggered by separability and risk, not the default loop.

## Decision

Pending.
