---
name: spectacular
description: |
  AI-native operational workspace for software projects. Stop losing context. Start shipping.
  Manages the full lifecycle of a .spectacular/ workspace: reads project state, proposes actions,
  scaffolds requests, manages lifecycle transitions, writes memory, archives completed work,
  and grills/refines/reviews any structured doc (PRD, PLAN, TASKS, PRINCIPLES, ARCHITECTURE,
  ROADMAP, STACK, AGENTS, DECISIONS) via a single registry-driven engine.
  Use when: opening /spectacular on any project, scaffolding a new request, archiving completed
  work, capturing a memory, snapshotting a canonical doc, onboarding to an existing workspace,
  or building any canonical doc from scratch.
  Triggers: /spectacular, spectacular status, spectacular new, spectacular archive, spectacular
  remember this, spectacular snapshot, spectacular promote, spectacular init,
  spectacular <doc>, spectacular <doc> grill, spectacular <doc> refine, spectacular <doc> review,
  spectacular prd, spectacular plan, spectacular tasks, spectacular decisions, spectacular
  principles, spectacular architecture, spectacular roadmap, spectacular stack, spectacular agents.
when_to_use: |
  Invoke on any project that has a .spectacular/ directory. Routes to reference docs based on
  the command — never loads full context, always loads minimally and progressively. The
  generalized doc verbs (grill/refine/review) apply to any doc type registered in doc-registry.md.
version: 0.3.0
category: devtools
status: published
tags: [workspace, project-management, context, agents, lifecycle, doc-writing]
---

# Spectacular Skill

AI-native operational workspace for software projects. Lean orchestrator — read this file to understand triggers and routing, then load the relevant reference doc for the actual work.

---

## Trigger detection

