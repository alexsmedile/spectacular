---
type: concept-piece
id: PZL-074
status: captured
domain: closure-lifecycle
sources: [source-006, source-007, source-015]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: [PZL-073]
overlaps_with: [PZL-019, PZL-066, PZL-073, PZL-080, PZL-088]
conflicts_with: []
tags: [closure, reconciliation, snapshot, archive, truth]
updated: 2026-08-08
---

# Closure reconciles current truth

## Core message

Closure must compare implementation and evidence with the approved delta, refuse
missing proof, update affected contracts with approval, preserve prior truth, and archive.

## Value

Prevents implemented behavior, documented intent, and current system truth from drifting.

## Assumptions

- Contract updates can be derived accurately from verified change.
- Snapshot/archive ownership is clear.

## Evidence and collisions

Archive/spec-sync and closure gates already implement this protected outcome in a
different shape. The source strengthens it rather than proving a new mission schema.
Source 007 formalizes reconciliation over graph baseline, delta, implementation, and
evidence, but must retain code/runtime authority and drift detection.
Source 015 identifies the missing executable edge between evidence and lifecycle, while its studied
targets do not themselves demonstrate contract reconciliation.

## Trade-offs and recommendation

Trustworthy truth maintenance versus closure ceremony. Keep the outcome load-bearing;
evaluate whether current or mission-based artifacts express it more simply.

## Decision

Pending.
