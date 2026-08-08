---
type: concept-piece
id: PZL-172
status: captured
domain: closure-lifecycle
sources: [source-006, source-007, source-015]
source_authority: proposal-plus-primary-comparative-study
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-066, PZL-073, PZL-083]
overlaps_with: [PZL-074, PZL-092, PZL-095, PZL-171]
conflicts_with: []
tags: [evidence, reconciliation, lifecycle, authority, freshness, deterministic-gate]
updated: 2026-08-08
---

# Machine-checkable evidence reconciliation envelope

## Core message

Before a lifecycle promotion, deterministically validate an envelope linking the accepted baseline
and contract delta to the requested transition, authorized actor, required evidence clauses,
evidence references, results, attribution, and freshness; escalate semantic uncertainty rather than
letting document presence imply proof.

## Value

Makes the evidence-to-lifecycle edge observable and prevents prompts, stale summaries, or orphaned
test output from silently authorizing closure.

## Assumptions

- Closed properties such as identity, attribution, timestamps, required fields, and referenced
  results can be checked without pretending to prove all semantics.
- S03 defines truth and freshness, S05 defines authority, and S06 defines evidence sufficiency.

## Evidence and collisions

Sources 006 and 007 require evidence-backed reconciliation but do not define a minimal executable
envelope. Source 015 shows real gate and review implementations plus a concrete prompt/manifest
drift example. None of the studied targets proves the proposed envelope's user value or minimum
schema. PZL-171 demonstrates the same boundary: structural checks prove declared properties, not
all behavioral meaning.

## Trade-offs and recommendation

Reliable promotion versus a box-checking schema. Adopt the responsibility in S06, but define the
minimum envelope only after truth and authority settle; test it with missing, stale, misattributed,
failed, and semantically disputed fixtures.

## Decision

Pending.
