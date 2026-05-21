---
version: 1.3
updated: 2026-05-21
summary: "AI-native operational workspace — directory convention, skill, and CLI bootstrap for software projects"
---

# Spectacular — Product Requirements Document

## Vision

Spectacular is an AI-native operational workspace for software projects.

It helps humans and coding agents maintain coherence across long-running product development by separating:

- strategic context
- current system truth
- active requests
- operational memory
- reusable execution skills

The framework is designed for modern AI-assisted development where the primary bottleneck is no longer code generation, but context management, decision continuity, and execution clarity.

---

# Deliverable

Spectacular ships as three complementary layers:

1. **Convention** — a documented directory structure and file contract that humans and agents follow
2. **Skill** — a Claude Code slash command (`/spectacular`) that operates the system: reads state, proposes actions, scaffolds files, manages lifecycle transitions, and writes memory
3. **CLI** — a one-time bootstrap tool (`spectacular init`) that scaffolds the directory, installs the skill, and generates stub root files

The skill is the primary interface for ongoing use. The CLI is used once per project.

---

# Core Principles

## 1. Context is the system

The project structure must optimize:
- retrieval quality
- context layering
- long-running continuity
- low context entropy

Agents should load:
- only relevant information
- progressively
- by capability and task

Never the entire repository.

---

## 2. Separate intent from truth

Strategic documents define:
- why the product exists
- goals
- philosophy
- constraints

Operational documents define:
- current behavior
- implementation truth
- active work

Temporary work should never pollute canonical truth.

---

## 3. Small files over giant documents

Prefer:
- modular capability files
- layered context
- local reasoning

Avoid:
- giant PRDs
- monolithic specs
- mega prompts

---

## 4. Humans and agents share the same workspace

The repository should:
- remain readable by humans
- remain structured for agents
- support retrieval systems
- support automation

---

## 5. Operational memory compounds

The system should preserve:
- lessons
- failures
- architectural traps
- recurring bugs
- implementation patterns

Agents should not repeatedly rediscover solved problems.

---

# Goals

## Primary Goals

- create a scalable AI-native project structure
- reduce context rot
- improve agent execution quality
- improve long-running project coherence
- separate stable truth from active work
- support multi-agent workflows (v2)
- enable reusable operational skills

---

# Non-Goals

Spectacular is NOT:
- a ticketing system
- a project management tool
- a documentation wiki
- a replacement for git
- a replacement for human-facing READMEs
- a rigid enterprise requirements framework

The system should remain lightweight and operational.

---

# Target Users

## Primary Users

- solo developers using AI coding agents
- startup teams
- AI-native product teams
- technical founders
- rapid iteration teams

---

## Secondary Users

- agencies
- open-source maintainers
- internal platform teams

---

# Repository Architecture

```txt
.spectacular/
|
| <root-layer>
│   ├── PRD.md
│   ├── STACK.md
│   ├── DECISIONS.md
│   ├── AGENTS.md
│   └── config.yaml
│
├── ideas/              # exploratory thinking, scratchpad — not acted on by skill
│
├── current/            # canonical system truth, capability specs
│
├── requests/           # active and planned work
│
├── skills/             # project-specific reusable skills
│
├── memory/             # long-term operational learning (git-committed, team-visible)
│
└── archive/            # completed requests, historical snapshots
```

`.spectacular.local/` — personal overrides, always gitignored, never committed.

---

# Root Layer

The root layer contains stable project grounding.

```txt
├── PRD.md              # product/business intent
├── STACK.md            # technologies + architecture
├── PRINCIPLES.md       # optional
├── DECISIONS.md        # optional
├── AGENTS.md           # spectacular-specific agent orchestration rules
└── config.yaml         # machine-readable project config
```

## Purpose

Defines:
- project intent
- engineering philosophy
- stack constraints
- agent context loading rules
- stable architecture decisions

## Rules

- should change infrequently
- should remain concise
- should avoid implementation details
- should optimize retrieval quality
- **never overwritten in place — snapshot before editing** (see Versioning)

---

### AGENTS.md

