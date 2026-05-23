---
status: verified
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Roadmap Richness v2

## M1 — `roadmap-overrides.md` updates
- [ ] Add Outcome slot prompt between Slot 2 (Phase) and Slot 3 (Scope-in)
- [ ] Document Outcome as required for `full` + `themed` tiers; absent for `vision`
- [ ] Update Themed-tier slot variant to include Outcome
- [ ] Add "Beginner pattern" section — start vision, graduate themed, unlock full when target_version: links
- [ ] Add "Icebox-promotion ritual" section — 4-step walk (pick → version → tier → slots → delete from icebox)
- [ ] Update gate check 12 (date patterns) — extend scope to entire themed/vision blocks, not just scope/exit slots
- [ ] Add gate check 16: Outcome required by tier (full + themed → ≥1 sentence; vision → absent)
- [ ] Add gate check 17: full-tier row count tiered warning (silent ≤7, info 8-10, warning 11+)
- [ ] Add gate check 18: scope-out warning when scope-in ≥4 AND scope-out empty
- [ ] Rename all instances of "Bucket list" → "Icebox" in this file

## M2 — Template rewrite
- [ ] `templates/roadmap/base.md` — add `**Outcome:**` line in full-tier example block (between Phase and Scope-in)
- [ ] Add `**Outcome:**` to themed-tier example block
- [ ] Rename `## Bucket list` → `## Icebox` in template
- [ ] Update comment block at top to mention Outcome + Icebox naming
- [ ] Update bucket-promotion comment to reference the icebox-promotion ritual section in overrides

## M3 — Live ROADMAP dogfood
- [ ] `spectacular snapshot .spectacular/ROADMAP.md` (snapshot before edit)
- [ ] Add Outcome paragraph to v0.7.1 block (full, currently active)
- [ ] Add Outcome paragraph to v0.7.x block (themed)
- [ ] Add Outcome paragraph to v0.11.x block (themed)
- [ ] Add Outcome paragraph to v1.0.0 block (themed)
- [ ] Rename `## Bucket list` heading → `## Icebox`
- [ ] Scan themed/vision blocks for date patterns; remove if found
- [ ] Confirm row count: how many full-tier blocks (should be 1)

## M4 — Meta-phase aliases
- [ ] Extend Phase taxonomy table in `roadmap-overrides.md` with meta-phase grouping:
      DISCOVER = intent | discover | prototype
      BUILD    = spec-refine | mvp | iterate
      RELEASE  = test | release-prep | release
- [ ] Update gate check 8 (Phase valid) to accept both individual phase names AND meta-phase names
- [ ] Document the "start coarse, refine as work crystallizes" guidance in roadmap-overrides.md
- [ ] Add example showing a version starting at `Phase: build` then refining to `Phase: mvp`

## M5 — Doctor extension (light)
- [ ] `check_workspace` in `cli/spectacular`: scan ROADMAP.md for `## Bucket list` heading
- [ ] If found AND no `## Icebox` heading present: emit info line suggesting rename
- [ ] Mechanical fix (mechanical tag): `sed -i 's/^## Bucket list/## Icebox/' .spectacular/ROADMAP.md`
- [ ] Skip silently if `## Icebox` already present (post-migration)
- [ ] Add scenario to `doctor.test.sh`: workspace with Bucket-list ROADMAP → info; with Icebox → silent

## M6 — Tests + v0.7.2 release
- [ ] `doctor.test.sh` scenario: Bucket-list → Icebox info line
- [ ] Manual verification of gate checks 16/17/18 — author a test ROADMAP triggering each, document expected output
- [ ] CHANGELOG.md entry for v0.7.2
- [ ] `.claude-plugin/plugin.json` version bump to 0.7.2
- [ ] SPEC.md capability bullet update — structured ROADMAP gets Outcome slot + Icebox + meta-phase aliases
- [ ] CLAUDE.md Active Requests: remove roadmap-richness-v2; add to Archived (shipped)
- [ ] Live doctor: clean
- [ ] Dogfood: archive this request via `spectacular promote roadmap-richness-v2 --to verified --force --archive`

## Verification (per 2-of-6 rule)

Hits: user-visible change (new slot, renamed section, new gate checks), multi-surface flow (template + overrides + doctor + live ROADMAP), external contract (Outcome slot is a public field; meta-phase aliases are public values). 3 of 6 → VERIFY.md required at review.

Carry to VERIFY.md when ready:
- [ ] Outcome slot present in template; gate flags missing Outcome in full/themed
- [ ] Icebox rename complete across template, overrides, live ROADMAP
- [ ] Gate emits warning on themed block with date pattern
- [ ] Gate emits info at 8 full-tier blocks; warning at 11
- [ ] Gate emits warning when scope-in ≥4 AND scope-out empty
- [ ] Meta-phase aliases accepted (`Phase: build` doesn't fail gate)
- [ ] Doctor flags Bucket-list ROADMAP; silent on Icebox
- [ ] All 7 existing test files still green
- [ ] Live workspace dogfood successful
