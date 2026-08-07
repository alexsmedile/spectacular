---
type: concept-piece
id: PZL-067
status: captured
domain: approval-authority
sources: [source-006, source-007]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-066]
overlaps_with: [PZL-072, PZL-076, PZL-087]
conflicts_with: [PZL-087]
tags: [approval, lock, autonomy, hitl]
updated: 2026-08-07
---

# One mandatory execution approval

## Core message

For an MVP, require one human confirmation of mission behavior, boundaries, and
business logic, then lock that baseline for autonomous execution.

## Value

Reduces approval fatigue and gives the run one clear authority artifact.

## Assumptions

- Stop conditions preserve safety after the initial lock.
- Spec, activation, irreversible, closure, and release authorities can remain distinct only when needed.

## Evidence and collisions

Current workflows deliberately separate direction approval, activation, HITL
gates, verification, closure, PR readiness, and merge. Their redundancy has not
been demonstrated. Source 007 directly replaces the single-lock proposal with a
human product lock plus independent engineering assurance.

## Trade-offs and recommendation

Low-friction autonomy versus consent expansion. Model each existing gate by risk
and consequence before collapsing any approvals.

## Decision

Pending.
