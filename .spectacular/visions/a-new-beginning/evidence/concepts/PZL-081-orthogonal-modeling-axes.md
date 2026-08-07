---
type: concept-piece
id: PZL-081
status: captured
domain: system-model
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-064]
overlaps_with: [PZL-013, PZL-082]
conflicts_with: []
tags: [capability, component, interface, state, axes]
updated: 2026-08-07
---

# Orthogonal capability and architecture axes

## Core message

Capabilities describe observable outcomes; components, interfaces, and state describe
how the system realizes them. Do not substitute one axis for the other.

## Value

Prevents user-facing contracts from becoming fake architecture maps and vice versa.

## Assumptions

- Both axes are relevant to the project's work.
- Cross-axis relationships can remain understandable in Markdown.

## Evidence and collisions

The distinction is logically sound and supported by the checkout example. Current
spec shapes already include capability and flat technical contracts without explicit edges.

## Trade-offs and recommendation

Clear modeling versus additional concepts. Adopt the distinction as a design rule;
first-class storage still depends on demonstrated retrieval need.

## Decision

Pending.
