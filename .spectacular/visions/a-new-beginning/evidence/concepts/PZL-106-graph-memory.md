---
type: concept-piece
id: PZL-106
status: captured
domain: context-retrieval
sources: [source-008, source-009]
source_authority: unsourced-protocol-synthesis
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-080, PZL-086]
overlaps_with: [PZL-090]
conflicts_with: [PZL-079]
tags: [graph-memory, retrieval, relationships, rag]
updated: 2026-08-07
---

# Graph memory for structural retrieval

## Core message

Represent authoritative entities and labeled relationships so retrieval can follow known
dependencies and provenance rather than similarity alone.

## Value

Supports explainable multi-hop traversal for questions about topology, impact, and ownership.

## Assumptions

- Edge extraction and freshness are accurate enough to be useful.
- A graph complements rather than silently outranks primary sources.

## Evidence and collisions

Microsoft's GraphRAG research supports gains for a class of global sensemaking queries,
not universal superiority over vector or routed retrieval. Source 006 defers semantic and
vector retrieval, and Spectacular has no thin-slice benchmark.

## Trade-offs and recommendation

Structural traversal can improve recall while adding extraction, correction, and freshness
cost. Trial a small accepted-contract projection against current routing before adoption.

## Decision

Pending.
