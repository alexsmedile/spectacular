---
type: concept-piece
id: PZL-120
status: captured
domain: information-architecture
sources: [source-009]
source_authority: owner-endorsed-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-003, PZL-086]
overlaps_with: [PZL-108]
conflicts_with: []
tags: [mermaid, visualization, progressive-disclosure, projection]
updated: 2026-08-07
---

# Derived Mermaid decision map

## Core message

Render the decision-critical concept and dependency graph as a compact Mermaid view whose
stable node IDs resolve to the authoritative Markdown cards.

## Value

Lets a human see convergence, forks, blockers, and decision order without loading 122 cards.

## Assumptions

- The diagram is generated or checked against card metadata.
- Only relationships that affect a decision appear; the map is not a decorative inventory.

## Evidence and collisions

The owner specifically identified Mermaid nodes as a useful summarization layer. PZL-108
supports pointer-first navigation. A manually maintained diagram could become another stale truth.

## Trade-offs and recommendation

Visual compression improves orientation but hides nuance and can become unreadable at scale.
Create one curated decision-frontier projection plus a linked table, capped around 10–20 nodes.

## Decision

Pending.
