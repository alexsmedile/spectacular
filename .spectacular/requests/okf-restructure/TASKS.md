---
status: planned
updated: 2026-07-07
related:
  - PLAN.md
---

# Tasks — okf-restructure

## v1

### M1 — Architecture and Schema Definition
- [ ] Define "project-wide anchors" and the OKF layout in `.spectacular/ARCHITECTURE.md`
- [ ] Update `skills/spectacular/references/doc-index.md` to register `index.md` file paths and the `specs/` consolidation
- [ ] Update workspace base templates (`skills/spectacular/templates/`) to reflect the new layout
- [ ] → check: ARCHITECTURE.md and doc-index.md describe the OKF structure

### M2 — CLI Restructuring
- [ ] Update path constants, array mappings, and file creation helpers in `cli/spectacular`
- [ ] Update `top_usage()` CLI help outputs to reflect plural folder names and indexing
- [ ] Update `spectacular doctor` validation areas to check for plural folders and `index.md` targets
- [ ] Update `doctor` checks to validate sequential prefix names (`D<N>-<slug>` and `M<N>-<slug>`)
- [ ] → check: bash -n cli/spectacular exits 0 and paths command returns new configurations

### M3 — Skill Restructuring
- [ ] Update Skill orchestrator triggers in `skills/spectacular/SKILL.md`
- [ ] Replace all index references in `skills/spectacular/references/` (e.g. `decisions/index.md` ➔ `decisions/index.md`)
- [ ] Update `grill`, `refine`, and `review` logic to load rules and rules overrides from the new paths
- [ ] → check: grep search finds no root-level index files (like decisions/index.md) referenced in skill files

### M4 — Workspace Schema Migration
- [ ] Write the `skills/spectacular/references/migrations/v0.6-to-v2.0.md` migration plan/frontmatter
- [ ] Add the `v0.6-to-v2.0` mechanical migration functions in `cli/spectacular`
- [ ] Implement folder renaming, index moving, and links rewriting inside the migration runner
- [ ] → check: running migrate on a mock workspace upgrades it successfully to v2.0

### M5 — Verification & E2E Validation
- [ ] Run the complete test suite (`bash tests/run.sh`)
- [ ] Verify that running `spectacular doctor` on the migrated repository reports zero errors
- [ ] Confirm git hooks are completely green
- [ ] → check: tests/run.sh exits 0 and pre-commit check passes

## v2 (deferred)
- [ ] Implement compose subcommands for multi-part specifications
- [ ] Add automated migration options for custom external convention packs
