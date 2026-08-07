---
type: concept-piece
id: PZL-045
status: captured
domain: wayfinding
sources: [source-004, source-005]
source_authority: proposal
assessment: disputed
evidence_status: supported
disposition: pending
depends_on: [PZL-033, PZL-037]
overlaps_with: [PZL-033, PZL-039, PZL-047]
conflicts_with: [PZL-061]
tags: [wayfinding, status, scheduler, deletion]
updated: 2026-08-07
---

# Retire wayfinding

## Core message

Delete the fog/frontier dependency sequencer and eight wayfinding verbs because a
solo project's next-work queue is small enough for status and human judgment.

## Value

Would remove a complex scheduling/discovery surface from the skill and CLI.

## Assumptions

- Request priority and spec/discovery dependency readiness are the same question.
- Strict dependency ordering has not prevented meaningful mistakes.

## Evidence and collisions

A live comparison refutes equivalence: wayfinding chose SPC-007 while status chose
stance-layer. Source 005 explicitly recommends keeping dependency-aware Wayfinder
while removing top-level `next`. A reproduced kind/type ranking defect means the
retention test must be rerun after stabilization. The behavior is unique; whether
it is valuable enough remains open.

## Trade-offs and recommendation

Much smaller surface versus loss of dependency/fog semantics. Evaluate cases where
the different answer changed action quality before retaining or removing it.
Disputed.

## Decision

Pending.
