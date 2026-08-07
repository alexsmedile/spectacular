---
type: concept-piece
id: PZL-054
status: captured
domain: provider-boundary
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-038, PZL-048, PZL-051]
overlaps_with: [PZL-044, PZL-050]
conflicts_with: []
tags: [git, github, gh, mutation, boundary]
updated: 2026-08-07
---

# Native provider mutation boundary

## Core message

Let Spectacular interpret ownership, blockers, readiness, authorization, and
reconciliation while native Git and `gh` perform branch and pull-request mutations.

## Value

Avoids reimplementing provider clients while preserving Spectacular's unique context.

## Assumptions

- Recipes can retain consent, idempotence, and recoverability guarantees.
- Required native tools are available where the affected workflow runs.

## Evidence and collisions

Git operations are duplicated across workspace, AFK, and GitHub paths. The current
remote-deletion regression shows the risk, but native delegation alone does not
guarantee safety or a coherent user experience.

## Trade-offs and recommendation

Less wrapper code versus distributed execution and weaker deterministic control.
Specify semantic outputs and safety invariants before choosing command ownership.

## Decision

Pending.
