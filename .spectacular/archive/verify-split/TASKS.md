---
status: verified
updated: 2026-07-12
related:
  - PLAN.md
---

# Tasks — verify-split

## v1

### M1 — Split verify.md
- [x] Create verify-authoring.md from Parts 2–3 (2-of-6 canonical, fold patterns, VERIFY.md shape, script promotion) with frontmatter + history note
- [x] Cut verify.md to Part 1 (check kinds, walk loop, recording, coherence pass, gate, retrospective, VERIFY-LOG shape)
- [x] Cross-pointers: verify.md → authoring for rule details; authoring → verify.md for the walk
- [x] → check: PLAN §Validation M1 — verify.md 13,994 bytes (≤14,000); all Part 2–3 headings grep in authoring only; walk sections grep in verify.md

### M2 — Link re-point sweep
- [x] SKILL.md verification rows: authoring-time rows → plan-rules § 2-of-6 / verify-authoring; walk rows unchanged
- [x] doc-index.md: not applicable — it catalogs grill-able *doc types*; verify was never listed (no row to update; skill references aren't indexed there)
- [x] lifecycle.md ×3, init-workflow.md, new-request.md: re-point 2-of-6 mentions
- [x] plan-rules.md canonical note → verify-authoring.md
- [x] `spectacular doctor links` clean
- [x] → check: PLAN §Validation M2 — `verify.md Part 2` greps 0 outside authoring's history note; also fixed policy-injection.md, docs/commands.md ×2, repo CLAUDE.md rows found in the sweep

### M3 — Template `- [~]` patch
- [x] templates/tasks/base.md v2 rows → `- [~]` (no second live copy: CLI heredoc scaffolds v1 only — no v2 section)
- [x] tasks-rules.md: one line — v2 items use `- [~]` until promoted to v1
- [x] Sandbox scaffold test, then remove sandbox
- [x] → check: PLAN §Validation M3 — sandbox card shows `0/9 (+4 def)`; `- [~]` rows counted deferred, never open

## v2 (deferred)

- [~] roadmap-rules.md core/doctrine split (7.2k — heaviest remaining reference; carried over from b28)
- [~] debug-trace.md example-JSON diet (carried over from b28)
