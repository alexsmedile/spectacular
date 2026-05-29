---
status: planned
updated: 2026-05-29
related:
  - PLAN.md
---

# Tasks — verify-walk

## v1

### M1 — Walk algorithm
- [ ] Write `references/verify.md`: locate VERIFY.md or fall back to PLAN § Validation
- [ ] Iterate each check; prompt for evidence; record pass / blocker per item
- [ ] Define the gate (all-pass vs any-blocker outcomes)

### M2 — Lifecycle wiring
- [ ] All-pass → flip `review → verified` via `promote`
- [ ] Any blocker → stay `review`, write the blocker list back to the request

### M3 — Retrospective + archive tie-in
- [ ] End-of-walk optional "what surprised you?" prompt → `memory/` entry
- [ ] `spectacular archive` warns when `verified` was never reached via the walk

### M4 — Surface + docs
- [ ] SKILL.md routing-table entry for `verify`
- [ ] CLI redirect: `spectacular verify <slug>` → skill-only message
- [ ] `docs/commands.md` agentic-verbs section covers `verify`

### M5 — Dogfood + ship
- [ ] Drive 1+ real request through the walk to `verified`
- [ ] CHANGELOG [1.11.0] entry; plugin bump to v1.11.0

## v2 (deferred)

- [ ] Auto-suggest the walk when a request hits `review` (proactive surface)
- [ ] Per-check evidence persistence (store evidence inline in VERIFY.md)
