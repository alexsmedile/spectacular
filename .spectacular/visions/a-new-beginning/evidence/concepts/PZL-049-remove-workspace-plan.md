---
type: concept-piece
id: PZL-049
status: captured
domain: cli-surface
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-053, PZL-058]
conflicts_with: []
tags: [workspace, duplicate, command, removal]
updated: 2026-08-07
---

# Remove duplicate workspace plan

## Core message

Retire `workspace plan` because it dispatches to the same implementation as
`workspace preflight`; retain one clearly named read projection.

## Value

Removes a literal duplicate and prevents two names from implying distinct contracts.

## Assumptions

- No external automation depends on the alias without a migration path.
- The retained name accurately describes the output.

## Evidence and collisions

The dispatcher shares one path for both commands. Documentation may still promise
different semantics, so public usage and help need inventory before removal.

## Trade-offs and recommendation

Small compatibility cost for clear surface reduction. Strong candidate for a
deprecation/removal batch, separate from broader workspace redesign.

## Decision

Pending.
