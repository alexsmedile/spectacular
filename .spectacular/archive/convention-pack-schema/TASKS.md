---
status: verified
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Convention Pack Schema

## v1

### M1 — Schema locked
- [x] Draft `references/packs-contract.md` mirroring `kits-contract.md` structure
- [x] Define `pack.md` frontmatter fields: `pack`, `version`, `description`, `extends`, `applies-to`, `rules`, `templates`, `references`
- [x] Document the 6 rule categories: `naming`, `taxonomy`, `root-files`, `gitignore`, `file-placement`, `project-types`
- [x] Each rule category gets a frontmatter example + doctor-enforcement note
- [x] Document the override convention (project-local → user → app-store → bundled precedence; `extends:` declared but resolved in v2)

### M2 — Folder shape ratified
- [x] Document pack folder layout (pack.md + templates/ + references/ + optional scripts/)
- [x] Explain how the folder shape maps to the existing skill convention (kit comparison table)
- [x] Decide naming: `pack.md` (chosen for consistency)

### M3 — `minimal` pack shipped
- [x] Create `skills/spectacular/templates/packs/minimal/`
- [x] Write `pack.md` with frontmatter: pack=minimal, version=1.0, description, applies-to=[any], rules=minimal (root-files + gitignore only)
- [x] Write `templates/.gitignore` — canonical gitignore defaults from archived repo-conventions PLAN
- [x] Write `templates/README.md` — README contract stub (Type/Stack/Run header + What/Setup/Usage sections)
- [x] Write `references/why-minimal.md` — philosophy: opt-in opinions, two essentials only, when to install heavier packs
- [x] Confirm pack parses via awk frontmatter parser (smoke tested)

### M4 — Registry entry
- [x] Add `convention-pack` entry to `doc-registry.md`
- [x] Fields: template=templates/packs/minimal/pack.md, mode=grill, location=~/.spectacular/packs/<name>/pack.md, scope=user, snapshot-on-edit=false, overrides=references/pack-overrides.md (forward-declared), description
- [x] Document the new `scope: user` value in registry schema + field-semantics section
- [x] Confirm registry stays parseable (no breaking changes — only adds a new entry + a new scope value)

### M5 — App-store folder established
- [x] Create `<repo-root>/packs/` directory
- [x] Write `packs/README.md` — explains app-store model, available packs table, install/contribute flow, mode semantics
- [x] Reserve `packs/alex-default/` slot with stub README pointing forward to request 2

### M6 — Schema coverage check
- [x] Walked all 10 convention sections from `archive/repo-conventions/PLAN.md`
- [x] Coverage table in `references/packs-contract.md` § Schema coverage check — 10/10 expressible
- [x] One judgment call captured: convention 6 ("most-specific AGENTS.md wins") is runtime behavior, lives in pack's `references/` narrative, not schema

### Validation (folded — verification score 1/6, no VERIFY.md per [[verification]])
- [x] All 6 milestone validations from PLAN § Validation confirmed
- [x] `packs-contract.md` exists with full schema spec
- [x] Minimal pack parses through awk frontmatter
- [x] `doc-registry.md` has convention-pack entry
- [x] All 10 archived conventions expressible in schema
- [x] `packs/` at repo root with README + alex-default placeholder

## v2 (deferred — handled by other requests)

- [ ] Fabricator skill that produces packs interactively → `convention-pack-fabricator`
- [ ] init / new-request / doctor wiring to consume active pack → `convention-pack-application`
- [ ] Opinionated `alex-default` pack → produced by fabricator dogfood
- [ ] Pack composition (multi-pack per repo) — single-pack-only in v1
- [ ] Auto-detection of project type from `package.json` / `pyproject.toml` / `SKILL.md` etc.
- [ ] `spectacular pack list` / `pack install <name>` / `pack remove <name>` CLI commands
- [ ] `extends:` cross-pack inheritance resolution (parsed-only in v1)
