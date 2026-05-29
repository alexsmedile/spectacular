---
status: planned
priority: high
owner: alex
updated: 2026-05-29
summary: "Skill-side /spectacular verify <slug> walk — guides each VERIFY.md check interactively, flips review→verified, closes Principle 7"
related:
  - PRD.md
  - ../../PRINCIPLES.md
  - ../../ROADMAP.md
target_version: v1.11.0
---

# Plan — verify-walk

## 1. Goal

Add a first-class skill-side validation walk, `/spectacular verify <slug>`, that reads a request's VERIFY.md (or PLAN § Validation), guides the user through each check interactively, captures evidence, and flips `status: review → verified` — turning verification from a static checklist into an executed ritual (closes Principle 7).

## 2. Constraints

- **Skill-only — no CLI counterpart.** Verification is judgment work, not mechanical mutation. The CLI must redirect `verify` to the skill with a friendly message (same pattern as `grill`/`refine`).
- **Reuses existing lifecycle.** State lives only in PLAN.md frontmatter; the walk flips `review → verified` via the existing `promote` verb, never by editing frontmatter directly.
- **Respects the 2-of-6 verification rule.** When no VERIFY.md exists (folded into PLAN § Validation), the walk reads the PLAN section instead — it must handle both shapes. See [[verification]].
- **No fabrication of evidence.** If a check can't be confirmed, it records a blocker, not a pass.

## 3. Milestones

- M1 — `references/verify.md` walk algorithm exists: reads VERIFY.md or PLAN § Validation, iterates each check, prompts for evidence, records pass/blocker.
- M2 — Walk outcome wired to lifecycle: all-pass flips `review → verified` (via `promote`); any blocker leaves `review` and records the blocker list.
- M3 — Retrospective + archive tie-in: optional end-of-walk "what surprised you vs the PRD/PLAN?" prompt captured to `memory/`; `spectacular archive` warns if `verified` was never reached via the walk.
- M4 — Surfaced + documented: SKILL.md routing table updated; `docs/commands.md` agentic-verbs section covers `verify`; CLI redirect message added.
- M5 — Dogfood + ship: 1+ real request driven through the walk; CHANGELOG entry; plugin bump to v1.11.0.

## 4. Tasks

See `TASKS.md`.

## 5. Dependencies

- None blocking. Builds on shipped lifecycle (`promote`), the 2-of-6 [[verification]] rule, and the soft-DB `memory/` substrate (for the retrospective capture). Independent of the other runway requests.

## 6. Validation

- M1 — Running the walk on a request with a VERIFY.md iterates every check; running on one without VERIFY.md falls back to PLAN § Validation cleanly.
- M2 — An all-pass walk lands the request at `verified`; a walk with one unmet check leaves it at `review` with the blocker recorded.
- M3 — Retrospective answer appears as a `memory/` entry; archiving an un-walked request emits the warning.
- M4 — `verify` appears in SKILL.md routing + docs; CLI `spectacular verify <slug>` prints the skill-redirect message.
- M5 — A real request's PLAN shows `verified` reached through the walk; CHANGELOG + manifests at v1.11.0.

## 7. Deliverables

- `skills/spectacular/references/verify.md` — the walk algorithm + gate
- SKILL.md routing-table entry for `verify`
- CLI redirect for `spectacular verify <slug>` (skill-only message)
- `docs/commands.md` update (agentic-verbs section)
- Archive-flow warning when `verified` was never reached via the walk
- CHANGELOG [1.11.0] entry
