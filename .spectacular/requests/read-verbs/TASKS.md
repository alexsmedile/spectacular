---
status: review
updated: 2026-05-26
related:
  - PLAN.md
---

# Tasks — read-verbs

## M1 — Listing verbs (plural-noun = read)

- [x] Write `_request_iter_all` — walk `.spectacular/requests/*/PLAN.md`, emit one path per line
- [x] Write `_decision_iter_all` — walk `.spectacular/decisions/*.md` (and DECISIONS.md entries if folder absent)
- [x] Write `_memory_iter_all` — walk `.spectacular/memory/*.md`
- [x] Write `_session_iter_all` — walk `.spectacular/sessions/*.md`
- [x] Write `_parse_since "<7d|24h|30d|2w>"` → epoch threshold (returns 0 on parse fail)
- [x] Write `_emit_table` / `_emit_json` shared formatters (header + rows or JSON array)
- [x] `cmd_requests` — flags: `--status`, `--active` (alias `--status active`), `--since`, `--limit`, `--all`, `--json`
- [x] `cmd_decisions` — flags: `--tag`, `--since`, `--limit`, `--all`, `--json`
- [x] `cmd_memories` — flags: `--tag`, `--since`, `--limit`, `--all`, `--json`
- [x] `cmd_sessions` (read; distinct from `cmd_session` mutator) — flags: `--status open|closed|all`, `--since`, `--limit`, `--all`, `--json`
- [x] Top-level dispatch: `requests|decisions|memories|sessions) shift; cmd_<name> "$@"; exit $? ;;`
- [x] Smoke test: each verb returns expected rows on this workspace; `--json` parses with `jq`

## M2 — Detail verbs (singular-noun + slug)

- [x] `cmd_request <slug>` — skim view: frontmatter + section headers from PLAN.md + milestone tick rate from TASKS.md; `--full` dumps PLAN.md raw
- [x] `cmd_decision <slug>` — skim: frontmatter + body summary; `--full` dumps file
- [x] `cmd_memory <slug>` — same shape as decision
- [x] `cmd_sessions_show <slug>` — sub-verb under `sessions` to avoid `session start|end` collision; same skim shape
- [x] Top-level dispatch: `request|decision|memory) shift; cmd_<name> "$@"; exit $? ;;`
- [x] `sessions show <slug>` dispatched inside `cmd_sessions` subcommand router
- [x] All detail verbs accept `--json`

## M3 — Workspace verbs

- [x] `cmd_show <doctype>` — hardcoded canonical list: `prd|spec|principles|architecture|roadmap|stack|agents|decisions|memory|sessions|personas|feedback|idea`. Dispatch to file at expected location. `--section <name>` filters to that H2 block. `--json` returns `{path, content}`.
- [x] `cmd_summary` — call `cmd_requests --json`, `cmd_decisions --json --since 7d`, `cmd_memories --json --limit 5`, `cmd_sessions --json --status open`, aggregate counts + recent. Single-screen table. `--json` for machine.
- [x] `cmd_progress <slug>` — parse `.spectacular/requests/<slug>/TASKS.md`; group `- [x]` / `- [ ]` by `## M<N> —` headings; emit `M1: 8/8 ✓, M2: 3/5, ...`. `--json` for machine.
- [x] `cmd_paths` — emit JSON path map: `{prd, spec, principles, architecture, roadmap, stack, agents, decisions_file, decisions_dir, memory_index, memory_dir, sessions_index, sessions_dir, feedback_dir, ideas_dir, requests_dir, archive_dir, snapshots_dir, packs_dir}`. Default JSON; `--text` for human.
- [x] Top-level dispatch entries for all four

## M4 — Polish + cross-verb consistency

- [x] `--limit` defaults to 20; `--all` overrides; `--limit 0` is rejected
- [x] All verbs print `(no entries)` info line when result set is empty (table mode) or `[]` (JSON mode)
- [x] `--since` accepts `Nd`, `Nh`, `Nw` (days/hours/weeks); rejects bare numbers with helpful error
- [x] `--json` output for every list verb is a JSON array of objects with consistent field names: `{slug, status, priority, updated, summary, ...}`
- [x] Help text for each new verb includes 1-2 example invocations
- [x] Top-level `--help` adds a "Read verbs" section grouping requests/decisions/memories/sessions/show/summary/progress/paths

## M5 — Skill + docs + release

- [x] `skills/spectacular/SKILL.md` — add "Read verbs (v1.8.0+)" trigger block; bump version to 1.8.0
- [ ] `skills/spectacular/references/status.md` — note that `spectacular summary` is the cheaper cold-start alternative (one CLI call vs. status's deeper walk)
- [x] `README.md` — add to CLI reference: all 8 new verb groups
- [ ] `README.md` — workspace tree comment doesn't need changes; add note in "What This Is" if a sentence about agent-friendly reads fits (deferred — not necessary for ship)
- [x] `CHANGELOG.md` — `[1.8.0]` entry: Added (read-verb family), Changed (SKILL routing), Notes (grammar decision rationale)
- [x] All 5 version manifests to 1.8.0 (.claude-plugin/marketplace.json, .claude-plugin/plugin.json, .codex-plugin/plugin.json, cli/spectacular, skills/spectacular/SKILL.md)
- [x] Dogfood: from a fresh terminal, run `spectacular summary`, `spectacular requests --active`, `spectacular progress read-verbs` — all return useful output in <500ms
- [x] `spectacular doctor` clean
- [x] Flip PLAN.md `status: review`

## Deferred

- [ ] `spectacular ls <doctype>` alias dispatcher (e.g. `ls requests` → `requests`). Tabled — adds a surface without clear win once plural-noun reads exist.
- [ ] `spectacular recent` aggregate dashboard. Tabled — `summary` covers the cold-start use case; `--since` covers the time-filter use case.
- [ ] `spectacular ideas` / `idea <slug>` read verbs. Skipped intentionally: `spectacular idea list` already covers it (v1.7.0). Adding plural form would be a rename for cosmetic consistency only. Re-evaluate if other v1.8.0 work makes ideas the odd one out.
- [ ] `spectacular feedback` plural-form alias. Same logic — `feedback-loop list` exists.
- [ ] `--watch` flag (re-emit on file changes). Out of scope; cron/agent loop covers it.
