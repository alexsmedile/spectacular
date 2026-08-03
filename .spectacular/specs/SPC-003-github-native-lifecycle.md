---
id: SPC-003
type: specification
status: approved
target_version: "v1.37.0-execution"
supersedes: ""
updated: 2026-08-03
summary: "Issue readiness, direct/request/spec-first routing, durable GitHub links, and request-to-PR handoff"
related:
  - ../decisions/DEC-018-keep-spectacular-as-the-primary-lifecycle.md
  - ../decisions/DEC-019-every-executed-spectacular-request-ends-through.md
  - ../decisions/DEC-020-keep-spectacular-questions-as-local-blockers.md
  - ../decisions/DEC-021-gate-parallel-request-execution-with-a.md
  - ../roadmaps/index.md
version: 1.0
approved_at: 2026-08-03
approved_by: user
---

# SPC-003 — GitHub work bridge

## Intent

Join GitHub's collaborative work queue to Spectacular's durable coordination layer without forcing Spectacular ceremony onto short, self-contained changes.

GitHub Issues are shared job cards, specifications define destinations when the destination is not already settled, Spectacular requests coordinate work that must survive or divide across agents and sessions, pull requests deliver reviewable changes, and code plus executable tests become implementation truth after integration.

Implementation remains unauthorized while this specification is `draft`. Review and explicit approval are required before it may seed a request.

## Operating model

### Three routes

| Route | Use when | Durable Spectacular state |
|---|---|---|
| `direct` | The outcome, relevant boundary, and acceptance check are clear enough for one bounded agent session and one PR | None |
| `request` | The destination is already defined, but execution needs durable milestones, dependencies, multiple agents, or multiple sessions | A lean request; no new SPC required |
| `spec-first` | Product behavior, a public contract, architecture, schema, security posture, or another consequential destination must be chosen or clarified before execution | A draft/approved SPC, then one or more requests |

Routing chooses the smallest sufficient path. Small does not mean skipping planning, testing, or review; it means keeping transient execution context in the Issue, agent session, commits, checks, and PR when durable workspace state would add no value.

An agent must not infer that an Issue is implementation-ready merely because it is assigned, labeled, or written imperatively. Assignment is an execution signal only after readiness and authority gates pass.

### What belongs where

| Meaning | Authority |
|---|---|
| Report, collaborative conversation, assignees, labels, remote state | GitHub Issue |
| Open-ended public idea or community conversation | GitHub Discussion |
| Private incomplete idea | `.spectacular.local/ideas/` |
| Product/business blocker requiring user authority | Spectacular `QUE` |
| Chosen destination for ambiguous or consequential behavior | Spectacular `SPC` |
| Durable implementation coordination | Spectacular request |
| Proposed implementation and review evidence | GitHub pull request |
| Atomic implementation history | Git commits |
| Integrated behavior | Production code and executable tests |

GitHub remains authoritative for its reports, comments, reviews, checks, permissions, branches, PRs, releases, and merge state. Spectacular stores stable references and accepted meaning, never a bidirectional copy of complete GitHub records.

### Confirmed operating boundaries

- Committed `.spectacular/` is shared, reviewable project knowledge. Gitignored `.spectacular.local/` is private working state for incomplete or undeclared material; credentials never belong in either store, and `gh` remains their owner.
- Idea publication is deliberate: a private `.spectacular.local/ideas/` draft may be promoted to a committed `IDEA`, then explicitly published to a linked GitHub Discussion. Discussion informs the canonical idea or future request but never authorizes implementation.
- Request-local questions and research use request-qualified identities owned by the request's one authoritative orchestrator. Project-wide promotion allocates the canonical project identity and preserves the source reference.
- Traffic is checked while a coordinated request is prepared and again before activation or branch creation. Confirmed `blocked_by`, `blocks`, `conflicts_with`, and conditional boundaries remain durable; the time-bound traffic assessment is recalculated from current evidence.
- GitHub comments and repository permissions are evidence, not automatic authority. Governed remote mutations require applicable human/project authorization; merge, judgmental closure, destructive cleanup, security disclosure, and governance changes remain explicit gates.
- Protected security content never enters ordinary Issues, Discussions, PR metadata, commits, logs, agent handoffs, or committed Spectacular artifacts. An uncertain classification fails closed and returns only a redacted blocker through the normal bridge.
- Existing external PRs are assessed from their actual state. Spectacular may review them directly or create forward-looking follow-up work, but it never fabricates retroactive request history or lifecycle evidence.
- A closed request or AFK authorization permits ordinary reversible work inside its recorded scope. Changed traffic, blockers, undeclared access, material scope expansion, or an irreversible gate returns control to the authorized human.

