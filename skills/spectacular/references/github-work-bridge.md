---
description: Route GitHub work to direct execution, a durable request, or spec-first discovery; hand coordinated work to a PR and reconcile state safely.
when_to_use: Triaging an Issue, creating Issue/goal-derived work, opening or readying a request PR, or reconciling GitHub with Spectacular.
---

# GitHub work bridge

GitHub is the collaborative work queue. Spectacular is the durable destination and coordination layer when a change earns it. Never copy remote bodies or comments into a second local database.

## Authority boundary

The bridge brings Spectacular's **method** to GitHub—capture, clarification,
research, decision, and route selection—without rebuilding Spectacular's local
data model in labels, comments, or automation.

| GitHub owns | Spectacular owns |
|---|---|
| Capture, discussion, ownership, Issues/Discussions, PRs, checks, merge, collaboration, notifications | Reasoning, durable local context, decisions, plans, validation, coordination |

An Issue is the canonical workbench while a shared idea is being discussed.
An SPC or request becomes canonical only for the explicit agreement or
execution plan it owns. Link the two at promotion; never maintain copied
bodies, duplicate checklists, or synchronized lifecycle labels.

### Issue capture page schema

Use this as a lightweight Issue form or body shape when GitHub is the capture
surface. It is a human-readable page schema, not a new CLI schema or a
mandatory workflow state machine.

```md
## Capture

What sparked this? What problem, observation, or opportunity might exist?

## Current hypothesis

What might we change or build? This is allowed to be wrong.

## Open questions

- What information is missing?
- Which decision needs a human owner?
- What evidence would change the direction?

## Working plan

Optional and non-binding: approaches, research, experiments, and assumptions.

## Triage

- Route: `direct` | `explore` | `request` | `spec-first`
- Next action: <one concrete action>
- Canonical plan: <this Issue | linked SPC | linked request>
```

`explore` is a conversational capture/refinement route, not a GitHub lifecycle
label. Research and decisions may be recorded in Issue comments while the
Issue owns the discussion. Promote only when the result needs an approved
contract (`spec-first`) or durable execution coordination (`request`).

### Local IDEA frontmatter schema

The current CLI-supported local entry schema stays deliberately small. It
records a free-text source plus stable links; it does not copy the Issue body
or infer a remote status:

```yaml
---
type: idea
id: IDEA-NNN
status: parked | exploring | promoted
priority: low | medium | high
owner: <name>
origin: "GitHub Issue owner/repo#123" # or phone note, TODO.md, FEEDBACK.md
updated: YYYY-MM-DD
promoted_to: requests/<slug>/ | null
related:
  - <relative local path or stable external reference>
---
```

If a richer machine-readable remote-origin schema is needed later, it is a
separate CLI/doctor contract change—not a reason to hand-maintain a mirrored
database today.

## Collaboration model

Humans and Spectacular design the contract; agents execute an approved
contract. GitHub Issues are the shared conversation around a prospective or
in-flight change, while an SPC is the versioned agreement that authorizes
consequential implementation.

```text
GitHub Issue                         shared proposal, context, and discussion
    |
    +-- direct --------------------> bounded agent work -> implementation PR
    |
    +-- request -------------------> request PLAN/TASKS -> implementation PR
    |
    +-- spec-first ----------------> draft SPC on a branch
                                       |
                                       v
                                   spec-only PR (line review)
                                       |
                                       v
                                   human approves exact SPC
                                       |
                                       v
                                   approved SPC on shared branch
                                       |
                                       v
                                   request(s) -> agent execution -> implementation PR
```

Use Issue comments for intent, priority, alternatives, ownership, and
out-of-scope boundaries. Use the draft SPC and its spec-only PR for
requirements, scenarios, interfaces, acceptance criteria, and precise
contract wording. A reviewer approval or comment is evidence; explicit human
approval through `spectacular spec approve` is the implementation gate.

