---
doc-id: tasks
mode: stub
location: .spectacular/requests/<slug>/TASKS.md
scope: per-request
template: templates/tasks/base.md
snapshot-on-edit: false
summary: "Executable checklist for one request"
status: active
---

# TASKS Rules — TASKS-specific rules consumed by the skill

Loaded by `refine.md` / `review.md` when the active doc is `tasks` (per doc-index). TASKS.md is `mode: stub` — no grill loop, but the skill still validates structure on review.

## Behavior

TASKS.md is created when a request is scaffolded (via `spectacular new`). The skill pre-populates it from `templates/tasks/base.md`. Users edit directly; the skill only intervenes when `spectacular tasks review` is called.

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

- [[doc-index]] — registry entry referencing this file
- [[refine]], [[review]] — skill flows that consume this
- [[plan-rules]] — companion override for PLAN.md
