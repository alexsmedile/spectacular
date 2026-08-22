---
type: Proposal
id: 01a029be-b7d3-703c-a7ee-50c6b8bae3a2
ref: P11
title: Context-sandwich compilation and decision steering protocol
status: draft
created_by: Alex
created: "2026-08-22T15:51:00Z"
updated: "2026-08-22T22:01:33Z"
scope:
    - v2
target_contract: Contract:019fe381-5d61-7223-b362-03a5f99a7b10
---

# Context-sandwich compilation and decision steering protocol

Exploration for the M15-M21 campaign. This Proposal is mutable and grants no
execution authority. The owner-resolved boundaries below become binding only when
an approved Mission freezes them.

**Preparation verdict: `sufficient for M15`.** The outcome, responsibility
boundaries, command-surface change, proof threshold, and Mission slices are settled.
M15's activation gate still asks the owner to approve the exact recommended
tokenizer and Run transition matrix below. M16 and M18 cannot activate before that
approval is frozen. M16 must then prove the claimed context reduction before either
new command can ship in M17.

## Problem

Delegated agents are commonly given either too little context to act correctly or
so much repository context that cost and attention rise while scope discipline
falls. Owner choices also remain trapped in chat, so later workers re-ask settled
questions or act from incomplete memory.

Spectacular needs a compact, decision-aware charter for one Objective at a time.
It must improve context efficiency without becoming a Git wrapper, an autonomous
scheduler, or a second source of authority.

## Settled boundaries

### Progressive planning

- Only the next Mission is prepared as an activation-ready plan. Later Campaign
  blocks remain fluid summaries until upstream Evidence is accepted.
- Existing downstream drafts may be retained as design sketches. They are not kept
  synchronized, repeatedly validated, or independently reviewed before promotion.
- When a block becomes next, the Orchestrator re-prepares it from current Contracts,
  Decisions, and Evidence, reusing still-valid material rather than preserving stale
  precision.
- Owner questions are limited to open semantic choices and batched once. A final
  plan review is the default; another review requires a material semantic repair or
  newly discovered conflict.

### Authority and lifecycle

- Repository-changing background work requires an active Mission. Before
  activation, agents may research or build disposable prototypes, but they may
  not produce mergeable governed implementation.
- One Orchestrator session is bounded to one Mission. Several Missions may be
  active in separate integration worktrees when their declared order and active
  writable perimeters permit it.
- A Mission owns multiple Objectives. A Run is one execution attempt for one
  Objective; sequential Runs preserve retries and recovery without changing the
  Objective. Runs on distinct eligible Objectives may execute concurrently.
- The Orchestrator plans concurrency from dependencies, writable perimeters, and
  available capacity. There is no numeric concurrency ceiling. A Run never
  selects its own parallel work or spawns another agent.
- Agents are external execution destinations, not canonical governance objects.
  Spectacular preserves the Objective, Run, frozen Handoff, and returned Evidence.

### Git boundary

- Native Git commands create and switch branches and worktrees before
  `mission start`. The command only inspects Git state and refuses unless the tree
  is clean, the branch is non-default, and no merge, rebase, or similar operation
  is interrupted.
- Each concurrently active Mission uses an integration branch/worktree. Objective
  branches/worktrees start from that Mission branch and merge back into it.
- Spectacular records and verifies isolation; it never creates, stashes, checks
  out, merges, deletes, or otherwise wraps Git operations.
- After a Run or Objective is merged and tested, the Orchestrator may propose
  removing its clean worktree. The owner confirms every cleanup. Final branch and
  remaining worktree cleanup is a Mission-closure responsibility.
- Dirty, unmerged, or unverifiable state is never proposed for cleanup. Forced
  worktree removal is forbidden.

This boundary supersedes D15's automatic allocation and D16's automatic pruning
while preserving their safety outcomes: governed work does not execute directly
on `main`, and cleanup occurs only after integration is proved.

### Charter size and file perimeter

