---
version: 1.14
updated: 2026-08-04
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
├── AGENTS.md           # how to operate inside .spectacular/
├── STACK.md            # host project's tech choices
├── config.yaml         # machine-readable project config
├── POLICY.md           # practice layer — policies under work-phase hooks (v1.12.0+)
│
├── specs/              # specifications — canonical SPC-NNN entries + index
├── roadmaps/           # time-ordered "what's next" — roadmaps/index.md + shipped v*.md (v2.0 OKF)
├── decisions/          # verified architectural decisions — DEC-NNN-<slug>.md
├── questions/          # unresolved, user-owned blockers — QUE-NNN-<slug>.md
├── research/           # sourced investigations — RES-NNN-<slug>.md
├── spikes/             # feasibility experiments — SPK-NNN-<slug>.md
├── memories/           # long-term operational learning — memories/index.md + M<N>-<slug>.md (v2.0 OKF)
├── sessions/           # work-session log — sessions/index.md + entries (v1.5.0+; v2.0 OKF)
├── ideas/              # parked inspiration — IDEA-NNN-<slug>.md
├── visions/            # opt-in pre-request direction workspaces — VISION.md + fragments + evidence
├── requests/           # active and planned work
├── skills/             # project-specific reusable skills
├── feedbacks/          # prototyping-mode feedback entries (v1.6.0+; v2.0 OKF)
├── audits/             # bug investigations — diagnose before planning a fix (A<N>.md; v1.25.0+)
├── fixes/              # verified-fix log — logged only once resolved (F<N>.md; v1.25.0+)
├── findings/           # optional reserved advanced findings — future FND-NNN entries
├── bugs/               # optional reserved defect registry — future BUG-NNN entries
├── security/           # optional reserved security registry — future SEC-NNN entries
├── benchmarks/         # optional reserved performance evidence — future BMK-NNN entries
├── debugs/             # in-flight debug-job traces — debugs/<slug>/ (v1.26.0+; v2.0 OKF)
├── afk/                # AFK branch provenance ledger (created only after authorized start)
├── _snapshots/         # versioned snapshots of canonical docs (store name configurable; default _snapshots since v1.24.0)
│   ├── PRD/            # one folder per canonical doc, uppercase preserved
│   │   └── @v1.2.md    # filename = the content's version (couples to version:)
│   ├── specs/cli/SPEC/ # capability-spec snapshots mirror their sub-path
│   │   └── @v1.md
│   └── ROADMAP/
│       └── @v4.md
├── .last-mutation      # undo breadcrumb — gitignored, session ephemera (v1.22.0+)
└── archive/            # completed requests, historical snapshots
    └── afk-branches/   # outcome/evidence records written before local branch cleanup
```

`.spectacular.local/` — personal overrides, always gitignored, never committed. `.spectacular/` itself is fully committed to git.

---

# Root layer

The root layer is the set of **project-wide anchors** — files that describe the whole project, not any single capability or request. They are the *only* files that live loose at the root of `.spectacular/`; everything else lives inside a category folder (`specs/`, `decisions/`, `requests/`, …). This minimal-root invariant is what the OKF v2.0 layout enforces: root holds anchors plus category directories, nothing else.

The anchors are stable project grounding. They change infrequently, stay concise, avoid implementation details, and are **never overwritten in place** — snapshot before editing (see § Versioning).

| File | Purpose |
|---|---|
| `PRD.md` | Product intent — what Spectacular (or the host project) is, for whom, why |
| `PRINCIPLES.md` | Operating principles + runtime enforcement hooks |
| `ARCHITECTURE.md` | This file — `.spectacular/` structure and conventions |
| `roadmaps/index.md` | Versioned future work |
| `AGENTS.md` | Onboarding doc for any agent landing in `.spectacular/` |
| `STACK.md` | **Host project** technology and architecture choices |
| `decisions/index.md` | ADR-style log — one decision per entry, immutable |
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

## decisions/index.md

ADR-style log. One decision per entry (`decisions/D<N>-<slug>.md`), indexed by `decisions/index.md`. Immutable once written. Each entry contains:

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
    - decisions/index.md
skills:
  symlink_on_init: []         # project skills to auto-symlink into .claude/skills/
```

---

# Frontmatter conventions

Frontmatter is the skill's primary signal for reading project state. Every canonical document includes frontmatter.

## Root files (PRD, PRINCIPLES, ARCHITECTURE, roadmaps/index, AGENTS, STACK, decisions/index)

