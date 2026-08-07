---
type: concept-piece
id: PZL-115
status: captured
domain: approval-authority
sources: [source-008, source-009, source-012, source-013, source-014]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-067, PZL-087]
overlaps_with: [PZL-040, PZL-097]
conflicts_with: []
tags: [human-review, risk, approval, consequence]
updated: 2026-08-07
---

# Risk-calibrated human gates

## Core message

Move human judgment earlier for high-consequence direction and architecture; allow bounded,
reversible work to proceed to an evidence-backed review point.

## Value

Spends attention where mistakes are expensive without applying identical ceremony to every change.

## Assumptions

- Risk classification includes security, privacy, compatibility, irreversibility, and scope.
- End review never retroactively authorizes forbidden action.

## Evidence and collisions

Converges with the existing intent gate and risk hooks. Source 008's “small bug fix” example
cannot bypass early review when a small diff crosses a sensitive boundary.

## Trade-offs and recommendation

Calibration reduces ceremony but misclassification can remove necessary consent. Define a
small conservative risk rubric and explicit escalation triggers before shifting gates.

## Decision

Pending.