### Workspace lifecycle

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
| `spectacular doctor` / `spectacular doctor <area>` / `spectacular doctor --fix` | → `references/doctor.md` |
| Skill operation hits substrate failure (registry won't parse, kit malformed, etc.) | → `references/doctor.md` § Skill-invoked checks (scoped) |
| Actively working on a request | → `references/active-request.md` |

### Doc-writing (generalized — works for any registered doc)

The generalized handler matches `spectacular <doc> [<verb>]` where `<doc>` is any entry in `references/doc-registry.md`. The verb defaults based on the doc's mode and current state.

| User says | Route to |
|---|---|
| `spectacular <doc>` (no verb) | → `references/doc-registry.md` to resolve mode, then to appropriate engine |
| `spectacular <doc> grill` | → `references/grill.md` (with registry context) |
| `spectacular <doc> refine` | → `references/refine.md` (with registry context) |
| `spectacular <doc> review` | → `references/review.md` (with registry context) |

**Doc IDs registered in v0.3.0:** `prd`, `plan`, `tasks`, `principles`, `architecture`, `roadmap`, `stack`, `agents`, `decisions`.

### Verification routing (when writing PLAN.md or moving requests to review)

When grilling, scaffolding, or finalizing a PLAN.md for any request, **also route to `references/verification.md`** to decide where verification lives for this request. Two distinct decisions:

| Decision point | Route to |
|---|---|
| Scaffolding a new request (`spectacular new`) | → `verification.md` — apply 2-of-6 rule. Default: no VERIFY.md. Add `### Verification` group to TASKS.md or fill PLAN § Validation instead. |
| Grilling/refining a PLAN.md | → `verification.md` § Decision flow — confirm 2-of-6 rule result; ask user if VERIFY.md needed |
| Moving request `active → review` | → `lifecycle.md` § Verification artifact detection — pick artifact (VERIFY.md > TASKS Verification > PLAN Validation) |
| Moving request `review → verified` | → `verification.md` + `lifecycle.md` — verify every check item in the chosen artifact. **Never skip.** If VERIFY.md exists, every `- [ ]` blocks transition. |

**Critical:** "VERIFY.md is opt-in" refers to *creating the file*, not *performing verification*. Verification always runs against *some* artifact. When VERIFY.md exists it is load-bearing; do not bypass it because it's "optional."

The doc-writer engine never auto-scaffolds VERIFY.md. It is created only when:
- The 2-of-6 rule triggers during request scaffolding, AND
- The user confirms.

### Legacy PRD triggers (backwards compatible)

These map to the generalized handler with `<doc> = prd`. Behavior is identical.

| Legacy trigger | Equivalent | Routes via |
|---|---|---|
| `spectacular prd` | `spectacular prd grill` (if empty) or `spectacular prd review` (if filled) | registry → grill or review |
| `spectacular prd grill` | same | registry → `grill.md` + `prd-overrides.md` |
| `spectacular prd refine` | same | registry → `refine.md` + `prd-overrides.md` |
| `spectacular prd review` | same | registry → `review.md` + `prd-overrides.md` |

The legacy `prd-grill.md` / `prd-refine.md` / `prd-review.md` references are kept for backwards compatibility but new behavior lives in the generic engine + `prd-overrides.md`.

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
| **Workspace lifecycle** | |
| `references/status.md` | No-arg invocation — read state, build briefing, surface next action |
| `references/new-request.md` | Scaffold new request, slug rules, templates |
| `references/active-request.md` | Continue work, session state, task tracking |
| `references/lifecycle.md` | State transitions, signal detection, proactive proposals |
| `references/verification.md` | When VERIFY.md is needed (2-of-6 rule) vs folded into PLAN § Validation or TASKS § Verification |
| `references/doctor.md` | Environment/infrastructure self-check (CLI detects, skill repairs); judgment-fix repair flow |
| `references/archive.md` | Archive a request, propose current/ sync + memory entries |
| `references/memory.md` | `remember this` command, write triggers, anti-collision rules |
| `references/versioning.md` | Snapshot-before-edit rules, naming convention |
| `references/current-sync.md` | Proposing current/ updates when archiving |
| `references/scaffold-reference.md` | Canonical file templates with frontmatter stubs |
| `references/onboarding.md` | First invocation on an existing project |
| `references/init-workflow.md` | CLI init + first-time project setup |
| **Doc-writing engine (v0.3.0+)** | |
| `references/doc-registry.md` | Registry: doc-type → template + slots + mode + location + overrides |
| `references/grill.md` | Generic interactive slot-filler (consumes registry + overrides) |
| `references/refine.md` | Generic vibe→spec rewriter + append-mode handler |
| `references/review.md` | Generic quality gate runner |
| `references/prd-overrides.md` | PRD-specific rules: kit selection, slot prompts, vague-word list, gate checks |
| `references/plan-overrides.md` | PLAN-specific rules: milestone ordering, dependency-link validation |
| `references/tasks-overrides.md` | TASKS-specific rules: checklist format, frontmatter sync |
| `references/kits-contract.md` | Kit extension schema: adds-slots, modifies-slots, triggers-docs; single-kit-only in v1 |
| **Legacy PRD references (deprecated, kept for backwards compat)** | |
| `references/prd-grill.md` | Legacy — superseded by `grill.md` + `prd-overrides.md` |
| `references/prd-refine.md` | Legacy — superseded by `refine.md` + `prd-overrides.md` |
| `references/prd-review.md` | Legacy — superseded by `review.md` + `prd-overrides.md` |

---

## Templates index

| Path | Purpose |
|---|---|
| `templates/prd/base.md` | Canonical 8-slot PRD template (general-purpose) |
| `templates/prd/kits/coding.md` | Coding kit — base + stack + interfaces |
| `templates/prd/kits/product.md` | Product kit — base + user stories + metrics + distribution |
| `templates/prd/kits/content.md` | Content kit — base + audience + format + distribution |
| `templates/prd/kits/research.md` | Research kit — base + hypothesis + method + decision-being-informed |
| `templates/prd/kits/blank.md` | Blank kit — pure 8-slot base, no extras |
| `templates/plan/base.md` | 7-slot PLAN template (per-request) |
| `templates/tasks/base.md` | TASKS checklist template (per-request) |
| `templates/principles/base.md` | Operating principles + enforcement hooks |
| `templates/architecture/base.md` | `.spectacular/` structure spec |
| `templates/roadmap/base.md` | Time-ordered roadmap |
| `templates/stack/base.md` | Host project tech choices |
| `templates/agents/base.md` | Onboarding doc for `.spectacular/` agents |
| `templates/decisions/entry.md` | Single ADR entry (append-mode template) |

Project may override by placing files at `.spectacular/templates/<doc>/...` — same filenames, project-local takes precedence.
