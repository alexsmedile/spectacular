---
type: concept-piece
id: PZL-077
status: captured
domain: workspace-shape
sources: [source-006, source-007]
source_authority: proposal
assessment: mixed
evidence_status: partial
disposition: pending
depends_on: [PZL-063, PZL-064, PZL-066, PZL-073]
overlaps_with: [PZL-006, PZL-007, PZL-038, PZL-101]
conflicts_with: [PZL-006, PZL-101]
tags: [workspace, minimal, capabilities, missions, archive]
updated: 2026-08-07
---

# Minimal mission workspace

## Core message

Reduce the canonical tree to PROJECT, SYSTEM, capability contracts, mission bundles
containing MISSION/RUN/EVIDENCE, and archive.

## Value

Makes the entire product model visible in one small directory tree.

## Assumptions

- Decisions, policy, migrations, indexes, and other durable knowledge fit inside
  these artifacts or Git history.
- Current users can migrate without losing retrieval or safety.

## Evidence and collisions

The shape is conceptually coherent but replaces multiple proven collections and
contracts. Source 001 proposes a different three-item floor centered on PRD/requests.
Source 007 expands it into typed system directories, decisions, learnings, DELTA,
and runs, so even the Mission-centered sources disagree on the minimum tree.

## Trade-offs and recommendation

Radical clarity versus migration risk and overloaded files. Evaluate protected
outcomes first; treat the tree as one candidate projection, not the starting decision.

## Decision

Pending.
