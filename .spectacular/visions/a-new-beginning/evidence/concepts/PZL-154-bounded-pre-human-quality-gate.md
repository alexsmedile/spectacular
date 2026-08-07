---
type: concept-piece
id: PZL-154
status: captured
domain: verification-architecture
sources: [source-012, source-013]
source_authority: user-provided-unsourced-proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-110, PZL-114, PZL-127]
overlaps_with: [PZL-104, PZL-115]
conflicts_with: [PZL-079]
tags: [pre-review, evaluator, critic, quality-gate, human-attention]
updated: 2026-08-07
---

# Bounded pre-human quality gate

## Core message

Before asking for human judgment, remove mechanically detectable breakage and material rubric
violations through deterministic checks plus bounded independent review.

## Value

Protects human attention for product and trade-off decisions rather than obvious defects.

## Assumptions

The gate distinguishes defects from product choices and has explicit pass, accepted-risk, blocked,
and failed terminal states.

## Evidence and collisions

Source 012 proposes an evaluator loop that auto-corrects until a high bar. PZL-127 already shows why
subjective satisfaction and unbounded retries are unsafe; the pattern also exceeds the MVP fence if
made a universal multi-agent requirement.

## Trade-offs and recommendation

Use it only when inspection cost or consequence justifies the latency. Cap attempts, retain exact
findings, and return unresolved product questions to the human rather than letting critics decide.

## Decision

Pending.
