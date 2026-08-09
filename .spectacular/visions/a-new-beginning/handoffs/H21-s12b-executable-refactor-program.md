---
type: decision-session-handoff
schema_version: spectacular.handoff.v2
handoff_id: H21
session: S12B
mode: read-only-interactive-owner-decision
authority: central-orchestration
status: authorized-for-dispatch
expected_parent_commit: 01a07c321952176af617ede38caa4c0d8d0c71d0
expected_parent_tree: 282774068f8f364202514eeb834be4c45970e85b
date: 2026-08-09
---

# H21 — S12B executable refactor program

## Mission

Compile the accepted five-spec topology into the smallest owner-approved executable refactor
program: draft implementation Missions, dependency waves, joins, authority envelopes, proof gates,
rollback/recovery points, and stop checkpoints. It must make a cold reader able to see what can be
implemented next and why. It must not activate a Mission, modify product code, migrate a workspace,
delete or archive legacy material, create a branch/PR, invoke a provider, or select Type-2
implementation details merely to make the program look complete.

This is a decision session, not implementation. It prepares the program that later authorizes
individual implementation Missions.

## Baseline gate

Before analysis, verify the detached clean checkout's commit and tree against the frontmatter;
verify this handoff and all fifteen accepted contracts against their sidecars; stop and report any
mismatch. Do not normalize baseline drift.

## Immutable inputs

Read all accepted foundation contracts, especially
`SPECIFICATION-TOPOLOGY-CONTRACT.md@1.0`, and verify their sidecars. The first fourteen immutable
input hashes are recorded in H19. The fifteenth is:

| Contract | SHA-256 |
|---|---|
| SPECIFICATION-TOPOLOGY-CONTRACT@1.0 | `d8eb174a5f6aea2356e60ba0acc7a9d3084df53a86e6b9c3694102af5c1d1ed1` |

Also read `AGENTS.md`, `.spectacular/AGENTS.md`, this handoff, `FOUNDATION-PLAN.md`,
`ORCHESTRATION.md`, `METHOD.md` Phase 7, the H19 return, and only the current source/tests needed
to estimate a proposed Mission boundary. Do not preload archives or unrelated evidence.

## Locked constraints

- Implement only through one bounded Mission at a time, with one branch/worktree only when mutation
  isolation is needed. A branch is not a Mission; agents do not get branches by default.
- Preserve the five-spec exclusive seams. Spec 1 defines semantic records; Spec 2 owns transitions,
  assessment, reconciliation, and continuity; Spec 3 owns guided/retrieval surfaces; Spec 4 owns
  deterministic Go/Markdown mechanics; Spec 5 owns cutover and recovery.
- The v2 core has no legacy parser, aliases, fallback reads, dual writes, lazy conversion, capsule
  import, or historical archive retrieval path. The capsule remains isolated.
- No generic scheduler, permanent agent fleet, provider capture, duplicate status system, or
  implementation DAG is introduced. Waves are program-planning relationships, not runtime control.
- The v1 mapping inventory is an explicit blocking Gap for migrated-candidate acceptance. No agent
  infers ambiguous conversion.
- Human gates cover consequential effects and Mission/Contract resolution, not routine already
  approved local edits. Native providers own their effects and receipts.
- Every Mission requires proportional preparation and its own authority, proof, recovery, and
  stop conditions. Retry/resume never widens authority.

## Interactive protocol — mandatory

Do not jump to a Mission list or return. For each cluster: inspect evidence; present 2–3 coherent
options with consequences; recommend one; ask one explicit owner decision; stop and wait for Alex;
record both verbatim answer and normalized disposition. Do not infer a decision from silence.

### A — Mission boundaries and vertical slices

Compare 2–3 programs that partition the five specifications into the fewest coherent, reviewable
implementation Missions. For each Mission candidate, state outcome, depends-on, non-goals, likely
authority/effect class, proof, recovery boundary, and why it is not too broad or fragmented.

Ask the owner to approve or change the Mission boundaries. Stop.

### B — Ordering, joins, and cutover gates

Turn the selected boundaries into ordered waves. Identify which Missions are serialized because
they share semantic/kernel/command/schema surfaces, which can be conditionally parallel after a
designed join, and which v1 release/capsule/mapping/candidate gates block later work. Include
rollback and a coherent-stop condition after each wave.

Ask the owner to approve or change the ordering and any allowed parallelism. Stop.

### C — Authority, evidence, and operating envelope

For each Mission, propose the minimum Mission charter: allowed effects, branch/worktree posture,
required deterministic checks, risk-triggered independent review, evidence receipts, owner gates,
retry budget class, and one cold-resume return. Identify where a prototype or design-sufficiency
review is required before code.

Ask the owner to approve or change the operating/evidence posture. Stop.

### D — Program skepticism and launch authorization

Run a fresh self-review of the complete program against the foundation: authority overlap,
impossible joins, provider capture, v1 leakage, lost recovery, under-specified migration, too-wide
Missions, false parallelism, unproven quality, and unclear owner gates. Present only material
findings and corrections.

Ask whether to approve the executable program, bounce named Missions/waves, or escalate a Type-1
conflict. Stop.

## Required return

Only after every cluster has an explicit owner disposition, return a complete
`spectacular.handoff-return.v2` packet with:

- H21/S12B baseline/tree, handoff hash, input verification, read set, and clean-state result;
- owner dispositions verbatim and normalized by cluster;
- every proposed Mission with outcome, dependencies, non-goals, authority/effects, proof,
  rollback/recovery, stop conditions, and approval state;
- ordered waves, conditional parallelism, designed joins, and explicit blocking Gaps;
- evidence/review/owner-gate matrix and cold-resume entry point;
- material skeptical findings and their resolution;
- Type-2 items deliberately deferred to Mission preparation or implementation; and
- `next_action: central accept|bounce|escalate; do not activate implementation`.

Do not write the program or return artifact in this session. Central orchestration alone accepts,
records, and authorizes an individual implementation Mission.
