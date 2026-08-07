---
type: concept-piece
id: PZL-130
status: captured
domain: context-retrieval
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-121, PZL-129]
overlaps_with: [PZL-090, PZL-106, PZL-124]
conflicts_with: [PZL-079]
tags: [semantic-retrieval, provider, optional, advisory]
updated: 2026-08-07
---

# Optional advisory retrieval provider

## Core message

Allow semantic or graph retrieval as a replaceable discovery aid after deterministic projections
are stable; it may suggest context but never establish authority.

## Value

Could improve discovery in large workspaces without coupling the core to one index technology.

## Assumptions

A provider contract can preserve Markdown, code, schemas, and accepted decisions as higher-ranked
sources.

## Evidence and collisions

Issue #13 is intentionally gated by #11 and #12. It converges with graph-first retrieval but
collides with the MVP deferral fence and currently has no measured retrieval failure to solve.

## Trade-offs and recommendation

Keep the interface hypothetical until deterministic retrieval has benchmarked misses. Then test
one optional adapter against the same corpus rather than adding a mandatory substrate.

## Decision

Pending.
