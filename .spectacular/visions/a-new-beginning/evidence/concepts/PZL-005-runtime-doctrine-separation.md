---
type: concept-piece
id: PZL-005
status: captured
domain: reference-retention
sources: [source-001, source-003, source-004]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-004]
overlaps_with: [PZL-004]
conflicts_with: []
tags: [doctrine, retention, docs, runtime-context]
updated: 2026-08-07
---

# Runtime/doctrine separation

## Core message

Keep explanatory doctrine, glossary, historical rationale, and human-maintainer
material outside the normal agent runtime retrieval path.

## Value

Allows rich maintainership documentation without charging every operational task
for it in context.

## Assumptions

- Runtime workflow docs remain understandable without nearby doctrine.
- A discoverable, non-runtime home exists for retained rationale.

## Evidence and collisions

The goal aligns with progressive disclosure. The proposed destination `docs/`
conflicts with the current pageworks-owned public documentation boundary, and
some doctrine may still be needed for rare routing judgments.

## Trade-offs and recommendation

Lean runtime context versus rationale becoming stale or undiscoverable. Accept the
separation goal for comparison, but choose placement and link policy separately.
Do not move material into `docs/` under the current boundary.

## Decision

Pending.
