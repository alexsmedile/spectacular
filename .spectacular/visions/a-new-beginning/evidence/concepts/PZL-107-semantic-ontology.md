---
type: concept-piece
id: PZL-107
status: captured
domain: semantic-contract
sources: [source-008, source-009, source-010]
source_authority: unsourced-protocol-synthesis
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-019, PZL-083, PZL-086]
overlaps_with: [PZL-080, PZL-082]
conflicts_with: []
tags: [ontology, semantics, constraints, truth]
updated: 2026-08-07
---

# Explicit semantic ontology

## Core message

Make important business meanings, relationships, and constraints explicit above the
physical data model so agents do not have to invent them from schema names.

## Value

Separates what stored data is from what the product accepts that data to mean.

## Assumptions

- The selected semantics are small, reviewable, and maintained with implementation change.
- Runtime truth remains distinguishable from accepted contract semantics.

## Evidence and collisions

Source 008's “infallible world model” is refuted: an ontology can be incomplete, stale,
ambiguous, or wrong. It overlaps the typed contract graph without proving another layer.

## Trade-offs and recommendation

Explicit semantics reduce inference but can become centralized schema bureaucracy. Embed a
minimal vocabulary and invariants in existing contracts before creating an ontology system.

## Decision

Pending.
