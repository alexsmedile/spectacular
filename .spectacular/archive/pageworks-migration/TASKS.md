---
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — pageworks migration

## M1 — Scaffold pageworks

- [ ] Create `skills_db/pageworks/` directory
- [ ] Write `pageworks/SKILL.md` — frontmatter: name, description, when_to_use, version `0.1.0`, category `devtools`, status `published`, compatible_with `spectacular >= 1.0.0` (informational)
- [ ] Write `pageworks/README.md` — what pageworks is, install, why it superseded spectacular's docs surface
- [ ] Write `pageworks/CHANGELOG.md` — first entry `[0.1.0] — 2026-05-23 — first release; superseded the docs surface from spectacular v1.x`
- [ ] Write `pageworks/.claude-plugin/plugin.json` — mirrors spectacular's structure with pageworks name + 0.1.0 version
- [ ] Write `pageworks/.codex-plugin/plugin.json` — mirrors spectacular's structure
- [ ] Write `pageworks/CLAUDE.md` — same boundary doc style as spectacular's CLAUDE.md, explaining what pageworks is and isn't
- [ ] Bootstrap `pageworks/.spectacular/` with `spectacular init` so future pageworks work is tracked spectacular-style
- [ ] Add `pageworks/.gitignore` (covers `_archive/`, `_backup/`, `.spectacular.local/`)

## M2 — Migrate references (lift-and-shift)

- [ ] Copy `spectacular/skills/spectacular/references/docs-contract.md` → `pageworks/references/contract.md`; rename internal references; remove "v2 renderer adapters" forward-references that were resolved in spectacular v1.1.0
- [ ] Copy `spectacular/skills/spectacular/references/docs-overrides.md` → `pageworks/references/authoring.md`; rewrite to remove spectacular-engine vocabulary (kits, packs, doc-registry); become pageworks-native
- [ ] Copy `spectacular/skills/spectacular/references/docs-renderer-adapters.md` → `pageworks/references/renderers.md`; rename CLI references from `spectacular docs export` to `pageworks export`
- [ ] Copy `spectacular/skills/spectacular/templates/docs/docs.yaml.tmpl` → `pageworks/templates/docs.yaml.tmpl`
- [ ] Copy `spectacular/skills/spectacular/templates/docs/index.md.tmpl` → `pageworks/templates/index.md.tmpl`
- [ ] Copy `spectacular/skills/spectacular/templates/docs/page.md.tmpl` → `pageworks/templates/page.md.tmpl`
- [ ] Add new: `pageworks/templates/pages/tutorial.md.tmpl` (Diátaxis tutorial template)
- [ ] Add new: `pageworks/templates/pages/how-to.md.tmpl` (Diátaxis how-to template)
- [ ] Add new: `pageworks/templates/pages/reference.md.tmpl` (Diátaxis reference template)
- [ ] Add new: `pageworks/templates/pages/explanation.md.tmpl` (Diátaxis explanation template)
- [ ] Add new: `pageworks/references/page-types.md` — Diátaxis quadrant guide, when to use which template
- [ ] Add new: `pageworks/references/prose-patterns.md` — voice/tone, callout conventions, code-block conventions, link patterns, headings, prerequisites blocks
- [ ] Add new: `pageworks/references/maintenance.md` — drift detection, spec→doc sync patterns, `updated:` freshness rules
- [ ] Update `pageworks/SKILL.md` to load these references on the right triggers

## M3 — pageworks CLI binary

- [ ] Create `pageworks/cli/pageworks` (executable Bash) — shebang, constants block, helpers block (lift `die`, `info`, `created`, `skip`, `write_if_missing` from spectacular)
- [ ] Add `PAGEWORKS_VERSION="0.1.0"` constant
- [ ] Port `docs_init` → `pageworks_init` (scaffold docs/ + docs.yaml + index.md + 3 default sections; `--minimal` flag)
- [ ] Port `docs_export_mkdocs` → `pageworks_export_mkdocs`
- [ ] Port `docs_export_docusaurus` → `pageworks_export_docusaurus`
- [ ] Port `docs_export_workflow` → `pageworks_export_workflow`
- [ ] Port `docs_yaml_*` parsers (scalar, renderer_scalar, sections, extras)
- [ ] Port `check_docs` → `pageworks doctor` (with the slimmer pageworks subcommand surface)
- [ ] Port idempotency helpers (`docs_write_adapter_file`, `docs_is_pinned`)
- [ ] Add top-level dispatcher: `init | export | doctor | --version | --help`
- [ ] Add `pageworks/cli/install.sh` (curl-installable, copies binary to `~/.local/bin/pageworks`)
- [ ] Add `pageworks export --help` listing supported renderers (mkdocs, docusaurus)
- [ ] Bash syntax check: `bash -n cli/pageworks`

## M4 — Spectacular deprecation banners

