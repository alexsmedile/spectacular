---
status: verified
updated: 2026-06-08
related:
  - PLAN.md
---

# Tasks — cross-request-links

## v1

### M1 — Schema extension
- [x] Document `depends-on:` / `blocks:` in ARCHITECTURE.md (siblings to `related:`)
- [x] Specify the computed-not-stored inverse rule

### M2 — Inverse-link resolver
- [x] Scan all PLAN frontmatter, build the relationship graph (helper: `_links_parse_list`)
- [x] Derive `blocked-by` from `blocks`, `required-by` from `depends-on` at read time (never write to target)
- [x] New `spectacular links [<slug>] [--json]` verb — whole-graph dump or per-request
- [x] Surface inverse links in `spectacular request <slug>` detail

### M3 — doctor links extension + path fix
- [x] **Root-aware path resolution** — slugs → requests/ then archive/; root-doc names → .spectacular/ (kills 7 false `related:` warnings)
- [x] Validate `depends-on`/`blocks` targets across requests + archive
- [x] Flag dangling references (missing/archived slugs) as warnings
- [x] Side-rider (FEEDBACKS.md 🟢): `doctor memory` staleness flag — age-check mirroring sessions/feedback/ideas convention (lean conservative — memory less time-sensitive)

### M4 — status + new surfaces
- [x] `spectacular summary` link advisory: related active requests
- [x] `spectacular new` prompts to declare relationships on keyword match

### M5 — Examples + ship
- [x] 2 example projects demonstrating the link graph (`tests/cli/links.test.sh`)
- [x] CHANGELOG [1.16.0] entry; plugin bump to v1.16.0

## v2 (deferred)

- [ ] ROADMAP-as-source-of-truth enforcement (every active request links to a version)
- [ ] Visual link-graph render (depends on [[visual-layer]])
- [ ] Memory contradiction-check (FEEDBACKS.md 🟢): flag a memory whose claim a later session/decision overturns — cross-doc semantic reasoning, likely skill-side not CLI
