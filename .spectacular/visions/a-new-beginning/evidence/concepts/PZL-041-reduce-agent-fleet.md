---
type: concept-piece
id: PZL-041
status: captured
domain: agent-architecture
sources: [source-004]
source_authority: proposal
assessment: disputed
evidence_status: unverified
disposition: pending
depends_on: [PZL-037]
overlaps_with: [PZL-001]
conflicts_with: []
tags: [agents, roles, delegation, maintenance]
updated: 2026-08-07
---

# Reduce the agent fleet

## Core message

Reduce nine specialist agent definitions to roughly five by removing rarely used
or overlapping research, spec-review, audit, and debug roles.

## Value

Shrinks routing surface and role maintenance while keeping the highest-value
discover/apply/review capabilities.

## Assumptions

- Definition count represents runtime or cognitive cost.
- Proposed removals lack unique safety or quality value.

## Evidence and collisions

Nine agent files exist, matching the source. No dispatch-frequency, outcome, or
context-cost evidence is provided. Agents load on demand, so file count alone does
not establish prompt cost.

## Trade-offs and recommendation

Simpler delegation versus collapsing valuable separation of investigation, mutation,
review, and verification. Inventory actual dispatches and unique contracts first.
Disputed.

## Decision

Pending.
