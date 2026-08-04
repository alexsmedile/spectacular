---
title: Scaffold Reference
description: Complete .spectacular/ directory spec — every file, frontmatter schema, creation rules, versioning.
section: ""
type: reference
status: stable
since: 0.1.0
updated: 2026-08-04
---

# Scaffold Reference

The `.spectacular/` directory is the operational workspace for a project. This document specifies every file and directory — what it is, what goes in it, when it is created, and what rules govern it.

---

## File-type catalog

Every file type Spectacular uses, grouped by scope. This table is the quick reference; the detailed per-file sections follow below.

### Canonical top-level docs — `.spectacular/`

| File | Role | Usage |
|---|---|---|
| `PRD.md` | Product intent | What / why / for whom. Project-wide — **never** per-request. |
| `specs/index.md` | System spec (index) | Cheap, always-on index of what's built now; links to `specs/<cap>.md`. |
| `PRINCIPLES.md` | Operating principles | The "why / how we work" beliefs + enforcement hooks. |
| `POLICY.md` | Practice layer | Prefilled policies under work-phase hooks; gates lifecycle transitions. |
| `ARCHITECTURE.md` | Workspace structure | Layout, frontmatter conventions, lifecycle, versioning. |
| `roadmaps/index.md` | What's next | Time-ordered direction (v1 / v2 / v3+) + Icebox. |
| `STACK.md` | Tech choices | Host project's frontend / backend / infra + engineering rules. |
| `AGENTS.md` | Agent onboarding | How to operate in `.spectacular/`; authoritative context-loading table. |
| `decisions/index.md` | Decision log (ADR) | Index of evidence-backed `DEC-NNN` records. |
| `PERSONAS.md` | Audience profiles | Opt-in — only with `product` / `content` kit or `--with personas`. |
| `config.yaml` | Machine config | Name, naming rules, agent context, `workspace_schema` + provenance. |

### Per-request files — `requests/<slug>/`

| File | Role | Usage |
|---|---|---|
| `PLAN.md` | Request decomposition | 7-slot plan; **owns the lifecycle `status:` field**. Required. |
| `TASKS.md` | Implementation checklist | Executable milestone blocks (`## M<N>`). Required. |
| `SESSION.md` | Request working-state | Current-state / blockers / next-actions. Created on `active`. **Singular.** |
| `VERIFY.md` | Verification checklist | Typed checks for user-visible / high-stakes changes. On demand. |
| `VERIFY-LOG.md` | Verification evidence | Append-only; one entry per `spectacular verify` walk. |
| `UNDERSTANDING.md` | Understand-before-change | Optional; alternative to PLAN's Understanding slot; gates `planned → active`. |
| `RISKS.md` | Risk register | On demand for auth / billing / migration / flagged-sensitive requests. |
| `specs/<cap>.md` | Per-capability truth | System-truth spec for one capability (distinct from top-level `specs/index.md`). |

### Soft-DB collections — folder of entries + index

Each is a folder of individually-addressable `.md` entries (frontmatter, git-committed, **appended never overwritten**). The index is regenerated from the folder.

| Collection | Index | Role | Write verb |
|---|---|---|---|
| `visions/` | — | Pre-request direction workspaces (`.spectacular/visions/<slug>/`) | `spectacular imagine` |
| `memories/` | `memories/index.md` | Durable standing facts / "always do X" | `spectacular remember` |
| `decisions/` | `decisions/index.md` | ADR — why A over B | `spectacular decide` |
| `questions/` | — | Active human-owned blockers and ambiguities (`QUE-NNN`); resolved history moves to `archive/questions/` | `spectacular question new\|resolve` |
| `research/` | — | Sourced discovery evidence (`RES-NNN`) | `spectacular research new\|resolve` |
| `spikes/` | — | Authorized feasibility evidence (`SPK-NNN`) | `spectacular spike new\|resolve` |
| `sessions/` | `sessions/index.md` | Work-session time-log | `spectacular session start\|end` |
| `ideas/` | — | Parked inspiration (`IDEA-NNN`) | `spectacular idea new` → `promote` |
| `feedbacks/` | — | Post-ship prototyping signal | `spectacular feedback-loop new` |
| `audits/` | `A<N>.md` | Bug diagnosis before a fix is planned | `spectacular audit new\|resolve` |
| `fixes/` | `F<N>.md` | Verified, signed, reusable fix corpus | `spectacular fix new\|list` |

