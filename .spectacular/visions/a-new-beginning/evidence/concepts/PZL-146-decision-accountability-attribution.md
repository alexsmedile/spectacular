---
type: concept-piece
id: PZL-146
status: captured
domain: decision-governance
sources: [source-011, source-014]
source_authority: owner-authored-proposal-corpus
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-086, PZL-115]
overlaps_with: [PZL-119]
conflicts_with: []
tags: [decision, decider, provenance, authority, accountability]
updated: 2026-08-07
---

# Decision accountability attribution

## Core message

Record who authorized a decision separately from which human or agent authored, recommended, or
recorded it.

## Value

Clarifies accountability and prevents agent provenance from being mistaken for approval authority.

## Assumptions

Identity fields have stable enough meanings to be useful across runtimes.

## Evidence and collisions

Issue #40 proposes decider attribution. The important collision is semantic: `author`, `proposed_by`,
`recorded_by`, and `decided_by` are not interchangeable.

## Trade-offs and recommendation

Define the authority semantics before field names. Require `decided_by` only where authorization is
material; keep execution provenance separately generated.

## Decision

Pending.