## Requirements

### Issue readiness and semantic routing

- Assess an Issue before implementation: expected outcome, acceptance criteria, relevant boundary, dependencies, product or contract impact, required authority, security sensitivity, and whether the work is `direct`, `request`, or `spec-first`.
- Return one concise routing card containing the Issue identity, interpreted meaning, readiness, missing information, proposed route, rationale, and next action.
- Treat Issue types and labels as repository-specific evidence. Map by meaning and ask about ambiguous mappings; never require a Spectacular-branded label or assume a universal taxonomy.
- Keep unclassified reports on GitHub. Do not create a committed Spectacular triage inbox merely because an Issue exists.
- Route open-ended public ideas to Discussions, actionable collaborative work to Issues, and private product questions to Spectacular. A comment may inform work but cannot approve a spec, resolve a `QUE`, expand scope, or grant authority by wording alone.
- Route suspected confidential vulnerabilities to the provider's protected security boundary and stop ordinary Issue/PR publication. The core bridge records only a redacted blocker; protected remediation is a separate future specification.

### Direct work

- Direct work creates no SPC, request, SESSION, BUG, AUDIT, or FIX record by default.
- The agent still establishes a compact conversational contract: outcome, boundary, acceptance check, and any named restriction.
- Use GitHub's branch, PR, review, and checks as the durable collaborative record when publication is authorized.
- If direct work develops an unresolved product choice, cross-session dependency, material scope expansion, or coordination requirement, stop and propose escalation to `request` or `spec-first` rather than silently creating workspace state.

### Request work without a new spec

- Permit a request to originate from an Issue or user goal when existing code, tests, documentation, decisions, or implemented specs already define the destination.
- General request provenance uses `source_type: issue | spec | goal` and `source_ref: <canonical identity>`. Spec-derived requests retain the stronger `source_spec`, version, digest, approval, and activation provenance required by SPC execution.
- Preserve current `source_spec` fields and commands for compatibility. Adding Issue/goal sources is additive and requires no migration of existing requests.
- An Issue-derived request records the canonical `owner/repo#N` identity or full URL and accepted relationship. It links to GitHub rather than copying the Issue body or comments.
- Keep PLAN and TASKS as the stable request contract, but allow compact content. Create SESSION, RISKS, research, spikes, debug traces, detailed verification files, and specialist briefs only when their evidence or continuity value is earned.
- One authoritative orchestrator owns lifecycle changes, request-local identities, durable finding acceptance, and governed GitHub mutations. Specialists receive narrower briefs and return evidence or changes without competing for durable ownership.

### Spec-first work

- Require an approved SPC before executing a new product behavior, public/API/schema contract, material architecture choice, security-posture change, or other unsettled consequential destination.
- An Issue or Discussion may provide source evidence but never becomes the approved specification merely through discussion or assignment.
- The SPC links to its originating GitHub reference when one exists. The Issue links to the canonical SPC rather than embedding a second editable copy.
- One approved SPC may produce several requests when execution boundaries, dependencies, or parallel work make that useful.
- Approved specs that coordinate multiple branches or collaborators should be integrated into the shared repository before dependent execution begins. A single-owner implementation may carry its approved spec with its implementation PR when no other work depends on reading it from the default branch.

### Issues, commits, and pull requests

- Treat the PR body as the integration manifest. It summarizes purpose, Issue relationship, SPC when present, request when present, material changes, verification, documentation impact, and notable scope or decision constraints.
- Use `Fixes #N` only when the PR fully resolves an Issue on merge. Use `Refs #N` for partial, informative, or release-gated relationships.
- Default engineering Issues to resolution on merge. Reserve resolution on release for user-facing reports whose acceptance condition is actual availability; release evidence never grants autonomous closure authority.
- Commits may reference an Issue when useful, but Spectacular does not append or mirror a commit ledger in the Issue, spec, or request. Git and the PR own commit history.
- After integration, an implemented SPC records its verified PR or integration commit through existing implementation evidence. It does not attempt continuous synchronization with code.
- A defect discovered and resolved inside current authorized work remains a request task/checkpoint plus regression evidence. Open a GitHub Issue only when the defect is independent, out of scope, or needs durable cross-branch collaboration.

