---
status: verified
updated: 2026-05-22
related:
  - PLAN.md
  - VERIFY.md
---

# Tasks — Smart Init

## v1

### M1 — Always-set defined
- [x] Document the 5-file always-set (`PRD.md`, `requests/`, `current/`, `config.yaml`, `AGENTS.md`) in `references/init-workflow.md`
- [x] Add rationale section: why these 5 and not the others
- [x] Update `.spectacular/ARCHITECTURE.md § Init flow` to match
- [x] Snapshot `ARCHITECTURE.md` before edit (saved at `ARCHITECTURE@v1.0.md` during verification-convention work; v1.1 increments)

### M2 — Flag interface
- [x] Add `--kit <name>` flag parsing to `cli/spectacular`
- [x] Add `--with <doc1,doc2,...>` flag parsing
- [x] Add `--minimal` flag (overrides kit defaults, scaffolds always-set only)
- [x] Update `--help` output
- [x] Error cleanly on invalid combinations (unknown kit, unknown doc ID in --with)

### M3 — Pre-flight check
- [x] Implement file-exists check before every write
- [x] Distinguish empty-file (fill) vs non-empty (skip) cases — new `is_empty_file()` helper
- [x] Stdout reporting per the 9-state table in PLAN.md § Pre-flight behavior (✓ created / ✓ filled / ⊘ skipped / ⊘ diagnostics)
- [x] Idempotency test: re-run on existing workspace exits 0 with all-skip report (scenario 4)
- [x] `.gitignore` append-only logic (existing behavior preserved)
- [x] Generic "run diagnostics via `spectacular doctor` once available" message for malformed cases
- [x] Confirmed `--force` flag is NOT added (decision: violates non-destructive principle)
- [x] "Add kit later" flow tested via re-run with `--kit` flag (scenario 4 baseline)

### M4 — Interactive mode
- [x] Refactor `-i` flow: ask for kit (numbered menu 1-5)
- [x] Auto-scaffold kit's `triggers-docs.always` list (via `resolve_doc_set()`)
- [x] Prompt y/n for each entry in kit's `triggers-docs.suggested` list (default: y)
- [x] Bare `init` (no -i, no --kit) uses blank kit unconditionally — no inference (scenario 1)

### M5 — Kit consumption
- [x] Parse selected kit's frontmatter from `templates/prd/kits/<name>.md` (`kit_triggers_always/suggested` awk parsers)
- [x] Resolve `triggers-docs.always` to scaffold list
- [x] Resolve `triggers-docs.suggested` to prompt list (interactive) or noop (non-interactive)
- [x] Fall back gracefully if kit metadata missing (returns empty)
- [x] Error on unknown doc IDs (parse_args validates before scaffold)
- [x] Project-local override (`.spectacular/templates/prd/kits/<id>.md`) wins over installed via `find_kit_file()`

### M6 — Tests + VERIFY.md
- [x] Create `tests/run.sh` harness — POSIX-portable, bash 3.2 compatible
- [x] Create `tests/cli/init.test.sh` — covers all 6 VERIFY scenarios
- [x] Each test creates isolated `/tmp/spectacular-test-<n>/`, runs CLI, asserts, cleans up
- [x] Create `requests/smart-init/VERIFY.md` from scaffold-reference.md stub
- [x] Populate VERIFY.md with 6 manual QA scenarios mirroring the test cases
- [x] Tests bypass GitHub fetch via local symlink seed — offline + fast

### M7 — Dogfood
- [x] Run `tests/run.sh`; all 50 asserts pass across 8 scenarios (6 core + abort + --global)
- [x] CLI smoke-tested against `/tmp/v[1-7]/` during review walkthrough
- [x] Manual walkthrough of VERIFY.md checklist (all items ticked; abort + --global exercised)
- [x] CHANGELOG.md v0.3.0 entry written

### Verification (folded into TASKS per [[verification]] — procedural complement to VERIFY.md)
- [x] All M6 test scenarios pass (`tests/run.sh` exit code 0)
- [x] All VERIFY.md scenarios manually walked (8 scenarios + 3 edge cases + 3 regression)
- [x] PLAN § Validation criteria each confirmed
- [x] CHANGELOG.md v0.3.0 entry written

### Post-review polish (delivered 2026-05-22)
- [x] Interactive final summary + abort prompt (`Proceed? [Y/n]:`) — gap 1 closed
- [x] `--global` flag exercised under fake HOME (no `~/` pollution) — gap 2 closed
- [x] Added test scenarios 7 (abort) + 8 (--global) to automated suite — permanent coverage

## v2 (deferred)

- [ ] Auto-detect project type from repo signals (`package.json`, `SKILL.md`, `.claude-plugin/`, `pyproject.toml`)
- [ ] `spectacular scaffold <type>` retrofit command for existing repos
- [ ] Project-type → kit-suggestion mapping
- [ ] Multi-kit application at init time
- [ ] Replace generic "run diagnostics" message with explicit doctor references once doctor ships
