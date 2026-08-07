---
type: concept-piece
id: PZL-171
status: captured
domain: contract-verification
sources: [source-014]
source_authority: user-provided-unsourced-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-073, PZL-083, PZL-095]
overlaps_with: [PZL-094, PZL-110, PZL-119]
conflicts_with: []
tags: [schema-first, compatibility, mocks, breaking-change, deterministic-gate]
updated: 2026-08-07
---

# Executable interface compatibility

## Core message

For shared public interfaces, derive mocks and compatibility checks from one accepted schema and gate
unannounced breaking changes mechanically where the schema can express them.

## Value

Lets consumers work in parallel while turning part of the compatibility promise into repeatable
evidence rather than review-time memory.

## Assumptions

- The interface is stable and consequential enough to justify schema ownership and tooling.
- Semantic, behavioral, security, and operational compatibility still receive appropriate review.

## Evidence and collisions

Source 014 names OpenAPI, Protobuf, GraphQL, generated mocks, and diff tooling but supplies no project
evidence. Generated mocks drift if they do not share schema authority, and syntactic checks cannot
detect every behavioral break or prove an implementation satisfies the contract.

## Trade-offs and recommendation

Strong deterministic guard versus schema/tooling cost and false confidence. Promote only established
shared interfaces; keep one schema authority, generate projections, and combine mechanical diffs
with clause-mapped evidence for changes the schema cannot prove.

## Decision

Pending.
