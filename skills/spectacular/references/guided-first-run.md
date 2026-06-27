---
description: Usher a brand-new (empty) .spectacular/ workspace through its first request, one step at a time.
when_to_use: /spectacular on a freshly-initialized workspace with no requests yet.
---

# Guided First-Run — From empty workspace to first request

Triggered when the skill detects a **fresh/empty** workspace: `.spectacular/` exists (init ran) but there are no requests, and the canonical docs are still template stubs. This is the Tier-0 onboarding moment.

> **Distinct from [[onboarding]].** Onboarding orients a new owner on an *existing* project with prior work. Guided first-run is for a *blank slate* — there's nothing to brief, so instead of an empty table the skill ushers the user into starting.

## Core rule — one step at a time

Never dump the verb surface. A new user shouldn't see 30+ verbs on first contact. Reveal exactly one next action, do it, then reveal the next. The path:

```
empty workspace → (offer) describe the project → (optional) PRD grill → first request → point at `spectacular next`
```

## The walk

**Step 0 — recognize the state.** No `requests/`, or `requests/` empty; PRD.md still the template stub. Don't brief — usher.

**Step 1 — orient in one line, offer the first move.**

> "Fresh Spectacular workspace — nothing tracked yet. Want to start by describing what you're building (I'll shape it into a PRD), or jump straight to your first request?"

Two doors, not a menu:
- **Describe the project** → go to Step 2 (PRD grill).
- **First request now** → skip to Step 3.

**Step 2 — PRD grill (optional).** If they describe the project, route to `spectacular prd grill` ([[grill]] + [[prd-overrides]]). Keep it light on a first run — the PRD can deepen later. When the PRD has a goal + intent + a success criterion or two, stop and move on. Don't exhaust the slot list on day one.

**Step 3 — first request.** Route to [[new-request]]: derive a slug, scaffold PLAN + TASKS, confirm. This is the same `spectacular new` flow, just reached by ushering rather than by command.

**Step 4 — hand off to the normal loop.** Once the first request exists, reveal the steady-state next action and stop ushering:

> "You're set. `spectacular next` will always tell you the single highest-priority thing to do — run it whenever you're unsure what's next."

From here the workspace is no longer empty, so future `/spectacular` invocations route to [[status]] (or [[onboarding]] on a new machine), not back here.

## What guided first-run does NOT do

- It does not auto-scaffold a PRD or request without confirmation — every write is user-approved (PRINCIPLES: agents propose, humans decide).
- It does not run the full PRD grill to completion — first contact stays light; depth comes later.
- It does not list the verb surface — one action at a time is the whole point.

## CLI entry

No dedicated flag needed in v1. The skill detects the empty state during [[status]] / [[onboarding]] and routes here; `spectacular next` already ushers an empty workspace toward `spectacular new`. A future `spectacular init --walk` (run the guided flow immediately after init) is an optional convenience, not required — left out until asked.

## Related

- [[onboarding]] — existing-workspace orientation (the non-empty counterpart)
- [[status]] — the steady-state briefing this hands back to
- [[new-request]] — the request scaffold this ushers into
- [[prd-overrides]] — PRD grill rules for the optional Step 2
- [[principles]] — agents propose, humans decide; one thing at a time
