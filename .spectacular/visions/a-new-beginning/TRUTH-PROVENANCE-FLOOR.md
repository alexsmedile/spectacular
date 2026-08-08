---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S03A
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept
upstream: PRODUCT-CONSTITUTION.md@1.0
next_session: S02
---

# Minimum Truth and Provenance Floor

This is the accepted S03A prerequisite for the `a-new-beginning` rebuild. It
defines the minimum conditions under which Spectacular may treat information as
usable evidence. It does not define the final contract graph, workspace
ontology, artifact taxonomy, storage architecture, or command surface.

## 1. Authority is claim-scoped

No source is universally authoritative. Authority depends on the question a
claim answers.

| Question | Provisional authority |
|---|---|
| What behavior or direction was accepted? | The latest accepted, in-scope owner decision or contract |
| What does a repository revision implement? | The cited revision, source, schemas, and fresh test results |
| What occurred in an environment? | An attributable observation scoped to environment and time |
| Why was a past choice made? | The attributed historical record |
| What should happen next? | An authorized decision; recommendations and assumptions are inputs only |

Direct discovery, inference, assumption, recommendation, accepted direction,
implementation fact, observation, and historical record remain distinguishable.
One authority cannot silently answer a different authority domain's question.

## 2. Conflict is a stop-and-expose state

When relevant authorities disagree, Spectacular states the mismatch instead of
selecting a convenient winner. For example:

> Accepted behavior is X; revision Y implements Z; environment E exhibited W.

Conflict, contrary evidence, mismatch, missing provenance, material staleness,
or unverifiability stops consequential conclusions, acceptance, lifecycle
promotion, and external effects until the appropriate authority reconciles it.
Coherence across the relevant, sufficiently fresh authorities is the condition
to proceed.

## 3. Provenance is lean and pointer-first

Evidence remains where it naturally belongs. A receipt points to it instead of
copying it into a second source of truth.

A consequential claim records only what is needed to locate, interpret, and
re-check it:

- the claim and question it answers;
- its authority class;
- a durable source pointer;
- relevant scope, revision, contract version, environment, or time;
- whether it is direct or derived;
- actor, observer, or method when material to trust.

A locator without the scope needed to re-check the claim is insufficient. This
floor does not choose Markdown, a database, a graph store, or any hybrid.

## 4. Freshness is claim-scoped

There is no universal expiry period. Historical records remain valid accounts
of their attributed history. Claims about current state require a freshness
basis appropriate to their authority class.

Revalidation is required when:

- the claim's relevant scope, revision, contract, authority, or environment changes;
- contrary evidence or a material mismatch appears;
- the source can no longer be located or interpreted;
- the claim lacks the freshness basis needed for its intended consequence.

Drift is an alert and a check. It never silently renews, supersedes, or replaces
authority.

## 5. Projections never become authority

Indexes, dashboards, status cards, diagrams, decision trees, generated
summaries, retrieval results, handoffs, and other compressed views are
non-authoritative projections.

Live generation may improve freshness; it does not create authority. A useful
projection:

- identifies its source set and generation time when material;
- drills down to authoritative inputs;
- exposes missing inputs, conflicts, and freshness state;
- does not resolve mismatches or promote decisions;
- labels stored views as dated snapshots.

Logs are attributed historical records of events. They do not become current
governing truth merely because they are chronological or generated live.

## 6. Unknowns are explicit, with bounded continuation

Missing, stale, contradictory, or unverifiable information is reported as
unknown with its affected question and consequence. The agent does not fill the
gap with an unlabeled inference.

While a consequential path is stopped, an agent may continue work demonstrably
unaffected by the gap, including bounded inspection, evidence collection,
option preparation, and recovery planning. Any assumption used as a decision
input is explicit, scoped, and accepted by the appropriate owner.

## What this contract does not decide

The following remain deliberately open:

- the complete truth hierarchy and capability/system contract model (S03B);
- evidence sufficiency, reconciliation envelopes, and closure mechanics (S06);
- execution, approval, and provider authority (S05);
- authoritative homes, workspace scaffold, databases, graphs, and projections (S08);
- public vocabulary, artifact names, and command grammar (S09).

The parked `live-decision-graph` idea may be evaluated in S03B/S08 only as a
derived, non-authoritative projection.

## Exit condition

S02 may now define success measures and evidence rules without circularly
deciding what counts as authoritative, attributable, fresh, or unknown.
