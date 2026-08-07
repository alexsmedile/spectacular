---
type: concept-piece
id: PZL-019
status: captured
domain: truth-model
sources: [source-002, source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-013, PZL-018, PZL-062, PZL-073, PZL-074, PZL-080]
conflicts_with: []
tags: [truth, code, tests, production, evidence, records]
updated: 2026-08-07
---

# Product Truth and Records

## Core message

Distinguish intended capability contracts from implementation truth and durable
records of why, evidence, outcomes, and lessons.

## Value

Prevents plans and specs from competing with code/tests while retaining historical
reasoning needed for future work.

## Assumptions

- “Product Truth” can distinguish implemented intent from observed behavior.
- Records remain retrievable without becoming live authority.

## Evidence and collisions

Current architecture already makes production code plus executable tests
implementation truth and archives detailed specs/history. Adding production
behavior to “truth” needs care: observed behavior may prove a defect, drift, or
incident rather than the intended system contract.
Source 006 makes evidence-backed closure update the capability contract, strongly
supporting reconciliation while leaving code/tests as final implementation truth.
Source 007 calls the system graph “current truth”; this must be narrowed to accepted
contract state so stale documents never outrank implementation or observed reality.

## Trade-offs and recommendation

Strong authority separation versus an overloaded truth label. Preserve at least
three layers: intended contract, implemented/tested truth, and observed runtime
reality. Strong concept; final names remain open.

## Decision

Pending.
