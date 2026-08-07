---
type: concept-piece
id: PZL-090
status: captured
domain: context-retrieval
sources: [source-007]
source_authority: synthesized-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-080, PZL-083]
overlaps_with: [PZL-001, PZL-010, PZL-070]
conflicts_with: []
tags: [retrieval, graph, semantic-search, authority]
updated: 2026-08-07
---

# Authoritative graph-first retrieval

## Core message

Compile context by following accepted contract relationships first; use semantic
search only to find candidates, never to determine authority.

## Value

Makes progressive disclosure explainable and less vulnerable to similar but stale text.

## Assumptions

- Graph edges are complete and fresh enough for safe retrieval.
- Source files and decisions can override graph summaries visibly.

## Evidence and collisions

The AI-SDLC attachment recommends vector indexing, while Source 007 correctly narrows
it to secondary discovery. No benchmark compares graph and current routed retrieval.

## Trade-offs and recommendation

Deterministic authority paths versus missed context from incomplete edges. Benchmark a
thin graph against current routing before adopting it as the primary compiler.

## Decision

Pending.
