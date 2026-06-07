---
status: active
updated: 2026-06-07
related:
  - PLAN.md
---

# Tasks — roadmap-ledger

## v1

### M1 — Ledger schema
- [ ] Define ledger table: `build | slug | title | tier | target-version | status`
- [ ] Ledger lives at the top of ROADMAP.md (above first version block)
- [ ] `spectacular new` stamps `build: bN` on new requests; increments `last_build:` in `config.yaml`
- [ ] Document the rule in ARCHITECTURE.md: version is derived from ledger, never hand-written in prose

### M2 — Remove target_version, add build id
- [ ] Remove `target_version:` from all active + planned request PLAN frontmatters
- [ ] Add `build: bN` to each (assign sequential ids to existing requests)
- [ ] Update `scaffold-reference.md` + PLAN template: `build:` in, `target_version:` out
- [ ] Update `spectacular new` to write `build: bN` instead of `target_version:`
- [ ] Update `plan-rules.md`: prose must not repeat version numbers; version lives in ledger only

### M3 — De-duplicate ROADMAP prose
- [ ] Convert ROADMAP block headings + dep chain prose to slug/label refs (no `v1.x` outside ledger)
- [ ] Verify: `grep -c "v1\.[0-9]" ROADMAP.md` outside the ledger table is ~0

### M4 — Insert/reorder is one edit
- [ ] Demonstrate: inserting a fixture request = one ledger row, zero prose touched
- [ ] Document before/after vs the policy-engine reslot (~14 refs → 1 row)

### M5 — Render + doctor check
- [ ] `spectacular roadmap` reads version blocks from the ledger (extend existing render)
- [ ] `doctor links` (from cross-request-links) flags stray hardcoded version refs outside the ledger

### M6 — Migrate + ship
- [ ] Convert live ROADMAP.md to ledger-driven (ledger table + slug refs in prose)
- [ ] Dogfood: reslot a real request, confirm one-row edit
- [ ] Coordinate final migration with cross-request-links M5 (both ship as v1.16.0)
- [ ] CHANGELOG entry; plugin bump to target release

## v2 (deferred)

- [ ] `target_version` computed field readable via `spectacular request <slug>` (reverse-lookup ledger)
- [ ] `spectacular new` UI shows the assigned build id + projected version at creation
