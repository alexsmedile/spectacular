---
type: decisions
doc: decisions
mode: index
updated: 2026-06-16
---

# Decisions

- **D1** — Request slugs stay unnumbered
- **D2** — 2026-05-21
- **D3** — 2026-05-11
- **D4** — PRD@v1.1.md snapshot gap
- **D5** — Use index mode for DECISIONS.md in all new projects past 50 entries — Agents load only the cheap index by default; full ADR prose loaded on demand per D<N>.md
- **D6** — Remove docs * verbs + deprecation_notice() + docs-* refs as MINOR (v1.17.0) — banner-warned since v1.2.0, pageworks is t — v2.0.0 major is left with a single breaking concern (the file-contract change). The CLI shrinks by ~700 lines. Any user still calling docs * gets a clear error pointing at pageworks instead of a deprecation banner.
