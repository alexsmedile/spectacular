---
status: verified
updated: 2026-05-22
verified_on: 2026-05-22
verified_via: "Manual walkthrough of all 6 core scenarios + 3 edge cases + 2 new scenarios (interactive abort, --global with fake HOME). All checklist items green. Cross-checked with tests/cli/init.test.sh (50/50 asserts pass across 8 scenarios)."
related:
  - PLAN.md
  - TASKS.md
---

# Verify — smart-init

> VERIFY answers "did we build it correctly and safely?" (PLAN § Validation answers "what does each milestone need to satisfy?", TASKS answers "what work needs doing?")

This file is **load-bearing**. Every `- [ ]` blocks the `review → verified` transition. Do not move smart-init to verified until every check is `- [x]`. Per [[verification]], the file is opt-in to scaffold but mandatory once it exists.

Scaffolded 2026-05-22 because smart-init hits 3 of 6 axes (user-visible change, multi-surface verification, external contract change).

## Manual QA checklist (6 core scenarios)

Each test scenario uses a fresh `/tmp/spectacular-test-<n>/` directory created at the start, cleaned up at the end.

### 1. Bare init produces only the always-set
- [x] `cd /tmp/test-1 && spectacular init` exits 0
- [x] Files present: `.spectacular/PRD.md`, `.spectacular/config.yaml`, `.spectacular/AGENTS.md`
- [x] Directories present: `.spectacular/requests/`, `.spectacular/current/`
- [x] Files NOT present: `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, `DECISIONS.md`
- [x] PRD.md frontmatter contains `kit: blank`
- [x] `.gitignore` contains `.spectacular.local/` (appended if existed, created if not)

### 2. `--kit coding` adds STACK + ARCHITECTURE
- [x] `cd /tmp/test-2 && spectacular init --kit coding` exits 0
- [x] Always-set files present (per scenario 1)
- [x] Additional files: `.spectacular/STACK.md`, `.spectacular/ARCHITECTURE.md`
- [x] PRD.md frontmatter contains `kit: coding`
- [x] PRD.md has **base 8 slot headings** (Vision → First milestone). Init does NOT apply kit `adds-slots` (Stack, Interfaces) — that's the doc-writer grill engine's job when the user later runs `spectacular prd grill`. Init only honors `triggers-docs.always`. Documented in `cli/spectacular doc_prd()`.
- [x] Suggested docs (PRINCIPLES, ROADMAP, DECISIONS) NOT scaffolded — non-interactive skips them

### 3. `-i` interactive mode walks kit + suggested prompts
- [x] `cd /tmp/test-3 && spectacular init -i`
- [x] Prompt 1: kit menu (5 options: blank, coding, content, product, research)
- [x] Pick `coding`; STACK + ARCHITECTURE scaffold automatically (always-docs)
- [x] Prompt 2-4: y/n for each suggested doc (PRINCIPLES, ROADMAP, DECISIONS)
- [x] Answer `y/n/y`; verify PRINCIPLES + DECISIONS scaffolded, ROADMAP not
- [x] Final summary shown before any writes; user can abort with `n` (verified scenario 7: `n` → "Aborted. No files written." + exit 0; no files created)

### 4. Idempotent re-run skips existing files
- [x] From end-state of scenario 2 (`/tmp/test-2/`), edit PRD.md to add a content line
- [x] Run `spectacular init` again
- [x] Exit 0; PRD.md content unchanged (skip, never overwrite)
- [x] Stdout shows `⊘ PRD.md already present, leaving alone` for each existing file
- [x] No new files created; no errors

### 5. `--with PRINCIPLES,ROADMAP` adds exactly those two
- [x] `cd /tmp/test-5 && spectacular init --with PRINCIPLES,ROADMAP` exits 0
- [x] Always-set present + PRINCIPLES.md + ROADMAP.md
- [x] No STACK / ARCHITECTURE / DECISIONS (not requested, no kit specified)
- [x] Unknown doc ID in --with errors cleanly (e.g. `--with FOOBAR` → error message + non-zero exit)

### 6. `--minimal` overrides kit defaults
- [x] `cd /tmp/test-6 && spectacular init --kit coding --minimal` exits 0
- [x] Always-set present only; STACK + ARCHITECTURE NOT scaffolded despite kit specifying them as `always`
- [x] PRD.md frontmatter contains `kit: coding` (kit identity preserved even though docs skipped)

## Edge cases to verify

- [x] Run init with an existing non-empty `.gitignore` — entry appended, prior content untouched
- [x] Run init when `.spectacular/PRD.md` exists but is empty (0 bytes) — gets filled with stub
- [x] Run init when `.spectacular/PRD.md` exists with garbage content — skipped with diagnostics-deferral message

## Regression checklist

- [x] `spectacular --help` still works (note: there is no `--version` flag in v0.2.x — only `--help`. Replaced original VERIFY item with the actually-existing flag.)
- [x] `spectacular init --global` exercised under fake HOME (`/tmp/spectacular-test-8-home/`) — skill resolves to `$HOME/.agents/skills/spectacular`, workspace scaffolded in project dir, real `~/` untouched. Covered by automated scenario 8.
- [x] `spectacular init --name foo --summary "bar" --agents-file CLAUDE.md` writes name + summary to config.yaml; CLAUDE.md created instead of AGENTS.md

## Rollback validation

Smart-init has no destructive operations — rollback = reverting the release tag. No data migration to undo.

- [x] Confirmed: no destructive paths in cli/spectacular changes
- [x] If rollback needed: `git revert <smart-init-commit>` returns CLI to v0.2.0 behavior; existing v0.3.0-initialized workspaces continue to work (skill is backwards compatible)
