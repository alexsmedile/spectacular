---
doc-id: plan
mode: grill
location: .spectacular/requests/<slug>/PLAN.md
scope: per-request
template: templates/plan/base.md
slots: [Goal, Constraints, Milestones, Tasks, Dependencies, Validation, Deliverables]
snapshot-on-edit: false
summary: "Request-scoped plan — 7-slot decomposition (owns lifecycle state)"
status: active
---

# PLAN Rules — PLAN-specific rules consumed by the skill

Loaded by `grill.md` / `refine.md` / `review.md` when the active doc is `plan` (per doc-index).

## Canonical section headings (enforced)

The 7 required sections use **unnumbered** headings and must appear **in order**:

```
## Goal
## Constraints
## Milestones
## Tasks
## Dependencies
## Validation
## Deliverables
```

Extra sections (`## Understanding`, `## Decisions`, or request-specific ones) may
appear **between** them. `doctor` (lifecycle area) **errors** on an active request
whose PLAN is missing a required heading or has them out of order; `archive/` is
skipped. `doctor --fix` rewrites the legacy numbered form (`## 1. Goal`) to the
unnumbered form. The "Slot N" labels below are authoring ordinals, not heading
text — the heading is `## Goal`, not `## 1. Goal`.

## Slot prompts

**Slot 1 — Goal**
> One sentence. What does this request change?
>
> Compress the request's intent into a single line. Names or references the PRD goal / success criterion it serves — this is a slice, not a restatement (gate check 11).
>
> *Example:* "Add a CLI command that scaffolds the minimal `.spectacular/` set, defaulting to PRD + SPEC + requests/ + specs/ + config.yaml + AGENTS.md only."

**Slot 2 — Constraints**
> What's fixed before you start? Inherited from PRD/STACK/PRINCIPLES + request-specific limits.
>
> *Example:* "Bash-only CLI (no new language deps). Backwards compatible with existing workspaces. No `--force` flag."

**Slot 3 — Milestones**
> Ordered, demoable checkpoints. Outcomes, not tasks.
>
> Each milestone is something someone can see working. 3-7 milestones for a typical request.
>
> *Example:*
> - M1 — Always-set defined (doc rationale shipped)
> - M2 — Flag interface (--kit, --with, --minimal)
> - M3 — Pre-flight non-overwrite
> - M4 — Interactive mode
> - M5 — Kit consumption
> - M6 — Dogfood

**Slot 4 — Tasks**
> Pointer to `TASKS.md`. Tasks are the executable checklist; this slot just confirms TASKS.md exists and groups by milestone.

**Slot 5 — Dependencies**
> Other requests, skills, blocking decisions. Use `[[request-slug]]` notation.
>
> *Example:* "Hard dep on [[doc-writer]] (needs registry). Touches [[cli-bootstrap]]."

**Slot 6 — Validation**
> How each milestone is verified. Per-milestone checks. Each check states its **authority**: a `run:` command, an assertable property, a judgable artifact, or a human-observable behavior (see [[verify]] check kinds). A check with no authority can't fail. Aspiration verbs (`improve`, `enhance`, `optimize`, `handle gracefully`) are not checks.
>
> *Example:* "M3 passes when re-running `spectacular init` on an initialized workspace exits 0 with all-skip report."

**Slot 7 — Deliverables**
> Artifacts that ship out of this request. Concrete files, docs, behaviors.
>
> *Example:* "Updated `cli/spectacular`, updated `references/init-workflow.md`, CHANGELOG entry."

## Understanding slot (policy-gated, not one of the 7)

