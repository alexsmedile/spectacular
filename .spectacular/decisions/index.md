---
type: decisions
doc: decisions
mode: index
version: 1.0
updated: 2026-06-28
summary: "Architectural and product decisions log"
---

# Decisions

- **D1** — Request slugs stay unnumbered
- **D2** — 2026-05-21
- **D3** — 2026-05-11
- **D4** — PRD@v1.1.md snapshot gap
- **D5** — Use index mode for DECISIONS.md in all new projects past 50 entries — Agents load only the cheap index by default; full ADR prose loaded on demand per D<N>.md
- **D6** — Remove docs * verbs + deprecation_notice() + docs-* refs as MINOR (v1.17.0) — banner-warned since v1.2.0, pageworks is t — v2.0.0 major is left with a single breaking concern (the file-contract change). The CLI shrinks by ~700 lines. Any user still calling docs * gets a clear error pointing at pageworks instead of a deprecation banner.
- **D7** — Skill description length check gates on `description` alone at 1024 chars (error), 1000 (warning) — not the description+w — Codex measures description ALONE: the v1.17.2 patch took description 1146→986 while description+when_to_use stayed 1253 (>1024) and the error cleared — proving the concatenation is not what Codex caps. check_skill() and the pre-commit guard both measure description alone.
- **D8** — Stub rules-file bodies thin to frontmatter + one pointer; the shared 3-verb default lives once in doc-index.md — architecture/principles/stack thin to pointer; agents keeps its top-level-AGENTS.md delta; spec and tasks keep real bodies (spec=index+sync role, tasks=full review/refine spec — mislabeled mode:stub). doc-index.md gains a 'stub default behavior' section. Frontmatter untouched everywhere — engine dispatch intact.
- **D9** — Undo v1: single-level breadcrumb (.last-mutation), staleness via timestamp-vs-mtime, idea-promote undo prompts before re — v1 scope stays tight: no breadcrumb stack, no git-reflog coupling, no destructive auto-delete. Each can graduate to v2 if real use demands. undo refuses on stale breadcrumb rather than mis-reverting.
