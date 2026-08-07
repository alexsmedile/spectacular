---
type: concept-piece
id: PZL-095
status: captured
domain: verification-method
sources: [source-007, source-014]
source_authority: synthesized-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-073]
overlaps_with: [PZL-092, PZL-093]
conflicts_with: []
tags: [evidence, tdd, tests, migrations, ui, operations]
updated: 2026-08-07
---

# Evidence-first strategy by change class

## Core message

Define proof before implementation, then select reproduction, contract tests,
regression tests, migration checks, UI evidence, operational observation, or validators by change type.

## Value

Preserves test-first intent without forcing artificial unit tests onto every artifact.

## Assumptions

- Change classification guides a suitable minimum proof set.
- Evidence plans can adapt when implementation discoveries invalidate them.

## Evidence and collisions

The attachments prescribe tests first; Source 007 improves that generic rule with
change-specific evidence. Current verification already varies by complexity.

## Trade-offs and recommendation

Senior judgment and relevance versus more routing rules. Adopt the principle and a
small example matrix; avoid a universal test-type bureaucracy.

## Decision

Pending.
