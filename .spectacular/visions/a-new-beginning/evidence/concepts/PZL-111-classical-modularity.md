---
type: concept-piece
id: PZL-111
status: captured
domain: implementation-quality
sources: [source-008, source-009, source-010, source-014]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-042, PZL-059]
conflicts_with: []
tags: [modularity, cohesion, coupling, ai-slop]
updated: 2026-08-07
---

# Classical modularity as anti-slop constraint

## Core message

Judge agent-produced designs and code by ordinary engineering qualities—cohesion,
encapsulation, separation of concerns, narrow interfaces, and maintainability.

## Value

Prevents fast specification-to-code output from becoming an excuse for duplicated,
shallow, or tightly coupled mechanisms.

## Assumptions

- Quality is evaluated in repository context rather than through slogans alone.
- “Deep module” guidance does not justify oversized hidden complexity.

## Evidence and collisions

The repo's monolithic CLI and duplicated command paths make modularity a live concern.
Source 008 provides principles rather than a measurable refactor boundary.

## Trade-offs and recommendation

Modularity improves change isolation but premature abstraction adds indirection. Require
each proposed module to own a stable responsibility and remove actual duplication.

## Decision

Pending.
