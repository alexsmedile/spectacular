---
type: concept-piece
id: PZL-028
status: captured
domain: identity-taxonomy
sources: [source-003, source-004]
source_authority: proposal
assessment: promising
evidence_status: supported
disposition: pending
depends_on: [PZL-027]
overlaps_with: [PZL-009]
conflicts_with: []
tags: [ids, reserved, taxonomy, compatibility]
updated: 2026-08-07
---

# Reserved identity pruning

## Core message

Remove reserved entity prefixes and init flags that advertise workflows the
product does not implement or use.

## Value

Shrinks conceptual surface and prevents placeholder taxonomy from looking like a
supported capability.

## Assumptions

- Reservation provides less compatibility value than its ongoing cognitive cost.
- No external workspace has allocated or depended on the reserved prefixes.

## Evidence and collisions

TSK, FND, FIX, BUG, SEC, BMK, and PRT are reserved in canonical IDs. Five live
fixes use legacy F1–F5; FIX is a successor reservation, not the active identity.
Removing the reservation and removing the fix workflow are separate decisions.

## Trade-offs and recommendation

Honest supported surface versus future migration freedom and external compatibility.
Evaluate each prefix independently; promising cleanup, not a bulk delete.

## Decision

Pending.
