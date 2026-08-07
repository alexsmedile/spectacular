---
type: concept-piece
id: PZL-153
status: captured
domain: discovery-method
sources: [source-012]
source_authority: user-provided-unsourced-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-095, PZL-112]
overlaps_with: [PZL-110]
conflicts_with: []
tags: [harness, state-machine, cli, edge-cases, prototype]
updated: 2026-08-07
---

# Executable logic decision harness

## Core message

Use a disposable executable simulator or CLI when state transitions, recovery, races, permissions,
or transformations are easier to decide by exercising scenarios than by reading prose.

## Value

Makes backend and workflow contradictions observable before production interfaces and storage harden.

## Assumptions

The harness models the decision-critical semantics faithfully enough and remains isolated from
production data, secrets, and irreversible effects.

## Evidence and collisions

This is a discovery artifact, unlike PZL-110's production deterministic harness. A toy simulator
can create false confidence if it omits the concurrency, persistence, or provider behavior at issue.

## Trade-offs and recommendation

Define the hypotheses and adversarial scenarios first. Keep the harness minimal, preserve its
results, and discard or deliberately promote code rather than letting it leak into production.

## Decision

Pending.
