---
status: planned
updated: 2026-05-30
related:
  - PLAN.md
---

# Tasks — roadmap-ledger

## M1 — Ledger schema
- [ ] Define the single table: `seq | slug | title | tier | planned-target | status`
- [ ] Decide: stable `seq` vs row-position for ordering
- [ ] Decide ledger location (top of ROADMAP.md)
- [ ] Document the rule: version is derived from the ledger, never hand-written elsewhere

## M2 — De-duplicate references
- [ ] Convert ROADMAP prose + dep chains to slug/label refs (no absolute `v1.x`)
- [ ] Demote/remove `target_version:` from request PLAN frontmatter (source → derived/advisory)
- [ ] Verify: version mentions outside the ledger table drop to ~0

## M3 — Insert/reorder is one edit
- [ ] Fixture: insert a request = one ledger row + re-render, zero prose touched
- [ ] Document the before/after vs the policy-engine reslot (~14 refs → 1 row)

## M4 — Render + check
- [ ] Roadmap render reads per-version view from the ledger (coordinate with visual-layer)
- [ ] `doctor` flags any hardcoded version reference outside the ledger

## M5 — Migrate + ship
- [ ] Convert live ROADMAP.md to ledger-driven
- [ ] Dogfood: reslot a real request, confirm one-row edit
- [ ] CHANGELOG entry; plugin bump

## Resolve before building (from PLAN open questions)
- [ ] Stable id: sequence number vs slug-only
- [ ] target_version: leave frontmatter entirely vs read-only mirror
- [ ] Merge with cross-request-links vs separate adjacent requests
- [ ] Buffer/gap representation in a position-derived scheme
- [ ] Include shipped history in the ledger, or planned-only
