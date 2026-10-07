---
type: Atlas
title: Operating levels and software factory architecture
version: '0.3'
created: '2026-08-31T17:59:30+02:00'
updated: '2026-10-07T18:49:34Z'
---

# Working levels and modular delivery

Choose a file boundary that matches the meaning and lifetime of the work.
The operating-level model is a way to discuss scale, not a required process.

| Scale | Examples | Useful durable context |
|---|---|---|
| Code primitive | Function, type, invariant | Focused code-linked explanation when independently useful |
| Module | Package, parser, service boundary | Linked Atlas object or interface Spec |
| Product slice | Feature, release phase, delivered tranche | One coherent plan and outcome |
| Repository | Accepted project direction and constraints | Anchors and Decisions |
| Coordinated effort | Several independent slices | Optional linking plan or campaign |

## Split by meaning

A parent plan may link a parser slice, installer slice, and delivery slice. Each
child is independently actionable and has its own verification. Do not create
one file per implementation step when those steps need the same readers and lifetime.

A reusable interface spec can outlive a release plan. A domain entity page can
outlive several specs. A Contract preserves accepted mechanical agreement across
Missions. Links make those different lifetimes navigable without copying text.

## Choose execution deliberately

Ordinary work uses the task’s existing authorization and repository rules. Native
plan mode owns its controls; the Spectacular Skill stores the project plan in
plans/ when writing is allowed. Risk can justify suggesting a governed route;
it does not select one automatically.

Coordination does not change Git policy. Main-branch commits, worktrees, reviews,
and publication follow the owner’s authorization and applicable repository rules.
A worker’s returned diff is checked after integration before claiming completion.

[Project truth](../PROJECT.md) · [Governance choice](governance-tiers-and-cadence.md) ·
[Workspace navigation](workspace-navigation.md)
