---
status: planned
updated: 2026-05-29
related:
  - PLAN.md
---

# Tasks — visual-layer

## v1

### M1 — ascii-render helper
- [ ] Bar rendering (`███░░ 60%`), box-drawing, width clamping
- [ ] NO_COLOR + non-TTY detection → strip color/escapes
- [ ] Unit tests (bar math, clamping, NO_COLOR)

### M2 — Data-backed renders
- [ ] `spectacular progress <slug>` → milestone bars + roll-up %
- [ ] `spectacular summary` → dashboard (request-state bars + substrate counts)
- [ ] Confirm `--format json` output byte-unchanged for both

### M3 — Roadmap render
- [ ] `spectacular roadmap` → version arc / timeline (runway → major → vision)
- [ ] Tier-aware rendering; plain-text degrade non-TTY

### M4 — App-UI mockup blocks
- [ ] Define + document the ASCII layout-block format
- [ ] Skill can drop a rendered mockup into a request PLAN/SPEC
- [ ] 1+ real example block in a request

### M5 — Docs + ship
- [ ] `docs/<visual-conventions>.md` page + register in docs.yaml
- [ ] CHANGELOG [1.14.0] entry; plugin bump to v1.14.0

## v2 (deferred)

- [ ] Visual link-graph render (consumes [[cross-request-links]])
- [ ] `spectacular status` visual treatment
- [ ] Burndown / exit-criteria % render per roadmap version