- The compiled governance envelope targets roughly 1,200 tokens. That budget
  covers frozen truth, selected Decisions, paths, stops, and proof commands. It
  does not include full source bodies or diagnostic output.
- From 1,201 through 1,400 tokens, compilation warns and returns judgment to the
  Orchestrator after attempting safe compaction. From 1,401 through 1,440 it emits
  a strong split recommendation and requires explicit Orchestrator approval.
  Above 1,440 tokens it refuses.
- A delegated slice normally edits two to four named files. This is a slicing
  default, not a universal limit. A larger coherent perimeter is allowed when the
  Orchestrator states why splitting would reduce correctness or coherence.
- Read-only sources and conditional manifest/lockfile allowances are listed
  separately from writable files.
- A Handoff carries a machine-readable `writes:` perimeter. Repository-relative
  file paths authorize exact files; a trailing `/` authorizes that directory
  subtree, including files created during scaffolding. Globs and `..` are
  forbidden. Renames require both source and destination to be in scope.
- The declared writable perimeter is hard. Editing outside it stops the Run and
  preserves the work for review. Expansion requires an Orchestrator-approved new
  Handoff because recorded Handoffs are frozen.
- Repair may add a bounded diagnostic allowance of roughly 3,500 tokens without
  expanding semantic scope or writable files.

### Decision selection and recording

- A quick live sweep of Decision frontmatter lets the Orchestrator select relevant
  choices by meaning. The charter compiler receives the selected refs and then
  behaves deterministically. No persistent preflight cache is required.
- The charter cites selected sources compactly, for example
  `Sources: [D12, D13]`. It never dumps the full Decision ledger.
- Durable Decisions cover architecture, security, data, dependencies, public
  interfaces, operational constraints, durable UX/API choices, and revisions of
  earlier Decisions. Routine implementation details and unchallenged hygiene
  assumptions do not enter the ledger.
- `spectacular decide` accepts a complete Decision package drafted by the
  Orchestrator from a file or stdin. It validates and records the package
  atomically, assigns identity and ref, and refreshes generated indexes.
- `decide` reports newly eligible work only when an Objective explicitly names
  the recorded Decision as its blocker. Otherwise the Orchestrator reevaluates
  eligibility. The command never infers eligibility by meaning or changes Run
  state.

### Proof and steering

- Every Mission records a baseline observation before activation. Fail-first is
  required only when the claimed change should make an existing behavior move
  from failing to passing. Performance, documentation, cleanup, and observational
  work record an appropriate non-failing baseline instead.
- `A`, `B`, and `C` select visible options. `M` combines explicitly named parts of
  visible options. Natural-language answers always remain valid.
- `G` asks the Orchestrator to prepare the next fidelity stage. `F` asks it to
  present final verification and any required authorization.
- Letters apply only to the visible card that defines them. `G` and `F` never
  silently authorize scope growth, Git mutation, deletion, provider effects, or
  spending.

### Live discovery boundary

A live frontmatter sweep provides the same correctness and retrieval capability
as D14's persistent cache without invalidation state or stale authority risk.
Caching adds no unique behavior, so P11 supersedes D14 and does not implement a
persistent cache. A later Mission may reconsider a disposable optimization only
after benchmarks demonstrate a material scale, latency, or token bottleneck.

## Compiled charter

The charter is a read-only retrieval helper that compiles a three-layer governance
envelope for one Objective. Mission and Objective frontmatter gain an explicit
`sources:` list of typed canonical refs. Compilation automatically retrieves the
bound Contract plus those directly declared `sources:` entries, in declaration
order with duplicates removed. The Orchestrator may add more explicit source refs
at invocation. The compiler never follows prose links, recursively walks a source's
own citations, or invents an uncited source by semantic inference. The charter
helps the Orchestrator write the complete Handoff; it is not itself the assignment
and is not stored as a canonical record.

```text
1. FROZEN TRUTH
   Project boundaries; active Mission and Objective; claim and proof boundary

2. OWNER STEERING
   AI-selected Decision refs and exact selected content; resolved Gaps; non-goals

3. EXECUTION PERIMETER
   Writable and read-only paths; allowed actions; stops; verification command
```

