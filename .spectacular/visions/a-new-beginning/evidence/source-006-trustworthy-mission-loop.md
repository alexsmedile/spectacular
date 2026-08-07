---
type: source-card
source: source-006
provided_as: source6
received: 2026-08-07
authority: proposal
status: ingested
scope: [product-loop, baseline, capability-contracts, missions, autonomy, evidence, resume, mvp]
duplicate_sections: []
completeness: substantial
---

# Source 006 — Trustworthy mission loop

## Thesis

Reduce Spectacular to one end-to-end loop: understand the repository, define an
end-to-end capability, approve a bounded contract delta, execute it autonomously,
prove it, reconcile documented truth, and resume without chat history. Defer every
entity, engine, and orchestration feature that does not strengthen that loop.

## Source integrity

This is a coherent replacement product model, not an incremental feature deck.
Its terminology overlaps Source 002, while its proposed minimal workspace and
mission lifecycle replace substantial current structure. It supplies no repository
audit or external evidence, so behavioral outcomes and artifact shapes are tracked
separately.

## Current-system comparison

| Proposed outcome | Existing analogue | Intake judgment |
|---|---|---|
| PROJECT/SYSTEM baseline | PRD, AGENTS, PRINCIPLES, ARCHITECTURE, STACK, specs index | consolidation candidate, not net-new outcome |
| capability contracts | `specs/` and archive/spec-sync | strong semantic overlap; path/schema replacement unproven |
| mission plan and lock | approved SPC plus request PLAN/TASKS and activation gates | lifecycle consolidation proposal |
| guarded mission run | skill build workflow, AFK authorization/runs, agent fleet, workspace preflight | existing pieces have different mutation authority |
| evidence-backed closure | PLAN validation, VERIFY/VERIFY-LOG, SPEC-DELTA, archive/spec-sync | protected outcome already central |
| status and resume | status brief, request brief, SESSION next actions | protected outcome already implemented in parts |

## Architectural collisions

1. The current PRD makes the skill the ongoing interface, calls the CLI a one-time
   bootstrap, and excludes multi-agent orchestration. Source 006 makes agentic
   `mission plan/run/close` shell commands the product center.
2. The current agent fleet makes the orchestrator the only general mutator. The
   proposed run gives one coding agent local edits, commits, and optional draft PRs.
3. The current mutation principle says the skill judges and the CLI performs
   deterministic reads/writes. Mission planning and closure require judgment.
4. `PROJECT.md` and `SYSTEM.md` could replace several current Anchors or duplicate
   them unless authority and migration are explicit.
5. `capabilities/` could clarify end-to-end semantics or duplicate/rename `specs/`.
6. One mission lifecycle would absorb current spec, request, discovery, fix,
   verification, session, AFK, and archive distinctions whose behaviors are not
   shown equivalent.
7. Relying on Git history conflicts with current snapshot and migration contracts.

## Proposed minimal deck

1. Repository-inspected, human-confirmed project/system baseline.
2. One compact contract per meaningful end-to-end capability.
3. A mission as an approved change to those contracts, expressed as outcomes.
4. One bounded, resumable coding-agent run using the host runtime.
5. Evidence mapped to contract clauses and closure-driven truth reconciliation.
6. Compact status and context reconstruction independent of chat history.

## Valuable ideas independent of the replacement shape

- Judge every feature by its contribution to one trustworthy closed loop.
- Keep contracts capability-centered until repeated complexity earns another type.
- Approve behavior and boundaries, not predicted file edits.
- Compile the smallest run context from accepted scope and relevant truth.
- Make stop conditions and retry budgets explicit before autonomy begins.
- Map evidence to exact promised behavior rather than recording generic test success.
- Treat resume by a cold agent as a primary acceptance criterion.
- Defer parallelism, schedulers, semantic retrieval, dashboards, and extra databases
  until the single-agent loop is excellent.

## Assumptions and contradictions to resolve

1. Autonomous execution is part of Spectacular's product rather than a convention
   operated by the host coding runtime.
2. One approval can preserve the current distinct spec, activation, irreversible,
   and closure authorities.
3. A capability is the right aggregation boundary for every supported project.
4. Mission objectives can replace persistent task and discovery records without
   losing useful evidence, ownership, or dependency semantics.
5. One active run and one coding agent satisfy the initial target users.
6. A draft PR is available and appropriate across project types.
7. Human-reviewed structured evidence is trustworthy enough without independent
   verification for the first release.
8. Git history alone supplies adequate document recovery and migration safety.

## Provisional assessment

**Strong protected outcomes:** capability-centered truth; outcome/boundary approval;
bounded context; explicit stops; clause-mapped evidence; closure reconciliation;
cold-agent resume; end-to-end acceptance testing.

**Promising but shape-dependent:** PROJECT/SYSTEM, capabilities/, missions/, one
MISSION/RUN/EVIDENCE bundle, and the compact lifecycle.

**Direct product decision:** whether Spectacular becomes an autonomous mission
runner. This explicitly reverses or supersedes current PRD and mutation-authority
boundaries; it cannot be adopted as a mere rename.

No replacement workspace, lifecycle, or execution authority is accepted here.

## New concept pieces

- [PZL-062 — One trustworthy closed loop](concepts/PZL-062-trustworthy-closed-loop.md)
- [PZL-063 — Repository-inspected project baseline](concepts/PZL-063-repository-inspected-baseline.md)
- [PZL-064 — End-to-end capability contract](concepts/PZL-064-end-to-end-capability-contract.md)
- [PZL-065 — Promote document types only when earned](concepts/PZL-065-earned-document-types.md)
- [PZL-066 — Mission as approved contract delta](concepts/PZL-066-mission-contract-delta.md)
- [PZL-067 — One mandatory execution approval](concepts/PZL-067-one-execution-approval.md)
- [PZL-068 — Outcome-oriented mission objectives](concepts/PZL-068-outcome-oriented-objectives.md)
- [PZL-069 — Reuse the host coding runtime](concepts/PZL-069-reuse-host-runtime.md)
- [PZL-070 — Compile bounded run context](concepts/PZL-070-bounded-run-context.md)
- [PZL-071 — Single-agent checkpointed run](concepts/PZL-071-single-agent-checkpointed-run.md)
- [PZL-072 — Explicit autonomous stop conditions](concepts/PZL-072-explicit-stop-conditions.md)
- [PZL-073 — Evidence mapped to contract clauses](concepts/PZL-073-clause-mapped-evidence.md)
- [PZL-074 — Closure reconciles current truth](concepts/PZL-074-closure-reconciles-truth.md)
- [PZL-075 — Cold-agent status and resume](concepts/PZL-075-cold-agent-resume.md)
- [PZL-076 — One compact mission lifecycle](concepts/PZL-076-compact-mission-lifecycle.md)
- [PZL-077 — Minimal mission workspace](concepts/PZL-077-minimal-mission-workspace.md)
- [PZL-078 — End-to-end MVP acceptance test](concepts/PZL-078-end-to-end-mvp-test.md)
- [PZL-079 — Explicit MVP deferral fence](concepts/PZL-079-mvp-deferral-fence.md)

## Decision packets seeded

- Is autonomous mission execution a product responsibility or a host-runtime recipe?
- Which existing artifacts should PROJECT/SYSTEM and capabilities replace, if any?
- Which approval authorities may safely collapse into one lock?
- Which lifecycle distinctions protect unique behavior versus incidental structure?
- What is the smallest real-project acceptance test for the protected loop?
