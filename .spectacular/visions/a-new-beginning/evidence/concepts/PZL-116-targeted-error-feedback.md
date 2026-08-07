---
type: concept-piece
id: PZL-116
status: captured
domain: diagnostic-loop
sources: [source-008, source-009, source-010]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-096, PZL-110]
overlaps_with: [PZL-095]
conflicts_with: []
tags: [failure, feedback, checks, repair]
updated: 2026-08-07
---

# Targeted error-feedback repair

## Core message

Feed a failed deterministic check and its exact evidence into a bounded diagnostic attempt;
require a new hypothesis and rerun the narrowest relevant check before broad verification.

## Value

Connects the harness to purposeful repair instead of repeated blind regeneration.

## Assumptions

- Check output is stable, scoped, and safe to expose.
- Retry budgets and known-good checkpoints exist.

## Evidence and collisions

This composes Source 008's deterministic feedback with PZL-096's hypothesis rule. Piping the
same error back without a changed hypothesis merely automates an infinite loop.

## Trade-offs and recommendation

Automated repair shortens feedback while consuming attempts on misleading failures. Treat
error output as evidence in the existing retry state machine, not a new autonomous subsystem.

## Decision

Pending.
