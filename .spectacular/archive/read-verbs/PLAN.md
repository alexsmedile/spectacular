---
status: archived
priority: high
owner: alex
updated: 2026-05-26
target_version: v1.8.0
summary: "Add read-only listing + detail verbs (requests, decisions, memories, sessions, show, summary, progress, paths) to collapse multi-step agent workflows into single deterministic calls."
related:
  - ../../ARCHITECTURE.md
  - ../../SPEC.md
archived: 2026-05-26
---

# Plan — read-verbs

## 1. Goal

Add a coherent family of read-only CLI verbs so AI agents (and humans) can answer common workspace questions in **one CLI call** instead of N file reads + frontmatter parses. Codifies the existing pragmatic grammar (bare verbs for high-frequency actions; noun-namespace for multi-verb lifecycles; **plural-noun = read** as the new pattern). Zero breaking changes — every existing verb stays.

## 2. Constraints

- **Locked grammar (2026-05-26):**
  - Bare top-level verb = high-frequency action on the implicit object (`new`, `archive`, `promote`, `remember`, `decide`, `snapshot`, `touch`)
  - Noun + subcommand = multi-verb lifecycle (`session start|end`, `idea new|list|promote`, `feedback-loop ...`, `pack ...`)
  - **Plural noun = list** (new in v1.8.0): `requests`, `decisions`, `memories`, `sessions`
  - **Singular noun + slug = detail** (new): `request <slug>`, `decision <slug>`, `memory <slug>`. Session detail is `sessions show <slug>` to avoid the `session start|end` collision.
- **Universal flags on every list verb:** `--status <s>`, `--since <duration>`, `--limit N` (default 20), `--all`, `--json`.
- **Default output: table** (SLUG / STATUS / PRIORITY / TARGET / SUMMARY-tail). `--json` flag for machine consumers.
- **Default limit: 20** (override with `--limit N` or `--all`).
- **Detail view = skim by default** (frontmatter + section headers + tasks-progress for requests), `--full` flag dumps raw markdown. Token-cheap by default.
- **`show <doctype>` hardcodes the canonical list** (not registry-driven). Adding doc-types is rare enough to not justify dynamic dispatch here; keeps the verb predictable.
- **`summary` aggregates by calling list verbs internally** — no duplicated parsing logic. Single source of truth per noun.
- **No new processes, no new dependencies.** Pure bash + the existing `fm_get` / yaml helpers.
- **All verbs read-only.** No state mutation. Safe to call any number of times.

## 3. Milestones

- **M1** — Listing verbs land: `requests`, `decisions`, `memories`, `sessions` (read) — table + `--json` + filter flags
- **M2** — Detail verbs land: `request <slug>`, `decision <slug>`, `memory <slug>`, `sessions show <slug>` — skim + `--full`
- **M3** — Workspace verbs: `show <doctype>`, `summary`, `progress <slug>`, `paths` (JSON path map)
- **M4** — Universal `--since <duration>` parser (e.g. `7d`, `24h`, `30d`) shared across all list verbs
- **M5** — SKILL.md routing rows + doc updates (README, ARCHITECTURE.md) + CHANGELOG entry; bump to 1.8.0 and release-ready

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- v1.7.0 ships first (already done — `ideas-doctype`).
- `fm_get` / `fm_set` helpers (already present).
- Existing `_idea_iter_all`, `_feedback_iter_all` patterns as templates for `_request_iter_all`, `_decision_iter_all`, etc.

## 6. Validation

- **M1** — `spectacular requests --active` returns one line per active request; `requests --json` parses with `jq`; `decisions --since 7d --limit 5` filters correctly.
- **M2** — `request ideas-doctype` prints frontmatter + section outline + milestone tick rate; `--full` dumps PLAN.md content verbatim.
- **M3** — `show prd` dumps PRD.md; `summary` returns single-screen overview with counts across all 4 plural nouns; `progress read-verbs` shows M1: 0/N, M2: 0/N etc.; `paths --json | jq .prd` returns `.spectacular/PRD.md`.
- **M4** — `--since 7d` filters all four list verbs consistently. `--since 24h`, `--since 30d` all parse.
- **M5** — README CLI reference includes every new verb; SKILL.md has trigger rows. Doctor clean. All 5 version manifests at 1.8.0.

## 7. Deliverables

- `cli/spectacular` — new functions:
  - `cmd_requests` / `cmd_request` (+ `_request_iter_all`, `_request_status`)
  - `cmd_decisions` / `cmd_decision`
  - `cmd_memories` / `cmd_memory`
  - `cmd_sessions` / `cmd_sessions_show` (subverb to avoid `session` mutator collision)
  - `cmd_show` (doctype dispatcher with hardcoded canonical list)
  - `cmd_summary` (aggregates from list verbs)
  - `cmd_progress` (parses TASKS.md milestone ticks)
  - `cmd_paths` (JSON path emitter)
  - `_parse_since` (universal duration parser: `Nd`, `Nh`, `Nw`)
  - `_emit_table` / `_emit_json` formatters (shared)
- Top-level dispatch entries for all new verbs
- `--help` updated to list new verbs grouped under "Read verbs"
- `skills/spectacular/SKILL.md` — new "Read verbs" trigger block
- `README.md` — CLI reference + summary doc
- `CHANGELOG.md` — `[1.8.0]` entry
- All 5 manifests bumped to 1.8.0
