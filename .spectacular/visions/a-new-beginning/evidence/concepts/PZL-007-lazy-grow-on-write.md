---
type: concept-piece
id: PZL-007
status: captured
domain: initialization
sources: [source-001, source-003]
source_authority: proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-006]
overlaps_with: [PZL-006, PZL-008]
conflicts_with: []
tags: [lazy-creation, collections, init, substrate]
updated: 2026-08-07
---

# Lazy grow-on-write substrates

## Core message

Create optional collections, indexes, and supporting directories on the first
command that actually writes their first record rather than during init.

## Value

Eliminates empty capability theater and lets workspace shape reveal real use.

## Assumptions

- Readers and doctor can distinguish “unused” from “unsupported” safely.
- First-write creation is deterministic, idempotent, and migration-safe.
- Git and packaging preserve empty-versus-absent semantics where needed.

## Evidence and collisions

Several current collections can remain empty. However, canonical specs and agent
rules are not equivalent to append-only soft databases. Arbitrary usage-count
suggestions such as third request or fifth archive have no evidence.

## Trade-offs and recommendation

Honest sparse structure versus more conditional paths in every writer and doctor
check. Evaluate first for append-only optional collections; decide canonical docs
separately. Promising when applied by artifact class, not as one blanket rule.

## Decision

Pending.
