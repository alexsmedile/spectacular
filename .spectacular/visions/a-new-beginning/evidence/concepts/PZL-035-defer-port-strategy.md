---
type: concept-piece
id: PZL-035
status: captured
domain: migration-strategy
sources: [source-003, source-004]
source_authority: proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-023, PZL-032]
overlaps_with: [PZL-004]
conflicts_with: []
tags: [freeze, extraction, rewrite, sequencing]
updated: 2026-08-07
---

# Defer port strategy until surface settles

## Core message

Choose freeze-and-wrap, incremental extraction, or rewrite only after the retained
product surface and contracts are known.

## Value

Avoids estimating or porting code that later reduction removes.

## Assumptions

- A port or implementation-language decision is genuinely in scope.
- The current system can remain stable until surface decisions finish.

## Evidence and collisions

Source 004 identifies the 16,507-line Bash CLI as the target and supplies the
missing session context. The three options still lack cost, portability, and
compatibility evidence.

## Trade-offs and recommendation

Better-informed cost judgment versus delaying architectural risk discovery.
Retain as a dependency rule and request the missing context before evaluating options.

## Decision

Pending.
