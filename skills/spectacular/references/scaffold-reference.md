---
description: Canonical file templates with frontmatter stubs for every .spectacular/ doc type.
when_to_use: Need the frontmatter shape for a file you are scaffolding by hand.
---

# Scaffold Reference — File Templates

Canonical templates for all Spectacular files. Used by the skill when creating new files.

---

## Root layer

### PRD.md
```md
---
version: 1.0
updated: <today>
summary: "<one-line product/business intent>"
related:
  - PRINCIPLES.md
  - ARCHITECTURE.md
  - ROADMAP.md
  - AGENTS.md
  - DECISIONS.md
  - STACK.md
---

# <Project name> — Product Requirements Document

## Vision
<One paragraph — what this product is and why it exists>

## Problem
<One sentence — the concrete pain this solves. Who hurts, in what situation, how often.>

## Target users
<One primary user. Not a list, not "everyone".>

## Deliverable
<What ships at v1. The artifacts a user can touch.>

## Goals & success criteria
- <Measurable success criterion with a number and a timeframe>

## Non-goals
- <Explicit exclusion you'd push back on>

## Constraints
- <What's fixed before starting — budget, tech, policy>

## First milestone
<One concrete, demoable near-term outcome.>

## Related docs
- See `PRINCIPLES.md`, `ARCHITECTURE.md`, `ROADMAP.md`, `AGENTS.md`
```

For full PRD crafting, prefer the kit-based templates in `templates/prd/` — run `spectacular prd grill` (generic engine via `grill.md` + `prd-rules.md`) to drive the interactive flow.

### PRINCIPLES.md
```md
---
version: 1.0
updated: <today>
summary: "Operating principles for this project + runtime enforcement hooks"
related:
  - PRD.md
  - ARCHITECTURE.md
---

# <Project> — Operating Principles

## 1. <Principle name>
<Statement>

**How the skill enforces this:**
- <concrete hook>

## 2. <Principle name>
...
```

Default principles (Spectacular-aligned, customize per project):
1. Context is the system
2. Separate intent from truth
3. Small files over giant documents
4. Humans and agents share the same workspace
5. Operational memory compounds
6. Progressive disclosure
7. Three layers: intent → execution → validation
8. Humans decide, agents propose
9. Feedback ≠ verification ≠ benchmark
10. Build the smallest verified slice, full scope in mind

### POLICY.md
```md
---
version: 1.0
updated: <today>
summary: "Operating policies — the practice layer paired with PRINCIPLES.md"
---

# <Project> — Operating Policies

## @<Hook>

### <verb>-<noun>
- principle: <N>        # optional — the PRINCIPLES.md § it enforces
- severity: <block|warn>
- check: <condition that must hold>
<prose: rationale + the instruction the skill follows when injected>
```

**Always-set** (every `spectacular init`) — the practice layer. Ships 19 prefilled policies (4 block / 15 warn) filed under the 9 work-phase hooks. Severity is opt-in to blocking: a policy blocks only if it explicitly says `severity: block`. Full spec: [policies-contract.md](policies-contract.md). Hooks (locked 9): `@Init @Planning @Implementation @Verification @Archive @Debugging @Remember @Snapshot @SessionEnd`.

### ARCHITECTURE.md
```md
---
version: 1.0
updated: <today>
summary: "<one-line description of this project's structure>"
related:
  - PRD.md
  - PRINCIPLES.md
---

# <Project> — Architecture

## Layout
<Tree of the project's key directories>

## Root layer
<Files at the project root and what each is for>

## <Domain> layer
<Repeat per major domain>

## Frontmatter conventions
<Schema for canonical files>

## Lifecycle
<State transitions, if applicable>

## Versioning
<Snapshot rules>
```

### ROADMAP.md
```md
---
version: 1.0
updated: <today>
summary: "<project> roadmap — v1 status, v2 features, v3+ direction"
related:
  - PRD.md
  - ARCHITECTURE.md
---

# <Project> — Roadmap

## v1 (current)
- [ ] <shipped or in-progress item>

## v2 — <feature group>
<What it is, why it's deferred>

## v3+ — <direction>
<Long-term direction, not commitments>
```

