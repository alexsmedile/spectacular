---
status: planned
priority: high
owner: alex
updated: 2026-07-07
build: b23
summary: "Restructure the .spectacular/ workspace to align with Open Knowledge Format (OKF) layout, rename folders to plural, move indexes inside category folders, and update the CLI, Skill, and migrations."
related:
  - ../../PRD.md
  - ../../ARCHITECTURE.md
  - ../../PRINCIPLES.md
---

# Plan — okf-restructure

## 1. Goal

Restructure the `.spectacular/` workspace to enforce and support the Open Knowledge Format (OKF) layout: move collection index files inside their respective category folders as `index.md`, flatten nested specification folders under a single plural `specs/` directory, standardise on plural directory naming with prefix identifiers (`decisions/D1-<slug>`, `memories/M1-<slug>`), formalise the concept of "project-wide anchors," and update the CLI and Skill to support these conventions.

## 2. Constraints

- **Backwards Compatibility:** Existing workspaces must be able to migrate cleanly from schema v0.6 to v2.0 without data loss via `spectacular migrate` or `spectacular doctor --fix`.
- **Minimal Root Footprint:** The root `.spectacular/` folder must contain only project-wide anchors (`PRD.md`, `PRINCIPLES.md`, `POLICY.md`, `ARCHITECTURE.md`, `STACK.md`, `PERSONAS.md`, `config.yaml`) and category directories.
- **Strict Pluralization:** All category directory names must be plural (`specs/`, `decisions/`, `memories/`, `sessions/`, `roadmaps/`, `feedbacks/`, `audits/`, `fixes/`, `debugs/`, `requests/`, `archive/`, `ideas/`).
- **No Overwriting during Init:** Restructuring must be idempotent and non-destructive.
- **Prefix Consistency:** Entries under `decisions/` and `memories/` must follow the sequential prefix pattern (`D<N>-<slug>.md` and `M<N>-<slug>.md`).

## Understanding

### How it works now

- Collection index files live at the root of `.spectacular/` (e.g. `decisions/index.md`, `memories/index.md`, `sessions/index.md`, `roadmaps/index.md`, `specs/index.md`).
- Collection entries live in folders at the root (e.g. `decisions/`, `memory/`, `sessions/`, `roadmap/`, `specs/`).
- Directory naming is a mix of singular (`memory/`, `roadmap/`, `feedback/`, `audit/`, `debug/`) and plural (`decisions/`, `sessions/`, `specs/`, `ideas/`, `fixes/`).
- Capability specs are nested under subdirectories (e.g. `specs/doc-engine/specs/index.md`, `specs/roadmap/specs/index.md`), rather than flat files.
- The concept of "project-wide anchors" is used implicitly but not formalised in the documentation or system layout rules.

### What changes

- **Index File Re-location:** Move index files inside their folders as `index.md` (e.g., `decisions/index.md`, `memories/index.md`, `specs/index.md`, `roadmaps/index.md`).
- **Specification Consolidation:** Eliminate the `capabilities/` folder and merge all specifications into a flat `specs/` directory (e.g., `specs/cli.md`, `specs/skill.md`, `specs/doc-engine.md`, `specs/roadmap.md`). No nested folders.
- **Strict Pluralization:** Rename `memory/` to `memories/`, `roadmap/` to `roadmaps/`, `feedback/` to `feedbacks/`, `audit/` to `audits/`, and `debug/` to `debugs/`.
- **Sequential Prefixes:** Entries in `decisions/` must use the format `D<N>-<slug>.md` (e.g., `D1-mode-count.md`), and entries in `memories/` must use `M<N>-<slug>.md` (e.g. `M1-v1-5-0-substrate.md`).
- **CLI & Skill Updates:** Update `cli/spectacular` (validation checks, doctor areas, file creators, path constants) and skill reference documents (briefings, new-request, archive, lifecycle) to point to the new paths and enforce the new naming patterns.
- **Migration Path:** Create a v0.6-to-v2.0 migration script to mechanically migrate old layouts.
- **Terminology:** Update `ARCHITECTURE.md` and related docs to define and document the top-level files as "project-wide anchors."

### What stays the same

- The core CLI binary location (`cli/spectacular`), installer (`cli/install.sh`), and repository hooks.
- Markdown frontmatter structures, types, and operational semantics (excluding index file paths).
- The five request lifecycle states (`planned`, `active`, `review`, `verified`, `archived`).

## Decisions

- **Decision 1:** Choose `specs/` as the consolidated directory for capability files instead of `capabilities/` because `specs` is the established term in the Spectacular codebase and matches the existing `specs/index.md` index nomenclature.
- **Decision 2:** Use flat markdown files (e.g., `specs/doc-engine.md`) rather than subfolders (e.g. `specs/doc-engine/index.md`) for capabilities because flat files prevent folder proliferation, simplify search, and mirror the standard soft-DB collection entry structure.

## 3. Milestones

- **M1 — Architecture and Schema Definition:** Formalise "project-wide anchors" and the OKF layout in `ARCHITECTURE.md`, `doc-index.md`, and update the templates.
- **M2 — CLI Restructuring:** Update `cli/spectacular` constants, help menus, `top_usage()`, file-creating logic, path-resolution functions, and `doctor` validators to handle the OKF plural folders and sequential prefixes.
- **M3 — Skill Restructuring:** Update all skill references (`skills/spectacular/references/`) to use the new folder/file paths and naming formats.
- **M4 — Workspace Schema Migration:** Implement the mechanical `v0.6-to-v2.0` migration step, validating that running `spectacular migrate` or `spectacular doctor --fix` upgrades an existing workspace without losing any files or history.
- **M5 — Verification & E2E Validation:** Execute tests across all CLI options, verify that no errors or warnings are flagged by `doctor`, and ensure that git hooks remain green.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None (this is a standalone architectural migration).

## 6. Validation

- **M1 — Architecture Definition:**
  * Assertable: `ARCHITECTURE.md` defines "project-wide anchors" and outlines the OKF plural directories.
  * Assertable: `skills/spectacular/references/doc-index.md` registers the new `index.md` locations and the `specs/` consolidation.
- **M2 — CLI Restructuring:**
  * Run: `bash -n cli/spectacular` exits 0.
  * Run: `./cli/spectacular paths` returns the updated plural directories and index file paths.
- **M3 — Skill Restructuring:**
  * Assertable: `grep -r "decisions/index.md" skills/spectacular/references/` returns no matches (they should all point to `decisions/index.md`).
- **M4 — Workspace Schema Migration:**
  * Run: Create a mock v0.6 workspace, run `./cli/spectacular migrate`, and assert that directories are renamed, indices are moved, and `workspace_schema` is bumped to `2.0`.
- **M5 — Verification & E2E Validation:**
  * Run: `spectacular doctor` on the migrated repository exits 0 with no errors.
  * Run: `bash tests/run.sh` passes all scenarios.

## 7. Deliverables

- Updated [cli/spectacular](../../../cli/spectacular)
- Updated [ARCHITECTURE.md](../../ARCHITECTURE.md) and related docs defining project-wide anchors
- Relocated index files and renamed plural folders in the spectacular workspace
- A validated `v0.6-to-v2.0.md` migration script
- Green test suite
