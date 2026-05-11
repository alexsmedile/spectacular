# New Request — Scaffolding

Triggered by: `spectacular new <description>`, or when conversational context clearly indicates a new piece of work to track.

---

## Autopilot vs interactive

| Context | Behavior |
|---|---|
| Clear description in conversation | Derive slug, draft PLAN.md, show user, confirm, write |
| Thin context | Ask targeted questions (goal? scope? why now?), then scaffold |

Always show the derived slug before writing. User can override before confirmation.

---

## Slug rules

1. Derive from the request description — kebab-case, concise (2–4 words)
2. Check `requests/` for existing slugs — if collision, propose appending `-2` or ask user to rename
3. Apply `naming.prefix` from `config.yaml` if set
4. User can override at any time before write

Example: "add team billing" → `add-team-billing`

---

## Scaffold sequence

On confirmation, create:

```
requests/<slug>/
├── PLAN.md       ← always created
└── TASKS.md      ← always created
```

Create only on demand (skill proposes, user confirms):
- `SESSION.md` — when request moves to `active`
- `RISKS.md` — when request touches auth, billing, migrations, or anything flagged sensitive in STACK.md
- `VERIFY.md` — when request has user-visible behavior changes or high-stakes implementation
- `artifacts/` — when screenshots, benchmarks, or research need storing

---

## PLAN.md template

```md
---
status: planned
priority: medium
owner: 
updated: <today>
summary: "<one-line description>"
related:
  - current/<capability>
---

# <Request title>

## Goal
<What are we trying to achieve?>

## Why
<Why now? What problem does this solve?>

## Scope
- 
- 

## Out of scope
- 

## Approach
<High-level implementation approach>

## Success criteria
- 
```

---

## TASKS.md template

```md
---
updated: <today>
---

# Tasks — <slug>

## <Group name>

- [ ] 
- [ ] 

## <Group name>

- [ ] 
- [ ] 
```

---

## Promoting from an idea

When user runs `spectacular promote <idea-file>`:

1. Read `ideas/<idea-file>.md` for content — use it to pre-fill PLAN.md goal, why, approach
2. Scaffold request as above (show slug, confirm)
3. Move `ideas/<idea-file>.md` → `archive/ideas/<idea-file>.md`
4. Note in PLAN.md: `promoted from ideas/<idea-file>.md`
