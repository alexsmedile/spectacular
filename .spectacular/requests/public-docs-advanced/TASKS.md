# Tasks — Public Docs Advanced

> **Status: planned, gated.** Do not begin until 2-of-3 activation triggers fire (see PLAN.md § Activation triggers).

## M1 — Renderer adapter spec

- [ ] Write `skills/spectacular/references/docs-renderer-adapters.md`
- [ ] Mintlify mapping: docs.yaml → mint.json schema
- [ ] Docusaurus mapping: docs.yaml → sidebars.js + docusaurus.config.js
- [ ] Fumadocs mapping: docs.yaml → meta.json per folder
- [ ] MkDocs mapping: docs.yaml → mkdocs.yml nav
- [ ] Frontmatter field mapping table per renderer

## M2 — Export CLI

- [ ] `spectacular docs export <renderer> [--out path]`
- [ ] Mintlify implementation
- [ ] Docusaurus implementation
- [ ] Fumadocs implementation
- [ ] MkDocs implementation
- [ ] Smoke test: each adapter produces a config that the renderer's CLI accepts

## M3 — Versioning

- [ ] `spectacular docs publish <version>` — snapshot docs/ → docs/versioned/v<x.y.z>/
- [ ] docs.yaml schema extension: `versions:` list
- [ ] Doctor: versioned snapshots are read-only, edits = warning
- [ ] Idempotent re-run (re-publishing same version is no-op)

## M4 — Spec sync

- [ ] `spectacular docs sync-from-spec <spec-path>` — draft user-facing page
- [ ] Page frontmatter additions: `synced_from`, `synced_at`
- [ ] `spectacular docs sync-ack <page>` — bump synced_at without content change
- [ ] Doctor drift check: spec mtime > synced_at → warn
- [ ] Skill prompts user before overwriting any existing page

## M5 — Pack integration

- [ ] Add `docs-layout` to convention-pack rule categories (extends packs-contract.md to 7 categories)
- [ ] Schema: `required-sections`, `renderer`, `required-frontmatter`, `changelog-handling`
- [ ] Doctor `conventions` area: validate docs/ against pack docs-layout
- [ ] alex-default pack: add docs-layout block as reference example
- [ ] minimal pack: leave docs-layout empty {}

## M6 — Tests + release

- [ ] tests/cli/docs-export.test.sh — 4 renderer smoke tests
- [ ] tests/cli/docs-versioning.test.sh — publish, re-publish, doctor checks
- [ ] tests/cli/docs-sync.test.sh — sync, ack, drift detection
- [ ] tests/cli/pack.test.sh — extend with docs-layout scenarios
- [ ] Version bump (TBD)
- [ ] CHANGELOG entry
- [ ] Docs for the new verbs (in docs/ itself — meta dogfood)
- [ ] Tag + release notes
