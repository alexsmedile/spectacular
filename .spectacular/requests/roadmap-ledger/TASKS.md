---
status: active
updated: 2026-06-14
related:
  - PLAN.md
---

# Tasks — roadmap-ledger

## v1

### M1 — Ledger schema
- [x] Document ledger schema in ARCHITECTURE.md: columns `build | slug | title | tier | target-version | status`; tier legend (`full` = near-term detailed, `themed` = mid-term directional, `vision` = long-horizon direction-only); status values (`planned | active | shipped`, release-level, distinct from request lifecycle); grouped-build rule (two rows same version = fine); computed-not-stored rule (version lives only here)
- [x] Document ledger placement rule in ARCHITECTURE.md: lives at top of ROADMAP.md above first version block; human adds rows manually when slotting a request; `spectacular new` does NOT auto-insert rows
- [x] CLI: `spectacular new` reads `last_build:` from `config.yaml` (treat missing as `0`), stamps `build: b(N+1)` on new PLAN.md, writes `last_build: N+1` back to `config.yaml`
- [x] CLI: `spectacular new` output prints `✓ build id: bN` and a "add a row to the ledger in ROADMAP.md when slotting" hint

### M2 — Remove target_version, add build id
- [x] Remove `target_version:` from all active + planned request PLAN frontmatters
- [x] Add `build: bN` to each (assign sequential ids to existing requests sorted by updated: date; b3–b8; last_build: 8)
- [x] Update PLAN template (`skills/spectacular/templates/plan/base.md`): `build: <BUILD>` in, no `target_version:`
- [x] Update `spectacular new`: substitutes `<BUILD>` from template; removed `--target-version` flag + help text
- [x] Update `plan-rules.md`: frontmatter schema note + version-in-prose rule (no hardcoded vX.y in milestone text)

### M3 — De-duplicate ROADMAP prose
- [x] Convert ROADMAP block headings + dep chain prose to slug/label refs (no `v1.x` outside ledger)
- [x] Verify: `grep -c "v1\.[0-9]" ROADMAP.md` outside the ledger table is ~0

### M4 — Insert/reorder is one edit
- [x] Demonstrate: inserting a fixture request = one ledger row, zero prose touched
- [x] Document before/after vs the policy-engine reslot (~14 refs → 1 row)

### M5 — Render + doctor check
- [ ] `spectacular roadmap` reads version blocks from the ledger (extend existing render)
- [ ] `doctor links` (from cross-request-links) flags stray hardcoded version refs outside the ledger

### M6 — Migrate + ship
- [ ] Convert live ROADMAP.md to ledger-driven (ledger table + slug refs in prose)
- [ ] Dogfood: reslot a real request, confirm one-row edit
- [ ] Coordinate final migration with cross-request-links (already shipped as v1.16.0; roadmap-ledger ships in the next available slot)
- [ ] CHANGELOG entry; plugin bump to target release

## v2 (deferred)

- [ ] `target_version` computed field readable via `spectacular request <slug>` (reverse-lookup ledger)
- [ ] `spectacular new` UI shows the assigned build id + projected version at creation
