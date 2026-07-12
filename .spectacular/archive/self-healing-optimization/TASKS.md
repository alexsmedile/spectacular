---
status: verified
updated: 2026-07-12
related:
  - PLAN.md
---

# Tasks — self-healing-optimization

## v1

### M1 — SKILL.md lean pass
- [x] Verify each doomed section's canonical copy exists in its ref (feedback-loop.md, imagine.md, verify.md, lifecycle.md, .spectacular/AGENTS.md) before cutting
- [x] Cut feedback-loop surfacing/auto-promotion prose to routing rows + one-line pointer
- [x] Cut imagine scope notes to routing rows + one-line pointer
- [x] Collapse verification-routing table to the decision rows; doctrine sentence points at verify.md
- [x] Cut task-tracking section to a two-line pointer at .spectacular/AGENTS.md § Task tracking
- [x] Rewrite State awareness to defer to AGENTS.md's table + the read-verbs cold-start pattern
- [x] → check: wc -c SKILL.md ≤ 16500 (re-baselined, see PLAN ## Decisions); all pre-trim route targets still grep

### M2 — Coherence fixes
- [x] status.md review signal: offer the verification walk against the resolved artifact (VERIFY.md > TASKS Verification > PLAN §Validation), never "create VERIFY.md"
- [x] active-request.md: fix stale `SPEC.md` → `specs/index.md`
- [x] active-request.md: replace duplicated context-loading table with a pointer to .spectacular/AGENTS.md
- [x] → check: greps in PLAN §Validation M2 pass

### M3 — 2-of-6 extraction
- [x] Add compact 2-of-6 rule table to plan-rules.md (names verify.md Part 2 as canonical)
- [x] Point new-request.md's verification step at plan-rules.md's table
- [x] Update SKILL.md verification rows: scaffold/grill → plan-rules table; review→verified + `spectacular verify` → verify.md Part 1
- [x] → check: greps in PLAN §Validation M3 pass

## v2 (deferred)

All four items graduated to real requests on 2026-07-12 after a /grill-me pass (decisions recorded in each PLAN):

- [~] → [[cli-gate-ergonomics]] (b29, high): dogfood findings 1–3 — policy gate block=full/warn=title, advance auto-scaffolds SESSION.md, doctor findings block
- [~] → [[verify-split]] (b30, medium): dogfood finding 5 (verify.md walk-only + verify-authoring.md) + finding 4 (template `- [~]` patch); carries the roadmap-rules split + debug-trace diet as its own v2
