---
type: concept-piece
id: PZL-098
status: captured
domain: execution-architecture
sources: [source-007]
source_authority: synthesized-proposal
assessment: mixed
evidence_status: unverified
disposition: pending
depends_on: [PZL-087, PZL-097, PZL-099]
overlaps_with: [PZL-041, PZL-071]
conflicts_with: [PZL-071, PZL-079]
tags: [review, ci, deployment, observation, risk, agents]
updated: 2026-08-07
---

# Risk-triggered staged run chain

## Core message

A Mission may use separate implementation, review, CI, deployment, observation, and
reconciliation runs, adding specialist review only when affected contracts justify it.

## Value

Matches permissions and evidence to lifecycle risk while suppressing low-value review noise.

## Assumptions

- Multiple agents/runs improve independence without creating orchestration overhead.
- Provider integration is available for later boundaries.

## Evidence and collisions

Source 006 limits MVP to one coding agent and defers scheduling; Source 007 says not
to generalize scheduling while proposing a typed multi-run chain.

## Trade-offs and recommendation

Independent assurance versus platform creep. Defer the generalized chain; thin-slice
one implementation attempt plus one read-only verification report.

## Decision

Pending.
