---
type: concept-piece
id: PZL-056
status: captured
domain: documentation-contract
sources: [source-005]
source_authority: code-audit-proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-055]
overlaps_with: [PZL-002, PZL-005]
conflicts_with: []
tags: [docs, help, generation, drift, contract]
updated: 2026-08-07
---

# Generated command documentation

## Core message

Generate executable help, public command reference, and contract tests from one
accepted command surface rather than maintaining contradictory descriptions.

## Value

Makes public documentation drift detectable and usually impossible by construction.

## Assumptions

- Generated sections can coexist with useful human explanations.
- The command registry becomes the genuine authority.

## Evidence and collisions

Public docs describe non-CLI lifecycle operations that exist as subcommands, and
the system spec preserves docs commands that an approved decision and dispatcher
show as removed.

## Trade-offs and recommendation

Consistency versus generator maintenance and less hand-shaped prose. Generate
only contract tables/usage; keep rationale in a clearly separate human layer.

## Decision

Pending.
