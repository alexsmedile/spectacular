---
version: 1.0
updated: <DATE>
summary: "Onboarding doc for any agent or human landing inside .spectacular/"
related:
  - PRD.md
  - SPEC.md
  - PRINCIPLES.md
  - ARCHITECTURE.md
---

# Working in `.spectacular/`

This is the onboarding doc for any agent or human landing inside this workspace. Read this first.

<!--
  AGENTS.md governs THIS folder (.spectacular/).
  Repo-root AGENTS.md / CLAUDE.md govern the whole repository.
  This file is mode: stub — edit directly, no slot grill.
-->

## What this folder is

`.spectacular/` is an AI-native operational workspace. It separates four things deliberately:

| Layer | Files | Purpose |
|---|---|---|
| Intent | `PRD.md` | What we want and why |
| Truth | `SPEC.md` + `specs/` | What's built right now (SPEC.md is the always-on index) |
| Work | `requests/<slug>/` | What's changing next (PLAN.md owns lifecycle state) |
| Memory | `memory/` | What we learned, kept across sessions |

Optional canonical docs (`PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `STACK.md`, `DECISIONS.md`) are scaffolded by kits — not every project needs them.

## How to operate

1. **Read frontmatter first**, file bodies second. `status:`, `summary:`, `version:`, `updated:`, `related:` are the navigation layer.
2. **Load progressively.** The top-level `SPEC.md` is cheap and always relevant. Don't pre-load all of `specs/` or `requests/` — load only what the current task needs (see Context loading below).
3. **Snapshot before overwrite** on any canonical doc — root layer, `SPEC.md`, `specs/<capability>/SPEC.md`, `config.yaml`. Naming: `<FILE>@vN.md`. The unversioned filename always points to current. Use `spectacular snapshot <file>` to do it.
4. **Propose, don't act** on irreversibles — archives, lifecycle promotions, memory writes, bulk edits >5 files. Confirm with the human first.
5. **Never read `archive/`** during normal operation. It's write-only from your perspective.
6. **Write to `memory/` only on explicit confirmation** via `spectacular remember this`. Never autonomously.

## Context loading by task

Load the minimum needed for the task — don't slurp the whole workspace.

| Task type | Load |
|---|---|
| Cold onboarding | `PRD.md`, `SPEC.md`, this file |
| Planning a new request | `PRD.md`, `SPEC.md`, `PRINCIPLES.md` (if present), `DECISIONS.md` (if present) |
| Refining intent / PRD work | `PRD.md` + relevant skill refs |
| Implementing a request | `STACK.md` (if present), `requests/<slug>/PLAN.md`, `TASKS.md`, `SPEC.md`, the specific `specs/<capability>/SPEC.md` files the work touches |
| Reviewing / QA | `requests/<slug>/VERIFY.md` (if present), affected `specs/<capability>/SPEC.md`, `RISKS.md` (if present) |
| Structural questions | `ARCHITECTURE.md` only (if present) |

## Task tracking — two layers

| Layer | Tool | Use for |
|---|---|---|
| Session micro-steps | Harness `TaskCreate` / `TaskUpdate` | Ephemeral per-session steps — drives the CLI live progress UI |
| Persistent milestones | On-disk `requests/<slug>/TASKS.md` | Stable milestone blocks that survive across sessions |

Anti-pattern: 1:1 duplication between the two layers. Use the harness for the granular checklist *within* a milestone block; use TASKS.md for the milestones themselves.

## Available skills

<!-- List skills available in this project. -->

- `/spectacular` — workspace orchestration (briefings, requests, lifecycle, archive, memory, doc verbs)
- <OTHER SKILL>

## Creating requests

```bash
spectacular new <description>
```

The skill derives a kebab-case slug, applies the verification 2-of-6 rule to decide whether to scaffold a `VERIFY.md`, and asks for confirmation before writing.

> **Anti-pattern:** never create `requests/<slug>/PRD.md`. Product intent is project-wide and lives at `.spectacular/PRD.md`. Per-request `PLAN.md` references it.

## Don'ts

- Don't touch `archive/` during normal operation
- Don't overwrite canonical docs in place — snapshot first
- Don't write to `memory/` autonomously
- Don't create per-request PRDs
- Don't commit `.spectacular.local/` (it's personal + gitignored)

## Handoff conventions

<!-- How agents hand off mid-task to another session or another agent. -->

- Update `requests/<slug>/SESSION.md` with current state + next action when stepping away mid-work
- Surface blockers in the active PLAN.md `Dependencies` slot
- If a decision was made that future sessions need to know about, propose a `DECISIONS.md` entry
