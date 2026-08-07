---
type: concept-piece
id: PZL-084
status: captured
domain: uncertainty-model
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-083]
overlaps_with: [PZL-039, PZL-072, PZL-100]
conflicts_with: []
tags: [unknowns, assumptions, questions, discovery, fail-closed]
updated: 2026-08-07
---

# Named missing-contract escalation

## Core message

An agent must not silently infer missing contract information; it records an
assumption, asks a question, or creates a bounded discovery objective.

## Value

Makes uncertainty visible before it becomes architectural invention or scope drift.

## Assumptions

- The system can distinguish harmless local inference from contract-level uncertainty.
- Named unknowns can be owned without recreating several global databases.

## Evidence and collisions

Current question/research/spike semantics already protect different unknowns. The
proposal changes their ownership, not the need for explicit uncertainty.

## Trade-offs and recommendation

Safer reasoning versus interruption and record churn. Preserve the invariant;
decide storage and escalation threshold separately.

## Decision

Pending.