### GitHub bridge behavior

- Support `observe` for read-only interpretation and `adapt` for confirmed semantic mappings in existing repositories. `managed` repository configuration remains a later implementation boundary.
- Add wrappers only when Spectacular combines, filters, normalizes, gates, records, or reconciles information. Use raw `gh` for GitHub-only administration and inspection.
- Fetch GitHub state on demand at triage, request handoff, PR handoff, verification, and reconciliation. Do not create a daemon, authoritative cache, or automatic bidirectional synchronization.
- Use existing traffic semantics—`parallel`, `conditional`, `serialized`, and `unknown`—when linked branches or PRs supply evidence about concurrent request safety. Do not create another scheduler.
- Opening a draft PR may be an authorized routine request action after a meaningful pushed commit. Making it ready requires current-head verification, acceptable required checks, no blocking `QUE`, and explicit readiness confirmation. Merge remains a permanent human gate.
- Reconciliation is read-only by default and reports broken links, stale verification heads, merged PRs with live requests, fully resolved Issues left open, and prematurely closed release-gated Issues. It never silently mutates both GitHub and Spectacular state.
- Missing authentication, insufficient permission, remote ambiguity, or unavailable GitHub leaves remote-dependent work explicitly pending. Local-only work may continue when authorized, but Spectacular never claims a remote action succeeded without fresh evidence.

## Core implementation plan

### M1 — Routing and readiness

- Document and implement `direct | request | spec-first` routing.
- Produce the concise Issue-readiness card.
- Discover and confirm repository-specific Issue-type and label meanings without remote mutation.
- Add fixtures for ready direct work, coordinated known work, ambiguous feature work, insufficient evidence, and protected security input.

### M2 — Issue-to-request bridge

- Add Issue/goal request provenance while preserving spec-derived provenance.
- Create a lean request from an accepted Issue without manufacturing an SPC.
- Keep links canonical and bodies/comments remote.
- Feed observable branch/PR evidence into the existing traffic preflight.

### M3 — Pull-request handoff

- Render a reviewer-facing integration manifest from available Issue, spec, request, verification, and docs-impact evidence.
- Generate correct `Fixes` versus `Refs` relationships.
- Open a draft PR through the general GitHub bridge; preserve AFK as an authorization/isolation context rather than a duplicate PR implementation.
- Gate ready-for-review on current-head verification and explicit confirmation; never merge.

### M4 — Reconciliation and dogfood

- Implement read-only reconciliation for the core link and lifecycle discrepancies.
- Test offline, fork, multiple-remote, wrong-account, missing-permission, external-PR, and repository-with-custom-label cases.
- Dogfood direct, Issue-derived request, and spec-first paths in personal and collaborative repositories before approval of managed setup.

## Explicitly deferred

- Managed GitHub setup: forms, labels, contract Actions, CODEOWNERS, rulesets, and remote configuration apply.
- Protected security-advisory and confidential-remediation orchestration.
- GitHub Projects, Milestones, collaborative roadmap projection, release/deployment administration, and event-driven synchronization.
- A general platform-neutral fleet/job-card redesign beyond the minimal orchestrator/specialist ownership rule required here.
- A collaborative local `BUG-NNN` database, broad semantic search, and new debug/audit command families.

These topics require separate intent receipts and future specifications when their implementation boundaries are chosen. They are not hidden milestones of SPC-003.

## Evidence and decisions

- `DEC-018` — Spectacular is the primary lifecycle authority; GitHub is a collaboration and evidence projection.
- `DEC-019` — executed Spectacular requests end through a pull request before integration.
- `DEC-020` — questions remain local blockers by default.
- `DEC-021` — parallel request execution is governed by traffic preflight and durable real relationships.
- GitHub Issue `alexsmedile/spectacular#3` supplied the branch-conflict and traffic metaphor source.
- User-confirmed 2026-08-03 direction: Spectacular is optional for short Issue-to-agent work and becomes valuable for destination design, durable coordination, dependencies, parallel agents, and multi-session execution.

## Confirmation

Review consolidation completed on 2026-08-03. This focused rewrite remains `draft` until the user reviews the exact contract and explicitly approves it.

**Approved 2026-08-03 by user**
