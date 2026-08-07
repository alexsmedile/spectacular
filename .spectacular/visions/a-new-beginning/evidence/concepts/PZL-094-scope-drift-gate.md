---
type: concept-piece
id: PZL-094
status: captured
domain: execution-gate
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-070, PZL-088]
overlaps_with: [PZL-072, PZL-092]
conflicts_with: []
tags: [scope, paths, dependencies, interfaces, delta]
updated: 2026-08-07
---

# Contract and repository scope gate

## Core message

Pause or propose a Mission delta when execution crosses allowed paths, adds dependencies,
refactors unrelated code, changes unexpected public interfaces, or alters undeclared contracts.

## Value

Turns scope boundaries into runtime detection rather than retrospective review.

## Assumptions

- Allowed paths can be predicted without blocking legitimate code discovery.
- Contract-impact detection is reliable enough to avoid false confidence.

## Evidence and collisions

Minimal-diff guidance and current traffic/AFK scope support the goal. Exact path fences
can conflict with outcome-oriented plans when implementation location is discovered late.

## Trade-offs and recommendation

Strong autonomy safety versus brittle manifests. Use contract boundaries as authority
and paths as expected scope with an explicit, reviewable expansion mechanism.

## Decision

Pending.
