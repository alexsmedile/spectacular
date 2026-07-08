---
status: archived
priority: medium
owner: alex
updated: 2026-06-28
build: b13
summary: "Reduce skill-reference doc sprawl: audit the 18 <doc>-rules.md files for empty/duplicated bodies (thin vs fold vs write), and collapse the verify-doc trio (verification.md + verify.md + verify-tests.md) into one."
related:
  - ../../ARCHITECTURE.md
  - ../../specs/index.md
archived: 2026-06-28
---

# Plan — rules-files-audit

> **Origin (2026-06-27):** `<doc>-rules.md` files carry two things: frontmatter
> (the machine dispatch the generic grill/refine/review engine reads) and a body
> (per-doc prompts + gate checks). The frontmatter is always load-bearing. But
> 6 of 18 files are `mode: stub` with near-identical ~21-line bodies whose only
> content is boilerplate ("grill → no-op + hint, refine → rewrite, review →
> structural"). That boilerplate is copy-pasted, drifts independently, and adds
> files to the registry surface for zero per-doc information. The empty-body /
> same-body case is the problem to investigate and resolve here.

## 1. Goal

Decide and apply the right treatment for low-content rules files so that frontmatter (load-bearing) is preserved while duplicated/empty bodies stop being maintained per-file.

## 2. Constraints

- **Frontmatter must stay per-file.** The engine resolves dispatch from each `<doc>-rules.md`'s frontmatter (mode, location, template, slots…). Removing a file removes its dispatch — not an option for any registered doc.
- **No behavior change for users.** Stub docs (`grill` = polite no-op, etc.) must behave identically after the change.
- **doc-index.md is the catalog, not the registry.** If shared boilerplate moves there, it documents the default stub behavior once — it does not become a second dispatch source.
- **Investigation first, edits second.** This request decides the approach before touching files; the decision is recorded in DECISIONS.md.

## Understanding

### How it works now

18 `<doc>-rules.md` files. By mode:
- **Rich bodies (keep as-is):** prd (266 lines), roadmap (406), pack (223), plan (128), personas (106), vision (130), the index-mode soft-DB docs (decisions/feedback/idea/memory/sessions, 49-99), policy (24, structured).
- **Stub bodies (the problem):** `agents` (23), `architecture` (21), `principles` (21), `spec` (23), `stack` (21), and `tasks` (51, partial). The 5 small stubs share an essentially identical body: a one-line "scaffolded from template, user edits directly" + the same 3-verb behavior list.

The generic engine reads frontmatter for dispatch; the body loads as context when a verb runs. For a stub, that body conveys no per-doc information the default couldn't.

### What changes

To be decided in M1 (options below), then applied:
- **Option A — Thin to frontmatter + 1 pointer line.** Each stub body becomes a single line: "Stub doc — default verb behavior, see `doc-index.md` § stub defaults." Shared boilerplate documented once.
- **Option B — Promote where warranted.** Some "stubs" may deserve a real body (e.g. `tasks` at 51 lines is borderline; `spec` arguably needs review-gate prompts). Audit each: thin it OR write the body it's missing.
- **Likely outcome:** a mix — most stubs thin to A, one or two graduate to B.

### What stays the same

- All frontmatter, all dispatch, all user-facing verb behavior.
- The rich rules files (prd, roadmap, pack, etc.) — untouched.
- The one-file-per-registered-doc invariant (we thin bodies, we don't delete files).

## 3. Milestones

- **M1 — Audit + decision.** Read all 6 stub bodies side by side, confirm what's truly shared vs doc-specific. Decide per file: thin-to-pointer (A) or write-real-body (B). Record the rule in DECISIONS.md.
- **M2 — Add the shared default to doc-index.md.** One "stub default behavior" section the thinned files point to.
- **M3 — Apply per-file treatment.** Thin the A-files to frontmatter + pointer; write bodies for any B-files identified.
- **M4 — Verify no behavior drift.** Confirm the engine still dispatches each doc correctly; `spectacular doctor docs` passes.
- **M5 — Collapse the verify-doc trio.** `verification.md` (where checks live — 2-of-6 rule), `verify.md` (the validation walk), and `verify-tests.md` (when to script checks) are three files for one concept users can't disambiguate. Merge into `verify.md` with clearly-labelled sections; the 2-of-6 rule and the script-vs-checklist guidance become subsections, not standalone files. Update SKILL.md routing (currently routes to all three) to point at the one file. Same skill-ref-sprawl reduction as the stub-body work — bundled here.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None blocking. Adjacent to [[onboarding-dedup]] (same "reference shared content instead of duplicating" principle) but independent.

## 6. Validation

- M1 — DECISIONS.md entry records the approach + per-file dispositions.
- M3 — each thinned file still has complete frontmatter (engine load works); diff shows only body changes.
- M4 — `spectacular <stubdoc> grill/refine/review` behaves identically pre/post; `doctor docs` clean.
- M5 — `verification.md` + `verify-tests.md` no longer exist as standalone files; their content is sectioned inside `verify.md`; SKILL.md verification routing points at the one file; the verify walk still runs identically.

## 7. Deliverables

- Decision in DECISIONS.md (stub-body policy).
- `doc-index.md` § stub-default-behavior.
- Thinned (or promoted) bodies for the 6 affected rules files.
- Single merged `verify.md` (absorbing verification.md + verify-tests.md); SKILL.md routing updated.
- SPEC sync if the rules-file contract changes shape.
