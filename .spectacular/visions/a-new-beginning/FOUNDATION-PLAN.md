---
type: refactor-foundation-plan
status: active
authority: recovery-capsule
vision: a-new-beginning
checkpoint: pre-s10-sdlc-coherence-review
decision_sessions: 14
foundational_priorities: 20
updated: 2026-08-09
---

# Foundation Plan — Spectacular refactor

## Goal

Decide Spectacular's product identity, truth model, work model, authority boundaries, evidence loop,
ecosystem responsibilities, retrieval/scaffold, public interface, surviving subsystems, and target
implementation architecture before producing the final refactor plan.

This is a recovery capsule for the agreed decision order. Detailed questions, PZL inputs, and exit
gates live in [`evidence/decision-sessions.md`](evidence/decision-sessions.md); detailed ranking and
recommendations live in
[`evidence/top-20-foundational-decisions.md`](evidence/top-20-foundational-decisions.md); companion
analysis lives in [`evidence/responsibility-boundaries.md`](evidence/responsibility-boundaries.md).
Distributed work and checkpoint authority live in [`ORCHESTRATION.md`](ORCHESTRATION.md). No item
below is an accepted product decision until its session records the owner's disposition and this
orchestration task reconciles the reviewed packet.

S01, S03A, S02, S03B, S04, S05, S06, S07, S08, S09, and the compatibility floor are accepted in [`PRODUCT-CONSTITUTION.md`](PRODUCT-CONSTITUTION.md),
[`TRUTH-PROVENANCE-FLOOR.md`](TRUTH-PROVENANCE-FLOOR.md), and
[`SUCCESS-EVIDENCE-CONSTITUTION.md`](SUCCESS-EVIDENCE-CONSTITUTION.md), and
[`PRODUCT-TRUTH-CONTRACT-MODEL.md`](PRODUCT-TRUTH-CONTRACT-MODEL.md), and
[`WORK-UNIT-LIFECYCLE-CONTRACT.md`](WORK-UNIT-LIFECYCLE-CONTRACT.md), and
[`EXECUTION-AUTHORITY-CONTRACT.md`](EXECUTION-AUTHORITY-CONTRACT.md), and
[`EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md`](EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md), and
[`RESPONSIBILITY-PLACEMENT-CONTRACT.md`](RESPONSIBILITY-PLACEMENT-CONTRACT.md), and
[`RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md`](RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md), and
[`PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md`](PUBLIC-LANGUAGE-AND-INTERFACE-CONTRACT.md), and
[`CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md`](CLEAN-BREAK-CUTOVER-AND-RECOVERY-CONTRACT.md).
H16 is the read-only SDLC coherence gate before S10; later recommendations remain conditional.

## Decision batches

| Batch | Sessions | Decisions | Required result |
|---|---:|---|---|
| I · Product constitution | S01–S03 | Product identity, protected core, success function, truth and contract model | A stable constitution for judging every later option |
| II · Operating model | S04–S06 | Work units, lifecycles, execution/human authority, evidence and continuity | One coherent Mission-to-closure contract |
| III · Responsibility and surfaces | S07–S09 | Core/companion ownership, retrieval/scaffold, vocabulary and interface | Sharp boundaries with one authority per responsibility |
| IV · Reduction and delivery | S10–S12 | Subsystem survival, target implementation/migration, final spec program | Approved specs and an executable refactor sequence |

## Fourteen decisive sittings across twelve domains

1. **S01 — Product identity, protected core, and non-goals.** Define what makes the result
   Spectacular and which attractive responsibilities it refuses.
2. **S03A — Minimum truth/provenance floor.** Define enough authority, provenance, freshness, and
   missing-information semantics for evidence to be meaningful.
3. **S02 — Refactor success function and evidence constitution.** Define protected gates,
   qualitative regression review, survival evidence, prototype rules, and anti-gaming checks.
4. **S03B — Full truth hierarchy and contract model.** Define authority, freshness, provenance,
   Capability Contract shape, and the graph adoption level.
5. **S04 — Work-unit ontology and lifecycle.** Define Capability, SPC/spec, Mission/request,
   portfolio goal, objective, run, task, decision, discovery node, and record.
6. **S05 — Execution authority, human authority, and side effects.** Allocate execution,
   lifecycle, approval, Git/GitHub, and provider mutations.
7. **S06 — Verification, evidence, closure, and continuity.** Define proof, independent review,
   bounded repair, reconciliation, durable resume, and terminal next action.
8. **S07 — Core, companions, agents, modes, and adapters.** Decide pageworks, specwright,
   bugworks, verifyworks, Wayfinder, AFK, roles, namespaces, and handoffs.
9. **S08 — Retrieval architecture, instruction layers, and earned workspace.** Define universal
   context, projections, registries, reference tiers, minimal scaffold, and grow-on-write.
10. **S09 — Public language and interface grammar.** Define canonical vocabulary, CLI/skill
   grammar, readable artifact leads, status views, visual conventions, and deprecations.
11. **Compatibility-floor checkpoint.** Fix the supported population, deprecation promise,
    recovery boundary, and minimum compatibility window before any retirement.
12. **S10 — Subsystem, collection, policy, and fleet survival.** Apply S01's rubric to every
    contested capability; this is the first legitimate deletion session.
13. **S11 — Implementation architecture and migration strategy.** Define module seams,
    Bash/port posture, tests, migration mechanics, compatibility window, and recovery.
14. **S12A — Specification topology and approval.** Produce, review, dependency-check, and approve
    the fewest coherent specs.
