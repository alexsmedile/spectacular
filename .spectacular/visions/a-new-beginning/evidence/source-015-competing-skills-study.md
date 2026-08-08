---
type: source-card
source: source-015
provided_as: H05-competing-skills-research-return
retrieved: 2026-08-08
authority: primary-repository-comparative-study
status: ingested
scope: [workflow-artifacts, gates, execution-runtime, verification, continuity, status-drift, companions]
completeness: three-target-study-complete
raw_location: returns/H05-competing-skills-study.md
---

# Source 015 — Competing-skills research study

## Study set

- [GSD Core](https://github.com/open-gsd/gsd-core) at `3f349e5`.
- [Superpowers](https://github.com/obra/superpowers) at `44c9b2d`.
- [sdd-skill](https://github.com/SpillwaveSolutions/sdd-skill) at `eba9606`.

The study used repository instructions, references, executable structure, tests, and manifests. It
compared mechanisms against Spectacular's proposed direction rather than ranking feature counts.

## Thesis

The strongest transferable behavior is not another universal SDD sequence. It is durable,
artifact-addressed continuity combined with explicit gates, bounded repair, independent evidence,
and a clean distinction between canonical project state and an execution attempt. These examples
support a small project control plane and warn against absorbing a model runtime, scheduler, status
dashboard, or expanding workflow taxonomy.

## Verified observations

1. GSD defines explicit pre-flight, bounded revision, escalation, and abort gates and carries phase
   decisions through persisted context and completion artifacts.
2. Superpowers uses fresh task workers, a plan-scoped scratch ledger, task review packages, capped
   repair rounds, and a final whole-branch review.
3. sdd-skill provides a teachable linear artifact sequence but relies heavily on prompt-directed
   state; its instruction and manifest versions differ, and its progress model shadows lifecycle
   through fixed percentages and artifact presence.
4. All three couple execution quality to host/runtime capabilities to different degrees. None
   demonstrates that Spectacular should own model hosting, general scheduling, or provider mutation.

## Transferable mechanisms

- Artifact-addressed handoffs with a terminal next action.
- Independent, fresh evidence before completion claims.
- Bounded repair with explicit escalation or abort.
- Brownfield inspection before durable planning.
- Run-local recovery state that references canonical intent instead of replacing it.

## Failure patterns to avoid

| Pattern | Cause | Spectacular guard |
|---|---|---|
| Prompt-enforced authority | Instructions are not a deterministic transition boundary | Validate closed authority/evidence properties before lifecycle mutation |
| Artifact-per-concern growth | Every recovery need becomes another file, marker, alias, or router | Require each artifact to own a unique query or transition |
| Shadow status | Percentages or summaries are stored separately from authoritative state | Derive projections from canonical artifacts and evidence |
| Runtime capture | Core coherence depends on one runtime's agents, worktrees, or interaction API | Keep execution mechanics behind a host adapter and portable handoff contract |

## Corrections and limits

- Repository mechanisms demonstrate existence, not user benefit or comparative effectiveness.
- A machine-checkable envelope can validate attribution, freshness, completeness, and declared
  results, but cannot mechanically establish every semantic claim.
- A separate execution/review companion is only a candidate; host-runtime behavior may be enough.
- GSD's artifact richness is both useful recovery evidence and a warning about taxonomy cost.
- sdd-skill's drift is one concrete example, not proof that every prose-driven workflow drifts.

## Foundation Plan impact

No decision order changes and no accepted contract is reversed. Source 015 strengthens S05's
host-runtime boundary, S06's evidence/closure/continuity contract, S07's companion-extraction test,
and S08's derived-projection rule. PZL-172 and PZL-173 become explicit S06 inputs.

## Concept routing

- New: PZL-172 — machine-checkable evidence reconciliation envelope.
- New: PZL-173 — run-local ledger subordinate to canonical Mission state.
- Reinforced: PZL-069, PZL-074, PZL-092, PZL-105, PZL-119, PZL-127, and PZL-129.
