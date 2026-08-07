---
type: concept-piece
id: PZL-083
status: captured
domain: contract-schema
sources: [source-007, source-014]
source_authority: synthesized-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-082]
overlaps_with: [PZL-064, PZL-086]
conflicts_with: []
tags: [guarantees, invariants, failures, freshness, provenance]
updated: 2026-08-07
---

# Shared contract envelope

## Core message

Every contract should expose guarantees, exclusions, inputs/outputs, invariants,
failures, dependencies, implementation pointers, probes, freshness, and provenance.

## Value

Makes authority, limits, proof, and staleness inspectable across contract types.

## Assumptions

- Fields apply meaningfully across all six types.
- Empty fields do not become template ceremony.

## Evidence and collisions

Current schemas already carry parts of the envelope, but no audit shows which fields
agents actually use or keep current.

## Trade-offs and recommendation

Consistent retrieval versus repetitive boilerplate. Define a minimal required core
and type-specific optional fields after testing real contracts.

## Decision

Pending.
