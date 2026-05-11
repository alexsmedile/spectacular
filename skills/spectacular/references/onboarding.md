# Onboarding — First Invocation on Existing Project

Triggered by: first time `/spectacular` is invoked in a project that already has a `.spectacular/` directory.

---

## Goal

Orient quickly. Don't overwhelm. Load the minimum needed to give an intelligent first briefing.

---

## Onboarding sequence

1. Read `config.yaml` — get project name and config
2. Read root layer frontmatter (PRD.md, STACK.md, DECISIONS.md, AGENTS.md summaries)
3. Read all `requests/*/PLAN.md` frontmatter — understand active/planned work
4. Read `current/` file list and frontmatter summaries — understand system state
5. Check `memory/` file list — note if anything exists
6. Produce a briefing (same format as status.md)

---

## What to flag on first look

| Observation | Action |
|---|---|
| No `requests/` directory or empty | Note it — project may be pre-work |
| `current/` has draft specs with no active request | Surface them |
| Root files missing frontmatter | Note it, offer to add |
| `config.yaml` missing | Offer to create from template |
| `AGENTS.md` missing | Offer to create from template |
| Requests with `active` status but no SESSION.md | Flag as setup gap |
| Very old `updated` dates | Flag as potentially stale |

Surface a max of **2-3 observations**. Don't list every gap — pick the most actionable ones.

---

## Tone

Treat this as taking over from someone mid-project. Be oriented, not confused. Ask one clarifying question if needed, not five.

Example first briefing:

> "Project: **my-app** — AI-native workspace manager.
>
> | Layer    | Items                                          |
> |----------|------------------------------------------------|
> | Current  | auth (stable), billing (draft)                 |
> | Requests | add-team-billing (active), dark-mode (planned) |
> | Memory   | 1 lesson                                       |
>
> `add-team-billing` is active with a SESSION.md — looks like billing implementation is in progress.
> `current/billing` is still `draft`.
>
> Want to pick up where billing left off, or start something new?"
