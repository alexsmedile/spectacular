# Status — No-arg invocation

Triggered by: `/spectacular` with no arguments, or `spectacular status`.

---

## Steps

1. Read `.spectacular/config.yaml` for project name and config.
2. Read frontmatter from all root layer files (`PRD.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`).
3. Read frontmatter from all `current/<capability>.md` files (or `current/**/` subdirs — read the top-level index if present).
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
| `current/` capability is `draft`, no active request | Suggest creating a request or promoting to stable |
| `requests/` has slug collision potential | Warn |

Only surface the highest-signal item in the briefing. Offer to show full details.

---

## Priority ranking (for "highest-priority next action")

1. `review` status requests with all tasks done — ready to verify or archive
2. `active` requests with blockers noted in SESSION.md
3. `planned` requests with no SESSION.md yet
4. Stale/stuck items

Pick one. State it clearly. Ask what the user wants to do.
