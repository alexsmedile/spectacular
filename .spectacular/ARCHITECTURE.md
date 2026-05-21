---
version: 1.0
updated: 2026-05-21
summary: "The .spectacular/ directory — layers, file roles, frontmatter conventions, lifecycle, versioning"
related:
  - PRD.md
  - PRINCIPLES.md
  - AGENTS.md
---

# Spectacular — Architecture

This document defines the **structure of `.spectacular/`** — what each folder and file is for, how they relate, and the conventions every file must follow. For the *intent* behind these choices, see [PRD.md](PRD.md); for the *principles*, see [PRINCIPLES.md](PRINCIPLES.md).

This is distinct from `STACK.md` — STACK describes the **host project's** technology choices (Next.js, Postgres, etc.); ARCHITECTURE describes **Spectacular's own** layout.

---

# Layout

```txt
.spectacular/
│
├── PRD.md              # product intent (this project)
├── PRINCIPLES.md       # operating principles
├── ARCHITECTURE.md     # this file
├── ROADMAP.md          # time-ordered "what's next"
├── AGENTS.md           # how to operate inside .spectacular/
├── STACK.md            # host project's tech choices
├── DECISIONS.md        # ADR-style decision log
├── config.yaml         # machine-readable project config
│
├── ideas/              # exploratory scratchpad — not acted on by skill
├── current/            # canonical system truth — capability specs
├── requests/           # active and planned work
├── skills/             # project-specific reusable skills
├── memory/             # long-term operational learning
└── archive/            # completed requests, historical snapshots
```

`.spectacular.local/` — personal overrides, always gitignored, never committed. `.spectacular/` itself is fully committed to git.

---

# Root layer

The root layer is stable project grounding. These files change infrequently, stay concise, avoid implementation details, and are **never overwritten in place** — snapshot before editing (see § Versioning).

| File | Purpose |
|---|---|
| `PRD.md` | Product intent — what Spectacular (or the host project) is, for whom, why |
| `PRINCIPLES.md` | Operating principles + runtime enforcement hooks |
| `ARCHITECTURE.md` | This file — `.spectacular/` structure and conventions |
| `ROADMAP.md` | Versioned future work |
| `AGENTS.md` | Onboarding doc for any agent landing in `.spectacular/` |
| `STACK.md` | **Host project** technology and architecture choices |
| `DECISIONS.md` | ADR-style log — one decision per entry, immutable |
| `config.yaml` | Machine-readable project config |

## STACK.md vs ARCHITECTURE.md

These are **two different docs at two different scopes**:

- **ARCHITECTURE.md** describes *Spectacular's own structure* — the `.spectacular/` layout that applies to every project that adopts it.
- **STACK.md** describes *the host project's tech choices* — Next.js, Postgres, deployment targets. Replaced per project.

For the Spectacular repo itself, `STACK.md` happens to describe Spectacular's tooling (Bash, markdown, Claude skill format). For any consumer project, it describes their own stack.

## AGENTS.md

Spectacular-specific. Distinct from the repo-root `CLAUDE.md` / `AGENTS.md`. Tells any agent landing inside `.spectacular/`:
- which context to load for which task type
- which skills are available
- handoff conventions
- what *not* to touch

Humans write it; the skill proposes updates when new skills or capabilities are added.

## DECISIONS.md

ADR-style log. One decision per entry. Immutable once written. Each entry contains:

```md
## YYYY-MM-DD — <short title>

Decision:
<what we decided>

Why:
<reasoning>

Tradeoffs:
<what we gave up>
```

---

# Configuration

`config.yaml` is the machine-readable project configuration. The skill reads it on every invocation.

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

# Frontmatter conventions

Frontmatter is the skill's primary signal for reading project state. Every canonical document includes frontmatter.

## Root files (PRD, PRINCIPLES, ARCHITECTURE, ROADMAP, AGENTS, STACK, DECISIONS)

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "One-sentence description of this file's purpose"
related:
  - <sibling-file>.md
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
- Other fields are optional; skill warns but does not block
- `PLAN.md` frontmatter is the **single source of lifecycle state** for a request
- Capability specs in `current/` track their own state independently

