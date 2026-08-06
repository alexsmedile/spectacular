---
id: IDEA-002
type: idea
status: exploring
priority: high
owner: alex
origin: GitHub Issue alexsmedile/spectacular#25 — AFK mode should gate commits/PRs on understood end goal + explicit approval
updated: 2026-08-06
promoted_to: null
related:
  - https://github.com/alexsmedile/spectacular/issues/25
---

# IDEA-002 — afk-orchestration-authority

## Hypothesis

A user-approved AFK orchestration run should let the main agent complete its declared goal graph autonomously inside an isolated branch, while gating unexpected judgment and new durable scope.

## Context

AFK is an orchestration-level delegation granted once by the user for a
declared session goal. The main agent owns completion of that goal graph and
may coordinate subagents and autopilot work. An isolated AFK branch is a
security containment boundary: it constrains where autonomous Git mutations
can happen, but it is not the source of authority.

The isolated branch is a short-lived, run-scoped **AFK integration branch**,
`afk/<run-id>/integration`; task branches use sibling names such as
`afk/<run-id>/node-verify-cli`. This Git-valid layout is a project-defined use
of the established integration-branch pattern, not a new Git primitive. Task
branches may merge into the integration branch after verification. The
orchestrator may commit, push, and open a PR from it under the one-time AFK
delegation; protected/default, release, and deployment branches remain
human-gated destinations. The integration branch is cleaned up when the run
ends under the existing retention rules.

Within the approved goal, session lifetime, and branch, the orchestrator may
create task-level prerequisite nodes, branches, commits, pushes, pull requests,
and merges needed to complete the job. It must add a discovered task node only
when necessary, record its dependency and rationale, and may not turn that
into a new durable request, specification, or idea without human approval.

Subagents inherit only the task-level authority delegated to them. They return
their result when their task is complete; the main orchestrator reviews it in
the main session, sends independent verification when appropriate, and accepts
or rejects the result before integration.

## Confirmed direction

The authorization record is one durable, per-run Markdown file. It has a
write-once authorization block written when the human grants AFK authority and
an append-only-by-CLI event log. Each event records an actor, time, scope,
rationale, and concrete evidence reference. The record is an audit trail, not
tamper-proof storage or proof that a change or grouping is semantically correct.

## Draft contract — AFK orchestration authority

This is a spec-first proposal. It does not change the current CLI behavior or
authorize Git mutation until a human approves a resulting SPC and its request.

### 1. Authorization record schema

Each run lives at `.spectacular/afk/runs/<run-id>.md`. CLI mutation paths treat
its start block as write-once. A correction records the original value and
reason as a new event rather than silently replacing it. The record contains:

```yaml
---
id: <run-id>
type: afk-run
status: active | gated | completed | cancelled
authorized_by: <human identity>
authorized_at: <RFC-3339 timestamp>
session: <session id>
goal: <declared outcome>
goal_graph:
  - id: <declared node id>
    outcome: <bounded result>
    depends_on: [<node id>]
known_work:
  - <request, SPC, issue, or stable reference>
integration_branch: afk/<run-id>/integration
allowed_actions: [delegate, create-branch, commit, push, open-pr, merge-to-integration]
excluded_actions: [merge-to-protected, release, deploy, remote-delete]
hitl_gates: [new-durable-work, external-account, undeclared-scope, protected-destination]
---
```

The body contains an append-only-by-CLI `## Events` log. Events use stable IDs
and link to their supporting artifacts; no CLI path may silently rewrite an
earlier authorization or outcome. Plain Markdown remains human-editable:
normal Git review/history is the audit boundary, not a claim of tamper-proof
storage. Doctor validation detects malformed fields, duplicate/non-monotonic
event IDs, invalid event chronology, and missing required evidence; it cannot
prove that a human did not edit a file.

```md
### E-004 — node-added
- at: <RFC-3339 timestamp>
- actor: orchestrator
- scope: N-verify-cli
- rationale: Needed to validate the integration acceptance criterion.
- depends_on: [N-implementation]
- evidence: [PLAN.md#validation]
- result: accepted | rejected | gated
```

Required event kinds are `authorization`, `node-added`, `delegated`,
`returned`, `verification`, `acceptance`, `integration`, `git-action`, `gate`,
`resume`, `handoff`, `completion`, `cancellation`, and `cleanup`. `git-action`
records the exact branch, commit SHA, remote/PR URL when applicable, and
command/result reference, but does not assert that the commit is a complete or
correct semantic unit.

### 2. Goal graph and node additions

The declared goal graph is the complete authority boundary at run start. A
node names a bounded outcome, known dependency, and acceptance evidence—not a
blanket permission for related work.

The orchestrator may add a task-level prerequisite node only when it is needed
to complete an existing authorized node. Before dispatching it, the record must
append `node-added` with the parent/dependency, rationale, bounded outcome, and
verification route. It may not expand the product goal or create a durable
request, SPC, idea, research, or spike record. If completing the node would
need one of those durable boundaries, the run enters `gated` and asks the human
to authorize it separately.

### 2a. Soft technical decisions

The AFK delegation covers an in-scope technical choice when it preserves the
declared goal, accepted specifications, interfaces, and constraints; is
reversible; and creates no material compatibility, security, cost, or
operational commitment. The orchestrator records it as a `technical-choice`
event linked to the node, with alternatives considered, rationale, and
evidence. That record explains the choice; it does not claim the choice is
objectively correct.

Examples include selecting an established project helper instead of equivalent
inline code, choosing test-fixture structure, or splitting an authorized node
into implementation and verification work. A choice that alters user-visible
behavior, an API/schema, architecture direction, security posture, external
account use, or the product goal is not soft and must gate.