- [ ] Add `deprecation_notice()` helper to `cli/spectacular`
- [ ] Wire `deprecation_notice "docs init"` into `docs_init` (prints to stderr; doesn't block execution)
- [ ] Wire same into `docs_export`, `docs new`, `docs review`, `docs status` paths
- [ ] Deprecation banner text: `⚠ 'spectacular docs <verb>' is deprecated in v1.2.0 and will be removed in v2.0.0. Use 'pageworks <verb>' — see https://github.com/alexsmedile/pageworks`
- [ ] Add `> **Deprecated in v1.2.0** — see [pageworks](https://github.com/alexsmedile/pageworks). This reference will be removed in spectacular v2.0.0.` banner to:
  - `skills/spectacular/references/docs-contract.md`
  - `skills/spectacular/references/docs-overrides.md`
  - `skills/spectacular/references/docs-renderer-adapters.md`
- [ ] Slim `check_docs` in `cli/spectacular`:
  - keep: docs/ folder presence check
  - keep: docs.yaml manifest presence check
  - add: pageworks-installation check (file exists? `which pageworks` succeeds?) — info-level only
  - remove: section/page declaration validation
  - remove: frontmatter field validation
  - remove: renderers block validation
  - remove: orphan detection
- [ ] Update `spectacular docs --help` to show deprecation in the header
- [ ] Update `top_usage()` in spectacular CLI — note that `docs` is deprecated

## M5 — Handoff wiring

- [ ] New ref: `skills/spectacular/references/pageworks-handoff.md` — boundary doc; when to delegate; install hint; example user-facing handoff message
- [ ] Wire handoff prompt in `cmd_archive`: if archived request has SPEC.md or specs/ changes, print prompt: *"This change may affect public docs/. Run `pageworks` to update them, or skip with --no-docs-prompt."* (Confirmation is informational; spectacular does not invoke pageworks itself.)
- [ ] Add `--no-docs-prompt` flag to `cmd_archive`
- [ ] Add `docs.prompt_on_archive: true` / `false` to `.spectacular/config.yaml` schema; respect during archive
- [ ] Update `skills/spectacular/SKILL.md` — small section on the boundary + handoff
- [ ] Update `.spectacular/AGENTS.md` — boundary section: spectacular owns `.spectacular/`, pageworks owns `docs/`
- [ ] Update `CLAUDE.md` (root) — Active Requests table; add a "Skill boundary" subsection

## M6 — Tests + release

**pageworks tests**

- [ ] `pageworks/tests/run.sh` — test runner (mirror spectacular's pattern)
- [ ] `pageworks/tests/cli/init.test.sh` — port from `spectacular/tests/cli/docs.test.sh` scenarios 1-3, 10
- [ ] `pageworks/tests/cli/export.test.sh` — port from `spectacular/tests/cli/docs-export.test.sh` (all 16 scenarios)
- [ ] `pageworks/tests/cli/doctor.test.sh` — port docs-related scenarios from `spectacular/tests/cli/docs.test.sh` (scenarios 4-9)
- [ ] Run all pageworks tests; confirm 100% pass

**spectacular tests**

- [ ] `spectacular/tests/cli/docs-deprecation.test.sh` — new: each docs verb prints deprecation banner; behavior otherwise unchanged
- [ ] Update `spectacular/tests/cli/docs.test.sh` — scenarios that validate now-removed schema checks become discovery-only assertions; or move to pageworks's port
- [ ] Update `spectacular/tests/cli/doctor.test.sh` — `doctor docs` only emits discovery-level info now
- [ ] Run full spectacular test suite; confirm no regressions

**Releases**

- [ ] Bump spectacular: `SPECTACULAR_VERSION` to `1.2.0`, both plugin manifests, SKILL.md frontmatter
- [ ] CHANGELOG entry under spectacular `[1.2.0]` — Deprecated: docs surface (moved to pageworks); Changed: doctor docs now discovery-only
- [ ] Copy spectacular CLI to `~/.local/bin/spectacular`
- [ ] Verify `pageworks --version` → `0.1.0`
- [ ] Tag spectacular v1.2.0, push, `gh release create v1.2.0 --generate-notes`
- [ ] Tag pageworks v0.1.0 (from inside its own repo if separated, or as a path-prefixed tag if monorepo-style)
- [ ] `gh release create v0.1.0 --generate-notes` for pageworks (if it has its own GitHub repo)
- [ ] Update CLAUDE.md Active Requests table — move pageworks-migration to Archived line
- [ ] Snapshot PLAN + TASKS, archive request → `.spectacular/archive/pageworks-migration/`
- [ ] `/plugin marketplace update spectacular`
- [ ] `/plugin marketplace add pageworks` (if applicable)

## Open questions to resolve during execution

- [ ] Does pageworks ship as its own GitHub repo, or live inside the spectacular monorepo as a sub-path? (Affects tag scheme + release tooling.) → **Decide at M1**.
- [ ] Does pageworks's `.spectacular/` itself need an `archive/` from day one or only when its first request lands? → **Lazy: create when needed**.
- [ ] Should `pageworks doctor` also do drift detection (spec mtime vs page `updated:`) in v0.1.0, or defer to a `pageworks-maintenance` follow-on? → **Defer; v0.1.0 ports existing checks only**.
