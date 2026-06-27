---
description: First-invocation walk on an existing .spectacular/ project.
when_to_use: Opening /spectacular on a project for the first time.
---

# Onboarding — First Invocation on Existing Project

Triggered by: first time `/spectacular` is invoked in a project that already has a `.spectacular/` directory.

Onboarding **is** a status briefing with a different posture — you're taking over an existing project, not checking in on your own. So it runs the same read + briefing flow as [[status]], plus the onboarding-specific deltas below. The shared sequence lives in status.md and is not restated here (one source of truth).

---

## Run the [[status]] flow, with these deltas

1. **Substrate check runs *unconditionally*** — not only on failure. Before the briefing, auto-run `spectacular doctor workspace frontmatter` to confirm the workspace is known-good. Surface findings first: errors → recommend `spectacular doctor --fix` / `/spectacular doctor --fix` before proceeding; warnings are noted, non-blocking. (status.md runs this *only* when a read fails; onboarding always runs it.)

2. **Read sequence + briefing format** — exactly status.md § Steps and § Build the briefing. Read frontmatter only: config → root canonical docs → `requests/*/PLAN.md` → SPEC.md (+ capability specs) → `memory/` file list → briefing. Don't read `archive/` bodies.

3. **Use the takeover tone** (see below) and surface the first-look observations table (below) instead of status.md's signal-detection table. Both pick the most actionable items; onboarding frames them as "what a new owner should know."

4. **Backwards compatibility:** older `.spectacular/` workspaces may lack `PRINCIPLES.md`, `ARCHITECTURE.md`, or `ROADMAP.md` (pre-PRD-split). Treat each as optional — fall back to `PRD.md` only, and surface it: *"This workspace predates the PRD split — offer to refactor?"*

---

## What to flag on first look

Onboarding's equivalent of status.md's signal table — same "surface the actionable" rule, framed for a new owner.

| Observation | Action |
|---|---|
| No `requests/` directory or empty | Note it — project may be pre-work (consider the guided first-run below) |
| `specs/` has draft capability specs with no active request | Surface them |
| Root files missing frontmatter | Note it, offer to add |
| `config.yaml` / `AGENTS.md` missing | Offer to create from template |
| `PRINCIPLES.md` / `ARCHITECTURE.md` / `ROADMAP.md` missing while `PRD.md` is large (>500 lines) | Pre-split workspace — offer canonical-docs-rework |
| Requests with `active` status but no SESSION.md | Flag as setup gap |
| Very old `updated` dates | Flag as potentially stale |

Surface a max of **2-3 observations**. Don't list every gap — pick the most actionable.

---

## Tone

Treat this as taking over from someone mid-project. Be oriented, not confused. Ask one clarifying question if needed, not five.

Example first briefing (status.md's briefing format + a takeover framing):

> "Project: **my-app** — AI-native workspace manager.
>
> | Layer    | Items                                          |
> |----------|------------------------------------------------|
> | Specs    | auth (stable), billing (draft)                 |
> | Requests | add-team-billing (active), dark-mode (planned) |
> | Memory   | 1 lesson                                       |
>
> `add-team-billing` is active with a SESSION.md — looks like billing implementation is in progress.
> `specs/billing/SPEC.md` is still `draft`.
>
> Want to pick up where billing left off, or start something new?"

---

## Empty / brand-new workspace → guided first-run

If onboarding finds a workspace with **no requests** (fresh `spectacular init`, nothing started), don't print an empty briefing — usher the user into [[guided-first-run]], which walks new → PRD grill → first request one step at a time.
