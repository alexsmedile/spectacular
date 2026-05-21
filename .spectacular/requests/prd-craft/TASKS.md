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

## v2 (deferred)

- [ ] Phase splitting (PRD → phases/01-...md, 02-...md)
- [ ] CLI command `spectacular prd` (Bash binary)
- [ ] Auto-detect kit from project context
- [ ] Quality-gate scoring (not just pass/fail)
