---
type: concept-piece
id: PZL-160
status: captured
domain: ux-safety
sources: [source-013]
source_authority: user-provided-unsourced-proposal-plus-owner-hypothesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-072, PZL-115]
overlaps_with: [PZL-094, PZL-155, PZL-157]
conflicts_with: []
tags: [consequence, failure-mode, fallback, reversibility, user-control]
updated: 2026-08-07
---

# Consequence-first UX audit

## Core message

For every AI-mediated action, identify the primary failure mode, irreversible user consequence,
fallback behavior, reversibility, and retained user control before choosing the interaction model.

## Value

Shifts design review from polished happy paths to the mistakes that determine trust and safety.

## Assumptions

Consequences can be ranked by affected user, severity, detectability, recoverability, and frequency.

## Evidence and collisions

Source 013 usefully asks for irreversible consequence but does not define risk scoring. Low-risk
reversible interactions should not inherit the ceremony of destructive or rights-affecting actions.

## Trade-offs and recommendation

Make consequence analysis proportional. Require explicit previews, confirmation, undo, or manual
fallback only where the risk model justifies their friction.

## Decision

Pending.
