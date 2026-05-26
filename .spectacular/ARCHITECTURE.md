---
version: 1.1
updated: 2026-05-22
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
├── feedback/           # prototyping-mode feedback entries (v1.6.0+; system-level)
├── snapshots/          # versioned snapshots of canonical docs (v1.6.0+)
│   ├── PRD/            # one folder per canonical doc, uppercase preserved
│   │   └── @v1.2.md
│   └── ROADMAP/
│       └── @v4.md
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

**CLI verbs (v1.7.0+):** `spectacular idea new <slug>`, `spectacular idea list [--status <s>]`, `spectacular idea promote <slug>`. Status enum: `parked | exploring | promoted`. Full spec in [[idea-rules]]; doctor area: [[doctor-areas]] § ideas.

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
├── feedback/           # request-scoped feedback-loop entries (v1.6.0+; see references/feedback-loop.md)
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

**On-demand only.** The skill proposes creation when the request hits the **2-of-6 rule** (see [[verification]] for full text):

1. User-visible change
2. High reversibility cost
3. Multi-surface verification
4. Risk surface non-trivial (auth/billing/security/data)
5. External contract change
6. Rollback plan exists

When fewer than 2 axes hit, **fold verification into PLAN § Validation or TASKS § Verification** — no separate file. Most internal/spec/refactor/doc requests don't need VERIFY.md.

**Purpose: execution proof** — how you confirm the implementation actually worked.

Distinct from PLAN.md and TASKS.md:
- PLAN § Validation answers "what does each milestone need to satisfy?"
- TASKS § Verification answers "what step-by-step checks confirm done?"
- VERIFY.md answers "did we build it correctly and safely, with risk-aware coverage?"

Contains (when scaffolded): step-by-step manual QA checklist, specific edge cases, regression checklist, rollback validation.

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
- `status:` in `specs/<capability>/SPEC.md` = capability state (`stable | draft | deprecated`); top-level `SPEC.md` carries no per-capability status — it's the index
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
- snapshot location (v1.6.0+): `snapshots/<DOC>/@v<N>.md` — folder per canonical doc, uppercase preserved
- sub-doc snapshots mirror their path: `specs/cli/SPEC.md` → `snapshots/specs/cli/SPEC/@v1.0.md`
- version tracked in frontmatter: `version: 1.0`
- the unversioned filename at root (`PRD.md`) always points to the current version
- applies to: root layer files, `SPEC.md`, `specs/<capability>/SPEC.md` capability specs, `config.yaml`
- this is **default behavior** — not opt-in
- legacy snapshots at root (`PRD@v1.0.md`) continue to be read; `spectacular doctor snapshots` warns until migrated via `--fix` (warn until v1.7, then info)

---

# Init flow

`spectacular init` is a one-time CLI bootstrap. Detailed implementation lives in [`requests/cli-bootstrap/PLAN.md`](requests/cli-bootstrap/PLAN.md) (v0.2.x) and [`requests/smart-init/PLAN.md`](requests/smart-init/PLAN.md) (v0.3.0+).

## v0.3.0 — smart init

As of v0.3.0, init scaffolds **only what the project needs**, not all root docs.

**Always-set** (6 files + 2 dirs, scaffolded unconditionally — v0.5.0+):
- `.spectacular/PRD.md` — anchor doc; every other doc references it
- `.spectacular/SPEC.md` — system spec index (what's built right now, present tense)
- `.spectacular/config.yaml` — project name, kit identity, naming rules
- `.spectacular/<agents-file>` — onboarding doc (defaults to `AGENTS.md`)
- `.spectacular/requests/` — request folders
- `.spectacular/specs/` — per-capability specs (optional content; only when a capability outgrows a one-liner in SPEC.md)

> v0.4.0 and earlier scaffolded `.spectacular/current/` instead of `SPEC.md` + `specs/`. The legacy folder is auto-migrated via `spectacular doctor specs --fix`.

**Kit-driven additions** (see [[kits-contract]]):
- The user picks a kit (`blank`, `coding`, `content`, `product`, `research`)
- Each kit declares `triggers-docs.always` (scaffolded automatically) and `triggers-docs.suggested` (interactive prompt y/n)
- Non-interactive default: `blank` kit, no extras

**Explicit additions** via `--with <doc1,doc2,...>` flag — additive over kit defaults.

**Suppression** via `--minimal` — scaffolds always-set only, ignoring kit's always-docs. Kit identity is still recorded in PRD frontmatter.

## Sequence

1. Parse flags + validate (`--kit` known, `--with` doc IDs in registry)
2. If `-i`: run interactive prompts (name, summary, agents-file, scope, kit menu, per-suggested-doc y/n)
3. Resolve doc-set: `always-set ∪ (kit always-docs unless --minimal) ∪ --with entries`
4. Scaffold directories
5. Per-doc dispatch via `write_if_missing` (pre-flight rules: skip if exists, fill if empty, diagnose if malformed)
6. Update `.gitignore` (append `.spectacular.local/` if absent)
7. Install skill into `.agents/skills/spectacular/` (or `~/.agents/skills/spectacular/` with `--global`)
8. Symlink `.claude/skills/spectacular/` → install location

## Idempotency + non-destructive

Re-running init on an initialized workspace is always safe — no file is ever overwritten. Adding a kit later (`spectacular init --kit coding` on an existing project) only scaffolds the kit's missing always-docs; existing files are left alone.

`.spectacular/` is always fully committed. `.spectacular.local/` is always gitignored.

---

# Related

- [PRD.md](PRD.md) — why Spectacular exists
- [PRINCIPLES.md](PRINCIPLES.md) — the principles this architecture implements
- [ROADMAP.md](ROADMAP.md) — v2+ structural additions (workspaces, nested workspaces, workflows)
- [AGENTS.md](AGENTS.md) — how to operate inside this structure
