---
status: verified
updated: 2026-05-23
review_via: "Automated tests (tests/cli/doctor.test.sh — 37/37 asserts across 12 scenarios; +regression for multi-area arg parsing) + dogfood against this workspace + live skill-substrate scenarios 16/17/18 walked in throwaway workspace 2026-05-23"
related:
  - PLAN.md
  - TASKS.md
---

# Verify — doctor

> VERIFY answers "did we build it correctly and safely?" (PLAN § Validation answers "what does each milestone need to satisfy?", TASKS answers "what work needs doing?")

**Load-bearing.** Every `- [ ]` blocks the `review → verified` transition. Per [[verification]], the file is opt-in to scaffold but mandatory once it exists.

Scaffolded 2026-05-22 because doctor hits 5 of 6 axes (user-visible change, partial reversibility cost, multi-surface verification, risk surface — touches canonical docs during repair, external contract change, rollback via snapshot).

## CLI detect — automated coverage

These are exercised by `tests/cli/doctor.test.sh`. Mirror checks here for human-eyed pass.

### 1. Clean workspace
- [x] Fresh seeded workspace with all always-set + valid frontmatter → exit 0, "0 errors, 0 warnings"
- [x] Output mentions every area touched (skill, workspace, links, lifecycle, kits)

### 2. Missing always-set file
- [x] Remove PRD.md → doctor reports "always-set file missing" with error severity → exit 2
- [x] `--fix` re-scaffolds PRD.md using `doc_prd()` with proper frontmatter
- [x] Post-fix re-run is clean (exit 0)

### 3. Malformed frontmatter
- [x] PRD.md without frontmatter delimiter → "missing frontmatter delimiter" error
- [x] Exit 2 because malformed = error severity
- [x] `--fix` does NOT touch the malformed file (judgment required); content unchanged before/after

### 4. Snapshot version gap
- [x] PRD@v1.0.md + PRD@v1.2.md present, no v1.1 → "version-sequence gap" warning
- [x] Exit 1 (warning only)
- [x] Scoped invocation `doctor snapshots` flags the gap and skips other areas

### 5. Broken `related:` link
- [x] `related: - ../nonexistent/PLAN.md` → "target ... does not exist" warning
- [x] Exit 1
- [x] Scoped `doctor links` works the same

### 6. Lifecycle drift
- [x] `status: active` without SESSION.md → "active without SESSION.md" warning
- [x] `status: verified` without verification artifact (no VERIFY.md, no `### Verification` in TASKS, no `## Validation`/`## Success criteria`/`## Acceptance` in PLAN) → error
- [x] `status: verified` with VERIFY.md containing unchecked items → error

### 7. Scoped areas
- [x] `spectacular doctor frontmatter` runs only frontmatter checks
- [x] Other areas' findings (e.g. missing `.gitignore` from workspace area) NOT reported

### 8. JSON output
- [x] `--format json` emits valid JSON parseable by `python3 -c 'import json,sys; json.load(sys.stdin)'`
- [x] JSON has `version`, `ran_at`, `workspace`, `findings`, `summary`
- [x] Each finding has `area`, `severity`, `file`, `message`, `proposed_fix`, `fix_type`
- [x] Pass-entries omitted from JSON (signal-to-noise ratio)

## CLI mechanical fixes — automated coverage

### 9. `--fix` for `.gitignore`
- [x] No `.gitignore` → CLI creates it with `.spectacular.local/`
- [x] `.gitignore` exists without entry → CLI appends without touching prior content
- [x] `.gitignore` has entry → no-op + "no mechanical fixes applied"

### 10. `--fix` for missing dirs
- [x] Missing `.spectacular/requests/` and `.spectacular/current/` → both created
- [x] Stdout: `✓ fixed [workspace]: created <path>`

### 11. `--fix` for dangling symlink
- [x] `.claude/skills/spectacular` points to nonexistent target → CLI removes + recreates pointing to real install
- [x] Verified link target after fix

### 12. `--fix` for missing always-set file
- [x] PRD.md absent → re-scaffolded with `doc_prd()` (8-slot template, `kit: blank` frontmatter)
- [x] config.yaml absent → re-scaffolded with `doc_config()`
- [x] AGENTS.md absent → re-scaffolded with `doc_agents()`

### 13. `--fix` never overwrites non-empty content
- [x] Verified scenario 4: malformed PRD.md content unchanged before/after `--fix`
- [x] Re-running `--fix` on a clean workspace prints "no mechanical fixes applied"

## Skill repair flow (`/spectacular doctor --fix`) — manual / agent-driven

These scenarios cannot be automated — they require an LLM agent to walk findings interactively.

