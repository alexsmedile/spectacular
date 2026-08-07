---
type: concept-piece
id: PZL-121
status: captured
domain: evaluation-method
sources: [source-009, source-012]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-010, PZL-078, PZL-110]
overlaps_with: [PZL-033]
conflicts_with: []
tags: [harness, model, benchmark, evaluation]
updated: 2026-08-07
---

# Evaluate the model-harness system

## Core message

Measure an agent as a model-plus-harness configuration; hold tasks and evidence criteria
constant before attributing reliability or cost differences to either part.

## Value

Prevents model upgrades from masking routing, tool, context, state, or verification defects.

## Assumptions

- Representative tasks are repeatable and resistant to benchmark-specific overfitting.
- Cost, latency, retrieval, retries, and correctness are recorded together.

## Evidence and collisions

Harness-Bench reports substantial outcome variation across model-harness pairings, while
Agents' Last Exam includes workloads where model variation is larger. Source 009's exact
5%/48-point slogan is therefore not a universal estimate.

## Trade-offs and recommendation

Controlled evaluation costs time but avoids architecture by anecdote. Extend the cold-start
test into three representative tasks and compare configurations before changing the fleet.

## Decision

Pending.
