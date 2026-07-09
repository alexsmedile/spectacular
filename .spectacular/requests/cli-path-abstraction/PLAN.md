---
status: planned
priority: medium
owner: alex
updated: 2026-07-07
build: b24
summary: "Introduce a centralized path variable registry at the top of the spectacular CLI to avoid hardcoding .spectacular/ paths in functions."
related:
  - ../../PRD.md
---

# Plan — cli-path-abstraction

## Goal

Introduce a centralized path variable registry at the top of the spectacular CLI to avoid hardcoding `.spectacular/...` paths inside commands and functions.

## Constraints

- **Bash 3.2 Compatibility**: Do not use associative arrays (`declare -A`), as they are not supported in Bash 3.2 (macOS default shell). Use standard uppercase variables instead.
- **Compatibility**: All CLI subcommands and options must function identically; no change to user-facing outputs or exit codes.
- **Test Integrity**: Must pass the entire E2E test suite without any modified assertions or regressions.

## Understanding

### How it works now

Currently, throughout `cli/spectacular`, commands and check functions directly hardcode and duplicate path strings like `.spectacular/roadmaps/index.md`, `.spectacular/decisions/`, `.spectacular/memories/`, etc.

### What changes

Centralized global variables will be defined at the top of `cli/spectacular` (e.g. `SPEC_ROADMAPS_DIR`, `SPEC_DECISIONS_DIR`, `SPEC_MEMORIES_DIR`, etc.) and used inside functions instead of literal strings.

### What stays the same

All logic checking for schemas, bootstrapping files, doctor validations, and snapshots remains exactly the same.

## Decisions

- **Flat variables over associative arrays**: Chose plain variables (e.g. `PATH_MEMORIES_DIR`) over associative arrays (`PATH_MAP[memories]`) because Spectacular requires Bash 3.2 compatibility.

## Milestones

- M1 — Centralized variable declaration & audit pass
- M2 — Migrate CLI commands to path variables
- M3 — Full verification and regression check

## Tasks

See `TASKS.md`.

## Dependencies

None.

## Validation

- M1 — run: `bash -n cli/spectacular` check passes.
- M2 — run: `bash -n cli/spectacular` passes and all paths are fully verified.
- M3 — run: `bash tests/run.sh` passes 100% green.

## Deliverables

- Centralized variable block in [spectacular](../../../cli/spectacular).