```yaml
---
version: 1.0
updated: 2026-05-11
summary: "One-sentence description of this file's purpose"
related:
  - <sibling-file>.md
---
```

## specs/<capability>.md capability specs

```yaml
---
status: draft | unconfirmed | approved | implemented | superseded | deprecated | archived
updated: 2026-05-11
summary: "What this capability does"
---
```

## requests/ PLAN.md (lifecycle state owner)

```yaml
---
status: planned | active | review | verified
hold: deferred | blocked        # optional — an orthogonal hold on the current stage
priority: high | medium | low
owner: alex
updated: 2026-05-11
summary: "What this request changes"
related:
  - specs/auth.md
depends-on:
  - other-slug
blocks:
  - another-slug
---
```

**Rules:**
- `status`, `updated`, `summary` are required
- Other fields are optional; skill warns but does not block
- `PLAN.md` frontmatter is the **single source of lifecycle state** for a request
- Capability specs at `specs/<capability>.md` track their own state independently

**Request source and activation provenance:** every newly sourced request stores `source_type: issue | spec | goal` and `source_ref`. Issue/goal sources also record explicit `sensitivity: normal | protected`; protected work cannot use the ordinary PR path. Spec-derived work additionally stores `source_spec`, `source_spec_version`, and `source_spec_digest`; the digest identifies the exact approved Markdown baseline. Activation adds `activated_at`, `activated_by`, and `activated_against` (Git commit or `uncommitted`). Do not copy an SPC, Issue body, or comments into the request. `request new --from SPC-NNN`, `--from-issue`, and `--from-goal` create PLAN/TASKS but leave them planned; the gated activation owns production authorization.

**`hold:` — orthogonal hold (b21-adjacent, 2026-07-10).** The lifecycle is a
**pure five-state chain** (`planned → active → review → verified → archived`). A
request that must pause — deferred by choice, or blocked on something external —
does **not** get a sixth status; instead it keeps its real lifecycle stage and
carries an optional `hold:` modifier (`deferred`, `blocked`, or any short reason).
This keeps the linear chain uncluttered, mirroring how the substrate keeps
non-lifecycle states (`ideas`, `feedback`) in their own collections rather than
in the request enum. Behavior: `spectacular status` renders the hold everywhere
(fleet column `planned(deferred)`, card `(hold: deferred)`, `--json` `"hold"`
field); `spectacular advance` **refuses** to move a held request until the hold
is cleared (remove the field or set `hold: none`). Sort rank keys off the raw
lifecycle status, so a held request sorts by its real stage.

### Cross-request relationship fields (v1.16.0+)

`related:`, `depends-on:`, and `blocks:` are sibling fields declaring inter-request relationships. All three accept lists of request slugs.

| Declared field | Meaning | Computed inverse |
|---|---|---|
| `related: [B]` | A and B touch each other (no ordering implied) | `related: [A]` on B |
| `depends-on: [B]` | A cannot ship before B | `required-by: [A]` on B |
| `blocks: [B]` | A must ship before B can proceed | `blocked-by: [A]` on B |

**Computed-not-stored rule:** inverse labels (`required-by:`, `blocked-by:`) are *never written* to a request's PLAN.md — they are derived at read time from the full graph of forward declarations. Storing inverses duplicates source and causes drift. A request only declares its own outbound edges.

**Archived dependencies = satisfied:** a `depends-on:` targeting an *archived* request resolves as met — shown as `depends-on: X ✓ (shipped)`, not a dangling warning. A slug matching nothing (active or archived) is dangling and flagged by `doctor links`.

**Advisory only:** these fields carry no enforcement. No locking, no auto-blocking. Conflict resolution is always human judgment.

### Roadmap ledger (v1.17.0+)

The ledger is a single markdown table at the **top of `roadmaps/index.md`**, above the first version block. It is the **only place a target version number is written** — request frontmatter carries a stable `build:` id instead; all prose references requests by slug or build id.

#### Schema

```
| build | slug | title | tier | target-version | status |
|-------|------|-------|------|----------------|--------|
| b1    | auth-backend | Auth backend | full | v1.10.0 | shipped |
| b2    | user-profile | User profile | full | v1.10.0 | shipped |
| b3    | cross-request-links | Cross-request awareness | full | v1.16.0 | active |
```

**Columns:**

