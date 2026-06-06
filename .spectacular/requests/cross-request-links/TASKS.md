---
status: active
updated: 2026-06-02
related:
  - PLAN.md
---

# Tasks — cross-request-links

## v1

### M1 — Schema extension
- [ ] Document `depends-on:` / `blocks:` in ARCHITECTURE.md (siblings to `related:`)
- [ ] Specify the computed-not-stored inverse rule

### M2 — Inverse-link resolver
- [ ] Scan all PLAN frontmatter, build the relationship graph (helper: `_links_graph`)
- [ ] Derive `blocked-by` from `blocks`, `required-by` from `depends-on` at read time (never write to target)
- [ ] New `spectacular links [<slug>] [--json]` verb — whole-graph dump or per-request
- [ ] Surface inverse links in `spectacular request <slug>` detail

### M3 — doctor links extension + path fix
- [ ] **Root-aware path resolution** — slugs → requests/ then archive/; root-doc names → .spectacular/ (kills 7 false `related:` warnings)
- [ ] Validate `depends-on`/`blocks` targets across requests + archive
- [ ] Flag dangling references (missing/archived slugs) as warnings
- [ ] Side-rider (FEEDBACKS.md 🟢): `doctor memory` staleness flag — age-check mirroring sessions/feedback/ideas convention (lean conservative — memory less time-sensitive)

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
