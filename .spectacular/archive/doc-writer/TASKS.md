---
status: verified
updated: 2026-05-21
related:
  - PLAN.md
---

# Tasks — Doc Writer

## v1

### M1 — Registry defined
- [x] Draft `references/doc-registry.md` with YAML schema + first entry (PRD)
- [x] Decide registry fields: `template`, `slots`, `mode`, `location`, `scope`, `snapshot-on-edit`, `overrides`, `kit-support`, `description`
- [x] Decide mode values: `grill` | `append` | `freeform`
- [x] Document override-file convention (when to create one, when not to)

### M2 — Engine extracted
- [x] Write `references/grill.md` — doc-agnostic slot-filler; reads registry to get slot list
- [x] Write `references/refine.md` — doc-agnostic vibe→spec engine + append-mode handler
- [x] Write `references/review.md` — doc-agnostic gate runner; checks base slots + override rules
- [x] Verify engine references contain zero hardcoded doc-type names (PRD/PLAN/etc. only appear in examples)

### M3 — PRD migrated
- [x] Write `prd-overrides.md` — extracted PRD-specific rules (kit selection, slot prompts, vague-word list, gate checks, vibe→spec tables)
- [x] Update registry entry to point to `prd-overrides.md`
- [x] Legacy `prd-grill.md` / `prd-refine.md` / `prd-review.md` kept for backwards compat per PLAN
- [ ] Dogfood: run `spectacular prd` (legacy trigger) → confirm identical output to pre-migration *(deferred to dogfood phase)*

### M4 — PLAN + TASKS added
- [x] Write `templates/plan/base.md` — 7-slot (Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables)
- [x] Write `templates/tasks/base.md` — checklist convention + frontmatter stubs
- [x] Add PLAN registry entry (mode: grill, scope: per-request)
- [x] Add TASKS registry entry (mode: freeform, scope: per-request)
- [x] Write `references/plan-overrides.md` — milestone-before-tasks ordering rule, slot prompts, gate checks
- [x] Write `references/tasks-overrides.md` — checklist format check, frontmatter sync
- [ ] Dogfood: grill a throwaway PLAN, confirm 7-slot output *(deferred to dogfood phase)*

### M5 — Remaining docs added
- [x] Write `templates/principles/base.md` (freeform — principles + enforcement hooks)
- [x] Write `templates/architecture/base.md` (freeform — layout + conventions)
- [x] Write `templates/roadmap/base.md` (freeform — time-ordered list)
- [x] Write `templates/stack/base.md` (freeform — tech choices list)
- [x] Write `templates/agents/base.md` (freeform — onboarding shape)
- [x] Write `templates/decisions/entry.md` (append — one ADR entry, not whole file)
- [x] Add registry entries for all 6
- [x] Decide which need overrides for v1 (decision: none — start clean, add later if patterns emerge)

### M6 — SKILL.md routing unified
- [x] Add generalized trigger handler: `spectacular <doc> <verb>` → engine
- [x] Preserve legacy `spectacular prd` / `spectacular prd grill` / `spectacular prd refine` / `spectacular prd review` aliases
- [x] Update SKILL.md References index + Templates index
- [x] Update SKILL.md trigger detection table
- [x] Bumped version 0.2.0 → 0.3.0; snapshot saved at versions/SKILL@0.2.0.md

### M7 — Dogfood
- [ ] Grill a fresh PLAN.md on a throwaway request, no PRD-specific paths invoked
- [ ] Grill a fresh PRINCIPLES.md, confirm freeform mode works (just scaffolds the template)
- [ ] Append a DECISIONS entry, confirm append mode works
- [ ] Run legacy `spectacular prd` and diff against pre-migration behavior

## v2 (deferred)

- [ ] Cross-doc consistency (PRD goals ↔ PLAN milestones)
- [ ] Project-local registry overrides (`.spectacular/doc-registry.yaml`)
- [ ] AI-suggested slot answers from existing project context
- [ ] Doc-to-doc translation (PRD success criteria → PLAN milestones generator)
