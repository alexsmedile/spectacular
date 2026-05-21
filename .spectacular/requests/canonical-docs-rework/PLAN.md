---
status: verified
priority: high
owner: alex
updated: 2026-05-21
summary: "Refactor PRD.md into 4 focused root docs (PRD, PRINCIPLES, ARCHITECTURE, ROADMAP) + rewrite AGENTS.md as the .spectacular/ onboarding doc"
related:
  - ../../PRD.md
  - ../../PRINCIPLES.md
  - ../../ARCHITECTURE.md
  - ../../ROADMAP.md
  - ../../AGENTS.md
  - ../../DECISIONS.md
  - ../prd-craft/PLAN.md
---

# Plan — Canonical Docs Rework

## Goal

Shrink `.spectacular/PRD.md` from 896 lines to ~200 lines by extracting architecture, principles, and roadmap content into focused companion docs. Rewrite `AGENTS.md` as the in-folder onboarding doc for any agent working inside `.spectacular/`. Eat our own dogfood: the new PRD must pass the `prd review` gate, and the final set must obey principle 3 (small files over giant documents).

## Why (intent)

`PRD.md` currently contains four documents wearing one filename: the actual PRD, the repo architecture spec, a roadmap, and skill interaction details. This violates principles 1 (context is the system — load only what's needed) and 3 (small files). Anyone planning ends up loading 550 lines of architecture they don't need; anyone needing structure detail has to scroll past goals and non-goals. Splitting also unblocks future work on agent spec, capabilities, and runtime-enforced principles — none of which can land cleanly while PRD is the catch-all.

## Constraints

- Zero context loss — every section in current PRD.md has a defined home in the new structure (see routing map below).
- Markdown-only, no new tooling, no scripts.
- All canonical docs follow the versioning rule: snapshot before overwrite (`PRD@v1.3.md` etc.).
- No new files beyond the 3 planned (PRINCIPLES, ARCHITECTURE, ROADMAP) — CONVENTIONS.md and LIFECYCLE.md fold into ARCHITECTURE.md.
- AGENT-SPEC.md and CAPABILITIES.md explicitly deferred — not in this request.
- Each new doc gets a `version: 1.0` frontmatter and the canonical Spectacular root-file frontmatter shape.

## Milestones

1. **PRD reshaped** — `PRD.md` reduced to ~200 lines, 6-slot PRD shape (problem, who, success, non-goals, constraints, milestone), passes `prd review`.
2. **PRINCIPLES.md live** — 8 principles (5 existing + 3 new) with operational-hook sub-bullets per principle.
3. **ARCHITECTURE.md live** — all "Layer" sections + frontmatter schema + versioning + lifecycle, in one focused doc.
4. **ROADMAP.md live** — v2/future content extracted, time-ordered.
5. **AGENTS.md rewritten** — onboarding-focused, includes context-loading rules per task type.
6. **Cross-links sane** — PRD's "Related docs" section points to all companion docs; each companion doc points back to PRD; AGENTS.md cites all the above.

## Tasks

See [TASKS.md](TASKS.md).

## Dependencies

- `prd-craft/` request — defines the review gate this PRD must pass. Should run that gate after milestone 1.
- No external dependencies. All edits inside `.spectacular/`.

## Validation

| Milestone | How we verify it passed |
|---|---|
| 1. PRD reshaped | Line count <250; `prd review` passes; routing map cross-checked (every old section has a home) |
| 2. PRINCIPLES.md | All 8 principles present; each has an enforcement hook; PRD links to it |
| 3. ARCHITECTURE.md | All old "Layer" sections present; frontmatter schema present; versioning + lifecycle present; PRD/AGENTS link to it |
| 4. ROADMAP.md | v2 workspaces, nested workspaces, workflows, future-vision bullets all present; ordered by version |
| 5. AGENTS.md | Self-contained for an agent landing cold in `.spectacular/`; context-loading rules per task; no duplicated content from PRD |
| 6. Cross-links | Each new doc has a "Related" frontmatter and tail section; no broken `[[]]` or relative links |

## Deliverables

- `.spectacular/PRD.md` — rewritten (snapshot saved as `PRD@v1.3.md`)
- `.spectacular/PRINCIPLES.md` — new
- `.spectacular/ARCHITECTURE.md` — new
- `.spectacular/ROADMAP.md` — new
- `.spectacular/AGENTS.md` — rewritten (snapshot saved as `AGENTS@v1.0.md`)
- Updated DECISIONS.md entry documenting this rework
- Routing map preserved in this PLAN as the authoritative checklist

## Routing map (authoritative)

Every section of current `PRD.md` v1.3 has exactly one new home. This is the contract: zero context loss.

| Current PRD section | New home | Notes |
|---|---|---|
| Vision | PRD § Vision | Keep |
| Deliverable (3 layers) | PRD § Deliverable | Keep |
| Core Principles 1–5 | PRINCIPLES § 1–5 | Move; add enforcement hooks |
| Goals | PRD § Goals & success criteria | Merge with measurable success |
| Non-Goals | PRD § Non-goals | Keep |
| Target Users | PRD § Target users | Keep |
| Repository Architecture (tree) | ARCHITECTURE § Layout | Move |
| Root Layer | ARCHITECTURE § Root layer | Move |
| AGENTS.md description | ARCHITECTURE § Root layer + AGENTS.md itself | Description in ARCH; AGENTS.md becomes the doc |
| STACK.md example | ARCHITECTURE § Root layer (mention only) | Example removed — belongs to host project |
| DECISIONS.md example | ARCHITECTURE § Root layer (mention only) | Same |
| config.yaml schema | ARCHITECTURE § Configuration | Move |
| Frontmatter Schema | ARCHITECTURE § Frontmatter conventions | Move (folded — no CONVENTIONS.md) |
| Ideas Layer | ARCHITECTURE § Ideas layer | Move |
| Current Specs Layer | ARCHITECTURE § Current layer | Move |
| Requests Layer | ARCHITECTURE § Requests layer | Move |
| Request Files (PLAN/TASKS/SESSION/RISKS/VERIFY) | ARCHITECTURE § Request files | Move |
| PRD vs PLAN distinction | ARCHITECTURE § Request files (top) + PRD § Deliverable (1-liner) | Primary in ARCH |
| Skills Layer | ARCHITECTURE § Skills layer | Move |
| Workflows Layer (v2 placeholder) | ROADMAP § v2 — Workflows | Move |
| Memory Layer | ARCHITECTURE § Memory layer | Move |
| Archive Layer | ARCHITECTURE § Archive layer | Move |
| Versioning | ARCHITECTURE § Versioning | Move (folded — no CONVENTIONS.md) |
| Lifecycle Model | ARCHITECTURE § Lifecycle | Move (folded — no LIFECYCLE.md) |
| Skill Interaction Model | `skills/spectacular/SKILL.md` or `references/skill-interaction.md` | Out of PRD entirely — it's runtime, not product |
| Init Flow (CLI) | `requests/cli-bootstrap/PLAN.md` + ARCHITECTURE § Init flow (brief) | Detail in CLI request; ARCH keeps 5-line summary |
| Agent Design Principles | PRINCIPLES § Agent principles | Move |
| Retrieval Principles (context layers + retrieval order) | AGENTS.md | This is exactly what AGENTS.md is for |
| Anti-Entropy Rules + "Most important insight" | PRINCIPLES § Anti-entropy + § Closing insight | Move |
| Future Vision (v2+) | ROADMAP § v2 overview | Move |
| v2 — Workspaces | ROADMAP § v2 — Workspaces | Move |
| v2 — Nested Workspaces | ROADMAP § v2 — Nested Workspaces | Move |

## Target shape per doc

**PRD.md** (~200 lines): Vision → Problem (NEW) → Target users → Deliverable → Goals & success criteria → Non-goals → Constraints (NEW) → First milestone (NEW) → Principles (summary) → Related docs.

**PRINCIPLES.md** (~150 lines): 5 existing principles + 3 new (progressive disclosure / three-layer model / humans decide) + Agent principles + Anti-entropy + Closing insight. Each principle gets a `How the skill enforces this:` sub-bullet.

**ARCHITECTURE.md** (~400–500 lines): The .spectacular/ tree → Root layer (file roles + config.yaml schema) → Frontmatter conventions → Ideas / Current / Requests / Skills / Memory / Archive layers → Request files (PLAN/TASKS/SESSION/RISKS/VERIFY + PRD-vs-PLAN) → Lifecycle → Versioning → Init flow summary.

**ROADMAP.md** (~100 lines): v1 (current) → v2 workflows / workspaces / nested workspaces / multi-agent / hook automation → v3+ context orchestration / RepoOS.

**AGENTS.md** (~80 lines, rewrite): What this folder is → How to operate → Context loading by task type → Available skills → Creating requests → Don'ts.

## Approach

Sequential, doc-by-doc, snapshot-before-edit on every canonical change. No parallel edits — each milestone's review must pass before the next starts. Mirrors the principle-7 (intent → execution → validation) flow internally for each milestone.

## Open questions

- Should `Skill Interaction Model` (old § 20) live in `skills/spectacular/SKILL.md` body, or in a new `references/skill-interaction.md`? Lean: SKILL.md is already lean orchestration; create the reference doc.
- Does `STACK.md` for this repo need updating in the same pass, or leave for a follow-up? Lean: leave; it's small and accurate.
- Should `DECISIONS.md` get a new entry *as part of* this request, or after? Lean: as part of milestone 6 — the rework itself is a decision worth logging.
