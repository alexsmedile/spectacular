---
status: active
updated: 2026-06-02
related:
  - PLAN.md
  - TASKS.md
---

# Session — imagine-mode

## Current state

Request just moved `planned → active` (2026-06-02). The four design questions that gated implementation are resolved and recorded in PLAN § Open questions:

- **Q1 mode** → distinct `imagine` mode
- **Q2 layout** → typed subfolders (`stories/` `ui/` `arch/`)
- **Q3 handoff** → auto-offer `→ plan` + pre-fill `## Understanding`
- **Q4 approval** → fragment frontmatter `approved:`

PLAN `## Understanding` is filled. M1 is now "write the contract" (no more decisions to make).

## Active task

**M1 — Write the contract. ✅ DONE (2026-06-02).** Shipped:
- `skills/spectacular/references/vision-rules.md` — full doc-type + `imagine` mode contract
- `doc-index.md` — `vision` row in per-request table + `imagine` row in mode taxonomy
- `ARCHITECTURE.md` — `vision/` in request anatomy tree + dedicated `vision/` section + thesis

**M2 — `vision/` soft-folder substrate. ✅ DONE (2026-06-02).** Shipped:
- `spectacular imagine <slug>` scaffolds `vision/` (spine + stories/ui/arch); skill stub for bare `imagine`
- `spectacular vision add <kind> <name> --slug <s> [--caption]` mutator (kind→folder), `approved: pending`, manifest auto-regen
- `_regen_vision_manifest` + `_vision_iter_fragments` helpers
- `check_vision` doctor area + `--fix` (manifest regen via `doctor_apply_mechanical_fixes`)
- `templates/vision/{spine,story,ui,arch}.md`
- All wiring: KNOWN_DOCS, DOC_AREAS, run_areas, doctor_parse_args, router
- Verified end-to-end in a temp workspace; suite 9/9; real-repo doctor clean (no new findings)

**Next: M3 — generative render engine (SKILL work, not CLI).**
- `imagine` mode behavior in a skill ref (leads with proposed artifacts)
- Render spine slots + ≥1 ASCII fragment per kind
- SKILL.md routing + triggers for `imagine`

## Blockers

None. Q8 (ASCII palette) decision lands at M3.

## Next actions

1. M3: write the `imagine` mode behavior — likely fold into `vision-rules.md` or a dedicated skill ref; wire SKILL.md triggers/routing.
2. M3: decide Q8 (ship a `templates/vision/` ASCII palette vs improvise).

## Note

7 pre-existing `doctor links` warnings (`related: PRD.md` resolving relative to request dir) are NOT from this work — every request PLAN has them; they're what the `cross-request-links` request's `doctor links` area will fix.
