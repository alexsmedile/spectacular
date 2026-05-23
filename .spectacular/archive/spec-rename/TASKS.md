# Tasks — Spec Rename

## M1 — Schema + scaffold

- [ ] Add `spec` doc-type entry to `skills/spectacular/references/doc-registry.md` (template + slots + mode + location)
- [ ] Create `skills/spectacular/templates/SPEC.md` template (index-style: capabilities list, "no capabilities yet" default state)
- [ ] Add `spec-overrides.md` only if grill/refine prompts genuinely differ from PRD — skip if generic engine suffices
- [ ] Update `cli/spectacular` init: scaffold `SPEC.md` (always-set) + `specs/` dir with `.gitkeep`
- [ ] Update `cli/spectacular` init: stop scaffolding `current/`
- [ ] Update always-set count in init logging + help text (5 files → 5 files, just renamed)

## M2 — Doctor migration

- [ ] Add `specs` area to `DOC_AREAS` in `cli/spectacular`
- [ ] Implement `check_specs()`: SPEC.md exists + parseable, specs/ dir exists, per-capability SPEC.md files valid if present
- [ ] Migration detection: `current/` present + no `specs/` → warning + `--fix` available
- [ ] Conflict detection: both `current/` and `specs/` present → error, refuse auto-fix, instruct manual merge
- [ ] Extend `doctor_apply_mechanical_fixes()` with rename handler: `mv current/ specs/`, preserve contents, log change
- [ ] Update `doctor_usage` help text with `specs` area

## M3 — Reference updates

- [ ] `ARCHITECTURE.md` — replace `current/` references with `specs/`; document SPEC.md as canonical index
- [ ] `AGENTS.md` — context-loading table: `current/<capability>` → `specs/<capability>/SPEC.md`
- [ ] `skills/spectacular/SKILL.md` — frontmatter version bump, references index, routing table for `spectacular spec` verbs
- [ ] `skills/spectacular/references/scaffold-reference.md` — SPEC.md frontmatter stub
- [ ] `skills/spectacular/references/doc-registry.md` — new entry confirmed
- [ ] `skills/spectacular/references/lifecycle.md` — any current/ mentions
- [ ] `skills/spectacular/references/active-request.md` — capability spec references
- [ ] `skills/spectacular/references/init-workflow.md` — scaffold list updated
- [ ] `skills/spectacular/references/onboarding.md` — load order references
- [ ] `skills/spectacular/references/doctor.md` — specs area documented
- [ ] `docs/scaffold.md` — full file-tree updated
- [ ] `docs/commands.md` — doctor area list + spec verbs
- [ ] `docs/configuration.md` — any current/ mentions
- [ ] `README.md` — workspace tree snippet, table of canonical docs
- [ ] `CLAUDE.md` (project) — workspace structure section
- [ ] `.spectacular/AGENTS.md` (this repo's copy) — context-loading table

## M4 — Dogfood

- [ ] Migrate `.spectacular/current/` (currently empty) → `.spectacular/specs/` in this repo
- [ ] Author `.spectacular/SPEC.md` for spectacular itself — what's built, what it does, link out to architecture
- [ ] Verify briefing flow reads new layout end-to-end

## M5 — Tests + release

- [ ] Create `tests/cli/specs.test.sh` — 6 scenarios: fresh init, init with kit, init on legacy current/, doctor migration --fix, conflict (both dirs), repeat init non-destructive
- [ ] All existing tests pass (init, doctor, pack)
- [ ] Bump version 0.4.0 → 0.5.0 across 8 sources (plugin.json ×2, marketplace.json ×2, codex-plugin/plugin.json, README badge, SKILL.md frontmatter, CHANGELOG top)
- [ ] CHANGELOG.md v0.5.0 entry — Breaking (current/ removed), Added (SPEC.md, specs/), Changed (doc-registry)
- [ ] Pre-commit hook clean
- [ ] Tag + release notes