Spectacular-specific agent orchestration rules. Distinct from repo-root `CLAUDE.md` / `AGENTS.md`.

Defines:
- which context to load for which task type
- which skills are available
- handoff conventions

The skill reads `AGENTS.md` to configure its own behavior. Humans write it; skill proposes updates when new skills or capabilities are added.

---

### STACK.md

Example:

```md
Frontend
- Next.js 15
- Tailwind v4
- Zustand

Backend
- Supabase
- Edge Functions

Rules
- Prefer server actions
- Avoid client-side fetching
```

---

### DECISIONS.md

Example:

```md
# Decisions

## 2026-05-11

Decision:
Use Postgres RLS instead of app-level permissions

Why:
Centralized security logic

Tradeoffs:
Harder local debugging
```

---

### config.yaml

Machine-readable project configuration.

```yaml
project:
  name: my-app
  summary: "One-liner about the project"

naming:
  requests: kebab-case        # enforced by skill on scaffold
  prefix: ""                  # optional prefix/suffix for request slugs

required_files:
  requests:
    - PLAN.md
    - TASKS.md

agents:
  default_context:
    - PRD.md
    - STACK.md
    - DECISIONS.md

skills:
  symlink_on_init: []         # project skills to auto-symlink into .claude/skills/
```

---

# Frontmatter Schema

Frontmatter is the skill's primary signal for reading project state. All canonical documents should include frontmatter.

## Root files (PRD.md, STACK.md, etc.)

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "One-sentence description of this file's purpose"
---
```

## current/ capability specs

```yaml
---
status: stable | draft | deprecated
updated: 2026-05-11
summary: "What this capability does"
---
```

## requests/ PLAN.md (lifecycle state owner)

```yaml
---
status: planned | active | review | verified
priority: high | medium | low
owner: alex
updated: 2026-05-11
summary: "What this request changes"
related:
  - current/auth
---
```

**Rules:**
- `status`, `updated`, `summary` are required
- Other fields are optional; skill warns but does not block on missing optional fields
- `PLAN.md` frontmatter is the single source of lifecycle state for a request
- Capability specs in `current/` and `specs/` track their own state independently

---

# Ideas Layer

The ideas layer is a **thinking scratchpad**. It is not a workflow stage. Nothing in `ideas/` is acted on automatically by the skill.

Use it for:
- raw thoughts
- market observations
- UX experiments
- discarded approaches
- future concepts
- unresolved brainstorming

```txt
ideas/
├── multiplayer-editor.md
├── ai-memory-system.md
└── growth-loops.md
```

## Rules

- low commitment
- speculative
- non-canonical
- skill will **propose** saving unresolved decisions here when a conversation has open branches

## Promotion to request

Ideas are not a required gate. A request can be created directly without an idea file. When an idea is deliberately promoted to a request, the skill scaffolds the request from the idea content and moves the idea file to `archive/ideas/`.

---

# Current Specs Layer

The current specs layer represents the canonical system truth.

```txt
current/
├── auth/
│   ├── login.md
│   ├── sessions.md
│   └── permissions.md
├── billing/
└── editor/
```

## Purpose

Defines:
- current behavior
- active capabilities
- security requirements
- performance expectations
- user-visible behavior

## Rules

- authoritative
- current only
- behavior-oriented
- modular
- **never overwritten in place — skill snapshots before proposing edits**
- skill proposes `current/` updates when a request is archived (human confirms before write)

## Specification structure

Each capability spec should contain:
- purpose
- requirements
- scenarios
- security considerations
- performance expectations

---

# Requests Layer

The requests layer contains proposed or active work.

```txt
requests/add-team-billing/
├── PLAN.md             # required — lifecycle state, goal, approach, success criteria
├── TASKS.md            # required — executable implementation checklist
├── SESSION.md          # created when request moves to active
├── RISKS.md            # proposed by skill for high-risk requests
├── VERIFY.md           # proposed by skill for user-visible or high-stakes changes
├── specs/              # per-request capability specs (track own frontmatter state)
└── artifacts/
       ├── screenshots/
       ├── benchmarks/
       ├── user-feedback/
       └── research/
