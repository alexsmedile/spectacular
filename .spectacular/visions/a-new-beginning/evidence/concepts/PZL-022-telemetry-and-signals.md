---
type: concept-piece
id: PZL-022
status: captured
domain: feedback-model
sources: [source-002]
source_authority: proposal
assessment: mixed
evidence_status: unverified
disposition: pending
depends_on: [PZL-018]
overlaps_with: [PZL-019]
conflicts_with: []
tags: [telemetry, signals, feedback, operations, evidence]
updated: 2026-08-07
---

# Telemetry and Signals

## Core message

Use Telemetry or the more neutral Signals for evidence produced during real-world operation.

## Value

Names the feedback channel from actual use back into contracts and project context.

## Assumptions

- Operational signals exist for all relevant project types.
- The label stays distinct from human feedback, verification, research, and benchmarks.

## Evidence and collisions

Spectacular currently treats feedback, verification, and benchmarks as orthogonal.
Telemetry commonly implies quantitative runtime instrumentation, while the source
also includes qualitative observations. A single term could erase these boundaries.

## Trade-offs and recommendation

Strong operations metaphor versus category collapse. Prefer Signals as an umbrella
only if subtypes retain their existing evidence contracts; otherwise keep the
current explicit terms. Mixed.

## Decision

Pending.
