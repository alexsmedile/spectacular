---
status: active
updated: 2026-05-23
related:
  - PLAN.md
---

# Tasks — Workspace Migrations

Two-stage shipping plan. Stage 1 = v0.6.1 hotfix that unblocks Octopus-shape projects. Stage 2 = full registry framework as v0.6.2.

## Stage 1 — v0.6.1 hotfix (active)

### M1 — `workspace_schema:` field
- [ ] Add `CURRENT_SCHEMA="0.6"` constant near top of `cli/spectacular`
- [ ] `doc_config()` writer includes `workspace_schema:` line in scaffolded config.yaml
- [ ] `init` on existing workspace: backfill missing `workspace_schema:` (only when explicit `--update` or via `doctor --fix`)
- [ ] Tests: scenario in `init.test.sh` covering fresh-init writes the field

### M2 — `spectacular status --against-latest`
- [ ] Flag added to `cmd_status` arg parser
- [ ] Reader: parse `workspace_schema:` from config.yaml (default `"0.4"` if absent)
- [ ] Output: one line with current vs latest + suggested verb
- [ ] Tests: extend `tests/cli/init.test.sh` (status path) or add `tests/cli/status.test.sh`

### M3 — `check_specs` upgrade (flat contract docs + capability subfolders both valid)
- [ ] Modify `check_specs` to enumerate top-level `.md` files in `specs/` (excluding `SPEC.md` which lives one level up)
- [ ] Emit `pass` for each top-level .md as "contract doc"
- [ ] Keep existing per-subfolder validation (`<cap>/SPEC.md`)
- [ ] Update `doctor-areas.md` row for `specs` area to document both shapes
- [ ] Tests: doctor scenario with flat .md only + scenario with mixed
- [ ] Document the convention in [[doctor-areas]]: "contract docs at root; capabilities in subfolders with SPEC.md"

### M4 — v0.6+ scaffold suggestion
- [ ] In `check_workspace` (or `check_specs`), detect missing PRINCIPLES.md + ARCHITECTURE.md + ROADMAP.md
- [ ] If 1+ missing: emit ONE info line summarizing which are missing + the `spectacular init --with ...` command
- [ ] Skip silently if all 3 present
- [ ] Tests: scenario where all 3 missing → info; scenario where 2 missing → info lists those 2; scenario where 0 missing → no output

### M5 — `spectacular migrate` CLI (minimal)
- [ ] New `cmd_migrate` function with subcommands: default (apply), `--dry-run`
- [ ] Parser: `migrate [--dry-run]`
- [ ] Migration list: hardcoded array in CLI (Stage 1 only — 2 entries)
- [ ] Dispatch: for each migration where `workspace_schema:` < migration's `to`, run; bump field on success
- [ ] On no-op: print "workspace is up to date (schema X)" + exit 0
- [ ] `--dry-run`: list planned migrations with descriptions, no writes
- [ ] Stub for judgment migrations: print "this migration requires judgment — run `/spectacular migrate` in your AI agent" + exit non-zero (not used in Stage 1; placeholder for Stage 2)
- [ ] Update `Unknown command` listing in CLI usage to include `migrate`

### M6 — Backfilled migrations
- [ ] **v0.4→0.5:** rename `.spectacular/current/` → `.spectacular/specs/` if current/ present AND specs/ absent. Preserve contents verbatim (flat or subfolder). Bump `workspace_schema` to `"0.5"`. Idempotent.
- [ ] **v0.5→0.6:** ensure `.spectacular/specs/` exists (mkdir + `.gitkeep` if empty). Bump `workspace_schema` to `"0.6"`. Idempotent.
- [ ] Both implemented inline in `cmd_migrate` (no registry yet — that's Stage 2)
- [ ] Edge case: both `current/` AND `specs/` present → migration refuses with clear error pointing at doctor (mirrors existing `check_specs` behavior)
- [ ] Smoke tests: fresh workspace fixture with v0.3-shape (current/ + flat .md inside); migrate → assert specs/ exists, current/ gone, contents preserved, workspace_schema bumped

### M7 — Tests + docs + v0.6.1 release
- [ ] `tests/cli/migrate.test.sh` — list, dry-run, apply, idempotence (run twice → second no-op), broken-state refusal
- [ ] `tests/cli/init.test.sh` — workspace_schema written on init
- [ ] `tests/cli/status.test.sh` or extension — `--against-latest` output
- [ ] `tests/cli/doctor.test.sh` — flat contract docs pass; v0.6+ suggestion info line
- [ ] CHANGELOG entry for v0.6.1
- [ ] Bump `version:` in `.claude-plugin/plugin.json` to 0.6.1
- [ ] SPEC.md: add "Workspace migrations (v0.6.1+)" capability bullet
- [ ] CLAUDE.md: update Active Requests row + Archived list
- [ ] Live workspace bumped: `.spectacular/config.yaml` gets `workspace_schema: "0.6"`
- [ ] Doctor self-check clean

## Stage 2 — v0.6.2 (planned, deferred until Stage 1 ships)

### M8 — Migration registry
- [ ] `skills/spectacular/references/migrations/` dir
- [ ] `migrations-contract.md` — frontmatter schema (`from`, `to`, `description`, `mechanical`, `reversible`, `affects`)
- [ ] Register `migrations-contract` in `doc-registry.md`
- [ ] Move Stage 1's two hardcoded migrations into the registry as `v0.4-to-v0.5.md` + `v0.5-to-v0.6.md`
- [ ] Loader in CLI: scan `references/migrations/*.md`; replace hardcoded array

### M9 — `--to` / `--from` flags
- [ ] `migrate --to <version>` (migrate up to specific schema, default: latest)
- [ ] `migrate --from <version>` (re-run a specific migration, for repair)
- [ ] Chain validation: every `from:` chains to a previous `to:`; no gaps; doctor flags broken chain

### M10 — Skill-side walk for judgment migrations
- [ ] `references/migrate.md` — skill flow for `/spectacular migrate`
- [ ] Mirror `doctor-repair.md` pattern: per-migration y/n/q confirm; snapshot canonical docs before edit
- [ ] SKILL.md routing table: add `/spectacular migrate`

### M11 — Init integration
- [ ] `init --update`: after skill update, check pending migrations; prompt user to run `migrate` (don't auto-run)
- [ ] Optional auto-run flag for unattended setups

### M12 — Chain validation in doctor
- [ ] New `migrations` area OR extension of existing `kits` — validates registry integrity
- [ ] Suggested fix when chain broken: point at maintainer (this is a Spectacular-internal bug, not a user bug)

## Verification (per 2-of-6 rule)

Stage 1 hits: user-visible change (new CLI verbs + status flag), schema change (workspace_schema field), risk surface (structural layout migration), external contract (workspace_schema becomes public field). **4 of 6 → needs VERIFY.md when transitioning to `review`.**

Carry to VERIFY.md when ready:
- [ ] Cold-start: fresh `init` → workspace_schema set, doctor clean
- [ ] Octopus-shape: workspace with flat SCHEMA-*.md in current/ → doctor warns about current/ → migrate → doctor clean (no nag on flat layout)
- [ ] Idempotence: migrate twice → second run is no-op
- [ ] Dry-run: --dry-run shows planned changes, doesn't write
- [ ] Broken state: workspace with both current/ AND specs/ → migrate refuses
- [ ] Status --against-latest: correct line for schema 0.4 / 0.5 / 0.6 workspaces
- [ ] Scaffold suggestion: missing all 3 → one info line; missing 0 → silent
- [ ] Stage 1 ships before any further structural change in Spectacular itself
