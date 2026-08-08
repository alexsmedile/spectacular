---
type: concept-piece
id: PZL-119
status: captured
domain: orchestration-contract
sources: [source-009, source-010, source-015]
source_authority: unsourced-expanded-synthesis
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-104, PZL-105]
overlaps_with: [PZL-070, PZL-091]
conflicts_with: []
tags: [handoff, schema, provenance, edges]
updated: 2026-08-08
---

# Typed inter-node handoff contract

## Core message

Pass a bounded, validated packet between planner, worker, reviewer, and merger containing
objective, authoritative inputs, exclusions, output schema, evidence, and unresolved questions.

## Value

Makes graph edges inspectable contracts rather than implicit conversation inheritance.

## Assumptions

- The schema remains small enough to use and rich enough to expose missing authority.
- Full primary evidence remains retrievable by stable reference.

## Evidence and collisions

Source 009 proposes validated JSON or Markdown state but supplies no schema or comparison.
This operationalizes PZL-105's shared envelope and PZL-091's run manifest at node boundaries.

## Trade-offs and recommendation

Typed packets reduce ambiguity but can become verbose shadow copies. Prototype one common
envelope with role-specific output fields and reject duplicated source content.

## Decision

Pending.