```

## Slug rules

- skill derives slug from conversation context, shows user before creating
- user can override at any time
- slugs are kebab-case by default (configurable in `config.yaml`)
- slugs are unique — if slug exists, skill proposes appending `-2` or asks user to rename

## Purpose

Track:
- planned changes
- execution work
- temporary reasoning
- implementation state

## Rules

- temporary
- operational
- archived on completion (not deleted)
- `PLAN.md` frontmatter owns lifecycle state

---

# Request Files

## PRD vs PLAN — scope distinction

These are **two different artifacts at two different layers**, not the same artifact at two scopes.

| Artifact | Location | Scope | Answers |
|---|---|---|---|
| `PRD.md` | `.spectacular/` root only | **Product** (whole project) | Why does this product exist? |
| `PLAN.md` | `requests/<slug>/` only | **Request** (one slice of work) | What are we building in this slice and why? |

**Rules:**
- A project has exactly one `PRD.md` (at the root). It's long-lived and snapshot-versioned.
- A request has exactly one `PLAN.md`. It owns the request's lifecycle state via frontmatter.
- Requests **never** carry a PRD.md. Product-level intent already lives at the root — duplicating it per-request creates drift.
- If a request needs to extend or revise product intent, edit root `PRD.md` (snapshot first) — don't fork it into a request.

## PLAN.md (required)

Defines intent for one request — what we're building in this slice and why.

- goal
- why
- scope / out of scope
- approach
- success criteria (product/goal level)
- rollout

Frontmatter owns `status` for the request lifecycle.

---

## TASKS.md (required)

Defines:
- executable implementation checklist
- grouped work items

Skill monitors task completion as a signal for lifecycle transition proposals.

---

## SESSION.md

Created automatically when request moves to `active`.

Defines:
- current execution state
- blockers
- next actions

Committed to git — part of the team's operational record.

---

## RISKS.md (on demand)

Skill proposes creation when request touches auth, billing, migrations, or anything flagged sensitive in STACK.md.

Defines:
- edge cases
- architectural risks
- mitigation plans

Agents rarely reason about failure modes unless explicitly prompted. This file improves implementation quality significantly.

---

## VERIFY.md (on demand)

Skill proposes creation for requests with user-visible behavior changes or high-stakes implementation.

**Purpose: execution proof** — how you confirm the implementation actually worked.

Distinct from PLAN.md:
- PLAN answers "did we build the right thing?"
- VERIFY answers "did we build it correctly and safely?"

Defines:
- step-by-step manual QA checklist
- specific edge cases to verify
- regression checklist
- rollback validation steps

---

# Skills Layer

`.spectacular/skills/` contains **project-specific** reusable skills.

```txt
skills/
├── review/
├── migration/
└── release/
```

## Rules

- project-specific skills live here, authored per repo
- symlinked into `.claude/skills/` only on demand, only if the skill is runnable
- `.spectacular/skills/` never contains the Spectacular skill itself

## Spectacular skill location

- **Global install**: `~/.claude/skills/spectacular/`
- **Project-local install**: `.claude/skills/spectacular/` (created by `spectacular init`)

## Skill architecture

The Spectacular skill is intentionally lean — a thin orchestration layer that routes to specific reference docs:

```txt
~/.claude/skills/spectacular/
├── SKILL.md                    # lean orchestrator — triggers, routing, state awareness
└── references/
    ├── init-workflow.md        # CLI init + first-time onboarding
    ├── new-request.md          # scaffolding requests, slug rules, templates
    ├── active-request.md       # continuing work, session state, task tracking
    ├── lifecycle.md            # state transitions, signal detection, proposals
    ├── memory.md               # remember command, archive reflection
    ├── current-sync.md         # proposing current/ updates on archive
    ├── scaffold-reference.md   # canonical file templates with frontmatter stubs
    └── onboarding.md           # first invocation on an existing project