### Optional advanced engineering collections

Heavy projects can reserve additional evidence stores without adding them to every workspace:

```bash
spectacular init --with findings,fixes,bugs,security,benchmarks
```

| Folder | Reserved ID | Intended use | Availability |
|---|---|---|---|
| `findings/` | `FND-NNN` | Research, audit, or interview takeaway | Path/ID reserved; workflow deferred |
| `fixes/` | `FIX-NNN` | Remediation/hotfix identity | Existing `F<N>` ledger stays active; migration deferred |
| `bugs/` | `BUG-NNN` | Known defect or regression | Path/ID reserved; workflow deferred |
| `security/` | `SEC-NNN` | Vulnerability or threat finding | Path/ID reserved; workflow deferred |
| `benchmarks/` | `BMK-NNN` | Performance/load/profiling result | Path/ID reserved; workflow deferred |

`security/` is singular because it names the engineering domain; `securities/` means financial instruments. Use `fnd1` for findings, `fix1` for fixes, and `bug1` for bugs. Ambiguous `f1` is refused; `b1` remains the roadmap build alias.

`audit → requests → fixes` form the self-learning bug loop (see the bug-workflow skill reference).

### ⚠ Naming conventions & traps

The rule: **a plural folder name is a category directory** — it holds either an `index.md` plus sequential entry files (soft-DB collections) or per-item sub-directories (execution trees).