15. **S12B — Executable refactor program.** Compile only approved specs into implementation
    Missions, dependency waves, joins, evidence gates, rollback, and stop checkpoints.

## Top 20 foundational priorities

Ranked by product blast radius, dependency fan-out, cost of reversal, and maintenance impact:

1. Product promise and boundary.
2. Protected behavioral loop.
3. Refactor success function.
4. Truth and authority hierarchy.
5. Contract information model.
6. Canonical work unit.
7. Mission, run, and lifecycle separation.
8. Execution authority.
9. Human authority model.
10. Side-effect and provider boundary.
11. Evidence and reconciliation contract.
12. Durable resume state.
13. Responsibility taxonomy.
14. Companion ecosystem and handoff.
15. Retrieval and instruction architecture.
16. Minimum workspace and growth rule.
17. Graph and multi-agent ceiling.
18. Public vocabulary and interface grammar.
19. Target implementation architecture.
20. Compatibility and migration strategy.

## Conditional responsibility hypothesis

The following remains a conditional responsibility hypothesis for S07. S01 has accepted the
control-plane product boundary, but detailed placement remains conditional until its named session.

- **Spectacular:** project control plane—durable context, accepted contracts, bounded work intent,
  decisions, lifecycle, evidence requirements, reconciliation, status, and resume.
- **Host coding runtime:** repository inspection, current-attempt planning, edits, tests, retries,
  and bounded delegation.
- **Spectacular CLI:** deterministic local projections, validation, scaffolding, lifecycle writes,
  and migration mechanics.
- **Native providers:** Git, GitHub, CI, deployment, permissions, and provider mutations.
- **Agent roles:** bounded workers with no independent lifecycle authority.
- **Companion skills:** optional complete standalone jobs connected through typed file/reference
  handoffs; never mandatory and never direct lifecycle mutators.

Companion starting position:

- retain pageworks as the proven optional companion;
- validate the decision multiplexer as the leading in-session companion experiment before turning
  it into a skill; route its artifact choice from impact, reversibility, uncertainty, and evidence
  threshold;
- treat AI UX stress testing as its first domain profile unless distinct substrate, tools, and users
  later justify an independent product;
- validate specwright as the first shared-engine extraction;
- validate bugworks as the strongest domain companion after Mission/evidence contracts settle;
- keep verification and closure in core for the MVP;
- defer Wayfinder extraction until discovery/Mission semantics settle;
- split AFK into core authorization/resume plus host-runtime and git-ops execution;
- classify most other candidates as roles, Mission profiles, validators, or adapters.

## Program rules

1. Constitutional order is S01 → S03A → S02 → S03B → S04–S06.
2. S07–S09 follow the accepted operating model and must remain mutually coherent.
3. No subsystem deletion before an accepted compatibility-floor checkpoint and S10.
4. No target code architecture before the final product surface is known.
5. No implementation task without an approved specification and validation path.
6. One implementation request per accepted boundary, not per decision sitting.
7. Graphs begin as relationships and derived views; scheduling/concurrency must earn promotion.
8. Spectacular remains standalone; companion integrations are optional.
9. Reopen upstream decisions only with named new evidence and an impact audit.
10. S12A may derive and approve specifications but cannot create implementation work; only S12B is
    authorized to produce the executable refactor program from approved spec versions.
11. This task is the program orchestrator: side sessions investigate, grill, spike, specify, build,
    or review; they return typed packets and never authorize successors or reconcile themselves.
12. Read-only evidence work may run in parallel when questions are independent. Type-1 decisions
    remain sequential. Concurrent mutations require disjoint scopes, explicit joins, separate
    worktrees, and an accepted traffic result.
13. Branches isolate mutations, not ideas, agents, or chat sessions. One implementation Mission
    owns one branch and reviewable PR; shared CLI/schema/registry surfaces are serialized.
14. Central checkpoint review records `accept | bounce | escalate | supersede`, reconciles only
    accepted results, and preserves one safe next action.

## Decision interaction protocol

Use the artifact that exposes the uncertainty most directly:

- classify impact, reversibility, uncertainty, and evidence threshold before choosing ceremony;
- optioned grilling for foundational dependent decisions;
- comparative matrices for several coupled but inspectable trade-offs;
- tension/contradiction mapping before generating options when goals oppose each other;
- two to four functional variants for UI/interaction alternatives;
- parametric sandboxes for continuous boundary conditions such as latency, scale, confidence,
  autonomy, and failure;
- simulated personas only as heuristic adversarial lenses followed by appropriate real evidence;
- executable CLI/simulators for state and recovery logic;
- deterministic checks and bounded independent review before expensive human inspection;
- direct execution for clear, reversible, evidence-gated work.

Decision density is a secondary S01 metric. Count only decisions that are understood, accepted,
retained, and not later reversed because the review surface hid relevant consequences.

Decision ownership follows an advice process: one authorized owner seeks affected and expert input,
records material dissent, and decides or escalates. Ownership never collapses separate product,
security, provider, irreversible-action, or policy authority.

## Immediate stabilization precondition

Before the large refactor, handle the verified `kind/type` compatibility-reader defect and the
implicit remote-branch deletion contract conflict as a narrow safety/correctness patch. This patch
must not prejudge the v2 product surface.

## Current checkpoint

S01, S03A, S02, S03B, S04, S05, S06, S07, S08, S09, and the compatibility floor are accepted and
hashed. Source 016 supplied no reversal-grade evidence, but H16 will stress the accepted model
against requirements, design, iterative delivery, deployment, maintenance, incidents, and
multi-actor coordination before S10 evaluates subsystem survival.
