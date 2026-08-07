---
type: concept-piece
id: PZL-110
status: captured
domain: enforcement-model
sources: [source-008, source-009, source-010, source-012]
source_authority: unsourced-protocol-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-092, PZL-095]
overlaps_with: [PZL-072, PZL-093, PZL-094]
conflicts_with: []
tags: [harness, tests, linters, determinism, agents]
updated: 2026-08-07
---

# Deterministic agent harness

## Core message

Contain nondeterministic agent judgment inside deterministic state, scope, syntax, test,
and policy checks wherever a property is mechanically verifiable.

## Value

Turns important constraints from prompt advice into reproducible gates and evidence.

## Assumptions

- Checks target meaningful properties rather than artifact presence alone.
- Semantic and product judgment retain explicit human or review authority.

## Evidence and collisions

Spectacular already uses parsers, doctor checks, lifecycle gates, and a test suite. Source
008's causal claim is broad: deterministic checks cannot establish all correctness.

## Trade-offs and recommendation

Harnesses increase reliability but can accumulate brittle policy and false confidence.
Keep a small risk-led gate set with observable failure value and named ownership.

## Decision

Pending.
