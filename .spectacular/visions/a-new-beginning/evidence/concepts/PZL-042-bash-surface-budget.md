---
type: concept-piece
id: PZL-042
status: captured
domain: implementation-architecture
sources: [source-004, source-005]
source_authority: proposal
assessment: promising
evidence_status: supported
disposition: pending
depends_on: [PZL-035, PZL-038]
overlaps_with: [PZL-010, PZL-035, PZL-059]
conflicts_with: []
tags: [bash, cli, size, maintenance, rewrite]
updated: 2026-08-07
---

# Bash surface budget

## Core message

Treat the single-file Bash CLI's size as a maintenance signal and choose its future
architecture only after the retained product surface is known.

## Value

Makes implementation complexity visible and avoids porting capability that will be removed.

## Assumptions

- Line count correlates with maintenance cost and defect risk.
- A smaller accepted surface can materially reduce the CLI.

## Evidence and collisions

The CLI is currently 16,507 lines. The proposed ~6k target has no behavior-derived
estimate; source sessions predict 9–11k before a later decision. Size must not drive
removal of load-bearing behavior. Source 005 proposes modular Bash sources and a
single assembled executable, which could improve maintainability without directly
reducing the line count.

## Trade-offs and recommendation

Useful budget and sequencing versus line-count gaming. Establish capability, test,
startup, portability, and maintainability criteria before freeze/extract/rewrite.
Promising metric, speculative target.

## Decision

Pending.
