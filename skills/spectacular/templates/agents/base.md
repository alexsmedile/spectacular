---
version: 1.0
updated: <DATE>
summary: "Onboarding doc for agents working in .spectacular/"
related:
  - PRD.md
  - ARCHITECTURE.md
  - PRINCIPLES.md
---

# Agents — operating inside `.spectacular/`

<!--
  AGENTS.md is the onboarding doc for any agent landing inside .spectacular/.
  Distinct from the repo-root CLAUDE.md / AGENTS.md (those govern the whole repo).
  This file is mode: freeform — no slot grill, edit directly.
-->

## What this folder is

<One-paragraph description of the project's purpose and current state. The cold-start frame.>

## How to operate

1. Read this file first.
2. Read frontmatter (not bodies) of root canonical docs to build state.
3. Load only the docs the current task needs (see Context loading table below).
4. Never overwrite canonical docs in place — snapshot first.

## Context loading by task type

| Task type | Load |
|---|---|
| Planning / design | PRD.md, PRINCIPLES.md, DECISIONS.md |
| Implementation | STACK.md, ARCHITECTURE.md, current/<capability>, requests/<active-slug>/ |
| Review / QA | VERIFY.md (if present), current/<capability>, RISKS.md |
| Onboarding cold | PRD.md, ARCHITECTURE.md, this file |

Load only the capability spec relevant to the current task, not all of `current/`.

## Available skills

<!-- List skills available in this project. -->

- `/spectacular` — workspace orchestration
- <OTHER SKILL>

## Creating requests

```bash
spectacular new <description>
```

Scaffolds `requests/<slug>/PLAN.md` + `TASKS.md`. Edit PLAN to refine intent, then start working.

## Don'ts

- Never read `archive/` during normal operation.
- Never edit root canonical docs without snapshotting first.
- Never write to `.spectacular/memory/` autonomously — only on explicit `spectacular remember this`.

## Handoff conventions

<!-- How agents hand off mid-task to another session or another agent. -->

- Update `SESSION.md` in the active request folder with current state + next action.
- Surface blockers in PLAN.md `Dependencies` slot.
