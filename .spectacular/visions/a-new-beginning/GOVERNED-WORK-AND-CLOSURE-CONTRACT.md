---
type: governed-work-and-closure-contract
contract: governed-work-and-closure
version: 1.0
status: accepted
accepted_by: owner
accepted_at: 2026-08-10
central_disposition: accept
scenario_a_join: compatible
next_action: prepare-b-c-governed-loop-mission
---

# Governed Work and Closure Contract

## Authority chain

The Capability Contract is the complete current owner-accepted behavior record. A Proposal is an
immutable base-bound change record, not the Contract moving through phases.

```text
current Contract + submitted Proposal delta
  → owner-accepted target delta
  → bounded Mission delivery
  → Evidence and owner Decision
  → authorized reconciliation
  → next complete Contract version
```

Draft Proposals may begin conversationally. Submission requires a target Contract, base version and
fingerprint, exact additions/modifications/removals, rationale, scope, and explicit Gaps. A generated
candidate Contract is disposable and non-authoritative. A new-capability Proposal supplies one
complete MCC-conformant candidate against explicit absence; reconciliation creates Contract v1.

## Governed creation and authority

Mission creation requires an accepted Proposal, matching base, bounded outcome, Objectives,
preparation verdicts, dependencies/Gaps, claim-mapped evidence plan, authority envelope, recovery
point, and activation Decision. Creation is atomic and identity-idempotent; stale base, collision,
or non-identical replay refuses without mutation.

The owner controls Proposal disposition, Mission activation/resolution, and current Contract
reconciliation. An approved Mission envelope may carry that authority across predeclared Objective
and readiness transitions; it does not require a fresh owner interruption at every Objective unless
the envelope says so. Run-local and reversible repair choices may be delegated.

Every mutation requires target identity, expected fingerprint, authorization Decision, and an
idempotency key. Conflict, drift, expiry, or failed prerequisites refuse with zero mutation.

## Handoff

A Handoff is an immutable intra-Mission execution dispatch. It binds Mission, Objective, Run,
actors/destination, scope, authoritative inputs, authority Decision, expected fingerprints,
allowed/forbidden effects, evidence, budgets, expiry, stops, recovery, and return destination.
Validation proves structure, references, same-Mission containment, authority-envelope subset, and
current baseline. It cannot prove actor competence, provider permission/effects, evidence truth or
sufficiency, or Mission success. Mission accountability never transfers.

## Evidence and assessment

Every material claim maps to attributable `direct | observation | proxy | judgment | unknown`
evidence with scope, method, actor, revision/Contract target, environment, time/freshness basis,
limits, contrary evidence, checks, and review state. Conflicts remain visible and block affected
claims rather than being averaged away.

Assessment returns `ready-for-owner | repair-required | escalated`. Each repair attempt records a
new hypothesis/evidence/narrower action, budget use, before/after evidence, checks, result, and
recovery point. Risk-triggered independent review follows the accepted evidence contract.

## Decision, reconciliation, and closure

A Decision records authority, question, scope, disposition, rationale, material alternatives,
targets/fingerprints, authorized effects, conditions/expiry, evidence, and supersession.
Reconciliation consumes the exact current Contract set, accepted Proposal delta, owner Decision,
and expected fingerprints as one logical all-or-nothing operation. Prior Contracts, Proposal,
Decision, Evidence, and reconciliation receipt remain inspectable.

`assessed`, `resolved`, `reconciled`, and `archived` are distinct. Resolution is
`completed | abandoned | superseded`; a zero-delta outcome requires explicit rationale.
Reconciliation either applies authorized changes or records `no-contract-delta` /
`reconciliation-not-required`. Archival follows resolution and reconciliation and proves nothing.

Scenario C must prove return → assessment → Decision → reconciliation → resolution → archival →
runtime replacement → cold resume without chat or provider-session history.

## Mechanical join

Scenario A's versioned result/refusal envelopes, exact references, source/fingerprint basis,
registry-owned mechanics, separated warnings, and exit classes `0/2/3` satisfy this contract without
redesign. B+C may add mutation, authorization, idempotency, and transaction fields without changing
Scenario A's read contract. Exact additive record encoding remains a reversible Mission-level detail.
