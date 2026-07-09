---
version: 2.0
updated: 2026-05-21
summary: "AI-native operational workspace for software projects — convention, skill, and CLI"
related:
  - PRINCIPLES.md
  - ARCHITECTURE.md
  - roadmaps/index.md
  - AGENTS.md
  - decisions/index.md
  - STACK.md
---

# Spectacular — Product Requirements Document

## Vision

Spectacular is an AI-native operational workspace for software projects.

It helps humans and coding agents maintain coherence across long-running product development by separating strategic context, current system truth, active requests, operational memory, and reusable execution skills.

The framework is designed for modern AI-assisted development where the primary bottleneck is no longer code generation, but context management, decision continuity, and execution clarity.

## Problem

Modern AI-assisted projects lose coherence fast. Agents and humans share the same repo but not the same operating model. Strategic intent (the PRD) drifts apart from current truth (what the code actually does); active work piles up in chat scrollback and ad-hoc markdown; decisions are made and immediately forgotten; memory of past failures lives in nobody's head and nowhere on disk.

Existing tools solve one slice — ticket trackers manage tasks, wikis hold documentation, ADR repos log decisions — but none provide a single workspace where intent, execution, validation, and memory coexist in a form both humans and agents can read, retrieve, and reason about progressively.

Spectacular fills that gap with a lightweight, file-based convention that any project can adopt in minutes and any agent can operate without bespoke configuration.

## Target users

**Primary**
- Solo developers using AI coding agents (Claude Code, Cursor, Cline, Aider)
- Startup teams and AI-native product teams
- Technical founders
- Rapid-iteration teams who outgrow chat-based context

**Secondary**
- Agencies running multiple parallel client projects
- Open-source maintainers coordinating contributors
- Internal platform teams

## Deliverable

Spectacular ships as three complementary layers:

1. **Convention** — a documented `.spectacular/` directory structure and file contract that humans and agents follow. See [ARCHITECTURE.md](ARCHITECTURE.md).
2. **Skill** — a Claude Code slash command (`/spectacular`) that reads state, proposes actions, scaffolds files, manages lifecycle transitions, and writes memory. The primary interface for ongoing use.
3. **CLI** — a one-time bootstrap tool (`spectacular init`) that scaffolds the directory, installs the skill, and generates stub root files. Used once per project.

PRD vs PLAN clarifier: `.spectacular/PRD.md` is project-wide (this file). Per-request work lives in `requests/<slug>/PLAN.md` + `TASKS.md` — requests never carry their own PRD. Full distinction in [ARCHITECTURE.md § Request files](ARCHITECTURE.md).

## Goals & success criteria

**Primary goals**
- Create a scalable AI-native project structure
- Reduce context rot across long-running projects
- Improve agent execution quality through structured retrieval
- Improve human-agent coherence in long-running work
- Separate stable truth from active work
- Enable reusable operational skills

**Success criteria** (measurable, time-boxed)
- Any project can adopt Spectacular by running `spectacular init` in under 2 minutes with zero prompts
- A coding agent landing cold in a Spectacular workspace can produce a correctly-scoped briefing from `.spectacular/AGENTS.md` + frontmatter alone, without reading full file bodies
- 90 days after a project adopts Spectacular, the team can answer "what's active / what's blocked / what was decided last month" without scrolling chat
- The framework remains usable by humans reading raw markdown — no required tooling beyond a text editor

## Non-goals

Spectacular is NOT:
- A ticketing system (Linear, Jira, GitHub Issues do this)
- A project management tool (Asana, Notion do this)
- A documentation wiki (separate concern)
- A replacement for git
- A replacement for human-facing READMEs
- A rigid enterprise requirements framework
- A multi-agent orchestration platform (v2+ may add coordination conventions; orchestration itself is out)

The system must remain lightweight and operational.

## Constraints

- **Markdown-only** — no databases, no servers, no compiled artifacts. Every canonical file is human-readable plain text.
- **No mandatory tooling** — the skill and CLI are accelerators; the convention works without them. A human with `vim` must be able to operate a Spectacular workspace.
- **Git-committed by default** — `.spectacular/` is fully committed; only `.spectacular.local/` is gitignored.
- **Snapshot before overwrite** — canonical documents are never overwritten in place. See [ARCHITECTURE.md § Versioning](ARCHITECTURE.md).
- **Small files over giant documents** — see [PRINCIPLES.md](PRINCIPLES.md) principle 3.
- **Humans decide on irreversibles** — the skill proposes; humans confirm anything destructive, lifecycle-advancing, or canonical. See [PRINCIPLES.md](PRINCIPLES.md) principle 8.

## First milestone

Ship Spectacular v1 such that:
- `spectacular init` scaffolds a working `.spectacular/` directory from a single command (see active request `requests/cli-bootstrap/`)
- The four canonical root docs (PRD / PRINCIPLES / ARCHITECTURE / ROADMAP) are split, focused, and each under their target line count
- The `/spectacular` skill operates the workspace end-to-end: status, new, archive, snapshot, remember, promote, prd
- The PRD-craft flow (`prd / prd refine / prd review`) produces a passing PRD on a fresh project in under 15 minutes (see active request `requests/prd-craft/`)

## Principles (summary)

Full text + enforcement hooks in [PRINCIPLES.md](PRINCIPLES.md).

1. **Context is the system** — load only what's needed, progressively, by task.
2. **Separate intent from truth** — strategic docs and operational docs are different files.
3. **Small files over giant documents** — modular, retrievable, locally reasonable.
4. **Humans and agents share the same workspace** — readable by both.
5. **Operational memory compounds** — lessons preserved across sessions and teams.
6. **Progressive disclosure** — references load on demand, never upfront.
7. **Three layers: intent → execution → validation** — every unit of work passes through all three.
8. **Humans decide, agents propose** — irreversibles require confirmation.

## Related docs

- [PRINCIPLES.md](PRINCIPLES.md) — operating principles + how the skill enforces each
- [ARCHITECTURE.md](ARCHITECTURE.md) — `.spectacular/` structure, layers, request files, lifecycle, versioning, frontmatter
- [roadmaps/index.md](roadmaps/index.md) — v2+ features (workspaces, nested workspaces, multi-agent, workflows)
- [AGENTS.md](AGENTS.md) — how agents operate inside `.spectacular/`; context loading per task type
- [decisions/index.md](decisions/index.md) — architectural decision log
- [STACK.md](STACK.md) — host-project technology choices (distinct from Spectacular's own architecture)
