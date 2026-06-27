---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — onboarding-dedup

## v1

### M1 — Identify shared spine vs deltas
- [x] Marked the shared spine (config→root docs→requests→SPEC→memory→briefing) vs onboarding deltas
- [x] Identified onboarding-only: always-run substrate check, takeover tone, gap-observations table, pre-split detection, example briefing

### M2 — Refactor onboarding.md
- [x] Replaced the 7-step read sequence with "Run the status.md flow, with these deltas"
- [x] Kept only the onboarding deltas + onboarding-specific sections
- [x] status.md confirmed as single owner of read sequence + briefing format (added empty-workspace branch)

### M3 — Verify
- [x] onboarding.md no longer restates the read sequence; sends reader to status.md
- [x] `doctor links` clean (status/guided-first-run wikilinks resolve); warm-workspace status unchanged

### M4 — Guided first-run
- [x] Detect empty/new workspace (status.md + onboarding.md both branch to guided-first-run)
- [x] New `guided-first-run.md`: usher describe→PRD-grill(optional)→first request→point at `spectacular next`
- [x] CLI entry decided: skill-driven (no flag needed — `next` already ushers empty; `init --walk` left as future optional)
- [x] One step at a time — "never dump the verb surface" rule explicit in the doc
- [x] Verify: empty→ushers, warm→unaffected; routing wired in SKILL.md (with-prior-work→onboarding, empty→guided-first-run)
