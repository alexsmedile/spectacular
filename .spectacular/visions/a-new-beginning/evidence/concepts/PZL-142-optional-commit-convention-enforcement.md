---
type: concept-piece
id: PZL-142
status: captured
domain: implementation-quality
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-110, PZL-141]
overlaps_with: [PZL-054]
conflicts_with: []
tags: [commitlint, ci, conventional-commits, pack, enforcement]
updated: 2026-08-07
---

# Optional commit-convention enforcement

## Core message

Repositories that choose a commit convention may enforce it deterministically at their selected
boundary, without making commit grammar a universal Spectacular lifecycle.

## Value

Provides machine-checkable history where releases or collaboration actually depend on it.

## Assumptions

The repository decides whether individual commits, PR titles, or merge commits are the contract.

## Evidence and collisions

Issue #36 proposes commitlint and CI but leaves the enforcement boundary open. Native CI and Git
hooks should perform the check; Spectacular may supply an optional pack.

## Trade-offs and recommendation

Do not scaffold dependencies by default. Decide the release/collaboration need first, then provide
a narrow preset with native tooling and clear escape behavior.

## Decision

Pending.
