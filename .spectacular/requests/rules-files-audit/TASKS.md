---
status: planned
updated: 2026-06-27
related:
  - PLAN.md
---

# Tasks — rules-files-audit

## v1

### M1 — Audit + decision
- [ ] Read all 6 stub bodies (agents/architecture/principles/spec/stack/tasks) side by side
- [ ] Classify each: thin-to-pointer (A) or write-real-body (B)
- [ ] Confirm `tasks-rules.md` (51 lines) — partial stub or real?
- [ ] Record stub-body policy + per-file dispositions in DECISIONS.md

### M2 — Shared default in doc-index.md
- [ ] Add a "stub default verb behavior" section (grill no-op / refine rewrite / review structural)

### M3 — Apply treatment
- [ ] Thin A-files to frontmatter + single pointer line
- [ ] Write real bodies for any B-files
- [ ] Confirm no frontmatter was touched (dispatch intact)

### M4 — Verify no drift
- [ ] `spectacular <stubdoc> grill/refine/review` behaves identically pre/post
- [ ] `spectacular doctor docs` clean

### M5 — Collapse verify-doc trio
- [ ] Merge verification.md (2-of-6 rule) + verify-tests.md (script-vs-checklist) into verify.md as sections
- [ ] Delete the two absorbed files
- [ ] Update SKILL.md verification routing table to point at verify.md only
- [ ] Confirm the verify walk (review→verified) runs identically