| Column | Values | Notes |
|---|---|---|
| `build` | `b1`, `b2`, … | Monotonic counter, stamped at `spectacular new`, immutable |
| `slug` | kebab-case | Human identity; used in `depends-on:`/`blocks:` |
| `title` | short label | Copied from PLAN `summary:` when slotting; may drift slightly |
| `tier` | `full` · `themed` · `vision` | See tier legend below |
| `target-version` | `v1.10.0` · `tbd` | **Only place this is written.** Editable; changing it is a one-row edit. Use **`tbd`** when the build is slotted + prioritized but not yet pinned to a release — a committed sentinel, not an unfilled blank (distinct from a `<TBD>` placeholder). A build moves `tbd → vX.Y.Z` when pinned. |
| `status` | `planned` · `active` · `shipped` | Release-level; distinct from request lifecycle (see below) |

#### Tier legend

| Tier | Meaning |
|---|---|
| `full` | Near-term — detailed milestones, spec'd, on the active runway |
| `themed` | Mid-term — directional theme known, details deferred |
| `vision` | Long-horizon — direction only, no committed scope |

#### Status values (release-level)

| Status | Meaning |
|---|---|
| `planned` | Version not yet started |
| `active` | Version in progress |
| `shipped` | Version tagged and released |

These are **distinct from request lifecycle** (`planned | active | review | verified` in PLAN.md frontmatter). A request can be `verified` (done) while the ledger row is still `planned` (release hasn't shipped yet). The ledger row flips to `shipped` when the version tags — a one-time write.

#### Rules

- **Version-is-derived:** the `target-version` column is the single source of truth. No version number is written anywhere else (not in PLAN frontmatter, not in prose, not in milestone text).
- **Grouped builds:** two requests targeting the same version = two rows with the same `target-version` value. Flat table; the render groups visually at read time.
- **Human-adds-rows:** `spectacular new` stamps `build: bN` on the PLAN.md and increments `last_build:` in `config.yaml`, but does **not** insert a ledger row. The human adds the row to `roadmaps/index.md` when slotting the request into a version.
- **Gaps are normal:** if a build id is skipped (request merged into another release, abandoned), that gap is fine — like skipped Xcode build numbers.
- **Planned runway only:** the ledger tracks future/in-progress work. Shipped history lives in `CHANGELOG.md` (facts) and, as per-version planning prose, under `roadmaps/v<X.Y.Z>.md` once migrated — not in the ledger. The ledger keeps one row per build (past + future) as the compact index; see `spectacular roadmap migrate` + `specs/roadmap.md` § Index mode for how shipped prose blocks are aged out of `roadmaps/index.md` to bound context cost.

#### `build:` in PLAN.md frontmatter

```yaml
---
status: active
build: b3
summary: "What this request changes"
---
```

`build:` replaces `target_version:`. It is stamped at `spectacular new` and never changes, even if the version shifts. The version is a ledger read, not a stored copy — so reslotting a request is a one-row edit in the ledger with zero changes to the request's own files.

---

# Wayfinding knowledge layer

Spectacular keeps distinct typed Markdown collections so uncertainty is visible instead of being smuggled into plans or specifications:

- `ideas/` parks out-of-scope inspiration (`IDEA-NNN`) without expanding the active milestone.
- `questions/` holds unresolved blockers (`QUE-NNN`); records requiring user judgment are surfaced at session start.
- `research/` holds sourced investigation (`RES-NNN`) and `spikes/` holds technical feasibility experiments (`SPK-NNN`). People say `R1`; `RES` remains the only persisted research prefix.
- `decisions/` contains verified, evidence-backed architectural decisions (`DEC-NNN`), not unresolved preferences.
- `specs/` contains feature specifications (`SPC-NNN`). Their evidence-gated lifecycle is defined by `skills/spectacular/references/lifecycle-contract.md`; only `approved` specs may seed execution requests, while code remains authoritative.

Discovery is progressive, not mandatory. Inspect code/tests/docs or ask directly first; create `RES` for a bounded fact gap, `SPK` for technical feasibility, and a Vision-owned prototype fragment when human reaction to a possible experience is the evidence. Post-build feedback may own a prototype used to evaluate something already shipped. `PRT` stays reserved. A tracer bullet is an approved `SPC` with `execution_mode: tracer`: a thin, production-quality vertical slice that is retained and extended, never throwaway discovery code. “Artifact” names an output owned by one of these records or a request, not an `ART` entity or catch-all database.

Technical debt also has no parallel database. Put in-scope remediation in request tasks, likely near-term work in a roadmap `candidate` with `tbd`, uncommitted possibilities in `ideas/`, and deliberate compromises in a decision linked to the cleanup owner. A production mock is routed this way; a spike/prototype mock is disposable evidence. See `skills/spectacular/references/discovery-protocol.md`.

## Artifact retention and freshness

Freshness is derived from entity, status, and path—never duplicated as a `retention:` field:

| Class | Includes | Contract |
|---|---|---|
| **Live** | code/tests, `roadmaps/index.md`, `specs/index.md`, active requests, unresolved questions, active-release specs before execution | Re-evaluate at implementation, archive, major roadmap/architecture change, and session boundaries |
| **Temporary** | active execution specs, SESSION/request artifacts, prototypes, AFK drafts | Bounded by an owner; promote durable output or archive on closure |
| **Stale-safe** | archived specs/questions/requests, DEC/FND, completed RES/SPK, shipped `roadmaps/vX.Y.Z.md` | Historical evidence, not current truth; check against code/vendor docs before reuse |
| **Throwaway** | spike/fork branches, prototype sandbox code, scratch/generated output | Preserve outcome/evidence and recovery pointer, then delete |

Production code and executable unit/invariant tests are implementation truth after verified integration. `roadmaps/index.md` is the only live roadmap entry point; no root `ROADMAP.md` or `ROADMAP_ARCHIVE.md` exists. Older shipped prose lives in `roadmaps/vX.Y.Z.md` behind the live index.

At every human-agent session start, unresolved `requires_user_input` questions are surfaced before other work. Resolution records answer provenance and moves the QUE to `archive/questions/`; a DEC is created only for a genuine choice. Detailed specs remain synchronized through approval/action, become temporary execution context once code generation starts, and move to `archive/specs/` after verified integration or rejection. See `skills/spectacular/references/artifact-retention.md`.

Cross-references always persist canonical IDs, even when users speak in aliases such as `D1`, `Q1`, or `SPEC1`. Explicit prefixes win; naked numbers require a collection context. IDs use at least three digits. Discovery targets use `vX.Y.Z-discovery`; approved-spec execution targets use `vX.Y.Z-execution`.

Readiness is derived from canonical dependencies. Unresolved records form the **fog**; dependency-ready records form the **frontier**. The sequencer rejects invalid graphs and uses strict dependency-first topological order. `wayfind next` ranks explicit priority first, then uncertainty: user-input question, spike, research, other question, specification. Strong dependency language across PRD, roadmap, plans, and specs produces advisory doctor findings only; explicit frontmatter remains authoritative and is never silently rewritten.

# Ideas layer

A **pre-commitment workbench**, not an execution stage. Nothing in `ideas/` is acted on automatically by the skill. A committed idea can be grilled, researched, revised, and given a non-binding working plan while information or decisions remain open.

Use it for: raw thoughts, market observations, UX experiments, discarded approaches, future concepts, unresolved brainstorming, and locally useful refinement of an idea first captured elsewhere.

```txt
ideas/
├── IDEA-001-multiplayer-editor.md
├── IDEA-002-ai-memory-system.md
└── IDEA-003-growth-loops.md
```

**Rules:**
- low commitment; speculative even when it has a working plan
- one of several capture entry points: GitHub Issues, GitHub Discussions, `TODO.md`, `FEEDBACK.md`, private notes, and local/committed ideas can all begin the loop
- a capture stays authoritative where it was first recorded until the human deliberately moves the work forward; link sources rather than automatically mirroring them
- skill **proposes** saving unresolved decisions here when conversations have open branches

**Promotion to request:** Ideas are not a required gate. A request can be created directly. A draft implementation plan does not itself promote an idea: promote only after a human accepts an execution outcome and durable coordination is warranted. When deliberately promoted, the skill scaffolds the request from the idea content and moves the idea file to `archive/ideas/`.

**CLI verbs (v1.7.0+):** `spectacular idea new <slug>`, `spectacular idea list [--status <s>]`, `spectacular idea promote <slug>`. Status enum: `parked | exploring | promoted`. Full spec in [[idea-rules]]; doctor area: [[doctor-areas]] § ideas.

---

# Spec layer (specs/index.md + specs/)

The spec layer carries ephemeral execution context, not permanent system truth. Each specification follows `draft|unconfirmed → approved → implemented`; only `approved` entries may seed an execution request. After verified integration, `implemented` detail normally archives, while superseded/deprecated remain optional intermediate history. Rejected or abandoned pre-implementation specs archive with their prior status and reason. Code and executable tests are authoritative.

```txt
specs/
├── index.md                         # dense capability index
├── SPC-001-user-onboarding.md       # canonical ID + descriptive slug
└── SPC-002-roadmap.md
```

**Purpose:** defines current behavior, active capabilities, security requirements, performance expectations, user-visible behavior.

**Rules:**
- authoritative — `specs/index.md` is always relevant + cheap to load; capability specs load on demand
- current only — no past state, no future plans (future lives in `roadmaps/index.md`)
- behavior-oriented, not implementation-oriented
- modular — one flat file per capability, `specs/<capability>.md`; promote a dense `specs/index.md` bullet into its own file only when it earns it (see the v1.10.0 density refactor)
- **never overwritten in place** — skill snapshots before proposing edits (`_snapshots/specs/<capability>/SPEC/@v<ver>.md`)
- skill proposes `specs/index.md` + `specs/` updates when a request is archived; humans confirm (the spec-sync flow)

**Capability spec structure** — each `specs/<capability>.md` contains:
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

### ID-namespace convention

Spectacular uses canonical padded IDs for durable Wayfinding entities while retaining older project-local mnemonic IDs as readable compatibility forms:

| Prefix | Entity | Scope | Lives in |
|---|---|---|---|
| `M<N>` | Milestone | per-request | `TASKS.md` headings, `PLAN.md` §3/§6 |
| `DEC-NNN` | Decision | project-wide | `decisions/DEC-NNN-<slug>.md` |
| `F<N>` | Legacy verified fix | project-wide | `fixes/F<N>.md` |
| `FND-NNN` | Finding (reserved) | project-wide | `findings/FND-NNN-<slug>.md` |
| `FIX-NNN` | Fix/remediation (reserved successor) | project-wide | `fixes/FIX-NNN-<slug>.md` |
| `BUG-NNN` | Bug/defect (reserved) | project-wide | `bugs/BUG-NNN-<slug>.md` |
| `SEC-NNN` | Security vulnerability (reserved) | project-wide | `security/SEC-NNN-<slug>.md` |
| `BMK-NNN` | Benchmark result (reserved) | project-wide | `benchmarks/BMK-NNN-<slug>.md` |
| `b<N>` | Roadmap build id | project-wide | `roadmaps/index.md` ledger table |
| `A<N>` | Debug audit finding | per-debug-job | `debugs/<slug>/` trace artifacts |

**Rule:** don't invent a new prefix for an existing entity type. Reserved advanced IDs stabilize names only; no entry is allocated until its workflow ships. Use `fix1` for fixes and `fnd1` for findings; ambiguous `f1` is refused. Roadmap build IDs own `b1`, and bugs use `bug1`. `RCH`/`RSC`/`SER`/`SRC`, `ART`, `TRC`, and `DEB` are not accepted aliases or entity prefixes.

**The ID is a mnemonic, not the link.** The real linkage between a `TASKS.md` milestone and its `PLAN.md` §3/§6 counterpart is the milestone's **name** (the text after the em-dash), not the `M<N>` token. If a heading ever does drift to a non-standard prefix or a renumbered `M<N>`, matching by name still works — `doctor lifecycle` checks both: it flags an off-standard prefix as a fixable warning, but only reports a genuine chain-break when the names *also* fail to line up (so a relettered-but-still-named-the-same milestone doesn't false-positive).

