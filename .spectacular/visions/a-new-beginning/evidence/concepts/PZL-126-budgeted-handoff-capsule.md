---
type: concept-piece
id: PZL-126
status: captured
domain: orchestration-contract
sources: [source-010]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-108, PZL-119]
overlaps_with: [PZL-105]
conflicts_with: []
tags: [handoff, token-budget, compression, state]
updated: 2026-08-07
---

# Budgeted handoff capsule

## Core message

Transfer the smallest complete state capsule—goal, baseline, paths, accepted decisions,
evidence, unresolved failures, and next action—while omitting conversational chronology.

## Value

Protects context budgets without forcing the next worker to reconstruct critical state.

## Assumptions

- Full evidence remains addressable by stable pointers.
- Completeness is tested on realistic resume and review tasks.

## Evidence and collisions

Source 010's 300-token limit has no supplied basis and cannot fit every task. Its omission
rule is valuable: failed trial chronology belongs in evidence unless it changes the next action.

## Trade-offs and recommendation

Aggressive compression reduces noise but can erase constraints and negative evidence. Use a
measured soft budget with mandatory fields and overflow references, not a universal token cap.

## Decision

Pending.
