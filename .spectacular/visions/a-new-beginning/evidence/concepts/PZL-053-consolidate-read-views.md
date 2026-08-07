---
type: concept-piece
id: PZL-053
status: captured
domain: cli-surface
sources: [source-005]
source_authority: code-audit-proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-051]
overlaps_with: [PZL-002, PZL-049, PZL-058]
conflicts_with: []
tags: [status, request, projections, consolidation]
updated: 2026-08-07
---

# Consolidate request read views

## Core message

Group fleet-level read projections under `status` and request-specific projections
under `request`, removing top-level synonyms over the same state.

## Value

Creates two obvious discovery points and reduces overlapping nouns and verbs.

## Assumptions

- Each current view can be mapped without losing a distinct user question.
- Namespace depth is cheaper than top-level discoverability.

## Evidence and collisions

Status, summary, next, requests, request, progress, links, and traffic overlap, but
no field-level equivalence matrix or usage data exists yet.

## Trade-offs and recommendation

Smaller grammar versus hiding valuable specialized projections. Inventory inputs,
fields, authority, and consumers before deciding each fold independently.

## Decision

Pending.