An autonomous durable technical `DEC` is allowed only when the write-once
authorization block includes `technical-decisions` and the existing narrow AFK
decision rule is satisfied: active run, technical and in scope, reversible,
evidence-backed, alternatives recorded, and no product/business trade-off.
Otherwise the orchestrator records a gate rather than a decision.

### 3. Branch containment

Each authorized run has exactly one short-lived integration branch,
`afk/<run-id>/integration`, created from a named, human-approved base. The
integration branch is a containment boundary, not authority by itself.

Task branches are sibling refs in the run namespace, for example
`afk/<run-id>/node-verify-cli`. They may merge only into their run's
integration branch after recorded verification and orchestrator acceptance.
The orchestrator may commit, push, and open a PR from the integration branch
only when those actions appear in the write-once authorization block.

Protected/default, release, deployment, and any configured protected
destination are excluded actions. The system refuses an automated merge into
them even if a provider would technically permit it. The human performs or
explicitly reauthorizes that final destination action outside this contract.

### 4. Subagent inheritance

The orchestrator delegates a closed, task-level brief containing: node ID,
bounded outcome, permitted branch, permitted actions, excluded actions,
dependencies, acceptance criteria, and return condition. A subagent inherits
only that listed scope; it cannot extend the goal graph, create durable work,
accept another agent's result, or operate outside the run branch namespace.

Subagent authority ends on `returned`, completion, cancellation, or the first
gate. The orchestrator remains responsible for recording the result and for
any subsequent Git action.

### 5. Verification and orchestrator acceptance

No task branch is eligible for integration merely because its agent reports
completion. The record must contain a `verification` event with the verifier,
independent scope, commands/checks, exit/result, evidence reference, and any
known limitations. The main orchestrator then records `acceptance` or
`rejection` after reviewing the change and verification evidence.

Only an accepted task branch may merge into `afk/<run-id>/integration`. The
integration event names both heads and its resulting SHA. The orchestrator may
commission an independent verifier where risk, coupling, or a delegated
implementation makes it appropriate; it must record why verification was
omitted when the node's declared acceptance evidence is otherwise sufficient.

### 6. Gates, failure, and recovery

The run moves `active → gated` before any action involving an unexpected
product/business choice, external account or browser/sensitive data,
undeclared scope/dependency, protected destination, breaking API/schema
decision, new durable work boundary, missing required verification, or a
conflict that cannot be resolved within the declared node.

A gate event records the blocked action, reason, safe current branch/HEAD,
options, and the exact human decision needed. `resume` records the human
authorization and any amended start-block correction as a new event; it never
silently broadens scope. A failed verification causes rejection or a new
in-scope repair node, not integration. A cancellation preserves the audit
record and branch/commit evidence; recovery uses the existing archive-ref
mechanism before any local branch deletion. Remote deletion remains human-only.

### 7. Completion and cleanup

The orchestrator may mark a run complete only when every declared or validly
added node is accepted or explicitly not-needed, all in-scope integrations are
recorded, the final verification/handoff evidence is present, and no gate is
open. Completion ends all inherited authority.

The final `handoff` event names the integration branch tip, commits, PR URL if
opened, validation evidence, outstanding human actions, and cleanup state.
Cleanup archives the integration/task branch tip to a verified local recovery
ref before deletion, records the restore command, and never deletes a remote
branch automatically. Issue #5's `session end` output remains a separate,
read-only working-tree review; it neither reads AFK authority as permission nor
records a Git mutation.

### 8. Migration from current AFK Git policy

The current model supplies useful foundations—one active/gated run, project
opt-in, explicit `--apply --yes`, non-primary branch refusal, archive-first
local cleanup, and draft-PR handoff—but it lacks a goal graph, run integration
branch, event ledger, delegated scopes, verifier acceptance, and active-run
binding for every Git action.

Migration must be explicit, preview-first, and backward-readable:

1. Keep legacy run files valid under the existing `active | gated | completed |
   cancelled` lifecycle; missing orchestration fields mean the run has only
   legacy authority.
2. Introduce a new opt-in policy version and a proposed `afk run authorize`
   dry run that renders the full write-once authorization block before requiring
   explicit human confirmation to write it.
3. Add the event ledger and branch namespaces before allowing any broader
   action. Current `afk start`, `afk pr`, and cleanup behavior remains
   compatible until a run uses the new authority version.
4. Bind commit, push, PR, and integration checks to an active new-version run,
   its integration branch, its `allowed_actions`, and its event evidence.
   Keep default/protected/release/deployment merges and remote deletion outside
   autonomous authority.
5. Retain #5's session-end commit review exactly as a no-mutation path for AFK
   and non-AFK sessions.

## Product-boundary gate

The current PRD says a multi-agent orchestration platform is out of scope.
This proposal therefore cannot be implemented merely by approving this idea.
An SPC must first make the deliberate product-boundary decision: either define
this as a minimal coordination convention consistent with that non-goal, or
explicitly revise the PRD and roadmap. The SPC must also define the exact
human confirmation syntax, audit validation rules, Git provider boundary, and
M1 test matrix before it can seed a request.

## Working plan

1. Resolve the PRD product-boundary gate and the minimal product definition.
2. Turn this contract into a draft SPC with explicit command grammar, audit
   validation, compatibility mapping, and a no-mutation implementation phase.
3. Review the SPC against normal sessions and Issue #5's read-only commit
   review before choosing any Git-mutating surface.
4. Promote only after a human accepts the authority, audit schema, and
   protected-boundary rules.

## Promoted to

—
