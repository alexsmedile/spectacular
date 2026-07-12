---
status: verified
updated: 2026-07-12
related:
  - PLAN.md
---

# Tasks — review-sweep

## v1

### M1 — Evidence freshness schema
- [x] verify.md § VERIFY-LOG shape: `against: <commit/build> · <identity>` stamp on `[manual]`/`[observe]` rows; `pending-reverify` marker semantics (what flips it, what clears it) — ✔ shape shows against: + ⟳ pending-reverify rows
- [x] verify.md walk loop: recording a manual/observe pass asks for the stamp — ✔ Evidence-stamp paragraph after manual kind
- [x] doctor lifecycle: warn on unstamped `[manual]`/`[observe]` row in a review/verified request's VERIFY-LOG; warn on any `pending-reverify` row in review/verified — ✔ sandbox: 2 warnings fire, stamped log clean
- [x] tests/cli/doctor.test.sh scenario: unstamped row → warning; pending-reverify row → warning; stamped clean log → no warning — ✔ scenario 23, 4 asserts green (74 passed, 4 pre-existing env fails)
- [x] → check: PLAN §Validation M1

### M2 — Sweep protocol + trigger
- [x] references/review-sweep.md: three tiers (review = full audit; ticked-active = audit + propose advance; planned = overlap check vs specs/index.md + archive), findings → VERIFY-LOG sweep entry, handoff → SESSION.md § Next actions, never-promotes rule, orchestrator-is-the-mutator rule, per-request parallel fan-out, sweep-entry shape — ✔ 5.3KB, three tiers + loop + entry shape
- [x] SKILL.md routing row: `spectacular sweep [<slug>]` → review-sweep.md — ✔ verification routing table
- [x] CLI `sweep` redirect stub (mirrors `verify` redirect) — ✔ exit 1 (verify convention)
- [x] docs/commands.md sweep section — ✔ before prd section
- [x] → check: PLAN §Validation M2

### M3 — request-auditor agent
- [x] agents/request-auditor.md: read-only (Read/Grep/Glob/Bash), small-model, closed audit brief, findings-block return shape, bounces on judgment outside brief — ✔ read-only tools, model: haiku
- [x] .claude/agents/request-auditor.md relative symlink — ✔ ../../agents/request-auditor.md
- [x] CLAUDE.md fleet grid list + repo AGENTS.md if it lists agents — ✔ both CLAUDE.md + AGENTS.md
- [x] → check: PLAN §Validation M3

### Wrap
- [x] Roadmap ledger row b31 → v1.35.0 — ✔ v3.14, status active
- [x] CHANGELOG [Unreleased] entry — ✔ Review sweep (b31) section
- [x] → check: doctor links clean; full doctor no new findings (env yaml excepted)
