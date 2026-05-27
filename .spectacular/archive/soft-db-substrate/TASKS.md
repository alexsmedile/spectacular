---
status: verified
updated: 2026-05-26
related:
  - PLAN.md
---

# Tasks — soft-db-substrate

## M1 — Schema lock ✅ DONE (2026-05-24)

- [x] Grill session: resolve 7 open questions + 1 follow-on (session linkage)
- [x] Draft frontmatter schema for `decisions/<slug>.md` entries (D8)
- [x] Draft frontmatter schema for `memory/<slug>.md` entries (D8)
- [x] Draft frontmatter schema for `sessions/<date>-<slug>.md` entries (D8)
- [x] Lock index file shape: MEMORY.md + SESSIONS.md new; DECISIONS.md unchanged (D3)
- [x] Worked example for each in `discovery.md`
- [x] `discovery.md` written

## M2 — Rules files ✅ DONE (2026-05-24)

- [x] `decisions-rules.md` — extended with v1.5.0 `spectacular decide` verb + Session: link field
- [x] `memory-rules.md` — created (mode: index, entries-dir convention introduced)
- [x] `sessions-rules.md` — created (mode: index, lifecycle invariant + auto-link mechanic)
- [x] New mode `index` added to `doc-index.md` taxonomy + verb × mode matrix
- [x] `doc-index.md` catalog rows added for memory + sessions
- [x] `templates/memory/entry.md` + `templates/sessions/entry.md` created
- [x] `SKILL.md` routing table + Doc IDs (v1.5.0) + references index + templates index updated
- [x] Doctor green: 0 errors / 0 warnings / 8 info (memory + sessions doctor areas land in M5)

## M3 — CLI writers ✅ DONE (2026-05-24)

- [x] `spectacular decide "<text>"` — appends ADR block to DECISIONS.md, auto-links Session: field
- [x] `spectacular remember "<text>" [--tag a,b]` — writes memory/<slug>.md, regenerates MEMORY.md, auto-links `session:` frontmatter
- [x] `spectacular session start [--tag ...]` — creates entry with `status: open`, refuses if one already open
- [x] `spectacular session end` — flips `status: closed`, sets `end_date`, scans for linked decisions+memories, appends Linked sections to body
- [x] Helpers: `_slug_from_text`, `_resolve_slug_collision`, `_summary_from_text`, `_active_session_slug`, `_resolve_template`, `_render_template`
- [x] Index regenerators: `_regen_memory_index`, `_regen_sessions_index` (table format, newest-first)
- [x] All verbs respect `--dry-run`
- [x] Dispatcher updated: `remember` routes to skill for bare/`this`, to CLI when text passed; `decide` + `session` added
- [x] Smoke test on fresh workspace: start → decide × 2 → remember → end. Multi-decision linkage, dry-run, lifecycle-invariant all verified
- [x] Bug fixed: `set -o pipefail` + empty `grep` exit-code chain in session end body-builder (wrapped with `|| true`, switched ADR scan to awk)
- [x] Template fixed: removed pre-baked "Linked decisions/memories" headers (session end writes them; template duplicate was confusing)

## M4 — Index files + kit changes ✅ DONE (2026-05-24)

- [x] Bootstrap `MEMORY.md` in this repo (table format with frontmatter incl. version + summary)
- [x] Bootstrap `SESSIONS.md` in this repo (table format)
- [x] DECISIONS.md left flat (per D3)
- [x] Coding kit triggers extended: `triggers-docs.suggested` += memory, sessions
- [x] `KNOWN_DOCS` extended (added `memory`, `sessions`)
- [x] `doc_memory()` + `doc_sessions()` scaffolders added; scaffold dispatch updated
- [x] Regenerators write `version: 1.0` + `summary:` so frontmatter doctor passes

## M5 — Doctor ✅ DONE (2026-05-24)

- [x] `check_memory()` — frontmatter, drift (entries vs index links)
- [x] `check_sessions()` — frontmatter, drift, lifecycle invariant (≤1 open), 4h stale warning
- [x] Both areas added to `DOC_AREAS`, `run_areas()`, `doctor_parse_args` allowlist, and `doctor_usage` help
- [x] Repo doctor: 0 errors / 0 warnings / 9 info (was 8 info pre-v1.5.0)

## M6 — Docs + ship ✅ DONE (2026-05-24)

- [x] `docs/commands.md` — added 4 mutator verbs (decide, remember, session start, session end) with examples
- [x] `doc-index.md` updated for memory + sessions doc-types + new `index` mode in taxonomy + verb × mode matrix column (done in M2)
- [x] CHANGELOG `[1.5.0]` entry: Added (2 doc-types, 4 verbs, 2 doctor areas, kit triggers, templates) / Changed / Fixed (em-dash awk bug, pipefail in session end) / Notes (v1.6.x deferrals)
- [x] Manifest bumps: cli SPECTACULAR_VERSION, .claude-plugin/plugin.json, .claude-plugin/marketplace.json (×2), .codex-plugin/plugin.json, SKILL.md frontmatter, README badge → all 1.5.0
- [x] Bug fixed: `_summary_from_text` no longer crashes on UTF-8 multibyte (em-dash) — switched awk char-iteration to sed
- [x] Bug fixed: session-end body builder pipefail (empty grep) — wrapped with `|| true`, switched to awk for ADR scan

## Remaining (user-triggered)

- [ ] Commit changes + create tag `v1.5.0` (user)
- [ ] `gh release create v1.5.0 --generate-notes` (user)
- [ ] `/plugin marketplace update spectacular` in Claude Code (user-triggered, can't be done remotely)
- [ ] Archive request after release

> M6 from original plan (decisions/ folder migration) dropped per D3 → moved to future request `memory-protocols-v2` (see `_research/agent-memory/REPORT.md`).

## Deferred / Open

- [x] All 7 PLAN open questions + 1 follow-on locked in `discovery.md`
- [ ] DECISIONS.md → decisions/ migration → v1.6.x request
- [ ] `--ignore-stale` flag for doctor sessions (decide if needed once 4h warning lands)
- [ ] `author/agent` provenance frontmatter field → v1.7.x multi-agent advisory work