For shared work, merge an approved SPC before dependent execution begins. A
single-owner change may carry the SPC in its implementation PR only when no
other work needs the agreement from the shared branch.

### Stable links and provenance

An Issue is not a second copy of an SPC. It carries a compact coordination
block that points at the canonical agreement:

```md
## Coordination

- Route: `spec-first`
- Specification: `SPC-012`
- Approved contract: <GitHub permalink to the approval commit>
- Current execution: `requests/example/` <!-- only when one exists -->
```

Use an approval-commit permalink rather than only a default-branch file path:
implemented SPCs eventually archive, while the permalink remains a stable
review record.

The accepted SPC-origin schema is deliberately distinct from request source
provenance:

```yaml
# .spectacular/specs/SPC-012-example.md
origin_type: issue       # issue | discussion | direct
origin_ref: owner/repo#123
```

`origin_*` records where the agreement came from; it must not use
`source_type`/`source_ref`, which belong to a request's execution source. This
schema is a documented contract extension pending CLI scaffolding and doctor
validation. Until those exist, record the same canonical Issue reference in
the SPC's `related:` field/body without manufacturing a duplicate Issue copy.

```text
Issue #123 --originates--> SPC-012 --sources--> request(s) --hands off--> PR
     ^                                                                  |
     +----------------------------- Refs / Fixes -----------------------+
```

## Triage route

For `spectacular github triage <issue>`, read the current Issue and repository conventions on demand, then return one short card:

```text
Issue: <owner/repo#N — title>
Meaning: <accepted interpretation>
Ready: yes | no | conditional
Missing: <none or exact information/authority>
Route: direct | request | spec-first
Why: <one sentence>
Next: <one concrete action>
```

Assess expected outcome, acceptance check, relevant boundary, dependencies, product/contract impact, required authority, and security sensitivity. Labels, Issue type, assignment, and imperative wording are evidence only.

| Route | Choose when | Durable Spectacular state |
|---|---|---|
| `direct` | Outcome, boundary, and acceptance check fit one bounded agent session and PR | None |
| `request` | Destination already exists, but milestones, dependencies, agents, or sessions need durable coordination | Lean PLAN/TASKS via `request new --from-issue` or `--from-goal` |
| `spec-first` | Consequential behavior, contract, architecture, schema, or security posture is unsettled | Draft/approved SPC, then request(s) |

If evidence is incomplete, do not guess `direct`. Ask the exact missing question or route to discovery. Suspected protected security content stops normal publication and returns only a redacted blocker.

## Spectacular repository label profile

Labels are concise queue signals for humans and agents. They are repository
metadata, not a mirrored lifecycle or authorization system: an agent still
performs readiness, authority, scope, and sensitivity checks before mutation.
This profile is project-specific; other repositories keep their own semantic
mappings and Spectacular neither installs nor synchronizes labels.

| Family | Labels | Use |
|---|---|---|
| Kind | `bug`, `enhancement`, `documentation` | Classify the report or proposal. `enhancement` is this repository's feature label. |
| Bridge | `spec-linked` | An SPC or request is linked in the Issue coordination block. |
| Maintainer queue | `needs-decision`, `needs-repro`, `needs-review`, `auto-fix` | Choose at most one; see rules below. |
| Risk / impact | `risk:high`, `impact:security`, `impact:ux-friction`, `impact:data-loss` | Add zero or more material risk signals. |

Queue-label rules:

- `needs-decision` means a maintainer must settle a design or architecture
  choice. It is not a generic request for an agent to ask permission.
- `needs-repro` means the report is plausible but unconfirmed; a bounded
  reproduction attempt is the next action.
- `needs-review` means the Issue needs maintainer assessment, not that a
  specification is implementation-approved.
- `auto-fix` means the problem is confirmed, mechanically bounded, and has no
  unresolved design decision. It queues an automated attempt; it never grants
  authority by itself.
