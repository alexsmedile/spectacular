---
type: synthesis-checkpoint
checkpoint: 003
sources: [source-001, source-002, source-003, source-004, source-005, source-006, source-007]
concepts: 102
human_dispositions: 0
updated: 2026-08-07
---

# Synthesis checkpoint 003

## What Source 007 changes

Source 007 does not introduce a fourth product world. It adds a potential system-
modeling and context-compilation layer to the guarded Mission runner (World B) or
host-owned execution hybrid (World C) from checkpoint 002.

Its central equation is useful when read narrowly:

> Accepted contract state after closure is reconciled from the prior accepted
> state, approved delta, actual implementation, and adequate evidence.

It is unsafe when “system graph” is treated as infallible current truth. Code and
tests remain implementation truth; production observations may reveal drift or
defects; the graph is the accepted contract model and may itself become stale.

## Strong additions independent of a graph platform

1. Capabilities and implementation architecture are orthogonal axes.
2. Repository facts, human requirements, assumptions, and agent recommendations
   require distinct provenance.
3. Missing contract information must become an explicit unknown rather than a silent inference.
4. Outcome objectives may declare dependencies and proof without predicting files.
5. Evidence strategy should be selected before implementation by change class.
6. Retries must add a hypothesis or evidence and end with a resumable failure packet.
7. Delivery completion must be declared explicitly rather than inferred from local tests.
8. Durable Mission intent should survive failed execution attempts.
9. Low-value AI review commentary should be suppressed; assurance is risk-triggered.

These can improve any of Worlds A, B, or C without adopting six contract registries,
a graph compiler, or a new workspace.

## Graph adoption levels

### Level 0 — Modeling rule

Do not conflate capabilities with components/interfaces/state. Record provenance
and named unknowns in existing specs and request artifacts.

### Level 1 — Linked contract projection

Add a small set of typed relationships and freshness/provenance fields to existing
contracts. Compile context by following those links, with source files and decisions
retaining authority.

### Level 2 — First-class system graph

Introduce graph versions, graph transactions, six semantic contract types, a Mission
compiler, reconciliation rules, and impact traversal.

### Level 3 — Graph-oriented product architecture

Replace the workspace with `system/`, typed subcollections, Mission DELTAs, multiple
runs, decisions, learnings, and archive; drive staged agent execution from it.

Only Level 0 currently has strong conceptual support. Levels 1–3 require the thin
vertical slice in PZL-102 and must earn each additional authority and collection.

## Direct contradictions now exposed

| Question | Source 006 | Source 007 |
|---|---|---|
| Approval | one mandatory human lock | product lock plus independent engineering assurance |
| Execution | one coding agent, one active run | multiple specialized runs and potentially different agents |
| Workspace | PROJECT/SYSTEM/capabilities/Missions/archive | VISION/system graph/Missions/decisions/learnings/archive |
| Taxonomy | one compact Mission lifecycle | Mission lifecycle plus separate run lifecycle and six contract types |
| Completion | tested draft PR in MVP | configurable local through deployed-and-observed |

The attachments support context, reconnaissance, bounded retries, testing, human
review, and draft PRs. They do not resolve these contradictions or validate the graph.

## Updated decision order after intake

1. Fix independent correctness/safety defects PZL-047 and PZL-048 if authorized.
2. Decide whether Spectacular owns execution, persists host-owned execution intent,
   or remains a lifecycle control layer.
3. Define accepted contract truth and its relationship to code/tests/production.
4. Choose graph adoption Level 0, 1, 2, or 3—subject to a Level-1 thin slice before
   approving higher levels.
5. Decide approval authority, Mission/run cardinality, and completion-boundary ownership.
6. Only then select workspace paths, CLI grammar, migration, and implementation architecture.

No decision in this checkpoint has been accepted.
