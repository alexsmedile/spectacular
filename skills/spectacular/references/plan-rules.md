# PLAN Overrides — PLAN-specific rules consumed by the generic engine

Loaded by `grill.md` / `refine.md` / `review.md` when the active doc is `plan` (per registry).

## Slot prompts

**Slot 1 — Goal**
> One sentence. What does this request change?
>
> Compress the request's intent into a single line. Should align with PRD's Vision or Goals — this is a slice, not a restatement.
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
> How each milestone is verified. Per-milestone checks.
>
> *Example:* "M3 passes when re-running `spectacular init` on an initialized workspace exits 0 with all-skip report."

**Slot 7 — Deliverables**
> Artifacts that ship out of this request. Concrete files, docs, behaviors.
>
> *Example:* "Updated `cli/spectacular`, updated `references/init-workflow.md`, CHANGELOG entry."

## Mini-refine patterns

| Pattern | Slots scope | Trigger | Proposed action |
|---|---|---|---|
| Vague adjective | 1, 3, 6 | Vague-word list hit | "What does '<word>' mean concretely?" |
| Unordered milestones | 3 | Items not prefixed with M1/M2/.../1./2. | "Number the milestones in order — readers need the sequence." |
| Tasks-as-milestones | 3 | Milestone text contains `implement`, `write`, `add`, `fix` as first verb | "That sounds like a task. What's the **outcome** that proves M<N> done?" |
| Missing dep link | 5 | Naked request name without `[[...]]` notation | "Wrap in `[[...]]` so the link is followable." |
| Empty validation | 6 | < 1 check per milestone | "How will you know M<N> passed?" |

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
| 7 | Validation per milestone | Slot 6 has ≥1 check per milestone in slot 3 |
| 8 | Frontmatter lifecycle field | `status:` is one of `planned | active | review | verified` |
| 9 | Deliverables non-empty | Slot 7 has ≥1 concrete artifact |

### Milestone-before-tasks ordering

PLAN structure check: milestones (slot 3) must logically precede tasks (slot 4 / TASKS.md). The engine doesn't enforce this beyond confirming both slots are present — the human ordering check happens during refine.

If TASKS.md has groupings, those groupings should reference the slot-3 milestones by name or number.

## Related

- [[doc-registry]] — registry entry referencing this file
- [[grill]], [[refine]], [[review]] — engines that consume this
- [[tasks-rules]] — companion override for TASKS.md
- [[prd-rules]] — reference example with kit-support
