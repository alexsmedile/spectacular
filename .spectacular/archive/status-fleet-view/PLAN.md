---
status: archived
priority: medium
owner: alex
updated: 2026-07-09
build: b23
summary: "Enforce a strict PLAN/TASKS structure on active requests (unnumbered fixed headings, ### M milestones, flush-left checkboxes), then render spectacular status / status <slug> / status --json from that structure — frontmatter plus grep-safe body signals (Goal line, x/total progress, current milestone). doctor errors on structure drift for active requests (archive skipped). status --json is the agent opt-in contract; retires the hand-cached CLAUDE.md Active Requests table."
related:
  - PRD.md
archived: 2026-07-09
---

# Plan — status-fleet-view

## Goal

Enforce one canonical PLAN/TASKS structure on active requests, then make
`spectacular status` render the request fleet directly from that structure —
frontmatter plus grep-safe body signals — so no doc ever hand-caches the table
again and any agent can discover the fleet via `status --json` without opening a
file.

## Constraints

- **Bash 3.2 only** — no associative arrays; iterate files + reuse the existing
  `fm_get` helper (`cli/spectacular:711`) and `printf` for column alignment.
- **Enforced structure, not best-effort.** doctor errors (not warns) on active
  requests that violate the schema, so `status` body-grep can *assume* the
  canonical headings + normalized checkboxes exist. No graceful-degradation
  fallback branches.
- **Archive is skipped.** doctor never touches `.spectacular/archive/` — only
  active/planned requests under `.spectacular/requests/` are enforced. The 5
  legacy loose requests are archived-or-kept per the sequencing below, never
  auto-rewritten by this request.
- **Progressive disclosure, cheapest-first:** frontmatter fields (no read cost)
  → `## Goal` line + task counts (one anchored grep each) → agent opens the
  file. The fleet table renders from frontmatter + at most a per-file grep.
- Must not break `status --against-latest` / `status --since`, nor the skill
  judgment layer downstream.
- No new **required** frontmatter fields — the core 6 (status/priority/updated/
  summary/related/owner) are 100% present. `build:` (22/51) tolerates absence.

## The enforced schema

**PLAN.md canonical sections** (unnumbered, fixed set — `## Goal / ## Constraints /
## Milestones / ## Tasks / ## Dependencies / ## Validation / ## Deliverables`).
The 7 are **required and must appear in order**; additional `##` sections
(`## Understanding`, `## Decisions`, or request-specific ones like this PLAN's
`## The enforced schema` / `## Corpus audit`) are **allowed between them** —
doctor enforces the required set's presence + order, not a closed list. Chosen
over the current numbered dialect (`## 1. Goal …`) because a digit-free heading
is cleaner to grep and the template becomes the single source new requests
inherit.

**TASKS.md milestones** — `### M<n> — <title>` headings (already the template
default).

**Checkboxes** — flush-left `- [ ]` / `- [x]` are the counted units. Indented
`  - [ ]` nested subtasks are **allowed** (nested acceptance checklists under a
task) but **not counted** — progress is top-level only, keeping `x/total`
comparable across requests. `- [~]` is documented as a first-class **deferred**
state, counted separately (`5/8 (+1 deferred)`), excluded from the open/done
denominator.

## Corpus audit (why these rules)

Measured across all 51 PLAN/TASKS (active + archive), 2026-07-09:
- Frontmatter core 6 → **51/51**. `build:` 22/51, `target_version` 16/51.
- PLAN headings split into **two dialects**: numbered `## 1. Goal …` (5/6 active,
  the template default) and unnumbered `## Goal …` (1/6). Neither is greppable by
  a single fixed anchor today ⇒ the schema must *pick one and enforce it*.
- TASKS `### M` milestone headings → 31/51 (missing on legacy requests; present
  in the current template).
- Indented checkboxes → 5/51 files, **all archived**, all nested acceptance
  checklists under a parent task (never notes/metadata). 0 in active requests.
- `- [~]` deferred → 3 uses, undocumented in the template.
The 6 current active requests are the **doctor test corpus** — after the
enforcement lands they must flag as non-compliant, proving the check bites.

## Understanding

### How it works now

`spectacular status` (bare) is a **skill stub**: the CLI intercepts only
`--against-latest` and `--since`, then calls `skill_verb_message` and hands
rendering to the model (`cli/spectacular:11147`). There is no deterministic
frontmatter table. CLAUDE.md carries a hand-maintained "Active Requests" table
edited by hand on every request cut/archive — it already drifted once (missing
`cli-path-abstraction`). Every column in it is already a PLAN frontmatter field;
the `(b21)` bug-ref is the `build:` field re-typed as prose. The PLAN **template**
(`skills/spectacular/templates/plan/base.md`) still emits numbered headings, and
`plan-rules.md` references the 7 slots by number — so enforcement means changing
the template + rules, not just adding a check.

### What changes

- **Template + rules** adopt the unnumbered fixed section set (PLAN template,
  `plan-rules.md` slot prompts, `tasks-rules.md` checkbox note).
- **doctor** (lifecycle/docs area) gains a structure check that **errors** on
  active requests violating the schema; archive skipped. `doctor --fix` does the
  safe mechanical fixes (de-number headings, register `[~]`).
