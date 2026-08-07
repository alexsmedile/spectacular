---
type: concept-piece
id: PZL-162
status: captured
domain: ai-ux-pattern
sources: [source-013]
source_authority: user-provided-unsourced-proposal-plus-owner-hypothesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-160]
overlaps_with: [PZL-095]
conflicts_with: []
tags: [graceful-degradation, fallback, structured-ui, hallucination, latency]
updated: 2026-08-07
---

# Graceful structured degradation

## Core message

When an AI path is unavailable, slow, ambiguous, or unsafe, preserve task completion through a clear
structured/manual path rather than trapping the user behind an error message or retry loop.

## Value

Keeps the product usable when probabilistic behavior fails and makes the non-AI contract explicit.

## Assumptions

A structured fallback can accomplish the essential outcome and shares compatible state with the AI path.

## Evidence and collisions

Source 013 identifies this as a resilience opportunity. Maintaining two complete interfaces can double
surface and drift unless they share one underlying contract.

## Trade-offs and recommendation

Define the essential non-AI completion path first, then layer assistance over it. Avoid duplicating
business rules or state transitions between AI and manual interfaces.

## Decision

Pending.
