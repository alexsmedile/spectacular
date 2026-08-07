---
type: concept-piece
id: PZL-150
status: captured
domain: decision-method
sources: [source-012, source-013]
source_authority: user-provided-unsourced-proposal
assessment: promising
evidence_status: unverified
disposition: pending
depends_on: [PZL-121]
overlaps_with: [PZL-010, PZL-126, PZL-133]
conflicts_with: []
tags: [decision-density, human-attention, evaluation, compression]
updated: 2026-08-07
---

# Validated decision density

## Core message

Measure how many meaningful, retained decisions a human can confidently resolve per unit of
attention—not how many prompts, options, or apparent choices an artifact contains.

## Value

Makes human review bandwidth an explicit design constraint for decision workflows.

## Assumptions

- Decisions can be identified without splitting one choice into artificial micro-decisions.
- Retention, correctness, and reversal cost can be checked after the review pass.

## Evidence and collisions

Source 012 claims 10–20 decisions per pass without evidence. Raw density can hide ambiguity,
encourage shallow agreement, and conflict with the cold-context correctness budget in PZL-010.

## Trade-offs and recommendation

Use validated decision density as a secondary S01 metric alongside decision error, later reversal,
time-to-understanding, fatigue, and artifact cost. Never optimize the count alone.

## Decision

Pending.
