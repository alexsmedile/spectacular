---
type: concept-piece
id: PZL-075
status: captured
domain: continuity
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-070]
overlaps_with: [PZL-010, PZL-071, PZL-078, PZL-096, PZL-099, PZL-102]
conflicts_with: []
tags: [status, resume, checkpoint, context, chat-independent]
updated: 2026-08-07
---

# Cold-agent status and resume

## Core message

A compact briefing and continue operation should reconstruct current state, last
proof, next objective, blockers, pending evidence, Git baseline, and bounded context.

## Value

Directly tests continuity without relying on the original chat or agent memory.

## Assumptions

- Durable records contain every load-bearing decision and checkpoint.
- Resume detects code and Git drift before acting.

## Evidence and collisions

Current status brief, request brief, and SESSION next actions support much of this,
but their complete cold-agent resume behavior has not been tested end to end.
Source 007 requires resumption from evidence and a known Git checkpoint across failed runs.

## Trade-offs and recommendation

High user value versus pressure to over-persist execution traces. Make cold-resume
an acceptance test and add state only where failures reveal a gap.

## Decision

Pending.
