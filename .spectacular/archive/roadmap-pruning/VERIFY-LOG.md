---
status: active
updated: 2026-06-29
related:
  - PLAN.md
  - TASKS.md
---

# Verify log — roadmap-pruning

CLI behavior change (new verb + doctor area) → executable tests + dogfood on the
live repo, all run before claiming.

## M1 — Approach + spec

- **{judge}** Design call: **B (roadmap-index mode)**, keep newest 3 inline (decisions 1-2 in PLAN). ✅
- **{assert}** `specs/roadmap.md` gained "## Index mode — shipped-history scaling" (migrate verb, per-version files, Shipped index, mode detection, safety, enforces history→CHANGELOG); snapshotted @v2, version 1.2. ✅
- **{assert}** `ARCHITECTURE.md` planned-runway rule updated: shipped prose → `roadmap/v*.md` once migrated, pointing at the migrate verb + spec. ✅

## M2 — Doctor `roadmap` area

- **{assert}** New area registered in `DOC_AREAS`, arg-validation case, and `run_areas` dispatch. `check_roadmap` written. ✅
- **`run: ./cli/spectacular doctor roadmap`** (live, index mode) → "all 7 Shipped index lines have a corresponding file" + "all 7 per-version files are indexed", 0 errors. ✅
- **{assert}** Orphan detection: copy with deleted file → "orphan Shipped index line(s) with no file: …". ✅
- **{assert}** Stale detection: copy with extra `v0.0.1.md` → "per-version file(s) with no Shipped index line: v0.0.1". ✅
- **{assert}** Flat-mode nudge: pre-migration snapshot (@v5, 10 shipped inline) → info "10 shipped prose blocks inline (flat mode) — beyond keep=3". ✅
- **{assert}** inline-shipped counter returns 3 on the live (kept) ROADMAP — no false prune nudge. ✅

## M3 — Migrate mechanism

- **`run: bash tests/cli/roadmap-migrate.test.sh`** → 22 passed, 0 failed (dry-run / migrate / idempotence / doctor clean+orphan+stale+flat). ✅
- **{assert}** Dry-run reports moves, writes nothing. ✅
- **{assert}** Migrate moves shipped-beyond-keep to `roadmap/v*.md`, keeps newest 3 inline, writes `## Shipped` index before Icebox, leaves planned/active/vision blocks inline. ✅
- **{assert}** Idempotent — re-run reports "nothing to migrate". ✅
- **{assert}** Per-version files written before ROADMAP.md rewrite (no data loss on partial run — verified by the failed-sed run during dev that left the live file intact). ✅
- **{assert}** Bash 3.2 compatible (no `declare -A`; caught + fixed during dev). ✅

## Dogfood (this repo)

- **{judge}** Snapshotted ROADMAP → @v5, then migrated: **528 → 410 lines**. 7 oldest shipped (v1.9–v1.19, excluding the kept v1.20/21/22 + the `active` v1.17) moved to `.spectacular/roadmap/v*.md`. ✅
- **{judge}** No information lost: moved prose lives in `roadmap/v*.md`; facts in CHANGELOG. Removed the "Recently shipped" mirror (pure CHANGELOG dup) + two stale reconciliation block-quotes; fixed 3 `roadmap-overrides`→`roadmap-rules` refs. ✅
- **`run: ./cli/spectacular doctor roadmap`** → 0 errors on the migrated live repo. ✅

## Full suite

- **`run: bash tests/run.sh`** → 11 areas, all pass (roadmap-migrate added; 10→11). ✅
- **`run: ./cli/spectacular doctor`** → 0 errors. (1 warning = the known transient SPEC-drift date heuristic, unrelated; clears on spec-sync.) ✅

## Result

All three milestones complete; approach B shipped + dogfooded. Ready for `review` → verify walk.
