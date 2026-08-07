---
type: concept-piece
id: PZL-039
status: captured
domain: discovery-model
sources: [source-004, source-006]
source_authority: proposal
assessment: disputed
evidence_status: partial
disposition: pending
depends_on: [PZL-037]
overlaps_with: [PZL-043, PZL-065, PZL-076]
conflicts_with: []
tags: [question, research, spike, prototype, discovery]
updated: 2026-08-07
---

# Collapse discovery types

## Core message

Replace QUE, RES, SPK, and prototype identities with one question record carrying
a `kind:` field.

## Value

Reduces routing, schemas, commands, IDs, and user choice at capture time.

## Assumptions

- The types differ mainly in label rather than execution authority and evidence contract.
- One lifecycle can safely express human answers, sourced facts, and throwaway experiments.

## Evidence and collisions

No live question/research/spike records currently appear, but historical usage was
not loaded. Current rules intentionally distinguish user judgment, read-only
research, and human-authorized feasibility code. Prototype is artifact-owned, not
an equivalent record type. Source 006 proposes a different collapse: research is
a mission objective, bugs are violations, fixes are repair missions, and session
history is the run record. That increases the simplification case but does not
prove the behaviors equivalent.

## Trade-offs and recommendation

Simpler capture versus erasing safety and evidence semantics. Compare behavior and
required gates, not counts or names. Disputed as a blanket collapse.

## Decision

Pending.