---

# Ideas layer

A **thinking scratchpad**, not a workflow stage. Nothing in `ideas/` is acted on automatically by the skill.

Use it for: raw thoughts, market observations, UX experiments, discarded approaches, future concepts, unresolved brainstorming.

```txt
ideas/
├── multiplayer-editor.md
├── ai-memory-system.md
└── growth-loops.md
```

**Rules:**
- low commitment
- speculative
- non-canonical
- skill **proposes** saving unresolved decisions here when conversations have open branches

**Promotion to request:** Ideas are not a required gate. A request can be created directly. When an idea is deliberately promoted, the skill scaffolds the request from the idea content and moves the idea file to `archive/ideas/`.

---

# Current layer

The current layer represents canonical system truth — what the product *actually does* right now.

```txt
current/
├── auth/
│   ├── login.md
│   ├── sessions.md
│   └── permissions.md
├── billing/
└── editor/
```

**Purpose:** defines current behavior, active capabilities, security requirements, performance expectations, user-visible behavior.

**Rules:**
- authoritative
- current only — no past state, no future plans
- behavior-oriented, not implementation-oriented
- modular — one capability per file or folder
- **never overwritten in place** — skill snapshots before proposing edits
- skill proposes `current/` updates when a request is archived; humans confirm

**Capability spec structure** — each spec contains:
- purpose
- requirements
- scenarios
- security considerations
- performance expectations

---

# Requests layer

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

**Slug rules:**
- skill derives slug from conversation context, shows user before creating
- user can override at any time
- slugs are kebab-case by default (configurable in `config.yaml`)
- slugs are unique — if slug exists, skill proposes `-2` suffix or asks user

**Rules:**
- temporary
- operational
- archived on completion (never deleted)
- `PLAN.md` frontmatter owns lifecycle state

---

# Request files

## PRD vs PLAN — scope distinction

These are **two different artifacts at two different layers**, not the same artifact at two scopes.

| Artifact | Location | Scope | Answers |
|---|---|---|---|
| `PRD.md` | `.spectacular/` root only | **Product** (whole project) | Why does this product exist? |
| `PLAN.md` | `requests/<slug>/` only | **Request** (one slice of work) | What are we building in this slice and why? |

**Rules:**
- A project has exactly one `PRD.md` (at the root). Long-lived, snapshot-versioned.
- A request has exactly one `PLAN.md`. Owns lifecycle state via frontmatter.
- Requests **never** carry a PRD.md. Product-level intent already lives at the root.
- If a request needs to extend or revise product intent, edit root `PRD.md` (snapshot first) — don't fork it into a request.

## PLAN.md (required)

Defines intent + plan for one request. 7-slot shape:

- **Goal** — one sentence; compressed intent from PRD
- **Constraints** — what's fixed before starting
- **Milestones** — ordered, demoable checkpoints (not tasks — outcomes)
- **Tasks** — pointer to `TASKS.md`
- **Dependencies** — other requests, skills, blocking decisions
- **Validation** — how each milestone is verified
- **Deliverables** — artifacts that ship out of this request

Frontmatter owns `status:` for the request lifecycle.

## TASKS.md (required)

Executable implementation checklist, grouped by milestone. The skill monitors task completion as a signal for lifecycle transition proposals.

