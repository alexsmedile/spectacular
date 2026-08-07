---
type: concept-piece
id: PZL-140
status: captured
domain: continuity
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-075, PZL-118]
overlaps_with: [PZL-113]
conflicts_with: []
tags: [terminal, next-action, blocked, handoff, status]
updated: 2026-08-07
---

# Terminal next-action contract

## Core message

When a run genuinely stops, its terminal output must state exactly one safe next action or the
specific condition required to unblock it.

## Value

Eliminates ambiguous endings and makes continuation possible without reconstructing the session.

## Assumptions

The contract applies only at a terminal or blocked boundary, not after every intermediate tool call.

## Evidence and collisions

Issue #33 records direct usage friction. It composes with cold resume and crash-resumable state but
must not cause an agent to stop while approved work remains.

## Trade-offs and recommendation

Define a small output schema: state, reason, next action, owner, and evidence pointer. Verify it on
blocked and completed scenarios.

## Decision

Pending.
