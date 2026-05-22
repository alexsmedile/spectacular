---
status: planned
updated: 2026-05-22
related:
  - PLAN.md
---

# Tasks — Doctor

## v1

### M1 — Detection taxonomy + report format
- [x] Draft `references/doctor.md` with full check list per area (skill, workspace, frontmatter, snapshots, links, lifecycle, kits)
- [x] Severity model locked — ✅ pass / ⚠️ warning / ❌ error / ℹ️ info
- [x] Report format locked: text default + `--format json`
- [x] JSON schema defined: `area`, `severity`, `file`, `message`, `proposed_fix`, `fix_type` (mechanical | judgment)
- [x] Repair flow documented in doctor.md § CLI vs skill split, § Mechanical fixes, § Judgment fixes

### M2 — CLI detect-only mode
- [x] Add `doctor` subcommand to `cli/spectacular` arg parser (dispatch in `main()`)
- [x] Implement `check_skill()` — install path + symlink + skills.lock parse + SKILL.md frontmatter (with source-repo detection downgrading missing skills.lock to info)
- [x] Implement `check_workspace()` — `.spectacular/` exists + always-set present + config.yaml parses
- [x] Implement `check_frontmatter()` — every canonical doc has required fields (`version`, `updated`, `summary`); ISO date check on `updated:`
- [x] Implement `check_snapshots()` — group `<doc>@vX.Y.md` by base, detect gaps in minor-version sequence per major
- [x] Implement `check_links()` — walk `related:` frontmatter, resolve relative to file's dir, flag missing targets
- [x] Implement `check_lifecycle()` — verification artifact detection (VERIFY.md OR TASKS § Verification OR PLAN § Validation/Success criteria/Acceptance); status mismatch detection; stale verified surfaces archive candidate
- [x] Implement `check_kits()` — bundled + local kit frontmatter parses; `extends: prd` only; `triggers-docs` references known doc IDs; PRD's declared kit resolvable
- [x] Wire scoped area as positional arg (`spectacular doctor frontmatter`)
- [x] Wire `--format text|json` flag
- [x] Exit codes: 0 (clean), 1 (warnings only), 2 (errors)
- [x] Tested against this workspace → 5 real findings surfaced including the missing PRD@v1.1.md snapshot

### M3 — CLI mechanical fixes
- [ ] Add `--fix` flag parsing
- [ ] `fix_gitignore_append()` — add `.spectacular.local/` entry if absent (no prompt — trivial)
- [ ] `fix_missing_dirs()` — create missing `requests/` / `current/` (no prompt — trivial)
- [ ] `fix_broken_symlink()` — repair `.claude/skills/spectacular` if dangling and `.agents/skills/spectacular/` exists (single-attempt, no prompt)
- [ ] `fix_missing_stub()` — re-scaffold missing always-set file with init's template (no prompt — restores baseline only)
- [ ] All `--fix` operations write to stdout: `✓ fixed <area>: <action>`
- [ ] Findings requiring judgment → still report, append "→ requires `/spectacular doctor --fix` for agent-driven repair"
- [ ] Confirm: `--fix` never modifies an existing non-empty canonical doc (judgment required → agent)

### M4 — Skill repair flow (`/spectacular doctor --fix`)
- [ ] Update SKILL.md routing: `spectacular doctor [<area>] [--fix]` → `references/doctor.md`
- [ ] In `references/doctor.md` § Repair flow: skill reads CLI's `--format json` report
- [ ] For each judgment-requiring finding: skill proposes a context-aware fix
- [ ] Per-finding `y/n/q` confirm flow
- [ ] Snapshot before any canonical-doc edit (route through existing versioning.md rule)
- [ ] On `q`: print remaining-findings summary + exit
- [ ] On all-applied: re-run detect, confirm clean state, summarize what changed

### M5 — Skill-invoked subset
- [ ] In `references/status.md` — if `doc-registry.md` parse fails during state read, auto-run doctor's `kits` + `frontmatter` checks, surface findings inline
- [ ] In `references/grill.md` — if kit file parse fails during pre-flight, auto-run doctor's `kits` check, surface findings inline, refuse to grill
- [ ] In `references/onboarding.md` — first-invocation flow runs doctor's `workspace` + `frontmatter` checks
- [ ] Document the skill-invoked surface in `references/doctor.md` § Skill-invoked checks

### M6 — Smart-init message update
- [ ] Update `cli/spectacular` `diag()` helper: replace generic message with area-specific pointer
- [ ] Frontmatter issue → `⊘ <file> (run \`spectacular doctor frontmatter\` for details)`
- [ ] Schema mismatch → `⊘ <file> (schema mismatch — run \`spectacular doctor frontmatter --fix\`)`
- [ ] Move corresponding v2-deferred item out of `smart-init/TASKS.md` deferred list (mark done)

### M7 — Tests + VERIFY.md
- [ ] Create `tests/cli/doctor.test.sh` — test harness reuses scenario pattern from init.test.sh
- [ ] Scenario: clean workspace → exit 0, no findings
- [ ] Scenario: missing always-set file → flagged + `--fix` re-creates it
- [ ] Scenario: malformed frontmatter → flagged + NOT auto-fixed (requires agent)
- [ ] Scenario: snapshot gap (synthetic v1.0 + v1.2 without v1.1) → flagged, exit 1
- [ ] Scenario: dangling `related:` link → flagged, exit 1
- [ ] Scenario: scoped area run (`spectacular doctor frontmatter`) skips other areas
- [ ] Scenario: `--format json` emits parseable JSON
- [ ] Scenario: `spectacular doctor --fix` mechanical-only path (no agent invocation needed)
- [ ] Create `requests/doctor/VERIFY.md` for agent-flow scenarios (per-finding judgment, snapshot before edit, q-to-quit) — these can't be automated, manual walkthrough required

### M8 — Dogfood
- [ ] Run `spectacular doctor` against this very workspace
- [ ] Confirm missing `PRD@v1.1.md` snapshot is surfaced
- [ ] Walk a real fix via `/spectacular doctor --fix` (best candidate: the snapshot gap or any frontmatter drift the scan finds)
- [ ] Run `spectacular doctor` post-fix → all findings resolved
- [ ] Run `tests/run.sh`; all green

### Verification (folded into TASKS per [[verification]] convention)
- [ ] All M7 test scenarios pass
- [ ] All VERIFY.md scenarios manually walked
- [ ] PLAN § Validation criteria each confirmed
- [ ] CHANGELOG.md entry written (v0.3.1 minor, or v0.4.0 if treated as feature release)

## v2 (deferred)

- [ ] Cross-doc semantic validation (PRD goals ↔ PLAN milestones) — belongs in doc-writer v2 anyway
- [ ] `prd diff <a> <b>` + `prd merge` subcommands — move to doc-writer's scope
- [ ] Performance benchmarks of skill operations
- [ ] Multi-workspace diagnostics
- [ ] Git history rewriting for snapshot backfill (currently proposes new commits only)
- [ ] `--yes-to-all` flag — explicitly rejected for v1, may revisit if user demand emerges
- [ ] JSON Schema validation of `doc-registry.md` itself
- [ ] Project-type-aware checks (skill project vs CLI project doctor behaviors)