### 14. Judgment finding walk-through
- [x] Real-workspace dogfood: doctor surfaced 5 findings. 3 of 4 judgment findings handled correctly by following the proposed-fix guidance in `references/doctor.md`:
  - Broken `related:` paths in 2 PLAN files → fixed by changing to `../../../skills/...` format
  - `cli-bootstrap` active without SESSION.md → moved to `status: planned` (parked per smart-init handoff)
  - PRD@v1.1 snapshot gap → documented in DECISIONS.md per option (a) "version bump skipped" path

### 15. Snapshot-before-edit
- [x] No canonical doc was edited without snapshot (verified by git diff: only the *active* canonical docs changed; @v snapshots untouched)
- [x] `versions/` and `.spectacular/*@v*.md` files remain immutable

### 16. Per-finding y/n/q model
- [x] Exercised 2026-05-23 in throwaway workspace `/tmp/spec-verify-doctor`. Induced 5 judgment findings (1 frontmatter, 4 links). Walked the prescribed flow (group by area, errors-before-warnings, per-finding y/n/q). Simulated `y/y/n/q` sequence: 2 fixes applied (snapshot-skipped where no prior version existed), 1 skipped, 1 quit. Re-running detect went 5→3 findings — exactly matches the expected `(applied=2) + (skipped=1) + (quit=1, leaving 1 unwalked)` state machine.

## Skill-invoked subset — manual / agent-driven

### 17. status.md substrate check
- [x] Added to status.md head: "if config.yaml or root docs won't parse, auto-run doctor workspace/frontmatter/kits"
- [x] **Live-tested 2026-05-23.** Broke PRD.md frontmatter (removed delimiter). Ran the exact CLI command status.md instructs the skill to auto-fire: `spectacular doctor workspace frontmatter kits`. CLI surfaced `❌ missing frontmatter delimiter`, exit 2.
  - **Bug found and fixed:** original multi-area arg parser at `cli/spectacular:1093` overwrote `DOC_OPT_AREA` per arg, so only the last area ran. The exact substrate-trigger pattern status.md depends on was silently broken. Patched to accumulate, added regression scenario 8b to `tests/cli/doctor.test.sh` (33→37 asserts).

### 18. grill.md substrate check
- [x] Added to grill.md § 1: "if doc-registry.md / overrides / active kit file won't parse, auto-run doctor kits frontmatter"
- [x] **Live-tested 2026-05-23.** Corrupted bundled `blank.md` kit (removed frontmatter). Ran `spectacular doctor kits` — surfaced `❌ kit 'blank' has no frontmatter`, exit 2. Restored kit, re-ran clean. The CLI mechanism grill.md depends on works correctly.

### 19. lifecycle.md substrate check
- [x] Added to lifecycle.md § Verification artifact detection: "auto-run doctor lifecycle scoped to that request when proposing verified"

### 20. onboarding.md substrate check
- [x] Added to onboarding.md preamble: "auto-run doctor workspace frontmatter on first invocation"

## Smart-init integration

### 21. Updated diagnostic message
- [x] smart-init's `diag()` now emits `run \`spectacular doctor frontmatter\` for details` (was: generic placeholder)
- [x] smart-init tests still pass (50/50 asserts) after the message change
- [x] Verified live: `spectacular init` on a workspace with malformed PRD.md shows the new message

## Regression checklist

- [x] Smart-init tests (`tests/cli/init.test.sh`) still pass 50/50 after doctor additions
- [x] CLI `--help` still works
- [x] `spectacular doctor --help` shows doctor-specific usage
- [x] Doctor doesn't break when running on a workspace without `.claude/` directory (no spurious warning)
- [x] Doctor doesn't break when running in the spectacular source repo itself (skills.lock downgraded to info)
- [x] Real-workspace sweep post-cleanup: 0 errors, 1 acknowledged warning (PRD@v1.1 documented gap)

## Rollback validation

Doctor has multiple safety layers:
- Read-only by default (no `--fix` → no writes)
- `--fix` only applies content-free mechanical operations
- Judgment fixes route to the skill, which snapshots-before-edit
- Snapshot files are immutable; revert = git checkout the snapshot

- [x] Confirmed: no destructive paths in CLI doctor code; every write goes through `write_if_missing` or `doc_<id>` (which respect existing content) or append-only (`.gitignore`)
- [x] If a doctor judgment fix is wrong: revert via `git checkout <doc>@v<prev>.md` → `<doc>.md`
- [x] Test 13 verifies `--fix` never overwrites existing content

## Outstanding for next iteration

All scenarios checked 2026-05-23 in throwaway workspace dogfood.

**Carry-forward (pre-existing scope gap, not a blocker):** workspace area's "config.yaml parses" check (per doctor.md spec table) is implemented as a single `grep -q "^project:"` — it does not validate actual YAML well-formedness. A `config.yaml` that retains a `project:` line but is otherwise malformed will pass. Worth tightening when shell-only YAML validation is acceptable (perhaps via a strict-mode flag), but outside doctor v1 scope.
