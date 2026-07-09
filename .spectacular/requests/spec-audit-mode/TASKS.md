---
status: review
updated: 2026-07-09
related:
  - PLAN.md
---

# Tasks — spec-audit-mode (frontmatter schema check)

<!--
  Pivoted 2026-07-09 (grill-me): semantic coverage audit → mechanical
  frontmatter schema check. Old M1–M4 (orphan bullets/files, stale specs)
  dropped. One milestone: the schema check.
-->

## v1

### M1 — Frontmatter schema check on flat capability specs
- [x] Extend `check_specs` flat-file loop: required keys `status, updated, summary, related` present → warning if missing
- [x] ISO-date check on `updated:` (reuse `check_frontmatter` pattern) → warning if not `YYYY-MM-DD`
- [x] Closed status enum `draft | published | deprecated` → warning if outside
- [x] Conditional version: `version:` required iff `status: published` → warning if published-without-version
- [x] Skip `index.md` (different doc-class — catalog, not capability)
- [x] Tests 11–17 in `tests/cli/specs.test.sh` (7 scenarios: one per rule + both conditional-version branches + index-skip guard)
- [x] Verify: full suite green (36/36); `doctor specs` on this repo flags nothing false (doc-engine draft clean, roadmap published clean, index skipped)
- [ ] ROADMAP ledger row: build b11 → target version
- [x] Doc note: `doctor` area reference (specs table) describes the schema
