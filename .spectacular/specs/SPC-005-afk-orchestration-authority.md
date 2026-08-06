---
id: SPC-005
type: specification
status: approved
target_version: "tbd"
supersedes: ""
updated: 2026-08-06
summary: "Bounded AFK coordination convention with auditable delegation and protected Git boundaries"
related:
  - ../ideas/IDEA-002-afk-orchestration-authority.md
  - https://github.com/alexsmedile/spectacular/issues/25
  - ../PRD.md
  - ../PRINCIPLES.md
  - ../POLICY.md
version: 1.0
approved_at: 2026-08-06
approved_by: alex
---

# SPC-005 — Bounded AFK coordination convention with auditable delegation and protected Git boundaries

## Intent

Define the smallest coordination convention that lets a human authorize one
AFK run to complete a declared goal graph inside a run-scoped integration
branch, while keeping exercised authority reconstructable and preserving human
control over material scope, durable work, and protected destinations.

This is not a general multi-agent orchestration platform. It is a bounded,
Markdown-first authorization and evidence contract for one orchestrator, its
closed-scope delegates, and one short-lived integration branch. It is designed
to remain compatible with the PRD's non-goal while making the existing AFK
Git-isolation feature more explicit and auditable.

This draft creates no implementation authorization. Its first request, if
approved, must implement only validation and read-only inspection of the
proposed record; it must not add Git-mutating authority.

## Requirements

### R1 — Per-run authorization and event record

- An AFK run uses one durable Markdown record at
  `.spectacular/afk/runs/<run-id>.md`.
- The record has a write-once authorization block for CLI mutation paths, with
  the human identity and timestamp; session; declared goal; initial goal nodes
  and dependencies;
  known work references; human-approved base and
  `afk/<run-id>/integration` integration branch; allowed and excluded actions;
  and declared human-in-the-loop gates.
- The record has an append-only-by-CLI event log. Every event contains a stable
  event ID, timestamp, actor, scope, rationale, and evidence reference. Event
  kinds are: `authorization`, `node-added`, `technical-choice`, `delegated`,
  `returned`, `verification`, `acceptance`, `integration`, `git-action`,
  `gate`, `resume`, `handoff`, `completion`, `cancellation`, and `cleanup`.
- A record documents claims and evidence. It must never claim that a diff,
  node, or commit grouping is semantically certain merely from paths, diffs, or
  an agent report.
- A correction preserves the original authorization value, reason, actor, and
  timestamp as an event; no authorization boundary silently broadens.
- Markdown remains human-editable. Normal Git review/history is the audit
  boundary, not tamper-proof storage. Doctor detects malformed fields,
  duplicate/non-monotonic event IDs, invalid chronology, and missing required
  evidence; it cannot prove that a human did not edit a record.

### R2 — Goal graph, prerequisite nodes, and soft decisions

- The declared goal graph is the authority boundary. Each node names a bounded
  outcome, dependencies, and acceptance evidence.
- The orchestrator may add a prerequisite node only to complete an authorized
  node. Before delegation, it records the parent/dependency, rationale,
  bounded outcome, and verification route in `node-added`.
- New durable requests, specifications, ideas, research, spikes, and product
  scope remain human-gated. The agent must gate rather than create them under
  the AFK delegation.
- A soft technical choice may proceed when it preserves the declared goal,
  accepted specifications, interfaces, and constraints; is reversible; and
  creates no material compatibility, security, cost, or operational commitment.
  It is recorded as `technical-choice` with alternatives, rationale, and
  evidence.
- A durable technical decision requires explicit `technical-decisions` in the
  authorization block and the existing narrow AFK decision rule: active run,
  technical and in scope, reversible, evidence-backed, alternatives recorded,
  and no product/business trade-off. Otherwise it gates.

### R3 — Branch containment and Git boundary

- An authorized new-format run has one short-lived integration branch,
  `afk/<run-id>/integration`, from a named human-approved base. The branch
  contains work; it does not grant authority on its own.
- A delegated task branch is a sibling ref in that run namespace, for example
  `afk/<run-id>/node-<slug>`, and may integrate only into its own integration
  branch after recorded verification and orchestrator acceptance.
- Commit, push, PR creation, and merge to the integration branch are possible
  only when separately present in the write-once `allowed_actions` block and
  only after a later, explicitly approved implementation phase adds those CLI
  operations.
- Default/protected, release, deployment, and configured protected branches;
  remote deletion; releases; and deployments remain excluded. No AFK command
  may merge to those destinations automatically, even if the Git provider
  would permit it.

### R4 — Delegation, verification, and acceptance

- Every delegate receives a closed brief: node ID, bounded outcome, permitted
  branch/actions, exclusions, dependencies, acceptance criteria, and return
  condition. It inherits no broader authority.