```

---

# Workflows Layer

Deferred to v2.

Intended purpose: document project-specific procedural sequences (release cycles, hotfix flows, migration procedures). Each project handles these differently; the value is clear but the design is not yet finalized.

---

# Memory Layer

`.spectacular/memory/` stores long-term operational learning.

```txt
memory/
├── failures.md
├── lessons.md
├── architecture-traps.md
└── recurring-bugs.md
```

## Rules

- **git-committed, team-visible** — survives agent changes, tool changes, team changes
- completely separate from `.claude/` personal memory (not Spectacular's concern)
- written by the skill, never by agents autonomously

## Write triggers

- **On archive**: skill reviews the completed request for notable blockers, risks hit, or lessons. Proposes memory entries; human confirms.
- **On demand**: `spectacular remember this` captures insights mid-session. Skill writes immediately on confirmation.

**Note:** the skill must avoid phrasing that triggers Claude Code's own auto-memory system to prevent double-capture.

---

# Archive Layer

Archive preserves completed requests and historical context.

```txt
archive/
├── add-team-billing/       # completed request, same slug
└── ideas/                  # promoted idea files
```

## Rules

- keep original slug/id
- never modify archived content
- skill does not read archive during normal operation (write-only from skill perspective)

---

# Versioning

Canonical documents are never overwritten in place.

## Rules

- skill always proposes a snapshot before editing any canonical document
- snapshot naming: `PRD@v1.0.md`, `STACK@v1.2.md` — `@version` suffix
- version tracked in frontmatter: `version: 1.0`
- the unversioned filename (`PRD.md`) always points to the current/latest version
- snapshots live alongside the current file (or in a subdirectory if preferred)
- applies to: root layer files, `current/` capability specs, `config.yaml`
- this is default behavior — not opt-in

---

# Lifecycle Model

```txt
idea (optional scratchpad)
  ↓
planned   → request scaffolded, PLAN.md + TASKS.md created
  ↓
active    → SESSION.md created, implementation underway
  ↓
review    → implementation complete, VERIFY.md checklist being run
  ↓
verified  → all checks passed
  ↓
archived  → moved to archive/, current/ updated, memory proposed
```

## State storage

- `status:` frontmatter in `PLAN.md` = request lifecycle state
- `status:` frontmatter in `current/<capability>.md` = capability state (`stable | draft | deprecated`)
- `status:` frontmatter in `requests/<slug>/specs/` = individual spec development state

## Transition rules

- skill detects signals and **proposes** transitions (e.g. all TASKS.md items checked → proposes moving to review)
- user can also force transitions explicitly
- skill is proactive on maintenance — surfaces stale state, blocked requests, missing updates

---

# Skill Interaction Model

## Invocation

`/spectacular` with no arguments triggers a **proactive conversational briefing**:

1. Skill reads frontmatter from root files, `current/`, and `requests/`
2. Surfaces a minimal state table (active requests, draft capabilities, what needs attention)
3. Identifies single highest-priority next action
4. Asks what the user wants to do

## Output format

Conversational briefing with an embedded minimal table. Example:

```
Project state as of 2026-05-11:

| Layer    | Items                                      |
|----------|--------------------------------------------|
| Current  | auth (stable), billing (draft)             |
| Requests | add-team-billing (review), dark-mode (planned) |
| Memory   | 2 lessons, 1 trap                          |