The compiler prunes only derived detail:

1. compress source summaries while retaining exact paths and required signatures;
2. shorten Decision rationale while retaining ref, disposition, and lineage;
3. condense repeated non-goal explanation;
4. never prune the Mission/Objective claim, proof boundary, authority, stops, or
   writable perimeter.

If safe compaction cannot meet the staged 1,200/1,400/1,440 boundaries,
compilation warns, requests judgment, or refuses as specified above. It never
emits an incomplete authority envelope. A broad but coherent Objective may split
into sequential Runs; independently provable outcomes split into Objectives.

**Recommended M15 owner choice:** use `spectacular-charter-tokenizer.v1`.
Byte-exact UTF-8 is counted
with a bundled `o200k_base` vocabulary and merge table identified by SHA-256 in the
generated interface and receipt. No Unicode normalization occurs; invalid UTF-8
refuses. Changing tokenizer data creates a new tokenizer version and re-runs every
boundary fixture instead of silently changing an existing limit.

## Dispatch and recovery

An asynchronous Run is eligible only when all of these are true:

1. its Mission is active;
2. its declared implementation dependencies are satisfied and any proof condition
   frozen in its Handoff has been accepted by the Orchestrator;
3. its bounded charter has compiled and a frozen Handoff names the Mission,
   Objective, Run, selected sources, and writable perimeter;
4. native Git isolation exists and its `writes:` perimeter is disjoint from every
   reserving Handoff across active Missions;
5. no open owner decision intersects its scope.

The Orchestrator may continue owner steering while eligible Runs work. It may
dispatch another Objective only when high-level planning already proves dependency
clearance and disjoint files. Missions may remain active despite potential future
overlap; only conflicting active Handoffs block dispatch.

Run states are `active`, `paused`, `blocked`, `awaiting-review`, `completed`, and
`stopped`. Active, paused, blocked, and awaiting-review Runs retain their writable
reservation. Completed and stopped Runs release it. Errors never select a state
automatically: the Orchestrator records the reason and next action, and asks the
owner when recovery changes scope, abandons work, or creates collision risk.

**Recommended M15 owner choice:** a new Run starts active. Legal transitions are:
active to paused, blocked,
awaiting-review, completed, or stopped; paused to active, blocked, or stopped;
blocked to active or stopped; and awaiting-review to active for rework, completed,
or stopped. Completed and stopped are terminal. Further work begins a new Run.

Evidence may cover one Run, one Objective, a cluster of Objectives, or the final
Mission. It names every Mission, Objective, Run, Handoff, and claim it answers.
It also records actor, executor, method, observed time, commit and tree, declared
checks and results, contrary evidence, and limitations.
The Handoff is never edited; a corrected or changed assignment is a new Handoff.
A completed Run means its execution ended successfully, not that proof has been
accepted. `objective finish` requires no reserving Run and at least one completed
Run; its claims may remain visibly unproven until their frozen Evidence or Review
gate. Small work does not earn a separate Evidence document merely for ceremony.

Normal `after:` dependencies wait for implemented work. Proof-sensitive sequencing
stays an explicit Orchestrator condition in the frozen Handoff: the Orchestrator
names the required Evidence or Review and does not dispatch until it is accepted.
P11 adds no second proof-dependency graph. If repeated Missions demonstrate that
this manual gate is unreliable, a later Proposal may justify mechanical proof edges
with real failure evidence.

## Mechanical interface

The campaign may grow the public surface from 14 to 18 commands in three reviewed
stages:

- `spectacular charter <mission-ref>/<objective-ref>`: read-only compilation and
  deterministic receipt for a selected Decision-ref set and declared perimeter.
- `spectacular decide <decision.md|->`: atomic durable Decision recording with
  owner wording, digest, targets, and optional supersession.
- `spectacular run transition <mission-ref>/<run-ref> --to <state> --reason
  <text> [--next-action <text>]`: atomic, attributable Run-state transition.
