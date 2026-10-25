---
type: Plan
version: '0.1'
created: '2026-10-25T11:44:34Z'
updated: '2026-10-25T11:44:34Z'
description: Evolve Spectacular into a guided software companion with anticipatory reasoning, explicit scope, and durable learning.
---

# Software companion upgrade

## Outcome and authority

The owner wants a method delivered through a skill: a companion for beginners,
hobbyists, and engineers to shape, build, operate, and potentially sell software.
It anticipates needs, helps choose boundaries, teaches through actual decisions,
and preserves reasoning across sessions. This plan records the requested direction;
delivery phases and defaults below are proposed, not shipped behavior.
Planning does not activate a Mission or authorize implementation/publication.

First usable outcome: a person shapes a small app, selects a bounded first version,
builds and exercises one useful workflow, then resumes in a fresh session without
losing boundaries, rationale, or the next action. Deployment is a later exposure
milestone when the selected pilot requires it, not a prerequisite for every project.

## Baseline and affected owners

- [Product](../PRODUCT.md), [project](../PROJECT.md), and [roadmap](../ROADMAP.md)
  currently foreground durable context and optional governance; preserve these foundations.
- [Skill](../../skills/spectacular/SKILL.md) already supports ordinary work, but its
  main body is dominated by governed procedures. Make the companion method the entry.
- [Orient](../../skills/spectacular/references/orient.md) and
  [prepare](../../skills/spectacular/references/prepare.md) assume CLI initialization
  and Mission-oriented outputs; reconcile these with ordinary work and explicit opt-in.
- [Requirements](../../skills/spectacular/references/requirements.md),
  [planning](../../skills/spectacular/references/plan.md), and
  [readiness](../../skills/spectacular/references/milestone-readiness.md) already
  supply progressive questions, prioritization, usable increments, and evidence distinctions.
- [Knowledge](../../skills/spectacular/references/knowledge.md) already preserves
  rationale, ownership, provenance, and conflicts. Extend it rather than adding a memory store.
- [Focused evaluations](../../test/benchmarks/focused-evals.json) cover routing and
  engineering tasks; companion behavior needs multi-turn and cold-resume coverage.

## Proposed delivery sequence

### 1. Define the method and reconcile entry routes

Update product/project direction and summarize this plan in the existing roadmap.
Define the reusable loop: understand → surface implications → choose boundaries →
build → verify → learn → reconcile context. Permit entry at any point for existing projects.
Refactor the Skill into a short router; put governed instructions behind explicit selection.
Patch orient/prepare and reconcile conflicting direct-build versus discovery instructions:
clear authorized tasks execute directly; ambiguous product choices get focused guidance.
Acceptance: idea discussion, direct implementation, existing-project resumption, and
explicit governed work each select the proper route without mandatory initialization.

### 2. Add anticipatory discovery and scope reasoning

Introduce one narrowly loaded companion reference for interaction and implication discovery;
extend requirements/readiness for scope rather than duplicating their definitions.
Capture users, problem, success, constraints, available tools, and consequential unknowns.
Suggest missing needs with a reason, consequence of omission, and recommendation.
Separate user requests, agent hypotheses, accepted choices, and verified observations.
Track Must now, Nice to have, Later, Dreams, and Out of scope in the existing content owner.
Later carries a revisit condition; Out of scope states whether the boundary is for this
increment or the product; Dreams carry no delivery commitment. Must explains necessity.
Flag contradictions and negotiate tradeoffs; never silently promote suggestions into scope.
Acceptance: a booking-app example surfaces cancellation/time-zone/collision questions as
applicable, preserves a deferred payments dream, and builds only the selected workflow.

### 3. Make guidance adaptive and memory useful

Use demonstrated expertise and user preferences to adjust explanation depth, not rigor.
Recommend reversible defaults; ask only questions that materially change the next outcome.
Teach a short principle at the relevant decision, with optional detail and no compulsory lesson.
Maintain rationale, rejected alternatives, uncertainties, evidence, and scope changes in
existing owners; reconcile affected summaries after material changes, not every message.
On return, retrieve relevant context and explain current outcome, changes, and next action;
surface conflicts rather than inventing a resolution. Preserve unknown metadata and provenance.
Acceptance: novice and engineer versions reach equivalent decisions without patronizing
the engineer; a fresh session recalls a rejected option and does not reintroduce it casually.

### 4. Pilot the complete workflow and measure behavior

Run a small internal-tool scenario with a beginner and an existing-project change with an
engineer, including an interruption, scope change, and failed assumption. Add a payment or
real-data exposure scenario to test readiness reasoning without performing live effects.
Use the existing scenario/evaluation infrastructure where suitable; inspect runner support
before promising automated multi-turn coverage. Preserve transcripts and artifact evidence.
Check implication quality, scope stability, question relevance, recovery, usability, and
truthful completion claims. Compare against current behavior using equivalent tasks.
Acceptance: authorized workflows work, deferred ideas stay deferred, cold resumption retains
decisions, and no unsupported readiness or unauthorized effect occurs. Report model/host
and limitations; human walkthrough evidence is needed for beginner-usability claims.

### 5. Package the proven experience and extend readiness

Update README, quickstart, installation, process, and architecture docs to match verified
behavior. Explain the skill, host agent, optional CLI, source access, and user-owned project.
Inspect actual license terms before promising commercial rights or redistribution freedoms.
Extend existing readiness guidance proportionally for deployment, costs, recovery, support,
and commercial use. Selling an app remains separate from proving it works or has demand.
Acceptance: a fresh user can install, begin, locate their project, and resume with the
documented capabilities; publication and commercial claims are backed by relevant evidence.

## Boundaries and priorities

Must for the first pilot: phases 1–3, focused behavior checks, and one exercised workflow.
Nice to have: reusable examples, richer progress views, and optional deeper explanations.
Later: broader domain playbooks and host/provider integrations after pilot evidence identifies needs.
Dreams: marketplace/distribution assistance and broader commercial launch support; no commitments.
Out of scope for this upgrade: hosted platform, compulsory dashboard/account, new database,
universal app stack, mandatory governance, historical record migration, or new public CLI commands.
Keep the current 26-command interface unchanged; any command change needs a separate owner proposal.

## Verification and continuation

Preserve the existing dirty working tree, including Skill, reference, ontology, and changelog edits.
Use focused skill checks, link checks, and behavior evaluations for runtime guidance changes.
Do not run verify.sh for this Markdown plan. If later work changes Go, use quick; before heavy
acceptance/release/all tiers run and pass preflight; all remains the final release/Mission gate.
Reconcile roadmap and product anchors as implementation lands, preserving frozen governed records.
Open choices: initial pilot app/host, guidance defaults, and supported deployment path;
choose these before their dependent work, not as a prerequisite to defining the method.
Next implementation batch: phase 1 plus one phase-2 discovery/scope example, verified against
ordinary and governed routing. Record outcomes here; do not create a second delivery roadmap.
