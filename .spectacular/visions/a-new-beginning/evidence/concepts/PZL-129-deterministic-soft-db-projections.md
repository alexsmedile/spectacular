---
type: concept-piece
id: PZL-129
status: captured
domain: context-retrieval
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-002, PZL-055, PZL-123]
overlaps_with: [PZL-053, PZL-070]
conflicts_with: []
tags: [soft-db, projections, json, briefing]
updated: 2026-08-07
---

# Deterministic decision-capable soft-DB projections

## Core message

Expose compact deterministic projections that contain enough state for the next decision, with
links or identifiers for drilling into authoritative Markdown only when needed.

## Value

Reduces repeated file scans without making an LLM-generated summary authoritative.

## Assumptions

The CLI can derive projections from canonical frontmatter and stable body signals.

## Evidence and collisions

Issue #12 extends the already-proven `status --json` pattern. It overlaps read-view consolidation
and must not duplicate the measurement or briefing schema proposed by #11.

## Trade-offs and recommendation

A projection can drift from its source if separately stored. Generate it on demand or validate
cached output, and define the decisions each view must support before adding fields.

## Decision

Pending.
