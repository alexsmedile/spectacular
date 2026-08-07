---
type: concept-piece
id: PZL-161
status: captured
domain: ai-ux-pattern
sources: [source-013]
source_authority: user-provided-unsourced-proposal-plus-owner-hypothesis
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-160]
overlaps_with: [PZL-072, PZL-149]
conflicts_with: []
tags: [intent-preview, ai-action, transparency, confirmation, anxiety]
updated: 2026-08-07
---

# AI-action intent preview

## Core message

Preview an AI action's intended effect, scope, and reversibility before consequential execution so
the user can correct misunderstanding without inspecting implementation detail.

## Value

Can preserve speed while reducing surprise and anxiety around opaque automated action.

## Assumptions

The preview faithfully represents the actual effect and is shown only when correction remains possible.

## Evidence and collisions

Source 013 presents preview as a general friction reducer. Mandatory previews on trivial reversible
actions would recreate micro-approval fatigue.

## Trade-offs and recommendation

Trigger previews by consequence, novelty, scope, or low confidence. For routine safe actions, use
lightweight progress and undo instead of a blocking confirmation.

## Decision

Pending.
