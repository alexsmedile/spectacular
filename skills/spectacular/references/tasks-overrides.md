# TASKS Overrides — TASKS-specific rules consumed by the generic engine

Loaded by `refine.md` / `review.md` when the active doc is `tasks` (per registry). TASKS.md is `mode: freeform` — no grill loop, but the engine still validates structure on review.

## Behavior

TASKS.md is created when a request is scaffolded (via `spectacular new`). The engine pre-populates it from `templates/tasks/base.md`. Users edit directly; the engine only intervenes when `spectacular tasks review` is called.

## Review gate checks (in addition to base)

| # | Check | How |
|---|---|---|
| 4 | Checklist format | All non-heading bullets use `- [ ]` or `- [x]` (not `- ` or `* `) |
| 5 | Milestone groupings | At least one `### M<N>` or `### <name>` heading present |
| 6 | Frontmatter status | `status:` matches the parent PLAN.md's status |
| 7 | Frontmatter related | `related:` includes `PLAN.md` |
| 8 | No abandoned checkboxes | No half-checked items like `- [.]`, `- [-]` (typos) |

## Refine patterns (freeform refine)

`spectacular tasks refine` walks the file and proposes:

| Pattern | Trigger | Proposed action |
|---|---|---|
| Flat task list | No `###` headings, ≥10 tasks | "Group by milestone using `### M1 — <name>`" |
| Wrong bullet syntax | `- task` (no checkbox), `* [ ]` | "Convert to `- [ ]` for consistency" |
| Stale completion mark | `[x]` task with `### M<N>` heading where M<N> isn't yet `verified` in PLAN | "Confirm this task is actually done — PLAN says M<N> is still in progress" |
| Frontmatter drift | `status:` differs from PLAN.md | "Sync TASKS.md status to PLAN.md or vice versa" |

## What TASKS refine does NOT do

- It does not add new tasks. The user owns the task list.
- It does not re-order tasks. Ordering is intentional (often dependency-driven).
- It does not delete completed `[x]` tasks. Completion history is part of the record.

## Related

- [[doc-registry]] — registry entry referencing this file
- [[refine]], [[review]] — engines that consume this
- [[plan-overrides]] — companion override for PLAN.md
