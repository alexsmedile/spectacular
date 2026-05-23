---
status: verified
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Roadmap Richness

## M1 — roadmap-overrides.md draft
- [ ] Copy `references/prd-overrides.md` as the structural model
- [ ] Slot 1 prompt: `Status` (planned | active | shipped | cancelled)
- [ ] Slot 2 prompt: `Phase` (mvp | testing-artifacts | atomic-concepts | spec-write | build-artifacts | versioning | release; optional alpha/beta/stable qualifier)
- [ ] Slot 3 prompt: `Scope (in)` — list of capability bullets
- [ ] Slot 4 prompt: `Scope (out)` — list of capability bullets explicitly deferred
- [ ] Slot 5 prompt: `Exit criteria` — checklist; ≥1 required
- [ ] Slot 6 (autopopulated): `Linked requests` — engine reads PLAN.md frontmatter `target_version:` field
- [ ] Mini-refine pattern: scope-in/out overlap detection
- [ ] Mini-refine pattern: phase regression (later version claims earlier phase than predecessor)
- [ ] Vibe → spec rewrite table for vague release descriptions
- [ ] Review gate checks: every version has all 5 required slots populated; ≥1 exit criterion

## M2 — Template rewrite
- [ ] `templates/roadmap/base.md` rewritten with structured per-version shape
- [ ] Comment hints inline showing slot purpose + good/bad examples
- [ ] Frontmatter unchanged (`version`, `updated`, `summary`, `related`)
- [ ] Old-shape detection sentinel left in old template if useful for migration

## M3 — Registry switch
- [ ] `doc-registry.md` `roadmap:` entry: `mode: freeform` → `mode: structured`
- [ ] Add `overrides: references/roadmap-overrides.md` to entry
- [ ] Add `snapshot-on-edit: true` (was false in freeform)

## M4 — Dogfood: rewrite live ROADMAP
- [ ] `spectacular snapshot .spectacular/ROADMAP.md` → ROADMAP@v1.0.md
- [ ] Rewrite ROADMAP.md against new structure for current versions
- [ ] Populate `target_version:` field in active request PLAN frontmatter where missing
- [ ] Re-render Linked-requests via `roadmap refine`

## M5 — Doctor extension
- [ ] `check_workspace` adds info-level check for ROADMAP shape
- [ ] Detection rule: ROADMAP.md exists AND no `## Scope (in)` heading found AND no `## Phase` heading found
- [ ] Output: "ROADMAP uses pre-v0.7 shape — run `spectacular roadmap refine`"
- [ ] Test scenario in doctor.test.sh

## M6 — Tests + VERIFY
- [ ] Doctor scenario: old-shape ROADMAP → info line; new-shape → silent
- [ ] Tests for engine handling of structured roadmap mode (may already be covered by generic engine tests)
- [ ] VERIFY.md if 2-of-6 fires (currently expects 2: user-visible change + multi-surface flow → likely yes)

## Verification (per 2-of-6 rule)
Expected hits: user-visible change (new grill verbs become structured), multi-surface flow (template + overrides + registry + doctor). 2 of 6 → VERIFY.md scaffold likely.
