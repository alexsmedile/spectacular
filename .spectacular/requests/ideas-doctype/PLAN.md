---
status: review
priority: medium
owner: alex
updated: 2026-05-26
target_version: v1.7.0
summary: "Register `ideas` as a first-class doc-type — SKILL triggers, idea-rules, doc-index row, template, CLI verbs, doctor area. Closes the gap where ARCHITECTURE.md documents `ideas/` but the skill can't create, list, or check ideas."
related:
  - ../../ARCHITECTURE.md
  - ../../ideas/memory-protocols.md
---

# Plan — ideas-doctype

## 1. Goal

Promote `.spectacular/ideas/` from "documented folder convention" to "first-class doc-type the skill + CLI know how to operate on" — covering the same surface as `feedback` (rules file, CLI mutators, doctor area, template, doc-index entry, SKILL.md triggers).

## 2. Constraints

- **No new substrate semantics.** `ideas/` is already documented in `ARCHITECTURE.md` as "thinking scratchpad, not a workflow stage, not acted on automatically." This request operationalizes that doc — it does not change the philosophy.
- **Mirror `feedback` shape.** `feedback-rules.md` + `cmd_feedback_loop` are the closest precedent (also soft-folder DB, mode: `index`, judgment-only doctor, hidden aliases pattern). Reuse the same file structure, frontmatter conventions, and CLI helpers (`_slug_from_text`, `_resolve_slug_collision`, `_render_template`, `fm_get`, `fm_set`).
- **Promotion flow already exists.** `references/new-request.md` lines 147-150 already document `ideas/<idea>.md` → request scaffolding → move to `archive/ideas/`. Do not re-design this — just wire `spectacular idea promote <slug>` as the entry point.
- **Doctor must be judgment-only.** No `--fix` for ideas — they're inherently human-curated. Mechanical fixes don't apply.
- **Design decisions locked (2026-05-26):**
  - Status enum: `parked` / `exploring` / `promoted`
  - Verb surface: both top-level (`spectacular idea ...`) and doc-form (`spectacular idea grill|refine|review`)
  - Doctor stale threshold: 90 days for `exploring` status
  - Promotion target: `archive/ideas/<slug>.md` (matches existing convention)

## 3. Milestones

- **M1** — `idea-rules.md` + template + doc-index row land; skill can route `spectacular idea grill` through the generic engine
- **M2** — CLI mutator verbs work: `spectacular idea new <slug>`, `idea list`, `idea promote <slug>`
- **M3** — `doctor ideas` area registered + judgment checks pass on the existing `.spectacular/ideas/memory-protocols.md`
- **M4** — SKILL.md triggers updated; `idea` added to registered doc-IDs list; ARCHITECTURE.md cross-references new rules file
- **M5** — Dogfood: scaffold one new test idea via `spectacular idea new`, verify doctor, promote it to a throwaway request, verify it moves to `archive/ideas/`

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- v1.6.0 substrate (`feedback` doc-type) — the canonical reference for soft-folder-DB shape. Already shipped.
- `new-request.md` reference doc (lines 147-150 cover idea→request promotion). Already exists; this request just wires CLI entry to it.

## 6. Validation

- **M1** — `spectacular idea grill` opens the generic grill engine with `idea-rules.md` loaded; doc-index lists `idea` row.
- **M2** — `spectacular idea new test-foo` creates `.spectacular/ideas/test-foo.md` with stub frontmatter (`status: parked`); `idea list` shows it; `idea promote test-foo` scaffolds a request and moves the file to `archive/ideas/`.
- **M3** — `spectacular doctor ideas` runs clean against `.spectacular/ideas/memory-protocols.md`; if I set `status: exploring` and backdate `updated:` >90 days, the area emits a warning.
- **M4** — `grep "idea" SKILL.md doc-index.md` returns the new rows; `idea-rules` is in the doc-IDs list at line ~63 of SKILL.md.
- **M5** — Full create→list→doctor→promote→archive cycle works without manual file edits.

## 7. Deliverables

- `skills/spectacular/references/idea-rules.md` (new)
- `skills/spectacular/templates/idea/base.md` (new)
- `skills/spectacular/references/doc-index.md` — one new row under "Project-wide canonical docs"
- `skills/spectacular/SKILL.md` — new trigger rows; updated doc-IDs registered list; bumped version to 1.7.0
- `cli/spectacular` — `cmd_idea` dispatcher (new/list/promote subverbs); `check_ideas` doctor function; `DOC_AREAS` extended with `ideas`; `doctor_parse_args` case added; help string bumped to "15 areas"; `SPECTACULAR_VERSION="1.7.0"`
- `skills/spectacular/references/doctor-areas.md` — new `## ideas (v1.7.0+)` section
- `.spectacular/ARCHITECTURE.md` — cross-reference to `[[idea-rules]]` in the existing `ideas/` section
- All version manifests bumped to 1.7.0 (`.claude-plugin/marketplace.json`, `.claude-plugin/plugin.json`, `.codex-plugin/plugin.json`, `cli/spectacular`, `skills/spectacular/SKILL.md`)
- CHANGELOG.md `[1.7.0]` entry
