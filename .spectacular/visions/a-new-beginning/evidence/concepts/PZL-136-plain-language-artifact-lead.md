---
type: concept-piece
id: PZL-136
status: captured
domain: intent-contract
sources: [source-011]
source_authority: owner-authored-proposal-corpus
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-119, PZL-131]
conflicts_with: []
tags: [summary, issue, pull-request, readability, handoff]
updated: 2026-08-07
---

# Plain-language artifact lead

## Core message

Lead generated issues and pull requests with one plain-language sentence describing the outcome or
change before exposing framework metadata and checklists.

## Value

Makes delivery artifacts understandable to collaborators who do not know Spectacular's taxonomy.

## Assumptions

The lead can be derived from or reviewed against the canonical goal and actual diff.

## Evidence and collisions

Issues #26 and #32 contain direct owner feedback that machine-shaped artifact openings obscure the
work. Issue #26's original pain is largely covered by the PR-focused #32.

## Trade-offs and recommendation

Keep goal, change summary, and title distinct: why, what changed, and identifier. Implement the PR
lead first; defer external issue creation ownership.

## Decision

Pending.
