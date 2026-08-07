---
type: source-card
source: source-010
provided_as: abstraction-layers-and-orchestration
received: 2026-08-07
authority: unsourced-expanded-synthesis
status: ingested
scope: [abstraction, semantic-layer, modules, context, memory, orchestration, handoffs, review, sandboxing]
duplicate_sections: [source-008-substantial, source-009-substantial, source-007-selected]
completeness: substantial-with-missing-visual
---

# Source 010 — Abstraction layers and orchestration protocol

## Thesis

Treat abstraction boundaries as context-survival mechanisms: agents should operate on
semantic data concepts, cohesive code interfaces, role-bounded prompts, compressed memory,
typed state handoffs, independent review, and deterministic validation rather than loading
the entire implementation and conversation.

## Source integrity

The source contains two tightly related syntheses and substantially restates Sources 008
and 009. It provides no author, date, citations, benchmark, schema, or observed project
results. It refers to an interactive visualization, but no visualization payload was
included. Named “300-token,” AC/DC, and Gauntlet rules were not traceable as established
standards; they remain proposal vocabulary rather than evidence.

The generic Diamond is a familiar fan-out/fan-in topology. A concrete implementation still
needs a join barrier, state versioning, retries, checkpoints, and resource bounds; an
example demonstrating those concerns is available in the [trpc-agent-go graph package](https://pkg.go.dev/trpc.group/trpc-go/trpc-agent-go/examples/graph/diamond).

## Claim audit

| Claim | Assessment | Correction or boundary |
|---|---|---|
| Abstraction layers preserve agent context | strong design hypothesis | A boundary helps only if it exposes sufficient authority, failure behavior, and escape hatches; hiding relevant constraints creates confident errors. |
| A semantic layer lets agents ignore physical storage | useful for stable business queries | Translation semantics, freshness, authorization, cost, and unsupported queries remain visible concerns. This is an adapter, not an infallible substrate. |
| Agents naturally create shallow modules | plausible tendency, unsupported universal claim | Evaluate cohesion, interface complexity, coupling, and change locality rather than agent authorship or file count. |
| Messy code inside a grey box is acceptable | rejected | Encapsulation limits blast radius but does not remove security, performance, debugging, operability, or maintenance costs. Internal quality has minimum invariants. |
| Orchestrator knows only what; specialists know only how | unsafe over-isolation | The orchestrator needs enough domain and risk context to decompose safely; specialists need the local purpose, constraints, and integration contract. |
| JSON is the abstraction boundary | too narrow | Typed JSON is good for structured state; Markdown and primary-artifact references remain better for nuanced human review and large evidence. |
| Always read L3 first | useful default, unsafe absolute | Summary-first retrieval needs freshness, provenance, conflict signals, and direct escalation to primary evidence for consequential decisions. |
| Handoffs must stay under 300 tokens | arbitrary heuristic | Budget handoffs by task and measure completeness. Preserve exact paths, baselines, decisions, evidence, failures, and next action; omit conversational chronology. |
| Skeptic rejection loops until entirely satisfied | unsafe/non-terminating | Use an explicit acceptance rubric, severity threshold, attempt budget, unchanged-finding detection, and escalation packet. |
| Deterministic sandboxing should supplement review | strong | Linters, parsers, AST checks, tests, permissions, and sandbox boundaries verify different properties and should return exact evidence into bounded repair. |

## How this changes the Spectacular refactor lens

The source suggests an **abstraction audit**, not an automatic multi-agent runtime:

1. **Intent layer** — Vision, accepted contracts, and Missions express why and what.
2. **Orchestration layer** — the skill routes judgment and compiles bounded work context.
3. **Mechanism layer** — the CLI performs deterministic reads, writes, schemas, and gates.
4. **Truth layer** — Markdown records, source code, tests, and observations retain explicit authority.
5. **Adapter layer** — Git, GitHub, databases, agent runtimes, and companion skills own external mechanics.

For every current command, reference, agent, and collection, review:

- Which layer owns it?
- Does it leak details upward or duplicate authority downward?
- What is the smallest sufficient interface?
- How does a caller reach deeper evidence when the abstraction fails?
- Which invariant is checked mechanically, and which requires judgment?

## Valuable ideas preserved

- Context budgets belong to abstraction boundaries, not just whole prompts.
- Semantic access should be stable over physical storage when the business query is stable.
- Runtime modules should hide cohesive complexity behind narrow interfaces.
- Delegated internals remain subject to quality, evidence, and observability invariants.
- Orchestrator and specialist contexts differ, but both require goal and safety alignment.
- Handoffs should be capsules, not transcripts, with a measured completeness budget.
- Review-repair cycles need independent evidence plus hard termination and escalation rules.
- Deterministic tools and sandboxes remain stronger than an agent's confidence.

## Provisional assessment

**Strong:** layer-specific context contracts; semantic adapters; deep cohesive modules;
typed state handoffs; transcript-free capsules; deterministic validation and sandboxing;
bounded independent review.

**Promising:** summary-first hierarchical memory and hierarchical agent roles when tested
against omission and coordination costs.

**Rejected as written:** messy internals are acceptable; orchestrators need no technical
context; specialists need no broader purpose; JSON is the only boundary; 300 tokens is a
universal limit; L3 is always sufficient first; the Gauntlet runs until subjective satisfaction.

No abstraction stack, agent hierarchy, memory hierarchy, handoff size, or review topology is
accepted as product architecture here.

## New concept pieces

- [PZL-123 — Layer-specific context contract](concepts/PZL-123-layer-context-contract.md)
- [PZL-124 — Semantic query adapter](concepts/PZL-124-semantic-query-adapter.md)
- [PZL-125 — Grey-box delegation with internal invariants](concepts/PZL-125-grey-box-invariants.md)
- [PZL-126 — Budgeted handoff capsule](concepts/PZL-126-budgeted-handoff-capsule.md)
- [PZL-127 — Bounded review-repair cycle](concepts/PZL-127-bounded-review-repair.md)

## Decision packets seeded

- Approve an abstraction-layer audit as a review method for the current surface.
- Define the minimum context and escape hatch at each Spectacular layer boundary.
- Decide whether a semantic adapter is a general product capability or only a pattern for
  future database-backed companions.
- Define internal quality invariants for any agent-owned grey box.
- Replace the 300-token slogan with a measurable handoff completeness budget.
- Define review termination, severity, retry, and escalation before any Gauntlet-style loop.
