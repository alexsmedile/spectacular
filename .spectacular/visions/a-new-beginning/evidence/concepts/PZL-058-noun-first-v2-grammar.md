---
type: concept-piece
id: PZL-058
status: captured
domain: cli-contract
sources: [source-005, source-006]
source_authority: code-audit-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-038, PZL-051]
overlaps_with: [PZL-009, PZL-052, PZL-053, PZL-055, PZL-066]
conflicts_with: []
tags: [v2, cli, grammar, nouns, namespaces]
updated: 2026-08-07
---

# Noun-first v2 grammar

## Core message

Organize the mechanical CLI around stable resource nouns with explicit operations,
folding top-level verb aliases and plural list commands into namespaces.

## Value

Makes discovery systematic and reduces memorized exceptions as the surface grows.

## Assumptions

- Noun namespaces match user mental models across all record and lifecycle types.
- A breaking v2 migration is acceptable after compatibility support.

## Evidence and collisions

Current grammar mixes nouns, verbs, plurals, aliases, and terminal redirects. An
existing noun-first direction supports the pattern, but the full proposed mapping
has not been tested for ergonomics or script migration. Source 006 reinforces a
`mission <operation>` namespace while disputing the “mechanical only” constraint.

## Trade-offs and recommendation

Coherent grammar versus more keystrokes and broad compatibility churn. Decide the
canonical matrix before implementing aliases or registry machinery.

## Decision

Pending.
