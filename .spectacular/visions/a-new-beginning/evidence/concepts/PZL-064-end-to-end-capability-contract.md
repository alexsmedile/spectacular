---
type: concept-piece
id: PZL-064
status: captured
domain: specification-model
sources: [source-006, source-007, source-014]
source_authority: proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-013, PZL-038]
overlaps_with: [PZL-013, PZL-019, PZL-065, PZL-080, PZL-081, PZL-083]
conflicts_with: []
tags: [capability, contract, flow, state, verification]
updated: 2026-08-07
---

# End-to-end capability contract

## Core message

Describe each meaningful capability in one compact contract covering outcome,
flow, invariants, state, failures, implementation links, and verification.

## Value

Connects user behavior to code and proof without fragmenting context by component.

## Assumptions

- End-to-end capability is a stable aggregation boundary.
- Implementation links can remain current enough to aid agents.

## Evidence and collisions

Current specs already describe capabilities and code remains authoritative. A new
`capabilities/` path may clarify semantics or merely rename and duplicate `specs/`.
Source 007 clarifies that capabilities cross architectural boundaries rather than
owning component topology themselves.

## Trade-offs and recommendation

High-cohesion contracts versus large files for cross-cutting systems. Test the
schema on several real capabilities before changing the canonical path.

## Decision

Pending.
