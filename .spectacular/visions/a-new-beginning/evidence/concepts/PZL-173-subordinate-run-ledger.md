---
type: concept-piece
id: PZL-173
status: captured
domain: continuity
sources: [source-007, source-015]
source_authority: synthesized-proposal-plus-primary-comparative-study
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-099, PZL-118]
overlaps_with: [PZL-075, PZL-119, PZL-126, PZL-140]
conflicts_with: []
tags: [run, ledger, resume, canonical-state, scratch, continuity]
updated: 2026-08-08
---

# Run-local ledger subordinate to Mission state

## Core message

An execution attempt may keep a compact recovery ledger for completed steps, checkpoints, repair
rounds, refs, and the next action, but it must identify and reference the owning Mission rather than
shadowing its goal, lifecycle, contract delta, or accepted evidence.

## Value

Enables cold resumption without turning attempt scratch into a second project-management authority.

## Assumptions

- Runs need recovery state often enough to justify a compact ledger.
- Canonical Mission state and disposable/runtime-owned scratch have an explicit join and retention
  rule.

## Evidence and collisions

Superpowers implements a plan-scoped ledger to prevent redispatch after context loss. GSD persists
several phase and continuation artifacts, showing both the recovery value and the taxonomy risk.
Current Spectacular AFK records and Source 007 already separate attempts from durable intent. The
study does not prove that every Mission needs a separate file or run identity.

## Trade-offs and recommendation

Reliable recovery versus duplicate status and artifact growth. Keep the conceptual ledger in S06;
start with the smallest append-only or runtime-local representation, derive shared status, and
promote only evidence that the Mission contract requires.

## Decision

Pending.