Created automatically when a request moves to `active`. Captures current execution state, blockers, next actions. Committed to git — part of the team's operational record.

## RISKS.md (on demand)

Skill proposes creation when a request touches auth, billing, migrations, or anything flagged sensitive in `STACK.md`. Defines edge cases, architectural risks, mitigation plans.

Agents rarely reason about failure modes unless explicitly prompted — this file improves implementation quality significantly.

## VERIFY.md (on demand)

**On-demand only.** The skill proposes creation when the request hits the **2-of-6 rule** (see [[verify-authoring]] for full text):

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

## visions/ (on demand — pre-request direction approval)

**Opt-in and pre-request.** `spectacular imagine <slug>` creates a bounded
`.spectacular/visions/<slug>/` workspace only when product, interaction, or
experience uncertainty benefits from something a human can react to. It never
creates a request or PLAN.

```text
visions/<slug>/
├── VISION.md       # intent, north star, understanding, chosen direction, approval
├── fragments/      # strategy/story/flow/ui/arch/prototype proposals
└── evidence/       # research, spike conclusions, screenshots, recordings, notes
```

Imagine is the generative operation; Vision is the durable decision artifact.
The loop is `Understand → Imagine → Probe → React → Confirm → Derive`.
Fragments carry an explicit human reaction (`pending`, `approved`, `revise`,
`rejected`, or `superseded`). A Vision moves `draft → proposed → approved`
(or to `rejected`) only after the chosen direction and fragment dispositions
form a coherent whole. Whole-Vision approval records actor and date and permits
derivation of a **draft SPC only**. The approved SPC later creates a request and
PLAN through the normal execution gate.