### STACK.md
```md
---
version: 1.0
updated: <today>
summary: "Technology stack and engineering rules"
---

# Stack

## Frontend
- 

## Backend
- 

## Infrastructure
- 

## Rules
- 
```

### DECISIONS.md
```md
---
version: 1.0
updated: <today>
summary: "Architectural and product decisions log"
---

# Decisions

## <YYYY-MM-DD>

**Decision:** 
**Why:** 
**Tradeoffs:** 
```

### AGENTS.md
```md
---
version: 1.0
updated: <today>
summary: "Onboarding doc for any agent or human landing inside .spectacular/"
related:
  - PRD.md
  - PRINCIPLES.md
  - ARCHITECTURE.md
---

# AGENTS.md — Working in `.spectacular/`

## What this folder is
`.spectacular/` is an AI-native operational workspace. Read `PRD.md` before planning, `ARCHITECTURE.md` for structure, `PRINCIPLES.md` for rules.

## How to operate
1. Read frontmatter first, file bodies second
2. Load progressively — don't pre-load `specs/` or `requests/` wholesale (the top-level `SPEC.md` is cheap)
3. Snapshot before overwrite on any canonical doc (`FILE@vN.md`)
4. Propose, don't act, on irreversibles
5. Never read `archive/` during normal operation
6. Write to `memory/` only on confirmation

## Context loading by task

| Task type | Load |
|---|---|
| Planning / design | `PRD.md`, `PRINCIPLES.md`, `DECISIONS.md` |
| Refining intent / PRD work | `PRD.md`, skill refs `grill.md` / `refine.md` / `review.md` + `prd-rules.md` |
| Implementing a request | `STACK.md`, `requests/<slug>/PLAN.md`, `TASKS.md`, `SPEC.md`, relevant `specs/<capability>/SPEC.md` |
| Reviewing / QA | `requests/<slug>/VERIFY.md`, relevant `specs/<capability>/SPEC.md`, `RISKS.md` |
| Onboarding cold | `PRD.md`, `ARCHITECTURE.md`, this file |

## Available skills
- `spectacular` — workspace management

## Creating requests
Use `spectacular new <description>` — never create `requests/<slug>/PRD.md` (anti-pattern; product intent is project-wide).

## Don'ts
- Don't touch `archive/`
- Don't duplicate truth
- Don't overwrite canonical docs in place
- Don't write to `memory/` autonomously
- Don't create per-request PRDs
```

### PERSONAS.md (opt-in)
```md
---
version: 1.0
updated: <today>
summary: "Audience profiles and the user stories that drive build decisions"
related:
  - PRD.md
---

# Personas

> Opt-in. Lightweight — each persona 6-10 lines. Stories live with the persona.

## <Persona name>

**Who** — One sentence. Role + context.
**Wants to** — One sentence. The outcome they're trying to reach.
**Pain** — 1-2 bullets.
- <pain>
**Stories** — "As X, I want Y, so Z."
- As <persona>, I want <action>, so that <outcome>
**Not for** — *(optional)* Who this is explicitly NOT serving.
```

Scaffolded only when `product` or `content` kit is selected, or via `spectacular init --with personas`. Validated by `spectacular doctor personas` when present; absence is never an error.

### config.yaml
```yaml
project:
  name: <project-name>
  summary: "<one-liner>"

naming:
  requests: kebab-case
  prefix: ""

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
  symlink_on_init: []
```

---

## Requests layer

### PLAN.md (7-slot decomposition)
```md
---
status: planned
priority: medium
owner: 
updated: <today>
summary: "<one-line description>"
related:
  - ../../PRD.md
  - specs/<capability>
---

# <Request title>

## Goal
<One sentence — compressed intent. Traces back to a success criterion in the root PRD.>

## Why (intent)
<Why now? What problem does this solve? Keep tight — the full why lives in `PRD.md`.>

## Constraints
- <What's fixed before starting>

## Milestones
1. **<Milestone name>** — <demoable outcome>
2. ...

## Tasks
See [TASKS.md](TASKS.md).

## Dependencies
- <Other requests, skills, blocking decisions>

## Validation
| Milestone | How we verify it passed |
|---|---|
| 1 | <test, demo, review> |

## Deliverables
- <Artifact that ships out of this request>

## Open questions
- <Things you don't know yet>
```

