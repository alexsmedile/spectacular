---
title: Workflow
description: The normal Spectacular loop after installation — init, briefing, requests, lifecycle, archive.
section: ""
status: stable
since: 0.1.0
updated: 2026-05-23
---

# Workflow Guide

This guide shows the normal Spectacular loop after installation. The CLI creates the workspace once. The skill operates it during day-to-day work.

---

## 1. Initialize a project

From the project root:

```bash
spectacular init
```

This creates the required workspace files:

```text
.spectacular/
├── PRD.md              # product intent
├── PRINCIPLES.md       # operating principles
├── ARCHITECTURE.md     # workspace structure
├── ROADMAP.md          # versioned future work
├── STACK.md            # host tech choices
├── DECISIONS.md        # ADR log
├── AGENTS.md           # onboarding for agents
├── config.yaml
├── current/
└── requests/
```

It also installs the Spectacular skill into `.agents/skills/spectacular/` and symlinks it for Claude Code at `.claude/skills/spectacular/`.

If you want prompts instead of defaults:

```bash
spectacular init -i
```

---

## 2. Fill the stable project context

Before asking an agent to do serious work, add concise grounding to the seven canonical root docs:

- `.spectacular/PRD.md` — product intent (Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone). Run `spectacular prd` for an interactive 8-slot grill if starting from scratch.
- `.spectacular/PRINCIPLES.md` — operating principles + how the skill enforces each at runtime
- `.spectacular/ARCHITECTURE.md` — the workspace structure itself (frontmatter, lifecycle, versioning)
- `.spectacular/ROADMAP.md` — versioned future work (v1 / v2 / v3+)
- `.spectacular/STACK.md` — host project's technology choices and engineering rules
- `.spectacular/DECISIONS.md` — ADR-style log of decisions and tradeoffs
- `.spectacular/AGENTS.md` — onboarding doc for any agent landing in `.spectacular/`; defines context loading by task type

These files should stay short and focused. They are not a wiki — their job is to keep agents oriented and let each one load only what the current task needs (progressive disclosure).

---

## 3. Open the workspace briefing

In Claude Code, Codex, or another agent that can load the skill, run:

```text
/spectacular
```

The skill reads the workspace state and reports:

- active and planned requests
- draft or deprecated capability specs
- recent operational memory
- stale or blocked state
- the single highest-priority next action

The skill should not dump the whole workspace. It loads context progressively.

---

## 4. Create a request

When you have work to track, tell the agent:

```text
spectacular new add team billing
```

The skill creates:

```text
.spectacular/requests/add-team-billing/
├── PLAN.md
└── TASKS.md
```

`PLAN.md` captures intent and lifecycle state. `TASKS.md` captures executable work.

The request starts as:

```yaml
status: planned
```

Move it to `active` when implementation begins. The skill may create `SESSION.md` to track handoff state across sessions.

---

## 5. Work from the request folder

During implementation, the request folder is the operational center:

- `PLAN.md` answers what is being built and why.
- `TASKS.md` tracks execution.
- `SESSION.md` records current state and handoff notes.
- `RISKS.md` is useful for sensitive work such as auth, billing, data migrations, or security changes.
- `VERIFY.md` is useful when user-visible behavior or regressions need explicit checks.

Keep request docs focused on the request. Do not prematurely rewrite `current/` while work is still in progress.

---

## 6. Use the lifecycle

Requests move through this lifecycle:

```text
planned → active → review → verified → archived
```

State lives in `PLAN.md` frontmatter:

```yaml
---
status: active
priority: high
updated: 2026-05-11
summary: "Add team billing"
---
```

Typical transition signals:

| State | Move when |
|---|---|
| `planned` → `active` | Implementation starts |
| `active` → `review` | `TASKS.md` is complete |
| `review` → `verified` | Verification checks pass |
| `verified` → `archived` | The request is complete and history can be moved out of active work |

The skill can propose transitions, but the human should confirm them.

---

## 7. Update system truth after completion

`SPEC.md` (and any per-capability `specs/<capability>/SPEC.md` files) describe what the system does now. They should be updated after a request changes real behavior.

When archiving a completed request, the skill should propose updates such as:

- add or update a bullet in `SPEC.md`
- create a new `specs/<capability>/SPEC.md` (only when the bullet outgrows one line)
- update an existing capability spec
- change a capability status from `draft` to `stable`
- leave unaffected specs unchanged

Canonical docs and `specs/` files should be snapshotted before edits:

```text
current/billing/plans.md
current/billing/plans@v1.0.md
```

---

## 8. Capture operational memory

Use memory for lessons the team should not rediscover:

```text
spectacular remember this
```

Good memory entries include:

- recurring failure modes
- migration traps
- integration quirks
- architectural lessons
- project-specific debugging patterns

Memory is team-visible and committed to git under `.spectacular/memory/`. Do not use it for personal notes or secrets.

---

## 9. Archive completed requests

After a request is verified:

```text
spectacular archive add-team-billing
```

The skill should:

1. review the request state
2. propose `current/` updates
3. propose memory entries if useful
4. move the request from `requests/` to `archive/`

Archived requests are not deleted. They are also not read during normal `/spectacular` status briefings.

---

## Practical rhythm

For long-running projects, a useful rhythm is:

1. Run `/spectacular` at the start of a session.
2. Work from one active request.
3. Keep `TASKS.md` and `SESSION.md` current.
4. Verify before changing request state to `verified`.
5. Archive completed work.
6. Update `current/` only when behavior has actually changed.
7. Write memory only for lessons with future value.