- **A mechanical `status)` path** (before `skill_verb_message`): bare → fleet
  table; `<slug>` → request card; `--json` → structured output.
- **CLAUDE.md** Active Requests table → one-line `spectacular status` pointer.
- **status.md** (skill) keeps its judgment layer but reads the CLI's `--json`.

### What stays the same

- `--against-latest` / `--since` untouched.
- The skill's proactive-signal detection (stale flags, "ready to archive") stays
  in `status.md` — CLI emits data, skill adds judgment (doctor pattern).
- AGENTS.md already omits the table — no change.

## Decisions

- Chose **enforced structure over best-effort/graceful degradation** — with the
  legacy loose files archived-or-flagged (never silently tolerated), `status`
  body-grep assumes canonical headings exist. Simpler code, stronger guarantee.
  (Reframed from the earlier warn-not-block design after deciding to enforce.)
- Chose **unnumbered fixed PLAN headings** (`## Goal …`) over numbered
  (`## 1. Goal …`) — digit-free is cleaner to grep and the template drives all
  new requests. Cost: re-heading the 6 active requests, done *after* the tooling
  ships (they are the test corpus first).
- Chose **`status --json` as the agent opt-in contract** over a separate
  `index.json` manifest — the CLI is the interface; a second on-disk artifact is
  exactly the drift risk this request kills. Agents opt in by running the command.
- Chose to **allow indented subtasks but count top-level only** — the nested
  acceptance-checklist pattern (5 archived uses) stays expressive, while `x/total`
  stays comparable across requests. (Rejected: forbid-indent — loses a real
  pattern; count-leaves — breaks cross-request comparability.)
- Chose to **document `[~]` as deferred**, counted separately, and **de-number
  headings in `doctor --fix`** as the mechanical repairs.
- Chose **CLI rendering over skill rendering** and **generator-as-source-of-truth
  over a shared partial + symlink** — deterministic, scriptable, drift-proof by
  construction (markdown can't transclude a fragment).

## Milestones

- M1 — Schema enforcement: PLAN template + `plan-rules.md` + `tasks-rules.md`
  adopt the unnumbered fixed section set; `doctor` (active requests only, archive
  skipped) **errors** on missing `## Goal` / `## Milestones` / `### M` and on
  disallowed checkbox states; `doctor --fix` de-numbers headings + registers
  `[~]`. The 6 active requests flag as non-compliant (proof the check bites).
- M2 — `spectacular status` prints an aligned fleet table from **frontmatter
  only** (slug · status · pri · build · updated · summary-trunc), sorted
  active→planned. Ships the drift-killing view.
- M3 — Fleet table gains **body signals** (assumes M1's enforced structure):
  `## Goal` line, `x/total` top-level task progress (`[~]` deferred shown), and
  current milestone (first `### M` with an open top-level task).
- M4 — `spectacular status <slug>` prints a request card: goal + full summary +
  `related:` deps + progress + current milestone + stale flag.
- M5 — `spectacular status --json` emits every field (frontmatter + body signals)
  — the agent opt-in contract for discovering the fleet without opening files.
- M6 — Retire the cache: CLAUDE.md Active Requests table → one-line pointer;
  `status.md` consumes `--json`. Then **convert the 6 active requests** to the
  new schema using `doctor --fix` + manual heading edits (they were the test
  corpus; now they conform).

## Tasks

See `TASKS.md`.

## Dependencies

- [[cli-path-abstraction]] — not blocking, but if it lands first the status
  command should read `requests/` via the centralized path var, not a literal.

## Validation

- M1 — run: `doctor lifecycle` (or docs) exits **non-zero** on a seeded active
  request missing `## Goal` and on one with a bad checkbox; exits 0 on a
  canonical one; **skips** an identically-broken file placed under `archive/`.
  `doctor --fix` rewrites `## 1. Goal` → `## Goal`.
- M2 — run: `tests/cli/status.test.sh` seeds ≥2 canonical requests at different
  statuses; asserts the table lists every slug with correct status/build and
  exits 0.
- M3 — run: seed a request with an indented subtask + a `[~]` line; assert the
  top-level count is correct (subtask excluded, `[~]` shown as deferred) and the
  current-milestone column matches the first `### M` with an open top-level task.
- M4 — run: `status <slug>` contains goal, full summary, each `related:` entry,
  the progress count, and the stale flag when `updated` is backdated >14d on an
  `active` request.
- M5 — run: `status --json | <json-validator>` parses; one object/request with
  the expected keys including body-signal fields.
- M6 — observable: CLAUDE.md no longer contains a per-request table (grep finds
  the pointer, not the rows); after conversion `doctor` on the 6 active requests
  exits 0.

## Deliverables

- Updated PLAN template (unnumbered headings), `plan-rules.md`, `tasks-rules.md`
  (indented-subtask + `[~]` note).
- `doctor` structure enforcement (errors, active-only, archive-skipped) +
  `--fix` heading de-numbering / `[~]` registration.
- `spectacular status` / `status <slug>` / `status --json` mechanical paths in
  `cli/spectacular`.
- `tests/cli/status.test.sh` + doctor structure-enforcement test (incl. the
  archive-skipped case).
- Updated `skills/spectacular/references/status.md` (consumes `--json`).
- CLAUDE.md Active Requests table replaced by a generator pointer.
- The 6 active requests converted to the new schema (M6).
