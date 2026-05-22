---
status: active
updated: 2026-05-21
related:
  - PLAN.md
  - ../../PRD.md
---

# Tasks — PRD Craft

## v1

- [x] Scaffold request (PLAN.md, TASKS.md)
- [x] Write `templates/prd/base.md` — 6 required slots, 3 optional
- [x] Write starter kits:
  - [x] `kits/coding.md`
  - [x] `kits/product.md`
  - [x] `kits/content.md`
  - [x] `kits/research.md`
  - [x] `kits/blank.md`
- [x] Write `references/prd-grill.md` — strict slot-order interview loop
- [x] Write `references/prd-refine.md` — vibe→spec patterns, `[NEEDS CLARIFICATION]` convention
- [x] Write `references/prd-review.md` — quality gate checklist + vague-word list
- [x] Wire SKILL.md routing — `prd`, `prd refine`, `prd review` triggers
- [x] Verify `references/prd-review.md` contains the vague-word list referenced by grill mini-refine (confirmed at prd-review.md:54, with `prd-refine.md` as canonical source)
- [x] ~~Update root `.spectacular/PRD.md` with PRD-vs-PLAN clarifier~~ — **superseded by `canonical-docs-rework`**. Clarifier now lives in `ARCHITECTURE.md § Request files` (primary copy) and `PRD.md § Deliverable` (one-liner).
- [ ] Dogfood: run `prd review` against `.spectacular/PRD.md`; snapshot to `PRD@v1.3.md` if edits land
- [ ] Test on a fresh blank project — produce a usable PRD in <15 minutes
- [ ] Investigate missing `PRD@v1.1.md` snapshot (v1.0 and v1.2 exist; v1.1 gap unexplained)

## v1.1 — Slot alignment + kit-prep

Lands before dogfood. Defines the stable 8-slot base for [[kits-as-plugins]] and [[smart-init]] to build on.

### Base template
- [x] Rename slot 2 in `templates/prd/base.md`: "Who it's for" → "Target users"
- [x] Rename slot 3 in `templates/prd/base.md`: "What success looks like" → "Goals & success criteria"
- [x] Add slot 0 "Vision" (narrative, one paragraph) to `templates/prd/base.md`
- [x] Add slot "Deliverable" between Target users and Goals & success criteria in `templates/prd/base.md`
- [x] Renumber inline comments + section headings to reflect 8-slot order
- [x] Update all 5 kits to match 8-slot order (kit refactor to diff-only still deferred to [[kits-as-plugins]])

### Grill
- [x] Update `references/prd-grill.md`: slot loop covers 8 slots (was 6)
- [x] Update slot intro language for Vision (more open-ended than the measurable slots)
- [x] Verify mini-refine still applies correctly to renamed slots

### Review gate
- [x] Update `references/prd-review.md`: pass criteria check all 8 slots
- [x] Keep "at least one number + verb + date/timeframe" rule on Goals & success criteria slot
- [x] Vague-word list unchanged; added Vision/Deliverable exemptions + concrete-deliverable check

### Refine patterns
- [x] Update `references/prd-refine.md`: any slot references use new names
- [x] Added vague-deliverable pattern; Vision intentionally exempt per design

### Dogfood — moved here from v1
- [ ] Run `prd review` against root `.spectacular/PRD.md` using new 8-slot gate
- [ ] Snapshot root PRD to `PRD@v2.0.md` if edits land
- [ ] Test on a fresh blank project — produce a usable 8-slot PRD in <15 minutes

### Cleanup carried from v1
- [ ] Investigate missing `PRD@v1.1.md` snapshot (v1.0 and v1.2 exist; v1.1 gap unexplained)

## v2 (deferred)

- [ ] Phase splitting (PRD → phases/01-...md, 02-...md)
- [ ] CLI command `spectacular prd` (Bash binary)
- [ ] Auto-detect kit from project context (covered by [[smart-init]] v2)
- [ ] Quality-gate scoring (not just pass/fail)
- [ ] Kit refactor to diff-only format → **moved to [[kits-as-plugins]]**
- [ ] Smart-init doc scaffolding → **moved to [[smart-init]]**
