---
status: planned
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Convention Pack Schema

## v1

### M1 — Schema locked
- [ ] Draft `references/packs-contract.md` mirroring `kits-contract.md` structure
- [ ] Define `pack.md` frontmatter fields: `name`, `version`, `description`, `extends`, `applies-to`, `mode` (suggest|scaffold|enforce hint default), `rules`, `templates`, `references`
- [ ] Document the 6 rule categories: `naming`, `taxonomy`, `root-files`, `gitignore`, `file-placement`, `project-types`
- [ ] Each rule category gets a frontmatter example + a worked use case
- [ ] Document the override convention (how project-local packs extend or override the active pack)

### M2 — Folder shape ratified
- [ ] Document pack folder layout (pack.md + templates/ + references/ + optional scripts/)
- [ ] Explain how the folder shape maps to the existing skill convention
- [ ] Decide naming: `pack.md` (main file) vs `PACK.md` vs `<name>.md` — recommend `pack.md` for consistency with how kits use the kit-id as filename

### M3 — `minimal` pack shipped
- [ ] Create `skills/spectacular/templates/packs/minimal/`
- [ ] Write `pack.md` with frontmatter: name=minimal, version=1.0, description, applies-to=any, rules=minimal (gitignore + README contract only)
- [ ] Write `templates/.gitignore` — the canonical gitignore defaults from the archived repo-conventions PLAN
- [ ] Write `templates/README.md` — the README contract stub (Type/Stack/Run header + sections)
- [ ] Write `references/why-minimal.md` — explains the philosophy: opinions are opt-in; minimal ships the essentials; download alex-default if you want stronger defaults
- [ ] Confirm pack parses via awk frontmatter parser (smoke test)

### M4 — Registry entry
- [ ] Add `convention-pack` entry to `doc-registry.md`
- [ ] Fields: template=templates/packs/minimal/pack.md, mode=grill, location=~/.spectacular/packs/<name>/pack.md, scope=user, snapshot-on-edit=false, overrides=references/pack-overrides.md (placeholder for now), description
- [ ] Document the new `scope: user` value — packs live under $HOME, not per-project (with project-local override layer)
- [ ] Confirm registry stays parseable (no breaking changes to existing entries)

### M5 — App-store folder established
- [ ] Create `<repo-root>/packs/` directory at the spectacular repo root
- [ ] Write `packs/README.md` — explains: this is the "app store"; each subfolder is a downloadable pack; users install via `spectacular pack install <name>` (lands in request 2)
- [ ] Reserve `packs/alex-default/` slot with a stub README pointing forward to request 2's dogfood that will produce it

### M6 — Schema coverage check
- [ ] Walk each of the 10 convention sections from `archive/repo-conventions/PLAN.md`
- [ ] For each section, write a mini example in `references/packs-contract.md` showing how that convention encodes in the schema
- [ ] Flag any section the schema can't express → schema adjustment OR explicit out-of-scope note
- [ ] Capture findings in `references/packs-contract.md` § Schema coverage

### Validation (folded into TASKS per [[verification]] convention)
- [ ] All 6 milestone validations from PLAN § Validation confirmed
- [ ] `packs-contract.md` exists with full schema spec
- [ ] Minimal pack parses through awk frontmatter
- [ ] `doc-registry.md` has convention-pack entry
- [ ] All 10 archived conventions expressible in schema
- [ ] `packs/` at repo root with README

## v2 (deferred — handled by other requests)

- [ ] Fabricator skill that produces packs interactively → `convention-pack-fabricator`
- [ ] init / new-request / doctor wiring to consume active pack → `convention-pack-application`
- [ ] Opinionated `alex-default` pack → produced by fabricator dogfood
- [ ] Pack composition (multi-pack per repo) — single-pack-only in v1
- [ ] Auto-detection of project type from `package.json` / `pyproject.toml` / `SKILL.md` etc.
- [ ] `spectacular pack list` / `pack install <name>` / `pack remove <name>` CLI commands
