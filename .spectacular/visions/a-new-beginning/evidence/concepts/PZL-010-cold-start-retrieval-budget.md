---
type: concept-piece
id: PZL-010
status: captured
domain: measurement
sources: [source-001, source-003]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-001, PZL-002, PZL-004, PZL-006]
conflicts_with: []
tags: [benchmark, context, tokens, correctness, cold-start]
updated: 2026-08-07
---

# Cold-start retrieval budget

## Core message

Measure how much generic Spectacular context an agent loads before its first
code-relevant action on representative tasks, and treat that cost as a product
budget.

## Value

Turns “lean” from an adjective into a regression-testable property and evaluates
the combined skill, CLI, and reference retrieval design.

## Assumptions

- Representative tasks and the “first real action” boundary can be defined.
- Token or line cost is measured alongside outcome quality.

## Evidence and collisions

The repo already has `tests/benchmarks/retrieval-baseline.sh`, which measures
reference count, CLI calls, output bytes, full-body reads, repeated reads, and
whether a next action is exposed. It is content-free and does not yet measure a
real agent's pre-action context, correctness, safety, or recovery tool calls.

## Trade-offs and recommendation

Clear budget and regression signal versus benchmark gaming and harness variance.
Extend the existing baseline with realistic task scenarios and paired quality
criteria. Strong evaluation goal; the proposed under-150-line target remains a
hypothesis until the baseline is measured.

## Decision

Pending.