PRD vs PLAN distinction: PRD is project-wide and lives at `.spectacular/PRD.md`; PLAN is per-request and lives at `requests/<slug>/PLAN.md`. **Never create `requests/<slug>/PRD.md`** — explicit anti-pattern.

### TASKS.md
```md
---
updated: <today>
---

# Tasks — <slug>

## <Group name>

- [ ] 
- [ ] 
```

### SESSION.md
```md
---
updated: <today>
---

# Session — <slug>

## Current state
<What's been done so far>

## Active task
<What's being worked on right now>

## Blockers
<Anything blocking progress — or "none">

## Next actions
- 
```

### RISKS.md
```md
---
updated: <today>
---

# Risks — <slug>

## <Risk title>

**Likelihood:** low | medium | high
**Impact:** low | medium | high

**Description:** 

**Mitigation:** 
```

### VERIFY.md
```md
---
updated: <today>
---

# Verify — <slug>

> VERIFY answers "did we build it correctly and safely?" (PLAN answers "did we build the right thing?")
> Checks are TYPED — each is verified by its own authority. Walked by `spectacular verify <slug>` (see [[verify]]).

<!--
Five check kinds (deterministic → judgment → human):
  executable  `run: <cmd>`   exit 0 = pass (external command)
  assertable  {assert}        agent checks a binary property of files/state
  judgable    {judge}         LLM reasons over named artifacts
  observable  {observable}    human looks & confirms (passive) — the default if untagged
  manual      {manual}        human performs an action, then confirms (active)

Two shapes (prefer section-grouping for uniform phases; inline for mixed):
  Section: `## Title {kind}` applies to ALL checks under it (absolute).
           `## Title {run}` → each line IS the command.
  Inline:  tag per line; `run:` or {assert}/{judge}/{manual}/{observable}.
-->

## Automated {run}
- [ ] 

## Properties {assert}
- [ ] 

## Manual QA {observable}
- [ ] 

## Actions {manual}
- [ ] 
```

### VERIFY-LOG.md
```md
---
updated: <today>
---

# Verify log — <slug>

<!-- Append-only. One ## entry per `spectacular verify` walk. The [kind] tag
     records which authority confirmed each check. Written by the walk; see [[verify]]. -->

## <YYYY-MM-DD HH:MM> — walk (<P> passed, <B> blocked, <S> skipped)

- ✓ [exec] <check> — `<cmd>` exit 0
- ✓ [assert] <check> — property holds: <what was checked>
- ✓ [judge] <check> — <agent reasoning + artifact seen>
- ✓ [observe] <check> — <evidence the human gave>
- ✓ [manual] <check> — performed <action>, result: <result>
- ✗ [exec] <check> — BLOCKED: `<cmd>` exit 1 — <stderr tail>
- ⊘ [observe] <check> — skipped
**Outcome:** <verified | stayed in review — N blockers>
```

---

## Specs layer (system truth)

### .spectacular/SPEC.md (always-on index)
```md
---
version: 1.0
updated: <today>
summary: "Index of what this system actually is and how it behaves right now"
related:
  - PRD.md
  - ARCHITECTURE.md
---

# <Project> — System Spec

## What this system is
<One paragraph, present tense.>

## Capabilities
- <one bullet per capability; link out to specs/<capability>/SPEC.md only when needed>
```

### Per-capability layer (optional)

### specs/<capability>/SPEC.md
```md
---
status: draft
updated: <today>
summary: "<one-line description of what this capability does>"
---

# <Capability name>

## Purpose
<What this capability does>

## Requirements
- 

## Scenarios
- 

## Security considerations
- 

## Performance expectations
- 
```
