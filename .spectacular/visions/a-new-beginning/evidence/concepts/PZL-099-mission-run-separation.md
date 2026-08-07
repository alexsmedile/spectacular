---
type: concept-piece
id: PZL-099
status: captured
domain: execution-model
sources: [source-007]
source_authority: synthesized-proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-066]
overlaps_with: [PZL-071, PZL-075, PZL-088]
conflicts_with: []
tags: [mission, run, attempt, failure, resume, lifecycle]
updated: 2026-08-07
---

# Mission/run attempt separation

## Core message

Keep durable Mission intent separate from one or more execution attempts so failed,
cancelled, discovery, repair, verification, and observation runs do not corrupt the Mission.

## Value

Preserves coherent intent and supports evidence/Git-based resume after failure.

## Assumptions

- Multiple attempts occur often enough to justify identity and lifecycle.
- Run records can stay compact and archive cleanly.

## Evidence and collisions

AFK run records and debug traces already distinguish attempts from durable work; Source
006's single RUN file is a simpler competing shape.

## Trade-offs and recommendation

Honest attempt history versus a second state machine. Adopt conceptual separation;
start with append-only attempts inside one run log before assigning IDs/directories.

## Decision

Pending.