add-team-billing is in review — all tasks checked. Ready to run VERIFY.md or archive?
```

## Request creation

- **Autopilot** when conversational context is clear: skill derives slug, drafts PLAN.md, shows user, confirms, writes
- **Interactive** when context is thin: skill asks targeted questions, then scaffolds
- User can always override slug before creation

## Explicit commands

All invocable by user directly or triggered by skill proactively:

- `spectacular new <description>` — scaffold new request
- `spectacular status` — project briefing
- `spectacular archive <slug>` — archive completed request, propose current/ sync + memory
- `spectacular remember this` — write to memory immediately
- `spectacular promote <idea-file>` — promote idea to request
- `spectacular snapshot <file>` — manually snapshot a canonical document

---

# Init Flow (CLI)

When `spectacular init` is run on a new project:

1. Scaffold `.spectacular/` directory structure
2. Prompt for project name and summary, write `config.yaml`
3. Create stub root files with frontmatter templates: `PRD.md`, `STACK.md`, `DECISIONS.md`, `AGENTS.md`
4. Install skill into `.claude/skills/spectacular/` (project-local) or prompt for global install
5. Add `.spectacular.local/` to `.gitignore`

`.spectacular/` is always fully committed to git. `.spectacular.local/` is always gitignored.

---

# Agent Design Principles

Agents should:
- load minimal viable context
- prefer local capability reasoning
- summarize before handoff
- avoid giant prompts
- preserve continuity through memory
- route to specific reference docs rather than loading everything

---

# Retrieval Principles

## Context layers

```txt
global context     → config.yaml, AGENTS.md
product context    → PRD.md, STACK.md, DECISIONS.md
capability context → current/<relevant-capability>/
change context     → requests/<slug>/PLAN.md, TASKS.md
task context       → requests/<slug>/SESSION.md, specs/
```

## Retrieval order by mode

### Planning
- PRD.md
- PRINCIPLES.md
- DECISIONS.md

### Implementation
- STACK.md
- PLAN.md
- TASKS.md
- local capability specs

### Review
- VERIFY.md
- capability specs
- RISKS.md

---

# Anti-Entropy Rules

- prefer spec files shorter than 500 lines
- avoid duplicated truth
- archive stale execution context
- split capabilities aggressively
- prefer explicit lifecycle states
- preserve operational memory
- never overwrite canonical documents — snapshot first

## The most important insight

You are NOT designing a documentation system.

You are designing a temporal operational model.

The names should express time and state.
`current` and `requests` do that beautifully.

Current truth → `load current/auth/`
Active work → `load requests/add-team-billing/`

---

# Future Vision (v2+)

- **Workflows layer** — project-specific procedural sequences (release, hotfix, migration)
- **Multi-agent orchestration** — subagent handoff conventions, parallel execution
- **Hook-driven automation** — auto-update SESSION.md on commit, auto-archive on merge
- **Context orchestration framework** — smart retrieval across the full spectacular structure
- **Repository operating system** — from solo builders to autonomous multi-agent engineering systems
- **Workspaces** — team-scoped workspace layers (see below)
- **Nested workspaces** — `.spectacular/` inside subdirectories for monorepo support (see below)

---

# v2 — Workspaces

Workspaces allow multiple teams or roles to maintain separate operational contexts within the same project.

## Naming convention

```
.spectacular/              ← default workspace (always present)
.spectacular.local/        ← personal overrides, always gitignored
.spectacular.<workspace>/  ← named team workspaces
```

Examples:
```
.spectacular.designteam/
.spectacular.devops/
.spectacular.builder/
```

## Rules

- Default workspace (`.spectacular/`) is always the base
- Named workspaces are fully committed and team-visible (same rules as default)
- `.spectacular.local/` is always gitignored regardless of workspace
- Skill reads the active workspace based on invocation context or explicit flag
- Named workspaces follow the same internal structure as the default workspace
- Workspaces do not inherit from each other — each is independent

## Invocation (proposed)

```
/spectacular                   ← operates on default workspace
/spectacular --workspace devops ← operates on .spectacular.devops/
spectacular status --workspace designteam
```

Design is not finalized. Workspace switching UX TBD.

---

# v2 — Nested Workspaces

Workspaces can live inside subdirectories of a monorepo, scoped to a specific app or package.

## Example

```
apps/builder/.spectacular/
apps/api/.spectacular/
packages/ui/.spectacular/
```

## Rules

- Nested workspaces are independent — they do not inherit from the repo-root `.spectacular/`
- The skill detects the nearest `.spectacular/` walking up from the current working directory
- Each nested workspace has its own `config.yaml`, root files, requests, memory, etc.
- Cross-workspace coordination is out of scope for v2 — each workspace operates independently

## Use case

Useful in monorepos where separate teams own separate apps and want independent operational context without cross-contamination.
