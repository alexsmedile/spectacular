---
type: concept-piece
id: PZL-132
status: captured
domain: information-architecture
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-120]
overlaps_with: [PZL-108]
conflicts_with: []
tags: [ascii, mermaid, renderer, projection, workflows]
updated: 2026-08-07
---

# Derived visual renderer

## Core message

Render compact ASCII or Mermaid views from canonical state; diagrams are navigational projections,
not another editable source of truth.

## Value

Makes lifecycle, dependencies, and current position legible without loading the underlying corpus.

## Assumptions

Every visual node can point back to a stable record or field.

## Evidence and collisions

Issue #18 provides the narrow renderer prototype; #24 wants the same capability over live Mission
state. PZL-120 is already a working static example in this refactor.

## Trade-offs and recommendation

Prototype one shared renderer over an existing lifecycle before binding it to the contested Mission
model. Never hand-edit generated diagrams.

## Decision

Pending.
