---
status: active
updated: 2026-06-02
related:
  - PLAN.md
---

# Tasks — imagine-mode

## v1

### M1 — Contract
- [x] Settle [blocking] Q1 (mode → distinct `imagine`)
- [x] Settle [blocking] Q2 (layout → typed subfolders)
- [x] Settle [blocking] Q3 (handoff → auto-offer + Understanding pre-fill)
- [x] Settle Q4 (approval → fragment frontmatter)
- [x] Fill PLAN `## Understanding`
- [x] Write `vision-rules.md` (frontmatter, slots, fragment kinds, `imagine` mode)
- [x] Add `vision` row to `doc-index.md`
- [x] Register `imagine` in the mode taxonomy
- [x] Add ARCHITECTURE.md section: `vision/` substrate + imagination-backed thesis

### M2 — vision/ soft-folder substrate ✅ (2026-06-02)
- [x] `spectacular imagine <slug>` scaffolds `requests/<slug>/vision/` (spine + subfolders)
- [x] `spectacular vision add <kind> <name>` mutator writes a fragment file
- [x] Spine `VISION.md` manifest regenerates from fragment files (index mode)
- [x] `templates/vision/` — spine + per-kind fragment scaffolds
- [x] `doctor vision` area: fragment frontmatter + kind/subfolder match + manifest drift (+`--fix`) + dangling persona refs + approval progress
- [x] Registered: KNOWN_DOCS, DOC_AREAS, run_areas, doctor_parse_args, router (`imagine`/`vision`)

### M3 — generative render engine ✅ (2026-06-02)
- [x] `imagine` mode behavior — `references/imagine.md` (render → react → derive; leads with artifacts)
- [x] Spine render spec: end-goal + macro dev phases + flow walk
- [x] Story / UI / arch fragment render spec (≥1 each; story pulls PERSONAS.md)
- [x] Q8 ASCII palette — shipped as a light convention in imagine.md (+ template scaffolds)
- [x] SKILL.md: Imagine-mode routing section + triggers; vision-rules links to imagine.md

> M4 + M5 are **specified** in `references/imagine.md` (steps 2 + 3 of the loop) — they
> need no code beyond the M2 mutators (`approved:` frontmatter, manifest regen) + PLAN
> review. They are exercised live by the M6 dogfood, which is their real validation.

### M4 — react-on-parts loop ✅ spec'd in imagine.md §2 (validated by M6)
- [x] Per-fragment `approved: true|false|pending` (M2 frontmatter + engine §2)
- [x] Approve/redirect/reject one fragment at a time (engine §2)
- [x] Regenerate only redirected fragments (engine §2)
- [x] Q4 approval substrate → fragment frontmatter (settled M1)

### M5 — Build derivation ✅ spec'd in imagine.md §3 (validated by M6)
- [x] Approved vision → draft PLAN mapping table (engine §3)
- [x] Pre-fill PLAN `## Understanding` from the vision (engine §3)
- [x] Hand off to PLAN review; draft never auto-accepted (engine §3)

### M6 — dogfood + ship
- [x] Run `imagine` on a real in-repo request end-to-end (`progress-view`; vision + 4 fragments, all approved)
- [x] Dogfood found + fixed 2 bugs: `--caption=` for dash-leading values; manifest drift now catches approval-state (feedback/2026-06-02-dogfood-progress-view, resolved)
- [x] Test suite green for new verbs + doctor area (9/9)
- [x] CHANGELOG [1.15.0] entry; plugin bump to v1.15.0

## v2 (deferred — see archive/ideas/explore-mode.md)

- [ ] Compare/reconcile derivation — diff an existing PRD/PLAN against the vision, surface gaps
- [ ] Project altitude — `imagine` before/around PRD; output at `.spectacular/vision/` (gated on Q5/Q6)
- [ ] Update PRD positioning copy to "spec-driven AND imagination-backed" (Q6)
- [ ] Rich diagram types beyond box/flow (sequence, state)
- [ ] Auto-promotion of fragments → tasks (stays human judgment)
