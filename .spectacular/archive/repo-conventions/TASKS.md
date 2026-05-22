---
status: planned
updated: 2026-05-21
related:
  - PLAN.md
---

# Tasks — Repo Conventions

## v1

### Encode conventions
- [ ] Write `skills/spectacular/references/repo-layout.md` covering:
  - [ ] Naming rules (kebab-case, anchor-descriptor-role, fixed role suffix set)
  - [ ] Top-level taxonomy for mono-collections
  - [ ] Per-project standard folders
  - [ ] Per-project root files
  - [ ] README contract (Type/Stack/Run header)
  - [ ] AGENTS.md + CLAUDE.md pairing pattern
  - [ ] `.gitignore` defaults (with tool-generated dirs as opt-in)
  - [ ] File placement rules
  - [ ] Project-type templates table
  - [ ] Mono-collection vs project-root detection heuristic

### Templates
- [ ] `templates/repo/.gitignore.tpl` — default gitignore
- [ ] `templates/repo/README.tpl.md` — README contract template
- [ ] `templates/repo/cli/` — minimal CLI scaffold
- [ ] `templates/repo/library/` — minimal library scaffold
- [ ] `templates/repo/webapp/` — minimal webapp scaffold
- [ ] `templates/repo/skill/` — skill scaffold (SKILL.md + folders)
- [ ] `templates/repo/plugin/` — plugin scaffold (.claude-plugin/, skills/, agents/, commands/)
- [ ] `templates/repo/content/` — content project scaffold
- [ ] `templates/repo/research/` — research project scaffold
- [ ] `templates/repo/vault-project/` — Obsidian-style project scaffold

### Wire into skill flows
- [ ] Update `references/init-workflow.md` — load repo-layout, ask for project type, scaffold
- [ ] Update `references/new-request.md` — apply file-placement rules for artifacts
- [ ] Update `SKILL.md` templates index to include `templates/repo/`

### Dogfood
- [ ] Test init on a fresh empty repo (each project type)
- [ ] Test `spectacular new test-slug` with research/screenshot artifacts
- [ ] Test mono-collection detection on `~/code` (should skip top-level scaffold)
- [ ] Verify `.gitignore` defaults don't auto-add tool dirs without asking

## v2 (deferred)

- [ ] `spectacular scaffold <type>` — retrofit existing repos
- [ ] Auto-detect project type from `package.json` / `pyproject.toml` / file presence
- [ ] `spectacular doctor` / `spectacular check` — lint existing repo against conventions
- [ ] `spectacular migrate` — rename `_archived/` → `_archive/`, `_backup/` → `_backups/`
- [ ] Deeper per-language conventions (Python packaging, Node monorepo)
- [ ] `.spectacular/CONVENTIONS.md` — per-project convention overrides
