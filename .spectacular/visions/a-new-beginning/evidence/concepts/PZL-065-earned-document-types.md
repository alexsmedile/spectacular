---
type: concept-piece
id: PZL-065
status: captured
domain: information-architecture
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-064]
overlaps_with: [PZL-003, PZL-004, PZL-039, PZL-082, PZL-101]
conflicts_with: []
tags: [progressive-disclosure, document-types, complexity]
updated: 2026-08-07
---

# Promote document types only when earned

## Core message

Keep component, interface, state, API, and operations concerns as sections inside
a capability until repeated complexity proves a separate canonical type is useful.

## Value

Prevents speculative taxonomies and retrieval hops at the start of a project.

## Assumptions

- Section growth reveals the right promotion boundary.
- Cross-capability reuse can be handled before a shared registry exists.

## Evidence and collisions

Current thin rules and sparse collections support the concern, but some contracts
are intentionally flat/shared. No promotion threshold is defined.
Source 007 proposes six stable semantic types while retaining embedded sections for
small projects, separating type identity from filesystem promotion.

## Trade-offs and recommendation

Low initial ceremony versus duplication inside capability files. Strong default;
define evidence-based promotion criteria rather than forbidding shared contracts.

## Decision

Pending.
