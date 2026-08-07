---
type: concept-piece
id: PZL-044
status: captured
domain: autonomy-coordination
sources: [source-004, source-005, source-006]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-037]
overlaps_with: [PZL-027, PZL-031, PZL-048, PZL-050, PZL-054, PZL-069, PZL-071, PZL-072]
conflicts_with: []
tags: [afk, worktree, traffic, autonomy, git]
updated: 2026-08-07
---

# Retire AFK coordination

## Core message

Remove AFK, workspace coordination, traffic preflight, and worktree infrastructure
as speculative autonomous-agent Git machinery.

## Value

Would cut a wide, safety-heavy subsystem from the CLI, references, records, and tests.

## Assumptions

- Users do not need bounded unattended work or concurrent worktree protection.
- Existing Git recovery/safety can be preserved more simply.

## Evidence and collisions

AFK has two local run records and a branch ledger; related requests are verified,
not in flight. The features are recently shipped, so sparse usage is weak evidence.
Their safety value may appear as prevented damage rather than record volume.
Source 005 offers a narrower alternative: keep AFK authorization records, move
read-only safety interpretation to workspace, and delegate mutations to Git/gh.
Source 006 offers another replacement: fold bounded autonomy into one mission run
using the host coding runtime and explicitly defer broader orchestration.

## Trade-offs and recommendation

Large simplification versus removing explicit authority and Git safety boundaries.
Compare real workflows and simpler substitutes before deciding. Disputed, not dead.

## Decision

Pending.
