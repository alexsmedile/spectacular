---
status: planned
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Convention Pack Modules

## v1 (= v2 of pack system)

### M1 — Module schema spec
- [ ] Extend `references/packs-contract.md` with § "Modular packs (v2)"
- [ ] Document `<pack-name>/modules/<category>.yaml` file shape (per-category schema body, no manifest fields)
- [ ] Document thin `pack.md` format: name + version + description + `modules: [list]`
- [ ] Document backwards-compat fallback: if `modules:` absent, parser reads inline `rules:` block (v1 path)
- [ ] Add 2 worked examples: thin pack.md + sample naming.yaml module

### M2 — Engine read path (CLI)
- [ ] Add `pack_load_modules()` helper in CLI — checks for `modules/` folder first, falls back to inline `rules:`
- [ ] Refactor `pack_gitignore_always_add()` (and any future per-rule helpers) to use `pack_load_modules()`
- [ ] Both v1 and v2 packs resolve transparently — no caller code changes needed

### M3 — alex-default migration (dogfood)
- [ ] Snapshot current `packs/alex-default/pack.md` to `packs/alex-default/pack@v1.md`
- [ ] Split current pack into `modules/naming.yaml`, `modules/taxonomy.yaml`, etc.
- [ ] Rewrite `pack.md` as thin manifest with `modules: [naming, taxonomy, root-files, gitignore, file-placement, project-types]`
- [ ] Regression test: existing pack.test.sh scenarios still pass against the migrated pack

### M4 — CLI module verbs
- [ ] `spectacular pack module list <pack>` — show modules declared by a pack
- [ ] `spectacular pack module show <pack> <module>` — print one module's contents
- [ ] `spectacular pack install <pack> --modules naming,gitignore` — install only declared modules (resulting user-scope pack has thin manifest pointing only at installed modules)
- [ ] Help text + error handling (unknown module, --modules on v1 monolithic pack)

### M5 — config.yaml composition
- [ ] Document `convention_pack.compose: [{ source: <pack>, modules: [<list>] }]` syntax in ARCHITECTURE.md
- [ ] CLI parser: `config_pack_compose()` returns list of (source, modules) tuples
- [ ] Conflict resolution: for each rule category, walk compose list in order; later entry wins for that category
- [ ] init + doctor consume composed pack — same surface as single-source pack, transparent merging

### M6 — Doctor extension
- [ ] `check_conventions()` validates per-module file existence when `modules:` declared
- [ ] Flags missing module files as error
- [ ] Validates that compose entries reference installable packs (resolves via scope precedence)

### M7 — Tests + VERIFY.md
- [ ] Regression: all 12 v1 pack.test.sh scenarios still pass against migrated alex-default
- [ ] New scenario: install pack with `--modules` subset → user pack has only requested modules
- [ ] New scenario: compose 2 packs → doctor reports merged rule set correctly
- [ ] New scenario: compose conflict → later entry's module wins
- [ ] New scenario: v1 pack.md (no modules:) still loads via fallback
- [ ] Create `requests/convention-pack-modules/VERIFY.md`

### Release
- [ ] CHANGELOG.md v0.5.0 entry — modular packs
- [ ] README.md updated — modules section + compose example
- [ ] docs/configuration.md — `compose:` syntax + conflict resolution

## v2 (deferred — further out)

- [ ] Per-module versioning (modules.<cat>.version)
- [ ] Module marketplace surface
- [ ] Module signing / verification
- [ ] Cross-pack module dependencies
- [ ] `pack split <pack>` fabricator verb — interactively decompose a v1 monolithic pack
