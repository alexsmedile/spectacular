---
type: handoff-return
schema_version: spectacular.handoff-return.v2
handoff_id: H09
session: S04
status: accepted
central_disposition: accept-with-normalization
baseline_commit: cc292e4ff54bf8aa160f8ad94fa3dede97e468a4
baseline_tree: f46fc57467cdcd042d94f085b8f87d83a0897506
source_thread: 019fe210-a59d-72a1-ba8e-5d00357cb7a4
date: 2026-08-08
---

# H09 return — S04 work-unit ontology and lifecycle

## Central disposition

**Accept with normalization.** H09 validated every immutable input, obtained
explicit owner dispositions for all S04 branches, introduced no mutation, and
kept execution authority, evidence sufficiency, storage, and public language
outside its scope.

The accepted ownership chain is normalized so an accepted Proposal produces an
authorized **target delta**. It does not replace the current Capability Contract
until S06's authorized reconciliation gate. This preserves S03B's separation of
accepted behavior, accountable delivery, and implementation fact.

## Accepted decisions

1. Proposal, current Capability Contract, Mission, Objective, Run, and Task own distinct questions.
2. Mission is the sole durable accountable work unit; Missions never nest.
3. Goals link Missions without owning progress; no portfolio Mission exists in the MVP.
4. Objectives are durable Mission progress; Tasks remain Run-local decomposition.
5. Runs are minimal boundary records for resume, retry, attribution, and handoff linkage.
6. Handoffs are immutable intra-Mission dispatch/return records, never progress authority.
7. Cross-Mission traffic uses typed links; a coordination record remains provisional and conditional.
8. One typed Gap replaces separate research/question/spike lifecycle families.
9. Bugs are contract violations, fixes are normally repair Missions, and findings are evidence.
10. Decisions record settled-current, rejected, and superseded resolutions only.
11. Compact state machines use orthogonal blocked/deferred conditions.
12. Direct changes may avoid Mission ceremony but not applicable protected gates.
13. Cold resume identifies the accountable Mission, current Objective, relevant Run, Gaps, evidence, and one safe next action or owner gate.

## Boundary check

S05 retains transition and effect authority. S06 retains evidence sufficiency,
assessment, checkpoint payloads, closure, and reconciliation. S07 retains
product placement, S08 storage/retrieval/schema, and S09 final public language.

## Result

The reconciled contract is
[`../../WORK-UNIT-LIFECYCLE-CONTRACT.md`](../../WORK-UNIT-LIFECYCLE-CONTRACT.md).
S05 becomes next-ready after validation and commit; it is not dispatched or
independently authorized by H09.
