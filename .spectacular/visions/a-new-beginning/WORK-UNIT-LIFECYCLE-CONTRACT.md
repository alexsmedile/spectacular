---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S04
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept-with-normalization
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
next_session: S05
---

# Work-Unit Ontology and Lifecycle Contract

This is the accepted S04 semantic model for proposed change, accountable work,
attempts, subordinate progress, decisions, gaps, handoffs, and durable history.
Its labels are internal placeholders; S09 owns final public vocabulary, IDs,
commands, aliases, and compatibility language.

It does not decide execution authority, evidence sufficiency, storage, product
packaging, public grammar, migration, or implementation.

## 1. One question per unit

| Semantic unit | Unique question it owns | It does not own |
|---|---|---|
| Capability Contract | What accepted observable behavior is currently required? | Delivery progress or implementation fact |
| Proposal | What behavior change is being considered or accepted as a target? | Current accepted behavior or execution progress |
| Mission/request | What one bounded approved outcome is accountable for delivery? | Current product behavior or attempt-local activity |
| Objective | Which durable outcome slice must this Mission achieve? | A second Mission lifecycle or executor actions |
| Run | What happened in one attributable, resumable attempt? | Mission intent or durable product behavior |
| Task/work step | What local action does this Run take next? | Durable cross-run progress authority |
| Decision | What was resolved or rejected, by whom, in what scope, and why? | An open question or execution progress |
| Gap | What consequential unknown, assumption, question, or discovery remains, and where does it matter? | A settled decision or separate activity taxonomy |
| Evidence/observation/history | What attributable fact, result, or past event was established? | Acceptance or lifecycle authority outside its own claim |
| Handoff | What bounded intra-Mission execution dispatch was issued, and what returned? | Mission progress or cross-Mission accountability transfer |
| Provisional coordination record | What immutable, attributable exchange must survive between independent Missions? | Either Mission's lifecycle or accountability |

No record gains authority merely because it links to another unit.

## 2. Ownership chain

```text
Proposal
  → accepted target contract delta
  → Mission/request
  → Objectives
  → Runs
  → run-local Tasks
```

A Proposal is not a Mission. Acceptance makes its delta an authorized target
for accountable work; it does not silently replace the current Capability
Contract. The current contract version remains authoritative until S06's
authorized reconciliation process accepts the result and updates or
supersedes it.

This preserves three distinct authorities:

- Capability Contract — current accepted behavior;
- Mission — accountable delivery of an accepted target delta or approved outcome;
- Run — attributable attempt history.

## 3. Canonical accountable work unit

A Mission is the sole durable accountable work unit. One Mission owns one
bounded approved outcome.

A Mission may:

- implement an accepted contract delta;
- repair a contract violation;
- preserve contracts while changing implementation;
- perform bounded research with an explicit zero-behavior delta.

Split work into separate Missions when outcomes can be independently accepted,
deferred, cancelled, evidenced, or handed off. Missions never nest.

A small, clear, reversible direct change may avoid durable Mission state, but
it does not bypass applicable intent, authority, evidence, reconciliation, or
safe-continuation gates.

## 4. Portfolio boundary

There is no portfolio Mission in the MVP. Goals may link independent Missions
but own no lifecycle or progress state.

Nested Missions and a separate portfolio entity remain deferred until linked
goals demonstrably fail cold recovery, accountability, or dependency
comprehension. Such failure is evidence for a new proposal, not automatic
promotion.

## 5. Mission, Objectives, Runs, and Tasks

- A Mission has one or more Objectives.
- Each Objective belongs to exactly one Mission.
- Objectives are durable, outcome-oriented progress slices with relevant
  dependencies and expected proof; they do not become nested Missions.
- A Mission has zero or more Runs.
- Every Run belongs to exactly one Mission.
- Tasks are local decomposition inside one Run and never become a second
  durable progress authority.

