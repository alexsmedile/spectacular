---
type: synthesis-checkpoint
checkpoint: 018
status: current
date: 2026-08-08
authority: central-orchestration
accepted_contracts: [S01, S03A, S02, S03B, S04]
next_session: S05
---

# Synthesis 018 — S04 work-unit model accepted

## Central disposition

H09/S04 is accepted with one upstream-preserving normalization: Proposal
acceptance authorizes a target delta but does not make it the current Capability
Contract before authorized reconciliation.

## Accepted spine

```text
current Capability Contract
  + accepted target delta
  → Mission
  → durable Objectives
  → boundary-based Runs
  → run-local Tasks
```

Missions never nest. Goals only link them. Handoffs remain intra-Mission
dispatch records; typed links carry cross-Mission relationships. Gaps collapse
uncertainty into one typed concept, and projections remain non-authoritative.

The authoritative result is
[`../WORK-UNIT-LIFECYCLE-CONTRACT.md`](../WORK-UNIT-LIFECYCLE-CONTRACT.md). Its
hash sidecar binds the accepted version.

## Next gate

S05 is next-ready to decide execution authority, human gates, lifecycle
transition authority, and external side effects. It must not reopen S04's
semantic ownership or smuggle storage/public-command choices into authority.
No S05 result or implementation is authorized by this checkpoint.