**Frontmatter conventions:**
- `depends_on:` — surface task dependencies
- `validates:` — link task groups to milestones (closes principle 7's validation loop)

## SESSION.md

Created automatically when a request moves to `active`. Captures current execution state, blockers, next actions. Committed to git — part of the team's operational record.

## RISKS.md (on demand)

Skill proposes creation when a request touches auth, billing, migrations, or anything flagged sensitive in `STACK.md`. Defines edge cases, architectural risks, mitigation plans.

Agents rarely reason about failure modes unless explicitly prompted — this file improves implementation quality significantly.

## VERIFY.md (on demand)

Skill proposes creation for requests with user-visible behavior changes or high-stakes implementation.

**Purpose: execution proof** — how you confirm the implementation actually worked.

Distinct from PLAN.md:
- PLAN answers "did we build the right thing?"
- VERIFY answers "did we build it correctly and safely?"

Contains: step-by-step manual QA checklist, specific edge cases, regression checklist, rollback validation.

---

# Skills layer

`.spectacular/skills/` contains **project-specific** reusable skills.

```txt
skills/
├── review/
├── migration/
└── release/
```

**Rules:**
- project-specific skills live here, authored per repo
- symlinked into `.claude/skills/` only on demand, only if runnable
- `.spectacular/skills/` never contains the Spectacular skill itself

**Spectacular skill location:**
- Global install: `~/.claude/skills/spectacular/`
- Project-local install: `.claude/skills/spectacular/` (created by `spectacular init`)

**Skill architecture** — the Spectacular skill is intentionally lean:

```txt
~/.claude/skills/spectacular/
├── SKILL.md                    # lean orchestrator — triggers, routing, state awareness
└── references/
    ├── init-workflow.md
    ├── new-request.md
    ├── active-request.md
    ├── lifecycle.md
    ├── memory.md
    ├── current-sync.md
    ├── prd-grill.md
    ├── prd-refine.md
    ├── prd-review.md
    ├── scaffold-reference.md
    └── onboarding.md
```

---

# Memory layer

`.spectacular/memory/` stores long-term operational learning.

```txt
memory/
├── failures.md
├── lessons.md
├── architecture-traps.md
└── recurring-bugs.md
```

**Rules:**
- **git-committed, team-visible** — survives agent changes, tool changes, team changes
- completely separate from `.claude/` personal memory
- written by the skill on confirmation, never by agents autonomously

**Write triggers:**
- **On archive:** skill reviews the completed request for notable blockers, risks hit, or lessons. Proposes memory entries; human confirms.
- **On demand:** `spectacular remember this` captures insights mid-session.

Skill must avoid phrasing that triggers Claude Code's own auto-memory to prevent double-capture.

---

# Archive layer

`.spectacular/archive/` preserves completed requests and historical context.

```txt
archive/
├── add-team-billing/       # completed request, same slug
└── ideas/                  # promoted idea files
```

**Rules:**
- keep original slug/id
- never modify archived content
- skill does not read `archive/` during normal operation (write-only from skill's perspective)

---

# Lifecycle

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

**State storage:**
- `status:` in `PLAN.md` frontmatter = request lifecycle state
- `status:` in `current/<capability>.md` = capability state (`stable | draft | deprecated`)
- `status:` in `requests/<slug>/specs/` = individual spec development state

**Transition rules:**
- skill detects signals and **proposes** transitions (e.g. all TASKS items checked → propose move to `review`)
- user can force transitions explicitly
- skill is proactive on maintenance — surfaces stale state, blocked requests, missing updates

---

# Versioning

Canonical documents are **never overwritten in place**.

**Rules:**
- skill always proposes a snapshot before editing any canonical document
- snapshot naming: `PRD@v1.0.md`, `STACK@v1.2.md` — `@version` suffix
- version tracked in frontmatter: `version: 1.0`
- the unversioned filename (`PRD.md`) always points to the current version
- snapshots live alongside the current file
- applies to: root layer files, `current/` capability specs, `config.yaml`
- this is **default behavior** — not opt-in

---

# Init flow

`spectacular init` is a one-time CLI bootstrap. Detailed implementation lives in [`requests/cli-bootstrap/PLAN.md`](requests/cli-bootstrap/PLAN.md).

Summary:
1. Scaffold `.spectacular/` directory structure
2. Prompt for project name and summary; write `config.yaml`
3. Create stub root files with frontmatter templates
4. Install the skill into `.claude/skills/spectacular/` (project-local) or `~/.claude/skills/spectacular/` (global)
5. Add `.spectacular.local/` to `.gitignore`

`.spectacular/` is always fully committed. `.spectacular.local/` is always gitignored.

---

# Related

- [PRD.md](PRD.md) — why Spectacular exists
- [PRINCIPLES.md](PRINCIPLES.md) — the principles this architecture implements
- [ROADMAP.md](ROADMAP.md) — v2+ structural additions (workspaces, nested workspaces, workflows)
- [AGENTS.md](AGENTS.md) — how to operate inside this structure
