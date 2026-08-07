---
type: concept-piece
id: PZL-103
status: captured
domain: orchestration-model
sources: [source-008, source-009, source-010]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-089, PZL-099]
overlaps_with: [PZL-096, PZL-098]
conflicts_with: []
tags: [dag, state-machine, retry, orchestration]
updated: 2026-08-07
---

# Dependency DAG plus bounded attempts

## Core message

Use a DAG to represent prerequisite, fan-out, and join structure; use explicit bounded
state transitions for each execution attempt and retry.

## Value

Preserves dependency visibility without pretending that diagnostics and recovery are acyclic.

## Assumptions

- The work is complex enough for explicit dependency structure to repay maintenance.
- Attempt state and graph-node state have defined ownership.

## Evidence and collisions

Source 008's “graph over loop” is a false dichotomy. Source 007 already proposes an
objective DAG and separate run states; neither source demonstrates a production graph.

## Trade-offs and recommendation

Graphs improve joins and parallel eligibility but add scheduling and persistence cost.
Adopt the composition rule first; require a thin-slice measurement before a graph engine.

## Decision

Pending.
