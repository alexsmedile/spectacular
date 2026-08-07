---
type: source-card
source: source-007
provided_as: source7
received: 2026-08-07
authority: synthesized-proposal
status: ingested
scope: [system-graph, contracts, missions, run-compiler, gates, evidence, integration, taxonomy]
supporting_materials: [ai-driven-sdlc, senior-agent-harness, vibe-to-spec]
duplicate_sections: [source-002-selected, source-006-expanded]
completeness: substantial
---

# Source 007 — System contract graph

## Thesis

Add a first-class typed system contract graph to the capability/Mission/run loop:
the graph describes accepted current truth, a Mission proposes a versioned graph
delta, one or more runs attempt it, evidence passes independent gates, and closure
reconciles the graph plus durable decisions and learnings.

## Source integrity

The source expands Source 006 rather than independently confirming it. Its three
local attachments were read in full and remain supporting materials, not additional
source-count entries. They have no visible author, date, citations, or empirical
results, so they support design ingredients rather than architectural conclusions.

| Supporting material | What it supports | What it does not establish |
|---|---|---|
| AI-driven SDLC | AI across planning/code/review/test/ops, human oversight, repo context, draft PRs | typed graph, graph reconciliation, dual locks, contract authority |
| Senior agent harness | system map, code boundaries, exact versions, secret rules, reconnaissance, minimal diffs, tests, bounded retry | durable graph model or superiority over existing workflows |
| Vibe-to-spec workflow | plain-language product framing, agent technical blueprint, human business review/lock, staged execution | independent engineering approval, multiple typed contracts, graph-first retrieval |

The attachments often recommend prompt instructions; Source 007 deliberately
converts those into structural gates. That conversion is an original proposal,
not evidence supplied by the attachments.

Attachment provenance:

- AI-driven SDLC — `/Users/alex/.codex/attachments/6a5d26df-92e7-4029-8599-b6de93238f7c/pasted-text.txt`
- Senior agent harness — `/Users/alex/.codex/attachments/ed5bc8a6-6003-431e-848c-96f5c382fb4e/pasted-text.txt`
- Vibe-to-spec workflow — `/Users/alex/.codex/attachments/caf4ca8a-3189-456c-bd96-5e31644c5a89/pasted-text.txt`

## Proposed graph model

The graph separates two axes:

- capabilities own observable end-to-end behavior;
- components, interfaces, state, operational contracts, and policies own system
  topology and constraints.

Typed relationships connect realization, invocation, state access, policy, and
verification. Every contract shares guarantees, exclusions, inputs/outputs,
invariants, failure modes, dependencies, implementation pointers, probes,
freshness, and provenance. Small projects may embed types as sections until
promotion is earned.

## Proposed design and approval pipeline

1. Human supplies outcomes, experience, business/privacy rules, constraints,
   non-goals, and examples without pretending to technical expertise.
2. Architect agent discovers repository facts and separately proposes graph deltas,
   alternatives, trade-offs, security, failure, and operational consequences.
3. Human product lock confirms behavior and business logic.
4. Independent agent performs engineering assurance before Mission approval.

The second gate's authority is underspecified: the same source says an AI reviewer's
confidence is not evidence. Engineering assurance can be a required report without
granting an agent final approval authority; that distinction needs a decision.

## Proposed Mission and run model

- A Mission is a versioned transaction over named contracts and declares a
  completion boundary, autonomy envelope, objective DAG, risks, assumptions, and proof.
- A compiler retrieves relevant graph neighbors and source-of-truth facts into a
  bounded run manifest; semantic search may locate candidates but cannot decide authority.
- Structural gates enforce reconnaissance, scope, evidence strategy, diagnostic
  retries, risk-triggered review, and completion-boundary evidence.
- A Mission owns several attempt records so failure does not corrupt durable intent.

## Current-system comparison

No current first-class system graph or proposed relationship vocabulary was found.
Existing capability specs, flat schema/axis contracts, architecture/stack docs,
policy, dependencies, request provenance, AFK runs, verification logs, and spec-sync
cover parts of the information. The proposal must show that explicit edges improve
retrieval and reconciliation enough to justify graph maintenance.

The equation `System₁ = reconcile(System₀, MissionDelta, Implementation, Evidence)`
is useful as an invariant but incomplete as authority language: code/tests remain
implementation truth, production may drift independently, and a Markdown graph may
be stale until reconciliation. “Current truth” should mean accepted contract state,
not infallible runtime reality.

## Architectural collisions

1. Source 006 says one coding agent and one run record; Source 007 proposes multiple
   specialized runs and agents within one Mission.
