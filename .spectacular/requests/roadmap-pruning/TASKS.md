---
status: planned
updated: 2026-06-28
related:
  - PLAN.md
---

# Tasks — roadmap-pruning

> Depends on b17 (roadmap-contract-docs) — ledger must be specced first.

## v1

### M1 — Decide approach + spec it
- [ ] Choose A (prune-to-ledger) vs B (roadmap-index mode); resolve M-questions 1-2
- [ ] Update specs/roadmap/SPEC.md + ARCHITECTURE.md with the chosen retention/pruning model (enforce the stated "history → CHANGELOG" principle)

### M2 — Detection (doctor)
- [ ] `doctor roadmap`: flag shipped prose blocks beyond keep-window + the "Recently shipped" mirror as prunable (info/warning), relayed by `status`

### M3 — Prune mechanism
- [ ] CLI prune verb / `doctor --fix roadmap`: snapshot ROADMAP.md, remove/move shipped blocks (A: delete after CHANGELOG-presence check; B: move to roadmap/v*.md + index), dry-run default
- [ ] Dogfood: prune this repo's 12 shipped blocks + mirror (confirm CHANGELOG covers each)
- [ ] Tests + VERIFY-LOG
