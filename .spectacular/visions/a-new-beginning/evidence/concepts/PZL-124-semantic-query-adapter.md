---
type: concept-piece
id: PZL-124
status: captured
domain: semantic-access
sources: [source-010]
source_authority: unsourced-expanded-synthesis
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-107, PZL-122]
overlaps_with: [PZL-106]
conflicts_with: [PZL-077]
tags: [semantic-layer, query, adapter, storage]
updated: 2026-08-07
---

# Semantic query adapter

## Core message

Expose stable business concepts and governed operations while an adapter translates them to
physical SQL, APIs, files, or graph traversal.

## Value

Keeps agent instructions stable when storage implementation changes and centralizes access rules.

## Assumptions

- Business meaning is stable and authorization, freshness, cost, and lineage are observable.
- Unsupported or ambiguous translations fail explicitly.

## Evidence and collisions

This is a concrete access pattern beyond PZL-107's semantic model. Spectacular's current
Markdown/Git core has no demonstrated database-query problem requiring such infrastructure.

## Trade-offs and recommendation

The adapter reduces repeated wiring but can become a leaky query language or stale shadow model.
Keep it as a companion/server pattern until a Spectacular use case earns it.

## Decision

Pending.
