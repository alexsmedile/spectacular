---
type: concept-piece
id: PZL-091
status: captured
domain: context-compilation
sources: [source-007]
source_authority: synthesized-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-070, PZL-088, PZL-090]
overlaps_with: [PZL-055, PZL-070]
conflicts_with: []
tags: [compiler, run-manifest, context, permissions, budgets]
updated: 2026-08-07
---

# Mission compiler and run manifest

## Core message

Compile durable Mission intent into a bounded manifest containing only relevant
contracts, exact repository facts, scope, tools, secrets rules, probes, Git, budgets, and gates.

## Value

Separates long-lived intent from attempt-specific executable context.

## Assumptions

- Compilation is reproducible and validates missing inputs.
- Generated manifests do not become stale shadow authority.

## Evidence and collisions

Current briefs and AFK context provide precedents; graph-based selection and exact
manifest schema remain unproven.

## Trade-offs and recommendation

Small safe prompts versus compiler complexity. Start as a generated projection with
source references and expiry, not a new hand-maintained truth document.

## Decision

Pending.
