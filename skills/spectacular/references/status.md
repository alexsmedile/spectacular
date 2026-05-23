# Status — No-arg invocation

Triggered by: `/spectacular` with no arguments, or `spectacular status`.

---

## Substrate check (auto-invoked on failure)

If any of the steps below fail to read or parse — `config.yaml` malformed, root doc frontmatter unparseable, `doc-registry.md` not loadable — **do not silently proceed with partial state**. Instead:

1. Auto-run `spectacular doctor workspace frontmatter kits` (the relevant subset)
2. Surface the doctor findings inline before the briefing
3. Suggest `spectacular doctor --fix` (mechanical) and/or `/spectacular doctor --fix` (judgment)
4. If errors block reading enough state to brief, abort with the doctor report; don't fabricate

See [[doctor-substrate]] for the full auto-invocation table.

## Steps

1. Read `.spectacular/config.yaml` for project name and config.
2. Read frontmatter from all root layer files: `PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `AGENTS.md`, `STACK.md`, `DECISIONS.md`. Older workspaces may only have `PRD.md` + `STACK.md` + `DECISIONS.md` + `AGENTS.md` — treat the four newer docs as optional.
3. Read `.spectacular/SPEC.md` (top-level index) and any `specs/<capability>/SPEC.md` files (per-capability, optional).
4. Read frontmatter from all `requests/*/PLAN.md` files.
5. Read `.spectacular/memory/` file list and rough counts (not full content).

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
| Status is `review`, VERIFY.md exists | Prompt to run VERIFY checklist |
| Status is `review`, no VERIFY.md | Offer to create VERIFY.md |
| Request `updated` date > 14 days ago, status `active` | Flag as potentially stale |
| `specs/` capability is `draft`, no active request | Suggest creating a request or promoting to stable |
| `requests/` has slug collision potential | Warn |

Only surface the highest-signal item in the briefing. Offer to show full details.

---

## Priority ranking (for "highest-priority next action")

1. `review` status requests with all tasks done — ready to verify or archive
2. `active` requests with blockers noted in SESSION.md
3. `planned` requests with no SESSION.md yet
4. Stale/stuck items

Pick one. State it clearly. Ask what the user wants to do.