PLAN has an **optional** `## Understanding` section with three subheads — `### How it works now`, `### What changes`, `### What stays the same`. It is **not** one of the 7 required authoring slots (grill/review don't demand it during planning), but the `understand-before-change` policy (`@Implementation`, severity `block`) **requires it filled before `planned → active`**.

> The gate is satisfied by **either** a filled `## Understanding` in PLAN.md **or** a `requests/<slug>/UNDERSTANDING.md` with the same three subheads (the VERIFY.md 2-of-N pattern). Escalate to the standalone file for a large request. There is no `ANALYSIS.md`.

When the skill is about to promote a request to `active` and the section is empty, it fills it by interviewing: *how does the touched system work today / what does this change / what does it deliberately leave alone?* See [policy-injection.md](policy-injection.md).

## Decisions section (not one of the 7)

PLAN carries an unnumbered `## Decisions` section for **request-scoped** design calls — the destination [[decisions-rules]]'s routing table points at. Format: *chose X over Y — because Z*. Rejected alternatives stay listed; deleting them re-litigates them later. Project-wide calls go to `DECISIONS.md` via `spectacular decide` instead. Empty is valid (gate check 10 only inspects entries that exist).

When findings invalidate part of a live PLAN, don't rewrite history — use the supersession convention in [[active-request]] § Superseding a live plan.

## Mini-refine patterns

| Pattern | Slots scope | Trigger | Proposed action |
|---|---|---|---|
| Vague adjective | 1, 3, 6 | Vague-word list hit | "What does '<word>' mean concretely?" |
| Unordered milestones | 3 | Items not prefixed with M1/M2/.../1./2. | "Number the milestones in order — readers need the sequence." |
| Tasks-as-milestones | 3 | Milestone text contains `implement`, `write`, `add`, `fix` as first verb | "That sounds like a task. What's the **outcome** that proves M<N> done?" |
| Missing dep link | 5 | Naked request name without `[[...]]` notation | "Wrap in `[[...]]` so the link is followable." |
| Empty validation | 6 | < 1 check per milestone | "How will you know M<N> passed?" |
| Authority-less check | 6 | Check has no run/assert/judge/observable anchor, or leads with an aspiration verb (`improve`, `enhance`, `optimize`) | "What command, property, or observable behavior decides this? A check that can't fail isn't a check." |
| Goal restates request | 1 | Goal ≈ frontmatter `summary` reworded | "That's the request restated. What does it *change*, traced to which PRD goal?" |

## Vibe → spec rewrite tables (refine mode)

### Tasks framed as milestones

| Vibe | Spec |
|---|---|
| "Implement the flag" | "M2 — `--kit`, `--with`, `--minimal` flags wired; `--help` updated" |
| "Write the docs" | "M5 — `init-workflow.md` reflects new behavior; example invocations included" |
| "Fix the bug" | "M3 — re-running init on existing workspace exits 0; no overwrites" |

### Vague milestones → demoable outcomes

| Vibe | Spec |
|---|---|
| "MVP" | "[NEEDS CLARIFICATION: which user can do which thing by M1?]" |
| "Initial version" | "[NEEDS CLARIFICATION: concrete demo scenario for the first milestone]" |
| "Working prototype" | "10 invited users complete the core flow without errors" |

## Review gate checks (in addition to base)

| # | Check | How |
|---|---|---|
| 4 | Milestones present | Slot 3 has ≥2 ordered items |
| 5 | TASKS.md exists | File at `requests/<slug>/TASKS.md` exists |
| 6 | Dependencies use `[[...]]` | All cross-request references wrapped |
| 7 | Validation per milestone | Slot 6 has ≥1 check per milestone in slot 3, and each check names an authority (a `run:` command, assertable property, judgable artifact, or observable behavior) — a bare aspiration ("works correctly", "improved") fails |
| 8 | Frontmatter lifecycle field | `status:` is one of `planned | active | review | verified` |
| 9 | Deliverables non-empty | Slot 7 has ≥1 concrete artifact |
| 10 | Decisions name alternatives | Each `## Decisions` entry names an alternative (contains "over", "not", or "instead of"). An empty Decisions section passes — no decisions yet is valid |
| 11 | Goal traces to PRD | Slot 1 names or references the PRD goal/success-criterion it serves; a Goal that only re-words the frontmatter `summary` fails |

### Milestone-before-tasks ordering

PLAN structure check: milestones (slot 3) must logically precede tasks (slot 4 / TASKS.md). The skill doesn't enforce this beyond confirming both slots are present — the human ordering check happens during refine.

If TASKS.md has groupings, those groupings should reference the slot-3 milestones by name or number.

## The 2-of-6 rule (VERIFY.md decision at scaffold/grill time)

> Compact copy for scaffold/grill journeys — canonical source: `verify-authoring.md` ("The 2-of-6
> rule"; authoring-time doc, split from verify.md in b30). If this table and verify-authoring.md
> ever disagree, verify-authoring.md wins.

Scaffold a standalone `VERIFY.md` when **at least two** are true:

| # | Criterion |
|---|---|
| 1 | **User-visible change** — behavior change observable outside the team |
| 2 | **High reversibility cost** — migrations, schema changes, destructive ops |
| 3 | **Multi-surface verification** — beyond tests: manual QA, browser, prod smoke |
| 4 | **Risk surface non-trivial** — auth, billing, payments, PII, data integrity |
| 5 | **External contract change** — public API, library shape, plugin interface, CLI flags |
| 6 | **Rollback plan exists** — a non-trivial undo that itself needs verification |

**Default: no file** — fewer than two → checks live in PLAN § Validation (descriptive properties)
or a TASKS `### Verification` group (step-by-step). "No file" ≠ "no verification": every request
still reaches `verified` through *some* artifact. The file is opt-in; the practice is not.

## Frontmatter schema (v1.17.0+)

Required: `status`, `updated`, `summary`.
Optional stable identity: `build: bN` — stamped by `spectacular new`, immutable. Do not edit.
Removed: `target_version:` — version is derived from the roadmap ledger, not stored in the request. Do not add it back.

**Version-in-prose rule:** milestone text, validation lines, and dependency chains must not contain hardcoded version numbers (`v1.x.y`). Reference requests by slug (`depends-on: cross-request-links`) or by build id (`b7`). The version lives only in the ledger table in `ROADMAP.md`. If a prose line says "plugin bump to v1.17.0" or "manifests at v1.16.0", flag it during refine — replace with "target release" or "the version in the ledger".

## Related

- [[doc-index]] — registry entry referencing this file
- [[grill]], [[refine]], [[review]] — skill flows that consume this
- [[tasks-rules]] — companion override for TASKS.md
- [[prd-rules]] — reference example with kit-support
