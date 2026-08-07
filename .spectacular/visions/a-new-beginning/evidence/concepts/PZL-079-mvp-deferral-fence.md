---
type: concept-piece
id: PZL-079
status: captured
domain: scope-governance
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: not-needed
disposition: pending
depends_on: [PZL-062]
overlaps_with: [PZL-023, PZL-038, PZL-082, PZL-090, PZL-098]
conflicts_with: []
tags: [mvp, deferral, parallelism, policy, packs, databases]
updated: 2026-08-07
---

# Explicit MVP deferral fence

## Core message

Defer parallel orchestration, schedulers, semantic retrieval, extra registries,
custom policy language, nested workspaces, roadmaps, general document engines,
packs, autonomous release, dashboards, extra databases, cross-project memory, and
complex migration until the single-agent loop succeeds.

## Value

Prevents platform breadth from masking weakness in the core experience.

## Assumptions

- Each deferred capability is unnecessary for the chosen MVP user and test.
- Existing shipped behavior can be frozen, hidden, or removed safely while evaluating the loop.

## Evidence and collisions

Several deferred items are already implemented and governed by decisions, so
“defer” must distinguish no new investment, default-off, deprecation, and deletion.
Source 007 keeps semantic search secondary and still rejects generalized scheduling,
but introduces six contract types and a possible multi-agent run chain inside the MVP.

## Trade-offs and recommendation

Strong focus versus abandoning current-user value. Use as an investment fence
immediately; make retention/removal decisions per subsystem with evidence.

## Decision

Pending.