- Delegated authority ends when the delegate returns, completes, is cancelled,
  or reaches a gate. A delegate cannot create durable work, add scope, accept
  another result, or act outside its run namespace.
- Completion reported by a delegate is not integration eligibility. A
  `verification` event must name the verifier, independent scope,
  commands/checks, result, evidence, and limitations.
- The main orchestrator independently records acceptance or rejection. Only an
  accepted result may be integrated; integration records both heads and the
  resulting SHA once the future mutation phase exists.

### R5 — Gates, failure, and recovery

- The run transitions `active → gated` before an unexpected product/business
  decision; external account, browser, or sensitive-data access; undeclared
  scope/dependency; protected destination; breaking API/schema decision; new
  durable work boundary; missing verification; or out-of-scope conflict.
- A gate records the blocked action, reason, safe branch/HEAD, options, and the
  exact human decision needed. Resume records that decision as a new event; it
  does not silently broaden prior authority.
- Failed verification causes rejection or an in-scope repair node, never
  integration. Cancellation preserves the record and all branch/commit
  evidence.
- Cleanup keeps the existing archive-first recovery property: a verified local
  recovery ref and restore command precede local branch deletion. Remote
  deletion remains human-only.

### R6 — Completion, session end, and non-goals

- A run completes only when all declared or validly added nodes are accepted or
  explicitly not-needed, integrations and handoff evidence are recorded, and
  no gate is open. Completion ends inherited delegate authority.
- The final handoff records integration-branch tip, commits, PR URL when one
  exists, validation evidence, outstanding human actions, and cleanup state.
- Issue #5's `spectacular session end` commit review remains independent and
  read-only. It does not inspect authority as permission and never stages,
  commits, amends, pushes, merges, resets, stashes, or records a Git mutation.
- This specification does not introduce a scheduler, a universal agent fleet,
  autonomous product planning, automatic protected-branch merge, remote branch
  deletion, release/deployment authority, or a replacement for Git/GitHub
  review controls.

### R7 — Compatibility and migration

- Existing AFK run files remain readable with their current lifecycle
  (`active | gated | completed | cancelled`). Missing new-format fields mean
  legacy, narrower authority only.
- Migration is preview-first and explicit. A future `afk run authorize` dry run
  renders the complete record before a human confirms it is written.
- Current project opt-in, explicit `--apply --yes`, non-primary refusal,
  archive-first cleanup, and draft-PR behavior remain compatible.
- New Git actions must bind to an active new-format run, its integration
  branch, write-once `allowed_actions`, and recorded event evidence. They must
  fail closed when any is missing.
- Implementation remains Bash 3.2-compatible and adds no dependencies.

## Phased implementation boundary

### M1 — Record and validate, no Git mutation

- Define the new-format AFK run schema, append-only-by-CLI event conventions,
  write-once authorization behavior, and doctor validation.
- Add read-only status/inspection output that makes a run's authority, event
  chronology, gates, and missing evidence visible.
- Add no CLI command that stages, commits, amends, pushes, opens a PR, merges,
  resets, stashes, deletes branches, or changes remote state.
- Test malformed records, event ordering, missing evidence, legacy-run
  compatibility, gated/resumed/completed/cancelled states, no repository, and
  #5 session-end independence.

### Later work — separately approved after M1 evidence

Any mutation-capable surface must be proposed in a later request after M1 is
verified. It must name exact command grammar, `--apply --yes` behavior,
per-action active-run checks, integration-branch-only enforcement, recovery,
and offline/permission failure tests. Approval of this draft does not approve
that later work.

## Evidence and decisions

- GitHub Issue [#25](https://github.com/alexsmedile/spectacular/issues/25)
  identifies the need to gate AFK commit/PR authority on a known goal and
  explicit approval.
- IDEA-002 contains the imagine-and-grill record: orchestration-level,
  one-time delegation; run-scoped `afk/<run-id>/integration` branch and sibling
  task branches; bounded prerequisite nodes; delegate inheritance; independent
  verification; and human-gated durable/product/protected boundaries.
- `afk-git-hygiene.md` currently provides the narrower three-layer AFK model,
  dry-run-first branch isolation, archive-first cleanup, and draft-PR handoff.
- `lifecycle-contract.md` supplies the existing AFK lifecycle and narrow,
  evidence-backed autonomous technical-decision exception.
- The PRD's multi-agent-platform non-goal constrains this to a minimal
  coordination convention rather than a general orchestration product.
- Issue #5's completed session-end review establishes the independent,
  no-mutation working-tree review boundary.

## Confirmation

Drafted after the human confirmed the minimal-coordination framing and the
authorization/event record model. It remains in review: explicit approval is
required before it can seed even the M1 no-mutation request, and later
Git-mutating authority needs a separate approved request after M1 evidence.

**Approved 2026-08-06 by alex** — Approved M1 record-validation and read-only-inspection scope; Git-mutating authority remains deferred
