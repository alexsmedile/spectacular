---
type: concept-piece
id: PZL-037
status: captured
domain: evidence-method
sources: [source-004]
source_authority: proposal
assessment: strong
evidence_status: supported
disposition: pending
depends_on: []
overlaps_with: [PZL-026, PZL-027]
conflicts_with: []
tags: [usage, records, files, telemetry, retention]
updated: 2026-08-07
---

# Normalize self-hosting usage evidence

## Core message

Use the self-hosting workspace as one evidence source, but count records, files,
indexes, textual occurrences, live state, and archived history separately.

## Value

Turns intuitive “we never use this” claims into reproducible subsystem evidence.

## Assumptions

- Self-hosting use is relevant to the target user profile.
- Each artifact class has a meaningful denominator and observation window.

## Evidence and collisions

Source counts mix units: sessions index as one use, SPC token occurrences as 1001
entities, DEC as unused despite 27 decisions, and memories index plus entries.
Unnormalized counts create false precision.

## Trade-offs and recommendation

Better decisions versus measurement work and incomplete external-user visibility.
Adopt the method as a workflow candidate; never use a single count as the survival rule.

## Decision

Pending.
