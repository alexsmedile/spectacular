---
type: concept-piece
id: PZL-123
status: captured
domain: abstraction-architecture
sources: [source-010]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-051, PZL-105, PZL-119]
overlaps_with: [PZL-001, PZL-070]
conflicts_with: []
tags: [abstraction, context, interface, escape-hatch]
updated: 2026-08-07
---

# Layer-specific context contract

## Core message

Each abstraction layer declares the minimum context it consumes, authority it exposes,
output it guarantees, failure it reports, and path to deeper evidence.

## Value

Reduces context load without turning hidden implementation into missing safety information.

## Assumptions

- Boundaries follow stable responsibilities rather than arbitrary file or agent divisions.
- Callers can detect insufficient abstraction and escalate deliberately.

## Evidence and collisions

Source 010 makes context survival the purpose of abstraction. Existing skill routing and
bounded run compilation support the idea, but current layers have not been audited for leaks.

## Trade-offs and recommendation

Strict layers improve local reasoning while hiding cross-cutting constraints. Audit every
boundary for minimum context plus an explicit provenance and drill-down escape hatch.

## Decision

Pending.
