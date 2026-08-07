---
type: concept-piece
id: PZL-072
status: captured
domain: autonomy-safety
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-066]
overlaps_with: [PZL-048, PZL-067, PZL-070, PZL-084, PZL-094, PZL-096]
conflicts_with: []
tags: [stop, retry, scope, irreversible, git]
updated: 2026-08-07
---

# Explicit autonomous stop conditions

## Core message

Predeclare stops for unresolved product decisions, scope expansion, forbidden or
irreversible operations, Git divergence, failed required checks, and retry exhaustion.

## Value

Turns autonomy from vague permission into bounded delegation with observable exits.

## Assumptions

- Each condition can be detected reliably enough to fail closed.
- Ordinary in-scope work need not request repeated confirmation.

## Evidence and collisions

Current AFK, migration, policy, Git, and lifecycle contracts already encode similar
stops. The remote-deletion regression shows declarative rules need enforced tests.
Source 007 adds scope-delta proposals and a retry rule requiring new hypotheses or evidence.

## Trade-offs and recommendation

Safety and flow versus false blocks and detection complexity. Treat stop semantics
as protected behavior independent of the final mission artifact shape.

## Decision

Pending.
