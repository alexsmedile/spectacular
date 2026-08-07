---
type: concept-piece
id: PZL-125
status: captured
domain: implementation-quality
sources: [source-010]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-110, PZL-111]
overlaps_with: [PZL-094, PZL-095]
conflicts_with: []
tags: [grey-box, modules, delegation, invariants]
updated: 2026-08-07
---

# Grey-box delegation with internal invariants

## Core message

Delegate implementation behind an accepted interface while still enforcing internal security,
cohesion, observability, performance, testability, and maintenance invariants.

## Value

Allows bounded autonomy without treating encapsulation as permission for hidden technical debt.

## Assumptions

- Interface and internal-quality evidence are defined before delegation.
- Review can inspect internals when an invariant or failure requires it.

## Evidence and collisions

Source 010 says messy internals are acceptable if the boundary holds. That conflicts directly
with the anti-slop goal and ignores debugging, security, and future modification costs.

## Trade-offs and recommendation

Grey-box scope reduces human micromanagement but can conceal accumulating debt. Approve
outcomes and interfaces while retaining a small internal-quality rubric and evidence gate.

## Decision

Pending.
