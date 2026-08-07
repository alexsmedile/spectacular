---
type: concept-piece
id: PZL-002
status: captured
domain: cli-contract
sources: [source-001, source-003, source-005]
source_authority: proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: []
overlaps_with: [PZL-001, PZL-003]
conflicts_with: []
tags: [cli, help, authority, context]
updated: 2026-08-07
---

# CLI command-help authority

## Core message

Make `spectacular help <verb>` the deterministic authority for command syntax,
options, projections, and safe mechanical next steps so model instructions do not
repeat them.

## Value

One executable source can serve humans, agents, tests, and generated docs while
removing syntax-only rows from universal prompt context.

## Assumptions

- Per-command help can expose enough semantics without becoming another giant
  output surface.
- The CLI can distinguish mechanical facts from agentic judgment boundaries.

## Evidence and collisions

Current `spectacular help status` exits successfully but prints only global help.
The authority does not exist yet. Source 005 independently found public command
docs and executable dispatch drifting apart, strengthening the case for one
generated contract. Some router rows still carry lifecycle or safety meaning, so
equivalence needs a field-by-field inventory.

## Trade-offs and recommendation

Deterministic discoverability versus more CLI surface and a possible second
documentation generator. Design and test per-command help before deleting skill
content. Promising; needs contract design.

## Decision

Pending.
