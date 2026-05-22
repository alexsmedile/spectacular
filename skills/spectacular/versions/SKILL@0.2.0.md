---
name: spectacular
description: |
  AI-native operational workspace for software projects. Stop losing context. Start shipping.
  Manages the full lifecycle of a .spectacular/ workspace: reads project state, proposes actions,
  scaffolds requests, manages lifecycle transitions, writes memory, and archives completed work.
  Use when: opening /spectacular on any project, scaffolding a new request, archiving completed
  work, capturing a memory, snapshotting a canonical doc, or onboarding to an existing workspace.
  Triggers: /spectacular, spectacular status, spectacular new, spectacular archive, spectacular
  remember this, spectacular snapshot, spectacular promote, spectacular init, spectacular prd,
  spectacular prd refine, spectacular prd review.
when_to_use: |
  Invoke on any project that has a .spectacular/ directory. Routes to reference docs based on
  the command — never loads full context, always loads minimally and progressively.
version: 0.2.0
category: devtools
status: published
tags: [workspace, project-management, context, agents, lifecycle]
---

# Spectacular Skill

AI-native operational workspace for software projects. Lean orchestrator — read this file to understand triggers and routing, then load the relevant reference doc for the actual work.

---

## Trigger detection

| User says / context | Route to |
|---|---|
| `/spectacular` with no args | → `references/status.md` |
| `spectacular status` | → `references/status.md` |
| `spectacular new <description>` | → `references/new-request.md` |
| `spectacular archive <slug>` | → `references/archive.md` |
| `spectacular remember this` | → `references/memory.md` |
| `spectacular promote <idea>` | → `references/new-request.md` |
| `spectacular snapshot <file>` | → `references/versioning.md` |
| First invocation on existing `.spectacular/` project | → `references/onboarding.md` |
| `spectacular init` (CLI context) | → `references/init-workflow.md` |
| `spectacular prd` or `spectacular prd grill` | → `references/prd-grill.md` |
| `spectacular prd refine` | → `references/prd-refine.md` |
| `spectacular prd review` | → `references/prd-review.md` |
| Actively working on a request | → `references/active-request.md` |

---

## State awareness

Before any action, read frontmatter from:
1. `.spectacular/config.yaml` — project config, naming rules
2. `.spectacular/AGENTS.md` — **authoritative** context-loading rules per task type; follow its table over guessing
3. Root canonical docs — `PRD.md` (intent), `PRINCIPLES.md` (rules), `ARCHITECTURE.md` (structure), `ROADMAP.md` (time), `STACK.md` (host tech), `DECISIONS.md` (ADR log)
4. `current/` — canonical capability specs (read summaries/status only unless task requires depth)
5. `requests/*/PLAN.md` — active work (read all frontmatter for status briefing)

Load **only** what the task needs (principle 6 — progressive disclosure). For planning, PRD + PRINCIPLES + DECISIONS. For implementation, STACK + PLAN + TASKS + relevant `current/`. For review, VERIFY + RISKS + capability specs. AGENTS.md owns the full table.

Never read `archive/` during normal operation.

---

## Canonical rules (always apply)

- **Never overwrite canonical documents in place** — snapshot first (`PRD@v1.0.md`). See `references/versioning.md`.
- **Lifecycle state** lives in `PLAN.md` frontmatter (`status: planned | active | review | verified`).
- **Capability state** lives in `current/<capability>.md` frontmatter (`status: stable | draft | deprecated`).
- **Slugs** are kebab-case, skill-derived, user-overridable, uniqueness enforced.
- **Memory** (`spectacular remember this`) writes to `.spectacular/memory/` — git-committed, team-visible. Never to `.claude/` memory.
- Be proactive: surface stale state, propose lifecycle transitions, flag blocked requests.

---

## Output format

Conversational briefing with a minimal embedded table. Never a raw dump. Identify the single highest-priority next action and ask what the user wants to do.

---

## References index

| File | Purpose |
|---|---|
| `references/status.md` | No-arg invocation — read state, build briefing, surface next action |
| `references/new-request.md` | Scaffold new request, slug rules, templates |
| `references/active-request.md` | Continue work, session state, task tracking |
| `references/lifecycle.md` | State transitions, signal detection, proactive proposals |
| `references/archive.md` | Archive a request, propose current/ sync + memory entries |
| `references/memory.md` | `remember this` command, write triggers, anti-collision rules |
| `references/versioning.md` | Snapshot-before-edit rules, naming convention |
| `references/current-sync.md` | Proposing current/ updates when archiving |
| `references/scaffold-reference.md` | Canonical file templates with frontmatter stubs |
| `references/onboarding.md` | First invocation on an existing project |
| `references/init-workflow.md` | CLI init + first-time project setup |
| `references/prd-grill.md` | Interactive 6-slot PRD crafting (one Q at a time) |
| `references/prd-refine.md` | Vibe→spec rewrite patterns + `[NEEDS CLARIFICATION]` markers |
| `references/prd-review.md` | Pass/fail quality gate for PRD.md |

---

## Templates index

| Path | Purpose |
|---|---|
| `templates/prd/base.md` | Canonical 6-slot PRD template (general-purpose) |
| `templates/prd/kits/coding.md` | Coding kit — base + stack + interfaces |
| `templates/prd/kits/product.md` | Product kit — base + user stories + metrics + distribution |
| `templates/prd/kits/content.md` | Content kit — base + audience + format + distribution |
| `templates/prd/kits/research.md` | Research kit — base + hypothesis + method + deliverable |
| `templates/prd/kits/blank.md` | Blank kit — pure base, no extras |

Project may override by placing files at `.spectacular/templates/prd/...` — same filenames, project-local takes precedence.
