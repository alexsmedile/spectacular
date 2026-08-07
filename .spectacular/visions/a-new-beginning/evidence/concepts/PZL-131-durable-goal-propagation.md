---
type: concept-piece
id: PZL-131
status: captured
domain: intent-contract
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-066, PZL-068]
overlaps_with: [PZL-119]
conflicts_with: []
tags: [goal, why, request, pull-request, propagation]
updated: 2026-08-07
---

# Durable goal propagation

## Core message

Carry the approved outcome and reason from request creation through run state, commits, and the
delivery artifact so each handoff retains the original intent.

## Value

Reviewers and resumed agents can judge whether a change solves the intended problem.

## Assumptions

Presence and bounded length can be validated mechanically; semantic correctness remains a review
judgment.

## Evidence and collisions

Issues #17 and #32 describe the same missing intent thread at different surfaces. The proposal
should not pretend that a script can validate whether prose is “meaningful.”

## Trade-offs and recommendation

Avoid copying mutable prose into many independent fields. Give the request one canonical goal and
derive concise delivery leads from it, recording deliberate overrides.

## Decision

Pending.
