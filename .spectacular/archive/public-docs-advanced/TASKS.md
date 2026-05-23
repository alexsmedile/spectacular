---
updated: 2026-05-23
supersedes: TASKS@v1.0.md
related:
  - PLAN.md
---

# Tasks — Public Docs Advanced (narrowed)

> Earlier 6-milestone scope preserved as `TASKS@v1.0.md`. This revision matches PLAN's M1–M5.

## M1 — Adapter reference doc ✓

- [x] Create `skills/spectacular/references/docs-renderer-adapters.md`
- [x] Mapping table: `docs.yaml` → `mkdocs.yml` (nav, theme, plugins, hooks)
- [x] Mapping table: `docs.yaml` → `docusaurus.config.js` + `sidebars.js`
- [x] Page frontmatter translation table (both renderers)
- [x] GitHub Pages deploy boilerplate documented inline for both
- [x] "Community contributions welcome" pointer for Mintlify + Fumadocs
- [x] Reference linked from `docs-contract.md` § Renderer-agnostic by design

## M2 — `docs export` CLI ✓

- [x] Add `docs export` subcommand to `cli/spectacular` dispatcher
- [x] `docs export mkdocs [--out <path>] [--force] [--no-workflow]` writes `mkdocs.yml`
- [x] `docs export docusaurus [--out <path>] [--force] [--no-workflow]` writes `docusaurus.config.js` + `sidebars.js`
- [x] Both: write `.github/workflows/docs.yml` (idempotent, skip if present, `--no-workflow` opts out)
- [x] `--force` overwrites; default skips with a clear "skipped (use --force)" report
- [x] `// spectacular: do-not-overwrite` magic-comment pin respected even with `--force`
- [x] Empty sections dropped from nav (avoids invalid mkdocs YAML)
- [x] Help text: `spectacular docs --help` lists `export` and renderers
- [x] Unknown renderer (mintlify/fumadocs): actionable error pointing to M1 reference
- [x] Top-level usage (`top_usage()` in CLI) gains the `docs export` entry
- [x] Sandbox-verified 2026-05-23: mkdocs + docusaurus generation, idempotency, --force, pin, error paths

## M3 — `docs.yaml` schema extension

- [ ] Extend `docs.yaml` schema with optional `renderers:` block — full schema section in `docs-contract.md`
- [ ] Update `templates/docs/docs.yaml` template (if present) with commented `renderers:` example
- [ ] Update `docs_init` in `cli/spectacular` to emit commented `renderers:` stub
- [ ] Doctor `docs` area: validate `renderers:` if present (non-required), warn on unknown renderer keys (anything outside `mkdocs`, `docusaurus`)
- [ ] Doctor scenario test for malformed `renderers:` block

## ~~M4 — Dogfood~~ DROPPED (architectural decision 2026-05-23)

Authoring spectacular's own docs/ is a job for a future `docs-writer` sub-skill, not for the spectacular orchestrator. Tracked downstream:

- New request: `spectacular-bundle-restructure` (convert spectacular into a bundle parent)
- New request: `docs-writer-skill` (depends on bundle restructure)
- New request: `public-docs-dogfood` (depends on docs-writer)

## M5 — Tests + release

- [ ] `tests/cli/docs-export.test.sh` — covers mkdocs + docusaurus output, idempotency, --force, pin, --no-workflow, unknown renderer, missing docs.yaml
- [ ] YAML structural validation for `mkdocs.yml` (parse with python `yaml.safe_load` if available, else basic key checks)
- [ ] JS syntax check for `docusaurus.config.js` + `sidebars.js` (`node --check` if available)
- [ ] Bump `cli/spectacular` `SPECTACULAR_VERSION` to `1.1.0`
- [ ] Bump `.claude-plugin/plugin.json` to `1.1.0`
- [ ] Bump `.codex-plugin/plugin.json` to `1.1.0`
- [ ] Bump `skills/spectacular/SKILL.md` frontmatter `version` to `1.1.0`
- [ ] Copy CLI to `~/.local/bin/spectacular`
- [ ] CHANGELOG entry under `[1.1.0]` — Added: `docs export mkdocs|docusaurus`, `renderers:` block in docs.yaml, adapter reference doc
- [ ] Update CLAUDE.md Active Requests table (remove this request, add to Archived line as v1.1.0)
- [ ] Snapshot PLAN + TASKS, archive request → `.spectacular/archive/public-docs-advanced/`
- [ ] Tag v1.1.0, push, `gh release create v1.1.0 --generate-notes`
- [ ] `/plugin marketplace update spectacular`

## Follow-on requests to scaffold after this ships

- [ ] Scaffold `spectacular-bundle-restructure` request — design conversion to bundle layout
- [ ] Scaffold `docuworks-skill` request — new sub-skill for docs authoring (toolbox: templates, prose patterns, tone, page-type contracts); reads `docs-contract.md` from spectacular as source of truth
- [ ] Scaffold `docs-writer-agent` request — Tier-4-style agent definition that loads `docuworks`; this is the role spectacular hands off to from `docs new` / page audits
- [ ] Scaffold `public-docs-dogfood` request — author spectacular's own docs/ using `docs-writer` + `docuworks` + deploy via the M2 adapters
