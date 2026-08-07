---
type: concept-piece
id: PZL-170
status: captured
domain: implementation-architecture
sources: [source-014]
source_authority: user-provided-unsourced-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-111, PZL-125]
overlaps_with: [PZL-059, PZL-123, PZL-124]
conflicts_with: []
tags: [optionality, adapter, simple-first, volatility, architecture]
updated: 2026-08-07
---

# Earned optionality seams

## Core message

Preserve reversibility through narrow interfaces at demonstrated volatile, external, or high-cost
boundaries while keeping the rest of the architecture as simple as current requirements allow.

## Value

Reduces vendor and migration lock-in without paying abstraction cost at every hypothetical boundary.

## Assumptions

- Volatility and replacement cost can be evidenced from contracts, history, or external ownership.
- A simple direct implementation may be wrapped later before its coupling becomes irreversible.

## Evidence and collisions

Classical modularity and adapter boundaries support isolation, but Source 014's blanket advice to
abstract third parties can produce speculative interfaces. “Monolith first” and “boring technology”
are useful defaults only when constraints, team fit, and failure modes support them.

## Trade-offs and recommendation

Optionality reduces future migration cost but creates present indirection. Require a named volatile
boundary and replacement scenario before adding an adapter; otherwise prefer the simplest deep
module that satisfies the accepted contract.

## Decision

Pending.