1. **Two `SESSION`s, opposite scope.** Per-request `SESSION.md` (singular — one request's working state, created on `active`) is **unrelated** to the top-level `sessions/` category folder and its `sessions/index.md` (the work-session time-log). Same word, different system — the biggest confusion trap.
2. **Consolidated specs.** Capability specs are flat files inside `specs/` (e.g. `specs/cli.md`, `specs/doc-engine.md`), indexed by `specs/index.md`. There are no nested spec folders.
3. **Plural folders.** Collection and execution folders are plural (`visions/`, `memories/`, `roadmaps/`, `decisions/`, `questions/`, `research/`, `spikes/`, `sessions/`, `audits/`, `fixes/`, `feedbacks/`, `ideas/`, `requests/`, `debugs/`, `findings/`, `bugs/`, `benchmarks/`) except `security/`, which is the singular domain name.
4. **Collections vs. execution trees.** Markdown collections hold flat, individually addressable entries; some also maintain an `index.md`. Execution trees (`requests/`, `debugs/`, `visions/`) hold per-item sub-directories with state files or run logs.
5. **Canonical identity.** Wayfinding filenames begin with `DEC-NNN`, `QUE-NNN`, `IDEA-NNN`, `RES-NNN`, `SPK-NNN`, or `SPC-NNN`. `PRT-NNN` and `TSK-NNN` remain reserved. Prototype outputs belong to a pre-spec Vision workspace (as a `prototype` fragment under `visions/<slug>/fragments/`), a post-build feedback entry (`feedbacks/`), or a spike (`spikes/`); generic artifacts, tracer bullets, and technical debt do not receive `ART`, `TRC`, or `DEB` IDs. Conversational aliases never replace canonical IDs in saved frontmatter.

---

## Directory structure

```
.spectacular/
├── PRD.md              # product intent — what & why & for whom (required)
├── PRINCIPLES.md       # operating principles + enforcement hooks (required)
├── ARCHITECTURE.md     # workspace structure, frontmatter, lifecycle, versioning (required)
├── roadmaps/index.md   # time-ordered "what's next" (required)
├── AGENTS.md           # onboarding doc for any agent in .spectacular/ (required)
├── STACK.md            # host project's technology choices (required)
├── decisions/index.md  # ADR-style decision log (required)
├── config.yaml         # machine-readable project config (required)
│
├── specs/index.md      # system spec — index of what's built now (v0.5.0+)
├── POLICY.md           # practice layer — work-phase hooks that gate transitions
│
├── specs/              # canonical system truth — one file per capability
│   │                   #   (renamed from current/ in v0.5.0)
│   ├── auth.md
│   ├── payments.md
│   └── subscriptions.md
│
├── memories/  decisions/  questions/  research/  spikes/  visions/
├── sessions/  ideas/  feedbacks/  audits/  fixes/
│                        # durable Markdown collections, created on demand
│
├── visions/            # pre-request direction workspaces — created on demand
│   └── <slug>/
│       ├── VISION.md   # required — 10-slot direction spine & approval gate
│       ├── fragments/  # flat proposal files (strategy/story/flow/ui/arch/prototype)
│       └── evidence/   # linked presentation artifacts or notes
│
├── requests/           # active and planned work — one folder per request
│   └── add-team-billing/
│       ├── PLAN.md     # required — 7-slot decomposition, owns lifecycle state
│       ├── TASKS.md    # required — executable checklist
│       ├── SESSION.md  # created on active
│       ├── RISKS.md    # on demand
│       ├── VERIFY.md   # on demand
│       └── specs/      # per-request capability specs
│
├── ideas/              # thinking scratchpad — not acted on automatically
├── memories/           # long-term operational learning (git-committed)
├── skills/             # project-specific reusable skills
└── archive/            # completed requests and promoted ideas
    └── add-team-billing/
```

**Requests never carry their own `PRD.md`** — product intent is project-wide and lives at `.spectacular/PRD.md`. Use `PLAN.md` for request-level intent.

`.spectacular.local/` — personal overrides at the repo root. Always gitignored, never committed.

---

## Creation rules

| Directory | When created |
|---|---|
| `specs/` | On `spectacular init` (was `current/` pre-v0.5.0) |
| `requests/` | On `spectacular init` |
| `visions/` | On first `spectacular imagine <slug>` |
| `ideas/` | On first idea file |
| `questions/` | On first `spectacular question new` |
| `research/` | On first `spectacular research new` |
| `spikes/` | On first `spectacular spike new` |
| `memories/` | On first `spectacular remember this` |
| `skills/` | On first project skill |
| `archive/` | On first archived request |
| `archive/ideas/` | On first promoted idea |

---

## Root layer files

These files form the stable grounding for the workspace. They should change infrequently and remain concise. **Never overwrite in place — snapshot before editing** (see [Versioning](#versioning)).

Spectacular splits root grounding across **seven focused docs** so agents can load only what a task needs (progressive disclosure). The PRD answers *what & why*; PRINCIPLES answers *how we operate*; ARCHITECTURE answers *what's where*; ROADMAP answers *when*; STACK answers *built with*; DECISIONS answers *why we chose what*; AGENTS answers *how to work in here*.

### `PRD.md`

Product/business intent. Why the project exists, what it is trying to achieve, who it is for. Uses the 8-slot canonical PRD shape (Vision / Problem / Target users / Deliverable / Goals & success criteria / Non-goals / Constraints / First milestone).

```yaml
---
version: 1.0
updated: 2026-05-21
summary: "One-sentence description of this file's purpose"
related:
  - PRINCIPLES.md
  - ARCHITECTURE.md
  - roadmaps/index.md
  - AGENTS.md
  - decisions/index.md
  - STACK.md
---
```

Sections (in order): Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone, Principles (summary), Related docs.

For interactive PRD building, use `spectacular prd` — the skill walks the 8-slot grill (Vision, Problem, Target users, Deliverable, Goals & success criteria, Non-goals, Constraints, First milestone) with kit-aware slot extensions.

---

### `PRINCIPLES.md`

Operating principles + runtime enforcement hooks. Each principle states a belief and pairs it with a concrete "How the skill enforces this:" sub-bullet — principles without enforcement are posters.

```yaml
---
version: 1.0
updated: 2026-05-21
summary: "Operating principles + runtime enforcement hooks"
related:
  - PRD.md
  - ARCHITECTURE.md
---
```

Default 8 principles (customize per project):
1. Context is the system
2. Separate intent from truth
3. Small files over giant documents
4. Humans and agents share the same workspace
5. Operational memory compounds
6. Progressive disclosure
7. Three layers — intent → execution → validation
8. Humans decide, agents propose

---

### `ARCHITECTURE.md`

The workspace structure itself — what each folder is for, frontmatter conventions, lifecycle states, versioning rules. Distinct from `STACK.md` (which describes the *host project's* tech choices); ARCHITECTURE describes the `.spectacular/` workspace.

```yaml
---
version: 1.0
updated: 2026-05-21
summary: ".spectacular/ structure, frontmatter, lifecycle, versioning"
related:
  - PRD.md
  - PRINCIPLES.md
  - AGENTS.md
---
```

Sections: Layout, Root layer, Frontmatter conventions, Ideas / Current / Requests / Skills / Memory / Archive layers, Request files, Lifecycle, Versioning, Configuration.

---

### `roadmaps/index.md`

Time-ordered "what's next". Coarse targets, not commitments. Detail for in-flight work lives in `requests/<slug>/`.

```yaml
---
version: 1.0
updated: 2026-05-21
summary: "Roadmap — v1 status, v2 features, v3+ direction"
related:
  - PRD.md
  - ARCHITECTURE.md
---
```

Sections: v1 (current), v2 features, v3+ direction.

---

### `STACK.md`

Technology choices and engineering rules. What the project is built on and the conventions that govern implementation decisions.

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "Technology stack and engineering rules"
---
```

Sections: Frontend, Backend, Infrastructure, Rules.

---

### `decisions/index.md`

Index of verified architectural and product decisions. Full ADRs live in `decisions/DEC-NNN-<slug>.md`; the index stays short enough for planning-time retrieval.

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "Architectural and product decisions log"
---
```

Index format:

```md
- **DEC-001** — Use Postgres RLS — Centralizes tenant isolation
```

Entry files use:

```yaml
---
id: DEC-001
type: decision
status: verified
origin: user
evidence: "User direction recorded by explicit command."
updated: 2026-08-01
---
```

```md
**Decision:** Use Postgres RLS instead of app-level permissions
**Context:** Tenant isolation must be enforced consistently
**Evidence:** Security review and prototype results
**Consequences:** Centralized policy; harder local debugging
```

---

### `AGENTS.md`

Onboarding doc for any agent or human landing inside `.spectacular/`. Defines what the folder is, how to operate in it, which context to load per task type, available skills, and don'ts. The `/spectacular` skill treats this file as the authoritative context-loading table.

```yaml
---
version: 1.0
updated: 2026-05-21
summary: "Onboarding doc for any agent or human landing inside .spectacular/"
related:
  - PRD.md
  - PRINCIPLES.md
  - ARCHITECTURE.md
---
```

Standard sections:

```md
## What this folder is
## How to operate
## Context loading by task

| Task type | Load |
|---|---|
| Planning / design | PRD.md, PRINCIPLES.md, decisions/index.md |
| Refining intent / PRD work | PRD.md, skill refs prd-grill.md / prd-refine.md / prd-review.md |
| Implementing a request | STACK.md, requests/<slug>/PLAN.md, TASKS.md, specs/index.md, relevant specs/<capability>.md |
| Reviewing / QA | requests/<slug>/VERIFY.md, relevant specs/<capability>.md, RISKS.md |
| Onboarding cold | PRD.md, ARCHITECTURE.md, this file |

## Available skills
- spectacular — workspace management

## Creating requests
Use `spectacular new <description>`. Never create requests/<slug>/PRD.md — anti-pattern.

## Don'ts
- Don't touch archive/
- Don't duplicate truth
- Don't overwrite canonical docs in place
- Don't write to memories/ autonomously
- Don't create per-request PRDs
```

> [!NOTE]
> For Claude-only teams, `AGENTS.md` can be replaced with `CLAUDE.md` by setting `agents.file: CLAUDE.md` in `config.yaml`. For multi-tool teams, set `agents.file: AGENTS.md` and add `tool_overrides.claude: CLAUDE.md` to also surface a Claude-specific file.

---

### `config.yaml`

Machine-readable project configuration. The skill reads this on every invocation.

```yaml
project:
  name: my-app
  summary: "One-liner about the project"

naming:
  requests: kebab-case        # slug format enforced on scaffold
  prefix: ""                  # optional prefix for request slugs

required_files:
  requests:
    - PLAN.md
    - TASKS.md

agents:
  file: AGENTS.md             # override to CLAUDE.md for Claude-only teams
  tool_overrides:             # per-tool supplementary files
    # claude: CLAUDE.md       # uncomment to also load CLAUDE.md for Claude
  default_context:
    - PRD.md
    - STACK.md
    - decisions/index.md
    # Full per-task context map lives in .spectacular/AGENTS.md

skills:
  symlink_on_init: []         # project skills to auto-symlink on init
```

---

## `specs/index.md` + `specs/` — capability specs

Ephemeral execution context. Code is authoritative; `specs/index.md` is the cheap, always-on index and flat specification files hold detail. Specs move `draft|unconfirmed → approved → implemented`, then normally archive after verified integration; implemented specs may first become superseded/deprecated. Rejected or abandoned draft/unconfirmed/approved specs archive with a reason instead of being purged. Legacy labels remain readable until an explicit, preview-first migration.

**Rules:**
- Keep the index concise; use one flat file per detailed capability (never nested folders)
- Authoritative and behavior-oriented — what the system does, not how it is implemented
- Never overwritten in place — skill snapshots before proposing edits
- Skill proposes updates when a request is archived, as a structured `SPEC-DELTA.md` the human confirms *(v1.28.0+)* — a missing delta blocks the archive closure gate

### Frontmatter

```yaml
---
status: draft | unconfirmed | approved | implemented | superseded | deprecated | archived
updated: 2026-05-11
summary: "What this capability does"
---
```

### File structure

```md
# Auth

## Purpose
What this capability does for the user.

## Requirements
- Requirement one
- Requirement two

## Scenarios
- Happy path
- Edge case

## Security considerations
- Note

## Performance expectations
- Note
```

---

## `requests/` — active and planned work

One folder per request, named with a kebab-case slug. Every request requires `PLAN.md` and `TASKS.md`. Other files are created on demand.

### Slug rules

- Derived from conversation context by the skill — shown to user before creating
- User can override at any time
- Kebab-case by default (configurable in `config.yaml`)
- Must be unique — skill proposes appending `-2` on collision

### Request lifecycle

```
planned → active → review → verified → archived
```

State lives exclusively in `PLAN.md` frontmatter. The skill reads it to surface status and propose transitions.

---

### `PLAN.md` (required)

Defines intent + plan for one request. Owns lifecycle state. Uses the **7-slot decomposition** that gives every request the same shape: goal → constraints → milestones → tasks (pointer) → dependencies → validation → deliverables.

```yaml
---
status: planned | active | review | verified
priority: high | medium | low
owner: alex
updated: 2026-05-21
summary: "One-line description of what this request changes"
related:
  - ../../PRD.md
  - specs/auth.md
---
```

Sections (in order): Goal, Why (intent), Constraints, Milestones, Tasks (pointer to TASKS.md), Dependencies, Validation, Deliverables. Optional trailing: Open questions.

**Required fields:** `status`, `updated`, `summary`. All others optional.

**Anti-pattern: never create `requests/<slug>/PRD.md`.** Product intent is project-wide and lives at the root. If a request needs to extend product intent, edit `.spectacular/PRD.md` (snapshot first).

---

### `TASKS.md` (required)

Executable implementation checklist. Grouped work items with checkboxes. The skill monitors completion as a signal for lifecycle transition proposals.

```yaml
---
updated: 2026-05-11
---
```

Format:

```md
## Group name

- [ ] Task one
- [x] Task two (completed)
```

---

### `SESSION.md`

Created automatically when a request moves to `active`. Tracks current execution state across sessions.

```yaml
---
updated: 2026-05-11
---
```

Sections: Current state, Active task, Blockers, Next actions.

Committed to git — part of the team's operational record.

---

### `RISKS.md` (on demand)

The skill proposes creating this file when a request touches auth, billing, data migrations, or anything flagged sensitive in `STACK.md`.

Sections: Risk title, Likelihood, Impact, Description, Mitigation.

---

### `VERIFY.md` (on demand)

The skill proposes creating this file for requests with user-visible behavior changes or high-stakes implementation.

Answers "did we build it correctly and safely?" — distinct from `PLAN.md` which answers "did we build the right thing?"

Sections: Manual QA checklist, Edge cases to verify, Regression checklist, Rollback validation.

---

### `specs/` (on demand)

Per-request capability specs. Same frontmatter schema as top-level `specs/<capability>.md`, but scoped to this request. Tracks its own state independently.

---

## `ideas/` — thinking scratchpad

Low-commitment, speculative. Nothing in `ideas/` is acted on automatically by the skill.

Use for: raw thoughts, market observations, UX experiments, discarded approaches, future concepts.

Each new idea is stored as `IDEA-NNN-<slug>.md`. **Promotion:** `spectacular idea promote <slug-or-alias>` scaffolds a request from the idea content and moves the idea file to `archive/ideas/`.

---

## `questions/`, `research/`, and `spikes/` — discovery records

These collections turn uncertainty into durable, linkable work:

- `QUE-NNN` records an ambiguity or blocker. It defaults to `requires_user_input: true` and uses `open | deferred | resolved`.
- `RES-NNN` records sources, evidence, and an outcome. It clears fog only when `verified`.
- `SPK-NNN` records a feasibility experiment. Its execution requires human authorization; prototype code is evidence, not implementation.

Dependencies are canonical IDs in `blocked_by:`. `spectacular wayfind status` computes unresolved fog and the dependency-ready frontier without storing either as a lifecycle status.

---

## `memories/` — operational learning

Long-term lessons from the project. Git-committed and team-visible — survives agent changes, tool changes, and team changes.

Common files: `lessons.md`, `failures.md`, `architecture-traps.md`, `recurring-bugs.md`.

**Write triggers:**
- On archive: skill reviews the completed request and proposes memory entries. Human confirms before writing.
- On demand: `spectacular remember this` — skill writes immediately on confirmation.

**Rule:** Written by the skill on human confirmation. Never by agents autonomously.

---

## `archive/` — completed work

Completed requests are moved here, not deleted. Original slug is preserved. Content is never modified after archiving.

```
archive/
├── add-team-billing/   # completed request
├── ideas/              # promoted idea files
├── questions/          # resolved QUE history
└── specs/              # implemented/rejected/obsolete detailed SPC history
```

The skill does not read `archive/` during normal operation.

Retention is lifecycle-derived: live artifacts are checkpoint-synchronized, temporary artifacts close with their owner, stale-safe archive is checked against code before reuse, and throwaway code is deleted only after outcome/evidence plus a recovery pointer are preserved.

---

## Versioning

Canonical documents (`PRD.md`, `PRINCIPLES.md`, `ARCHITECTURE.md`, `roadmaps/index.md`, `STACK.md`, `decisions/index.md`, `AGENTS.md`, `config.yaml`, `specs/index.md`, and capability specs) are never overwritten in place.

**Convention:** snapshot before editing, using the `@version` suffix:

```
PRD.md          # current — always the latest
PRD@v1.0.md     # snapshot taken before the v1.1 edit
PRD@v1.1.md     # snapshot taken before the v1.2 edit
```

Version is tracked in frontmatter (`version: 1.0`). The skill always proposes a snapshot before any edit to a canonical document — this behavior is not opt-in.

---

## `.spectacular.local/`

Personal override layer. Never committed to git. Gitignored by `spectacular init`.

Use for: local dev overrides, personal config variations, sensitive local paths.

The skill merges `.spectacular.local/` config over `.spectacular/` config when present, with local taking precedence.
