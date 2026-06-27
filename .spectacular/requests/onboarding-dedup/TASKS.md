---
status: planned
updated: 2026-06-27
related:
  - PLAN.md
---

# Tasks — onboarding-dedup

## v1

### M1 — Identify shared spine vs deltas
- [ ] Mark onboarding.md lines that duplicate status.md's read sequence
- [ ] Mark onboarding-specific content (always-run substrate check, takeover tone, gap-observations table, pre-split detection, example briefing)

### M2 — Refactor onboarding.md
- [ ] Replace the 6-step read sequence with a reference to status.md's flow
- [ ] Keep only the onboarding deltas + onboarding-specific sections
- [ ] Confirm status.md reads as the single owner of the briefing flow (minor edits only if needed)

### M3 — Verify
- [ ] First-invocation simulation: substrate check runs, takeover-tone briefing, ≤3 observations
- [ ] `/spectacular` (status) on a warm workspace unchanged

### M4 — Guided first-run
- [ ] Detect empty/new workspace (no requests, fresh init)
- [ ] Skill flow: usher new/PRD-grill → first request → point at `spectacular next`
- [ ] Decide CLI entry: `init --walk` vs auto-usher on empty `/spectacular`
- [ ] One step at a time — never dump the verb surface
- [ ] Verify: empty workspace ushers; warm workspace unaffected
