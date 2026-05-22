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
- [x] Add `--fix` flag parsing (within `doctor_parse_args()`)
- [x] `.gitignore` append/create logic (handles 3 cases: no file → create, missing entry → append, has entry → no-op)
- [x] Missing-dir fix (`mkdir -p .spectacular/requests` and `current`)
- [x] Dangling-symlink fix (remove broken `.claude/skills/spectacular` + recreate pointing to `.agents/skills/spectacular`)
- [x] Missing always-set file fix (re-scaffold PRD/config/AGENTS via existing `doc_<id>` functions)
- [x] All `--fix` operations write `✓ fixed [<area>]: <action>` to stdout
- [x] Judgment findings still report; ANSI suffix shows `(judgment — \`/spectacular doctor --fix\`)`
- [x] Verified: `--fix` never touches existing non-empty canonical doc (uses `write_if_missing` semantics from smart-init)
- [x] Added `.gitignore` missing-entry check to `check_workspace()` so the fix path actually fires

### M4 — Skill repair flow (`/spectacular doctor --fix`)
- [x] Update SKILL.md routing: `spectacular doctor [<area>] [--fix]` → `references/doctor.md`
- [x] In `references/doctor.md` § Repair flow: skill reads CLI's `--format json` report; documented walkable procedure (8 steps)
- [x] For each judgment-requiring finding: skill proposes a context-aware fix (4 worked examples in doctor.md)
- [x] Per-finding `y/n/q` confirm flow documented
- [x] Snapshot-before-edit rule documented (canonical docs only; per-request files have `snapshot-on-edit: false`)
- [x] On `q`: documented exit behavior
- [x] On all-applied: re-run detect documented
- [ ] **Live end-to-end test of the flow** — deferred to next workspace with real drift requiring agent walkthrough (see VERIFY scenario 16)

### M5 — Skill-invoked subset
- [x] In `references/status.md` — added substrate-check preamble (auto-run doctor workspace/frontmatter/kits on parse failure)
- [x] In `references/grill.md` — added substrate-check to § 1 (auto-run doctor kits frontmatter on parse failure)
- [x] In `references/onboarding.md` — added substrate-check to header (always runs on first invocation)
- [x] In `references/lifecycle.md` — added substrate-check to verification-artifact detection (auto-run scoped lifecycle when proposing verified)
- [x] Documented the skill-invoked surface in `references/doctor.md` § Skill-invoked checks
- [ ] **Live end-to-end test** of skill-invoked auto-checks — deferred to natural workspace corruption (VERIFY scenarios 17, 18)

### M6 — Smart-init message update
- [x] Update `cli/spectacular` `diag()` helper: replaced generic placeholder with `run \`spectacular doctor frontmatter\` for details`
- [x] Verified live: malformed PRD.md scenario shows new message
- [x] Smart-init's 50/50 tests still pass after the message change

### M7 — Tests + VERIFY.md
- [x] Create `tests/cli/doctor.test.sh` — 11 scenarios covering detect + mechanical-fix + scoped + json
- [x] Scenario: clean workspace → exit 0, no findings
- [x] Scenario: missing always-set file → flagged + `--fix` re-creates it
- [x] Scenario: malformed frontmatter → flagged + NOT auto-fixed (requires agent)
- [x] Scenario: snapshot gap (synthetic v1.0 + v1.2 without v1.1) → flagged, exit 1
- [x] Scenario: dangling `related:` link → flagged, exit 1
- [x] Scenario: lifecycle active without SESSION.md → flagged, exit 1
- [x] Scenario: scoped area run (`spectacular doctor frontmatter`) skips other areas
- [x] Scenario: `--format json` emits parseable JSON (python json.load validates)
- [x] Scenario: `spectacular doctor --fix` mechanical-only path (gitignore + missing files)
- [x] Scenario: `doctor --help` shows usage with --fix + --format
- [x] All 33 asserts across 11 scenarios passing
- [x] Create `requests/doctor/VERIFY.md` — load-bearing checklist mirroring tests + agent-flow scenarios
- [x] Smart-init tests still 50/50 green after doctor additions
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
