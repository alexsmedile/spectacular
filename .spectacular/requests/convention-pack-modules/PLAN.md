---
status: planned
priority: low
owner: alex
updated: 2026-05-23
summary: "v2: split monolithic packs into composable rule-category modules (naming/taxonomy/root-files/gitignore/file-placement/project-types as standalone files)"
related:
  - ../../archive/convention-pack-schema/PLAN.md
  - ../../archive/convention-pack-fabricator/PLAN.md
  - ../../archive/convention-pack-application/PLAN.md
---

# Plan — Convention Pack Modules

## Goal

Replace the v1 monolithic `pack.md` (all 6 rule categories inline) with a **modular pack** architecture where each rule category lives in its own loadable module. Enables mix-and-match composition — install alex-default's `naming` module + your team's `gitignore` module + a community `project-types` module without forking the whole pack.

## Why

The v1 schema bundles all 6 rule categories into one `pack.md` manifest. This forces an **all-or-nothing** install decision — a user who loves alex-default's naming but wants different gitignore defaults must either fork the whole pack or override individual rules in config.yaml (which v1 supports only via `overrides:` list, not via mixing module sources).

Modular packs solve three concrete pain points surfaced during v1 build:

1. **Forking tax** — small disagreements with a pack require copy-paste of the entire `pack.md`
2. **Discovery friction** — there's no way to find "just the gitignore rules" or "just the naming conventions"; everything is a whole-pack download
3. **Composition impossible** — a project that wants alex-default's structure + a Python-specific project-types block has no way to express that

The v1 decision was deliberate: ship monolithic to prove the schema, fabricator, and application layer work end-to-end, then design composition from production feedback rather than assumption.

## Scope

**In scope (v1 of this request, = v2 of pack system)**
- New folder layout: `<pack-name>/modules/{naming,taxonomy,root-files,gitignore,file-placement,project-types}.yaml` (one file per rule category)
- `pack.md` becomes a thin manifest: name + version + description + `modules:` list referencing the module files
- Backwards compatibility: v1 monolithic packs continue to work — engine reads inline `rules:` block as v1 fallback when no `modules:` declared
- CLI: `spectacular pack install <name> --modules naming,gitignore` (install only declared modules)
- CLI: `spectacular pack module list <pack>` / `pack module show <pack> <module>`
- config.yaml extension: `convention_pack.modules: [naming, gitignore]` (per-repo module subset)
- Composition: `convention_pack.compose: [{ source: alex-default, modules: [naming] }, { source: python-team, modules: [gitignore, project-types] }]`
- Conflict resolution rules: later compose entry overrides earlier for the same module; explicit module list always wins over fallback
- Doctor: `conventions` area extended to validate module file existence + per-module schema

**Out of scope (v3+)**
- Per-rule-category module versioning (modules inherit parent pack version in v2)
- Module marketplace / discovery surface beyond the existing `<repo>/packs/` app-store
- Module signing
- Cross-pack module dependencies (module A requires module B from pack C)

**Explicit anti-patterns**
- Breaking v1 monolithic packs — every v1 pack must continue working unmodified
- Forcing module split — packs may remain monolithic; modular is opt-in
- Auto-decomposing existing v1 packs — fabricator may offer a "split into modules" verb, but never silently rewrites

## Constraints

- Must preserve `packs-contract.md` schema fields — modules use the same per-category schema, just split across files
- Must work with the existing awk-based YAML parser (no new dependencies)
- Backwards compat: `spectacular pack install alex-default` (v1 monolithic) and `spectacular pack install python-team` (v2 modular) must both work transparently
- No `node_modules`-style nested install — module reuse happens at config.yaml level, not on-disk

## Verification routing

2-of-6 rule applied:
1. User-visible change — ✓ new CLI verbs + config.yaml composition syntax
2. Reversibility cost — ✗ low (additive; v1 packs unchanged)
3. Multi-surface verification — ✓ schema + fabricator + application + doctor all touched
4. Risk surface — ⚠️ partial (rewrites pack loading logic; v1 fallback path is the safety net)
5. External contract change — ✓ new config.yaml composition keys + module CLI verbs
6. Rollback — ✗ trivial (revert + v1 monolithic still works)

**Score: 3 of 6** → VERIFY.md scaffolded.

## Milestones

1. **Module schema spec** — extend `packs-contract.md` with module file shape + `pack.md` thin-manifest format; backwards-compat fallback rule
2. **Engine read path** — pack loader prefers `modules/<cat>.yaml` over inline `rules.<cat>:`; resolves both transparently
3. **`alex-default` migrated** — split alex-default into modules as the dogfood; v1 inline-rules version snapshot preserved as `pack.md@v1.md` for backwards-compat regression test
4. **CLI module verbs** — `pack module list / show`; `pack install --modules <subset>`
5. **config.yaml composition** — `convention_pack.compose:` list; conflict resolution rules (later wins for same module)
6. **Doctor extension** — `conventions` validates per-module schema; flags missing module files
7. **Tests + VERIFY.md** — regression coverage that v1 monolithic packs still work; composition scenarios; module-subset install

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on v1 ship** — modular only makes sense after monolithic proves the schema and surfaces composition pain
- **Touches [[convention-pack-schema]]** — extends the schema with module file shape (additive)
- **Touches [[convention-pack-fabricator]]** — new fabricator mode: produce modular pack OR split-existing
- **Touches [[convention-pack-application]]** — CLI + doctor + init wiring all gain module awareness

## Triggers

This request stays `planned` until at least one of these signals fires:

- 3+ users report wanting to "mix and match" pack pieces in feedback
- A community pack lands in the app-store that's effectively a fork of alex-default with one rule changed (= forking tax is real)
- Internal use surfaces a config.yaml `overrides:` list with 10+ entries (= the workaround is the wrong shape)

Until one of these fires, the v1 monolithic + `overrides:` escape hatch is the simpler design.

## Deliverables

- Updated `skills/spectacular/references/packs-contract.md` — module file shape + thin-manifest format + backwards-compat rule
- Updated `skills/spectacular/references/pack-overrides.md` — fabricator support for modular packs + split-existing verb
- Updated `cli/spectacular` — pack loader prefers modules over inline rules; new `pack module` verbs; config.yaml `compose:` parser
- Migrated `packs/alex-default/` — modules layout + v1 snapshot for regression
- Updated `tests/cli/pack.test.sh` — regression coverage (v1 still works) + new composition scenarios
- `requests/convention-pack-modules/VERIFY.md`
