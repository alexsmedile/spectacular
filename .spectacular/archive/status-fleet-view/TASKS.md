---
status: verified
updated: 2026-07-09
related:
  - PLAN.md
---

# Tasks — status-fleet-view

## v1

### M1 — Schema enforcement (structure doctor)
- [x] Update PLAN template `skills/spectacular/templates/plan/base.md`: unnumbered fixed sections (`## Goal / ## Constraints / ## Milestones / ## Tasks / ## Dependencies / ## Validation / ## Deliverables`), drop `## 1.`–`## 7.` prefixes.
- [x] Update `plan-rules.md` slot prompts + `tasks-rules.md` checkbox note (allow indented subtasks; document `- [~]` deferred; flush-left is counted).
- [x] `doctor` (active requests only, `.spectacular/archive/` skipped): **error** on missing/mis-ordered required PLAN sections, missing TASKS `### M`, and disallowed checkbox states.
- [x] `doctor --fix`: de-number PLAN headings (`## 1. Goal` → `## Goal`); leave `[~]` in place (now documented).
- [x] → check: doctor exits non-zero on a seeded active request missing `## Goal`; exits 0 on a canonical one; **skips** the same breakage under `archive/`; `--fix` de-numbers headings. The 6 current active requests flag as non-compliant.

### M2 — Fleet table (frontmatter-only)
- [x] Add mechanical `status)` path: intercept bare `status` before `skill_verb_message`.
- [x] Iterate `requests/*/PLAN.md`, `fm_get` status/priority/build/updated/summary; sort active→planned; `printf` aligned columns; truncate summary; tolerate missing `build:`.
- [x] → check: `tests/cli/status.test.sh` seeds ≥2 canonical requests, asserts every slug + status/build column, exit 0.

### M3 — Body signals (assumes enforced structure)
- [x] Extract `## Goal` line per PLAN (exact anchored grep — schema guarantees it).
- [x] Count top-level progress: anchor `^- \[`, count `[x]` vs `[ ]`, show `[~]` as deferred; ignore indented `  - [ ]`.
- [x] Current milestone: first `### M` with an open top-level `- [ ]`.
- [x] → check: seed a request with an indented subtask + a `[~]`; assert top-level count excludes the subtask, shows the deferred, and the milestone column matches.

### M4 — Request card
- [x] Intercept `status <slug>`: resolve `requests/<slug>/PLAN.md`, error if missing.
- [x] Print card: goal + full summary + `related:` deps + progress + current milestone + stale flag (updated >14d & active).
- [x] → check: test asserts goal, summary, related deps, progress, and stale flag on a backdated active request.

### M5 — JSON contract (agent opt-in)
- [x] Add `--json`: one object/request with frontmatter + body-signal fields.
- [x] Document in AGENTS.md that agents discover the fleet via `spectacular status --json`.
- [x] → check: `status --json` parses; expected keys present including body-signal fields.

### M6 — Retire the cache + convert the corpus
- [x] Replace CLAUDE.md Active Requests table with a one-line `spectacular status` pointer.
- [x] Update `skills/spectacular/references/status.md` to consume `--json`, keep judgment layer.
- [x] Convert the 6 active requests to the new schema (`doctor --fix` + manual heading edits).
- [x] → check: grep finds pointer not rows in CLAUDE.md; `doctor` on the 6 converted requests exits 0.
