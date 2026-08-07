---
type: concept-piece
id: PZL-118
status: captured
domain: continuity
sources: [source-008, source-009, source-010]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-075, PZL-099]
overlaps_with: [PZL-108]
conflicts_with: []
tags: [state, resume, checkpoint, crash-recovery]
updated: 2026-08-07
---

# Crash-resumable external state

## Core message

Persist authoritative lifecycle position, baseline, completed evidence, current objective,
and next safe action so execution can resume after process or context loss.

## Value

Makes continuity a testable state contract instead of a property of chat history.

## Assumptions

- Durable state is atomically updated and distinguishes canonical fields from raw artifacts.
- Resume verifies Git and environmental baselines before continuing.

## Evidence and collisions

Converges strongly with PZL-075 and PZL-099. Source 008 does not distinguish committed
project state from sensitive, local, transient run material.

## Trade-offs and recommendation

External state enables recovery but can drift or leak operational detail. Define the minimum
canonical resume manifest and store raw artifacts under explicit retention classes.

## Decision

Pending.
