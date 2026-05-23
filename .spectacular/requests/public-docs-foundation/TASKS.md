# Tasks — Public Docs Foundation

## M1 — Schema

- [ ] Write `skills/spectacular/references/docs-contract.md` — docs.yaml schema, page frontmatter contract, section/page model, audience semantics, examples
- [ ] Add `docs-page` doc-type entry to `skills/spectacular/references/doc-registry.md`
- [ ] Add `docs-manifest` doc-type entry to `skills/spectacular/references/doc-registry.md`
- [ ] Write `skills/spectacular/references/docs-overrides.md` — slot prompts for `docs new`, review gate checks, vibe→spec patterns for page rewrites
- [ ] Create `skills/spectacular/templates/docs/` — `docs.yaml.tmpl`, `index.md.tmpl`, `page.md.tmpl` (frontmatter stub)

## M2 — CLI surface

- [ ] Add `docs` subcommand to `cli/spectacular` (init verb only in M2)
- [ ] `spectacular docs init` — scaffold `docs/docs.yaml` + `docs/index.md` + 3 section dirs with placeholder pages
- [ ] Idempotent: re-run is non-destructive, fills empty files, skips non-empty
- [ ] `--minimal` flag — `docs.yaml` + `index.md` only, no default sections
- [ ] Help text + usage update

## M3 — Skill verbs

- [ ] SKILL.md routing table — add `docs` verbs (init/new/review/status)
- [ ] SKILL.md frontmatter triggers — add `spectacular docs *`
- [ ] Skill: `spectacular docs new <page>` — engine-driven via doc-registry; asks for section, scaffolds page from template, updates docs.yaml
- [ ] Skill: `spectacular docs new --section <name>` — appends section to docs.yaml, scaffolds dir + placeholder page
- [ ] Skill: `spectacular docs review` — quality gate using docs-overrides.md gate checks
- [ ] Skill: `spectacular docs status` — same shape as `/spectacular` briefing but scoped to docs/

## M4 — Doctor integration

- [ ] Add `docs` to `DOC_AREAS` in cli/spectacular
- [ ] Implement `check_docs()` — docs/ exists (skip if not), docs.yaml parseable, every yaml-declared page exists on fs, every fs page declared in yaml (warn on orphan), every page has required frontmatter
- [ ] Severity: missing required frontmatter = warning; broken nav = error; missing audience = error
- [ ] Extend `doctor_apply_mechanical_fixes()` — inject missing frontmatter stubs (audience defaults to `[user]`, status to `draft`, updated to today)
- [ ] Update `doctor_usage` help with `docs` area

## M5 — Dogfood

- [ ] Author `docs/docs.yaml` for this repo — sections: getting-started, guides, reference; pages from existing 5 files mapped appropriately
- [ ] Add frontmatter to `docs/commands.md` (audience: [user, agent], status: stable, etc.)
- [ ] Add frontmatter to `docs/configuration.md`
- [ ] Add frontmatter to `docs/scaffold.md`
- [ ] Add frontmatter to `docs/troubleshooting.md`
- [ ] Add frontmatter to `docs/workflow.md`
- [ ] Create `docs/index.md` — landing page pointing into the three sections
- [ ] Decide: split `docs/configuration.md` (currently mixes user setup + internal schema) — user-facing → `docs/reference/configuration.md`, internal → spec or skill reference. Re-evaluate per page.
- [ ] Run `spectacular docs review` → clean
- [ ] Run `spectacular doctor docs` → clean

## M6 — Tests + release

- [ ] Create `tests/cli/docs.test.sh` — 8 scenarios:
  - fresh init (default 3 sections)
  - init --minimal
  - re-init non-destructive
  - new page assigns to section + updates yaml
  - new --section appends correctly
  - review passes on clean tree
  - review fails on missing frontmatter
  - doctor --fix injects stubs
- [ ] All existing tests pass (init, doctor, pack, specs once spec-rename lands)
- [ ] Version bump (TBD based on ship order with spec-rename)
- [ ] CHANGELOG entry — Added: docs subcommand, docs verbs, docs-contract.md, docs-overrides.md, doctor docs area
- [ ] README.md — add `docs/` section to repo structure tree, link new docs-contract reference
- [ ] CLAUDE.md (project) — update reference docs list, add docs verbs to commands
- [ ] AGENTS.md and .spectacular/AGENTS.md — context-loading table entry for docs work
- [ ] Pre-commit hook clean
- [ ] Tag + release notes