2. Source 006 proposes one mandatory human approval; Source 007 adds independent
   engineering assurance.
3. Source 006 defers separate registries; Source 007 makes six stable contract types
   and a graph first-class, though not necessarily separate files.
4. Source 006's minimal tree is replaced by a larger `system/` graph plus decisions,
   learnings, deltas, and run collections.
5. Source 002 keeps Missions through Operations; Source 006 closes around proof;
   Source 007 makes completion boundary configurable from local to observed production.
6. Mission/run lifecycles simplify some taxonomies while adding two explicit state machines.

## Valuable ideas independent of graph adoption

- Never conflate end-to-end capabilities with implementation components.
- Separate discovered facts from invented recommendations.
- Treat missing information as a named assumption/question/discovery objective.
- Retrieve along authoritative relationships before broad semantic similarity.
- Define evidence before implementation and select proof strategy by change class.
- Permit retries only when a new hypothesis or evidence is introduced.
- Make completion boundaries explicit and verification risk-triggered.
- Keep durable intent coherent across failed attempts and resume from evidence plus Git.

## Provisional assessment

**Strong:** orthogonal capability/architecture axes; fact/recommendation provenance;
named unknowns; graph-aware bounded context; structural reconnaissance/scope/evidence
gates; hypothesis-driven retries; completion boundaries; Mission/run separation.

**Needs a thin-slice proof:** graph maintenance, six contract types, graph-first
retrieval, compiler, freshness, reconciliation, and the larger workspace.

**Disputed authority:** independent-agent “approval,” multi-agent run chain, and
calling the graph current truth without retaining code/runtime authority distinctions.

No graph schema, new approval authority, lifecycle, or workspace is accepted here.

## New concept pieces

- [PZL-080 — System contract graph](concepts/PZL-080-system-contract-graph.md)
- [PZL-081 — Orthogonal capability and architecture axes](concepts/PZL-081-orthogonal-modeling-axes.md)
- [PZL-082 — Stable typed system contracts](concepts/PZL-082-typed-system-contracts.md)
- [PZL-083 — Shared contract envelope](concepts/PZL-083-shared-contract-envelope.md)
- [PZL-084 — Named missing-contract escalation](concepts/PZL-084-named-missing-contracts.md)
- [PZL-085 — Vibe-to-system authority split](concepts/PZL-085-vibe-to-system-authority.md)
- [PZL-086 — Separate facts from recommendations](concepts/PZL-086-facts-versus-recommendations.md)
- [PZL-087 — Product lock and engineering assurance](concepts/PZL-087-dual-review-lock.md)
- [PZL-088 — Mission as versioned graph transaction](concepts/PZL-088-mission-graph-transaction.md)
- [PZL-089 — Objective DAG with proof prerequisites](concepts/PZL-089-objective-proof-dag.md)
- [PZL-090 — Authoritative graph-first retrieval](concepts/PZL-090-graph-first-retrieval.md)
- [PZL-091 — Mission compiler and run manifest](concepts/PZL-091-mission-compiler.md)
- [PZL-092 — Structural gates over prompt advice](concepts/PZL-092-structural-gates.md)
- [PZL-093 — Reconnaissance gate](concepts/PZL-093-reconnaissance-gate.md)
- [PZL-094 — Contract and repository scope gate](concepts/PZL-094-scope-drift-gate.md)
- [PZL-095 — Evidence-first strategy by change class](concepts/PZL-095-evidence-first-strategy.md)
- [PZL-096 — Hypothesis-driven retry budget](concepts/PZL-096-hypothesis-retry-loop.md)
- [PZL-097 — Explicit completion boundary](concepts/PZL-097-completion-boundary.md)
- [PZL-098 — Risk-triggered staged run chain](concepts/PZL-098-risk-triggered-run-chain.md)
- [PZL-099 — Mission/run attempt separation](concepts/PZL-099-mission-run-separation.md)
- [PZL-100 — Mission-owned local semantics](concepts/PZL-100-mission-owned-semantics.md)
- [PZL-101 — Graph-oriented workspace candidate](concepts/PZL-101-graph-workspace.md)
- [PZL-102 — Thin vertical graph slice](concepts/PZL-102-thin-vertical-graph-slice.md)

## Decision packets seeded

- Is the graph accepted contract state, runtime truth, or a retrieval projection?
- Which contract types and edges earn first-class identity in the first slice?
- Is engineering assurance advisory evidence, a blocking gate, or approval authority?
- Does one Mission own one run or a chain of typed attempts?
- Which completion boundaries belong to Spectacular versus release/deployment systems?
