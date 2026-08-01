---
description: Canonical request namespace, compiled implementation brief, approved-spec handoff, and native agent planning boundary.
when_to_use: Creating, inspecting, activating, resuming, advancing, verifying, or archiving a request.
---

# Request Workflow

Spectacular owns durable intent and evidence. The coding agent owns the current session's implementation plan. Do not copy one layer into the other one-for-one.

| Layer | Owner | Lifetime | Granularity |
|---|---|---|---|
| `PLAN.md` + `TASKS.md` | Spectacular | Cross-session, committed | Goal, constraints, milestones, durable checkpoints |
| Native Codex/Claude plan | Coding agent | Current session | Concrete edits, checks, and sequencing for the active milestone |
| Subagent brief | Dispatching agent | One bounded delegation | Narrow task, named files, acceptance check, return contract |

## Canonical commands

Mechanical entity operations are noun-first:

```text
spectacular request new [<slug>] --from SPC-001
spectacular request list
spectacular request <slug>
spectacular request <slug> --brief [-m2]
spectacular request <slug> --full
spectacular request advance <slug>
spectacular request archive <slug>
```

The older top-level forms (`new`, `requests`, `advance`, `archive`, `progress`) remain compatibility aliases.

Agentic document operations are verb-first: `/spectacular grill <doc> [target]`, `/spectacular refine ...`, and `/spectacular review ...`. Document-first forms remain aliases. Before changing anything, print the resolved document and target.

## Views

- Bare `request <slug>` is the cheap overview: frontmatter outline, source spec, lifecycle mode, current milestone/task, blockers, and progress.
- `--brief` is allowed only for `status: active`. It compiles the implementation start prompt from PLAN goal/constraints/validation/deliverables, the selected open TASKS milestone, SESSION handoff, activation provenance, and the execution boundary. It is generated, never stored.
- `--milestone M2`, `-m 2`, and `-m2` are identical. Without a selector, use the first milestone containing an open top-level task. Refuse a missing milestone or a request with no open milestone.
- `--full` emits request-owned Markdown in stable order: PLAN, TASKS, SESSION, UNDERSTANDING, RISKS, VERIFY, VERIFY-LOG, SPEC-DELTA, then remaining Markdown. Linked external records are listed or referenced, never expanded into the bundle.
- `--json` wraps the chosen view in `schema: spectacular.request.v1`. Arbitrary `--artifact` and `--section` selectors are deferred; read a named file directly when necessary.

## Approved-spec handoff

`spectacular request new [slug] --from SPC-001` is mechanical. It requires an approved SPC, refuses a second live request for the same source, derives both PLAN and TASKS in approved requirement order, records the source version/digest/scaffold Git baseline, and leaves the request `planned`.

`/spectacular act SPC-001` is agentic. `/spectacular SPC-001` is the unambiguous short form; `/spectacular spec act SPC-001` remains a compatibility form. The terminal only redirects to the skill—it never pretends to authorize or begin work.

The act flow:

1. Resolve the SPC and require `status: approved`.
2. Find zero or one live request with `source_spec: SPC-NNN`; refuse ambiguity. If none exists, run `spectacular request new --from SPC-NNN`.
3. Review PLAN/TASKS against the approved spec. Refuse incomplete structure, a held request, unresolved required-user questions, declared HITL gates, or silent scope additions/removals/reordering.
4. Run the `@Implementation` policy gate and satisfy understand-before-change.
5. Transition to active and record flat provenance: `source_spec_version`, `source_spec_digest`, `activated_at`, `activated_by`, and `activated_against` (Git commit or `uncommitted`). Never copy the spec body.
6. Retrieve `spectacular request <slug> --brief`, inspect the named code and tests, create the native session plan at finer granularity, and begin production work.

## Phase ownership

`PLAN` means discovery/planning and does not authorize production code. Active TASKS execution means build against the approved request/spec baseline. `VERIFY` runs tests and invariants and records `VERIFY-LOG.md`. The normal `review → verified` owner is `/spectacular verify`; a direct CLI transition requires passing evidence, or the explicit recorded exception `--override verify --reason <why>`.

Documentation impact is closure bookkeeping used during verification/archive, not an everyday navigation step. `spectacular afk cleanup` is advanced branch hygiene: it archives the recovery boundary before deleting an eligible local branch; it does not clean request artifacts.
