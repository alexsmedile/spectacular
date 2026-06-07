---
status: verified
updated: 2026-06-07
related:
  - PLAN.md
---

# Tasks — visual-layer

## v1

### M1 — ascii-render helper ✅ (2026-06-07)
- [x] Bar rendering (`███░░ 60%`), box-drawing, width clamping
- [x] NO_COLOR + non-TTY detection → strip color/escapes
- [x] Unit tests (bar math, clamping, NO_COLOR)

### M2 — Data-backed renders ✅ (2026-06-07)
- [x] `spectacular progress <slug>` → milestone bars + roll-up %
- [x] `spectacular summary` → dashboard (request-state bars + substrate counts)
- [x] Confirm `--format json` output byte-unchanged for both

### M3 — Roadmap render ✅ (2026-06-07)
- [x] `spectacular roadmap` → version arc / timeline (runway → major → vision)
- [x] Tier-aware rendering; plain-text degrade non-TTY

### M4 — App-UI mockup blocks ✅ (2026-06-07)
- [x] Define + document the ASCII layout-block format
- [x] Skill can drop a rendered mockup into a request PLAN/SPEC
- [x] 1+ real example block in a request (docs/visual-conventions.md)

### M5 — Docs + ship ✅ (2026-06-07)
- [x] `docs/visual-conventions.md` page + register in docs.yaml
- [x] CHANGELOG [1.15.0] entry; plugin bump to v1.15.0

## v2 (deferred)

- [ ] Visual link-graph render (consumes [[cross-request-links]])
- [ ] `spectacular status` visual treatment
- [ ] Burndown / exit-criteria % render per roadmap version
