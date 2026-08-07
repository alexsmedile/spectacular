---
type: concept-piece
id: PZL-008
status: captured
domain: kits
sources: [source-001, source-003]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-006]
overlaps_with: [PZL-007, PZL-009]
conflicts_with: []
tags: [kits, prd, suggestions, scaffold]
updated: 2026-08-07
---

# Kits as PRD presets

## Core message

Limit kits to shaping the PRD and providing contextual suggestions; do not let a
kit eagerly materialize additional canonical documents.

## Value

Keeps init consistent and makes kit choice about product-question relevance
rather than pre-committing the workspace structure.

## Assumptions

- Contextual suggestions appear at reliable moments without becoming nagging.
- Kit-specific docs are not required before the first real operation.

## Evidence and collisions

The coding kit declares STACK and ARCHITECTURE as always-triggered, yet a fresh
observed coding init created only the always-set. Current behavior, kit contract,
and documentation already disagree, so compatibility cannot be inferred.

## Trade-offs and recommendation

Lower scaffold burden versus delayed architectural context and suggestion-state
complexity. Diagnose the existing trigger drift and gather actual kit use before
choosing whether to repair or retire structural triggers. Mixed.

## Decision

Pending.
