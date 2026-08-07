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
spectacular request new [<slug>] --from <spec-slug-or-uuidv7>
spectacular request new <slug> --from-issue owner/repo#123 --summary "Accepted outcome" --sensitivity normal
spectacular request new <slug> --from-goal <goal-ref> [--summary "Accepted outcome"] --sensitivity normal
spectacular request list
spectacular request <slug>
spectacular request <slug> --brief [-m2]
spectacular request <slug> --full
spectacular request advance <slug>
spectacular request archive <slug>
spectacular traffic preflight <slug> [--against <slug>] [--json]
```

The older top-level forms (`new`, `requests`, `advance`, `archive`, `progress`) remain compatibility aliases.

## Request traffic preflight

`spectacular traffic preflight` is a local, read-only launch assessment. It reuses
SPC-003 / DEC-021's four outcomes: `parallel`, `conditional`, `serialized`, and
`unknown`; it neither schedules work nor changes lifecycle state. The result is
explicitly assessed as of the invocation date and must be rerun when evidence changes.

The only inputs are durable PLAN frontmatter declarations. `depends-on:`, `blocks:`,
and `conflicts-with:` serialize. A shared `release-constraints:` entry serializes
because the release or migration boundary is shared. A cross-request `related:` link
or shared `traffic-boundaries:` entry is conditional. `parallel` requires both
requests to declare complete `traffic-boundaries:` lists and for those named sets to
be disjoint. Anything else—including absent declarations—is `unknown`, never an
inference about file-level safety.

Record a confirmed relationship in the requests' own PLAN files, then rerun. The
preflight does not rewrite either request; GitHub branch/PR evidence remains optional
and is not required for the local result.

Agentic document operations are verb-first: `/spectacular grill <doc> [target]`, `/spectacular refine ...`, and `/spectacular review ...`. Document-first forms remain aliases. Before changing anything, print the resolved document and target.

## Views

- Bare `request <slug>` is the cheap overview: frontmatter outline, source spec, lifecycle mode, current milestone/task, blockers, and progress.
- `--brief` is allowed only for `status: active`. It compiles the implementation start prompt from PLAN goal/constraints/validation/deliverables, the selected open TASKS milestone, SESSION handoff, activation provenance, and the execution boundary. It is generated, never stored.
- `--milestone M2`, `-m 2`, and `-m2` are identical. Without a selector, use the first milestone containing an open top-level task. Refuse a missing milestone or a request with no open milestone.
- `--full` emits request-owned Markdown in stable order: PLAN, TASKS, SESSION, UNDERSTANDING, RISKS, VERIFY, VERIFY-LOG, SPEC-DELTA, then remaining Markdown. Linked external records are listed or referenced, never expanded into the bundle.
- `--json` wraps the chosen view in `schema: spectacular.request.v1`. Arbitrary `--artifact` and `--section` selectors are deferred; read a named file directly when necessary.

## Source and activation handoff

Every new request has one source. A spec-derived request uses only `contract: <UUIDv7>`; the merge ancestry is the approval evidence, so it does not copy a digest, revision, or provider field. Issue/goal sources use `source_type: issue | goal`, `source_ref`, and an absolute `origin:` when applicable. See [[github-work-bridge]].

`spectacular request new [slug] --from <spec-slug-or-uuidv7>` is mechanical. It requires the spec commit to be merged into `forge.shared_base` and present in the current execution branch ancestry, refuses a second live request for the same contract, derives both PLAN and TASKS, records `contract: <UUIDv7>`, and leaves the request `planned`.

`request new <slug> --from-issue <owner/repo#N|URL> --summary <outcome> --sensitivity normal|protected` creates a lean planned request and canonicalizes the Issue identity without copying its body/comments. `--from-goal <ref>` does the same for an already-defined goal. Both require explicit sensitivity classification; protected work cannot use the ordinary PR path. Review the closed outcome and boundaries before advancing either request.

`/spectacular act SPC-001` is agentic. `/spectacular SPC-001` is the unambiguous short form; `/spectacular spec act SPC-001` remains a compatibility form. The terminal only redirects to the skill—it never pretends to authorize or begin work.

The act flow:

1. Resolve the spec and require its commit to be merged into the shared base and present in the execution branch ancestry.
2. Find zero or one live request with `contract: <UUIDv7>`; refuse ambiguity. If none exists, run `spectacular request new --from <spec-slug-or-uuidv7>`.
3. Review PLAN/TASKS against the approved spec. Refuse incomplete structure, a held request, unresolved required-user questions, declared HITL gates, or silent scope additions/removals/reordering.
4. Run the `@Implementation` policy gate and satisfy understand-before-change.
5. Transition to active and record `activated_at`, `activated_by`, and `activated_against` (Git commit or `uncommitted`). Never copy a spec or remote body.
6. Retrieve `spectacular request <slug> --brief`, inspect the named code and tests, create the native session plan at finer granularity, and begin production work.

## Phase ownership

`PLAN` means discovery/planning and does not authorize production code. Active TASKS execution means build against the approved request/spec baseline. `VERIFY` runs tests and invariants and records `VERIFY-LOG.md`. The normal `review → verified` owner is `/spectacular verify`; a direct CLI transition requires passing evidence, or the explicit recorded exception `--override verify --reason <why>`.

Documentation impact is closure bookkeeping used during verification/archive, not an everyday navigation step. `spectacular afk cleanup` is advanced branch hygiene: it archives the recovery boundary before deleting an eligible local branch; it does not clean request artifacts.
