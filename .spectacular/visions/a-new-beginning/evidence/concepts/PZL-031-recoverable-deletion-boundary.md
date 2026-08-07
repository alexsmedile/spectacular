---
type: concept-piece
id: PZL-031
status: captured
domain: git-safety
sources: [source-003, source-004]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-023]
conflicts_with: []
tags: [git, deletion, recovery, branch, tag]
updated: 2026-08-07
---

# Recoverable deletion boundary

## Core message

Preserve a verified Git recovery boundary before large deletions so removal does
not carry the emotional or operational cost of irreversible destruction.

## Value

Encourages honest simplification while retaining exact restoration evidence.

## Assumptions

- The recovery ref is named, durable, and verified before deletion.
- Removed code does not remain in live retrieval paths.

## Evidence and collisions

The current work already uses an isolated refactor branch, and Spectacular has an
archive-first Git recovery principle. A tag/ref is cleaner than moving dead code
into `contrib/`, which keeps it in the working tree and agent search surface.

## Trade-offs and recommendation

Safe experimentation versus accumulating stale recovery refs. Strong rule; choose
a retention policy separately and avoid live-tree dead-code shelves.

## Decision

Pending.
