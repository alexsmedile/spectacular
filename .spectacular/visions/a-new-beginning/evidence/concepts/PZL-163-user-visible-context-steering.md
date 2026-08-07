---
type: concept-piece
id: PZL-163
status: captured
domain: ai-ux-pattern
sources: [source-013]
source_authority: user-provided-unsourced-proposal-plus-owner-hypothesis
assessment: mixed
evidence_status: unverified
disposition: pending
depends_on: [PZL-123, PZL-124]
overlaps_with: [PZL-109]
conflicts_with: []
tags: [context, memory, transparency, steering, privacy]
updated: 2026-08-07
---

# User-visible context steering

## Core message

Expose the consequential context an AI is using in a form users can inspect, remove, correct, or
scope—without pretending token counts or hidden model state are understandable product semantics.

## Value

Can improve correction, privacy, and trust when context selection materially affects outcomes.

## Assumptions

The product can name context in user-meaningful units and changes actually influence subsequent behavior.

## Evidence and collisions

Source 013 proposes a context-budget UI without user evidence. Raw tokens, embeddings, or retrieval
internals may confuse users and expose sensitive information without giving real control.

## Trade-offs and recommendation

Expose sources, scopes, remembered facts, and permissions rather than infrastructure metrics. Prototype
only where users need correction or consent, and verify that controls have observable effect.

## Decision

Pending.
