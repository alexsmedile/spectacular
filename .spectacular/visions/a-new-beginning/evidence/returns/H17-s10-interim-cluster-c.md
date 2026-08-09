---
type: handoff-checkpoint
schema_version: spectacular.handoff-checkpoint.v1
handoff_id: H17
session: S10
phase: after-cluster-C-before-cluster-D
status: active-incomplete
central_disposition: filed-not-accepted
baseline_commit: 6a0c6c3d3b165d4216c209a5fb8052fbdb4604c6
baseline_tree: 111e66734856574e4040258c9abaf824b94f1c7d
baseline_dirty: false
supersedes_for-current-state: H17-s10-interim-clusters-a-b.md
date: 2026-08-09
---

# H17 delta checkpoint — S10 Cluster C

## Authority boundary

This checkpoint extends the previously filed A/B checkpoint with explicit owner decisions from
Cluster C. It is not the final H17 return and does not accept S10 or authorize S11, lifecycle
changes, implementation, migration, deletion, or archival. Cluster D remains unresolved.

The H17 baseline remains commit
`6a0c6c3d3b165d4216c209a5fb8052fbdb4604c6`, tree
`111e66734856574e4040258c9abaf824b94f1c7d`. All twelve immutable inputs remain verified and H17
performed no repository mutations.

## C1 — Autopilot and execution topology

**Owner disposition: adopted.**

- Autopilot is Mission consent, not an executor. Spectacular is the Mission Director/control
  plane; host runtimes execute.
- Small jobs remain inline. Delegation must be earned.
- Execution topology is prepared before activation. Objectives may expose candidate lanes, but
  dependency, independence, cancellation, evidence, and join analysis establish whether work may
  run concurrently.
- Missions never nest. An Objective may have an optional Lead; Investigator and Builder are
  operator roles; Reviewer/Verifier is independent when required.
- Nested delegation is disabled unless depth and budget are explicitly authorized.
- The owner has one human-facing channel.
- A mutating Mission receives a branch when needed; workers do not receive branches merely because
  they are workers.
- Git and provider effects remain bounded by Mission consent.
- AFK, workspace preflight, and traffic responsibilities merge into Autopilot, start/resume,
  closure, and live orchestration. The GitHub bridge becomes a provider adapter.

## C2 — Project Guardrails

**Owner disposition: adopted.**

The fixed POLICY hook/schema subsystem is replaced by human-written Markdown **Project
Guardrails**, mechanically retrieved and interpreted by the agent:

- A Markdown heading identifies a policy; no explicit policy ID is required.
- The first non-empty line under the heading is an order-independent selector line containing one
  or more defined `@Event` tokens and optional extensible `$domain.verb` command tokens, such as
  `@BeforeAction $git.commit`.
- Remaining Markdown is unrestricted and injected verbatim.
- Events are defined; command vocabulary is guided and extensible.
- Editable Guardrails cannot create authority, weaken accepted invariants, prove provider facts,
  or impose unrelated harnesses.
- Native hooks are optional and are never edited automatically.
- Exact parser, path, catalog, and adapter design remains S11 work.

## C3 — Integrity, repair, and migration

**Owner disposition: adopted as recommended.**

| Current responsibility | Disposition | v2 boundary |
|---|---|---|
| Integrity checking | keep-simplify | Scoped, read-only v2 integrity checks support trustworthy prepare, resume, execute, and resolve operations. |
| Universal doctor repair | retire | Deterministic corrections belong to their owning operations. Judgment repair routes through Gap, Proposal, Decision, or Mission. Code repair belongs to Mission execution; Git/provider repair belongs to those providers. |
| v1→v2 conversion | extract | The S11 migration capsule owns conversion. Clean v2 core carries no legacy parser, fallback read, dual write, lazy conversion, historical archive tree, or general legacy migration registry. |

Retirement remains blocked until recovery is proven through the immutable v1 release/tag and a
validated snapshot. The current broad doctor—approximately 3,000 CLI lines and 23 areas—and
scattered migration machinery are direct maintenance evidence, not proof of external value.

## Open S10 work

Cluster D remains: generalized document engine, templates/rules, verification/review, current
fleet, Mission Slice Advisor, and Design Sufficiency Reviewer. H17 must continue one explicit
owner-decision cluster at a time and return the complete packet only after every remaining choice
is disposed.

## Next action

Continue H17 interactively with Cluster D. Central orchestration must wait for the complete H17
return before applying `accept | bounce | escalate` to S10 or authorizing S11.
