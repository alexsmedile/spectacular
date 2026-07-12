---
status: review
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

- [~] roadmap-rules.md core/doctrine split (7.2k — heaviest reference; same recipe as build/bug workflows)
- [~] debug-trace.md example-JSON diet (schemas keep one example each, trim field commentary)
- [~] verify.md walk-only split — Part 1 stands alone now that Parts 2–3 route elsewhere (dogfood finding 5)
- [~] Follow-up request `cli-gate-ergonomics` for the CLI-side dogfood findings: policy gate title+one-liner default (finding 1), advance auto-scaffolds SESSION.md (finding 2), doctor findings adjacent to summary (finding 3), v2 template items scaffold as `- [~]` (finding 4)
