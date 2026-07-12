---
description: No-arg invocation — read state, build briefing, surface the next action.
when_to_use: /spectacular with no arguments.
---

# Status — No-arg invocation

Triggered by: `/spectacular` with no arguments, or `spectacular status`.

---

## Substrate check (auto-invoked on failure)

If any of the steps below fail to read or parse — `config.yaml` malformed, root doc frontmatter unparseable, `doc-index.md` not loadable — **do not silently proceed with partial state**. Instead:

1. Auto-run `spectacular doctor workspace frontmatter kits` (the relevant subset)
2. Surface the doctor findings inline before the briefing
3. Suggest `spectacular doctor --fix` (mechanical) and/or `/spectacular doctor --fix` (judgment)
4. If errors block reading enough state to brief, abort with the doctor report; don't fabricate

See [[doctor-substrate]] for the full auto-invocation table.

## Empty workspace → usher, don't brief

If the read below finds **no requests** and the canonical docs are still template stubs (a fresh `spectacular init`), don't print an empty briefing — route to [[guided-first-run]], which walks the user into their first request one step at a time.

## Steps

1. Read `.spectacular/config.yaml` for project name and config.
2. Read frontmatter from all root layer files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `roadmaps/index.md`, `AGENTS.md`, `STACK.md`, `decisions/index.md`. Older workspaces may only have `PRD.md` + `STACK.md` + `decisions/index.md` + `AGENTS.md` — treat the four newer docs as optional.
3. Read `.spectacular/specs/index.md` (top-level index).
4. Run `spectacular status --json` to get the request fleet — one object per request with frontmatter (status/priority/build/updated/summary) **plus** grep-safe body signals (goal line, `progress: {done, open, deferred, total}`, `current_milestone`). This is the deterministic mechanical read; you layer judgment (signal detection below) on top. Don't hand-parse `requests/*/PLAN.md` frontmatter yourself — the CLI owns that extraction (status-fleet-view, b23). Fall back to reading the files only if the CLI is unavailable.
5. Read `.spectacular/memories/` file list and rough counts (not full content).
6. Run `spectacular doctor specs` and capture any specs/index.md drift warning (the CLI compares specs/index.md's `updated` against the newest archived request — see signal table below). This is the one archive read you delegate to the CLI; you still don't read `archive/` bodies yourself.

Do NOT read `archive/` unless explicitly asked.

---

## Build the briefing

Output a conversational briefing with a minimal embedded table:

```
Project state as of <date>:

| Layer    | Items                                               |
|----------|-----------------------------------------------------|
| Current  | auth (stable), billing (draft)                      |
| Requests | add-team-billing (review), dark-mode (planned)      |
| Memory   | 2 lessons, 1 trap                                   |

<Single sentence identifying the highest-priority next action.>

What would you like to work on?
```

---

## Signal detection — surface proactively

While reading state, flag any of the following:

| Signal | Proactive proposal |
|---|---|
| All TASKS.md items checked, status still `active` | Propose moving to `review` |
| Status is `review` | Offer to run the verification walk against the resolved artifact — VERIFY.md if present, else TASKS `### Verification` group, else PLAN § Validation (see [[lifecycle]] § Artifact detection). Never offer to *create* VERIFY.md here — the file is opt-in via the 2-of-6 rule at scaffold/grill time. |
| Request `updated` date > 14 days ago, status `active` | Flag as potentially stale |
| `specs/` capability is `draft`, no active request | Suggest creating a request or promoting to stable |
| specs/index.md drift (from `spectacular doctor specs`) | The CLI computes this — run `spectacular doctor specs` and relay any "specs/index.md … may be stale" warning. Offer to run spec-sync against the named archive. Don't re-derive the date math here; the CLI owns it. See [[spec-sync]]. |
| `requests/` has slug collision potential | Warn |

Only surface the highest-signal item in the briefing. Offer to show full details.

---

## Priority ranking (for "highest-priority next action")

1. `review` status requests with all tasks done — ready to verify or archive
2. `active` requests with blockers noted in SESSION.md
3. `planned` requests with no SESSION.md yet
4. Stale/stuck items

Pick one. State it clearly. Ask what the user wants to do.
