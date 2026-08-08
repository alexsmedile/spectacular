---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S09
accepted_by: owner
accepted_at: 2026-08-09
central_disposition: accept-with-normalization
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
  - EXECUTION-AUTHORITY-CONTRACT.md@1.0
  - EVIDENCE-CLOSURE-CONTINUITY-CONTRACT.md@1.0
  - RESPONSIBILITY-PLACEMENT-CONTRACT.md@1.0
  - RETRIEVAL-AND-EARNED-WORKSPACE-CONTRACT.md@1.0
next_session: compatibility-floor
---

# Public Language and Interface Contract

This accepted S09 contract chooses public semantic language and interaction
grammar. It does not authorize implementation, storage/schema choices,
compatibility removal, migration, or subsystem retirement.

## Product posture

Spectacular uses a constrained Mission vocabulary over a neutral semantic
model. `Mission` is the sole metaphorical canonical work noun. `Anchor` names a
scope-entry retrieval role. Contracts, authority, evidence, decisions, release,
reconciliation, and closure remain neutral.

Mission Control, Atlas, Journey, Waypoint, Telemetry, Signal, Beacon, Letter,
and universal Launch/Operations lifecycle metaphors are not canonical language.

## Canonical glossary

| Term | Unique public meaning |
|---|---|
| Capability Contract / Contract | Current owner-accepted observable behavior |
| Proposal | Candidate behavior change or bounded target delta; never current Contract truth by acceptance alone |
| Mission | One bounded approved outcome accountable for delivery; Missions never nest |
| Objective | Durable outcome-oriented result within a Mission |
| Run | One attributable resumable attempt within a Mission |
| Task | One local action within a Run |
| Checkpoint | Meaningful Run boundary: interruption, handoff, return, or terminal result |
| Gap | Consequential unresolved unknown, assumption, question, or discovery need |
| Question | A Gap kind when durable; otherwise an ephemeral utterance |
| Fog | Non-authoritative projection of Gaps, dependencies, conflicts, and authority blockers |
| Evidence | Attributable result or observation; never acceptance itself |
| Decision | Attributable resolution or agreement |
| Handoff | Bounded intra-Mission execution-responsibility dispatch and optional return; Mission accountability does not transfer |
| Mission Link / Link | Static typed relationship between independently accountable Missions |
| Mission Message / Message | Immutable attributable communication across a Link; emission cannot mutate the receiver |
| Anchor | Authoritative minimal context loaded when entering project or Mission scope; not a work-unit entity |
| Plan | Temporary host/session explanation within a Run; not durable Spectacular progress authority |
| ADR | Familiar technical format/name for a Decision record |
| Issue | Provider-owned conversation/work record; Spectacular retains pointers but not provider authority |

Link kinds are `depends-on`, `contributes-to`, `escalates-to`, `initiates`, and
`references`. Message kinds are `request`, `contribution`, `notice`,
`escalation`, and `response`; delivery may be direct, multicast, or broadcast.
A Message earns durable identity only under S08's durable-record rule. Its
content hash may be its identity, but hash algorithm and representation remain
implementation decisions.

`Spec` is not the target canonical public noun: use Contract for current
accepted behavior and Proposal for a candidate change. This is a target-language
decision only; the compatibility floor decides supported aliases and removal.

## Identity grammar

Stable identity, readable slug, accepted version, and content fingerprint are
separate. Target public shapes are:

```text
CC:<stable-key>@<version> <slug>
P:<stable-key> <slug>
M4/O1
M4/R2
M4/R2/T3
M4/R2/C1
M4/H1
L7: M4/O2 depends-on M9/O1
msg:<content-hash>
```

A Proposal names its target Contract, base version, and base fingerprint;
reconciliation stops on mismatch. A hash binds an evolving record revision but
is not its durable identity. Exact key encoding and hash algorithm remain open.

## Interaction grammar

The human-facing Skill is use-case-derived and judgment-bearing:

```text
/spectacular <judgment-verb> <noun-or-target>
```

Guided verbs are `orient`, `propose`, `define`, `decide`, `start`, `resume`,
`resolve`, `handoff`, `message`, `assess`, `reconcile`, and `audit`. Imagine,
grill, refine, review, verify, and generic next may remain internal techniques,
not required top-level modes.

The mechanical CLI is noun-first:

```text
spectacular <noun> <operation> [reference]
```

Canonical mechanical operations include list/show, stable-reference resolution,
schema/link validation, projection calculation, confirmed immutable-record
writes, Proposal-base checks, explicitly authorized transitions, and archival
after closure. Interpretation, Mission boundaries, evidence sufficiency, owner
acceptance, semantic reconciliation, and safe-continuation judgment remain in
the Skill/owner workflow. Native providers retain their own effects and verbs.

Primary help leads with guided jobs and one canonical mechanical path per
operation. Aliases appear only in a separate deprecated section and report an
exact canonical replacement; warnings stay separate from machine output.

## Orientation and projections

The default is a source-backed card plus a typed authority spine. Contextual
dashboard or narrative shapes may be requested.

A Mission Card shows a plain-language outcome lead, recorded state, current
Objective, latest relevant Run/Checkpoint, Fog, evidence/freshness pointers,
Decision/reconciliation pointers or absence, exactly one safe continuation or
owner gate, and source/generation basis.

It may state recorded lifecycle/disposition, explicit relationships, declared
freshness/absence, deterministic counts, latest attributable Run event, recorded
owner Decisions, and mechanically closed properties. It must link evidence and
authority before claiming implementation correctness, evidence sufficiency,
owner acceptance, Mission completion, Contract reconciliation, or next-action
safety.

The typed authority spine uses left-to-right arrows for provenance/authorized
progression and vertical branches for containment. Nodes display literal ID and
state, unknowns are explicit Gap nodes, and every node drills down. Color,
percentage, or icons never carry meaning alone. The visual is always a generated
projection.

## Collision policy

One concept has one canonical public name and operation path. Familiar words may
remain prose but cannot become competing framework labels. Aliases are exact
redirects without distinct semantics or primary documentation. Generic `next`,
`test`, `try`, and `iterate` are not permanent concepts. Objective replaces
milestone/slice/waypoint; Gap remains distinct from Fog; Evidence replaces
signal/telemetry; Message replaces letter/beacon, while broadcast is delivery.

Support windows, warning duration, removal releases, recovery promises, and
migration mechanics remain reserved for the compatibility checkpoint and S11.