- `spectacular evidence record <mission-ref> <evidence.md|->`: atomic Evidence
  recording for one Run, one Objective, a cluster, or a Mission-final gate.

The existing `run start` signature changes from `<mission-ref>` to
`<mission-ref>/<objective-ref>` and refuses when that Objective already has a
reserving Run. No new canonical noun is introduced. Observable Contract and Run
model behavior changes through Contract version bumps, not amendments.

M17 may expose `charter` and `decide` only after M16's paired proof passes, growing
14 commands to 16. M18 may expose `run transition` and change `run start`, growing
16 to 17. M20 may expose `evidence record`, growing 17 to 18. M15, M16, M19, and
M21 add no public command.

Completed Missions remain byte-identical historical plans under the Contract and
Run model they activated against. New semantics apply only to a new Mission bound
to the new Contract version. A decoder defect is repaired rather than hidden by
rewriting history; owner-authorized archival may remove retired history from
routine orientation but never from recovery. While active, a Mission mutates only
operational state. Evidence that its frozen semantic plan is mistaken stops work
and returns an owner gate for a new activation version that preserves the prior
activation.

## Proof strategy

M14's immutable paired benchmark fixtures are the release gate for this campaign.

- M16 must demonstrate at least 40% lower total context ingestion than the
  full-scan baseline while preserving task success, safety, recovery, and decision
  fidelity.
- Measurements separate governance-envelope tokens, named source tokens, and
  repair diagnostics so the 1,200-token target cannot hide cost elsewhere.
- M17 proves atomic Decision and index recording before the first command growth.
- M18 proves Objective-scoped lifecycle and historical decoding before changing
  `run start` or exposing `run transition`.
- M19 proves refusal-before-effect for every dispatch gate and zero data loss for
  every stop path.
- M20 proves atomic clustered Evidence, complete attribution, fresh commit/tree
  coverage, declared checks, limitations, and no self-certification.
- M21 measures candidate file/scope guardrails and promotes only deterministic rules
  that reject real violations without rejecting coherent work. Zero new rules is a
  valid result. The concurrent-timeline Gap stays open until actual demand earns it.
- Every campaign Mission receives independent review because it changes public
  interfaces, authority routing, multi-worktree behavior, or benchmark policy.

## Campaign roadmap

```text
M15  Reconcile governance and Contract baselines
     Restore inherited proof gates and freeze the Contract transition map
  |
  v
M16  Build and benchmark the read-only charter engine
     Prove the >=40% paired gate; add no public command
  |
  v
M17  Expose charter and Decision recording
     Grow 14 to 16 only after M16 proof
  |
  v
M18  Introduce the Objective-scoped Run lifecycle
     Change `run start`; add `run transition`; grow 16 to 17
  |
  v
M19  Dispatch frozen Handoffs in native-Git isolation
     Enforce cross-Mission write reservations and owner-confirmed cleanup
  |
  v
M20  Record basic clustered Evidence
     Add `evidence record`; grow 17 to 18
  |
  v
M21  Measure and harden scope guardrails
     Promote only proven deterministic checks; finish the Campaign regression
```

P10's preparation-judgment checkpoint remains a separate Mission candidate. It
does not block M15, and P11 does not absorb it into the worker runtime.

## Non-goals

- No external scheduler, daemon, SaaS dependency, or hidden background service.
- No automatic Git branch creation, checkout, merge, push, or deletion by the
  Spectacular CLI.
- No concurrent reserving Runs on the same Objective or overlapping writable
  perimeter.
- No second proof-dependency graph or automatic proof scheduler.
- No generated concurrent timeline built merely to close an existing Gap.
- No general semantic quality or "anti-slop" classifier.
- No requirement to fully draft, synchronize, validate, or review every downstream
  Campaign Mission before the next Mission begins.
- No whole-workspace scan by a delegated Runner.
- No mandatory Decision record for trivial implementation details.
- No fixed file, line, or token proxy treated as proof of quality.
- No shortcut letter that carries unstated authority.
