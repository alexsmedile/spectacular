---
type: concept-piece
id: PZL-001
status: captured
domain: skill-routing
sources: [source-001, source-003]
source_authority: proposal
assessment: promising
evidence_status: partial
disposition: pending
depends_on: [PZL-002]
overlaps_with: [PZL-004]
conflicts_with: []
tags: [context, routing, skill]
updated: 2026-08-07
---

# Domain-routed lean skill

## Core message

Keep only model-judgment instructions in `SKILL.md`: intent classification, a
small domain router, canonical safety rules, and execution terminal conditions.

## Value

Reduces universal context and makes routing scale by user intent rather than by
enumerating every command.

## Assumptions

- The chosen domains cover real work without frequent ambiguous fallthrough.
- Deterministic command syntax and projections are available elsewhere.
- Safety and lifecycle semantics are not accidentally classified as syntax.

## Evidence and collisions

The skill is 309 lines and contains 132 table rows across all tables. A precise
inventory must separate judgment-bearing rows from replaceable CLI documentation.
Deleting first and measuring later would risk semantic loss.

## Trade-offs and recommendation

Smaller context and clearer ownership versus greater dependence on correct domain
classification and follow-up retrieval. Prototype the router against a realistic
task corpus before choosing a line target. Promising; not yet adopted.

## Decision

Pending.
