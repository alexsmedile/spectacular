---
status: planned
updated: 2026-05-23
related:
  - PLAN.md
  - ../convention-pack-schema/PLAN.md
---

# Tasks — Convention Pack Fabricator

## v1

### M1 — Pack registry hooked up
- [ ] Update `doc-registry.md`'s convention-pack entry with `overrides: references/pack-overrides.md`
- [ ] Add SKILL.md routing: `spectacular pack new <name>` / `pack grill <name>` / `pack refine <name>` / `pack review <name>`
- [ ] Update SKILL.md References + Templates indexes

### M2 — Slot prompts (pack-overrides.md)
- [ ] Write `pack-overrides.md` mirroring `prd-overrides.md` structure
- [ ] Slot 1: Name + scope + applies-to (what project types?)
- [ ] Slot 2: Naming rules (folder case, file case, role suffixes, forbidden words)
- [ ] Slot 3: Top-level taxonomy (required folders, opt-in folders, mono-collection vs project root detection)
- [ ] Slot 4: Root files (README contract, AGENTS pattern, LICENSE, CHANGELOG)
- [ ] Slot 5: Gitignore defaults (always-add, opt-in, never-auto-add)
- [ ] Slot 6: File placement rules (where new files of each kind go)
- [ ] Slot 7: Project types (which `--type` values the pack supports; templates pointer)
- [ ] Mini-refine patterns for each slot

### M3 — Grill flow tested cold
- [ ] Run `spectacular pack new throwaway-test` against a clean state
- [ ] Confirm pack.md frontmatter writes correctly per slot answer
- [ ] Confirm templates/ + references/ folders scaffold (empty initially; user fills in via templates slot)
- [ ] Confirm review gate flags missing required categories

### M4 — Source-ingestion mode
- [ ] `--from <path1>,<path2>` flag accepted by `pack new`
- [ ] Skill reads source files, pre-fills the most-confident slots
- [ ] User reviews each pre-filled slot during grill (default y, edit-on-n)
- [ ] Document supported source types in pack-overrides.md (gitignore file, NAMING_RULES.md, README.md, observed folder structure)

### M5 — alex-default dogfood
- [ ] Run `spectacular pack new alex-default --from ~/code/NAMING_RULES.md,~/code/README.md,~/.claude/CLAUDE.md`
- [ ] Grill through all 6 categories
- [ ] Produce `~/.spectacular/packs/alex-default/`
- [ ] Copy/move to `<repo>/packs/alex-default/` for distribution
- [ ] Verify the produced pack covers all 10 sections from archived repo-conventions PLAN
- [ ] Commit to repo

### M6 — Review gate
- [ ] `spectacular pack review` checks: required slots filled, naming rules consistent (no conflicts), at least one project type declared, gitignore defaults non-empty
- [ ] Doctor extends `kits` check area to also validate packs in `~/.spectacular/packs/` and `<repo>/packs/` — OR a new `packs` doctor area (decide during implementation)

### M7 — Tests + VERIFY.md
- [ ] Create `tests/cli/pack.test.sh` — scenarios: new pack scaffolds correctly, review flags incomplete, --from pre-fills, alex-default smoke test
- [ ] Create `requests/convention-pack-fabricator/VERIFY.md` per 2-of-6 (scored 3)
- [ ] Manual scenarios: full interactive grill walkthrough, source-ingestion accuracy spot-check, dogfood validation

### Validation (folded)
- [ ] All M tasks above
- [ ] alex-default pack present at `<repo>/packs/alex-default/`
- [ ] Smart-init + doctor + pack tests all green together
- [ ] Schema (request 1) confirmed sufficient — any schema gaps surfaced and either patched or documented

## v2 (deferred)

- [ ] `spectacular pack list` / `pack install` / `pack remove` CLI commands — owned by [[convention-pack-application]]
- [ ] Pack diff/merge
- [ ] Multi-pack composition
- [ ] Auto-fabricate from a fully-formed repo (analyze + propose)
