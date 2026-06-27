---
status: verified
updated: 2026-06-28
related:
  - PLAN.md
---

# Verify Log — onboarding-dedup

Validation walk for v1.21.0. Doc/skill-ref request — no VERIFY.md (2-of-6 not met); checks map to PLAN § 6.

| Check | Kind | Evidence | Result |
|---|---|---|---|
| M2 — onboarding.md no longer restates the read sequence | observable | onboarding.md "Run the [[status]] flow, with these deltas"; the 7-step list is gone; reader sent to status.md | ✅ |
| M2 — status.md is the single briefing-flow owner | observable | status.md § Steps + § Build the briefing unchanged as the canonical sequence; added empty-workspace branch only | ✅ |
| M3 — onboarding deltas preserved | observable | always-run substrate check, takeover tone, first-look observations (≤3), pre-split detection, example briefing all retained | ✅ |
| M3 — warm-workspace status unchanged | observable | status.md read sequence + signal table + priority ranking untouched | ✅ |
| M4 — empty workspace ushers, doesn't brief | observable | guided-first-run.md walk (describe→PRD→first request→`next`); status.md + onboarding.md both branch to it on empty | ✅ |
| M4 — one step at a time, no verb dump | observable | guided-first-run.md "Core rule — one step at a time" + "does NOT list the verb surface" | ✅ |
| M4 — empty/existing distinction explicit | observable | onboarding = existing+work; guided-first-run = blank slate; both SKILL.md rows + each doc's intro state it | ✅ |
| Links resolve | executable | `doctor links` clean — [[status]], [[guided-first-run]], [[onboarding]] all resolve | ✅ |
| Regression | executable | `./tests/run.sh` → 9/9 areas pass | ✅ |

All checks pass. onboarding.md dedup + new guided-first-run.md; status.md gains one empty-workspace branch.