A prototype is one possible Vision fragment, not a parallel entity or required
roadmap phase. A spike contributes technical feasibility evidence
(`supported|refuted|inconclusive`) but never approves product direction.
Non-visual work can use strategy, story, flow, or architecture fragments; fully
specified backend and maintenance work should skip Vision entirely.

Legacy `requests/<slug>/vision/` folders remain readable and diagnosable. New
writes use `visions/<slug>/`; no automatic migration rewrites historical human
reactions. Full rules: [`references/vision-rules.md`](../../skills/spectacular/references/vision-rules.md).

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

`.spectacular/memories/` stores long-term operational learning — `memories/index.md` + `M<N>-<slug>.md` entries.

```txt
memories/
├── index.md
├── M1-failures.md
├── M2-lessons.md
├── M3-architecture-traps.md
└── M4-recurring-bugs.md
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
archived  → moved to archive/, specs/index.md + specs/ updated, memory proposed
```

**State storage:**
- `status:` in `PLAN.md` frontmatter = request lifecycle state
- `status:` in `specs/<capability>.md` = the specification lifecycle from `skills/spectacular/references/lifecycle-contract.md`; top-level `specs/index.md` carries no per-capability status — it is the index
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
- snapshot location (v1.24.0+): `_snapshots/<DOC>/@v<ver>.md` — store dir configurable via `config.yaml` `snapshots.folder` (default `_snapshots`; pre-v1.24 default was `snapshots/`), folder per canonical doc, uppercase preserved
- **filename couples to the content's version** (v1.24.0+): a doc at `version: 1.3` snapshots to `@v1.3.md`, *then* the live doc bumps to `1.4` — the `@v` label and `version:` never drift. Docs without a `version:` field (e.g. `DESIGN.md`) use a plain `@v<N>` counter and are not version-bumped.
- sub-doc snapshots mirror their path: `specs/cli.md` → `_snapshots/specs/cli/SPEC/@v1.0.md`
- version tracked in frontmatter: `version: 1.0`
- the unversioned filename at root (`PRD.md`) always points to the current version
- applies to: root layer files, `specs/index.md`, `specs/<capability>.md` capability specs, `DESIGN.md`, `config.yaml`
- this is **default behavior** — not opt-in
- legacy snapshots at root (`PRD@v1.0.md`) continue to be read; `spectacular doctor snapshots` warns until migrated via `--fix`
- **retention (v1.24.0+):** snapshots are bounded by tiered retention — origin (`@v1`) + periodic (newest per `month`/`week` bucket) + recent (newest `keep`, default 3). `spectacular snapshot prune` removes the rest (git-rm if tracked, else `.trash/`); configured via the `snapshots:` block. `doctor snapshots` nudges when prunable snapshots accumulate.

---

# Init flow

`spectacular init` is a one-time CLI bootstrap. Detailed implementation lives in [`requests/cli-bootstrap/PLAN.md`](requests/cli-bootstrap/PLAN.md) (v0.2.x) and [`requests/smart-init/PLAN.md`](requests/smart-init/PLAN.md) (v0.3.0+).

## v0.3.0 — smart init

As of v0.3.0, init scaffolds **only what the project needs**, not all root docs.

**Always-set** (6 files + 2 dirs, scaffolded unconditionally — v0.5.0+):
- `.spectacular/PRD.md` — anchor doc; every other doc references it
- `.spectacular/specs/index.md` — system spec index (what's built right now, present tense)
- `.spectacular/config.yaml` — project name, kit identity, naming rules
- `.spectacular/<agents-file>` — onboarding doc (defaults to `AGENTS.md`)
- `.spectacular/requests/` — request folders
- `.spectacular/specs/` — per-capability specs (optional content; only when a capability outgrows a one-liner in `specs/index.md`)

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
- [roadmaps/index.md](roadmaps/index.md) — v2+ structural additions (workspaces, nested workspaces, workflows)
- [AGENTS.md](AGENTS.md) — how to operate inside this structure
