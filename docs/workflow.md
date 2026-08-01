---
title: Workflow
description: The normal Spectacular loop after installation — init, briefing, requests, lifecycle, archive.
section: ""
status: stable
since: 0.1.0
updated: 2026-08-01
---

# Workflow Guide

This guide shows the normal Spectacular loop after installation. The CLI creates the workspace once. The skill operates it during day-to-day work.

---

## 1. Initialize a project

From the project root:

```bash
spectacular init
```

This creates the always-set workspace files:

```text
.spectacular/
├── PRD.md              # product intent
├── POLICY.md           # practice-layer policies
├── AGENTS.md           # onboarding for agents
├── config.yaml
├── specs/
│   └── index.md        # current system truth
└── requests/
```

Kits or `--with` add focused context such as `STACK.md`, `ARCHITECTURE.md`, `PRINCIPLES.md`, `roadmaps/index.md`, and `decisions/index.md` when the project needs them.

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
- `.spectacular/roadmaps/index.md` — versioned future work (v1 / v2 / v3+)
- `.spectacular/STACK.md` — host project's technology choices and engineering rules
- `.spectacular/decisions/index.md` — ADR-style log of decisions and tradeoffs
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

## 4. Clear discovery fog when the path is uncertain

Not every idea should immediately become implementation work. Use typed records to keep exploration explicit:

```text
idea → questions / research / spikes → user decision → draft|unconfirmed spec → approved spec → request
```

- Park out-of-scope inspiration with `spectacular idea new <slug>`.
- Record a human-owned ambiguity with `spectacular question new <slug>`.
- Gather sources and evidence with `spectacular research new <slug>`.
- Validate feasibility with `spectacular spike new <slug>`; execution requires human approval.
- Inspect unresolved fog and dependency-ready work with `spectacular wayfind status` and `spectacular wayfind next`.
- Inspect strict dependency order with `spectacular wayfind order`; invalid graphs refuse sequencing until repaired.

Each record receives a stable canonical ID such as `QUE-001` or `RES-001`. Compact aliases such as `q1` and `r1` work as input, while saved cross-references use canonical IDs.

When discovery converges, create a collaborative draft (or an AFK unconfirmed specification), approve it with the human, then act on it:

```bash
spectacular spec new team-billing --summary "Bill teams by active seat" --target-version v1.0.0-discovery
spectacular spec approve s1 --evidence "Pricing and architecture approved" --target-version v1.0.0-execution
spectacular spec act s1
```

Only an `approved` specification can seed an implementation request. After verified integration, `spec implement --verified-against <commit|build>` records the historical evidence point; code remains authoritative. If an open loop is not relevant now, defer it with a reason and optional review point.

Wayfinder language maps to the same gates: “park this idea” creates an idea, “put it on ice” defers with a reason, “find your way to…” shows prerequisites, and “act on goal…” still requires an approved specification. During implementation, park unexpected discoveries instead of adding them to the active milestone. Run `spectacular doctor wayfinding` to surface inferred dependency gaps and discovery/execution target inversions; findings are proposals, never automatic roadmap edits.

For explicitly authorized AFK work, create a durable goal-scoped `spectacular afk run`, inspect `spectacular afk status`, and propose the branch class before creating it. Draft specs, spikes, forks, and approved execution stay isolated. Cleanup creates a durable Git archive ref before local deletion; remote deletion and merge remain human actions. A verified handoff may open `[Spectacular] Executed: <version> - <name>` only after the request, source spec, verification log, and fresh tests satisfy their gates.

## 5. Create a request

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

## 6. Work from the request folder

During implementation, the request folder is the operational center:

- `PLAN.md` answers what is being built and why.
- `TASKS.md` tracks execution.
- `SESSION.md` records current state and handoff notes.
- `RISKS.md` is useful for sensitive work such as auth, billing, data migrations, or security changes.
- `VERIFY.md` is useful when user-visible behavior or regressions need explicit checks.

Keep request docs focused on the request. Do not prematurely rewrite `specs/index.md` or capability specs while work is still in progress.

---

## 7. Use the lifecycle

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

## 8. Record implementation evidence after completion

Code is the source of implemented behavior. Specs are execution context and may become stale. After a verified request, record `implemented_at` and `verified_against`; do not treat the status as continuous synchronization.

When archiving a completed request, the skill should propose updates such as:

- add or update a bullet in `specs/index.md`
- create a new `specs/<capability>.md` (only when the bullet outgrows one line)
- update an existing capability spec
- approve a draft before execution, or mark an approved spec implemented against a concrete commit/build
- leave unaffected specs unchanged

Canonical docs and `specs/` files should be snapshotted before edits:

```text
specs/billing.md
_snapshots/specs/billing/SPEC/@v1.0.md
```

---

## 9. Capture operational memory

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

Memory is team-visible and committed to git under `.spectacular/memories/`. Do not use it for personal notes or secrets.

---

## 10. Archive completed requests

After a request is verified:

```text
spectacular archive add-team-billing
```

The skill should:

1. review the request state
2. write a spec delta to `SPEC-DELTA.md` (`### ADDED` / `### MODIFIED` / `### REMOVED`, or `NONE — <why>`) and propose the matching `specs/index.md` / `specs/` updates
3. propose memory entries if useful
4. move the request from `requests/` to `archive/`

Archiving runs a **closure gate** *(v1.28.0+)*: it blocks on open `TASKS.md` boxes, an unwalked `VERIFY.md`, or a missing `SPEC-DELTA.md` — each bypassable once with `--override <check> --reason "<text>"`, recorded on the archived plan. Archived requests are not deleted. They are also not read during normal `/spectacular` status briefings.

---

## Practical rhythm

For long-running projects, a useful rhythm is:

1. Run `/spectacular` at the start of a session.
2. Answer surfaced human blockers or defer them deliberately.
3. Clear high-uncertainty discovery nodes before routine implementation.
4. Confirm the spec, then act on it to create a request.
5. Work from one active request and keep `TASKS.md` / `SESSION.md` current.
6. Verify before changing request state to `verified`.
7. Archive completed work.
8. Update `specs/index.md` / `specs/` (via a `SPEC-DELTA.md`) only when behavior has actually changed.
9. Write memory only for lessons with future value.
