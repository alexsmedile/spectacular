---
type: refactor-foundation-contract
version: 1.0
status: accepted
decision_session: S05
accepted_by: owner
accepted_at: 2026-08-08
central_disposition: accept-with-clarification
upstream:
  - PRODUCT-CONSTITUTION.md@1.0
  - TRUTH-PROVENANCE-FLOOR.md@1.0
  - SUCCESS-EVIDENCE-CONSTITUTION.md@1.0
  - PRODUCT-TRUTH-CONTRACT-MODEL.md@1.0
  - WORK-UNIT-LIFECYCLE-CONTRACT.md@1.0
next_session: S06
---

# Execution Authority and Effects Contract

This accepted S05 contract allocates authority for bounded work, delegated
Autopilot, external effects, and lifecycle transitions. It does not define
evidence sufficiency, reconciliation mechanics, storage, public vocabulary,
or provider adapters.

## Accountability boundary

| Actor | Owns | Does not own |
|---|---|---|
| Owner | Product direction, target deltas, Mission boundaries, charters, reserved effects, final Mission disposition, and current Capability Contract changes | Provider mechanics or executor attempt state |
| Spectacular | Compiling and validating authority envelopes; exposing conflicts and Gaps; recording attributable intent, decisions, and receipts | Authority merely by recording, provider credentials, or provider mutation |
| Host runtime | Bounded execution mechanics | Product/lifecycle authority beyond the envelope |
| Native provider | Its own Git, Git-hosting, deployment, and other mutations and attestations | Mission interpretation or contract authority |
| Executor | Work and delegated reversible choices inside a current envelope | Ambient authority or final Mission/contract resolution |
| Independent reviewer | Attributable assurance and blocking findings | Replacing owner authority |

Every consequential decision has one accountable owner. Proposer, adviser,
reviewer, executor, recorder, and decider are separately attributable.
Security, privacy, rights, and cross-domain objections must reach the owner.

## Default Mission envelope

One owner approval may cover declared safe work: inspection, local edits,
checks, coherent commits, and explicitly named branch/worktree, push, and
draft-PR actions. The executor stops and proposes a delta on a material scope,
contract, authority, baseline, policy, risk, or provider change.

Expected paths are scope signals, not a brittle fence. A newly discovered
supporting file may be changed only when outcome, contracts, risk, permissions,
and effects remain unchanged.

## Autopilot delegation

Autopilot is an explicit, non-default delegation mode. Before it starts, an
agent must clarify unclear intent, risk, policy, or authority through focused
questions or concrete options.

An Autopilot charter identifies outcome and non-goals; delegated decision domain
and authoritative sources; allowed and forbidden providers/effects; budgets,
required checks, expiry, and stop conditions; and return destination.

There is no ambient agent authority. The current charter, baseline, policy,
provider permissions, and stop rules constrain every action. An Autopilot may
complete its **authorized execution** and return attributable decision and
evidence receipts; only the owner may resolve the Mission or authorize a
change to the current Capability Contract.

## Initial effect ceiling

When expressly named in a charter, Autopilot may perform local work,
constrained dependency additions, commits, pushes, draft PRs, and staging
deployments.

It initially excludes merge, production release, production/configuration/
secret changes, remote deletion, destructive data actions, and
security/privacy/rights-sensitive effects. Remote deletion always requires
separate explicit owner consent; local cleanup never implies it. Public or data
changes outside the accepted target require a proposed Mission delta and owner
review before execution.

## Stops, drift, and resumption

The following require a stop and proposed delta: a new dependency, unrelated
refactor, public or data change, policy impact, expanded outcome, irreversible
action, baseline conflict, failed required check, exhausted retry budget, or
unknown provider fact. Inspection, evidence collection, option preparation,
and unaffected recovery may continue.

An interrupted Run resumes only after revalidating its unchanged charter,
baseline, policies, provider facts, budgets, and stop conditions. Failure,
handoff, and resumption never broaden authority.

## Provider and lifecycle boundaries

Native providers perform and attest to their own mutations. Spectacular
interprets Mission relevance, compiles envelopes, and records pointers and
receipts; it does not become a provider or credential holder.

Anyone may propose a Proposal, Gap, or Decision. The owner accepts or rejects
product direction, activates and resolves Missions, and authorizes Capability
Contract transitions. Runtime/executor controls only Run start, return,
interruption, and failure. S06 defines proof sufficiency, closure, and
reconciliation mechanics.
