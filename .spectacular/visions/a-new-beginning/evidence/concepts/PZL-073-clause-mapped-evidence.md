---
type: concept-piece
id: PZL-073
status: captured
domain: verification-contract
sources: [source-006, source-007]
source_authority: proposal
assessment: strong
evidence_status: partial
disposition: pending
depends_on: [PZL-064, PZL-066, PZL-068]
overlaps_with: [PZL-019, PZL-074, PZL-083, PZL-089, PZL-095]
conflicts_with: []
tags: [evidence, verification, contract, checks]
updated: 2026-08-07
---

# Evidence mapped to contract clauses

## Core message

Record commits/diffs, executed checks, exit status, demonstrated scenarios,
limitations, delivery artifact, and which exact contract clauses are satisfied.

## Value

Makes verification traceable to promised behavior rather than a generic green suite.

## Assumptions

- Contract clauses are precise and stable enough to reference.
- Human review can initially supply independent judgment where automation cannot.

## Evidence and collisions

VERIFY, VERIFY-LOG, SPEC-DELTA, and PR gates hold similar evidence in multiple
places. A single EVIDENCE file may consolidate them or become another copy.
Source 007 makes verification probes part of contracts and evidence strategy depend
on change class, improving traceability while increasing schema coupling.

## Trade-offs and recommendation

Strong traceability versus clerical mapping. Define one evidence authority and
derive projections rather than duplicating check results.

## Decision

Pending.