Durable Run identity exists only to support resume, retry, evidence
attribution, and handoff linkage. It is not a diary, schedule, queue, or
duplicate plan.

Record a Run only at meaningful boundaries: start, checkpoint, handoff or
return, interruption, or terminal result. The exact checkpoint payload and
evidence sufficiency remain S06 decisions.

## 6. Handoffs and cross-Mission traffic

A Handoff is an intra-Mission bounded execution dispatch to an actor, Run, or
session, with an optional attributable return. It is an immutable dispatch and
return record, not mutable work state and not progress authority.

Cross-Mission flow uses typed, pointer-first relationships:

- dependency;
- reference or contribution;
- decision or escalation;
- initiation.

When the exchange itself must survive independently, a provisional
cross-Mission coordination record may capture it. This candidate does not
transfer lifecycle or accountability between Missions. S07 decides product
placement; S08 decides whether it earns separate storage/schema; S09 decides
its public name.

No central mutable inbox, scheduler, executable graph, or handoff-as-progress
system is accepted.

## 7. Gaps, discovery, bugs, and decisions

One durable Gap concept has four kinds:

- unknown;
- assumption;
- question;
- discovery task.

A Gap earns independent identity only when it is consequential, cross-session,
attributable, cross-Mission, or requires bounded investigation. Otherwise it
remains a field or Objective within its owning Mission.

Research, spikes, and prototypes are methods for resolving Gaps, not automatic
parallel lifecycle systems. A bug is an observed Capability Contract violation.
A fix is normally a repair Mission. A finding is evidence or observation.

Decision records exist only for resolutions: settled-current, rejected, or
superseded. Open matters remain Proposals or Gaps. A materially new question or
choice creates a linked new record rather than silently reopening or
overwriting settled history.

## 8. Minimal semantic lifecycles

State exists only where it changes downstream behavior. Labels remain
provisional until S09.

```text
Proposal:
  draft → submitted → accepted | rejected | withdrawn

Capability Contract version:
  current → superseded

Mission:
  defined → active → awaiting-assessment → resolved
  resolved disposition: completed | abandoned | superseded

Objective:
  pending → ready → satisfied

Run:
  active → returned | interrupted | failed | cancelled

Gap:
  open → resolved | withdrawn

Decision:
  settled-current | rejected → superseded
```

`blocked` and `deferred` are orthogonal conditions, not multiplied lifecycle
states. S05 decides who may cause transitions and external effects. S06 decides
assessment, evidence sufficiency, closure, and reconciliation mechanics.

## 9. Cardinality and nesting rules

- One Capability Contract may receive many proposed deltas over time.
- One accepted target delta may seed one or more Missions only when each owns an
  independently accountable outcome; otherwise use one Mission.
- One Mission may affect multiple Capability Contracts while preserving one
  bounded accountable outcome.
- Missions do not nest; cross-Mission dependencies are typed links.
- One Mission owns one or more Objectives and zero or more Runs.
- One Run owns zero or more local Tasks and Handoffs.
- Evidence, decisions, and Gaps point to the narrowest owning scope without
  copying its progress or lifecycle.
- A Handoff stays within one Mission; durable peer exchange uses typed links or
  the provisional coordination record.

These relationships support dependency comprehension without defining a
scheduler or executable graph.

## 10. Cold-resume semantic requirement

A cold reader must be able to derive:

- the accountable Mission;
- its current Objective;
- the latest relevant Run checkpoint or return;
- blocking Gaps and cross-Mission dependencies;
- evidence and decision pointers;
- one safe next action or the exact owner gate preventing one.

This is an ontology requirement. S06 defines sufficient continuity evidence
and checkpoint semantics; S08 defines storage, retrieval, and projections.

## Exit condition

Every surviving unit owns one unique question. No two artifacts own the same
accepted behavior, proposed delta, accountable intent, attempt state, or
progress. Cold resume identifies one accountable unit and one safe next action
or explicit owner gate.
