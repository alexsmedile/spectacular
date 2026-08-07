---
type: concept-piece
id: PZL-040
status: captured
domain: policy-governance
sources: [source-004, source-006]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-034, PZL-037]
overlaps_with: [PZL-006, PZL-034, PZL-072, PZL-079]
conflicts_with: []
tags: [policy, hooks, risk, enforcement]
updated: 2026-08-07
---

# Three risk-gate policy model

## Core message

Replace phase-specific policy hooks with three risk boundaries: before write,
before archive, and before irreversible action.

## Value

Could retain meaningful safety while removing a mandatory policy lookup from every phase.

## Assumptions

- Planning, debugging, implementation, verification, and memory risks map cleanly
  onto those three boundaries.
- Policy injection cost is material in actual agent runs.

## Evidence and collisions

Current POLICY defines nine hooks, not eight. The source provides no block/warn
activation history. Some current gates protect correctness rather than write irreversibility.
Source 006 defers custom policy language while retaining mission permissions and
explicit stop conditions, offering a contract-level alternative rather than
evidence specifically for three hooks.

## Trade-offs and recommendation

Simpler model and fewer calls versus loss of phase-specific safeguards. Map every
current directive and observed failure prevented before choosing consolidation.
Mixed.

## Decision

Pending.
