---
status: archived
priority: medium
owner: alex
updated: 2026-06-28
build: b9
summary: "Split DECISIONS.md into a cheap index + per-entry files when it grows large — agents load only the index by default, fetch individual entries on demand."
related:
  - ../../PRD.md
  - ../../roadmaps/index.md
archived: 2026-06-28
---

# Plan — decisions-index

> **Origin (2026-06-15):** Octopus project hit a 2028-line DECISIONS.md (D1–D110) after ~6 months of use. Loading it whole before every planning session is a significant context tax. The `decisions-rules.md` schema already documents index mode; the CLI support is missing.

## Understanding

### How it works now

`spectacular decide "<text>"` appends a full ADR block directly to `.spectacular/DECISIONS.md`. After 50+ entries the file becomes expensive to load — it's read in full by agents before any planning-phase work. There is no way to split it or load selectively. `decisions-rules.md` already documents an index-mode layout and per-entry file format, but nothing in the CLI or skill enforces or migrates to it.

### What changes

- `DECISIONS.md` can operate in two modes: **flat** (current, backwards compatible) or **index** (new, one line per decision).
- **Index mode** is detected by presence of a `decisions/` subfolder next to `DECISIONS.md`.
- In index mode, `spectacular decide` writes `decisions/D<N>.md` (full ADR prose) and appends one line to the index.
- `spectacular decisions migrate` splits a flat file into index mode in one shot.
- `decisions-rules.md` frontmatter gains `mode: index` when migrated (detected at doc-load time).
- `doctor` gains awareness: if `decisions/` exists but `DECISIONS.md` is still flat prose (not an index), it flags it.

### What stays the same

- Flat `DECISIONS.md` remains valid — projects that never grow past ~50 entries need no action.
- Entry schema unchanged: `## YYYY-MM-DD — title`, Context / Decision / Consequences / Session.
- `spectacular decision <slug>` read verb is already correct (reads individual file by slug).
- `spectacular decide --dry-run` behavior unchanged.
- The append-only immutability convention still holds; `migrate` is the only write-many operation.

## 1. Goal

Ship the **DECISIONS.md index mode** — a soft-folder split that keeps the root file cheap (one line per decision) and loads full ADR prose on demand — plus the CLI migration verb that converts any existing flat file.

## 2. Constraints

- **Backwards compatible.** Flat `DECISIONS.md` must continue to work with zero changes. Mode detection is structural (folder presence), not a config flag.
- **Bash 3.2.** All CLI code targets macOS default bash 3.2 (no associative arrays, no `mapfile`).
- **No data loss.** `decisions migrate` must write all per-entry files before touching `DECISIONS.md`. `--dry-run` previews and writes nothing.
- **Dovetails with roadmap-ledger.** Both ship in the v1.17.0 window; `roadmap-ledger` owns the ledger; this owns the decisions split. No shared code.

## 3. Milestones

- M1 — **Schema + detection rule.** `decisions-rules.md` updated: `mode: index | flat`, detection rule (folder presence), index line format canonical, per-entry file format canonical. No CLI changes yet.
- M2 — **`decisions migrate` verb.** Reads flat `DECISIONS.md`, extracts each `## YYYY-MM-DD —` block into `decisions/D<N>.md`, rewrites root as index. `--dry-run` previews. Idempotent on already-migrated workspace.
- M3 — **`decide` writes index mode when active.** In index mode, `decide` writes `decisions/D<N>.md` + appends one-liner to index (not to the full prose). In flat mode, behavior unchanged.
- M4 — **`doctor` area.** `doctor decisions` checks: flat vs index detected correctly; no orphan entries (index line missing a file); no stale files (file missing an index line); correct `D<N>` numbering (no gaps, no duplicates).
- M5 — **Dogfood + ship.** Run `decisions migrate` on this repo's `.spectacular/DECISIONS.md`; add a new decision via `decide` in index mode; verify one-row write; CHANGELOG + plugin bump.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- No hard dependencies on other active requests.
- Soft: `cli-debt-removal` (b4) — both target v1.17.0. No shared code; order-independent.
- Soft: `roadmap-ledger` (b7) — also targeting v1.17.0. Ledger must be updated when this request is slotted.

## 6. Validation

- M1 — `decisions-rules.md` frontmatter has `mode: index | flat`; a reader can derive the detection rule and both file formats from the doc alone.
- M2 — `spectacular decisions migrate --dry-run` on a 5-entry flat file prints N preview lines + N filenames; no files written. Real run writes `decisions/D1.md`…`decisions/D5.md` + rewrites `DECISIONS.md` as 5-line index; original prose gone from root.
- M3 — After migration, `spectacular decide "New decision"` writes `decisions/D6.md` + appends `- **D6** — New decision — …` to index; flat `DECISIONS.md` not touched.
- M4 — `doctor decisions` flags a deliberately-orphaned entry (index line with no file) and a deliberate gap in numbering.
- M5 — This repo's DECISIONS.md migrated; new decision added; `grep -c "^\*\*Context" .spectacular/DECISIONS.md` returns 0 (no prose in index); `ls .spectacular/decisions/ | wc -l` returns the correct entry count.

## 7. Deliverables

- Updated `skills/spectacular/references/decisions-rules.md` (mode field, detection rule, both formats)
- `spectacular decisions migrate [--dry-run]` CLI verb
- `spectacular decide` updated to write index mode when `decisions/` exists
- `doctor decisions` area (new or extended)
- This repo's DECISIONS.md migrated (M5 dogfood)
- CHANGELOG entry
