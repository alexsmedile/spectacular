---
type: concept-piece
id: PZL-137
status: captured
domain: git-boundary
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-044, PZL-054]
overlaps_with: []
conflicts_with: [PZL-044]
tags: [afk, commits, checkpoints, autonomy, git]
updated: 2026-08-07
---

# Goal-scoped commit discipline

## Core message

An autonomous run should checkpoint coherent, goal-linked progress rather than leave an opaque
working tree or manufacture commits on a timer.

## Value

Improves recovery, reviewability, and traceability during long unattended work.

## Assumptions

Spectacular retains authority over autonomous execution rather than delegating Git mutations to a
companion skill or native tooling.

## Evidence and collisions

Issue #29 identifies a valid commit-layer gap, but AFK ownership is contested by PZL-044 and the
native-provider boundary in PZL-054.

## Trade-offs and recommendation

Preserve the invariant as an execution contract. Defer its implementation location until the AFK
survival and Git-authority decisions are made.

## Decision

Pending.
