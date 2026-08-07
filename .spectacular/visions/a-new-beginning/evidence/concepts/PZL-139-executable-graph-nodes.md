---
type: concept-piece
id: PZL-139
status: captured
domain: orchestration-model
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: disputed
evidence_status: unverified
disposition: pending
depends_on: [PZL-103, PZL-118, PZL-119]
overlaps_with: [PZL-104]
conflicts_with: [PZL-071, PZL-079]
tags: [dag, nodes, claiming, concurrency, reintegration]
updated: 2026-08-07
---

# Executable graph nodes

## Core message

Represent work as dependency-aware nodes with explicit claiming, eligibility, state transitions,
evidence, and reintegration contracts so independent nodes may run safely.

## Value

Could make complex, multi-session programs resumable and selectively parallel.

## Assumptions

Real workloads justify concurrency and node state beyond a single checkpointed run.

## Evidence and collisions

Issue #31 amplifies graph-orchestration sources but conflicts with the single-agent MVP and explicit
deferral fence. Its node contract also determines #20 and #24.

## Trade-offs and recommendation

Do not implement a scheduler. Jointly model one historical dependency graph, define claiming and
merge semantics, and prove that a simpler sequential objective DAG fails first.

## Decision

Pending.
