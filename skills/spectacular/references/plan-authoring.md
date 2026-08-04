---
description: PLAN authoring contract: seven sections, scope boundaries, and verification-artifact decision.
---
# PLAN Authoring

Use for PLAN scaffold/grill only. Write the required unnumbered sections in order:
`Goal`, `Constraints`, `Milestones`, `Tasks`, `Dependencies`, `Validation`, `Deliverables`.
Extra `Understanding`/`Decisions` sections may appear between them.

- Goal: one observable change tied to PRD intent, not a summary rewrite.
- Milestones: ordered demoable outcomes; TASKS is the executable checklist.
- Validation: one check per milestone with an authority (`run:`, assert, judge, or observable).
- Understanding is optional during planning but required before `planned → active`; use its three subheads or `UNDERSTANDING.md`.
- Decisions are request-scoped “chose X over Y because Z”; retain rejected alternatives.

Create `VERIFY.md` only when at least two apply: user-visible change, costly reversal,
multi-surface verification, non-trivial risk, external contract, or rollback plan.
Otherwise keep checks in PLAN Validation or TASKS Verification. Verification is always required.
