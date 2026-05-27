---
status: archived
priority: high
owner: alex
updated: 2026-05-26
target_version: v1.5.0
summary: "Promote decisions, memory, and sessions from flat files to soft-folder databases (index.md + entries/), and add mutator verbs: spectacular decide, spectacular remember, spectacular session start|end"
related:
  - ../../ROADMAP.md
  - ../../SPEC.md
  - ../../PRINCIPLES.md
archived: 2026-05-26
---

# Plan — soft-db-substrate

## 1. Goal

Spectacular's "what was decided / what should I remember / what happened in this session" surface currently lives in three places of varying maturity:

- **DECISIONS.md** — a flat ADR-style log (works, doesn't scale past ~50 entries)
- **memory/** — a directory the skill writes to via `spectacular remember this`, no index, no schema
- **sessions** — no doc-type at all today

This request unifies all three under the **soft-folder database pattern** already proven by `SPEC.md` + `specs/`: a thin index file at the root, one subfolder per entity, frontmatter as the signal layer. It then adds three CLI write verbs (`decide`, `remember`, `session start|end`) so capture is one keystroke instead of a paragraph of instructions.

This is the foundation block for v1.5.x → v1.6.x (query verbs). It does **not** add retrieval/query in this request — only structure + writes.

## 2. Constraints

- **Doc-writing engine reuse.** All three new doc-types must register via rules-file frontmatter (per v1.4.0 substrate-clarity work). No bespoke handlers.
- **Mutator vs agentic split.** Writes are mechanical (CLI). Grill/refine/review on the index files remain skill-only.
- **No breaking change to existing memory/.** Current `memory/` files are claimed by this system as-is; migration is additive (introduce MEMORY.md index, leave existing files).
- **DECISIONS.md migration is opt-in.** Existing flat DECISIONS.md keeps working; promotion to `decisions/` is offered by doctor, never forced.
- **No retrieval logic.** Query verbs (`decisions --7d`, `recall`, `sessions`) are deferred to v1.6.x. This block only writes well-shaped data.
- **Frontmatter shape must be RAG-ready.** Even though v1.5.x doesn't query, the schema (tags, related, date, type) must be the one a future embedding/RAG layer will read.

## 3. Milestones

- M1 — Schema lock: design index + entry frontmatter for all three doc-types; verify against existing `SPEC.md`/`specs/` precedent; document in `discovery.md`.
- M2 — Rules files: write `decisions-rules.md`, `memory-rules.md`, `sessions-rules.md` with mode + slots + triggers-docs.
- M3 — CLI writers: implement `spectacular decide "..."`, `spectacular remember "..."`, `spectacular session start|end`; all three should be ~30 lines of bash each.
- M4 — Index files: bootstrap `DECISIONS.md` (migrate existing format), `MEMORY.md`, `SESSION.md` (singular — most recent session at root; history in `sessions/`).
- M5 — Doctor: add `doctor memory|sessions` areas — orphan entries, missing index links, malformed frontmatter, 4h stale-session warning, session inverse-link.
- M6 — Docs + ship: `docs/commands.md` updated, CHANGELOG entry, release as v1.5.0.

> Migration of flat `DECISIONS.md` → `decisions/` folder shape was dropped from this request per `discovery.md` D3 — moves to v1.6.x alongside query verbs.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- v1.4.0 substrate-clarity (DONE) — rules-file frontmatter schema is what these three new doc-types plug into.

## 6. Validation

- M1 — `discovery.md` contains locked frontmatter schemas for all three doc-types with one worked example each
- M2 — Three new rules files exist; `doc-index.md` regenerated catalog includes them
- M3 — Three CLI commands work end-to-end on a fresh test workspace; each writes a valid entry + updates the index
- M4 — `spectacular init` on a fresh project scaffolds the three index files (or leaves them out per kit, TBD in grill)
- M5 — `spectacular doctor memory` / `doctor sessions` exit clean on this repo; 4h stale warning triggers on synthetic open session
- M6 — Manual smoke: `spectacular session start` → `spectacular decide "test"` → entry has `session:` set → `spectacular session end` → SESSIONS.md shows 1 decision linked

## 7. Deliverables

- 3 new rules files in `skills/spectacular/references/`
- 3 new CLI verbs in `cli/spectacular`
- 3 index files in `.spectacular/` (DECISIONS.md exists today; MEMORY.md + SESSION.md are new)
- Updated `doc-index.md`, `SPEC.md`, `docs/commands.md`, `CHANGELOG.md`
- This repo's own `.spectacular/decisions/`, `.spectacular/memory/`, `.spectacular/sessions/` dogfooded

## Open questions

All 7 PLAN-level questions + 1 follow-on resolved in `discovery.md` (D1–D9, 2026-05-24).

## Out of scope

- Query/retrieval verbs (`decisions --7d`, `recall`, `sessions`) — v1.6.x
- Embedding / RAG implementation — future, schema-ready only
- Multi-agent session coordination — v1.7.x advisory work, not here
- UI for browsing — never in CLI scope

## Notes

This is the largest single-version change since v1.4.0 substrate-clarity. Three new doc-types + three new mutator verbs + migration tooling. Grill heavily before opening M2.
