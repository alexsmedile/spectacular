---
status: planned
updated: 2026-05-29
related:
  - PLAN.md
---

# Tasks — cross-request-links

## v1

### M1 — Schema extension
- [ ] Document `depends-on:` / `blocks:` in ARCHITECTURE.md (siblings to `related:`)
- [ ] Specify the computed-not-stored inverse rule

### M2 — Inverse-link resolver
- [ ] Scan all PLAN frontmatter, build the relationship graph
- [ ] Derive `blocked-by` from `blocks` at read time (never write to target)

### M3 — doctor links area
- [ ] Validate link integrity across requests + archive
- [ ] Flag dangling references (missing/archived slugs) as warnings
- [ ] Side-rider (FEEDBACKS.md 🟢): `doctor memory` staleness flag — age-check mirroring sessions/feedback/ideas convention (pick threshold; memory is less time-sensitive so lean conservative)

### M4 — status + new surfaces
- [ ] `spectacular status` advisory: related active requests
- [ ] `spectacular new` prompts to declare relationships on keyword match

### M5 — Examples + ship
- [ ] 2 example projects demonstrating the link graph
- [ ] CHANGELOG [1.13.0] entry; plugin bump to v1.13.0

## v2 (deferred)

- [ ] ROADMAP-as-source-of-truth enforcement (every active request links to a version)
- [ ] Visual link-graph render (depends on [[visual-layer]])
- [ ] Memory contradiction-check (FEEDBACKS.md 🟢): flag a memory whose claim a later session/decision overturns — cross-doc semantic reasoning, likely skill-side not CLI
