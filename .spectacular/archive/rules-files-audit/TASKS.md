---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — rules-files-audit

## v1

### M1 — Audit + decision
- [x] Read all 6 stub bodies (agents/architecture/principles/spec/stack/tasks) side by side
- [x] Classify each: thin-to-pointer (A) or write-real-body (B)
- [x] Confirm `tasks-rules.md` (51 lines) — real body (review-gate + refine patterns), NOT a stub despite mode label
- [x] Record stub-body policy + per-file dispositions in DECISIONS.md (D8)

### M2 — Shared default in doc-index.md
- [x] Add "Stub default behavior" section (grill no-op / refine rewrite / review structural + snapshot rule)

### M3 — Apply treatment
- [x] Thin A-files (architecture/principles/stack) to frontmatter + single pointer; agents keeps its top-level-AGENTS.md delta
- [x] Keep B-files: spec (index role + sync override) trimmed to deltas; tasks untouched (real body)
- [x] Confirm no frontmatter was touched (doc-id/mode/template intact on all 6)

### M4 — Verify no drift
- [x] Frontmatter dispatch preserved; bodies 160→129 lines
- [x] `spectacular doctor docs` clean (0 errors, 0 warnings)

### M5 — Collapse verify-doc trio
- [x] Merge verification.md (2-of-6, Part 2) + verify-tests.md (scripts, Part 3) into verify.md
- [x] Delete the two absorbed files (git rm)
- [x] Update SKILL.md routing + all inbound [[verification]]/[[verify-tests]] wikilinks + path refs → verify.md
- [x] Confirm verify walk runs identically; doctor links clean; 9/9 tests pass
