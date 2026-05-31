---
description: Continue active work — session state, task tracking, milestone progress.
when_to_use: Actively working on a request that already exists.
---

# Active Request — Continuing Work

Triggered by: skill detects user is actively implementing a request, or user asks to continue work on a specific slug.

> **@Implementation policy gate.** Before moving a request `planned → active`, run `spectacular policy @Implementation` and follow every active policy. The default blocker is `understand-before-change`: PLAN.md must have a filled `## Understanding` section (or a `UNDERSTANDING.md` with the same three subheads) — satisfy it or stop. See [policy-injection.md](policy-injection.md).

---

## On entering active state

When a request transitions from `planned` → `active`:

1. Update `PLAN.md` frontmatter: `status: active`, `updated: <today>`
2. Create `SESSION.md` if it doesn't exist (see template below)
3. Load the full request context: `PLAN.md`, `TASKS.md`, `SESSION.md`, `SPEC.md`, relevant `specs/<capability>/SPEC.md`

---

## SESSION.md template

```md
---
updated: <today>
---

# Session — <slug>

## Current state
<What's been done so far>

## Active task
<What's being worked on right now>

## Blockers
<Anything blocking progress>

## Next actions
- 
- 
```

SESSION.md is **committed to git** — it's the team's operational record of in-progress work.

Update SESSION.md at natural breakpoints: after a meaningful chunk of work, when blocked, or when handing off.

---

## Task tracking

Monitor `TASKS.md` for completion signals:

- When a task item is checked off, acknowledge it
- When **all items in a group** are checked, note it proactively
- When **all items in TASKS.md** are checked → propose moving to `review` (see `lifecycle.md`)

---

## Context loading during active work

Load by task type:

| Task type | Load |
|---|---|
| Planning/design | `PLAN.md`, `PRD.md`, `DECISIONS.md`, `SPEC.md`, relevant `specs/<capability>/SPEC.md` |
| Implementation | `STACK.md`, `PLAN.md`, `TASKS.md`, local capability specs |
| Review/QA | `VERIFY.md`, capability specs, `RISKS.md` |

Always prefer loading targeted per-capability files over the full `specs/` tree. The top-level `SPEC.md` is cheap and always relevant.

---

## Handling blockers

When a blocker is encountered:
1. Note it in `SESSION.md` under `## Blockers`
2. Surface it in the next briefing
3. Propose resolution path if possible
