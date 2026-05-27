---
status: verified
updated: 2026-05-26
related:
  - PLAN.md
---

# Tasks — ideas-doctype

## M1 — Rules file + template + doc-index

- [x] Write `skills/spectacular/references/idea-rules.md` (model on `feedback-rules.md`)
  - frontmatter: `doc-id: idea`, `mode: index`, `location: .spectacular/ideas/`, `entries-dir: .spectacular/ideas/`, `scope: project-wide`, `template: templates/idea/base.md`, `snapshot-on-edit: false`
  - body: verbs (grill/refine/review), entry frontmatter schema (type, status: parked|exploring|promoted, priority, owner, origin, updated, related), promotion rule (`promoted` → `archive/ideas/<slug>.md`), doctor area summary
- [x] Write `skills/spectacular/templates/idea/base.md` — frontmatter stub + body sections (Hypothesis / Context / Open questions / Promoted-to placeholder)
- [x] Add `idea` row to `references/doc-index.md` under "Project-wide canonical docs"
- [x] Verify generic engine routes `spectacular idea grill` via `idea-rules.md` (no skill code change needed — engine is registry-driven)

## M2 — CLI mutator verbs

- [x] Add `cmd_idea` dispatcher in `cli/spectacular` (model on `cmd_feedback_loop`)
- [x] `spectacular idea new <slug>` — scaffold entry at `.spectacular/ideas/<slug>.md` from template, set `status: parked`, `updated:` today
- [x] `spectacular idea list [--status <state>]` — list entries with status, last-updated date
- [x] `spectacular idea promote <slug>` — invoke existing scaffold-from-idea flow (`new-request.md` lines 147-150), move source to `archive/ideas/<slug>.md`, set its `status: promoted`
- [x] Top-level dispatch: `idea) shift; cmd_idea "$@"; exit $? ;;`
- [x] Update `--help` to list idea verbs

## M3 — Doctor area

- [x] Add `ideas` to `DOC_AREAS`
- [x] Add `ideas` case to `doctor_parse_args`
- [x] Write `check_ideas()` — judgment-only:
  - missing required frontmatter fields → warn
  - `status: exploring` + `updated:` >90 days → warn (stale)
  - `status: promoted` but file still in `.spectacular/ideas/` (not `archive/ideas/`) → warn (orphan)
  - unknown status value → warn
- [x] Register in `run_areas`: `ideas) check_ideas ;;`
- [x] Bump help string from "14 areas" to "15 areas"
- [x] Add `## ideas (v1.7.0+)` section to `references/doctor-areas.md`

## M4 — SKILL.md triggers + cross-refs

- [x] Add trigger rows to SKILL.md "Workspace lifecycle" table:
  - `spectacular idea new <slug>` → CLI verb
  - `spectacular idea list` → CLI verb
  - `spectacular idea promote <slug>` → CLI verb (handoff to `new-request.md`)
  - `spectacular idea grill\|refine\|review` → generic engine via `idea-rules.md`
- [x] Add `idea` to the registered doc-IDs list (currently ends `...feedback, convention-pack, docs-manifest, docs-page`)
- [x] Bump SKILL.md version frontmatter to 1.7.0
- [x] Add cross-reference to `[[idea-rules]]` from `.spectacular/ARCHITECTURE.md` § ideas/

## M5 — Dogfood + release

- [x] `spectacular idea new test-idea-foo` → verify file shape
- [x] `spectacular doctor ideas` → expect clean on existing `memory-protocols.md` + `test-idea-foo.md`
- [x] Backdate `test-idea-foo.md` `updated:` to 100 days ago, set `status: exploring` → doctor warns
- [x] `spectacular idea promote test-idea-foo` → request scaffolded, file moved to `archive/ideas/`
- [x] Doctor clean after promotion
- [x] Delete the throwaway test request
- [x] Bump all version manifests to 1.7.0; run `check-manifests.sh` audit
- [x] CHANGELOG.md `[1.7.0]` entry — Added/Changed sections
- [x] Flip PLAN.md `status: review`

## Deferred

- [ ] `spectacular idea refine <slug>` interactive (M1 ships engine wiring only — explicit refine flow can wait)
- [ ] Auto-tagging idea entries with PRINCIPLES.md or current-session links (out of scope)
- [ ] `spectacular idea park-from-request <slug>` — the "abandon a request into an idea" flow done manually 2026-05-26; explicitly out of scope per user ("very particular case, not worth documenting")
