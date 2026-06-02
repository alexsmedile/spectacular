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

### M3 — generative render engine
- [ ] `imagine` mode behavior (skill ref) — leads with proposed artifacts
- [ ] Render spine: end-goal + macro dev phases + flow walk
- [ ] Render ≥1 story fragment from PERSONAS.md
- [ ] Render ≥1 ASCII UI/output fragment
- [ ] Render ≥1 ASCII architecture sketch
- [ ] ASCII palette decision (Q8) wired into templates if shipped

### M4 — react-on-parts loop
- [ ] Per-fragment `approved: true|false|pending` frontmatter
- [ ] Human approve/redirect/reject one fragment at a time
- [ ] Regenerate only redirected fragments (not the whole vision)
- [ ] Approval substrate decision (Q4) — frontmatter vs feedback/ entries

### M5 — Build derivation
- [ ] Approved vision → draft PLAN (stories→goals, flow→milestones, fragments→acceptance)
- [ ] Pre-fill PLAN `## Understanding` from the vision
- [ ] Hand off to existing PLAN grill/review (draft never auto-accepted)

### M6 — dogfood + ship
- [ ] Run `imagine` on a real in-repo request end-to-end
- [ ] Test suite green for new verbs + doctor area
- [ ] CHANGELOG [1.15.0] entry; plugin bump to v1.15.0

## v2 (deferred — see archive/ideas/explore-mode.md)

- [ ] Compare/reconcile derivation — diff an existing PRD/PLAN against the vision, surface gaps
- [ ] Project altitude — `imagine` before/around PRD; output at `.spectacular/vision/` (gated on Q5/Q6)
- [ ] Update PRD positioning copy to "spec-driven AND imagination-backed" (Q6)
- [ ] Rich diagram types beyond box/flow (sequence, state)
- [ ] Auto-promotion of fragments → tasks (stays human judgment)
