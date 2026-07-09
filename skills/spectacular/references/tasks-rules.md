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

## Checkbox schema (enforced)

Three checkbox states, all **flush-left** (no leading indent):

| Syntax | Meaning | Counted in progress |
|---|---|---|
| `- [ ]` | open | yes — denominator |
| `- [x]` | done | yes — numerator + denominator |
| `- [~]` | deferred | shown separately (`5/8 (+1 deferred)`), excluded from the open/done split |

**Indented `  - [ ]` sub-bullets are allowed** as a nested acceptance checklist
under a parent task, but are **not counted** — `status` progress counts top-level
checkboxes only, so `x/total` stays comparable across requests. Milestones group
tasks with `### M<N> — <name>` headings. `doctor` (lifecycle area) **errors** on an
active request missing `### M` headings or using a malformed checkbox (`- [.]`,
`- [-]`); `archive/` is skipped.

## Review gate checks (in addition to base)

| # | Check | How |
|---|---|---|
| 4 | Checklist format | Flush-left task bullets use `- [ ]`, `- [x]`, or `- [~]` (not `- ` or `* `). Indented `  - [ ]` sub-bullets are allowed |
| 5 | Milestone groupings | At least one `### M<N>` heading present |
| 6 | Frontmatter status | `status:` matches the parent PLAN.md's status |
| 7 | Frontmatter related | `related:` includes `PLAN.md` |
| 8 | No abandoned checkboxes | No malformed items like `- [.]`, `- [-]` (typos). `- [~]` is a valid deferred state, not a typo |

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
