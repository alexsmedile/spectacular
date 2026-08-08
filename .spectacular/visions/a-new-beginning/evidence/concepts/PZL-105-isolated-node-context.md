---
type: concept-piece
id: PZL-105
status: captured
domain: context-compilation
sources: [source-008, source-009, source-010, source-015]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-070]
overlaps_with: [PZL-091, PZL-104]
conflicts_with: []
tags: [context, isolation, specialization, contracts]
updated: 2026-08-08
---

# Isolated bounded node context

## Core message

Give a specialized worker the smallest authoritative context that supports its bounded
job, plus a shared output and safety contract.

## Value

Reduces inherited noise, stale conversation assumptions, and accidental scope expansion.

## Assumptions

- The compiler can identify all safety-critical context.
- Outputs carry provenance sufficient for integration.

## Evidence and collisions

Converges with PZL-070 and PZL-091. Isolation alone is unsafe when the bounded context
omits an applicable policy, interface, dependency, or decision.

## Trade-offs and recommendation

Clean context reduces rot but repeats orientation and can hide cross-cutting constraints.
Specify a minimum shared envelope and test omission failures before broad adoption.

## Decision

Pending.