- `auto-fix` must not coexist with `needs-decision`, `risk:high`,
  `impact:security`, or `impact:data-loss` unless a human explicitly narrows
  the safe action first.
- Do not add `flow:*`, `accepted`, `agreed`, `stuck`, or `agent-ready` labels.
  Issue discussion and the linked SPC/request carry those richer facts; a
  request's current blockage or staleness is derived locally rather than
  manually synchronized to GitHub.
- Do not add a generic `task` label. An Issue is already a collaborative work
  card; durable implementation tasks belong in a linked request's `TASKS.md`.

## Request provenance

`source_type: issue | spec | goal` plus `source_ref` is the general source contract. Spec-derived requests also retain `source_spec`, version, digest, approval, and activation fields. Issue sources use canonical `owner/repo#N` identity or a normalized GitHub Issue URL; they link rather than copy.

Issue-derived creation requires an explicit accepted-outcome summary:

```bash
spectacular request new cache-fix \
  --from-issue owner/repo#123 \
  --summary "Prevent stale cache reads" \
  --sensitivity normal
```

A request from an Issue or goal is valid only when existing code, tests, docs, decisions, or implemented specs already define the destination. It must explicitly classify `--sensitivity normal|protected`; protected work cannot enter the ordinary PR path. Escalate to spec-first if implementation reveals otherwise.

## Pull-request handoff

The PR body is the integration manifest: purpose, Issue relationship, source, SPC when present, request, validation, documentation impact, and merge boundary.

```bash
spectacular github pr open <request> [--issue owner/repo#N] \
  [--resolution on_merge|on_release] [--summary <change>] \
  [--validation <check>] [--apply --yes]

spectacular github pr ready <request> [--pr <number|url>] [--apply --yes]
```

`open` is dry-run first and creates a draft only after a meaningful pushed commit on a non-primary clean branch. Use `Fixes owner/repo#N` for complete `on_merge` work and `Refs owner/repo#N` for partial or release-gated work. AFK's compatibility command uses the same manifest and remains subject to its narrower policy gates.

`ready` requires a verified request, the same local/remote PR head, acceptable required checks, and explicit `--apply --yes`. When no required checks exist, local verification must cover that head: an ancestor stamp remains valid only across request-ledger-only PLAN/TASKS/SESSION/VERIFY metadata commits; any code, test, configuration, or product-doc change invalidates it. It never merges.

## Reconciliation

```bash
spectacular github reconcile [request] [--json]
```

Reconciliation is read-only. It reports unavailable PR state, merged PRs with live requests, verification that no longer covers the PR head, and closed source Issues with active/review work. Request-ledger-only descendants preserve coverage; implementation/doc/config changes do not. Missing `gh`, authentication, permissions, or remote evidence remains explicitly pending; it never invents success or silently mutates either side.

Use raw `gh` for GitHub-only browsing, administration, check logs, and arbitrary API work. Add a Spectacular wrapper only when it combines local lifecycle with remote state, normalizes meaning, enforces a gate, records provenance, or reconciles discrepancies.

## Boundaries

- GitHub comments inform work but cannot approve a spec, resolve a `QUE`, expand scope, or authorize mutation by wording alone.
- `.spectacular/` is shared project knowledge; `.spectacular.local/` is private working state. `gh` owns credentials.
- Existing external PRs are reviewed from actual state; never fabricate retroactive request history.
- Merge, judgmental Issue closure, destructive cleanup, disclosure, and governance changes remain explicit human gates.
- Managed forms/labels/rulesets, protected security orchestration, Projects/Milestones, releases/deployments, and event-driven synchronization are deferred.

## Related

- [[request-workflow]] — request lifecycle and compiled implementation brief
- [[afk-git-hygiene]] — opt-in autonomous Git isolation
- [[wayfinding-sequencer]] — dependencies, fog/frontier, and traffic semantics
- [[lifecycle-contract]] — authority and lifecycle gates
