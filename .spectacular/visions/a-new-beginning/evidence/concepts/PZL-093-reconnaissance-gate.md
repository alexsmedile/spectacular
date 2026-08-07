---
type: concept-piece
id: PZL-093
status: captured
domain: execution-gate
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-070]
overlaps_with: [PZL-068, PZL-092]
conflicts_with: []
tags: [reconnaissance, patterns, dependencies, tests, git]
updated: 2026-08-07
---

# Reconnaissance gate

## Core message

Before editing, require evidence that the run inspected relevant boundaries, nearby
patterns, dependency versions, tests, Git state, and project instructions.

## Value

Reduces invented architecture, deprecated APIs, and accidental inconsistency.

## Assumptions

- Inspection can be evidenced without logging performative file lists.
- Bounded direct changes can satisfy the gate cheaply.

## Evidence and collisions

The senior-harness attachment explicitly supports these inputs; current implementation
policy already requires understanding before change.

## Trade-offs and recommendation

Safer changes versus preflight ceremony. Scale required evidence by risk and novelty,
not by a universal checklist length.

## Decision

Pending.
