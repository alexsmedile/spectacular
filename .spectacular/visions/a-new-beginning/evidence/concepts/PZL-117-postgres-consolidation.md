---
type: concept-piece
id: PZL-117
status: captured
domain: infrastructure-option
sources: [source-008]
source_authority: unsourced-protocol-synthesis
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: []
overlaps_with: [PZL-106]
conflicts_with: [PZL-038, PZL-069, PZL-077]
tags: [postgresql, jsonb, vector, queue, property-graph]
updated: 2026-08-07
---

# PostgreSQL consolidation option

## Core message

For a server-backed system already operating PostgreSQL, evaluate whether relational data,
JSONB, vector search, and queue-like coordination can share one operational substrate.

## Value

May reduce the number of deployed state services and transactional boundaries.

## Assumptions

- The workload, scale, availability, and operations fit PostgreSQL.
- Extensions and version-specific features are acceptable dependencies.

## Evidence and collisions

The universal claim is unsupported. JSONB is core; `SKIP LOCKED` supports queue-like access
but yields an inconsistent view; pgvector is an extension; property-graph syntax appears in
PostgreSQL 19 documentation but not current stable 18. Spectacular's portable core is
Markdown/Git-native and does not yet require a server database.

## Trade-offs and recommendation

Consolidation simplifies operations until workloads, extensions, or failure domains diverge.
Keep this as a stack-specific architecture option, not a Spectacular refactor requirement.

## Decision

Pending.
