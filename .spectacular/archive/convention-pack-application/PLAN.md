---
status: archived
shipped_in: v0.4.0
archived: 2026-05-23
priority: medium
owner: alex
updated: 2026-05-23
summary: "Wire packs into init/new-request/doctor; add pack install/list/remove CLI commands; per-repo mode (suggest/scaffold/enforce)"
related:
  - ../convention-pack-fabricator/PLAN.md
  - ../../archive/doctor/PLAN.md
---

# Plan — Convention Pack Application

## Goal

Make packs *do something*. Wire the active pack into `spectacular init`, `spectacular new`, and `spectacular doctor`. Add CLI commands for pack lifecycle (`install`, `list`, `remove`). Per-repo selection lives in `config.yaml` with three application modes: `suggest`, `scaffold`, `enforce`.

## Why

The schema (request 1) defines packs; the fabricator (request 2) produces them. This request connects packs to runtime — without it, packs are inert files. This is also where the "app store" install flow lands: users discover packs in `<repo>/packs/`, install with one command, and start using them.

## Scope

**In scope (v1)**
- CLI: `spectacular pack install <name>` — fetch from `<github-repo>/packs/<name>/` to `~/.spectacular/packs/<name>/`
- CLI: `spectacular pack list` — show installed packs (bundled + user + project-local)
- CLI: `spectacular pack remove <name>` — uninstall user-scope pack (never touches bundled or project-local without --force)
- `config.yaml` schema extension: `convention_pack: { source, mode, overrides }`
- Mode: `suggest` — skill mentions pack opinions but never blocks or scaffolds automatically
- Mode: `scaffold` — init/new-request scaffold per pack rules; doctor adds info-level notes
- Mode: `enforce` — doctor adds `conventions` check area; warnings/errors block lifecycle on drift
- Init wiring: if `convention_pack:` in config.yaml, init consults pack for additional scaffold beyond always-set
- New-request wiring: pack's file-placement rules drive artifact directory creation
- Doctor wiring: new `conventions` check area when mode=enforce
- Project-local override: `<project>/.spectacular/packs/<name>/` shadows user-scope same-name pack

**Out of scope (v2)**
- Pack composition (multi-pack stack with precedence rules)
- Pack auto-update (re-fetch from GitHub on schedule)
- `spectacular pack publish` — upload a user pack to the app-store folder via PR or API
- Pack signing/verification
- Migration command — apply a pack's rules retroactively to an existing repo

**Explicit anti-patterns**
- Auto-installing packs without user request — `install` is always explicit
- Enforce mode blocking with no escape — every enforce-failure surfaces an override path (skip-once flag, downgrade to warning)
- Hardcoded pack content in CLI — packs are read at runtime from one of the four scope locations
- Scaffolding the same file twice when pack + always-set overlap — always-set wins; pack only adds what's not already in always-set

## Verification routing

2-of-6 rule applied:
1. User-visible change — ✓ CLI gains new subcommands + init/new-request/doctor behavior changes
2. Reversibility cost — ⚠️ partial (per-file scaffold is reversible; mode changes are config-only)
3. Multi-surface verification — ✓ CLI (install/list/remove) + init + new-request + doctor + 3 modes
4. Risk surface — ⚠️ partial (modifies workspace structure based on pack content — packs are user-trusted)
5. External contract change — ✓ new CLI surface + config.yaml schema addition
6. Rollback — ⚠️ partial (per-fix snapshot per existing rules; mode changes are config-only revert)

**Score: 4 of 6** → VERIFY.md required. Comprehensive — touches CLI + skill + multiple workflows.

## Milestones

1. **`pack install` / `pack list` / `pack remove`** — CLI lifecycle for packs; install fetches from GitHub tarball (mirror skill install logic)
2. **config.yaml schema** — `convention_pack:` block parseable by CLI; documented
3. **Init wiring** — if pack declared, scaffold its `applies-to` types alongside always-set
4. **new-request wiring** — pack's file-placement rules drive artifact dirs
5. **Doctor `conventions` area** — only active when mode=enforce; flags drift from pack rules
6. **Three modes proved** — suggest / scaffold / enforce each have a working flow
7. **Tests + VERIFY.md** — comprehensive coverage given the surface

## Tasks

See `TASKS.md`.

## Dependencies

- **Hard dependency on [[convention-pack-schema]]** — needs the schema to read packs
- **Hard dependency on [[convention-pack-fabricator]]** — needs at least one real pack (alex-default) to install/apply during dogfood
- **Touches [[doctor]]** — adds `conventions` check area
- **Touches [[smart-init]]** — extends init's resolver to consult active pack
- **Touches [[cli-bootstrap]]** — install verb mirrors skill install logic

## Deliverables

- Updated `cli/spectacular` — `doctor` subcommand sibling: `pack` subcommand with install/list/remove
- Updated `cli/spectacular` — init resolver consults active pack
- Updated `references/init-workflow.md` — pack integration documented
- Updated `references/new-request.md` — pack's file-placement rules consulted
- Updated `references/doctor.md` — `conventions` check area
- Updated `.spectacular/ARCHITECTURE.md` — config.yaml schema extension; convention-pack mode docs
- Updated CHANGELOG.md — v0.4.0 entry (pack system lands as a feature release)
- `tests/cli/pack.test.sh` extended with install/list/remove + mode scenarios
- VERIFY.md
