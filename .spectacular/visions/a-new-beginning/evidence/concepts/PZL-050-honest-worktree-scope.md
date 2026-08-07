---
type: concept-piece
id: PZL-050
status: captured
domain: worktree-coordination
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-044, PZL-054]
conflicts_with: []
tags: [git, worktree, preflight, naming, scope]
updated: 2026-08-07
---

# Honest worktree inspection scope

## Core message

Either inspect linked worktrees when claiming cross-worktree coordination or name
the feature honestly as active-checkout preflight.

## Value

Aligns the safety promise with observable coverage and prevents false confidence.

## Assumptions

- Sibling worktrees are relevant to the advertised coordination outcome.
- Git's worktree inventory is reliable enough for deterministic inspection.

## Evidence and collisions

No `git worktree list` use exists in the CLI, while workspace language claims
concurrent-worktree protection. The current implementation protects one checkout.

## Trade-offs and recommendation

Real inspection adds state and tests; renaming admits a narrower capability.
Choose based on the protected coordination outcome, not the current feature name.

## Decision

Pending.
